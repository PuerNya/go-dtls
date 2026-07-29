package dtls13

import (
	"net"
	"sync"
)

// amplificationGuard enforces RFC 9147 section 5.1's recommendation that a
// server send at most three times the bytes received before address validation.
type amplificationGuard struct {
	mu             sync.Mutex
	received, sent uint64
	validated      bool
}

func (g *amplificationGuard) recordReceived(n int) {
	if n <= 0 {
		return
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if ^uint64(0)-g.received < uint64(n) {
		g.received = ^uint64(0)
	} else {
		g.received += uint64(n)
	}
}
func (g *amplificationGuard) allowSend(n int) bool {
	if n < 0 {
		return false
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.validated {
		return true
	}
	budget := g.received
	if budget > ^uint64(0)/3 {
		budget = ^uint64(0)
	} else {
		budget *= 3
	}
	if uint64(n) > budget-g.sent {
		return false
	}
	g.sent += uint64(n)
	return true
}
func (g *amplificationGuard) validate() { g.mu.Lock(); g.validated = true; g.mu.Unlock() }

// amplificationConn accounts for every datagram before peer address
// validation. Writes beyond the current budget are intentionally suppressed;
// later received datagrams increase the budget and allow retransmission.
type amplificationConn struct {
	net.Conn
	guard *amplificationGuard
}

func (c *amplificationConn) Read(p []byte) (int, error) {
	n, err := c.Conn.Read(p)
	c.guard.recordReceived(n)
	return n, err
}

func (c *amplificationConn) Write(p []byte) (int, error) {
	if !c.guard.allowSend(len(p)) {
		return len(p), nil
	}
	return c.Conn.Write(p)
}
