package dtls13

import (
	"crypto"
	"crypto/tls"
	"crypto/x509"
)

// clientHandshakeStage is the server's view of the order the client's final
// flight must arrive in: Certificate and CertificateVerify only when a
// CertificateRequest was sent, then Finished. It mirrors serverHandshakeStage,
// which is the client's view of the server flight.
type clientHandshakeStage uint8

const (
	clientExpectCertificate clientHandshakeStage = iota
	clientExpectCertificateVerify
	clientExpectFinished
	clientHandshakeComplete
)

// accept checks that typ is permitted at the current stage. It does not
// advance the stage: whether Certificate is followed by CertificateVerify
// depends on whether the message carried a certificate, which only the caller
// knows after parsing it.
func (s clientHandshakeStage) accept(typ uint8) error {
	switch typ {
	case handshakeTypeCertificate, handshakeTypeCompressedCertificate:
		if s != clientExpectCertificate {
			return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected client Certificate"})
		}
	case handshakeTypeCertificateVerify:
		if s != clientExpectCertificateVerify {
			return alertError(alertUnexpectedMessage, &ProtocolError{"client CertificateVerify without certificate"})
		}
	case handshakeTypeFinished:
		if s != clientExpectFinished {
			return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected client Finished"})
		}
	}
	return nil
}

// newHandshakeCiphers creates the two epoch-2 record ciphers for a handshake
// and binds the negotiated connection IDs to them. receiveSecret and sendSecret
// are the traffic secrets for this endpoint's receive and send directions.
func (c *Conn) newHandshakeCiphers(suite *cipherSuite, receiveSecret, sendSecret []byte) (receive, send *recordCipher, err error) {
	if receive, err = newRecordCipher(suite, receiveSecret, 2, c.config.ReplayWindow); err != nil {
		return nil, nil, err
	}
	if send, err = newRecordCipher(suite, sendSecret, 2, c.config.ReplayWindow); err != nil {
		return nil, nil, err
	}
	if c.connectionIDNegotiated {
		if err = receive.setConnectionID(c.receiveConnectionID); err != nil {
			return nil, nil, err
		}
		if err = send.setConnectionID(c.sendConnectionID); err != nil {
			return nil, nil, err
		}
	}
	return receive, send, nil
}

// newEarlyRecordCipher derives the epoch-1 cipher from the PSK and the first
// ClientHello. The caller passes the ClientHello body that the early traffic
// secret must cover: the inner body when ECH is in use.
func newEarlyRecordCipher(suite *cipherSuite, psk, clientHelloBody []byte, recordSizeLimit uint16, replayWindow int) (*recordCipher, error) {
	var digest [maxSupportedHashSize]byte
	schedule := newKeySchedule(suite, psk)
	transcript := newTranscriptHash(suite.hash.New())
	_ = transcript.add(handshakeTypeClientHello, 0, clientHelloBody)
	cipher, err := newRecordCipher(suite, schedule.earlyTrafficSecret(transcript.sumInto(digest[:0])), 1, replayWindow)
	if err != nil {
		return nil, err
	}
	cipher.setPlaintextLimit(recordSizeLimit)
	return cipher, nil
}

// recordConnectionIDNegotiation stores the connection IDs each side will use.
// local is the ID this endpoint receives under; peer is the ID it sends with.
func (c *Conn) recordConnectionIDNegotiation(local, peer []byte) {
	c.connectionIDNegotiated = true
	c.sendConnectionID = append([]byte(nil), peer...)
	c.receiveConnectionID = append([]byte(nil), local...)
	c.localCIDUpdatesAllowed = len(c.receiveConnectionID) > 0
	c.peerCIDUpdatesAllowed = len(c.sendConnectionID) > 0
}

// recordRecordSizeLimit stores the RFC 8449 limits negotiated in each direction.
func (c *Conn) recordRecordSizeLimit(local, peer uint16) {
	c.recordSizeLimitNegotiated = true
	c.localRecordSizeLimit = local
	c.peerRecordSizeLimit = peer
}

// markEarlyDataRejected records that any early data already sent will not be
// accepted by the peer.
func (c *Conn) markEarlyDataRejected() {
	c.stopEarlyWrites()
	c.earlyStatus.Store(uint32(EarlyDataRejected))
}

// buildCertificateMessages encodes this endpoint's Certificate (compressed when
// the peer offered it and it helps) and, when signer is non-nil, a
// CertificateVerify over the transcript. Both are added to the transcript.
// It returns the messages and the next free message sequence.
func (c *Conn) buildCertificateMessages(certificate *certificateMessage, compression *certificateCompressionAlgorithms, signer crypto.Signer, scheme tls.SignatureScheme, server bool, transcript *transcriptHash, sequence uint16, digest []byte) ([]handshakeMessage, uint16, error) {
	certBody, err := certificate.marshal()
	if err != nil {
		return nil, sequence, err
	}
	certificateType, certificateBody, err := certificateHandshakeMessage(certBody, compression, c.config.EnableCertificateCompression, c.config.certificateCompressionCache())
	if err != nil {
		return nil, sequence, err
	}
	return c.buildCertificateBodyMessages(certificateType, certificateBody, signer, scheme, server, transcript, sequence, digest)
}

func (c *Conn) buildCertificateBodyMessages(certificateType uint8, certificateBody []byte, signer crypto.Signer, scheme tls.SignatureScheme, server bool, transcript *transcriptHash, sequence uint16, digest []byte) ([]handshakeMessage, uint16, error) {
	messages := []handshakeMessage{{typ: certificateType, sequence: sequence, body: certificateBody}}
	_ = transcript.add(certificateType, sequence, certificateBody)
	sequence++
	if signer == nil {
		return messages, sequence, nil
	}
	signature, err := signCertificateVerify(c.config.Rand, signer, scheme, transcript.sumInto(digest), server)
	if err != nil {
		return nil, sequence, err
	}
	cvBody, err := (&certificateVerifyMessage{algorithm: scheme, signature: signature}).marshal()
	if err != nil {
		return nil, sequence, err
	}
	messages = append(messages, handshakeMessage{typ: handshakeTypeCertificateVerify, sequence: sequence, body: cvBody})
	_ = transcript.add(handshakeTypeCertificateVerify, sequence, cvBody)
	sequence++
	return messages, sequence, nil
}

// handshakeCompletion is the input to finishHandshake: everything both
// handshakes publish once their final flight is acknowledged.
type handshakeCompletion struct {
	suite                   *cipherSuite
	schedule                *keySchedule
	transcript              *transcriptHash
	receiveCipher           *recordCipher
	peerFlightStart         uint16
	peerFlightEnd           uint16
	externalPSK             *externalPSKSelection
	resumed                 bool
	echAccepted             bool
	negotiated              string
	peerCerts               []*x509.Certificate
	peerRawPublicKey        []byte
	peerDelegatedCredential []byte
	chains                  [][]*x509.Certificate
	ocspResponse            []byte
	serverName              string
	// promoteEarly releases 0-RTT data buffered before the handshake completed;
	// only a server has any.
	promoteEarly bool
	// finishedACKCipher, when set, is retained so the server can keep
	// acknowledging a retransmitted client Finished (RFC 9147 section 5.8).
	finishedACKCipher *recordCipher
}

// finishHandshake performs the completion steps shared by both roles: derive
// the resumption secret, record the post-handshake transcript, install the
// receive epoch, note the bounds of the peer's final flight, build the
// exporter, and publish ConnectionState. Callers have already installed the
// application traffic keys.
func (c *Conn) finishHandshake(done handshakeCompletion) error {
	var digest [maxSupportedHashSize]byte
	if err := done.schedule.deriveResumption(done.transcript.sumInto(digest[:0])); err != nil {
		return err
	}
	c.writeMu.Lock()
	c.postHandshakeTranscript = done.transcript.clone()
	c.writeMu.Unlock()
	if err := c.receiveEpochs.install(done.receiveCipher); err != nil {
		return err
	}
	if done.promoteEarly {
		if err := c.promoteEarlyApplicationData(); err != nil {
			return err
		}
	}
	if done.finishedACKCipher != nil {
		c.writeMu.Lock()
		c.finishedACKCipher = done.finishedACKCipher
		c.writeMu.Unlock()
		c.finishedFlightStart = done.peerFlightStart
		c.finishedMessageSequence = done.peerFlightEnd
	}
	c.completedPeerFlightStart = done.peerFlightStart
	c.completedPeerFlightEnd = done.peerFlightEnd
	c.hasCompletedPeerFlight = true
	exporter := newExporter(done.suite, done.schedule.exporterMasterSecret)
	exporter.externalPSK = done.externalPSK
	c.mu.Lock()
	c.state = ConnectionState{
		Version: VersionDTLS13, HandshakeComplete: !c.isClient || !c.earlyIO, DidResume: done.resumed, ECHAccepted: done.echAccepted,
		CipherSuite: done.suite.id, NegotiatedProtocol: done.negotiated, ServerName: done.serverName,
		PeerCertificates: done.peerCerts, VerifiedChains: done.chains,
		PeerRawPublicKey:        append([]byte(nil), done.peerRawPublicKey...),
		PeerDelegatedCredential: append([]byte(nil), done.peerDelegatedCredential...),
		ServerCertificateType:   c.serverCertificateType, ClientCertificateType: c.clientCertificateType,
		OCSPResponse:      append([]byte(nil), done.ocspResponse...),
		LocalConnectionID: append([]byte(nil), c.receiveConnectionID...), PeerConnectionID: append([]byte(nil), c.sendConnectionID...),
		ReturnRoutabilityCheck:    c.returnRoutabilityCheckNegotiated,
		RecordSizeLimitNegotiated: c.recordSizeLimitNegotiated, LocalRecordSizeLimit: c.localRecordSizeLimit, PeerRecordSizeLimit: c.peerRecordSizeLimit,
		exporter: exporter,
	}
	c.mu.Unlock()
	return nil
}
