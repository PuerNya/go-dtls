package dtls13

import (
	"crypto"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"slices"
)

// clientHandshakeState carries the values that cross phase boundaries in the
// client handshake. Each field is set by the step named in its comment and read
// by later steps; nothing here outlives the handshake. The struct deliberately
// has no back-pointer to *Conn: steps are *Conn methods that receive the state
// explicitly, matching postHandshakeAuthState and returnRoutabilityState.
type clientHandshakeState struct {
	handshakeMetadata []byte
	finalized         bool
	// transcriptDigest is scratch space for every transcript sum.
	transcriptDigest [maxSupportedHashSize]byte

	// clientPrepareHello
	key          clientKeyExchange
	hello        *clientHello
	ech          *echClientContext
	greaseECH    bool
	session      *ClientSessionState
	sessionSuite *cipherSuite
	pskOffers    []clientPSKOffer

	// clientMarshalHello
	helloBody []byte

	// clientSendHello
	out *flight

	// clientReceiveServerHello / clientHandleHelloRetryRequest
	inbox                *handshakeInbox
	messages             completedHandshakeBatch
	serverHelloBody      []byte
	serverHelloSequence  uint16
	clientFinishedSeq    uint16
	serverHandshakeStart uint16
	helloRetrySuite      uint16
	transcript           *transcriptHash

	// clientConfirmECH
	echAccepted bool
	sh          *serverHello
	suite       *cipherSuite

	// clientValidateServerHello
	usingPSK      bool
	resumed       bool
	externalPSK   *externalPSKSelection
	selectedOffer *clientPSKOffer
	// singleTicketOffer backs selectedOffer when the client offered a single
	// legacy ticket rather than a pskOffers list.
	singleTicketOffer clientPSKOffer

	// clientDeriveHandshakeKeys
	schedule      *keySchedule
	receiveCipher *recordCipher
	sendCipher    *recordCipher

	// clientProcessServerFlight
	peerCerts                     []*x509.Certificate
	cachedInfo                    *clientCachedInformation
	peerRawPublicKey              []byte
	peerDelegatedCredential       []byte
	chains                        [][]*x509.Certificate
	negotiated                    string
	certificateRequest            *certificateRequestMessage
	certReqCompression            *certificateCompressionAlgorithms
	ocspResponse                  []byte
	clientCertificateTypeSelected bool
	serverFinishedSequence        uint16
}

// clientHandshake runs the DTLS 1.3 client handshake as a sequence of steps.
// Every step returns an error exactly as the monolithic implementation did, so
// the deferred alert sender in runHandshake sees the same error at the same
// point; in particular c.sendCipher is installed inside
// clientDeriveHandshakeKeys, which is where it was before the split.
func (c *Conn) clientHandshake() error {
	s := &clientHandshakeState{
		serverHelloSequence:  0,
		clientFinishedSeq:    1,
		serverHandshakeStart: 1,
	}
	steps := []func(*clientHandshakeState) error{
		c.clientPrepareHello,
		c.clientMarshalHello,
		c.clientSendHello,
		c.clientSendEarlyData,
		c.clientReceiveServerHello,
		c.clientHandleHelloRetryRequest,
		c.clientConfirmECH,
		c.clientValidateServerHello,
		c.clientDeriveHandshakeKeys,
		c.clientProcessServerFlight,
		c.clientDeriveApplicationSecrets,
		c.clientSendAuthAndFinished,
		c.clientAwaitFinalACK,
		c.clientFinalize,
	}
	for _, step := range steps {
		if err := step(s); err != nil {
			return err
		}
	}
	return nil
}

// clientPrepareHello generates the ephemeral key and key shares, builds the
// ClientHello from the configuration, selects the client session and PSK
// offers.
func (c *Conn) clientPrepareHello(s *clientHandshakeState) error {
	if c.config.HandshakeMetadata != nil && c.config.EncryptedClientHelloConfigList == nil {
		return &ConfigError{"HandshakeMetadata requires EncryptedClientHelloConfigList"}
	}
	key, err := generateEphemeralKey(c.config.CurvePreferences[0], c.config.Rand)
	if err != nil {
		return alertError(alertInternalError, err)
	}
	s.key = key
	keyShares := []keyShareEntry{{group: key.groupID(), data: key.publicBytes()}}
	if group, public, ok := key.fallbackPublicBytes(); ok && slices.Contains(c.config.CurvePreferences, group) {
		keyShares = append(keyShares, keyShareEntry{group: group, data: public})
	}
	hello := &clientHello{cipherSuites: append([]uint16(nil), c.config.CipherSuites...), keyShares: keyShares, supportedGroups: c.config.CurvePreferences, signatureSchemes: defaultSignatureSchemes(), serverName: c.config.ServerName, alpn: c.config.NextProtos, postHandshakeAuth: c.config.PostHandshakeAuth, recordSizeLimit: c.config.RecordSizeLimit, hasRecordSizeLimit: true, statusRequest: c.config.EnableOCSPStapling}
	hello.pskDHE = !c.config.SessionTicketsDisabled && c.config.ClientSessionCache != nil
	hello.handshakeMetadata = c.config.HandshakeMetadata
	if c.config.EnableDelegatedCredentials {
		hello.delegatedCredentialSchemes = delegatedCredentialSchemes()
	}
	if err := c.offerCertificateTypes(hello); err != nil {
		return err
	}
	c.offerCachedInformation(s, hello)
	if err = hello.setCertificateAuthorities(c.config.ServerCertificateAuthorities); err != nil {
		return err
	}
	if c.config.SessionTicketRequest.Enabled {
		hello.ticketRequest = c.config.SessionTicketRequest
	}
	if c.config.EnableCertificateCompression {
		hello.certificateCompressionOffered = true
	}
	if connectionID, offered := c.config.clientConnectionIDOffer(); offered {
		hello.hasConnectionID = true
		hello.connectionID = connectionID
		hello.returnRoutability = !c.config.DisableReturnRoutabilityCheck
	}
	if _, err = io.ReadFull(c.config.Rand, hello.random[:]); err != nil {
		return err
	}
	if c.config.EnableGREASE {
		hello.grease = true
	}
	if c.config.EncryptedClientHelloConfigList != nil {
		s.ech, err = newECHClientContext(c.config.EncryptedClientHelloConfigList)
		if err != nil {
			return err
		}
		hello.setEncryptedClientHello([]byte{echInnerType})
	}
	s.greaseECH = s.ech == nil && c.config.EncryptedClientHelloGrease
	s.session, s.sessionSuite = usableClientSession(c.config, c.conn)
	var externalOffers []clientPSKOffer
	if len(c.config.ExternalPSKs) > 0 {
		externalOffers = configuredExternalPSKOffers(c.config)
	}
	if len(externalOffers) > 0 {
		if s.session != nil {
			s.pskOffers = append(s.pskOffers, clientPSKOffer{
				identity: s.session.ticket, psk: s.session.psk, suite: s.sessionSuite,
				binderLabel: labelResumptionBinder, age: obfuscatedTicketAge(s.session, c.config.Time()),
				session: s.session, external: s.session.externalPSK,
			})
		}
		s.pskOffers = append(s.pskOffers, externalOffers...)
	} else if s.session != nil {
		hello.pskIdentity = append([]byte(nil), s.session.ticket...)
		hello.obfuscatedAge = obfuscatedTicketAge(s.session, c.config.Time())
	}
	if s.session != nil {
		if slices.Contains(hello.cipherSuites, s.sessionSuite.id) && hello.cipherSuites[0] != s.sessionSuite.id {
			prioritized := []uint16{s.sessionSuite.id}
			for _, id := range hello.cipherSuites {
				if id != s.sessionSuite.id {
					prioritized = append(prioritized, id)
				}
			}
			hello.cipherSuites = prioritized
		}
		if c.earlyIO && s.session.maxEarlyData > 0 && slices.Contains(hello.cipherSuites, s.sessionSuite.id) {
			hello.earlyData = true
		}
	}
	s.hello = hello
	return nil
}

// marshalClientHelloForOffers encodes hello with whichever PSK material the
// client is offering: a list of offers, a single legacy session binder, or
// none. binderTranscript and hrrBody are nil for the first ClientHello and set
// for the second one so the binder covers the HelloRetryRequest exchange.
func marshalClientHelloForOffers(hello *clientHello, offers []clientPSKOffer, session *ClientSessionState, sessionSuite *cipherSuite, binderTranscript, hrrBody []byte) ([]byte, error) {
	if len(offers) > 0 {
		return marshalClientHelloWithPSKOffers(hello, offers, binderTranscript, hrrBody)
	}
	if session != nil {
		return marshalClientHelloWithPSKBinder(hello, sessionSuite, session.psk, binderTranscript, hrrBody)
	}
	return hello.marshal()
}

// clientMarshalHello encodes the first ClientHello. With GREASE ECH the hello is
// encoded twice so the GREASE extension sees the final layout; with real ECH the
// encoded hello becomes the inner body and an outer hello replaces s.hello.
func (c *Conn) clientMarshalHello(s *clientHandshakeState) error {
	helloBody, err := marshalClientHelloForOffers(s.hello, s.pskOffers, s.session, s.sessionSuite, nil, nil)
	if err != nil {
		return err
	}
	if s.greaseECH {
		grease, greaseErr := generateGREASEECH(s.hello, c.config.Rand)
		if greaseErr != nil {
			return greaseErr
		}
		s.hello.setEncryptedClientHello(grease)
		if helloBody, err = marshalClientHelloForOffers(s.hello, s.pskOffers, s.session, s.sessionSuite, nil, nil); err != nil {
			return err
		}
	}
	if s.ech != nil {
		s.ech.innerHello = s.hello
		s.ech.innerBody = helloBody
		outer, outerErr := makeECHOuter(s.hello, s.ech.config, c.config.Rand)
		if outerErr != nil {
			return outerErr
		}
		if helloBody, err = computeOuterECH(outer, s.hello, s.ech, true); err != nil {
			return err
		}
		s.ech.outerHello = outer
		s.hello = outer
	}
	s.helloBody = helloBody
	return nil
}

// clientSendHello builds flight 1 from the encoded ClientHello and writes it.
func (c *Conn) clientSendHello(s *clientHandshakeState) error {
	out, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeClientHello, sequence: 0, body: s.helloBody}}, c.currentMTU(), 0, 0)
	if err != nil {
		return err
	}
	c.plainSendSequence.Store(out.nextRecordSequence())
	if err = c.writeFlight(c.conn, out); err != nil {
		return err
	}
	s.out = out
	return nil
}

// clientSendEarlyData installs epoch-1 keys when ClientHello advertised early
// data. The caller can then send datagrams before ServerHello arrives.
func (c *Conn) clientSendEarlyData(s *clientHandshakeState) error {
	if !s.hello.earlyData {
		return nil
	}
	c.earlyStatus.Store(uint32(EarlyDataPending))
	earlyHelloBody := s.helloBody
	if s.ech != nil {
		earlyHelloBody = s.ech.innerBody
	}
	earlyCipher, err := newEarlyRecordCipher(s.sessionSuite, s.session.psk, earlyHelloBody, s.session.recordSizeLimit, c.config.ReplayWindow)
	if err != nil {
		return err
	}
	c.writeMu.Lock()
	maxRecord := c.maxApplicationDatagramForCipher(earlyCipher)
	if maxRecord < 1 {
		c.writeMu.Unlock()
		return &ConfigError{"MTU is too small for early application data"}
	}
	c.earlyWriteCipher = earlyCipher
	c.earlyWriteRemaining = s.session.maxEarlyData
	c.writeMu.Unlock()
	if c.earlyIO {
		c.signalWriteReady()
	}
	return nil
}

// clientReceiveServerHello waits for exactly one ServerHello (or
// HelloRetryRequest) at epoch 0.
func (c *Conn) clientReceiveServerHello(s *clientHandshakeState) error {
	s.inbox = newHandshakeInbox(0, c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
	messages, err := c.receiveHandshakeWithRetransmit(s.inbox, nil, s.out)
	if err != nil {
		return err
	}
	if messages.len() != 1 || messages.at(0).typ != handshakeTypeServerHello {
		return alertError(alertUnexpectedMessage, &ProtocolError{"expected ServerHello"})
	}
	s.messages = messages
	s.serverHelloBody = messages.at(0).body
	c.stopEarlyWrites()
	return nil
}

// isHelloRetryRequestBody reports whether a ServerHello body carries the
// RFC 8446 HelloRetryRequest random sentinel.
func isHelloRetryRequestBody(body []byte) bool {
	return len(body) >= 34 && string(body[2:34]) == string(helloRetryRequestRandom[:])
}

// clientHandleHelloRetryRequest handles a HelloRetryRequest: it withdraws any
// early data, validates the request, generates a new key share if asked,
// re-encodes the ClientHello with the cookie and refreshed binders, sends
// flight 2, and receives the real ServerHello. It is a no-op when the first
// message was already a ServerHello.
func (c *Conn) clientHandleHelloRetryRequest(s *clientHandshakeState) error {
	if !isHelloRetryRequestBody(s.serverHelloBody) {
		return nil
	}
	hello := s.hello
	if hello.earlyData {
		c.markEarlyDataRejected()
		hello.earlyData = false
	}
	if s.ech != nil {
		s.ech.innerHello.earlyData = false
	}
	hrr, err := parseHelloRetryRequest(s.serverHelloBody)
	if err != nil {
		return err
	}
	suite, err := cipherSuiteForID(hrr.cipherSuite)
	if err != nil {
		return err
	}
	if !slices.Contains(hello.cipherSuites, hrr.cipherSuite) {
		return &ProtocolError{"HelloRetryRequest selected an unoffered cipher suite"}
	}
	s.helloRetrySuite = hrr.cipherSuite
	initial := newTranscriptHash(suite.hash.New())
	_ = initial.add(handshakeTypeClientHello, 0, s.helloBody)
	s.transcript = newTranscriptHash(suite.hash.New())
	_ = s.transcript.addHelloRetryRequest(initial.sumInto(s.transcriptDigest[:0]), s.serverHelloBody)
	if s.ech != nil {
		innerInitial := newTranscriptHash(suite.hash.New())
		_ = innerInitial.add(handshakeTypeClientHello, 0, s.ech.innerBody)
		s.ech.innerTranscript = newTranscriptHash(suite.hash.New())
		innerHash := innerInitial.sumInto(s.transcriptDigest[:0])
		accepted := false
		if hrr.hasECHConfirmation {
			zeroHRR := *hrr
			clear(zeroHRR.echConfirmation[:])
			zeroBody, marshalErr := zeroHRR.marshal()
			if marshalErr != nil {
				return marshalErr
			}
			confirmationTranscript := newTranscriptHash(suite.hash.New())
			_ = confirmationTranscript.addHelloRetryRequest(innerHash, zeroBody)
			want := echAcceptConfirmation(suite, s.ech.innerHello.random, "hrr ech accept confirmation", confirmationTranscript.sumInto(s.transcriptDigest[:0]))
			accepted = subtle.ConstantTimeCompare(want, hrr.echConfirmation[:]) == 1
		}
		_ = s.ech.innerTranscript.addHelloRetryRequest(innerInitial.sumInto(s.transcriptDigest[:0]), s.serverHelloBody)
		if accepted {
			s.ech.acceptedAtHRR = true
			hello = s.ech.innerHello
			initial = innerInitial
			s.transcript = s.ech.innerTranscript
		} else {
			s.ech.rejected = true
		}
	} else if hrr.hasECHConfirmation && !s.greaseECH {
		return alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited ECH confirmation in HelloRetryRequest"})
	}
	if hrr.selectedGroup != 0 {
		supported := slices.Contains(hello.supportedGroups, hrr.selectedGroup)
		alreadyOffered := slices.ContainsFunc(hello.keyShares, func(share keyShareEntry) bool { return share.group == hrr.selectedGroup })
		if !supported || alreadyOffered {
			return alertError(alertIllegalParameter, &ProtocolError{"HelloRetryRequest selected an invalid key share group"})
		}
		s.key, err = generateEphemeralKey(hrr.selectedGroup, c.config.Rand)
		if err != nil {
			return alertError(alertInternalError, err)
		}
		hello.keyShares = []keyShareEntry{{group: s.key.groupID(), data: s.key.publicBytes()}}
	}
	hello.cookie = hrr.cookie
	if len(s.pskOffers) > 0 {
		s.pskOffers = filterPSKOffersByHash(s.pskOffers, suite.hash)
		for i := range s.pskOffers {
			if s.pskOffers[i].session != nil {
				s.pskOffers[i].age = obfuscatedTicketAge(s.pskOffers[i].session, c.config.Time())
			}
		}
	}
	var helloBody []byte
	switch {
	case s.ech != nil && s.ech.rejected:
		helloBody, err = hello.marshal()
	case len(s.pskOffers) > 0:
		helloBody, err = marshalClientHelloWithPSKOffers(hello, s.pskOffers, initial.sumInto(s.transcriptDigest[:0]), s.serverHelloBody)
	case s.session != nil && s.sessionSuite.hash == suite.hash && len(hello.pskIdentity) > 0:
		hello.obfuscatedAge = obfuscatedTicketAge(s.session, c.config.Time())
		helloBody, err = marshalClientHelloWithPSKBinder(hello, s.sessionSuite, s.session.psk, initial.sumInto(s.transcriptDigest[:0]), s.serverHelloBody)
	default:
		hello.pskIdentity = nil
		hello.pskBinder = nil
		hello.pskIdentities = nil
		hello.pskBinders = nil
		helloBody, err = hello.marshal()
	}
	if err != nil {
		return err
	}
	wireHelloBody := helloBody
	if s.ech != nil && s.ech.acceptedAtHRR {
		s.ech.innerHello = hello
		s.ech.innerBody = helloBody
		s.ech.outerHello.cookie = append([]byte(nil), hello.cookie...)
		s.ech.outerHello.earlyData = false
		s.ech.outerHello.keyShares = cloneClientHello(hello).keyShares
		if wireHelloBody, err = computeOuterECH(s.ech.outerHello, hello, s.ech, false); err != nil {
			return err
		}
		_ = s.ech.innerTranscript.add(handshakeTypeClientHello, 1, helloBody)
		s.transcript = s.ech.innerTranscript
	} else {
		_ = s.transcript.add(handshakeTypeClientHello, 1, helloBody)
	}
	second, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeClientHello, sequence: 1, body: wireHelloBody}}, c.currentMTU(), 0, s.out.nextRecordSequence())
	if err != nil {
		return err
	}
	c.plainSendSequence.Store(second.nextRecordSequence())
	if err = c.writeFlight(c.conn, second); err != nil {
		return err
	}
	s.messages, err = c.receiveHandshakeWithRetransmit(s.inbox, nil, second)
	if err != nil {
		return err
	}
	if s.messages.len() != 1 || s.messages.at(0).typ != handshakeTypeServerHello {
		return alertError(alertUnexpectedMessage, &ProtocolError{"expected ServerHello after HelloRetryRequest"})
	}
	s.serverHelloBody = s.messages.at(0).body
	if isHelloRetryRequestBody(s.serverHelloBody) {
		return alertError(alertUnexpectedMessage, &ProtocolError{"received a second HelloRetryRequest"})
	}
	s.hello = hello
	s.helloBody = helloBody
	s.serverHelloSequence = 1
	s.clientFinishedSeq = 2
	s.serverHandshakeStart = 2
	return nil
}

// clientConfirmECH parses the ServerHello and, when real ECH was offered,
// checks the acceptance confirmation in the last eight bytes of random. On
// acceptance the inner hello and transcript become authoritative; otherwise
// ECH is recorded as rejected and the outer hello stays in effect.
func (c *Conn) clientConfirmECH(s *clientHandshakeState) error {
	sh, err := parseServerHello(s.serverHelloBody)
	if err != nil {
		return err
	}
	suite, err := cipherSuiteForID(sh.cipherSuite)
	if err != nil {
		return err
	}
	s.sh = sh
	s.suite = suite
	if s.ech == nil || s.ech.rejected {
		return nil
	}
	if s.ech.innerTranscript == nil {
		s.ech.innerTranscript = newTranscriptHash(suite.hash.New())
		_ = s.ech.innerTranscript.add(handshakeTypeClientHello, 0, s.ech.innerBody)
	}
	if len(s.serverHelloBody) < 34 {
		return alertError(alertDecodeError, &ProtocolError{"truncated ServerHello ECH confirmation"})
	}
	zeroServerHello := append([]byte(nil), s.serverHelloBody...)
	clear(zeroServerHello[26:34])
	confirmationTranscript := s.ech.innerTranscript.clone()
	_ = confirmationTranscript.add(handshakeTypeServerHello, s.serverHelloSequence, zeroServerHello)
	want := echAcceptConfirmation(suite, s.ech.innerHello.random, "ech accept confirmation", confirmationTranscript.sumInto(s.transcriptDigest[:0]))
	confirmed := subtle.ConstantTimeCompare(want, sh.random[24:]) == 1
	if s.ech.acceptedAtHRR && !confirmed {
		return alertError(alertIllegalParameter, &ProtocolError{"ServerHello did not confirm ECH after accepted HelloRetryRequest"})
	}
	if confirmed {
		s.echAccepted = true
		s.hello = s.ech.innerHello
		s.helloBody = s.ech.innerBody
		s.transcript = s.ech.innerTranscript
	} else {
		s.ech.rejected = true
	}
	return nil
}

// clientValidateServerHello enforces the ServerHello invariants the client
// must check: the suite was offered and matches any HRR, the legacy session ID
// is empty, the connection-ID and return-routability responses are valid, the
// PSK selection is consistent, and the key share group was offered. It also
// records the negotiated connection IDs on the connection.
func (c *Conn) clientValidateServerHello(s *clientHandshakeState) error {
	hello, sh, suite := s.hello, s.sh, s.suite
	if !slices.Contains(hello.cipherSuites, sh.cipherSuite) {
		return &ProtocolError{"server selected an unoffered cipher suite"}
	}
	if s.helloRetrySuite != 0 && sh.cipherSuite != s.helloRetrySuite {
		return alertError(alertIllegalParameter, &ProtocolError{"ServerHello changed cipher suite after HelloRetryRequest"})
	}
	if len(sh.sessionID) != 0 {
		return &ProtocolError{"ServerHello legacy session ID must be empty"}
	}
	if err := validateServerHelloConnectionID(hello, sh); err != nil {
		return err
	}
	if err := validateServerHelloReturnRoutability(hello, sh); err != nil {
		return err
	}
	if sh.hasConnectionID {
		c.recordConnectionIDNegotiation(hello.connectionID, sh.connectionID)
	}
	c.returnRoutabilityCheckNegotiated = sh.returnRoutability
	s.usingPSK = sh.selectedIdentity != nil
	if s.ech != nil && s.ech.rejected && s.usingPSK {
		return alertError(alertIllegalParameter, &ProtocolError{"server selected an outer GREASE PSK after rejecting ECH"})
	}
	if s.usingPSK {
		if len(s.pskOffers) == 0 && s.session != nil && *sh.selectedIdentity == 0 {
			s.singleTicketOffer = clientPSKOffer{psk: s.session.psk, suite: s.sessionSuite, session: s.session, external: s.session.externalPSK}
			s.selectedOffer = &s.singleTicketOffer
		} else {
			if int(*sh.selectedIdentity) >= len(s.pskOffers) {
				return &ProtocolError{"server selected an invalid PSK identity"}
			}
			s.selectedOffer = &s.pskOffers[*sh.selectedIdentity]
		}
		if suite.hash != s.selectedOffer.suite.hash {
			return &ProtocolError{"server selected a cipher suite incompatible with the PSK"}
		}
		s.resumed = s.selectedOffer.session != nil
		s.externalPSK = s.selectedOffer.external
	}
	if !slices.ContainsFunc(hello.keyShares, func(share keyShareEntry) bool { return share.group == sh.keyShare.group }) {
		return alertError(alertIllegalParameter, &ProtocolError{"ServerHello selected an unoffered key share"})
	}
	return nil
}

// clientDeriveHandshakeKeys computes the shared secret, completes the
// transcript through ServerHello, derives the handshake traffic secrets, and
// installs the epoch-2 ciphers. c.sendCipher is set here so that from this
// point the deferred alert sender in runHandshake encrypts alerts; moving this
// assignment would change which alerts go out in plaintext.
func (c *Conn) clientDeriveHandshakeKeys(s *clientHandshakeState) error {
	shared, err := s.key.sharedSecret(s.sh.keyShare.group, s.sh.keyShare.data)
	if err != nil {
		return err
	}
	if s.transcript == nil {
		s.transcript = newTranscriptHash(s.suite.hash.New())
		_ = s.transcript.add(handshakeTypeClientHello, 0, s.helloBody)
	}
	_ = s.transcript.add(handshakeTypeServerHello, s.serverHelloSequence, s.serverHelloBody)
	var psk []byte
	if s.usingPSK {
		psk = s.selectedOffer.psk
	}
	s.schedule = newKeySchedule(s.suite, psk)
	if err = s.schedule.deriveHandshake(shared, s.transcript.sumInto(s.transcriptDigest[:0])); err != nil {
		return err
	}
	s.receiveCipher, s.sendCipher, err = c.newHandshakeCiphers(s.suite, s.schedule.serverHandshakeTraffic, s.schedule.clientHandshakeTraffic)
	if err != nil {
		return err
	}
	c.writeMu.Lock()
	c.sendCipher = s.sendCipher
	c.writeMu.Unlock()
	s.inbox = newHandshakeInbox(s.serverHandshakeStart, c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
	return nil
}

// clientProcessServerFlight receives and validates the server's encrypted
// flight: EncryptedExtensions, optional CertificateRequest, Certificate,
// CertificateVerify, and Finished. Message order is enforced by
// serverHandshakeStage; the per-message checks live in the client*Message
// helpers so the loop body stays a dispatch table.
func (c *Conn) clientProcessServerFlight(s *clientHandshakeState) error {
	if s.resumed {
		s.peerRawPublicKey = append([]byte(nil), s.selectedOffer.session.peerRawPublicKey...)
		s.peerCerts = append([]*x509.Certificate(nil), s.selectedOffer.session.peerCertificates...)
		s.peerDelegatedCredential = append([]byte(nil), s.selectedOffer.session.peerDelegatedCredential...)
		s.chains = make([][]*x509.Certificate, len(s.selectedOffer.session.verifiedChains))
		for i := range s.selectedOffer.session.verifiedChains {
			s.chains[i] = append([]*x509.Certificate(nil), s.selectedOffer.session.verifiedChains[i]...)
		}
		s.ocspResponse = append([]byte(nil), s.selectedOffer.session.ocspResponse...)
	}
	verifiedServerSignature := false
	stage := serverExpectEncryptedExtensions
	for stage != serverHandshakeComplete {
		messages, err := receiveHandshakeMessageWithEarlyBatch(c.conn, s.inbox, s.receiveCipher, handshakeReceiveOptions{
			ackCipher: s.sendCipher,
			mtu:       c.currentMTU(),
			owner:     c,
		})
		if err != nil {
			return err
		}
		for index := 0; index < messages.len(); index++ {
			message := messages.at(index)
			if err = stage.accept(message.typ, s.usingPSK); err != nil {
				return err
			}
			switch message.typ {
			case handshakeTypeEncryptedExtensions:
				err = c.clientEncryptedExtensions(s, message)
			case handshakeTypeCertificate, handshakeTypeCompressedCertificate:
				err = c.clientServerCertificate(s, message)
			case handshakeTypeCertificateRequest:
				err = c.clientCertificateRequest(s, message)
			case handshakeTypeCertificateVerify:
				err = c.clientServerCertificateVerify(s, message)
				verifiedServerSignature = err == nil
			case handshakeTypeFinished:
				if !s.usingPSK && ((len(s.peerCerts) == 0 && len(s.peerRawPublicKey) == 0) || !verifiedServerSignature) {
					return &ProtocolError{"server authentication messages are incomplete"}
				}
				err = c.clientServerFinished(s, message)
			default:
				return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected server handshake message"})
			}
			if err != nil {
				return err
			}
		}
	}
	if s.clientCertificateTypeSelected && s.certificateRequest == nil && !s.hello.postHandshakeAuth {
		return alertError(alertIllegalParameter, &ProtocolError{"server selected client certificate type without CertificateRequest"})
	}
	return nil
}

func (c *Conn) clientEncryptedExtensions(s *clientHandshakeState, message completedHandshake) error {
	ee, err := parseEncryptedExtensions(message.body)
	if err != nil {
		return err
	}
	negotiated, acceptedEarly, retryConfigs, err := validateEncryptedExtensions(s.hello, &ee)
	if err != nil {
		return err
	}
	s.negotiated = negotiated
	if s.echAccepted && c.config.HandshakeMetadata != nil && ee.handshakeMetadata == nil {
		return alertError(alertMissingExtension, &ProtocolError{"server did not accept handshake metadata"})
	}
	if ee.handshakeMetadata != nil && !s.echAccepted {
		return alertError(alertIllegalParameter, &ProtocolError{"handshake metadata requires ECH acceptance"})
	}
	s.handshakeMetadata = ee.handshakeMetadata
	if ee.cachedInformation != 0 && (s.cachedInfo == nil || s.usingPSK) {
		return alertError(alertIllegalParameter, &ProtocolError{"cached_info selected without certificate authentication"})
	}
	if s.cachedInfo != nil {
		s.cachedInfo.selected = ee.cachedInformation
	}
	c.serverCertificateType, c.clientCertificateType = ee.serverCertificateType, ee.clientCertificateType
	s.clientCertificateTypeSelected = ee.hasClientCertificateType
	if !s.usingPSK && (s.ech == nil || !s.ech.rejected) && !supportsCertificateType(c.config.ServerCertificateTypes, c.serverCertificateType) {
		return alertError(alertUnsupportedCertificate, &ProtocolError{"server did not negotiate an acceptable certificate type"})
	}
	if ee.hasTicketRequest {
		requested := s.hello.ticketRequest.NewSessionCount
		if s.resumed {
			requested = s.hello.ticketRequest.ResumptionCount
		}
		c.sessionTicketRequest = &sessionTicketRequestState{limit: min(ee.expectedTicketCount, requested)}
	}
	if retryConfigs != nil {
		if s.ech != nil && s.echAccepted {
			return alertError(alertUnsupportedExtension, &ProtocolError{"server sent ECH retry configurations after accepting ECH"})
		}
		if s.ech != nil && s.ech.rejected {
			s.ech.retryConfigs = append([]byte(nil), retryConfigs...)
		}
		// GREASE clients validate but never retain retry configurations.
	}
	if ee.hasRecordSizeLimit {
		c.recordRecordSizeLimit(c.config.RecordSizeLimit, ee.recordSizeLimit)
		s.receiveCipher.setPlaintextLimit(c.localRecordSizeLimit)
		s.sendCipher.setPlaintextLimit(c.peerRecordSizeLimit)
	}
	if err = validateEarlyDataSelection(acceptedEarly, s.sh.selectedIdentity); err != nil {
		return err
	}
	if acceptedEarly {
		c.earlyMu.Lock()
		c.earlyAccepted = true
		c.earlyMu.Unlock()
	} else if s.hello.earlyData {
		c.markEarlyDataRejected()
	}
	_ = s.transcript.add(message.typ, message.sequence, message.body)
	return nil
}

func (c *Conn) clientServerCertificate(s *clientHandshakeState, message completedHandshake) error {
	if s.usingPSK {
		return &ProtocolError{"server sent Certificate in a PSK handshake"}
	}
	if s.cachedInfo != nil && s.cachedInfo.selected&cachedCertRequest != 0 && s.certificateRequest == nil {
		return alertError(alertIllegalParameter, &ProtocolError{"cached_info selected without CertificateRequest"})
	}
	certMsg, err := c.parseCachedServerCertificate(s, message)
	if err != nil {
		return err
	}
	if err = validateCertificateMessageWithRequests(certMsg, nil, s.hello.statusRequest && c.serverCertificateType == CertificateTypeX509, len(s.hello.delegatedCredentialSchemes) > 0 && c.serverCertificateType == CertificateTypeX509); err != nil {
		return err
	}
	if len(certMsg.certificates) > 0 {
		s.ocspResponse = append([]byte(nil), certMsg.certificates[0].ocspResponse...)
	}
	certificateSchemes := s.hello.certificateSignatureSchemes
	if len(certificateSchemes) == 0 {
		certificateSchemes = s.hello.signatureSchemes
	}
	if c.serverCertificateType == CertificateTypeRawPublicKey {
		if s.ech != nil && s.ech.rejected {
			return alertError(alertBadCertificate, &ProtocolError{"ECH rejection requires X.509 public-name authentication"})
		}
		s.peerRawPublicKey, err = verifyRawPublicKeyMessage(c.config, certMsg)
	} else if s.ech != nil && s.ech.rejected {
		s.peerCerts, s.chains, err = verifyCertificateChainForECHRejection(c.config, certMsg, certificateSchemes, s.ech.config.publicName)
	} else {
		s.peerCerts, s.chains, err = verifyCertificateChain(c.config, certMsg, true, certificateSchemes)
	}
	if err != nil {
		return err
	}
	s.peerDelegatedCredential, err = peerDelegatedCredential(certMsg, s.peerCerts, c.config, true, s.hello.signatureSchemes, s.hello.delegatedCredentialSchemes)
	if err != nil {
		return err
	}
	if s.ech != nil && s.ech.rejected && c.config.EncryptedClientHelloRejectionVerify != nil {
		rejectionState := ConnectionState{Version: VersionDTLS13, CipherSuite: s.suite.id, NegotiatedProtocol: s.negotiated, ServerName: s.ech.config.publicName, PeerCertificates: s.peerCerts, VerifiedChains: s.chains, OCSPResponse: append([]byte(nil), s.ocspResponse...), ECHAccepted: false}
		if verifyErr := c.config.EncryptedClientHelloRejectionVerify(rejectionState); verifyErr != nil {
			return alertError(alertAccessDenied, verifyErr)
		}
	}
	_ = s.transcript.add(message.typ, message.sequence, message.body)
	return nil
}

func (c *Conn) clientCertificateRequest(s *clientHandshakeState, message completedHandshake) error {
	if s.usingPSK {
		return &ProtocolError{"server requested a certificate in a PSK handshake"}
	}
	if s.certificateRequest != nil {
		return alertError(alertUnexpectedMessage, &ProtocolError{"duplicate CertificateRequest"})
	}
	body, err := s.cachedInfo.resolve(handshakeTypeCertificateRequest, message.body, c.config.MaxHandshakeMessage)
	if err != nil {
		return err
	}
	request, compression, err := parseCertificateRequestWithCompression(body)
	if err != nil {
		return err
	}
	if len(request.requestContext) != 0 {
		return alertError(alertIllegalParameter, &ProtocolError{"initial CertificateRequest context must be empty"})
	}
	s.certificateRequest = request
	s.certReqCompression = compression
	_ = s.transcript.add(message.typ, message.sequence, message.body)
	return nil
}

func (c *Conn) clientServerCertificateVerify(s *clientHandshakeState, message completedHandshake) error {
	if len(s.peerCerts) == 0 && len(s.peerRawPublicKey) == 0 {
		return alertError(alertUnexpectedMessage, &ProtocolError{"CertificateVerify before Certificate"})
	}
	cv, err := parseCertificateVerify(message.body)
	if err != nil {
		return err
	}
	if err = verifyPeerAuthenticationSignature(s.peerCerts, s.peerRawPublicKey, s.peerDelegatedCredential, s.hello.signatureSchemes, cv, s.transcript.sumInto(s.transcriptDigest[:0]), true); err != nil {
		return err
	}
	_ = s.transcript.add(message.typ, message.sequence, message.body)
	return nil
}

func (c *Conn) clientServerFinished(s *clientHandshakeState, message completedHandshake) error {
	verify, err := parseFinished(message.body, s.suite.hash.Size())
	if err != nil {
		return err
	}
	if !s.schedule.verifyFinished(s.schedule.serverHandshakeTraffic, s.transcript.sumInto(s.transcriptDigest[:0]), verify) {
		return alertError(alertDecryptError, &ProtocolError{"server Finished verification failed"})
	}
	_ = s.transcript.add(message.typ, message.sequence, message.body)
	s.serverFinishedSequence = message.sequence
	return nil
}

// clientDeriveApplicationSecrets derives the application traffic secrets from
// the transcript through the server Finished.
func (c *Conn) clientDeriveApplicationSecrets(s *clientHandshakeState) error {
	if c.earlyAccepted {
		c.earlyStatus.Store(uint32(EarlyDataAccepted))
	}
	return s.schedule.deriveApplication(s.transcript.sumInto(s.transcriptDigest[:0]))
}

// clientSendAuthAndFinished answers a CertificateRequest with Certificate and
// CertificateVerify, appends the client Finished, and writes the protected
// flight. The client never sends a certificate on a rejected-ECH connection.
func (c *Conn) clientSendAuthAndFinished(s *clientHandshakeState) error {
	var clientMessages []handshakeMessage
	nextClientSequence := s.clientFinishedSeq
	if s.certificateRequest != nil {
		certMessage := &certificateMessage{requestContext: s.certificateRequest.requestContext}
		var clientCertificate *tls.Certificate
		if s.ech == nil || !s.ech.rejected {
			var err error
			clientCertificate, err = c.selectClientCertificate(s.certificateRequest)
			if err != nil {
				return err
			}
			if clientCertificate != nil {
				for _, der := range clientCertificate.Certificate {
					certMessage.certificates = append(certMessage.certificates, certificateEntry{data: der})
				}
				if s.certificateRequest.statusRequest && len(clientCertificate.OCSPStaple) > 0 && len(certMessage.certificates) > 0 {
					response, responseErr := marshalOCSPResponse(clientCertificate.OCSPStaple)
					if responseErr != nil {
						return responseErr
					}
					certMessage.certificates[0].extensions = map[uint16][]byte{extStatusRequest: response}
				}
				c.addDelegatedCredential(certMessage, clientCertificate)
			}
		}
		var signer crypto.Signer
		var scheme tls.SignatureScheme
		if clientCertificate != nil && len(clientCertificate.Certificate) > 0 {
			var ok bool
			signer, ok = clientCertificate.PrivateKey.(crypto.Signer)
			if !ok {
				return errors.New("dtls13: client private key does not implement crypto.Signer")
			}
			var err error
			if scheme, err = c.certificateSigningScheme(clientCertificate, s.certificateRequest.signatureSchemes); err != nil {
				return err
			}
		}
		messages, next, err := c.buildCertificateMessages(certMessage, s.certReqCompression, signer, scheme, false, s.transcript, nextClientSequence, s.transcriptDigest[:0])
		if err != nil {
			return err
		}
		clientMessages = append(clientMessages, messages...)
		nextClientSequence = next
	}
	s.clientFinishedSeq = nextClientSequence
	clientFinished := s.schedule.finishedVerifyData(s.schedule.clientHandshakeTraffic, s.transcript.sumInto(s.transcriptDigest[:0]))
	_ = s.transcript.add(handshakeTypeFinished, s.clientFinishedSeq, clientFinished)
	clientMessages = append(clientMessages, handshakeMessage{typ: handshakeTypeFinished, sequence: s.clientFinishedSeq, body: clientFinished})
	clientFlight, err := buildProtectedFlight(clientMessages, c.currentMTU(), s.sendCipher)
	if err != nil {
		return err
	}
	if err = c.writeFlight(c.conn, clientFlight); err != nil {
		return err
	}
	s.out = clientFlight
	return nil
}

// clientAwaitFinalACK waits for the server to acknowledge every record of the
// client's final flight, accepting the ACK at either the handshake epoch or the
// server's application epoch, then installs the application keys.
func (c *Conn) clientAwaitFinalACK(s *clientHandshakeState) error {
	var applicationACKCipher *recordCipher
	if c.earlyIO && (s.ech == nil || !s.ech.rejected) {
		if err := c.installApplicationKeysAt(s.suite, s.schedule.clientApplicationTraffic, s.schedule.serverApplicationTraffic, s.clientFinishedSeq+1); err != nil {
			return err
		}
		if err := c.clientFinalize(s); err != nil {
			return err
		}
		c.signalApplicationReady()
	} else {
		var err error
		applicationACKCipher, err = newRecordCipher(s.suite, s.schedule.serverApplicationTraffic, 3, c.config.ReplayWindow)
		if err != nil {
			return err
		}
		applicationACKCipher.setPlaintextLimit(c.localRecordSizeLimit)
		if c.connectionIDNegotiated {
			if err = applicationACKCipher.setConnectionID(c.receiveConnectionID); err != nil {
				return err
			}
		}
	}
	acknowledged, err := c.receiveACKWithRetransmit(s.out, s.receiveCipher, applicationACKCipher)
	if err != nil {
		return err
	}
	if len(s.out.records) == 0 {
		return &ProtocolError{"empty client Finished flight"}
	}
	for _, record := range s.out.records {
		if !record.acknowledgedBy(acknowledged) {
			return &ProtocolError{"server ACK did not cover complete client flight"}
		}
	}
	return c.installApplicationKeysAt(s.suite, s.schedule.clientApplicationTraffic, s.schedule.serverApplicationTraffic, s.clientFinishedSeq+1)
}

// clientFinalize reports an authenticated ECH rejection, otherwise derives the
// resumption secret, installs the receive epoch, publishes the connection
// state, and discards the unused session group.
func (c *Conn) clientFinalize(s *clientHandshakeState) error {
	if s.finalized {
		return nil
	}
	if s.ech != nil && s.ech.rejected {
		return alertError(alertECHRequired, &ECHRejectionError{RetryConfigList: append([]byte(nil), s.ech.retryConfigs...)})
	}
	if err := c.finishHandshake(handshakeCompletion{
		suite:                   s.suite,
		schedule:                s.schedule,
		transcript:              s.transcript,
		receiveCipher:           s.receiveCipher,
		peerFlightStart:         s.serverHandshakeStart,
		peerFlightEnd:           s.serverFinishedSequence,
		externalPSK:             s.externalPSK,
		resumed:                 s.resumed,
		echAccepted:             s.echAccepted,
		negotiated:              s.negotiated,
		peerCerts:               s.peerCerts,
		peerRawPublicKey:        s.peerRawPublicKey,
		peerDelegatedCredential: s.peerDelegatedCredential,
		handshakeMetadata:       s.handshakeMetadata,
		chains:                  s.chains,
		serverName:              c.config.ServerName,
		ocspResponse:            s.ocspResponse,
	}); err != nil {
		return err
	}
	// A client retains the resumption secret so it can derive PSKs from the
	// NewSessionTickets the server sends after the handshake.
	c.resumptionSuite = s.suite
	c.resumptionMasterSecret = append([]byte(nil), s.schedule.resumptionMasterSecret...)
	if c.sessionTicketRequest != nil {
		if s.resumed {
			c.sessionTicketRequest.group = s.selectedOffer.session.ticketGroup
		} else {
			c.sessionTicketRequest.group = sha256.Sum256(c.resumptionMasterSecret)
		}
	}
	if s.session != nil && !s.resumed {
		discardClientSessionGroup(c.config, c.conn, s.session.ticketGroup)
	}
	if s.cachedInfo != nil && !s.usingPSK {
		c.config.CachedInformationCache.put(clientSessionCacheKey(c.config, c.conn), s.cachedInfo.received)
	}
	s.finalized = true
	return nil
}
