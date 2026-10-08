package dtls13

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"slices"
)

// serverHandshakeState carries the values that cross phase boundaries in the
// server handshake. Fields are grouped by the step that first sets them. The
// struct holds no back-pointer to *Conn; steps are *Conn methods that receive
// the state explicitly.
type serverHandshakeState struct {
	// transcriptDigest is scratch space for every transcript sum.
	transcriptDigest [maxSupportedHashSize]byte

	// serverReceiveClientHello
	amplification     amplificationGuard
	preValidationConn *amplificationConn
	inbox             *handshakeInbox
	messages          completedHandshakeBatch
	helloBody         []byte
	outerHello        *clientHello
	ch                *clientHello

	// serverProcessECH
	echContext  *echServerContext
	echKeys     []EncryptedClientHelloKey
	echAccepted bool
	echRejected bool

	// serverSelectInitial
	initialClientHello    clientHello
	suite                 *cipherSuite
	hrrUsed               bool
	serverSequenceOffset  uint16
	clientFinishedSeq     uint16
	initialHashTranscript *transcriptHash

	// serverSendHelloRetryRequest
	hrrBody        []byte
	hrrFlight      *flight
	requestedGroup tls.CurveID
	cookieAddress  []byte

	// PSK selection (serverTryEarlyDataWithoutCookie or serverRenegotiateAfterHRR)
	psk                 []byte
	resumed             bool
	usingPSK            bool
	resumedSession      *sessionTicketState
	externalPSK         *externalPSKSelection
	selectedPSKIdentity uint16

	// key share and ALPN
	share      keyShareEntry
	negotiated string

	// serverNegotiateExtensionsAndKeyShare
	shared      []byte
	serverShare []byte
	cert        *tls.Certificate
	signer      crypto.Signer
	scheme      tls.SignatureScheme

	// serverBuildServerHello
	sh                       *serverHello
	shBody                   []byte
	transcript               *transcriptHash
	serverHelloSequence      uint16
	firstPlainRecordSequence uint64

	// serverDeriveHandshakeKeys
	schedule     *keySchedule
	serverCipher *recordCipher
	clientCipher *recordCipher
	earlyCipher  *recordCipher

	// serverSendFlight
	serverFlight   *flight
	ticketCount    uint8
	serverSequence uint16
	// What the CertificateRequest told the client, kept to validate its reply.
	clientSignatureSchemes       []tls.SignatureScheme
	clientCertificateSchemes     []tls.SignatureScheme
	clientCertificateOIDFilters  []CertificateOIDFilter
	clientCertificateCompression *certificateCompressionAlgorithms
	clientStatusRequest          bool

	// serverProcessClientFlight
	clientFinalFlightStart uint16
	clientCerts            []*x509.Certificate
	clientRawPublicKey     []byte
	clientTypeNegotiated   bool
	clientChains           [][]*x509.Certificate
	ocspResponse           []byte
	clientAuthAt           int64
	clientRecords          []recordNumber
}

// serverHandshake runs the DTLS 1.3 server handshake as a sequence of steps.
// Error semantics are unchanged from the monolithic implementation; c.sendCipher
// is set inside serverSendFlight, immediately before the server flight is
// written, exactly where it was before the split.
func (c *Conn) serverHandshake() error {
	s := &serverHandshakeState{
		hrrUsed:              true,
		serverSequenceOffset: 1,
		clientFinishedSeq:    2,
	}
	s.preValidationConn = &amplificationConn{Conn: c.conn, guard: &s.amplification}
	steps := []func(*serverHandshakeState) error{
		c.serverReceiveClientHello,
		c.serverProcessECH,
		c.serverSelectInitial,
		c.serverTryEarlyDataWithoutCookie,
		c.serverCookieExchange,
		c.serverNegotiateExtensionsAndKeyShare,
		c.serverBuildServerHello,
		c.serverDeriveHandshakeKeys,
		c.serverSendFlight,
		c.serverDeriveApplicationSecrets,
		c.serverProcessClientFlight,
		c.serverSendFinalACK,
		c.serverFinalize,
		c.serverIssueTickets,
	}
	for _, step := range steps {
		if err := step(s); err != nil {
			return err
		}
	}
	return nil
}

// serverReceiveClientHello reads exactly one ClientHello through the
// amplification-limited connection and parses it.
func (c *Conn) serverReceiveClientHello(s *serverHandshakeState) error {
	s.inbox = newHandshakeInbox(0, c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
	messages, err := receiveHandshakeMessageWithEarlyBatch(s.preValidationConn, s.inbox, nil, handshakeReceiveOptions{owner: c, mtu: c.currentMTU()})
	if err != nil {
		return err
	}
	if messages.len() != 1 || messages.at(0).typ != handshakeTypeClientHello {
		return alertError(alertUnexpectedMessage, &ProtocolError{"expected ClientHello"})
	}
	s.messages = messages
	s.helloBody = messages.at(0).body
	s.outerHello, err = parseClientHello(s.helloBody)
	if err != nil {
		return err
	}
	s.ch = s.outerHello
	return nil
}

// serverProcessECH decrypts an outer-type ECH extension, replacing the working
// ClientHello with the inner one on success and recording acceptance or
// rejection.
func (c *Conn) serverProcessECH(s *serverHandshakeState) error {
	ech := s.outerHello.encryptedClientHello()
	if len(ech) == 0 {
		return nil
	}
	typ, _, _, _, _, err := parseECHExt(ech)
	if err != nil {
		return err
	}
	outerOffered := typ == echOuterType
	if outerOffered {
		s.echKeys = c.config.EncryptedClientHelloKeys
		if c.config.GetEncryptedClientHelloKeys != nil {
			if s.echKeys, err = c.config.GetEncryptedClientHelloKeys(c.clientHelloInfo(s.outerHello)); err != nil {
				return err
			}
		}
	}
	s.ch, s.helloBody, s.echContext, err = processECHClientHello(s.outerHello, s.helloBody, s.echKeys)
	if err != nil {
		return err
	}
	s.echAccepted = s.echContext != nil
	s.echRejected = outerOffered && !s.echAccepted
	return nil
}

// selectServerCipherSuite picks the cipher suite for a ClientHello, preferring
// one compatible with a configured external PSK.
func (c *Conn) selectServerCipherSuite(ch *clientHello) (*cipherSuite, error) {
	suite, err := selectCipherSuite(c.config.CipherSuites, ch.cipherSuites)
	if err != nil {
		return nil, err
	}
	if len(c.config.ExternalPSKs) > 0 {
		suite = preferExternalPSKCipherSuite(c.config, ch, suite)
	}
	return suite, nil
}

// serverSelectInitial records the client's post-handshake-auth offer, snapshots
// ClientHello1 without cookie or binders for the HRR comparison, picks the
// initial cipher suite, and hashes ClientHello1.
func (c *Conn) serverSelectInitial(s *serverHandshakeState) error {
	c.postHandshakeAuthOffered = s.ch.postHandshakeAuth
	s.initialClientHello = *s.ch
	s.initialClientHello.cookie = nil
	s.initialClientHello.pskBinder = nil
	s.initialClientHello.pskBinders = nil
	suite, err := c.selectServerCipherSuite(s.ch)
	if err != nil {
		return err
	}
	s.suite = suite
	s.initialHashTranscript = newTranscriptHash(suite.hash.New())
	_ = s.initialHashTranscript.add(handshakeTypeClientHello, 0, s.helloBody)
	return nil
}

// applyPSKSelection copies an accepted PSK selection into the state.
func (s *serverHandshakeState) applyPSKSelection(selection selectedPSK) {
	s.psk = selection.psk
	s.usingPSK = true
	s.selectedPSKIdentity = selection.identity
	s.resumedSession = selection.session
	s.externalPSK = selection.external
	s.resumed = s.resumedSession != nil
}

// serverTryEarlyDataWithoutCookie implements the trusted-environment shortcut
// that skips the cookie exchange for a resumed PSK handshake. This is the only
// path on which epoch-1 data is accepted; the default remains the
// RFC-recommended HRR/cookie exchange.
func (c *Conn) serverTryEarlyDataWithoutCookie(s *serverHandshakeState) error {
	if !c.config.AllowEarlyDataWithoutCookie || !s.ch.earlyData || len(s.ch.pskIdentity) == 0 {
		return nil
	}
	if _, err := selectKeyShare(c.config.CurvePreferences, s.ch.keyShares); err != nil {
		return nil
	}
	candidateProtocol, err := negotiateALPN(c.config.NextProtos, s.ch.alpn)
	if err != nil {
		return nil
	}
	candidate, accepted, err := c.acceptPSK(s.ch, s.helloBody, s.initialHashTranscript.sumInto(s.transcriptDigest[:0]), nil, s.suite, candidateProtocol)
	if err != nil {
		return err
	}
	if accepted && candidate.session != nil {
		s.applyPSKSelection(candidate)
		s.hrrUsed = false
		s.serverSequenceOffset = 0
		s.clientFinishedSeq = 1
	}
	return nil
}

// serverCookieExchange runs the HelloRetryRequest/cookie round trip when the
// no-cookie shortcut did not apply, then negotiates the key share, ALPN, and
// PSK against the ClientHello that will drive the rest of the handshake.
func (c *Conn) serverCookieExchange(s *serverHandshakeState) error {
	if !s.hrrUsed {
		share, err := selectKeyShare(c.config.CurvePreferences, s.ch.keyShares)
		if err != nil {
			return err
		}
		s.share = share
		s.negotiated, err = negotiateALPN(c.config.NextProtos, s.ch.alpn)
		return err
	}
	if err := c.serverSendHelloRetryRequest(s); err != nil {
		return err
	}
	if err := c.serverReceiveSecondClientHello(s); err != nil {
		return err
	}
	return c.serverRenegotiateAfterHRR(s)
}

// serverSendHelloRetryRequest seals a cookie over ClientHello1, chooses a
// replacement group when no offered key share is acceptable, and sends the
// HelloRetryRequest in plaintext through the amplification guard.
func (c *Conn) serverSendHelloRetryRequest(s *serverHandshakeState) error {
	if err := ensureCookieProtector(c.config); err != nil {
		return err
	}
	protector := &c.config.state.cookieProtector
	s.cookieAddress = []byte(c.conn.RemoteAddr().String())
	cookie, err := protector.seal(s.cookieAddress, s.initialHashTranscript.sumInto(s.transcriptDigest[:0]))
	if err != nil {
		return err
	}
	if _, shareErr := selectKeyShare(c.config.CurvePreferences, s.ch.keyShares); shareErr != nil {
		for _, preference := range c.config.CurvePreferences {
			if slices.Contains(s.ch.supportedGroups, preference) {
				s.requestedGroup = preference
				break
			}
		}
		if s.requestedGroup == 0 {
			return shareErr
		}
	}
	hrr := &helloRetryRequest{cipherSuite: s.suite.id, cookie: cookie, selectedGroup: s.requestedGroup}
	if s.echAccepted {
		hrr.hasECHConfirmation = true
		zeroBody, marshalErr := hrr.marshal()
		if marshalErr != nil {
			return marshalErr
		}
		confirmationTranscript := newTranscriptHash(s.suite.hash.New())
		_ = confirmationTranscript.addHelloRetryRequest(s.initialHashTranscript.sumInto(s.transcriptDigest[:0]), zeroBody)
		copy(hrr.echConfirmation[:], echAcceptConfirmation(s.suite, s.ch.random, "hrr ech accept confirmation", confirmationTranscript.sumInto(s.transcriptDigest[:0])))
	}
	if s.hrrBody, err = hrr.marshal(); err != nil {
		return err
	}
	if s.hrrFlight, _, err = buildPlainFlight([]handshakeMessage{{typ: handshakeTypeServerHello, sequence: 0, body: s.hrrBody}}, c.currentMTU(), 0, 0); err != nil {
		return err
	}
	return c.writeFlight(s.preValidationConn, s.hrrFlight)
}

// serverReceiveSecondClientHello receives ClientHello2, reconstructs the inner
// hello when ECH was accepted, and checks that it differs from ClientHello1 only
// as the HelloRetryRequest permitted and that the cookie is valid. A valid
// cookie proves address reachability, so the amplification limit is lifted.
func (c *Conn) serverReceiveSecondClientHello(s *serverHandshakeState) error {
	messages, err := c.receiveSecondClientHello(s.preValidationConn, s.inbox, s.hrrFlight)
	if err != nil {
		return err
	}
	if messages.len() != 1 || messages.at(0).typ != handshakeTypeClientHello {
		return alertError(alertUnexpectedMessage, &ProtocolError{"expected second ClientHello"})
	}
	s.messages = messages
	secondBody := messages.at(0).body
	second, err := parseClientHello(secondBody)
	if err != nil {
		return err
	}
	if s.echAccepted {
		if second, secondBody, err = processSecondECHClientHello(second, secondBody, s.echContext); err != nil {
			return err
		}
	}
	if !equalClientHelloAfterHRR(&s.initialClientHello, second, s.requestedGroup) {
		return &ProtocolError{"second ClientHello changed fields other than cookie"}
	}
	compatibleRemoval, err := removedPSKsAreIncompatible(c.config, &s.initialClientHello, second, s.suite)
	if err != nil {
		return err
	}
	if !compatibleRemoval {
		return alertError(alertIllegalParameter, &ProtocolError{"second ClientHello removed a PSK compatible with the HelloRetryRequest cipher suite"})
	}
	protector := &c.config.state.cookieProtector
	cookieHash, err := protector.open(s.cookieAddress, second.cookie)
	if err != nil {
		return alertError(alertIllegalParameter, &ProtocolError{"invalid HelloRetryRequest cookie"})
	}
	if string(cookieHash) != string(s.initialHashTranscript.sumInto(s.transcriptDigest[:0])) {
		return alertError(alertIllegalParameter, &ProtocolError{"HelloRetryRequest cookie transcript mismatch"})
	}
	s.amplification.validate()
	s.helloBody = secondBody
	s.ch = second
	return nil
}

// serverRenegotiateAfterHRR reselects the cipher suite, key share, ALPN, and
// PSK against ClientHello2.
func (c *Conn) serverRenegotiateAfterHRR(s *serverHandshakeState) error {
	suite, err := c.selectServerCipherSuite(s.ch)
	if err != nil {
		return err
	}
	s.suite = suite
	if s.share, err = selectKeyShare(c.config.CurvePreferences, s.ch.keyShares); err != nil {
		return err
	}
	if s.negotiated, err = negotiateALPN(c.config.NextProtos, s.ch.alpn); err != nil {
		return err
	}
	selection, accepted, err := c.acceptPSK(s.ch, s.helloBody, s.initialHashTranscript.sumInto(s.transcriptDigest[:0]), s.hrrBody, s.suite, s.negotiated)
	if err != nil {
		return err
	}
	if accepted {
		s.applyPSKSelection(selection)
	}
	return nil
}

// serverNegotiateExtensionsAndKeyShare records the record-size-limit
// negotiation, generates the server key share, and, for a certificate
// handshake, selects the certificate, signer, and signature scheme.
func (c *Conn) serverNegotiateExtensionsAndKeyShare(s *serverHandshakeState) error {
	if s.ch.hasRecordSizeLimit {
		c.recordRecordSizeLimit(c.config.RecordSizeLimit, effectiveRecordSizeLimit(s.ch.recordSizeLimit))
	}
	var err error
	if s.serverShare, s.shared, err = generateServerKeyShare(s.share.group, s.share.data, c.config.Rand); err != nil {
		return err
	}
	// TLS 1.3 applies a certificate-type selection in EncryptedExtensions to
	// every later Certificate message. A PSK handshake cannot send the initial
	// CertificateRequest, but can request PHA when the client offered it.
	if c.config.ClientAuth != tls.NoClientCert && (!s.usingPSK || s.ch.postHandshakeAuth) {
		typ, selectErr := selectCertificateType(s.ch.unknownExtensions[extClientCertificateType], c.config.ClientCertificateTypes, true, c.config.VerifyPeerRawPublicKey != nil)
		if selectErr == nil {
			c.clientCertificateType, s.clientTypeNegotiated = typ, true
		} else {
			alert, ok := protocolAlert(selectErr)
			if !s.usingPSK || !ok || alert != alertUnsupportedCertificate {
				return selectErr
			}
		}
	}
	if s.usingPSK {
		return nil
	}
	if err = requireCertificateSignatureAlgorithms(s.ch, false); err != nil {
		return err
	}
	if c.serverCertificateType, err = selectCertificateType(s.ch.unknownExtensions[extServerCertificateType], c.config.ServerCertificateTypes, len(c.config.Certificates) > 0 || c.config.GetCertificate != nil, c.config.RawPublicKeySigner != nil); err != nil {
		return err
	}
	if c.serverCertificateType == CertificateTypeRawPublicKey {
		s.cert, err = c.rawPublicKeyCertificate()
	} else {
		s.cert, err = c.serverCertificate(s.ch)
	}
	if err != nil {
		return err
	}
	signer, ok := s.cert.PrivateKey.(crypto.Signer)
	if !ok {
		return errors.New("dtls13: server private key does not implement crypto.Signer")
	}
	s.signer = signer
	certificateSchemes := s.ch.certificateSignatureSchemes
	if len(certificateSchemes) == 0 {
		certificateSchemes = s.ch.signatureSchemes
	}
	if c.serverCertificateType == CertificateTypeX509 {
		if validateErr := validateConfiguredCertificate(s.cert, certificateSchemes, true); validateErr != nil {
			return alertError(alertHandshakeFailure, validateErr)
		}
	}
	s.scheme, err = selectSignatureScheme(signer, s.ch.signatureSchemes)
	return err
}

// serverBuildServerHello fills the ServerHello, records connection-ID and
// return-routability negotiation, starts the handshake transcript (HRR-aware),
// writes the ECH acceptance confirmation into random, and encodes the message.
// DTLS 1.3 does not use TLS compatibility mode; the server MUST NOT echo
// legacy_session_id (RFC 9147 section 5).
func (c *Conn) serverBuildServerHello(s *serverHandshakeState) error {
	sh := &serverHello{cipherSuite: s.suite.id, keyShare: keyShareEntry{group: s.share.group, data: s.serverShare}}
	if s.ch.hasConnectionID && c.config.ConnectionID != nil {
		c.recordConnectionIDNegotiation(c.config.ConnectionID, s.ch.connectionID)
		sh.hasConnectionID = true
		sh.connectionID = append([]byte(nil), c.receiveConnectionID...)
	}
	if c.connectionIDNegotiated && s.ch.returnRoutability && !c.config.DisableReturnRoutabilityCheck {
		c.returnRoutabilityCheckNegotiated = true
		sh.returnRoutability = true
	}
	if s.usingPSK {
		selected := s.selectedPSKIdentity
		sh.selectedIdentity = &selected
	}
	if _, err := io.ReadFull(c.config.Rand, sh.random[:]); err != nil {
		return err
	}
	s.transcript = newTranscriptHash(s.suite.hash.New())
	if s.hrrUsed {
		_ = s.transcript.addHelloRetryRequest(s.initialHashTranscript.sumInto(s.transcriptDigest[:0]), s.hrrBody)
		_ = s.transcript.add(handshakeTypeClientHello, 1, s.helloBody)
		s.serverHelloSequence = 1
		if s.hrrFlight != nil {
			s.firstPlainRecordSequence = s.hrrFlight.nextRecordSequence()
		}
	} else {
		_ = s.transcript.add(handshakeTypeClientHello, 0, s.helloBody)
	}
	if s.echAccepted {
		clear(sh.random[24:])
		zeroBody, marshalErr := sh.marshal()
		if marshalErr != nil {
			return marshalErr
		}
		confirmationTranscript := s.transcript.clone()
		_ = confirmationTranscript.add(handshakeTypeServerHello, s.serverHelloSequence, zeroBody)
		copy(sh.random[24:], echAcceptConfirmation(s.suite, s.ch.random, "ech accept confirmation", confirmationTranscript.sumInto(s.transcriptDigest[:0])))
	}
	shBody, err := sh.marshal()
	if err != nil {
		return err
	}
	_ = s.transcript.add(handshakeTypeServerHello, s.serverHelloSequence, shBody)
	s.sh = sh
	s.shBody = shBody
	return nil
}

// serverDeriveHandshakeKeys derives the handshake traffic secrets and creates
// the epoch-2 ciphers, plus the epoch-1 cipher when early data was accepted on
// the no-cookie path.
func (c *Conn) serverDeriveHandshakeKeys(s *serverHandshakeState) error {
	s.schedule = newKeySchedule(s.suite, s.psk)
	if err := s.schedule.deriveHandshake(s.shared, s.transcript.sumInto(s.transcriptDigest[:0])); err != nil {
		return err
	}
	var err error
	if s.clientCipher, s.serverCipher, err = c.newHandshakeCiphers(s.suite, s.schedule.clientHandshakeTraffic, s.schedule.serverHandshakeTraffic); err != nil {
		return err
	}
	s.serverCipher.setPlaintextLimit(c.peerRecordSizeLimit)
	s.clientCipher.setPlaintextLimit(c.localRecordSizeLimit)
	if !s.hrrUsed && c.earlyAccepted {
		s.earlyCipher, err = newEarlyRecordCipher(s.suite, s.psk, s.helloBody, s.resumedSession.recordSizeLimit, c.config.ReplayWindow)
		if err != nil {
			return err
		}
	}
	return nil
}

// serverFlightConn returns the connection the server flight and the client's
// reply travel over: the amplification-limited one until the address is
// validated, which on the HRR path happened when the cookie was checked.
func (s *serverHandshakeState) serverFlightConn(c *Conn) net.Conn {
	if s.hrrUsed {
		return c.conn
	}
	return s.preValidationConn
}

// serverSendFlight builds EncryptedExtensions, an optional CertificateRequest,
// Certificate and CertificateVerify for a certificate handshake, and Finished;
// combines them with the plaintext ServerHello; and writes the flight.
// c.sendCipher is set here so that from this point the deferred alert sender in
// runHandshake encrypts alerts; moving this assignment would change which
// alerts go out in plaintext.
func (c *Conn) serverSendFlight(s *serverHandshakeState) error {
	plain, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeServerHello, sequence: s.serverHelloSequence, body: s.shBody}}, c.currentMTU(), 0, s.firstPlainRecordSequence)
	if err != nil {
		return err
	}
	s.ticketCount = 1
	if s.ch.ticketRequest.Enabled {
		s.ticketCount = s.ch.ticketRequest.NewSessionCount
		if s.resumed {
			s.ticketCount = s.ch.ticketRequest.ResumptionCount
		}
		s.ticketCount = min(s.ticketCount, c.config.MaxSessionTickets)
	}
	if c.config.SessionTicketsDisabled {
		s.ticketCount = 0
	}
	ee := &encryptedExtensions{recordSizeLimit: c.config.RecordSizeLimit, hasRecordSizeLimit: s.ch.hasRecordSizeLimit}
	if s.negotiated != "" || c.earlyAccepted || s.echRejected || s.ch.ticketRequest.Enabled {
		ee.extensions = make(map[uint16][]byte, 4)
	}
	if s.ch.ticketRequest.Enabled {
		ee.extensions[extTicketRequest] = []byte{s.ticketCount}
	}
	if s.negotiated != "" {
		if ee.extensions[extALPN], err = marshalALPN([]string{s.negotiated}); err != nil {
			return err
		}
	}
	if c.earlyAccepted {
		ee.extensions[extEarlyData] = nil
	}
	if s.echRejected && len(s.echKeys) > 0 {
		if ee.extensions[extECH], err = buildRetryConfigList(s.echKeys); err != nil {
			return err
		}
	}
	for _, selection := range []struct {
		typ     uint16
		value   CertificateType
		enabled bool
	}{
		{extServerCertificateType, c.serverCertificateType, !s.usingPSK},
		{extClientCertificateType, c.clientCertificateType, s.clientTypeNegotiated},
	} {
		if selection.enabled && s.ch.unknownExtensions[selection.typ] != nil {
			if ee.extensions == nil {
				ee.extensions = make(map[uint16][]byte)
			}
			ee.extensions[selection.typ] = []byte{byte(selection.value)}
		}
	}
	eeBody, err := ee.marshal()
	if err != nil {
		return err
	}
	s.serverSequence = 1 + s.serverSequenceOffset
	serverMessages := []handshakeMessage{{typ: handshakeTypeEncryptedExtensions, sequence: s.serverSequence, body: eeBody}}
	_ = s.transcript.add(handshakeTypeEncryptedExtensions, s.serverSequence, eeBody)
	s.serverSequence++
	if s.requestsClientCertificate(c) {
		request := c.newCertificateRequest(nil)
		s.clientStatusRequest = request.statusRequest
		if c.config.EnableCertificateCompression {
			s.clientCertificateCompression = &certificateCompressionZlibOffer
		}
		s.clientSignatureSchemes = append([]tls.SignatureScheme(nil), request.signatureSchemes...)
		s.clientCertificateSchemes = request.certificateSignatureSchemes
		s.clientCertificateOIDFilters = request.oidFilters
		if len(s.clientCertificateSchemes) == 0 {
			s.clientCertificateSchemes = append([]tls.SignatureScheme(nil), request.signatureSchemes...)
		}
		greaseExtension := uint16(0)
		if c.config.EnableGREASE {
			greaseExtension = greaseValue(s.sh.random[0])
		}
		requestBody, requestErr := request.marshalWithCertificateCompression(s.clientCertificateCompression, greaseExtension)
		if requestErr != nil {
			return requestErr
		}
		serverMessages = append(serverMessages, handshakeMessage{typ: handshakeTypeCertificateRequest, sequence: s.serverSequence, body: requestBody})
		_ = s.transcript.add(handshakeTypeCertificateRequest, s.serverSequence, requestBody)
		s.serverSequence++
	}
	if !s.usingPSK {
		certMsg := &certificateMessage{}
		for i, der := range s.cert.Certificate {
			entry := certificateEntry{data: der}
			if i == 0 && c.serverCertificateType == CertificateTypeX509 && c.config.serverCertificateEntryExtensions != nil {
				entry.extensions = make(map[uint16][]byte, len(c.config.serverCertificateEntryExtensions)+1)
				for typ, value := range c.config.serverCertificateEntryExtensions {
					entry.extensions[typ] = append([]byte(nil), value...)
				}
			}
			if i == 0 && c.serverCertificateType == CertificateTypeX509 && s.ch.statusRequest && len(s.cert.OCSPStaple) > 0 {
				if entry.extensions == nil {
					entry.extensions = make(map[uint16][]byte, 1)
				}
				if _, injected := entry.extensions[extStatusRequest]; !injected {
					entry.extensions[extStatusRequest], err = marshalOCSPResponse(s.cert.OCSPStaple)
					if err != nil {
						return err
					}
				}
			}
			certMsg.certificates = append(certMsg.certificates, entry)
		}
		messages, next, buildErr := c.buildCertificateMessages(certMsg, s.ch.certificateCompressionAlgorithms(), s.signer, s.scheme, true, s.transcript, s.serverSequence, s.transcriptDigest[:0])
		if buildErr != nil {
			return buildErr
		}
		serverMessages = append(serverMessages, messages...)
		s.serverSequence = next
	}
	finishedBody := s.schedule.finishedVerifyData(s.schedule.serverHandshakeTraffic, s.transcript.sumInto(s.transcriptDigest[:0]))
	serverMessages = append(serverMessages, handshakeMessage{typ: handshakeTypeFinished, sequence: s.serverSequence, body: finishedBody})
	_ = s.transcript.add(handshakeTypeFinished, s.serverSequence, finishedBody)
	protected, err := buildProtectedFlight(serverMessages, c.currentMTU(), s.serverCipher)
	if err != nil {
		return err
	}
	s.serverFlight = combineFlights(plain, protected)
	c.sendCipher = s.serverCipher
	return c.writeFlight(s.serverFlightConn(c), s.serverFlight)
}

// requestsClientCertificate reports whether this handshake sends a
// CertificateRequest: never on a PSK handshake, otherwise per Config.ClientAuth.
func (s *serverHandshakeState) requestsClientCertificate(c *Conn) bool {
	return !s.usingPSK && c.config.ClientAuth != tls.NoClientCert
}

// serverDeriveApplicationSecrets derives the application traffic secrets from
// the transcript through the server Finished.
func (c *Conn) serverDeriveApplicationSecrets(s *serverHandshakeState) error {
	return s.schedule.deriveApplication(s.transcript.sumInto(s.transcriptDigest[:0]))
}

// serverProcessClientFlight receives the client's final flight. The expected
// order is enforced by clientHandshakeStage; the per-message checks live in the
// server*Message helpers so the loop body stays a dispatch table.
func (c *Conn) serverProcessClientFlight(s *serverHandshakeState) error {
	s.clientFinalFlightStart = s.clientFinishedSeq
	s.inbox = newHandshakeInbox(s.clientFinalFlightStart, c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
	if s.resumedSession != nil {
		s.clientCerts = s.resumedSession.peerCertificates
		s.clientRawPublicKey = append([]byte(nil), s.resumedSession.peerRawPublicKey...)
		s.clientChains = s.resumedSession.verifiedChains
		s.clientAuthAt = s.resumedSession.clientAuthAt
	}
	stage := clientExpectFinished
	if s.requestsClientCertificate(c) {
		stage = clientExpectCertificate
	}
	sawClientCertificate := false
	verifiedClientSignature := false
	receiveConn := s.serverFlightConn(c)
	for stage != clientHandshakeComplete {
		messages, err := c.receiveHandshakeWithRetransmitOnEarly(receiveConn, s.inbox, s.clientCipher, s.serverFlight, s.earlyCipher, c.queueEarlyApplicationData, s.serverCipher)
		if err != nil {
			return err
		}
		s.clientRecords = append(s.clientRecords, recordNumber{epoch: 2, sequence: s.clientCipher.lastOpened})
		for index := 0; index < messages.len(); index++ {
			message := messages.at(index)
			if err = stage.accept(message.typ); err != nil {
				return err
			}
			switch message.typ {
			case handshakeTypeCertificate, handshakeTypeCompressedCertificate:
				sawClientCertificate = true
				var hasCertificates bool
				if hasCertificates, err = c.serverClientCertificate(s, message); err != nil {
					return err
				}
				if hasCertificates {
					stage = clientExpectCertificateVerify
				} else {
					stage = clientExpectFinished
				}
			case handshakeTypeCertificateVerify:
				if err = c.serverClientCertificateVerify(s, message); err != nil {
					return err
				}
				verifiedClientSignature = true
				stage = clientExpectFinished
			case handshakeTypeFinished:
				s.clientFinishedSeq = message.sequence
				if s.requestsClientCertificate(c) && !sawClientCertificate {
					return alertError(alertUnexpectedMessage, &ProtocolError{"client omitted Certificate message"})
				}
				required := !s.usingPSK && (c.config.ClientAuth == tls.RequireAnyClientCert || c.config.ClientAuth == tls.RequireAndVerifyClientCert)
				if required && len(s.clientCerts) == 0 && len(s.clientRawPublicKey) == 0 {
					return alertError(alertCertificateRequired, &ProtocolError{"client certificate is required"})
				}
				if !s.usingPSK && (len(s.clientCerts) > 0 || len(s.clientRawPublicKey) > 0) && !verifiedClientSignature {
					return alertError(alertUnexpectedMessage, &ProtocolError{"client omitted CertificateVerify"})
				}
				if err = c.serverClientFinished(s, message); err != nil {
					return err
				}
				stage = clientHandshakeComplete
			default:
				return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected client handshake message"})
			}
		}
	}
	return nil
}

// serverClientCertificate processes the client Certificate and reports whether
// it carried any certificates, which decides whether CertificateVerify follows.
func (c *Conn) serverClientCertificate(s *serverHandshakeState, message completedHandshake) (bool, error) {
	certMessage, err := parseCertificateHandshakeMessage(message.typ, message.body, s.clientCertificateCompression, c.config.MaxHandshakeMessage)
	if err != nil {
		return false, err
	}
	if err = validateCertificateMessageWithStatusRequest(certMessage, nil, s.clientStatusRequest); err != nil {
		return false, err
	}
	hasCertificates := len(certMessage.certificates) > 0
	if hasCertificates {
		s.ocspResponse = append([]byte(nil), certMessage.certificates[0].ocspResponse...)
		if c.clientCertificateType == CertificateTypeRawPublicKey {
			s.clientRawPublicKey, err = verifyRawPublicKeyMessage(c.config, certMessage)
			if err != nil {
				return false, err
			}
		} else {
			if s.clientCerts, s.clientChains, err = verifyClientCertificate(c.config, certMessage, s.clientCertificateSchemes); err != nil {
				return false, err
			}
			if err = matchCertificateOIDFilters(s.clientCerts[0], s.clientCertificateOIDFilters); err != nil {
				return false, alertError(alertUnsupportedCertificate, err)
			}
		}
	}
	_ = s.transcript.add(message.typ, message.sequence, message.body)
	return hasCertificates, nil
}

func (c *Conn) serverClientCertificateVerify(s *serverHandshakeState, message completedHandshake) error {
	if len(s.clientCerts) == 0 && len(s.clientRawPublicKey) == 0 {
		return alertError(alertUnexpectedMessage, &ProtocolError{"client CertificateVerify without certificate"})
	}
	cv, err := parseCertificateVerify(message.body)
	if err != nil {
		return err
	}
	if !slices.Contains(s.clientSignatureSchemes, cv.algorithm) {
		return alertError(alertIllegalParameter, &ProtocolError{"client selected an unoffered signature scheme"})
	}
	if err = verifyPeerCertificateSignature(s.clientCerts, s.clientRawPublicKey, cv.algorithm, s.transcript.sumInto(s.transcriptDigest[:0]), cv.signature, false); err != nil {
		return alertError(alertDecryptError, err)
	}
	_ = s.transcript.add(message.typ, message.sequence, message.body)
	return nil
}

func (c *Conn) serverClientFinished(s *serverHandshakeState, message completedHandshake) error {
	verify, err := parseFinished(message.body, s.suite.hash.Size())
	if err != nil {
		return err
	}
	if !s.schedule.verifyFinished(s.schedule.clientHandshakeTraffic, s.transcript.sumInto(s.transcriptDigest[:0]), verify) {
		return alertError(alertDecryptError, &ProtocolError{"client Finished verification failed"})
	}
	_ = s.transcript.add(message.typ, message.sequence, message.body)
	return nil
}

// serverSendFinalACK lifts the amplification limit on the no-cookie path, where
// a valid Finished proves address reachability, and acknowledges every record
// of the client's final flight.
func (c *Conn) serverSendFinalACK(s *serverHandshakeState) error {
	if !s.hrrUsed {
		s.amplification.validate()
	}
	ackRecords, _, err := buildACKRecords(s.clientRecords, c.currentMTU(), 0, s.serverCipher)
	if err != nil {
		return err
	}
	for _, wire := range ackRecords {
		if err = c.writeRecord(wire); err != nil {
			return err
		}
	}
	return nil
}

// serverFinalize installs the application keys, derives the resumption secret,
// promotes buffered early data, records the Finished ACK cipher and flight
// bounds, publishes the connection state, and notifies a wrapping transport
// that the handshake validated the peer.
func (c *Conn) serverFinalize(s *serverHandshakeState) error {
	if err := c.installApplicationKeysAt(s.suite, s.schedule.clientApplicationTraffic, s.schedule.serverApplicationTraffic, s.serverSequence+1); err != nil {
		return err
	}
	if err := c.finishHandshake(handshakeCompletion{
		suite:             s.suite,
		schedule:          s.schedule,
		transcript:        s.transcript,
		receiveCipher:     s.clientCipher,
		peerFlightStart:   s.clientFinalFlightStart,
		peerFlightEnd:     s.clientFinishedSeq,
		externalPSK:       s.externalPSK,
		resumed:           s.resumed,
		echAccepted:       s.echAccepted,
		negotiated:        s.negotiated,
		peerCerts:         s.clientCerts,
		peerRawPublicKey:  s.clientRawPublicKey,
		chains:            s.clientChains,
		ocspResponse:      s.ocspResponse,
		promoteEarly:      true,
		finishedACKCipher: s.serverCipher,
	}); err != nil {
		return err
	}
	if validated, ok := c.conn.(interface{ handshakeValidated() }); ok {
		validated.handshakeValidated()
	}
	return nil
}

// serverIssueTickets sends the NewSessionTicket flight.
func (c *Conn) serverIssueTickets(s *serverHandshakeState) error {
	if !s.usingPSK && (len(s.clientCerts) > 0 || len(s.clientRawPublicKey) > 0) {
		s.clientAuthAt = c.config.Time().Unix()
	}
	return c.sendNewSessionTickets(s.schedule, s.suite, s.ticketCount, s.ch.serverName, s.negotiated, s.clientAuthAt, s.clientCerts, s.clientChains, s.externalPSK)
}

func verifyClientCertificate(config *Config, message *certificateMessage, signatureSchemes []tls.SignatureScheme) ([]*x509.Certificate, [][]*x509.Certificate, error) {
	copyConfig := cloneConfig(config)
	if config.ClientAuth == tls.RequestClientCert || config.ClientAuth == tls.RequireAnyClientCert {
		copyConfig.InsecureSkipVerify = true
	}
	return verifyCertificateChain(copyConfig, message, false, signatureSchemes)
}

func (c *Conn) serverCertificate(ch *clientHello) (*tls.Certificate, error) {
	if c.config.GetCertificate != nil {
		certificate, err := c.config.GetCertificate(c.clientHelloInfo(ch))
		if err != nil {
			return nil, err
		}
		if certificate == nil {
			return nil, errors.New("dtls13: GetCertificate returned nil")
		}
		return certificate, nil
	}
	if len(c.config.Certificates) == 0 {
		return nil, errors.New("dtls13: server has no certificates")
	}
	if len(c.config.Certificates) > 1 {
		info := c.clientHelloInfo(ch)
		for i := range c.config.Certificates {
			if info.SupportsCertificate(&c.config.Certificates[i]) == nil {
				return &c.config.Certificates[i], nil
			}
		}
	}
	return &c.config.Certificates[0], nil
}
