package dtls13

import (
	"bytes"
	"context"
	"crypto/tls"
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
				if _, ok := errors.AsType[*localAlertError](result); ok {
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
		records := f.nextUnsentWire(storage[:0])
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
		messages, err := receiveHandshakeMessageWithEarlyBatch(conn, inbox, cipher, handshakeReceiveOptions{
			early:     early,
			onEarly:   onEarly,
			outgoing:  outgoing,
			ackCipher: ackCipher,
			mtu:       c.currentMTU(),
			owner:     c,
		})
		if err == nil {
			_ = conn.SetReadDeadline(c.handshakeDeadline)
			return messages, nil
		}
		networkErr, ok := errors.AsType[net.Error](err)
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

func (c *Conn) receiveSecondClientHello(conn net.Conn, inbox *handshakeInbox, hrr *flight) (completedHandshakeBatch, error) {
	buffer := acquireDatagramBuffer()
	defer releaseDatagramBuffer(buffer)
	for {
		datagram, err := readDatagramWithPending(conn, buffer[:], &c.pendingDatagram)
		if err != nil {
			return completedHandshakeBatch{}, err
		}
		for len(datagram) > 0 {
			record, consumed, parseErr := parsePlainRecordView(datagram)
			if parseErr != nil {
				break
			}
			datagram = datagram[consumed:]
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
				c.pendingDatagram = bytes.Clone(datagram)
				return delivered, nil
			}
		}
	}
}

type handshakeReceiveOptions struct {
	early     *recordCipher
	onEarly   func([]byte) error
	outgoing  *flight
	ackCipher *recordCipher
	mtu       int
	owner     *Conn
}

func receiveHandshakeMessageWithEarlyBatch(conn net.Conn, inbox *handshakeInbox, cipher *recordCipher, options handshakeReceiveOptions) (completedHandshakeBatch, error) {
	early, onEarly := options.early, options.onEarly
	outgoing, ackCipher := options.outgoing, options.ackCipher
	mtu, owner := options.mtu, options.owner
	pending := &inbox.pendingDatagram
	if owner != nil {
		pending = &owner.pendingDatagram
	}
	buffer := acquireDatagramBuffer()
	defer releaseDatagramBuffer(buffer)
	for {
		datagram, err := readDatagramWithPending(conn, buffer[:], pending)
		if err != nil {
			return completedHandshakeBatch{}, err
		}
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
				r, recordWireLen, parseErr := parsePlainRecordView(datagram)
				if parseErr != nil {
					break
				}
				typ = r.typ
				payload = r.payload
				consumed = recordWireLen
			} else {
				var openErr error
				if datagram[0] == recordTypeACK {
					r, recordWireLen, parseErr := parsePlainRecordView(datagram)
					if parseErr == nil && r.typ == recordTypeACK {
						consumed = recordWireLen
						numbers, ackErr := parseACK(r.payload)
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
				*pending = bytes.Clone(datagram)
				return delivered, nil
			}
		}
	}
}

func readDatagramWithPending(conn net.Conn, buffer []byte, pending *[]byte) ([]byte, error) {
	if pending != nil && len(*pending) > 0 {
		datagram := *pending
		*pending = nil
		return datagram, nil
	}
	n, err := conn.Read(buffer)
	return buffer[:n], err
}

func receiveACKRecordWithPending(conn net.Conn, dst []recordNumber, pending *[]byte, ciphers ...*recordCipher) ([]recordNumber, error) {
	buffer := acquireDatagramBuffer()
	defer releaseDatagramBuffer(buffer)

	for {
		datagram, err := readDatagramWithPending(conn, buffer[:], pending)
		if err != nil {
			return nil, err
		}
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
					if pending != nil {
						*pending = bytes.Clone(datagram[consumed:])
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
		numbers, err := receiveACKRecordWithPending(c.conn, ackScratch[:0], &c.pendingDatagram, ciphers...)
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
		networkErr, ok := errors.AsType[net.Error](err)
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
		equalStatusRequest(initial, second) &&
		initial.recordSizeLimit == second.recordSizeLimit &&
		initial.hasRecordSizeLimit == second.hasRecordSizeLimit &&
		initial.grease == second.grease &&
		equalExtensionMaps(initial.unknownExtensions, second.unknownExtensions)
}

func equalStatusRequest(initial, second *clientHello) bool {
	if initial.statusRequest != second.statusRequest {
		return false
	}
	if len(initial.statusRequestRaw) == 0 && len(second.statusRequestRaw) == 0 {
		return true
	}
	if len(initial.statusRequestRaw) == 0 {
		return initial.statusRequest && equalBytes(second.statusRequestRaw, marshalStatusRequest())
	}
	if len(second.statusRequestRaw) == 0 {
		return second.statusRequest && equalBytes(initial.statusRequestRaw, marshalStatusRequest())
	}
	return equalBytes(initial.statusRequestRaw, second.statusRequestRaw)
}

func equalByteSlices(left, right [][]byte) bool {
	return slices.EqualFunc(left, right, equalBytes)
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
