package dtls13

import "net"

// This file holds the per-record-type handlers that dispatchDatagramFrom
// delegates to. Every handler runs with dispatchMu held by the caller and owns
// at most one writeMu critical section. Handlers that must act *after*
// releasing writeMu (starting a KeyUpdate retransmission timer) are split into
// a helper that acquires writeMu and returns the decision to a wrapper that acts
// on it, so the lock is never held while a goroutine is started.

// dispatchApplicationData files an application-data record: it is rejected if
// it interleaves a protected handshake message, buffered while a fragmented
// post-handshake message or a post-handshake auth response is incomplete, and
// otherwise queued for ReadDatagram.
func (c *Conn) dispatchApplicationData(content []byte, number recordNumber, from net.Addr) error {
	buffered, err := c.bufferApplicationRecord(content, number, from)
	if err != nil {
		return err
	}
	if buffered {
		return nil
	}
	return c.queueApplicationData(content, from)
}

// bufferApplicationRecord records the application record for ordering
// checks and buffers it if a protected handshake message is still in flight.
// It takes and releases writeMu itself so queueApplicationData can take
// inputMu without writeMu held.
func (c *Conn) bufferApplicationRecord(content []byte, number recordNumber, from net.Addr) (bool, error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := c.rememberApplicationRecordLocked(number); err != nil {
		return false, err
	}
	buffered, err := c.bufferIncompleteHandshakeApplicationLocked(content, number, from)
	if err != nil || buffered {
		return buffered, err
	}
	return c.bufferPostHandshakeAuthApplicationLocked(content, number, from)
}

// dispatchAlert handles an authenticated alert. user_canceled is ignored and
// close_notify half-closes the read side; any other alert is fatal.
func (c *Conn) dispatchAlert(content []byte, number recordNumber) error {
	alert, err := parseAlert(content)
	if err != nil {
		description, _ := protocolAlert(err)
		return alertError(description, err)
	}
	if alert.isUserCanceled() {
		return nil
	}
	if alert.isCloseNotify() {
		c.closure.receive(number)
		c.inputMu.Lock()
		c.peerReadClosed = true
		c.inputMu.Unlock()
		c.notifyRead()
		return nil
	}
	return AlertError(alert.description)
}

// dispatchACK applies an ACK to every outstanding post-handshake flight and,
// when the ACK completes a KeyUpdate that a peer's update_requested was
// waiting on, begins the requested update. The retransmission timer is started
// only after writeMu is released.
func (c *Conn) dispatchACK(content []byte, epoch uint64, onACK func([]recordNumber)) error {
	var scratch [1]recordNumber
	numbers, err := parseACKInto(content, scratch[:0])
	if err != nil {
		description, _ := protocolAlert(err)
		return alertError(description, err)
	}
	if err = validateACKEpoch(numbers, epoch); err != nil {
		return err
	}
	if onACK != nil {
		onACK(numbers)
	}
	startKeyUpdate, err := c.applyACK(numbers)
	if err != nil {
		return err
	}
	if startKeyUpdate {
		c.startKeyUpdateRetransmission()
	}
	return nil
}

// applyACK is the writeMu critical section of dispatchACK. It returns
// whether a deferred KeyUpdate response was just sent and now needs its
// retransmission timer.
func (c *Conn) applyACK(numbers []recordNumber) (startKeyUpdate bool, err error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.sendingTraffic != nil && c.sendingTraffic.processACK(numbers) {
		c.sendCipher = c.sendingTraffic.cipher
		switch {
		case c.keyUpdateResponsePending && c.sendingTraffic.canBeginKeyUpdate():
			wire, _, beginErr := c.sendingTraffic.beginKeyUpdate(false)
			if beginErr != nil {
				return false, beginErr
			}
			if beginErr = c.writeRecord(wire); beginErr != nil {
				return false, beginErr
			}
			c.keyUpdateResponsePending = false
			startKeyUpdate = true
		case c.keyUpdateResponsePending && c.sendingTraffic.cipher.epoch >= maxSendingEpoch:
			c.keyUpdateResponsePending = false
		}
	}
	for _, flight := range []**flight{&c.ticketFlight, &c.clientAuthRequestFlight, &c.clientAuthResponseFlight} {
		if err = c.ackPostHandshakeFlightLocked(flight, numbers); err != nil {
			return false, err
		}
	}
	if err = c.processCIDACKsLocked(numbers); err != nil {
		return false, err
	}
	return startKeyUpdate, nil
}

// ackPostHandshakeFlightLocked applies numbers to one outstanding
// post-handshake flight, clearing it when complete and retransmitting the
// unacknowledged remainder otherwise. Requires writeMu.
func (c *Conn) ackPostHandshakeFlightLocked(slot **flight, numbers []recordNumber) error {
	f := *slot
	if f == nil {
		return nil
	}
	f.ack(numbers)
	if f.complete() {
		c.observeFlightRTT(f)
		*slot = nil
		return nil
	}
	return c.retransmitPartialFlight(c.conn, f)
}

// dispatchHandshake handles an authenticated handshake record after the
// handshake has completed. A retransmission of the peer's final handshake
// flight is acknowledged and dropped; everything else is a post-handshake
// message routed by fragment type.
func (c *Conn) dispatchHandshake(content []byte, number recordNumber, epoch uint64) error {
	var fragmentScratch [1]handshakeFragment
	fragments, err := parseHandshakeFragmentsViewInto(content, fragmentScratch[:0])
	if err != nil {
		return err
	}
	for _, fragment := range fragments {
		if fragment.typ == handshakeTypeKeyUpdate && (len(fragments) != 1 || fragment.offset != 0 || int(fragment.length) != len(fragment.body)) {
			return alertError(alertUnexpectedMessage, &ProtocolError{"KeyUpdate is not aligned to a record boundary"})
		}
	}
	if epoch == 2 && c.hasCompletedPeerFlight {
		if c.isCompletedPeerFlightRetransmit(fragments) {
			c.ackCompletedPeerFlightRecord(number)
		}
		return nil
	}
	if c.receivingTraffic == nil {
		return nil
	}
	for _, fragment := range fragments {
		if err = c.dispatchPostHandshakeFragment(fragment, number); err != nil {
			return err
		}
	}
	return nil
}

// isCompletedPeerFlightRetransmit reports whether every fragment falls inside
// the message-sequence range of the peer's already-completed final flight.
func (c *Conn) isCompletedPeerFlightRetransmit(fragments []handshakeFragment) bool {
	if len(fragments) == 0 {
		return false
	}
	for _, fragment := range fragments {
		if fragment.messageSequence < c.completedPeerFlightStart || fragment.messageSequence > c.completedPeerFlightEnd {
			return false
		}
	}
	return true
}

// ackCompletedPeerFlightRecord re-acknowledges a record of the peer's
// completed final flight so a peer whose ACK was lost stops retransmitting.
// Write errors are deliberately ignored: the connection is already up and the
// peer will simply retransmit again.
func (c *Conn) ackCompletedPeerFlightRecord(number recordNumber) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	var ackScratch [1][]byte
	acks, _, err := buildACKRecordsInto(ackScratch[:0], []recordNumber{number}, c.currentMTU(), 0, c.sendCipher)
	if err != nil {
		return
	}
	for _, wire := range acks {
		_, _ = c.conn.Write(wire)
	}
}

// dispatchPostHandshakeFragment routes one post-handshake fragment by type and
// role. Fragments that are not valid post-handshake messages are fatal.
func (c *Conn) dispatchPostHandshakeFragment(fragment handshakeFragment, number recordNumber) error {
	switch {
	case fragment.typ == handshakeTypeCertificateRequest && c.isClient:
		return c.dispatchPostHandshakeCertificateRequest(fragment, number)
	case !c.isClient && (fragment.typ == handshakeTypeCertificate || fragment.typ == handshakeTypeCompressedCertificate || fragment.typ == handshakeTypeCertificateVerify || fragment.typ == handshakeTypeFinished):
		return c.processPostHandshakeAuthFragment(fragment, number)
	case fragment.typ == handshakeTypeNewSessionTicket && c.isClient:
		return c.dispatchNewSessionTicket(fragment, number)
	case fragment.typ == handshakeTypeNewConnectionID || fragment.typ == handshakeTypeRequestConnectionID:
		return c.dispatchConnectionIDMessage(fragment, number)
	case fragment.typ == handshakeTypeKeyUpdate && fragment.offset == 0 && int(fragment.length) == len(fragment.body):
		return c.dispatchKeyUpdate(fragment, number)
	default:
		return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected post-handshake message"})
	}
}

// reassembleProtectedHandshake adds a fragment to the post-handshake
// reassembler, creating it on first use. When the message completes it records
// the record range so interleaved application data is detected and any data
// buffered during reassembly is released. The reassembler is guarded by
// dispatchMu, which every caller holds.
func (c *Conn) reassembleProtectedHandshake(fragment handshakeFragment, number recordNumber) (body []byte, complete bool, err error) {
	if c.postHandshakeReassembly == nil {
		c.postHandshakeReassembly = newReassemblerWithLimits(c.config.MaxHandshakeMessage, c.config.MaxBufferedHandshakeMessages, c.config.MaxBufferedHandshakeBytes)
	}
	body, complete, firstRecord, lastRecord, err := c.postHandshakeReassembly.addProtectedRecord(fragment, number)
	if err != nil || !complete {
		return nil, false, err
	}
	if err = c.rememberProtectedHandshakeRange(firstRecord, lastRecord); err != nil {
		return nil, false, err
	}
	return body, true, nil
}

// dispatchPostHandshakeCertificateRequest acknowledges the record, reassembles
// the server's post-handshake CertificateRequest, and answers it.
func (c *Conn) dispatchPostHandshakeCertificateRequest(fragment handshakeFragment, number recordNumber) error {
	if err := c.ackProtectedRecord(number); err != nil {
		return err
	}
	body, complete, err := c.reassembleProtectedHandshake(fragment, number)
	if err != nil || !complete {
		return err
	}
	return c.processPostHandshakeCertificateRequest(fragment.messageSequence, body)
}

// dispatchNewSessionTicket acknowledges the record, reassembles the ticket,
// and stores it.
func (c *Conn) dispatchNewSessionTicket(fragment handshakeFragment, number recordNumber) error {
	if err := c.ackProtectedRecord(number); err != nil {
		return err
	}
	body, complete, err := c.reassembleProtectedHandshake(fragment, number)
	if err != nil || !complete {
		return err
	}
	return c.processNewSessionTicket(fragment.messageSequence, body)
}

// dispatchConnectionIDMessage reassembles a NewConnectionId or
// RequestConnectionId message, applies it, and then acknowledges the record.
// Unlike the ticket and certificate-request paths the ACK follows processing,
// so a message the peer must not see acknowledged is never acknowledged.
func (c *Conn) dispatchConnectionIDMessage(fragment handshakeFragment, number recordNumber) error {
	body, complete, err := c.reassembleProtectedHandshake(fragment, number)
	if err != nil {
		return err
	}
	var requestCount uint8
	respond := false
	if complete {
		if fragment.typ == handshakeTypeNewConnectionID {
			err = c.processNewConnectionID(fragment.messageSequence, body)
		} else {
			requestCount, respond, err = c.processRequestConnectionID(fragment.messageSequence, body)
		}
		if err != nil {
			return err
		}
	}
	if err = c.ackProtectedRecord(number); err != nil {
		return err
	}
	if respond {
		go c.respondToConnectionIDRequest(requestCount)
	}
	return nil
}

// dispatchKeyUpdate installs the peer's new receive epoch, acknowledges the
// KeyUpdate, and, if the peer asked for one, begins our own update. The
// retransmission timer is started only after writeMu is released.
func (c *Conn) dispatchKeyUpdate(fragment handshakeFragment, number recordNumber) error {
	startRetransmission, err := c.processKeyUpdateFragment(fragment, number)
	if err != nil {
		return err
	}
	if startRetransmission {
		c.startKeyUpdateRetransmission()
	}
	return nil
}

// processKeyUpdateFragment is the writeMu critical section of
// dispatchKeyUpdate. It returns whether a KeyUpdate response was sent and
// needs its retransmission timer.
func (c *Conn) processKeyUpdateFragment(fragment handshakeFragment, number recordNumber) (startRetransmission bool, err error) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if c.postHandshakeAuthState != nil && c.postHandshakeAuthState.hasResponseEpoch {
		return false, alertError(alertUnexpectedMessage, &ProtocolError{"KeyUpdate interleaved with post-handshake authentication response"})
	}
	if c.hasIncompleteProtectedHandshakeLocked() {
		return false, alertError(alertUnexpectedMessage, &ProtocolError{"KeyUpdate followed an incomplete handshake message"})
	}
	message, updated, err := c.receivingTraffic.processKeyUpdate(fragment.messageSequence, fragment.body)
	if err != nil {
		return false, err
	}
	var ackScratch [1][]byte
	acks, _, err := buildACKRecordsInto(ackScratch[:0], []recordNumber{number}, c.currentMTU(), 0, c.sendCipher)
	if err != nil {
		return false, err
	}
	for _, wire := range acks {
		if err = c.writeRecord(wire); err != nil {
			return false, err
		}
	}
	if !updated || !message.requestUpdate {
		return false, nil
	}
	if !c.sendingTraffic.canBeginKeyUpdate() {
		if c.sendingTraffic.cipher.epoch < maxSendingEpoch {
			c.keyUpdateResponsePending = true
		}
		return false, nil
	}
	wire, _, err := c.sendingTraffic.beginKeyUpdate(false)
	if err != nil {
		return false, err
	}
	c.sendCipher = c.sendingTraffic.cipher
	if err = c.writeRecord(wire); err != nil {
		return false, err
	}
	return true, nil
}
