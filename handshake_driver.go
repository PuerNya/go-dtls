package dtls13

import (
	"context"
	"crypto"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"slices"
	"time"
)

const (
	handshakeTypeClientHello    uint8 = 1
	handshakeTypeServerHello    uint8 = 2
	handshakeTypeEndOfEarlyData uint8 = 5
)

type serverHandshakeStage uint8

const (
	serverExpectEncryptedExtensions serverHandshakeStage = iota
	serverExpectCertificateRequestOrCertificate
	serverExpectCertificate
	serverExpectCertificateVerify
	serverExpectFinished
	serverHandshakeComplete
)

func (s *serverHandshakeStage) accept(typ uint8, resumed bool) error {
	switch *s {
	case serverExpectEncryptedExtensions:
		if typ != handshakeTypeEncryptedExtensions {
			break
		}
		if resumed {
			*s = serverExpectFinished
		} else {
			*s = serverExpectCertificateRequestOrCertificate
		}
		return nil
	case serverExpectCertificateRequestOrCertificate:
		if typ == handshakeTypeCertificateRequest {
			*s = serverExpectCertificate
			return nil
		}
		if typ == handshakeTypeCertificate || typ == handshakeTypeCompressedCertificate {
			*s = serverExpectCertificateVerify
			return nil
		}
	case serverExpectCertificate:
		if typ == handshakeTypeCertificate || typ == handshakeTypeCompressedCertificate {
			*s = serverExpectCertificateVerify
			return nil
		}
	case serverExpectCertificateVerify:
		if typ == handshakeTypeCertificateVerify {
			*s = serverExpectFinished
			return nil
		}
	case serverExpectFinished:
		if typ == handshakeTypeFinished {
			*s = serverHandshakeComplete
			return nil
		}
	}
	return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected server handshake message order"})
}

func (c *Conn) runHandshake(ctx context.Context) (result error) {
	c.localRecordSizeLimit = defaultRecordSizeLimit
	c.peerRecordSizeLimit = defaultRecordSizeLimit
	deadline := c.config.Time().Add(c.config.HandshakeTimeout)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	if err := c.conn.SetDeadline(deadline); err != nil {
		return err
	}
	c.handshakeDeadline = deadline
	defer func() { _ = c.conn.SetDeadline(time.Time{}) }()
	defer func() {
		if description, ok := outboundAlert(result); ok {
			if c.sendCipher != nil {
				c.sendFatalAlert(description)
			} else {
				var local *localAlertError
				if errors.As(result, &local) {
					c.sendPlainFatalAlert(description)
				}
			}
		}
		if result != nil && !errors.Is(result, io.EOF) {
			c.clearTrafficSecrets(result)
		}
	}()
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			_ = c.conn.SetDeadline(time.Now())
		case <-done:
		}
	}()
	if c.isClient {
		result = c.clientHandshake()
		return result
	}
	result = c.serverHandshake()
	return result
}

func (c *Conn) writeFlight(conn io.Writer, f *flight) error {
	var storage [10][]byte
	for {
		var err error
		records := f.nextUnsentWire(10, storage[:0])
		if len(records) > 0 {
			f.noteSend(c.config.Time(), false)
		}
		for _, record := range records {
			if _, err = conn.Write(record); err != nil {
				break
			}
		}
		if err == nil || !isMessageTooLong(err) {
			return err
		}
		mtu, reduced := c.reducePathMTU()
		if !reduced {
			return normalizeDatagramWriteError(err, c.RemoteAddr())
		}
		if err = f.resize(mtu); err != nil {
			return err
		}
	}
}

func (c *Conn) retransmitFlight(conn io.Writer, f *flight) error {
	var storage [10][]byte
	records := f.retransmitWire(10, storage[:0])
	if len(records) > 0 {
		f.noteSend(c.config.Time(), true)
	}
	for _, record := range records {
		if _, err := conn.Write(record); err != nil {
			return normalizeDatagramWriteError(err, c.RemoteAddr())
		}
	}
	return nil
}

func (c *Conn) retransmitPartialFlight(conn io.Writer, f *flight) error {
	if f.claimPartialRetransmission() {
		if err := f.refreshPending(); err != nil {
			return err
		}
		if err := c.retransmitFlight(conn, f); err != nil {
			return err
		}
	}
	return c.writeFlight(conn, f)
}

func (c *Conn) prepareFlightRetransmission(f *flight, timeoutCount int) (bool, error) {
	if timeoutCount > 0 && timeoutCount%3 == 0 {
		if mtu, reduced := c.reducePathMTU(); reduced {
			return true, f.resize(mtu)
		}
	}
	return false, f.refreshPending()
}

func (c *Conn) receiveHandshakeWithRetransmit(inbox *handshakeInbox, cipher *recordCipher, outgoing *flight) (completedHandshakeBatch, error) {
	return c.receiveHandshakeWithRetransmitOn(c.conn, inbox, cipher, outgoing)
}

func (c *Conn) receiveHandshakeWithRetransmitOn(conn net.Conn, inbox *handshakeInbox, cipher *recordCipher, outgoing *flight) (completedHandshakeBatch, error) {
	return c.receiveHandshakeWithRetransmitOnEarly(conn, inbox, cipher, outgoing, nil, nil, nil)
}

// receiveHandshakeWithRetransmitOnEarly is the handshake receive loop with an
// optional epoch-1 cipher. DTLS 1.3 permits early application records to be
// interleaved with the protected handshake flight, so they must be consumed by
// the same datagram reader rather than by a second goroutine reading conn.
func (c *Conn) receiveHandshakeWithRetransmitOnEarly(conn net.Conn, inbox *handshakeInbox, cipher *recordCipher, outgoing *flight, early *recordCipher, onEarly func([]byte) error, ackCipher *recordCipher) (completedHandshakeBatch, error) {
	interval := c.flightInterval()
	if interval <= 0 {
		interval = time.Second
	}
	max := c.config.MaxFlightInterval
	if max < interval {
		max = interval
	}
	timeoutCount := 0
	for {
		next := c.config.Time().Add(interval)
		if next.After(c.handshakeDeadline) {
			next = c.handshakeDeadline
		}
		if err := conn.SetReadDeadline(next); err != nil {
			return completedHandshakeBatch{}, err
		}
		messages, err := receiveHandshakeMessageWithEarlyBatch(conn, inbox, cipher, early, onEarly, outgoing, ackCipher, c.currentMTU(), c)
		if err == nil {
			_ = conn.SetReadDeadline(c.handshakeDeadline)
			return messages, nil
		}
		networkErr, ok := err.(net.Error)
		if !ok || !networkErr.Timeout() || !c.config.Time().Before(c.handshakeDeadline) {
			return completedHandshakeBatch{}, err
		}
		if outgoing != nil {
			timeoutCount++
			resized, prepareErr := c.prepareFlightRetransmission(outgoing, timeoutCount)
			if prepareErr != nil {
				return completedHandshakeBatch{}, prepareErr
			}
			if resized {
				err = c.writeFlight(conn, outgoing)
			} else {
				err = c.retransmitFlight(conn, outgoing)
			}
			if err != nil {
				return completedHandshakeBatch{}, err
			}
		}
		if interval < max {
			interval *= 2
			if interval > max {
				interval = max
			}
		}
	}
}
func receiveHandshakeMessage(conn net.Conn, inbox *handshakeInbox, cipher *recordCipher) ([]completedHandshake, error) {
	batch, err := receiveHandshakeMessageBatch(conn, inbox, cipher)
	if err != nil {
		return nil, err
	}
	return batch.slice(), nil
}

func receiveHandshakeMessageBatch(conn net.Conn, inbox *handshakeInbox, cipher *recordCipher) (completedHandshakeBatch, error) {
	return receiveHandshakeMessageWithEarlyBatch(conn, inbox, cipher, nil, nil, nil, nil, 0, nil)
}

func (c *Conn) receiveSecondClientHello(conn net.Conn, inbox *handshakeInbox, hrr *flight) (completedHandshakeBatch, error) {
	buffer := acquireDatagramBuffer()
	defer releaseDatagramBuffer(buffer)
	for {
		datagram := buffer[:]
		n, err := conn.Read(datagram)
		if err != nil {
			return completedHandshakeBatch{}, err
		}
		var recordScratch [1]record
		records, err := parsePlainRecordsViewInto(datagram[:n], recordScratch[:0])
		if err != nil {
			continue
		}
		for _, record := range records {
			if record.typ != recordTypeHandshake {
				continue
			}
			var fragmentScratch [1]handshakeFragment
			fragments, parseErr := parseHandshakeFragmentsViewInto(record.payload, fragmentScratch[:0])
			if parseErr != nil {
				return completedHandshakeBatch{}, parseErr
			}
			var delivered completedHandshakeBatch
			retransmitHRR := false
			for _, fragment := range fragments {
				if fragment.typ == handshakeTypeClientHello && fragment.messageSequence < inbox.expected {
					retransmitHRR = true
					continue
				}
				if addErr := inbox.addBatch(&delivered, fragment); addErr != nil {
					return completedHandshakeBatch{}, addErr
				}
			}
			if retransmitHRR {
				if err = hrr.refreshPending(); err != nil {
					return completedHandshakeBatch{}, err
				}
				if err = c.retransmitFlight(conn, hrr); err != nil {
					return completedHandshakeBatch{}, err
				}
			}
			if delivered.len() > 0 {
				return delivered, nil
			}
		}
	}
}

func receiveHandshakeMessageWithEarly(conn net.Conn, inbox *handshakeInbox, cipher *recordCipher, outgoing *flight, ackCipher *recordCipher, mtu int, owner *Conn) ([]completedHandshake, error) {
	batch, err := receiveHandshakeMessageWithEarlyBatch(conn, inbox, cipher, nil, nil, outgoing, ackCipher, mtu, owner)
	if err != nil {
		return nil, err
	}
	return batch.slice(), nil
}

func receiveHandshakeMessageWithEarlyBatch(conn net.Conn, inbox *handshakeInbox, cipher, early *recordCipher, onEarly func([]byte) error, outgoing *flight, ackCipher *recordCipher, mtu int, owner *Conn) (completedHandshakeBatch, error) {
	buffer := acquireDatagramBuffer()
	defer releaseDatagramBuffer(buffer)
	for {
		datagram := buffer[:]
		n, err := conn.Read(datagram)
		if err != nil {
			return completedHandshakeBatch{}, err
		}
		datagram = datagram[:n]
		for len(datagram) > 0 {
			var payload []byte
			var typ uint8
			var consumed int
			var recordEpoch uint64
			if cipher == nil {
				if isUnifiedRecord(datagram) && owner != nil {
					sequence := owner.plainSendSequence.Add(1) - 1
					var ackScratch [1][]byte
					acks, _, ackErr := buildACKRecordsInto(ackScratch[:0], nil, mtu, sequence, nil)
					if ackErr != nil {
						return completedHandshakeBatch{}, ackErr
					}
					for _, wire := range acks {
						if _, ackErr = conn.Write(wire); ackErr != nil {
							return completedHandshakeBatch{}, ackErr
						}
					}
					break
				}
				var recordScratch [1]record
				records, parseErr := parsePlainRecordsViewInto(datagram, recordScratch[:0])
				if parseErr != nil {
					break
				}
				if len(records) == 0 {
					break
				}
				recordWireLen := plainRecordHeaderLen + len(records[0].payload)
				typ = records[0].typ
				payload = records[0].payload
				consumed = recordWireLen
			} else {
				var openErr error
				if datagram[0] == recordTypeACK {
					var recordScratch [1]record
					records, parseErr := parsePlainRecordsViewInto(datagram, recordScratch[:0])
					if parseErr == nil && len(records) > 0 && records[0].typ == recordTypeACK {
						consumed = plainRecordHeaderLen + len(records[0].payload)
						numbers, ackErr := parseACK(records[0].payload)
						if ackErr != nil {
							return completedHandshakeBatch{}, ackErr
						}
						if len(numbers) == 0 && outgoing != nil && owner != nil {
							if ackErr = outgoing.refreshPending(); ackErr != nil {
								return completedHandshakeBatch{}, ackErr
							}
							if ackErr = owner.retransmitFlight(conn, outgoing); ackErr != nil {
								return completedHandshakeBatch{}, ackErr
							}
						}
						datagram = datagram[consumed:]
						continue
					}
				}
				// Epoch bits are unambiguous during the initial handshake. Try
				// epoch 1 first for early records and epoch 2 with the normal
				// handshake cipher otherwise.
				if early != nil && len(datagram) > 0 && datagram[0]&unifiedHeaderEpochMask == byte(early.epoch&unifiedHeaderEpochMask) {
					payload, typ, consumed, openErr = early.openInPlace(datagram)
					recordEpoch = early.epoch
					if openErr == nil && typ == recordTypeApplicationData && onEarly != nil {
						if callbackErr := onEarly(payload); callbackErr != nil {
							return completedHandshakeBatch{}, callbackErr
						}
					}
				} else {
					payload, typ, consumed, openErr = cipher.openInPlace(datagram)
					recordEpoch = cipher.epoch
				}
				if openErr != nil {
					if fatalErr := protectedRecordReceiveError(openErr); fatalErr != nil {
						return completedHandshakeBatch{}, fatalErr
					}
					break
				}
			}
			datagram = datagram[consumed:]
			if typ == recordTypeAlert {
				alert, parseErr := parseAlert(payload)
				if parseErr != nil {
					return completedHandshakeBatch{}, parseErr
				}
				if alert.isUserCanceled() {
					continue
				}
				if alert.isCloseNotify() {
					return completedHandshakeBatch{}, io.EOF
				}
				return completedHandshakeBatch{}, AlertError(alert.description)
			}
			if typ == recordTypeACK && outgoing != nil {
				var ackScratch [1]recordNumber
				numbers, parseErr := parseACKInto(payload, ackScratch[:0])
				if parseErr != nil {
					return completedHandshakeBatch{}, parseErr
				}
				if parseErr = validateACKEpoch(numbers, recordEpoch); parseErr != nil {
					return completedHandshakeBatch{}, parseErr
				}
				outgoing.ack(numbers)
				if owner != nil && !outgoing.complete() {
					if sendErr := owner.retransmitPartialFlight(conn, outgoing); sendErr != nil {
						return completedHandshakeBatch{}, sendErr
					}
				}
				continue
			}
			if typ != recordTypeHandshake {
				continue
			}
			var fragmentScratch [1]handshakeFragment
			fragments, parseErr := parseHandshakeFragmentsViewInto(payload, fragmentScratch[:0])
			if parseErr != nil {
				return completedHandshakeBatch{}, parseErr
			}
			var delivered completedHandshakeBatch
			acknowledge := true
			peerRetransmitted := false
			for _, fragment := range fragments {
				if fragment.messageSequence < inbox.expected {
					acknowledge = false
					peerRetransmitted = true
				}
				if addErr := inbox.addProtectedBatch(&delivered, fragment, recordEpoch); addErr != nil {
					return completedHandshakeBatch{}, addErr
				}
			}
			if peerRetransmitted && outgoing != nil && !outgoing.hasAcknowledgedRecord() && owner != nil {
				if retransmitErr := outgoing.refreshPending(); retransmitErr != nil {
					return completedHandshakeBatch{}, retransmitErr
				}
				if retransmitErr := owner.retransmitFlight(conn, outgoing); retransmitErr != nil {
					return completedHandshakeBatch{}, retransmitErr
				}
			}
			if cipher != nil && ackCipher != nil && acknowledge {
				number := recordNumber{epoch: cipher.epoch, sequence: cipher.lastOpened}
				var ackScratch [1][]byte
				acks, _, ackErr := buildACKRecordsInto(ackScratch[:0], []recordNumber{number}, mtu, 0, ackCipher)
				if ackErr != nil {
					return completedHandshakeBatch{}, ackErr
				}
				for _, wire := range acks {
					if _, ackErr = conn.Write(wire); ackErr != nil {
						return completedHandshakeBatch{}, ackErr
					}
				}
			}
			if delivered.len() > 0 {
				return delivered, nil
			}
		}
	}
}

func receiveACKRecord(conn net.Conn, dst []recordNumber, ciphers ...*recordCipher) ([]recordNumber, error) {
	buffer := acquireDatagramBuffer()
	defer releaseDatagramBuffer(buffer)

	for {
		datagram := buffer[:]
		n, err := conn.Read(datagram)
		if err != nil {
			return nil, err
		}
		datagram = datagram[:n]
		for len(datagram) > 0 && isUnifiedRecord(datagram) {
			lastCipher := -1
			for i, cipher := range ciphers {
				if recordCipherMatchesUnifiedEpoch(cipher, datagram[0]) {
					lastCipher = i
				}
			}
			if lastCipher < 0 {
				break
			}
			consumed := 0
			for i, cipher := range ciphers {
				if !recordCipherMatchesUnifiedEpoch(cipher, datagram[0]) {
					continue
				}
				var content []byte
				var typ uint8
				var openErr error
				if i == lastCipher {
					content, typ, consumed, openErr = cipher.openInPlace(datagram)
				} else {
					content, typ, consumed, openErr = cipher.open(datagram)
				}
				if openErr != nil {
					consumed = 0
					if fatalErr := protectedRecordReceiveError(openErr); fatalErr != nil {
						return nil, fatalErr
					}
					continue
				}
				switch typ {
				case recordTypeACK:
					numbers, parseErr := parseACKInto(content, dst)
					if parseErr != nil {
						return nil, parseErr
					}
					if parseErr = validateACKEpoch(numbers, cipher.epoch); parseErr != nil {
						return nil, parseErr
					}
					return numbers, nil
				case recordTypeAlert:
					alert, parseErr := parseAlert(content)
					if parseErr != nil {
						return nil, parseErr
					}
					if alert.isCloseNotify() {
						return nil, io.EOF
					}
					if !alert.isUserCanceled() {
						return nil, AlertError(alert.description)
					}
				}
				break
			}
			if consumed == 0 {
				break
			}
			datagram = datagram[consumed:]
		}
	}
}

func (c *Conn) receiveACKWithRetransmit(outgoing *flight, ciphers ...*recordCipher) ([]recordNumber, error) {
	interval := c.flightInterval()
	if interval <= 0 {
		interval = time.Second
	}
	max := c.config.MaxFlightInterval
	if max < interval {
		max = interval
	}
	var acknowledged []recordNumber
	var ackScratch [1]recordNumber
	timeoutCount := 0
	for {
		next := c.config.Time().Add(interval)
		if next.After(c.handshakeDeadline) {
			next = c.handshakeDeadline
		}
		if err := c.conn.SetReadDeadline(next); err != nil {
			return nil, err
		}
		numbers, err := receiveACKRecord(c.conn, ackScratch[:0], ciphers...)
		if err == nil {
			acknowledged = append(acknowledged, numbers...)
			outgoing.ack(numbers)
			if outgoing.complete() {
				c.observeFlightRTT(outgoing)
				return canonicalRecordNumbers(acknowledged), nil
			}
			if err = c.retransmitPartialFlight(c.conn, outgoing); err != nil {
				return nil, err
			}
			continue
		}
		networkErr, ok := err.(net.Error)
		if !ok || !networkErr.Timeout() || !c.config.Time().Before(c.handshakeDeadline) {
			return nil, err
		}
		timeoutCount++
		resized, prepareErr := c.prepareFlightRetransmission(outgoing, timeoutCount)
		if prepareErr != nil {
			return nil, prepareErr
		}
		if resized {
			err = c.writeFlight(c.conn, outgoing)
		} else {
			err = c.retransmitFlight(c.conn, outgoing)
		}
		if err != nil {
			return nil, err
		}
		if interval < max {
			interval *= 2
			if interval > max {
				interval = max
			}
		}
	}
}

func (c *Conn) clientHandshake() error {
	var transcriptDigest [maxSupportedHashSize]byte
	key, err := generateEphemeralKey(c.config.CurvePreferences[0], c.config.Rand)
	if err != nil {
		return alertError(alertInternalError, err)
	}
	keyShares := []keyShareEntry{{group: key.groupID(), data: key.publicBytes()}}
	if group, public, ok := key.fallbackPublicBytes(); ok && slices.Contains(c.config.CurvePreferences, group) {
		keyShares = append(keyShares, keyShareEntry{group: group, data: public})
	}
	hello := &clientHello{cipherSuites: append([]uint16(nil), c.config.CipherSuites...), keyShares: keyShares, supportedGroups: c.config.CurvePreferences, signatureSchemes: defaultSignatureSchemes(), serverName: c.config.ServerName, alpn: c.config.NextProtos, postHandshakeAuth: c.config.PostHandshakeAuth, recordSizeLimit: c.config.RecordSizeLimit, hasRecordSizeLimit: true}
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
	var ech *echClientContext
	if c.config.EncryptedClientHelloConfigList != nil {
		ech, err = newECHClientContext(c.config.EncryptedClientHelloConfigList)
		if err != nil {
			return err
		}
		hello.setEncryptedClientHello([]byte{echInnerType})
	}
	greaseECH := ech == nil && c.config.EncryptedClientHelloGrease
	clientSession, sessionSuite := usableClientSession(c.config, c.conn)
	var pskOffers []clientPSKOffer
	var externalOffers []clientPSKOffer
	if len(c.config.ExternalPSKs) > 0 {
		externalOffers = configuredExternalPSKOffers(c.config)
	}
	if len(externalOffers) > 0 {
		if clientSession != nil {
			pskOffers = append(pskOffers, clientPSKOffer{
				identity: clientSession.ticket, psk: clientSession.psk, suite: sessionSuite,
				binderLabel: labelResumptionBinder, age: obfuscatedTicketAge(clientSession, c.config.Time()),
				session: clientSession, external: clientSession.externalPSK,
			})
		}
		pskOffers = append(pskOffers, externalOffers...)
	} else if clientSession != nil {
		hello.pskIdentity = append([]byte(nil), clientSession.ticket...)
		hello.obfuscatedAge = obfuscatedTicketAge(clientSession, c.config.Time())
	}
	c.earlyMu.Lock()
	queuedEarly := append([]byte(nil), c.earlyPending...)
	if len(queuedEarly) > 0 && (clientSession == nil || clientSession.maxEarlyData == 0) {
		c.earlyPending = nil
	}
	c.earlyMu.Unlock()
	if clientSession != nil {
		if slices.Contains(hello.cipherSuites, sessionSuite.id) && hello.cipherSuites[0] != sessionSuite.id {
			prioritized := []uint16{sessionSuite.id}
			for _, id := range hello.cipherSuites {
				if id != sessionSuite.id {
					prioritized = append(prioritized, id)
				}
			}
			hello.cipherSuites = prioritized
		}
		if len(queuedEarly) > 0 && clientSession.maxEarlyData > 0 {
			hello.earlyData = true
		}
	}
	var helloBody []byte
	if len(pskOffers) > 0 {
		helloBody, err = marshalClientHelloWithPSKOffers(hello, pskOffers, nil, nil)
	} else if clientSession != nil {
		helloBody, err = marshalClientHelloWithPSKBinder(hello, sessionSuite, clientSession.psk, nil, nil)
	} else {
		helloBody, err = hello.marshal()
	}
	if err != nil {
		return err
	}
	if greaseECH {
		grease, greaseErr := generateGREASEECH(hello, c.config.Rand)
		err = greaseErr
		if err != nil {
			return err
		}
		hello.setEncryptedClientHello(grease)
		if len(pskOffers) > 0 {
			helloBody, err = marshalClientHelloWithPSKOffers(hello, pskOffers, nil, nil)
		} else if clientSession != nil {
			helloBody, err = marshalClientHelloWithPSKBinder(hello, sessionSuite, clientSession.psk, nil, nil)
		} else {
			helloBody, err = hello.marshal()
		}
		if err != nil {
			return err
		}
	}
	if ech != nil {
		ech.innerHello = hello
		ech.innerBody = helloBody
		outer, outerErr := makeECHOuter(hello, ech.config, c.config.Rand)
		if outerErr != nil {
			return outerErr
		}
		helloBody, err = computeOuterECH(outer, hello, ech, true)
		if err != nil {
			return err
		}
		ech.outerHello = outer
		hello = outer
	}
	out, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeClientHello, sequence: 0, body: helloBody}}, c.currentMTU(), 0, 0)
	if err != nil {
		return err
	}
	c.plainSendSequence.Store(out.nextRecordSequence())
	if err = c.writeFlight(c.conn, out); err != nil {
		return err
	}
	if hello.earlyData {
		// Epoch 1 is derived from the PSK and the complete first ClientHello,
		// including its binder, and can be sent before ServerHello arrives.
		earlySchedule := newKeySchedule(sessionSuite, clientSession.psk)
		transcriptHash := newTranscriptHash(sessionSuite.hash.New())
		earlyHelloBody := helloBody
		if ech != nil {
			earlyHelloBody = ech.innerBody
		}
		_ = transcriptHash.add(handshakeTypeClientHello, 0, earlyHelloBody)
		earlyCipher, cipherErr := newRecordCipher(sessionSuite, earlySchedule.earlyTrafficSecret(transcriptHash.sumInto(transcriptDigest[:0])), 1, c.config.ReplayWindow)
		if cipherErr != nil {
			return cipherErr
		}
		earlyCipher.setPlaintextLimit(clientSession.recordSizeLimit)
		if len(queuedEarly) > int(clientSession.maxEarlyData) {
			c.earlyMu.Lock()
			c.earlyPending = nil
			c.earlyRejected = true
			c.earlyMu.Unlock()
		} else {
			c.writeMu.Lock()
			maxRecord := c.maxApplicationDatagramForCipher(earlyCipher)
			if maxRecord < 1 {
				c.writeMu.Unlock()
				return &ConfigError{"MTU is too small for early application data"}
			}
			if len(queuedEarly) > maxRecord {
				c.writeMu.Unlock()
				c.earlyMu.Lock()
				c.earlyPending = nil
				c.earlyMu.Unlock()
				return datagramTooLargeError(c.conn.RemoteAddr())
			}
			wire, sealErr := earlyCipher.seal(recordTypeApplicationData, queuedEarly)
			if sealErr != nil {
				c.writeMu.Unlock()
				return sealErr
			}
			if writeErr := c.writeRecord(wire); writeErr != nil {
				c.writeMu.Unlock()
				return writeErr
			}
			c.writeMu.Unlock()
			c.earlyMu.Lock()
			c.earlyPending = nil
			c.earlySent = true
			c.earlyMu.Unlock()
		}
	}
	inbox := newHandshakeInbox(0, c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
	messages, err := c.receiveHandshakeWithRetransmit(inbox, nil, out)
	if err != nil {
		return err
	}
	if messages.len() != 1 || messages.at(0).typ != handshakeTypeServerHello {
		return &ProtocolError{"expected ServerHello"}
	}
	serverHelloBody := messages.at(0).body
	serverHelloSequence := uint16(0)
	clientFinishedSequence := uint16(1)
	serverHandshakeStart := uint16(1)
	var helloRetrySuite uint16
	var transcript *transcriptHash
	echAccepted := false
	if len(serverHelloBody) >= 34 && string(serverHelloBody[2:34]) == string(helloRetryRequestRandom[:]) {
		if hello.earlyData {
			c.earlyMu.Lock()
			c.earlyRejected = c.earlySent
			c.earlyMu.Unlock()
			hello.earlyData = false
		}
		if ech != nil {
			ech.innerHello.earlyData = false
		}
		hrr, parseErr := parseHelloRetryRequest(serverHelloBody)
		if parseErr != nil {
			return parseErr
		}
		suite, parseErr := cipherSuiteForID(hrr.cipherSuite)
		if parseErr != nil {
			return parseErr
		}
		offered := false
		for _, id := range hello.cipherSuites {
			if id == hrr.cipherSuite {
				offered = true
			}
		}
		if !offered {
			return &ProtocolError{"HelloRetryRequest selected an unoffered cipher suite"}
		}
		helloRetrySuite = hrr.cipherSuite
		hrrSuite, suiteErr := cipherSuiteForID(hrr.cipherSuite)
		if suiteErr != nil {
			return suiteErr
		}
		initial := newTranscriptHash(suite.hash.New())
		_ = initial.add(handshakeTypeClientHello, 0, helloBody)
		transcript = newTranscriptHash(suite.hash.New())
		_ = transcript.addHelloRetryRequest(initial.sumInto(transcriptDigest[:0]), serverHelloBody)
		if ech != nil {
			innerInitial := newTranscriptHash(suite.hash.New())
			_ = innerInitial.add(handshakeTypeClientHello, 0, ech.innerBody)
			ech.innerTranscript = newTranscriptHash(suite.hash.New())
			innerHash := innerInitial.sumInto(transcriptDigest[:0])
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
				want := echAcceptConfirmation(suite, ech.innerHello.random, "hrr ech accept confirmation", confirmationTranscript.sumInto(transcriptDigest[:0]))
				accepted = subtle.ConstantTimeCompare(want, hrr.echConfirmation[:]) == 1
			}
			_ = ech.innerTranscript.addHelloRetryRequest(innerInitial.sumInto(transcriptDigest[:0]), serverHelloBody)
			if accepted {
				ech.acceptedAtHRR = true
				hello = ech.innerHello
				initial = innerInitial
				transcript = ech.innerTranscript
			} else {
				ech.rejected = true
			}
		} else if hrr.hasECHConfirmation && !greaseECH {
			return alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited ECH confirmation in HelloRetryRequest"})
		}
		if hrr.selectedGroup != 0 {
			supported := false
			alreadyOffered := false
			for _, group := range hello.supportedGroups {
				supported = supported || group == hrr.selectedGroup
			}
			for _, share := range hello.keyShares {
				alreadyOffered = alreadyOffered || share.group == hrr.selectedGroup
			}
			if !supported || alreadyOffered {
				return alertError(alertIllegalParameter, &ProtocolError{"HelloRetryRequest selected an invalid key share group"})
			}
			key, parseErr = generateEphemeralKey(hrr.selectedGroup, c.config.Rand)
			if parseErr != nil {
				return alertError(alertInternalError, parseErr)
			}
			hello.keyShares = []keyShareEntry{{group: key.groupID(), data: key.publicBytes()}}
		}
		hello.cookie = hrr.cookie
		if len(pskOffers) > 0 {
			pskOffers = filterPSKOffersByHash(pskOffers, hrrSuite.hash)
			for i := range pskOffers {
				if pskOffers[i].session != nil {
					pskOffers[i].age = obfuscatedTicketAge(pskOffers[i].session, c.config.Time())
				}
			}
		}
		if ech != nil && ech.rejected {
			helloBody, parseErr = hello.marshal()
		} else if len(pskOffers) > 0 {
			helloBody, parseErr = marshalClientHelloWithPSKOffers(hello, pskOffers, initial.sumInto(transcriptDigest[:0]), serverHelloBody)
		} else if clientSession != nil && sessionSuite.hash == hrrSuite.hash && len(hello.pskIdentity) > 0 {
			hello.obfuscatedAge = obfuscatedTicketAge(clientSession, c.config.Time())
			helloBody, parseErr = marshalClientHelloWithPSKBinder(hello, sessionSuite, clientSession.psk, initial.sumInto(transcriptDigest[:0]), serverHelloBody)
		} else {
			hello.pskIdentity = nil
			hello.pskBinder = nil
			hello.pskIdentities = nil
			hello.pskBinders = nil
			helloBody, parseErr = hello.marshal()
		}
		if parseErr != nil {
			return parseErr
		}
		wireHelloBody := helloBody
		if ech != nil && ech.acceptedAtHRR {
			ech.innerHello = hello
			ech.innerBody = helloBody
			ech.outerHello.cookie = append([]byte(nil), hello.cookie...)
			ech.outerHello.earlyData = false
			ech.outerHello.keyShares = cloneClientHello(hello).keyShares
			wireHelloBody, parseErr = computeOuterECH(ech.outerHello, hello, ech, false)
			if parseErr != nil {
				return parseErr
			}
			_ = ech.innerTranscript.add(handshakeTypeClientHello, 1, helloBody)
			transcript = ech.innerTranscript
		} else {
			_ = transcript.add(handshakeTypeClientHello, 1, helloBody)
		}
		second, _, parseErr := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeClientHello, sequence: 1, body: wireHelloBody}}, c.currentMTU(), 0, out.nextRecordSequence())
		if parseErr != nil {
			return parseErr
		}
		c.plainSendSequence.Store(second.nextRecordSequence())
		if parseErr = c.writeFlight(c.conn, second); parseErr != nil {
			return parseErr
		}
		messages, parseErr = c.receiveHandshakeWithRetransmit(inbox, nil, second)
		if parseErr != nil {
			return parseErr
		}
		if messages.len() != 1 || messages.at(0).typ != handshakeTypeServerHello {
			return &ProtocolError{"expected ServerHello after HelloRetryRequest"}
		}
		serverHelloBody = messages.at(0).body
		if len(serverHelloBody) >= 34 && string(serverHelloBody[2:34]) == string(helloRetryRequestRandom[:]) {
			return alertError(alertUnexpectedMessage, &ProtocolError{"received a second HelloRetryRequest"})
		}
		serverHelloSequence = 1
		clientFinishedSequence = 2
		serverHandshakeStart = 2
	}
	sh, err := parseServerHello(serverHelloBody)
	if err != nil {
		return err
	}
	suite, err := cipherSuiteForID(sh.cipherSuite)
	if err != nil {
		return err
	}
	if ech != nil && !ech.rejected {
		if ech.innerTranscript == nil {
			ech.innerTranscript = newTranscriptHash(suite.hash.New())
			_ = ech.innerTranscript.add(handshakeTypeClientHello, 0, ech.innerBody)
		}
		if len(serverHelloBody) < 34 {
			return alertError(alertDecodeError, &ProtocolError{"truncated ServerHello ECH confirmation"})
		}
		zeroServerHello := append([]byte(nil), serverHelloBody...)
		clear(zeroServerHello[26:34])
		confirmationTranscript := ech.innerTranscript.clone()
		_ = confirmationTranscript.add(handshakeTypeServerHello, serverHelloSequence, zeroServerHello)
		want := echAcceptConfirmation(suite, ech.innerHello.random, "ech accept confirmation", confirmationTranscript.sumInto(transcriptDigest[:0]))
		confirmed := subtle.ConstantTimeCompare(want, sh.random[24:]) == 1
		if ech.acceptedAtHRR && !confirmed {
			return alertError(alertIllegalParameter, &ProtocolError{"ServerHello did not confirm ECH after accepted HelloRetryRequest"})
		}
		if confirmed {
			echAccepted = true
			hello = ech.innerHello
			helloBody = ech.innerBody
			transcript = ech.innerTranscript
		} else {
			ech.rejected = true
		}
	}
	offered := false
	for _, id := range hello.cipherSuites {
		if id == sh.cipherSuite {
			offered = true
		}
	}
	if !offered {
		return &ProtocolError{"server selected an unoffered cipher suite"}
	}
	if helloRetrySuite != 0 && sh.cipherSuite != helloRetrySuite {
		return alertError(alertIllegalParameter, &ProtocolError{"ServerHello changed cipher suite after HelloRetryRequest"})
	}
	if len(sh.sessionID) != 0 {
		return &ProtocolError{"ServerHello legacy session ID must be empty"}
	}
	if err = validateServerHelloConnectionID(hello, sh); err != nil {
		return err
	}
	if err = validateServerHelloReturnRoutability(hello, sh); err != nil {
		return err
	}
	if sh.hasConnectionID {
		c.connectionIDNegotiated = true
		c.sendConnectionID = append([]byte(nil), sh.connectionID...)
		c.receiveConnectionID = append([]byte(nil), hello.connectionID...)
		c.localCIDUpdatesAllowed = len(c.receiveConnectionID) > 0
		c.peerCIDUpdatesAllowed = len(c.sendConnectionID) > 0
	}
	c.returnRoutabilityCheckNegotiated = sh.returnRoutability
	usingPSK := sh.selectedIdentity != nil
	if ech != nil && ech.rejected && usingPSK {
		return alertError(alertIllegalParameter, &ProtocolError{"server selected an outer GREASE PSK after rejecting ECH"})
	}
	var resumed bool
	var externalPSK *externalPSKSelection
	var selectedOffer *clientPSKOffer
	var singleTicketOffer clientPSKOffer
	if usingPSK {
		if len(pskOffers) == 0 && clientSession != nil && *sh.selectedIdentity == 0 {
			singleTicketOffer = clientPSKOffer{psk: clientSession.psk, suite: sessionSuite, session: clientSession, external: clientSession.externalPSK}
			selectedOffer = &singleTicketOffer
		} else {
			if int(*sh.selectedIdentity) >= len(pskOffers) {
				return &ProtocolError{"server selected an invalid PSK identity"}
			}
			selectedOffer = &pskOffers[*sh.selectedIdentity]
		}
		if suite.hash != selectedOffer.suite.hash {
			return &ProtocolError{"server selected a cipher suite incompatible with the PSK"}
		}
		resumed = selectedOffer.session != nil
		externalPSK = selectedOffer.external
	}
	if !slices.ContainsFunc(hello.keyShares, func(share keyShareEntry) bool { return share.group == sh.keyShare.group }) {
		return alertError(alertIllegalParameter, &ProtocolError{"ServerHello selected an unoffered key share"})
	}
	shared, err := key.sharedSecret(sh.keyShare.group, sh.keyShare.data)
	if err != nil {
		return err
	}
	if transcript == nil {
		transcript = newTranscriptHash(suite.hash.New())
		_ = transcript.add(handshakeTypeClientHello, 0, helloBody)
	}
	_ = transcript.add(handshakeTypeServerHello, serverHelloSequence, serverHelloBody)
	var psk []byte
	if usingPSK {
		psk = selectedOffer.psk
	}
	schedule := newKeySchedule(suite, psk)
	if err = schedule.deriveHandshake(shared, transcript.sumInto(transcriptDigest[:0])); err != nil {
		return err
	}
	receiveCipher, err := newRecordCipher(suite, schedule.serverHandshakeTraffic, 2, c.config.ReplayWindow)
	if err != nil {
		return err
	}
	sendCipher, err := newRecordCipher(suite, schedule.clientHandshakeTraffic, 2, c.config.ReplayWindow)
	if err != nil {
		return err
	}
	if c.connectionIDNegotiated {
		if err = receiveCipher.setConnectionID(c.receiveConnectionID); err != nil {
			return err
		}
		if err = sendCipher.setConnectionID(c.sendConnectionID); err != nil {
			return err
		}
	}
	c.sendCipher = sendCipher
	inbox = newHandshakeInbox(serverHandshakeStart, c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
	var peerCerts []*x509.Certificate
	var chains [][]*x509.Certificate
	if resumed {
		peerCerts = append([]*x509.Certificate(nil), selectedOffer.session.peerCertificates...)
		chains = make([][]*x509.Certificate, len(selectedOffer.session.verifiedChains))
		for i := range selectedOffer.session.verifiedChains {
			chains[i] = append([]*x509.Certificate(nil), selectedOffer.session.verifiedChains[i]...)
		}
	}
	var negotiated string
	var certificateRequest *certificateRequestMessage
	var certificateRequestCompressionAlgorithms *certificateCompressionAlgorithms
	verifiedServerSignature := false
	finished := false
	serverStage := serverExpectEncryptedExtensions
	var serverFinishedSequence uint16
	for !finished {
		messages, err = receiveHandshakeMessageWithEarlyBatch(c.conn, inbox, receiveCipher, nil, nil, nil, sendCipher, c.currentMTU(), c)
		if err != nil {
			return err
		}
		for index := 0; index < messages.len(); index++ {
			message := messages.at(index)
			if err = serverStage.accept(message.typ, usingPSK); err != nil {
				return err
			}
			switch message.typ {
			case handshakeTypeEncryptedExtensions:
				ee, parseErr := parseEncryptedExtensions(message.body)
				if parseErr != nil {
					return parseErr
				}
				var acceptedEarly bool
				var retryConfigs []byte
				negotiated, acceptedEarly, retryConfigs, parseErr = validateEncryptedExtensions(hello, &ee)
				if parseErr != nil {
					return parseErr
				}
				if ee.hasTicketRequest {
					requested := hello.ticketRequest.NewSessionCount
					if resumed {
						requested = hello.ticketRequest.ResumptionCount
					}
					c.sessionTicketRequest = &sessionTicketRequestState{limit: min(ee.expectedTicketCount, requested)}
				}
				if retryConfigs != nil {
					if ech != nil && echAccepted {
						return alertError(alertUnsupportedExtension, &ProtocolError{"server sent ECH retry configurations after accepting ECH"})
					}
					if ech != nil && ech.rejected {
						ech.retryConfigs = append([]byte(nil), retryConfigs...)
					}
					// GREASE clients validate but never retain retry configurations.
				}
				if ee.hasRecordSizeLimit {
					c.recordSizeLimitNegotiated = true
					c.localRecordSizeLimit = c.config.RecordSizeLimit
					c.peerRecordSizeLimit = ee.recordSizeLimit
					receiveCipher.setPlaintextLimit(c.localRecordSizeLimit)
					sendCipher.setPlaintextLimit(c.peerRecordSizeLimit)
				}
				if parseErr = validateEarlyDataSelection(acceptedEarly, sh.selectedIdentity); parseErr != nil {
					return parseErr
				}
				if acceptedEarly {
					c.earlyMu.Lock()
					c.earlyAccepted = true
					c.earlyMu.Unlock()
				} else if hello.earlyData {
					c.earlyMu.Lock()
					c.earlyRejected = c.earlySent
					c.earlyMu.Unlock()
				}
				_ = transcript.add(message.typ, message.sequence, message.body)
			case handshakeTypeCertificate, handshakeTypeCompressedCertificate:
				if usingPSK {
					return &ProtocolError{"server sent Certificate in a PSK handshake"}
				}
				certMsg, parseErr := parseCertificateHandshakeMessage(message.typ, message.body, hello.certificateCompressionAlgorithms(), c.config.MaxHandshakeMessage)
				if parseErr != nil {
					return parseErr
				}
				if parseErr = validateCertificateMessage(certMsg, nil); parseErr != nil {
					return parseErr
				}
				certificateSchemes := hello.certificateSignatureSchemes
				if len(certificateSchemes) == 0 {
					certificateSchemes = hello.signatureSchemes
				}
				if ech != nil && ech.rejected {
					peerCerts, chains, parseErr = verifyCertificateChainForECHRejection(c.config, certMsg, certificateSchemes, ech.config.publicName)
				} else {
					peerCerts, chains, parseErr = verifyCertificateChain(c.config, certMsg, true, certificateSchemes)
				}
				if parseErr != nil {
					return parseErr
				}
				if ech != nil && ech.rejected && c.config.EncryptedClientHelloRejectionVerify != nil {
					rejectionState := ConnectionState{Version: VersionDTLS13, CipherSuite: suite.id, NegotiatedProtocol: negotiated, ServerName: ech.config.publicName, PeerCertificates: peerCerts, VerifiedChains: chains, ECHAccepted: false}
					if verifyErr := c.config.EncryptedClientHelloRejectionVerify(rejectionState); verifyErr != nil {
						return alertError(alertAccessDenied, verifyErr)
					}
				}
				_ = transcript.add(message.typ, message.sequence, message.body)
			case handshakeTypeCertificateRequest:
				if usingPSK {
					return &ProtocolError{"server requested a certificate in a PSK handshake"}
				}
				if certificateRequest != nil {
					return &ProtocolError{"duplicate CertificateRequest"}
				}
				certificateRequest, certificateRequestCompressionAlgorithms, err = parseCertificateRequestWithCompression(message.body)
				if err != nil {
					return err
				}
				if len(certificateRequest.requestContext) != 0 {
					return alertError(alertIllegalParameter, &ProtocolError{"initial CertificateRequest context must be empty"})
				}
				_ = transcript.add(message.typ, message.sequence, message.body)
			case handshakeTypeCertificateVerify:
				if len(peerCerts) == 0 {
					return &ProtocolError{"CertificateVerify before Certificate"}
				}
				cv, parseErr := parseCertificateVerify(message.body)
				if parseErr != nil {
					return parseErr
				}
				offered := false
				for _, scheme := range hello.signatureSchemes {
					if scheme == cv.algorithm {
						offered = true
					}
				}
				if !offered {
					return &ProtocolError{"server selected an unoffered signature scheme"}
				}
				if parseErr = verifyCertificateVerify(peerCerts[0].PublicKey, cv.algorithm, transcript.sumInto(transcriptDigest[:0]), cv.signature, true); parseErr != nil {
					return alertError(alertDecryptError, parseErr)
				}
				verifiedServerSignature = true
				_ = transcript.add(message.typ, message.sequence, message.body)
			case handshakeTypeFinished:
				if !usingPSK && (len(peerCerts) == 0 || !verifiedServerSignature) {
					return &ProtocolError{"server authentication messages are incomplete"}
				}
				verify, parseErr := parseFinished(message.body, suite.hash.Size())
				if parseErr != nil {
					return parseErr
				}
				if !schedule.verifyFinished(schedule.serverHandshakeTraffic, transcript.sumInto(transcriptDigest[:0]), verify) {
					return alertError(alertDecryptError, &ProtocolError{"server Finished verification failed"})
				}
				_ = transcript.add(message.typ, message.sequence, message.body)
				serverFinishedSequence = message.sequence
				finished = true
			default:
				return &ProtocolError{"unexpected server handshake message"}
			}
		}
	}
	if err = schedule.deriveApplication(transcript.sumInto(transcriptDigest[:0])); err != nil {
		return err
	}
	var clientMessages []handshakeMessage
	nextClientSequence := clientFinishedSequence
	if certificateRequest != nil {
		certMessage := &certificateMessage{requestContext: certificateRequest.requestContext}
		var clientCertificate *tls.Certificate
		if (ech == nil || !ech.rejected) && len(c.config.Certificates) > 0 {
			candidate := &c.config.Certificates[0]
			certificateSchemes := certificateRequest.certificateSignatureSchemes
			if len(certificateSchemes) == 0 {
				certificateSchemes = certificateRequest.signatureSchemes
			}
			if validateConfiguredCertificate(candidate, certificateSchemes, false) == nil {
				clientCertificate = candidate
				for _, der := range clientCertificate.Certificate {
					certMessage.certificates = append(certMessage.certificates, certificateEntry{data: der})
				}
			}
		}
		certBody, marshalErr := certMessage.marshal()
		if marshalErr != nil {
			return marshalErr
		}
		certificateType, certificateBody, selectErr := certificateHandshakeMessage(certBody, certificateRequestCompressionAlgorithms, c.config.EnableCertificateCompression, c.config.certificateCompressionCache())
		if selectErr != nil {
			return selectErr
		}
		clientMessages = append(clientMessages, handshakeMessage{typ: certificateType, sequence: nextClientSequence, body: certificateBody})
		_ = transcript.add(certificateType, nextClientSequence, certificateBody)
		nextClientSequence++
		if clientCertificate != nil {
			signer, ok := clientCertificate.PrivateKey.(crypto.Signer)
			if !ok {
				return errors.New("dtls13: client private key does not implement crypto.Signer")
			}
			scheme, selectErr := selectSignatureScheme(signer, certificateRequest.signatureSchemes)
			if selectErr != nil {
				return selectErr
			}
			signature, signErr := signCertificateVerify(c.config.Rand, signer, scheme, transcript.sumInto(transcriptDigest[:0]), false)
			if signErr != nil {
				return signErr
			}
			cvBody, marshalErr := (&certificateVerifyMessage{algorithm: scheme, signature: signature}).marshal()
			if marshalErr != nil {
				return marshalErr
			}
			clientMessages = append(clientMessages, handshakeMessage{typ: handshakeTypeCertificateVerify, sequence: nextClientSequence, body: cvBody})
			_ = transcript.add(handshakeTypeCertificateVerify, nextClientSequence, cvBody)
			nextClientSequence++
		}
	}
	clientFinishedSequence = nextClientSequence
	clientFinished := schedule.finishedVerifyData(schedule.clientHandshakeTraffic, transcript.sumInto(transcriptDigest[:0]))
	_ = transcript.add(handshakeTypeFinished, clientFinishedSequence, clientFinished)
	clientMessages = append(clientMessages, handshakeMessage{typ: handshakeTypeFinished, sequence: clientFinishedSequence, body: clientFinished})
	clientFlight, err := buildProtectedFlight(clientMessages, c.currentMTU(), sendCipher)
	if err != nil {
		return err
	}
	if err = c.writeFlight(c.conn, clientFlight); err != nil {
		return err
	}
	applicationACKCipher, err := newRecordCipher(suite, schedule.serverApplicationTraffic, 3, c.config.ReplayWindow)
	if err != nil {
		return err
	}
	applicationACKCipher.setPlaintextLimit(c.localRecordSizeLimit)
	if c.connectionIDNegotiated {
		if err = applicationACKCipher.setConnectionID(c.receiveConnectionID); err != nil {
			return err
		}
	}
	acknowledged, err := c.receiveACKWithRetransmit(clientFlight, receiveCipher, applicationACKCipher)
	if err != nil {
		return err
	}
	if len(clientFlight.records) == 0 {
		return &ProtocolError{"empty client Finished flight"}
	}
	for _, record := range clientFlight.records {
		if !record.acknowledgedBy(acknowledged) {
			return &ProtocolError{"server ACK did not cover complete client flight"}
		}
	}
	if err = c.installApplicationKeysAt(suite, schedule.clientApplicationTraffic, schedule.serverApplicationTraffic, clientFinishedSequence+1); err != nil {
		return err
	}
	if ech != nil && ech.rejected {
		return alertError(alertECHRequired, &ECHRejectionError{RetryConfigList: append([]byte(nil), ech.retryConfigs...)})
	}
	if err = schedule.deriveResumption(transcript.sumInto(transcriptDigest[:0])); err != nil {
		return err
	}
	c.resumptionSuite = suite
	c.resumptionMasterSecret = append([]byte(nil), schedule.resumptionMasterSecret...)
	if c.sessionTicketRequest != nil {
		if resumed {
			c.sessionTicketRequest.group = selectedOffer.session.ticketGroup
		} else {
			c.sessionTicketRequest.group = sha256.Sum256(c.resumptionMasterSecret)
		}
	}
	c.postHandshakeTranscript = transcript.clone()
	if err = c.receiveEpochs.install(receiveCipher); err != nil {
		return err
	}
	c.completedPeerFlightStart = serverHandshakeStart
	c.completedPeerFlightEnd = serverFinishedSequence
	c.hasCompletedPeerFlight = true
	exporter := newExporter(suite, schedule.exporterMasterSecret)
	exporter.externalPSK = externalPSK
	c.mu.Lock()
	c.state = ConnectionState{Version: VersionDTLS13, HandshakeComplete: true, DidResume: resumed, ECHAccepted: echAccepted, CipherSuite: suite.id, NegotiatedProtocol: negotiated, ServerName: c.config.ServerName, PeerCertificates: peerCerts, VerifiedChains: chains, LocalConnectionID: append([]byte(nil), c.receiveConnectionID...), PeerConnectionID: append([]byte(nil), c.sendConnectionID...), ReturnRoutabilityCheck: c.returnRoutabilityCheckNegotiated, RecordSizeLimitNegotiated: c.recordSizeLimitNegotiated, LocalRecordSizeLimit: c.localRecordSizeLimit, PeerRecordSizeLimit: c.peerRecordSizeLimit, exporter: exporter}
	c.mu.Unlock()
	if clientSession != nil && !resumed {
		discardClientSessionGroup(c.config, c.conn, clientSession.ticketGroup)
	}
	return nil
}

func equalClientHelloAfterHRR(initial, second *clientHello, requestedGroup tls.CurveID) bool {
	if initial == nil || second == nil || second.earlyData {
		return false
	}
	// The second list may only remove identities; retained identities keep
	// their original order. Ticket ages and binders are recomputed after HRR.
	matched := 0
	for _, candidate := range initial.pskIdentities {
		if matched < len(second.pskIdentities) && equalBytes(candidate.identity, second.pskIdentities[matched].identity) {
			matched++
		}
	}
	if matched != len(second.pskIdentities) {
		return false
	}
	if requestedGroup != 0 {
		if len(second.keyShares) != 1 || second.keyShares[0].group != requestedGroup {
			return false
		}
	} else if !equalKeyShareEntries(initial.keyShares, second.keyShares) {
		return false
	}
	return initial.random == second.random &&
		equalBytes(initial.sessionID, second.sessionID) &&
		slices.Equal(initial.cipherSuites, second.cipherSuites) &&
		slices.Equal(initial.signatureSchemes, second.signatureSchemes) &&
		slices.Equal(initial.certificateSignatureSchemes, second.certificateSignatureSchemes) &&
		slices.Equal(initial.supportedGroups, second.supportedGroups) &&
		initial.serverName == second.serverName &&
		slices.Equal(initial.alpn, second.alpn) &&
		initial.pskDHE == second.pskDHE &&
		equalBytes(initial.connectionID, second.connectionID) &&
		initial.hasConnectionID == second.hasConnectionID &&
		initial.returnRoutability == second.returnRoutability &&
		initial.postHandshakeAuth == second.postHandshakeAuth &&
		initial.ticketRequest == second.ticketRequest &&
		initial.certificateCompressionOffered == second.certificateCompressionOffered &&
		initial.recordSizeLimit == second.recordSizeLimit &&
		initial.hasRecordSizeLimit == second.hasRecordSizeLimit &&
		equalExtensionMaps(initial.unknownExtensions, second.unknownExtensions)
}

func equalKeyShareEntries(left, right []keyShareEntry) bool {
	return slices.EqualFunc(left, right, func(a, b keyShareEntry) bool {
		return a.group == b.group && equalBytes(a.data, b.data)
	})
}

func equalExtensionMaps(left, right map[uint16][]byte) bool {
	if len(left) != len(right) {
		return false
	}
	for typ, value := range left {
		other, ok := right[typ]
		if !ok || !equalBytes(value, other) {
			return false
		}
	}
	return true
}

func removedPSKsAreIncompatible(config *Config, initial, second *clientHello, selectedSuite *cipherSuite) (bool, error) {
	if initial == nil || second == nil || selectedSuite == nil || len(initial.pskIdentities) == 0 {
		return true, nil
	}
	var protector *sessionTicketProtector
	if !config.SessionTicketsDisabled {
		var err error
		protector, err = newSessionTicketProtector(config.SessionTicketKey, config.Rand, config.Time)
		if err != nil {
			return false, err
		}
	}
	retained := make([]bool, len(initial.pskIdentities))
	next := 0
	for _, identity := range second.pskIdentities {
		for next < len(initial.pskIdentities) {
			index := next
			next++
			if equalBytes(initial.pskIdentities[index].identity, identity.identity) {
				retained[index] = true
				break
			}
		}
	}
	for index, identity := range initial.pskIdentities {
		if retained[index] {
			continue
		}
		if findExternalPSK(config, identity.identity, selectedSuite.hash) != nil {
			return false, nil
		}
		if protector == nil {
			continue
		}
		state, openErr := protector.open(identity.identity)
		if openErr != nil {
			continue
		}
		ticketSuite, suiteErr := cipherSuiteForID(state.suite)
		if suiteErr == nil && ticketSuite.hash == selectedSuite.hash {
			return false, nil
		}
	}
	return true, nil
}

func requireCertificateSignatureAlgorithms(hello *clientHello, resumed bool) error {
	if !resumed && (hello == nil || len(hello.signatureSchemes) == 0) {
		return alertError(alertMissingExtension, &ProtocolError{"certificate authentication requires signature_algorithms"})
	}
	return nil
}

func (c *Conn) serverHandshake() error {
	var transcriptDigest [maxSupportedHashSize]byte
	var amplification amplificationGuard
	preValidationConn := &amplificationConn{Conn: c.conn, guard: &amplification}
	inbox := newHandshakeInbox(0, c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
	messages, err := receiveHandshakeMessageBatch(preValidationConn, inbox, nil)
	if err != nil {
		return err
	}
	if messages.len() != 1 || messages.at(0).typ != handshakeTypeClientHello {
		return &ProtocolError{"expected ClientHello"}
	}
	helloBody := messages.at(0).body
	outerHello, err := parseClientHello(helloBody)
	if err != nil {
		return err
	}
	ch := outerHello
	var echContext *echServerContext
	var echKeys []EncryptedClientHelloKey
	echOuterOffered := false
	if len(outerHello.encryptedClientHello()) > 0 {
		typ, _, _, _, _, parseErr := parseECHExt(outerHello.encryptedClientHello())
		if parseErr != nil {
			return parseErr
		}
		echOuterOffered = typ == echOuterType
		if echOuterOffered {
			echKeys = c.config.EncryptedClientHelloKeys
			if c.config.GetEncryptedClientHelloKeys != nil {
				echKeys, parseErr = c.config.GetEncryptedClientHelloKeys(&ClientHelloInfo{ServerName: outerHello.serverName, SupportedProtos: outerHello.alpn, Conn: c})
				if parseErr != nil {
					return parseErr
				}
			}
		}
		ch, helloBody, echContext, err = processECHClientHello(outerHello, helloBody, echKeys)
		if err != nil {
			return err
		}
	}
	echAccepted := echContext != nil
	echRejected := echOuterOffered && !echAccepted
	c.postHandshakeAuthOffered = ch.postHandshakeAuth
	initialClientHello := *ch
	initialClientHello.cookie = nil
	initialClientHello.pskBinder = nil
	initialClientHello.pskBinders = nil
	suite, err := selectCipherSuite(c.config.CipherSuites, ch.cipherSuites)
	if err != nil {
		return err
	}
	if len(c.config.ExternalPSKs) > 0 {
		suite = preferExternalPSKCipherSuite(c.config, ch, suite)
	}
	hrrUsed := true
	serverSequenceOffset := uint16(1)
	clientFinishedSequence := uint16(2)
	initialHashTranscript := newTranscriptHash(suite.hash.New())
	_ = initialHashTranscript.add(handshakeTypeClientHello, 0, helloBody)
	var hrrBody []byte
	var hrrFlight *flight
	var requestedGroup tls.CurveID
	var psk []byte
	var resumed bool
	var usingPSK bool
	var resumedSession *sessionTicketState
	var externalPSK *externalPSKSelection
	var share keyShareEntry
	var negotiated string
	var shared []byte
	var selectedPSKIdentity uint16

	// A server may skip the cookie exchange for a resumed PSK handshake in a
	// trusted environment. This is the only path on which epoch-1 data is
	// accepted; the default remains the RFC-recommended HRR/cookie exchange.
	if c.config.AllowEarlyDataWithoutCookie && ch.earlyData && len(ch.pskIdentity) > 0 {
		_, selectErr := selectKeyShare(c.config.CurvePreferences, ch.keyShares)
		if selectErr == nil {
			candidateProtocol, protocolErr := negotiateALPN(c.config.NextProtos, ch.alpn)
			if protocolErr == nil {
				candidate, accepted, acceptErr := c.acceptPSK(ch, helloBody, initialHashTranscript.sumInto(transcriptDigest[:0]), nil, suite, candidateProtocol)
				if acceptErr != nil {
					return acceptErr
				}
				if accepted && candidate.session != nil {
					resumedSession = candidate.session
					psk, resumed, usingPSK, selectedPSKIdentity = candidate.psk, true, true, candidate.identity
					externalPSK = candidate.external
					hrrUsed = false
					serverSequenceOffset = 0
					clientFinishedSequence = 1
				}
			}
		}
	}

	if hrrUsed {
		if err = ensureCookieProtector(c.config); err != nil {
			return err
		}
		protector := &c.config.state.cookieProtector
		address := []byte(c.conn.RemoteAddr().String())
		cookie, cookieErr := protector.seal(address, initialHashTranscript.sumInto(transcriptDigest[:0]))
		if cookieErr != nil {
			return cookieErr
		}
		if _, shareErr := selectKeyShare(c.config.CurvePreferences, ch.keyShares); shareErr != nil {
			for _, preference := range c.config.CurvePreferences {
				for _, supported := range ch.supportedGroups {
					if preference == supported {
						requestedGroup = preference
						break
					}
				}
				if requestedGroup != 0 {
					break
				}
			}
			if requestedGroup == 0 {
				return shareErr
			}
		}
		hrr := &helloRetryRequest{cipherSuite: suite.id, cookie: cookie, selectedGroup: requestedGroup}
		if echAccepted {
			hrr.hasECHConfirmation = true
			zeroBody, marshalErr := hrr.marshal()
			if marshalErr != nil {
				return marshalErr
			}
			confirmationTranscript := newTranscriptHash(suite.hash.New())
			_ = confirmationTranscript.addHelloRetryRequest(initialHashTranscript.sumInto(transcriptDigest[:0]), zeroBody)
			copy(hrr.echConfirmation[:], echAcceptConfirmation(suite, ch.random, "hrr ech accept confirmation", confirmationTranscript.sumInto(transcriptDigest[:0])))
		}
		hrrBody, err = hrr.marshal()
		if err != nil {
			return err
		}
		hrrFlight, _, err = buildPlainFlight([]handshakeMessage{{typ: handshakeTypeServerHello, sequence: 0, body: hrrBody}}, c.currentMTU(), 0, 0)
		if err != nil {
			return err
		}
		if err = c.writeFlight(preValidationConn, hrrFlight); err != nil {
			return err
		}
		messages, err = c.receiveSecondClientHello(preValidationConn, inbox, hrrFlight)
		if err != nil {
			return err
		}
		if messages.len() != 1 || messages.at(0).typ != handshakeTypeClientHello {
			return &ProtocolError{"expected second ClientHello"}
		}
		secondBody := messages.at(0).body
		second, parseErr := parseClientHello(secondBody)
		if parseErr != nil {
			return parseErr
		}
		if echAccepted {
			second, secondBody, parseErr = processSecondECHClientHello(second, secondBody, echContext)
			if parseErr != nil {
				return parseErr
			}
		}
		if !equalClientHelloAfterHRR(&initialClientHello, second, requestedGroup) {
			return &ProtocolError{"second ClientHello changed fields other than cookie"}
		}
		compatibleRemoval, removalErr := removedPSKsAreIncompatible(c.config, &initialClientHello, second, suite)
		if removalErr != nil {
			return removalErr
		}
		if !compatibleRemoval {
			return alertError(alertIllegalParameter, &ProtocolError{"second ClientHello removed a PSK compatible with the HelloRetryRequest cipher suite"})
		}
		cookieHash, cookieErr := protector.open(address, second.cookie)
		if cookieErr != nil {
			return alertError(alertIllegalParameter, &ProtocolError{"invalid HelloRetryRequest cookie"})
		}
		if string(cookieHash) != string(initialHashTranscript.sumInto(transcriptDigest[:0])) {
			return alertError(alertIllegalParameter, &ProtocolError{"HelloRetryRequest cookie transcript mismatch"})
		}
		amplification.validate()
		helloBody = secondBody
		ch = second
		suite, err = selectCipherSuite(c.config.CipherSuites, ch.cipherSuites)
		if err != nil {
			return err
		}
		if len(c.config.ExternalPSKs) > 0 {
			suite = preferExternalPSKCipherSuite(c.config, ch, suite)
		}
		share, err = selectKeyShare(c.config.CurvePreferences, ch.keyShares)
		if err != nil {
			return err
		}
		negotiated, err = negotiateALPN(c.config.NextProtos, ch.alpn)
		if err != nil {
			return err
		}
		selection, accepted, acceptErr := c.acceptPSK(ch, helloBody, initialHashTranscript.sumInto(transcriptDigest[:0]), hrrBody, suite, negotiated)
		err = acceptErr
		if err != nil {
			return err
		}
		if accepted {
			psk, usingPSK, selectedPSKIdentity = selection.psk, true, selection.identity
			resumedSession, externalPSK = selection.session, selection.external
			resumed = resumedSession != nil
		}
	} else {
		serverKeyShare, selectErr := selectKeyShare(c.config.CurvePreferences, ch.keyShares)
		if selectErr != nil {
			return selectErr
		}
		share = serverKeyShare
		negotiated, err = negotiateALPN(c.config.NextProtos, ch.alpn)
		if err != nil {
			return err
		}
	}
	if ch.hasRecordSizeLimit {
		c.recordSizeLimitNegotiated = true
		c.localRecordSizeLimit = c.config.RecordSizeLimit
		c.peerRecordSizeLimit = effectiveRecordSizeLimit(ch.recordSizeLimit)
	}
	var serverShare []byte
	serverShare, shared, err = generateServerKeyShare(share.group, share.data, c.config.Rand)
	if err != nil {
		return err
	}
	var cert *tls.Certificate
	var signer crypto.Signer
	var scheme tls.SignatureScheme
	if !usingPSK {
		if err = requireCertificateSignatureAlgorithms(ch, false); err != nil {
			return err
		}
		cert, err = c.serverCertificate(ch)
		if err != nil {
			return err
		}
		var ok bool
		signer, ok = cert.PrivateKey.(crypto.Signer)
		if !ok {
			return errors.New("dtls13: server private key does not implement crypto.Signer")
		}
		certificateSchemes := ch.certificateSignatureSchemes
		if len(certificateSchemes) == 0 {
			certificateSchemes = ch.signatureSchemes
		}
		if validateErr := validateConfiguredCertificate(cert, certificateSchemes, true); validateErr != nil {
			return alertError(alertHandshakeFailure, validateErr)
		}
		scheme, err = selectSignatureScheme(signer, ch.signatureSchemes)
		if err != nil {
			return err
		}
	}
	// DTLS 1.3 does not use TLS compatibility mode; the server MUST NOT
	// echo legacy_session_id (RFC 9147 section 5).
	sh := &serverHello{cipherSuite: suite.id, keyShare: keyShareEntry{group: share.group, data: serverShare}}
	if ch.hasConnectionID && c.config.ConnectionID != nil {
		c.connectionIDNegotiated = true
		c.sendConnectionID = append([]byte(nil), ch.connectionID...)
		c.receiveConnectionID = append([]byte(nil), c.config.ConnectionID...)
		c.localCIDUpdatesAllowed = len(c.receiveConnectionID) > 0
		c.peerCIDUpdatesAllowed = len(c.sendConnectionID) > 0
		sh.hasConnectionID = true
		sh.connectionID = append([]byte(nil), c.receiveConnectionID...)
	}
	if c.connectionIDNegotiated && ch.returnRoutability && !c.config.DisableReturnRoutabilityCheck {
		c.returnRoutabilityCheckNegotiated = true
		sh.returnRoutability = true
	}
	if usingPSK {
		selected := selectedPSKIdentity
		sh.selectedIdentity = &selected
	}
	if _, err = io.ReadFull(c.config.Rand, sh.random[:]); err != nil {
		return err
	}
	transcript := newTranscriptHash(suite.hash.New())
	serverHelloSequence := uint16(0)
	firstPlainRecordSequence := uint64(0)
	if hrrUsed {
		_ = transcript.addHelloRetryRequest(initialHashTranscript.sumInto(transcriptDigest[:0]), hrrBody)
		_ = transcript.add(handshakeTypeClientHello, 1, helloBody)
		serverHelloSequence = 1
		if hrrFlight != nil {
			firstPlainRecordSequence = hrrFlight.nextRecordSequence()
		}
	} else {
		_ = transcript.add(handshakeTypeClientHello, 0, helloBody)
	}
	if echAccepted {
		clear(sh.random[24:])
		zeroBody, marshalErr := sh.marshal()
		if marshalErr != nil {
			return marshalErr
		}
		confirmationTranscript := transcript.clone()
		_ = confirmationTranscript.add(handshakeTypeServerHello, serverHelloSequence, zeroBody)
		copy(sh.random[24:], echAcceptConfirmation(suite, ch.random, "ech accept confirmation", confirmationTranscript.sumInto(transcriptDigest[:0])))
	}
	shBody, err := sh.marshal()
	if err != nil {
		return err
	}
	_ = transcript.add(handshakeTypeServerHello, serverHelloSequence, shBody)
	schedule := newKeySchedule(suite, psk)
	if err = schedule.deriveHandshake(shared, transcript.sumInto(transcriptDigest[:0])); err != nil {
		return err
	}
	serverCipher, err := newRecordCipher(suite, schedule.serverHandshakeTraffic, 2, c.config.ReplayWindow)
	if err != nil {
		return err
	}
	clientCipher, err := newRecordCipher(suite, schedule.clientHandshakeTraffic, 2, c.config.ReplayWindow)
	if err != nil {
		return err
	}
	serverCipher.setPlaintextLimit(c.peerRecordSizeLimit)
	clientCipher.setPlaintextLimit(c.localRecordSizeLimit)
	if c.connectionIDNegotiated {
		if err = serverCipher.setConnectionID(c.sendConnectionID); err != nil {
			return err
		}
		if err = clientCipher.setConnectionID(c.receiveConnectionID); err != nil {
			return err
		}
	}
	var earlyCipher *recordCipher
	if !hrrUsed && c.earlyAccepted {
		earlySchedule := newKeySchedule(suite, psk)
		earlyCipher, err = newRecordCipher(suite, earlySchedule.earlyTrafficSecret(initialHashTranscript.sumInto(transcriptDigest[:0])), 1, c.config.ReplayWindow)
		if err != nil {
			return err
		}
		earlyCipher.setPlaintextLimit(resumedSession.recordSizeLimit)
	}
	plain, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeServerHello, sequence: serverHelloSequence, body: shBody}}, c.currentMTU(), 0, firstPlainRecordSequence)
	if err != nil {
		return err
	}
	serverFlightConn := net.Conn(c.conn)
	if !hrrUsed {
		serverFlightConn = preValidationConn
	}
	ticketCount := uint8(1)
	if ch.ticketRequest.Enabled {
		ticketCount = ch.ticketRequest.NewSessionCount
		if resumed {
			ticketCount = ch.ticketRequest.ResumptionCount
		}
		ticketCount = min(ticketCount, c.config.MaxSessionTickets)
	}
	if c.config.SessionTicketsDisabled {
		ticketCount = 0
	}
	ee := &encryptedExtensions{recordSizeLimit: c.config.RecordSizeLimit, hasRecordSizeLimit: ch.hasRecordSizeLimit}
	if negotiated != "" || c.earlyAccepted || echRejected || ch.ticketRequest.Enabled {
		ee.extensions = make(map[uint16][]byte, 4)
	}
	if ch.ticketRequest.Enabled {
		ee.extensions[extTicketRequest] = []byte{ticketCount}
	}
	if negotiated != "" {
		ee.extensions[extALPN], err = marshalALPN([]string{negotiated})
		if err != nil {
			return err
		}
	}
	if c.earlyAccepted {
		ee.extensions[extEarlyData] = nil
	}
	if echRejected && len(echKeys) > 0 {
		ee.extensions[extECH], err = buildRetryConfigList(echKeys)
		if err != nil {
			return err
		}
	}
	eeBody, err := ee.marshal()
	if err != nil {
		return err
	}
	serverSequence := uint16(1 + serverSequenceOffset)
	serverMessages := []handshakeMessage{{typ: handshakeTypeEncryptedExtensions, sequence: serverSequence, body: eeBody}}
	_ = transcript.add(handshakeTypeEncryptedExtensions, serverSequence, eeBody)
	serverSequence++
	var clientSignatureSchemes []tls.SignatureScheme
	var clientCertificateSchemes []tls.SignatureScheme
	var clientCertificateCompressionAlgorithms *certificateCompressionAlgorithms
	if !usingPSK && c.config.ClientAuth != tls.NoClientCert {
		request := &certificateRequestMessage{signatureSchemes: defaultSignatureSchemes()}
		if c.config.EnableCertificateCompression {
			clientCertificateCompressionAlgorithms = &certificateCompressionZlibOffer
		}
		clientSignatureSchemes = append([]tls.SignatureScheme(nil), request.signatureSchemes...)
		clientCertificateSchemes = request.certificateSignatureSchemes
		if len(clientCertificateSchemes) == 0 {
			clientCertificateSchemes = append([]tls.SignatureScheme(nil), request.signatureSchemes...)
		}
		requestBody, requestErr := request.marshalWithCertificateCompression(clientCertificateCompressionAlgorithms)
		if requestErr != nil {
			return requestErr
		}
		serverMessages = append(serverMessages, handshakeMessage{typ: handshakeTypeCertificateRequest, sequence: serverSequence, body: requestBody})
		_ = transcript.add(handshakeTypeCertificateRequest, serverSequence, requestBody)
		serverSequence++
	}
	if !usingPSK {
		certMsg := &certificateMessage{}
		for _, der := range cert.Certificate {
			certMsg.certificates = append(certMsg.certificates, certificateEntry{data: der})
		}
		certBody, err := certMsg.marshal()
		if err != nil {
			return err
		}
		certificateType, certificateBody, selectErr := certificateHandshakeMessage(certBody, ch.certificateCompressionAlgorithms(), c.config.EnableCertificateCompression, c.config.certificateCompressionCache())
		if selectErr != nil {
			return selectErr
		}
		serverMessages = append(serverMessages, handshakeMessage{typ: certificateType, sequence: serverSequence, body: certificateBody})
		_ = transcript.add(certificateType, serverSequence, certificateBody)
		serverSequence++
		signature, err := signCertificateVerify(c.config.Rand, signer, scheme, transcript.sumInto(transcriptDigest[:0]), true)
		if err != nil {
			return err
		}
		cvBody, err := (&certificateVerifyMessage{algorithm: scheme, signature: signature}).marshal()
		if err != nil {
			return err
		}
		serverMessages = append(serverMessages, handshakeMessage{typ: handshakeTypeCertificateVerify, sequence: serverSequence, body: cvBody})
		_ = transcript.add(handshakeTypeCertificateVerify, serverSequence, cvBody)
		serverSequence++
	}
	finishedBody := schedule.finishedVerifyData(schedule.serverHandshakeTraffic, transcript.sumInto(transcriptDigest[:0]))
	serverMessages = append(serverMessages, handshakeMessage{typ: handshakeTypeFinished, sequence: serverSequence, body: finishedBody})
	_ = transcript.add(handshakeTypeFinished, serverSequence, finishedBody)
	protected, err := buildProtectedFlight(serverMessages, c.currentMTU(), serverCipher)
	if err != nil {
		return err
	}
	serverFlight := combineFlights(plain, protected)
	c.sendCipher = serverCipher
	if err = c.writeFlight(serverFlightConn, serverFlight); err != nil {
		return err
	}
	if err = schedule.deriveApplication(transcript.sumInto(transcriptDigest[:0])); err != nil {
		return err
	}
	clientFinalFlightStart := clientFinishedSequence
	inbox = newHandshakeInbox(clientFinalFlightStart, c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
	var clientCerts []*x509.Certificate
	var clientChains [][]*x509.Certificate
	var clientAuthAt int64
	if resumedSession != nil {
		clientCerts = resumedSession.peerCertificates
		clientChains = resumedSession.verifiedChains
		clientAuthAt = resumedSession.clientAuthAt
	}
	var clientRecords []recordNumber
	sawClientCertificate := false
	verifiedClientSignature := false
	clientDone := false
	const (
		clientExpectCertificate = iota
		clientExpectCertificateVerify
		clientExpectFinished
		clientHandshakeComplete
	)
	clientStage := clientExpectFinished
	if !usingPSK && c.config.ClientAuth != tls.NoClientCert {
		clientStage = clientExpectCertificate
	}
	for !clientDone {
		receiveConn := c.conn
		if !hrrUsed {
			receiveConn = preValidationConn
		}
		messages, err = c.receiveHandshakeWithRetransmitOnEarly(receiveConn, inbox, clientCipher, serverFlight, earlyCipher, c.queueEarlyApplicationData, serverCipher)
		if err != nil {
			return err
		}
		clientRecords = append(clientRecords, recordNumber{epoch: 2, sequence: clientCipher.lastOpened})
		for index := 0; index < messages.len(); index++ {
			message := messages.at(index)
			switch message.typ {
			case handshakeTypeCertificate, handshakeTypeCompressedCertificate:
				if clientStage != clientExpectCertificate {
					return &ProtocolError{"unexpected client Certificate"}
				}
				certMessage, parseErr := parseCertificateHandshakeMessage(message.typ, message.body, clientCertificateCompressionAlgorithms, c.config.MaxHandshakeMessage)
				if parseErr != nil {
					return parseErr
				}
				if parseErr = validateCertificateMessage(certMessage, nil); parseErr != nil {
					return parseErr
				}
				sawClientCertificate = true
				if len(certMessage.certificates) > 0 {
					clientCerts, clientChains, parseErr = verifyClientCertificate(c.config, certMessage, clientCertificateSchemes)
					if parseErr != nil {
						return parseErr
					}
				}
				if len(certMessage.certificates) > 0 {
					clientStage = clientExpectCertificateVerify
				} else {
					clientStage = clientExpectFinished
				}
				_ = transcript.add(message.typ, message.sequence, message.body)
			case handshakeTypeCertificateVerify:
				if clientStage != clientExpectCertificateVerify || len(clientCerts) == 0 {
					return &ProtocolError{"client CertificateVerify without certificate"}
				}
				cv, parseErr := parseCertificateVerify(message.body)
				if parseErr != nil {
					return parseErr
				}
				offered := false
				for _, scheme := range clientSignatureSchemes {
					offered = offered || scheme == cv.algorithm
				}
				if !offered {
					return alertError(alertIllegalParameter, &ProtocolError{"client selected an unoffered signature scheme"})
				}
				if parseErr = verifyCertificateVerify(clientCerts[0].PublicKey, cv.algorithm, transcript.sumInto(transcriptDigest[:0]), cv.signature, false); parseErr != nil {
					return alertError(alertDecryptError, parseErr)
				}
				verifiedClientSignature = true
				clientStage = clientExpectFinished
				_ = transcript.add(message.typ, message.sequence, message.body)
			case handshakeTypeFinished:
				if clientStage != clientExpectFinished {
					return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected client Finished"})
				}
				clientFinishedSequence = message.sequence
				if !usingPSK && c.config.ClientAuth != tls.NoClientCert && !sawClientCertificate {
					return &ProtocolError{"client omitted Certificate message"}
				}
				required := !usingPSK && (c.config.ClientAuth == tls.RequireAnyClientCert || c.config.ClientAuth == tls.RequireAndVerifyClientCert)
				if required && len(clientCerts) == 0 {
					return alertError(alertCertificateRequired, &ProtocolError{"client certificate is required"})
				}
				if !usingPSK && len(clientCerts) > 0 && !verifiedClientSignature {
					return &ProtocolError{"client omitted CertificateVerify"}
				}
				verify, parseErr := parseFinished(message.body, suite.hash.Size())
				if parseErr != nil {
					return parseErr
				}
				if !schedule.verifyFinished(schedule.clientHandshakeTraffic, transcript.sumInto(transcriptDigest[:0]), verify) {
					return alertError(alertDecryptError, &ProtocolError{"client Finished verification failed"})
				}
				_ = transcript.add(message.typ, message.sequence, message.body)
				clientStage = clientHandshakeComplete
				clientDone = true
			default:
				return &ProtocolError{"unexpected client handshake message"}
			}
		}
	}
	if !hrrUsed {
		// A valid Finished proves address reachability for the no-cookie path.
		amplification.validate()
	}
	ackRecords, _, err := buildACKRecords(clientRecords, c.currentMTU(), 0, serverCipher)
	if err != nil {
		return err
	}
	for _, wire := range ackRecords {
		if err = c.writeRecord(wire); err != nil {
			return err
		}
	}
	if err = c.installApplicationKeysAt(suite, schedule.clientApplicationTraffic, schedule.serverApplicationTraffic, serverSequence+1); err != nil {
		return err
	}
	if err = schedule.deriveResumption(transcript.sumInto(transcriptDigest[:0])); err != nil {
		return err
	}
	c.postHandshakeTranscript = transcript.clone()
	if err = c.receiveEpochs.install(clientCipher); err != nil {
		return err
	}
	if err = c.promoteEarlyApplicationData(); err != nil {
		return err
	}
	c.finishedACKCipher = serverCipher
	c.finishedFlightStart = clientFinalFlightStart
	c.finishedMessageSequence = clientFinishedSequence
	c.completedPeerFlightStart = clientFinalFlightStart
	c.completedPeerFlightEnd = clientFinishedSequence
	c.hasCompletedPeerFlight = true
	exporter := newExporter(suite, schedule.exporterMasterSecret)
	exporter.externalPSK = externalPSK
	c.mu.Lock()
	c.state = ConnectionState{Version: VersionDTLS13, HandshakeComplete: true, DidResume: resumed, ECHAccepted: echAccepted, CipherSuite: suite.id, NegotiatedProtocol: negotiated, PeerCertificates: clientCerts, VerifiedChains: clientChains, LocalConnectionID: append([]byte(nil), c.receiveConnectionID...), PeerConnectionID: append([]byte(nil), c.sendConnectionID...), ReturnRoutabilityCheck: c.returnRoutabilityCheckNegotiated, RecordSizeLimitNegotiated: c.recordSizeLimitNegotiated, LocalRecordSizeLimit: c.localRecordSizeLimit, PeerRecordSizeLimit: c.peerRecordSizeLimit, exporter: exporter}
	c.mu.Unlock()
	if validated, ok := c.conn.(interface{ handshakeValidated() }); ok {
		validated.handshakeValidated()
	}
	if !usingPSK && len(clientCerts) > 0 {
		clientAuthAt = c.config.Time().Unix()
	}
	if err = c.sendNewSessionTickets(schedule, suite, ticketCount, ch.serverName, negotiated, clientAuthAt, clientCerts, clientChains, externalPSK); err != nil {
		return err
	}
	return nil
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
		certificate, err := c.config.GetCertificate(&ClientHelloInfo{ServerName: ch.serverName, SupportedProtos: ch.alpn, Conn: c})
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
	return &c.config.Certificates[0], nil
}
