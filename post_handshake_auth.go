package dtls13

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"time"
)

type postHandshakeAuthState struct {
	context                          []byte
	transcript                       *transcriptHash
	inbox                            *handshakeInbox
	peerCertificates                 []*x509.Certificate
	verifiedChains                   [][]*x509.Certificate
	sawCertificate                   bool
	verifiedSignature                bool
	done                             chan error
	startSequence                    uint16
	hasStartSequence                 bool
	signatureSchemes                 []tls.SignatureScheme
	certificateSchemes               []tls.SignatureScheme
	certificateCompressionAlgorithms *certificateCompressionAlgorithms
	stage                            uint8
	responseEpoch                    uint64
	hasResponseEpoch                 bool
	firstResponseRecord              recordNumber
	lastResponseRecord               recordNumber
	pendingApplication               []pendingPostAuthApplication
	pendingApplicationBytes          int
}

type pendingPostAuthApplication struct {
	number  recordNumber
	content []byte
	from    net.Addr
}

const (
	postAuthExpectCertificate uint8 = iota
	postAuthExpectCertificateVerify
	postAuthExpectFinished
	postAuthComplete
)

func (c *Conn) newPostHandshakeAuthContext() ([]byte, error) {
	requestContext := make([]byte, 32)
	if _, err := io.ReadFull(c.config.Rand, requestContext); err != nil {
		return nil, err
	}
	counter := c.postHandshakeAuthCounter.Add(1)
	if counter == 0 {
		return nil, errors.New("dtls13: post-handshake authentication context counter exhausted")
	}
	binary.BigEndian.PutUint64(requestContext[len(requestContext)-8:], counter)
	return requestContext, nil
}

// RequestClientCertificate starts post-handshake client authentication and
// waits until the client's response has been received and verified. It is a
// server-only operation and performs the initial handshake first if necessary.
//
// The client must have advertised Config.PostHandshakeAuth. The server's
// Config.ClientAuth must request or require a certificate, and ClientCAs and
// VerifyPeerCertificate are applied according to that policy. At most one
// post-handshake authentication exchange may be active on a connection.
//
// Canceling ctx stops this call from waiting but does not retract a
// CertificateRequest that is already on the wire; the protocol exchange may
// continue in the background. A nil context is treated as context.Background.
func (c *Conn) RequestClientCertificate(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := c.HandshakeContext(ctx); err != nil {
		return err
	}
	if c.isClient {
		return errors.New("dtls13: only a server can request a client certificate")
	}
	if c.config.ClientAuth == tls.NoClientCert {
		return &ConfigError{"ClientAuth must request or require a client certificate"}
	}
	requestContext, err := c.newPostHandshakeAuthContext()
	if err != nil {
		return err
	}
	request := &certificateRequestMessage{requestContext: requestContext, signatureSchemes: defaultSignatureSchemes()}
	var certificateCompressionAlgorithms *certificateCompressionAlgorithms
	if c.config.EnableCertificateCompression {
		certificateCompressionAlgorithms = &certificateCompressionZlibOffer
	}
	body, err := request.marshalWithCertificateCompression(certificateCompressionAlgorithms)
	if err != nil {
		return err
	}

	c.writeMu.Lock()
	if !c.postHandshakeAuthOffered {
		c.writeMu.Unlock()
		return &ProtocolError{"client did not offer post_handshake_auth"}
	}
	if c.postHandshakeAuthState != nil {
		c.writeMu.Unlock()
		return errors.New("dtls13: post-handshake client authentication is already active")
	}
	if c.sendingTraffic == nil || c.receivingTraffic == nil || c.postHandshakeTranscript == nil {
		c.writeMu.Unlock()
		return errors.New("dtls13: traffic state is unavailable")
	}
	if err = c.sendingTraffic.canAllocateMessageSequences(1); err != nil {
		c.writeMu.Unlock()
		return err
	}
	sequence := c.sendingTraffic.messageSequence
	transcript := c.postHandshakeTranscript.clone()
	if err = transcript.add(handshakeTypeCertificateRequest, sequence, body); err != nil {
		c.writeMu.Unlock()
		return err
	}
	flight, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeCertificateRequest, sequence: sequence, body: body}}, c.currentMTU(), c.sendCipher)
	if err == nil {
		flight.setIntervals(c.config.FlightInterval, c.config.MaxFlightInterval)
		err = c.writeFlight(c.conn, flight)
	}
	if err != nil {
		c.writeMu.Unlock()
		return err
	}
	c.sendingTraffic.commitMessageSequences(1)
	state := &postHandshakeAuthState{
		context: requestContext, transcript: transcript,
		inbox:            newHandshakeInbox(c.finishedMessageSequence+1, c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes),
		done:             make(chan error, 1),
		signatureSchemes: append([]tls.SignatureScheme(nil), request.signatureSchemes...),
	}
	state.certificateSchemes = append([]tls.SignatureScheme(nil), request.certificateSignatureSchemes...)
	if len(state.certificateSchemes) == 0 {
		state.certificateSchemes = append([]tls.SignatureScheme(nil), request.signatureSchemes...)
	}
	state.certificateCompressionAlgorithms = certificateCompressionAlgorithms
	c.postHandshakeAuthState = state
	c.clientAuthRequestFlight = flight
	c.writeMu.Unlock()
	c.startClientAuthRetransmission(true)

	select {
	case err = <-state.done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Conn) processPostHandshakeCertificateRequest(sequence uint16, body []byte) error {
	request, certificateCompressionAlgorithms, err := parseCertificateRequestWithCompression(body)
	if err != nil {
		return err
	}
	if !c.isClient || !c.config.PostHandshakeAuth || len(request.requestContext) == 0 {
		return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected post-handshake CertificateRequest"})
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.hasClientAuthRequestSeq && sequence <= c.lastClientAuthRequestSeq {
		return nil
	}
	if c.clientAuthResponseFlight != nil && !c.clientAuthResponseFlight.complete() {
		return alertError(alertUnexpectedMessage, &ProtocolError{"overlapping post-handshake CertificateRequest"})
	}
	if c.postHandshakeTranscript == nil || c.sendingTraffic == nil {
		return &ProtocolError{"post-handshake authentication traffic state is unavailable"}
	}
	transcript := c.postHandshakeTranscript.clone()
	if err = transcript.add(handshakeTypeCertificateRequest, sequence, body); err != nil {
		return err
	}

	certificate := &certificateMessage{requestContext: append([]byte(nil), request.requestContext...)}
	var local *tls.Certificate
	if len(c.config.Certificates) > 0 {
		candidate := &c.config.Certificates[0]
		certificateSchemes := request.certificateSignatureSchemes
		if len(certificateSchemes) == 0 {
			certificateSchemes = request.signatureSchemes
		}
		if validateConfiguredCertificate(candidate, certificateSchemes, false) == nil {
			local = candidate
			for _, der := range local.Certificate {
				certificate.certificates = append(certificate.certificates, certificateEntry{data: der})
			}
		}
	}
	certificateBody, err := certificate.marshal()
	if err != nil {
		return err
	}
	messageCount := uint32(2)
	if local != nil && len(local.Certificate) > 0 {
		messageCount++
	}
	if err = c.sendingTraffic.canAllocateMessageSequences(messageCount); err != nil {
		return err
	}
	next := c.sendingTraffic.messageSequence
	certificateType, certificateBody, err := certificateHandshakeMessage(certificateBody, certificateCompressionAlgorithms, c.config.EnableCertificateCompression, nil)
	if err != nil {
		return err
	}
	messages := []handshakeMessage{{typ: certificateType, sequence: next, body: certificateBody}}
	if err = transcript.add(certificateType, next, certificateBody); err != nil {
		return err
	}
	next++
	if local != nil && len(local.Certificate) > 0 {
		signer, ok := local.PrivateKey.(crypto.Signer)
		if !ok {
			return errors.New("dtls13: client certificate private key is not a signer")
		}
		scheme, selectErr := selectSignatureScheme(signer, request.signatureSchemes)
		if selectErr != nil {
			return selectErr
		}
		signature, signErr := signCertificateVerify(c.config.Rand, signer, scheme, transcript.sum(), false)
		if signErr != nil {
			return signErr
		}
		verifyBody, marshalErr := (&certificateVerifyMessage{algorithm: scheme, signature: signature}).marshal()
		if marshalErr != nil {
			return marshalErr
		}
		messages = append(messages, handshakeMessage{typ: handshakeTypeCertificateVerify, sequence: next, body: verifyBody})
		if err = transcript.add(handshakeTypeCertificateVerify, next, verifyBody); err != nil {
			return err
		}
		next++
	}
	schedule := &keySchedule{suite: c.sendingTraffic.suite}
	finishedBody := schedule.finishedVerifyData(c.sendingTraffic.secret, transcript.sum())
	messages = append(messages, handshakeMessage{typ: handshakeTypeFinished, sequence: next, body: finishedBody})
	if err = transcript.add(handshakeTypeFinished, next, finishedBody); err != nil {
		return err
	}
	flight, err := buildProtectedFlight(messages, c.currentMTU(), c.sendCipher)
	if err == nil {
		flight.setIntervals(c.config.FlightInterval, c.config.MaxFlightInterval)
		err = c.writeFlight(c.conn, flight)
	}
	if err != nil {
		return err
	}
	c.sendingTraffic.commitMessageSequences(messageCount)
	c.postHandshakeTranscript = transcript
	c.clientAuthResponseFlight = flight
	c.lastClientAuthRequestSeq = sequence
	c.hasClientAuthRequestSeq = true
	c.startClientAuthRetransmission(false)
	return nil
}

func (c *Conn) processPostHandshakeAuthFragment(fragment handshakeFragment, number recordNumber) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	state := c.postHandshakeAuthState
	if !c.isClient && state == nil && c.hasCompletedClientAuth && fragment.messageSequence >= c.completedClientAuthStart && fragment.messageSequence <= c.completedClientAuthEnd {
		return c.writeACKLocked(number)
	}
	if c.isClient || state == nil {
		return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected post-handshake authentication response"})
	}
	if !state.hasStartSequence && (fragment.typ == handshakeTypeCertificate || fragment.typ == handshakeTypeCompressedCertificate) {
		if fragment.messageSequence <= c.finishedMessageSequence {
			return alertError(alertUnexpectedMessage, &ProtocolError{"post-handshake authentication reused a handshake message sequence"})
		}
		state.startSequence = fragment.messageSequence
		state.hasStartSequence = true
		state.inbox.expected = fragment.messageSequence
	}
	if state.hasResponseEpoch && state.responseEpoch != number.epoch {
		return alertError(alertUnexpectedMessage, &ProtocolError{"post-handshake authentication response changed epoch"})
	}
	if !state.hasResponseEpoch {
		state.responseEpoch = number.epoch
		state.hasResponseEpoch = true
		state.firstResponseRecord = number
		state.lastResponseRecord = number
	} else {
		if recordNumberLess(number, state.firstResponseRecord) {
			state.firstResponseRecord = number
		}
		if recordNumberLess(state.lastResponseRecord, number) {
			state.lastResponseRecord = number
		}
	}
	messages, err := state.inbox.addProtected(fragment, number.epoch)
	if err != nil {
		return err
	}
	if err = c.writeACKLocked(number); err != nil {
		return err
	}
	for _, message := range messages {
		if err = c.processPostHandshakeAuthMessageLocked(state, message); err != nil {
			return c.finishPostHandshakeAuthLocked(err)
		}
		if message.typ == handshakeTypeFinished {
			c.completedClientAuthStart = state.startSequence
			c.completedClientAuthEnd = message.sequence
			c.hasCompletedClientAuth = true
			return c.finishPostHandshakeAuthLocked(nil)
		}
	}
	return nil
}

func (c *Conn) processPostHandshakeAuthMessageLocked(state *postHandshakeAuthState, message completedHandshake) error {
	switch message.typ {
	case handshakeTypeCertificate, handshakeTypeCompressedCertificate:
		if state.stage != postAuthExpectCertificate {
			return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected post-handshake Certificate"})
		}
		certificate, err := parseCertificateHandshakeMessage(message.typ, message.body, state.certificateCompressionAlgorithms, c.config.MaxHandshakeMessage)
		if err != nil {
			return err
		}
		if !equalBytes(certificate.requestContext, state.context) {
			return &ProtocolError{"post-handshake Certificate context mismatch"}
		}
		if err = validateCertificateMessage(certificate, state.context); err != nil {
			return err
		}
		state.sawCertificate = true
		if len(certificate.certificates) > 0 {
			state.peerCertificates, state.verifiedChains, err = verifyClientCertificate(c.config, certificate, state.certificateSchemes)
			if err != nil {
				return err
			}
			state.stage = postAuthExpectCertificateVerify
		} else {
			state.stage = postAuthExpectFinished
		}
		return state.transcript.add(message.typ, message.sequence, message.body)
	case handshakeTypeCertificateVerify:
		if state.stage != postAuthExpectCertificateVerify || len(state.peerCertificates) == 0 {
			return &ProtocolError{"post-handshake CertificateVerify without certificate"}
		}
		verify, err := parseCertificateVerify(message.body)
		if err != nil {
			return err
		}
		offered := false
		for _, scheme := range state.signatureSchemes {
			offered = offered || scheme == verify.algorithm
		}
		if !offered {
			err = alertError(alertIllegalParameter, &ProtocolError{"post-handshake client selected an unoffered signature scheme"})
		}
		if err == nil {
			err = verifyCertificateVerify(state.peerCertificates[0].PublicKey, verify.algorithm, state.transcript.sum(), verify.signature, false)
			if err != nil {
				err = alertError(alertDecryptError, err)
			}
		}
		if err != nil {
			return err
		}
		state.verifiedSignature = true
		state.stage = postAuthExpectFinished
		return state.transcript.add(message.typ, message.sequence, message.body)
	case handshakeTypeFinished:
		if state.stage != postAuthExpectFinished {
			return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected post-handshake Finished"})
		}
		required := c.config.ClientAuth == tls.RequireAnyClientCert || c.config.ClientAuth == tls.RequireAndVerifyClientCert
		if !state.sawCertificate || (required && len(state.peerCertificates) == 0) || (len(state.peerCertificates) > 0 && !state.verifiedSignature) {
			return &ProtocolError{"post-handshake client authentication is incomplete"}
		}
		verify, err := parseFinished(message.body, c.receivingTraffic.suite.hash.Size())
		schedule := &keySchedule{suite: c.receivingTraffic.suite}
		secret, ok := c.receivingTraffic.secretForEpoch(state.responseEpoch)
		if err == nil && !ok {
			err = alertError(alertUnexpectedMessage, &ProtocolError{"post-handshake authentication epoch secret is unavailable"})
		}
		if err == nil && !schedule.verifyFinished(secret, state.transcript.sum(), verify) {
			err = alertError(alertDecryptError, &ProtocolError{"post-handshake client Finished verification failed"})
		}
		if err != nil {
			return err
		}
		if err = c.rememberProtectedHandshakeRangeLocked(state.firstResponseRecord, state.lastResponseRecord); err != nil {
			return err
		}
		if err = validatePostHandshakeAuthApplicationOrder(state); err != nil {
			return err
		}
		if err = state.transcript.add(message.typ, message.sequence, message.body); err != nil {
			return err
		}
		state.stage = postAuthComplete
		c.postHandshakeTranscript = state.transcript
		c.mu.Lock()
		c.state.PeerCertificates = append([]*x509.Certificate(nil), state.peerCertificates...)
		c.state.VerifiedChains = state.verifiedChains
		c.mu.Unlock()
		return nil
	default:
		return &ProtocolError{"unexpected post-handshake authentication message"}
	}
}

func validatePostHandshakeAuthApplicationOrder(state *postHandshakeAuthState) error {
	for _, application := range state.pendingApplication {
		if recordNumberLess(state.firstResponseRecord, application.number) && recordNumberLess(application.number, state.lastResponseRecord) {
			return alertError(alertUnexpectedMessage, &ProtocolError{"application data interleaved with post-handshake authentication response"})
		}
	}
	return nil
}

func (c *Conn) finishPostHandshakeAuthLocked(err error) error {
	state := c.postHandshakeAuthState
	if err == nil && state != nil {
		for _, application := range state.pendingApplication {
			if queueErr := c.queueApplicationData(application.content, application.from); queueErr != nil {
				err = queueErr
				break
			}
		}
	}
	c.postHandshakeAuthState = nil
	c.clientAuthRequestFlight = nil
	if state != nil {
		select {
		case state.done <- err:
		default:
		}
	}
	return err
}

func (c *Conn) bufferPostHandshakeAuthApplicationLocked(content []byte, number recordNumber, from net.Addr) (bool, error) {
	if c.isClient || c.postHandshakeAuthState == nil || !c.postHandshakeAuthState.hasResponseEpoch {
		return false, nil
	}
	state := c.postHandshakeAuthState
	if number.epoch != state.responseEpoch {
		return false, nil
	}
	if len(state.pendingApplication) >= c.maxPendingOrderingRecords() {
		return false, &ProtocolError{"too many application records buffered during post-handshake authentication"}
	}
	if len(content) > c.config.MaxBufferedApplicationData-state.pendingApplicationBytes {
		return false, &ProtocolError{"buffered post-handshake authentication application data limit exceeded"}
	}
	state.pendingApplication = append(state.pendingApplication, pendingPostAuthApplication{
		number: number, content: append([]byte(nil), content...), from: from,
	})
	state.pendingApplicationBytes += len(content)
	return true, nil
}

func (c *Conn) writeACKLocked(number recordNumber) error {
	var recordScratch [1][]byte
	records, _, err := buildACKRecordsInto(recordScratch[:0], []recordNumber{number}, c.currentMTU(), 0, c.sendCipher)
	if err != nil {
		return err
	}
	for _, wire := range records {
		if err = c.writeRecord(wire); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conn) startClientAuthRetransmission(request bool) {
	interval := c.flightInterval()
	if interval <= 0 {
		interval = time.Second
	}
	max := c.config.MaxFlightInterval
	if max < interval {
		max = interval
	}
	go func() {
		timer := time.NewTimer(interval)
		defer timer.Stop()
		timeoutCount := 0
		for range timer.C {
			c.writeMu.Lock()
			flight := c.clientAuthResponseFlight
			if request {
				flight = c.clientAuthRequestFlight
			}
			if flight == nil || flight.complete() {
				c.writeMu.Unlock()
				return
			}
			timeoutCount++
			resized, err := c.prepareFlightRetransmission(flight, timeoutCount)
			if err == nil && resized {
				err = c.writeFlight(c.conn, flight)
			} else if err == nil {
				err = c.retransmitFlight(c.conn, flight)
			}
			c.writeMu.Unlock()
			if err != nil {
				c.failConnection(err)
				return
			}
			if interval < max {
				interval *= 2
				if interval > max {
					interval = max
				}
			}
			timer.Reset(interval)
		}
	}()
}
