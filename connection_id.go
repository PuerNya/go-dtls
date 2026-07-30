package dtls13

import (
	"errors"
	"time"
)

const (
	handshakeTypeRequestConnectionID uint8 = 9
	handshakeTypeNewConnectionID     uint8 = 10

	connectionIDImmediate uint8 = 0
	connectionIDSpare     uint8 = 1
)

type newConnectionIDMessage struct {
	connectionIDs [][]byte
	usage         uint8
}

func (m newConnectionIDMessage) marshal() ([]byte, error) {
	if m.usage > connectionIDSpare {
		return nil, &ProtocolError{"invalid NewConnectionId usage"}
	}
	if m.usage == connectionIDImmediate && len(m.connectionIDs) == 0 {
		return nil, &ProtocolError{"immediate NewConnectionId has no connection IDs"}
	}
	encodedLength := 0
	for _, cid := range m.connectionIDs {
		if len(cid) > 255 {
			return nil, &ProtocolError{"8-bit vector overflow"}
		}
		if encodedLength > 65535-1-len(cid) {
			return nil, &ProtocolError{"16-bit vector overflow"}
		}
		encodedLength += 1 + len(cid)
	}
	w := newWireBuilder(2 + encodedLength + 1)
	start := w.startVector16()
	for _, cid := range m.connectionIDs {
		w.bytes8(cid)
	}
	w.endVector16(start)
	w.u8(int(m.usage))
	return w.b, w.err
}

func parseNewConnectionID(b []byte) (newConnectionIDMessage, error) {
	p := wireParser{b: b}
	encoded := p.bytes16()
	usage := uint8(p.u8())
	if err := p.done(); err != nil {
		return newConnectionIDMessage{}, err
	}
	if usage > connectionIDSpare {
		return newConnectionIDMessage{}, &ProtocolError{"invalid NewConnectionId usage"}
	}
	var connectionIDs [][]byte
	q := wireParser{b: encoded}
	for q.off < len(q.b) {
		cid := q.bytes8()
		if q.err != nil {
			return newConnectionIDMessage{}, q.err
		}
		connectionIDs = append(connectionIDs, append([]byte{}, cid...))
	}
	if usage == connectionIDImmediate && len(connectionIDs) == 0 {
		return newConnectionIDMessage{}, &ProtocolError{"immediate NewConnectionId has no connection IDs"}
	}
	return newConnectionIDMessage{connectionIDs: connectionIDs, usage: usage}, nil
}

type requestConnectionIDMessage struct{ count uint8 }

func (m requestConnectionIDMessage) marshal() []byte { return []byte{m.count} }

func parseRequestConnectionID(b []byte) (requestConnectionIDMessage, error) {
	if len(b) != 1 {
		return requestConnectionIDMessage{}, alertError(alertDecodeError, &ProtocolError{"invalid RequestConnectionId length"})
	}
	return requestConnectionIDMessage{count: b[0]}, nil
}

type cidFlightKind uint8

const (
	newCIDFlight cidFlightKind = iota
	requestCIDFlight
)

// SendNewConnectionIDs reliably advertises new local Connection IDs to the
// peer. Each ID is a value the peer may place in records sent to this endpoint.
// IDs must be at most 255 bytes and, together with existing local IDs, unique
// and prefix-free. Their total number is bounded by Config.MaxConnectionIDs.
//
// When immediate is true, connectionIDs must be non-empty and the peer is
// instructed to switch immediately to one of the supplied IDs. Otherwise the
// IDs become spares. Only one locally initiated NewConnectionId flight may
// await acknowledgement at a time. The method fails if non-empty Connection
// IDs and local CID updates were not negotiated.
func (c *Conn) SendNewConnectionIDs(connectionIDs [][]byte, immediate bool) error {
	if err := c.Handshake(); err != nil {
		return err
	}
	usage := uint8(connectionIDSpare)
	if immediate {
		usage = connectionIDImmediate
	}
	message := newConnectionIDMessage{connectionIDs: connectionIDs, usage: usage}
	body, err := message.marshal()
	if err != nil {
		return err
	}

	c.writeMu.Lock()
	if !c.connectionIDNegotiated || !c.localCIDUpdatesAllowed {
		c.writeMu.Unlock()
		return &ProtocolError{"NewConnectionId is not permitted for this connection"}
	}
	if c.newConnectionIDFlight != nil && !c.newConnectionIDFlight.complete() {
		c.writeMu.Unlock()
		return errors.New("dtls13: NewConnectionId is already awaiting acknowledgement")
	}
	if c.receivingTraffic == nil || c.sendingTraffic == nil {
		c.writeMu.Unlock()
		return &ProtocolError{"application traffic state is not installed"}
	}
	if err = c.sendingTraffic.canAllocateMessageSequences(1); err != nil {
		c.writeMu.Unlock()
		return err
	}
	merged, mergeErr := mergeConnectionIDs(c.receivingTraffic.acceptedCIDs, connectionIDs)
	if mergeErr != nil {
		c.writeMu.Unlock()
		return mergeErr
	}
	if len(merged) > c.config.MaxConnectionIDs {
		c.writeMu.Unlock()
		return &ConfigError{"NewConnectionId exceeds MaxConnectionIDs"}
	}
	sequence := c.sendingTraffic.messageSequence
	flight, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeNewConnectionID, sequence: sequence, body: body}}, c.currentMTU(), c.sendCipher)
	if err != nil {
		c.writeMu.Unlock()
		return err
	}
	flight.setIntervals(c.config.FlightInterval, c.config.MaxFlightInterval)

	previousCIDs := cloneConnectionIDs(c.receivingTraffic.acceptedCIDs)
	newCIDs := connectionIDDifference(merged, previousCIDs)
	if err = c.receivingTraffic.addConnectionIDs(connectionIDs); err != nil {
		c.writeMu.Unlock()
		return err
	}
	type cidRegistrar interface {
		registerConnectionIDs([][]byte) error
		unregisterConnectionIDs([][]byte)
	}
	registrar, hasRegistrar := c.conn.(cidRegistrar)
	if hasRegistrar {
		if err = registrar.registerConnectionIDs(newCIDs); err != nil {
			c.receivingTraffic.restoreConnectionIDs(previousCIDs)
			c.writeMu.Unlock()
			return err
		}
	}
	err = c.writeFlight(c.conn, flight)
	if err != nil {
		if hasRegistrar {
			registrar.unregisterConnectionIDs(newCIDs)
		}
		c.receivingTraffic.restoreConnectionIDs(previousCIDs)
		c.writeMu.Unlock()
		return err
	}
	c.sendingTraffic.commitMessageSequences(1)
	c.newConnectionIDFlight = flight
	c.writeMu.Unlock()
	c.startCIDFlightRetransmission(newCIDFlight)
	return nil
}

func connectionIDDifference(connectionIDs, existing [][]byte) [][]byte {
	var difference [][]byte
	for _, cid := range connectionIDs {
		found := false
		for _, current := range existing {
			if equalBytes(cid, current) {
				found = true
				break
			}
		}
		if !found {
			difference = append(difference, append([]byte{}, cid...))
		}
	}
	return difference
}

// RequestConnectionIDs reliably asks the peer to advertise up to count spare
// Connection IDs. It is valid only when a non-empty CID is active for records
// sent to the peer. A new request cannot be sent until the preceding request
// has been acknowledged and fulfilled.
//
// The peer controls how many IDs it returns, subject to its
// Config.MaxConnectionIDs and Config.GetConnectionID policy.
func (c *Conn) RequestConnectionIDs(count uint8) error {
	if err := c.Handshake(); err != nil {
		return err
	}
	body := (requestConnectionIDMessage{count: count}).marshal()
	c.writeMu.Lock()
	if !c.connectionIDNegotiated || len(c.sendConnectionID) == 0 {
		c.writeMu.Unlock()
		return &ProtocolError{"RequestConnectionId cannot be sent with an empty connection ID"}
	}
	if c.connectionIDRequestOpen {
		c.writeMu.Unlock()
		return errors.New("dtls13: RequestConnectionId is still unfulfilled")
	}
	if c.requestCIDFlight != nil && !c.requestCIDFlight.complete() {
		c.writeMu.Unlock()
		return errors.New("dtls13: RequestConnectionId is already awaiting acknowledgement")
	}
	if c.sendingTraffic == nil {
		c.writeMu.Unlock()
		return &ProtocolError{"application traffic state is not installed"}
	}
	if err := c.sendingTraffic.canAllocateMessageSequences(1); err != nil {
		c.writeMu.Unlock()
		return err
	}
	sequence := c.sendingTraffic.messageSequence
	flight, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeRequestConnectionID, sequence: sequence, body: body}}, c.currentMTU(), c.sendCipher)
	if err == nil {
		flight.setIntervals(c.config.FlightInterval, c.config.MaxFlightInterval)
		err = c.writeFlight(c.conn, flight)
	}
	if err != nil {
		c.writeMu.Unlock()
		return err
	}
	c.sendingTraffic.commitMessageSequences(1)
	c.requestCIDFlight = flight
	c.connectionIDRequestOpen = true
	c.writeMu.Unlock()
	c.startCIDFlightRetransmission(requestCIDFlight)
	return nil
}

// UseNextConnectionID switches outgoing protected records to the next spare
// CID provided by the peer. It returns an error if no spare CID is available.
//
// Switching a CID changes record routing but does not validate, select, or
// rebind a network path. Applications implementing migration need a separate
// path-validation policy.
func (c *Conn) UseNextConnectionID() error {
	if err := c.Handshake(); err != nil {
		return err
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if len(c.peerSpareConnectionIDs) == 0 {
		return errors.New("dtls13: peer has not provided a spare connection ID")
	}
	cid := c.peerSpareConnectionIDs[0]
	c.peerSpareConnectionIDs = c.peerSpareConnectionIDs[1:]
	return c.activatePeerConnectionIDLocked(cid)
}

func (c *Conn) activatePeerConnectionIDLocked(cid []byte) error {
	if c.sendingTraffic == nil {
		return &ProtocolError{"application traffic state is not installed"}
	}
	if err := c.sendingTraffic.setConnectionID(cid); err != nil {
		return err
	}
	c.sendConnectionID = append([]byte{}, cid...)
	c.sendCipher = c.sendingTraffic.cipher
	c.mu.Lock()
	c.state.PeerConnectionID = append([]byte{}, cid...)
	c.mu.Unlock()
	return nil
}

func (c *Conn) processNewConnectionID(sequence uint16, body []byte) error {
	message, err := parseNewConnectionID(body)
	if err != nil {
		return err
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if !c.connectionIDNegotiated || !c.peerCIDUpdatesAllowed {
		return alertError(alertUnexpectedMessage, &ProtocolError{"unexpected NewConnectionId"})
	}
	if c.hasNewCIDSequence && sequence <= c.lastNewCIDSequence {
		return nil
	}
	existing := make([][]byte, 0, 1+len(c.peerSpareConnectionIDs))
	existing = append(existing, append([]byte(nil), c.sendConnectionID...))
	existing = append(existing, c.peerSpareConnectionIDs...)
	_, mergeErr := mergeConnectionIDs(existing, message.connectionIDs)
	if mergeErr != nil {
		return alertError(alertIllegalParameter, mergeErr)
	}
	c.lastNewCIDSequence = sequence
	c.hasNewCIDSequence = true
	c.connectionIDRequestOpen = false
	for _, cid := range message.connectionIDs {
		duplicate := equalBytes(cid, c.sendConnectionID)
		for _, existing := range c.peerSpareConnectionIDs {
			duplicate = duplicate || equalBytes(cid, existing)
		}
		if !duplicate && 1+len(c.peerSpareConnectionIDs) < c.config.MaxConnectionIDs {
			c.peerSpareConnectionIDs = append(c.peerSpareConnectionIDs, append([]byte{}, cid...))
		}
	}
	if message.usage == connectionIDImmediate {
		cid := message.connectionIDs[0]
		for i, spare := range c.peerSpareConnectionIDs {
			if equalBytes(spare, cid) {
				c.peerSpareConnectionIDs = append(c.peerSpareConnectionIDs[:i], c.peerSpareConnectionIDs[i+1:]...)
				break
			}
		}
		if err = c.activatePeerConnectionIDLocked(cid); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conn) processRequestConnectionID(sequence uint16, body []byte) (uint8, bool, error) {
	message, err := parseRequestConnectionID(body)
	if err != nil {
		return 0, false, err
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	c.mu.RLock()
	sendingEmpty := len(c.state.LocalConnectionID) == 0
	c.mu.RUnlock()
	if !c.connectionIDNegotiated || sendingEmpty {
		return 0, false, alertError(alertUnexpectedMessage, &ProtocolError{"unexpected RequestConnectionId"})
	}
	if c.hasRequestCIDSequence && sequence <= c.lastRequestCIDSequence {
		return message.count, false, nil
	}
	c.lastRequestCIDSequence = sequence
	c.hasRequestCIDSequence = true
	return message.count, true, nil
}

func (c *Conn) respondToConnectionIDRequest(count uint8) {
	c.cidGenMu.Lock()
	defer c.cidGenMu.Unlock()
	c.writeMu.Lock()
	available := c.config.MaxConnectionIDs
	if c.receivingTraffic != nil {
		available -= len(c.receivingTraffic.acceptedCIDs)
	}
	c.writeMu.Unlock()
	if available < 0 {
		available = 0
	}
	requested := int(count)
	if requested > available {
		requested = available
	}
	connectionIDs := make([][]byte, 0, requested)
	for i := 0; i < requested && c.config.GetConnectionID != nil; i++ {
		cid, err := c.config.GetConnectionID()
		if err != nil {
			break
		}
		connectionIDs = append(connectionIDs, append([]byte{}, cid...))
	}
	_ = c.SendNewConnectionIDs(connectionIDs, false)
}

func (c *Conn) cidFlightLocked(kind cidFlightKind) **flight {
	if kind == newCIDFlight {
		return &c.newConnectionIDFlight
	}
	return &c.requestCIDFlight
}

func (c *Conn) processCIDACKsLocked(numbers []recordNumber) error {
	for _, kind := range []cidFlightKind{newCIDFlight, requestCIDFlight} {
		flightPointer := c.cidFlightLocked(kind)
		if *flightPointer == nil {
			continue
		}
		(*flightPointer).ack(numbers)
		if (*flightPointer).complete() {
			c.observeFlightRTT(*flightPointer)
			*flightPointer = nil
		} else if err := c.retransmitPartialFlight(c.conn, *flightPointer); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conn) startCIDFlightRetransmission(kind cidFlightKind) {
	go func() {
		interval := c.flightInterval()
		if interval <= 0 {
			interval = time.Second
		}
		max := c.config.MaxFlightInterval
		if max < interval {
			max = interval
		}
		timer := time.NewTimer(interval)
		defer timer.Stop()
		timeoutCount := 0
		for {
			<-timer.C
			c.readerMu.Lock()
			closed := c.readerClosed
			c.readerMu.Unlock()
			if closed {
				return
			}
			c.writeMu.Lock()
			flightPointer := c.cidFlightLocked(kind)
			flight := *flightPointer
			if flight == nil || flight.complete() {
				*flightPointer = nil
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
