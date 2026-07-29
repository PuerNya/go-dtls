package dtls13

import (
	"io"
	"net"
	"time"
)

const (
	returnRoutabilityPathChallenge uint8 = iota
	returnRoutabilityPathResponse
	returnRoutabilityPathDrop
	returnRoutabilityMessageLen = 9
)

type returnRoutabilityMessage struct {
	typ    uint8
	cookie [8]byte
}

func (m returnRoutabilityMessage) marshal() ([returnRoutabilityMessageLen]byte, error) {
	var wire [returnRoutabilityMessageLen]byte
	if m.typ > returnRoutabilityPathDrop {
		return wire, &ProtocolError{"invalid return routability message type"}
	}
	wire[0] = m.typ
	copy(wire[1:], m.cookie[:])
	return wire, nil
}

func parseReturnRoutabilityMessage(wire []byte) (returnRoutabilityMessage, bool, error) {
	var message returnRoutabilityMessage
	if len(wire) == 0 {
		return message, false, &ProtocolError{"truncated return routability message"}
	}
	message.typ = wire[0]
	if message.typ > returnRoutabilityPathDrop {
		return message, false, nil
	}
	if len(wire) != returnRoutabilityMessageLen {
		return message, false, &ProtocolError{"invalid return routability message length"}
	}
	copy(message.cookie[:], wire[1:])
	return message, true, nil
}

type returnRoutabilityPhase uint8

const (
	returnRoutabilityIdle returnRoutabilityPhase = iota
	returnRoutabilityCheckingOldPath
	returnRoutabilityCheckingNewPath
)

type returnRoutabilityState struct {
	phase                  returnRoutabilityPhase
	generation             uint64
	timer                  *time.Timer
	oldAddress, newAddress net.Addr
	cookie                 [8]byte
	probeConnectionID      []byte
	received, sent         uint64
}

type returnRoutabilityWriter interface {
	WriteTo([]byte, net.Addr) (int, error)
}

type returnRoutabilityRebinder interface {
	rebindRemote(net.Addr)
}

func sameNetworkAddress(left, right net.Addr) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	leftUDP, leftOK := left.(*net.UDPAddr)
	rightUDP, rightOK := right.(*net.UDPAddr)
	if leftOK && rightOK {
		return leftUDP.Port == rightUDP.Port && leftUDP.Zone == rightUDP.Zone && leftUDP.IP.Equal(rightUDP.IP)
	}
	return left.Network() == right.Network() && left.String() == right.String()
}

func (c *Conn) returnRoutabilityTimeout() time.Duration {
	if c.retransmitNanos.Load() <= 0 {
		return time.Second
	}
	interval := c.flightInterval()
	if c.retransmitNanos.Load() <= 0 {
		return time.Second
	}
	const maxDuration = time.Duration(1<<63 - 1)
	if interval > maxDuration/2 {
		return maxDuration
	}
	// flightInterval is 1.5x the sampled RTT, so twice it is RFC 9853's 3xRTT.
	return 2 * interval
}

func (c *Conn) observeReturnRoutabilityRecord(from net.Addr, received int) error {
	if !c.returnRoutabilityCheckNegotiated || from == nil || received <= 0 {
		return nil
	}
	oldAddress := c.RemoteAddr()
	if oldAddress == nil || sameNetworkAddress(from, oldAddress) {
		return nil
	}
	if _, ok := c.conn.(returnRoutabilityWriter); !ok {
		return nil
	}
	if _, ok := c.conn.(returnRoutabilityRebinder); !ok {
		return nil
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	state := c.returnRoutability
	if state == nil {
		state = new(returnRoutabilityState)
		c.returnRoutability = state
	}
	if state.phase != returnRoutabilityIdle {
		if sameNetworkAddress(from, state.newAddress) {
			state.received = addSaturated(state.received, uint64(received))
		}
		return nil
	}
	if _, err := io.ReadFull(c.config.Rand, state.cookie[:]); err != nil {
		return err
	}
	state.phase = returnRoutabilityCheckingOldPath
	state.oldAddress = oldAddress
	state.newAddress = from
	state.received = uint64(received)
	state.sent = 0
	if err := c.sendReturnRoutabilityLocked(oldAddress, returnRoutabilityPathChallenge, state.cookie); err != nil {
		c.clearReturnRoutabilityLocked()
		return err
	}
	c.armReturnRoutabilityTimerLocked()
	return nil
}

func (c *Conn) handleReturnRoutability(content []byte, from net.Addr) error {
	if !c.returnRoutabilityCheckNegotiated || from == nil {
		return nil
	}
	message, known, err := parseReturnRoutabilityMessage(content)
	if err != nil || !known {
		return nil
	}
	if message.typ == returnRoutabilityPathChallenge {
		responseType := returnRoutabilityPathDrop
		if sameNetworkAddress(from, c.RemoteAddr()) {
			responseType = returnRoutabilityPathResponse
		}
		c.writeMu.Lock()
		err = c.sendReturnRoutabilityLocked(from, responseType, message.cookie)
		c.writeMu.Unlock()
		return err
	}

	c.writeMu.Lock()
	state := c.returnRoutability
	if state == nil {
		c.writeMu.Unlock()
		return nil
	}
	switch state.phase {
	case returnRoutabilityCheckingOldPath:
		if sameNetworkAddress(from, state.oldAddress) && equalBytes(message.cookie[:], state.cookie[:]) {
			switch message.typ {
			case returnRoutabilityPathResponse:
				c.clearReturnRoutabilityLocked()
			case returnRoutabilityPathDrop:
				err = c.startBasicReturnRoutabilityLocked()
			}
		}
	case returnRoutabilityCheckingNewPath:
		if message.typ == returnRoutabilityPathResponse && sameNetworkAddress(from, state.newAddress) && equalBytes(message.cookie[:], state.cookie[:]) {
			rebindAddress := state.newAddress
			probeConnectionID := state.probeConnectionID
			c.clearReturnRoutabilityLocked()
			err = c.rebindValidatedPathLocked(rebindAddress, probeConnectionID)
		}
	}
	c.writeMu.Unlock()
	return err
}

func (c *Conn) rebindValidatedPathLocked(address net.Addr, probeConnectionID []byte) error {
	cidIndex := -1
	if len(probeConnectionID) > 0 && !equalBytes(probeConnectionID, c.sendConnectionID) {
		for i, cid := range c.peerSpareConnectionIDs {
			if equalBytes(cid, probeConnectionID) {
				cidIndex = i
				break
			}
		}
	} else if len(probeConnectionID) == 0 && len(c.peerSpareConnectionIDs) > 0 {
		cidIndex = 0
	}
	if cidIndex >= 0 {
		cid := c.peerSpareConnectionIDs[cidIndex]
		if err := c.activatePeerConnectionIDLocked(cid); err != nil {
			return err
		}
		c.peerSpareConnectionIDs = append(c.peerSpareConnectionIDs[:cidIndex], c.peerSpareConnectionIDs[cidIndex+1:]...)
	}
	c.conn.(returnRoutabilityRebinder).rebindRemote(address)
	return nil
}

func (c *Conn) startBasicReturnRoutabilityLocked() error {
	state := c.returnRoutability
	if state == nil || state.phase == returnRoutabilityIdle || state.newAddress == nil {
		return nil
	}
	if _, err := io.ReadFull(c.config.Rand, state.cookie[:]); err != nil {
		c.clearReturnRoutabilityLocked()
		return err
	}
	state.phase = returnRoutabilityCheckingNewPath
	if err := c.sendReturnRoutabilityLocked(state.newAddress, returnRoutabilityPathChallenge, state.cookie); err != nil {
		c.clearReturnRoutabilityLocked()
		return err
	}
	c.armReturnRoutabilityTimerLocked()
	return nil
}

func (c *Conn) sendReturnRoutabilityLocked(address net.Addr, typ uint8, cookie [8]byte) error {
	payload, err := (returnRoutabilityMessage{typ: typ, cookie: cookie}).marshal()
	if err != nil {
		return err
	}
	if c.sendCipher == nil {
		return &ProtocolError{"return routability traffic keys are not installed"}
	}
	activeConnectionID := c.sendCipher.connectionID
	state := c.returnRoutability
	if typ == returnRoutabilityPathChallenge && state != nil && state.phase == returnRoutabilityCheckingNewPath && sameNetworkAddress(address, state.newAddress) && len(c.peerSpareConnectionIDs) > 0 {
		state.probeConnectionID = c.peerSpareConnectionIDs[0]
		c.sendCipher.connectionID = state.probeConnectionID
	}
	wireLength := c.sendCipher.headerLen16() + c.sendCipher.aead.Overhead() + 1 + len(payload)
	if !c.allowReturnRoutabilitySendLocked(address, wireLength) {
		c.sendCipher.connectionID = activeConnectionID
		return nil
	}
	wire, err := c.sendCipher.seal(recordTypeReturnRoutability, payload[:])
	c.sendCipher.connectionID = activeConnectionID
	if err == nil {
		var n int
		n, err = c.writeRecordTo(wire, address)
		if err == nil && n != len(wire) {
			err = io.ErrShortWrite
		}
		if err == nil {
			c.recordReturnRoutabilitySendLocked(address, n)
		}
	}
	startUpdate := false
	if err == nil {
		startUpdate, err = c.maybeStartAutomaticKeyUpdateLocked()
	}
	if startUpdate {
		go c.startKeyUpdateRetransmission()
	}
	return err
}

func (c *Conn) writeRecordTo(wire []byte, address net.Addr) (int, error) {
	var n int
	var err error
	if sameNetworkAddress(address, c.RemoteAddr()) {
		n, err = c.conn.Write(wire)
	} else if writer, ok := c.conn.(returnRoutabilityWriter); ok {
		n, err = writer.WriteTo(wire, address)
	} else {
		return 0, &net.OpError{Op: "write", Net: "dtls", Addr: address, Err: &ProtocolError{"transport cannot write to the return routability address"}}
	}
	return n, normalizeDatagramWriteError(err, address)
}

func (c *Conn) allowReturnRoutabilitySendLocked(address net.Addr, size int) bool {
	state := c.returnRoutability
	if size < 0 || state == nil || state.phase == returnRoutabilityIdle || !sameNetworkAddress(address, state.newAddress) {
		return size >= 0
	}
	budget := state.received
	if budget > ^uint64(0)/3 {
		budget = ^uint64(0)
	} else {
		budget *= 3
	}
	return state.sent <= budget && uint64(size) <= budget-state.sent
}

func (c *Conn) recordReturnRoutabilitySendLocked(address net.Addr, size int) {
	state := c.returnRoutability
	if size > 0 && state != nil && state.phase != returnRoutabilityIdle && sameNetworkAddress(address, state.newAddress) {
		state.sent = addSaturated(state.sent, uint64(size))
	}
}

func addSaturated(left, right uint64) uint64 {
	if ^uint64(0)-left < right {
		return ^uint64(0)
	}
	return left + right
}

func (c *Conn) armReturnRoutabilityTimerLocked() {
	state := c.returnRoutability
	if state == nil {
		return
	}
	if state.timer != nil {
		state.timer.Stop()
	}
	state.generation++
	generation := state.generation
	state.timer = time.AfterFunc(c.returnRoutabilityTimeout(), func() {
		c.returnRoutabilityTimerExpired(generation)
	})
}

func (c *Conn) returnRoutabilityTimerExpired(generation uint64) {
	c.writeMu.Lock()
	state := c.returnRoutability
	if state == nil || state.generation != generation || state.phase == returnRoutabilityIdle {
		c.writeMu.Unlock()
		return
	}
	var err error
	if state.phase == returnRoutabilityCheckingOldPath {
		err = c.startBasicReturnRoutabilityLocked()
	} else {
		c.clearReturnRoutabilityLocked()
	}
	c.writeMu.Unlock()
	if err != nil {
		c.failConnection(err)
	}
}

func (c *Conn) clearReturnRoutability() {
	c.writeMu.Lock()
	c.clearReturnRoutabilityLocked()
	c.writeMu.Unlock()
}

func (c *Conn) clearReturnRoutabilityLocked() {
	state := c.returnRoutability
	if state == nil {
		return
	}
	if state.timer != nil {
		state.timer.Stop()
	}
	state.generation++
	state.phase = returnRoutabilityIdle
	state.timer = nil
	state.oldAddress = nil
	state.newAddress = nil
	clear(state.cookie[:])
	state.probeConnectionID = nil
	state.received = 0
	state.sent = 0
}
