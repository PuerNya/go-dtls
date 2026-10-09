package dtls13

import (
	"context"
	"net"
)

// EarlyDataStatus reports the disposition of this connection's 0-RTT offer.
// Acceptance does not acknowledge delivery of individual datagrams.
type EarlyDataStatus uint8

const (
	// EarlyDataNotAttempted means no 0-RTT offer was made.
	EarlyDataNotAttempted EarlyDataStatus = iota
	// EarlyDataPending means the offer has not yet been accepted or rejected.
	EarlyDataPending
	// EarlyDataAccepted means the peer accepted the offer, or the server accepted it.
	EarlyDataAccepted
	// EarlyDataRejected means the offer was rejected; sent datagrams are not retried.
	EarlyDataRejected
)

// ClientEarly is like Client, but permits WriteDatagram to send 0-RTT data
// using an eligible cached session. It does not perform I/O or take a ticket
// from the cache until the handshake starts. Without an eligible ticket,
// writes wait for the handshake. Sent 0-RTT datagrams are replayable and are
// never automatically retransmitted, including when the server rejects them.
func ClientEarly(conn net.Conn, config *Config) *Conn {
	c := Client(conn, config)
	c.earlyIO = true
	return c
}

// DialEarly connects to a UDP address and starts a handshake. It returns when
// 0-RTT writing is possible, or after a successful handshake if no eligible
// early-data session exists. The handshake continues after an early return;
// ctx governs it until completion. Canceling ctx afterwards has no effect.
// The caller owns the returned connection even if its handshake later fails.
func DialEarly(ctx context.Context, network, address string, config *Config) (*Conn, error) {
	c, err := dialClient(ctx, &net.Dialer{}, network, address, config)
	if err != nil {
		return nil, err
	}
	c.earlyIO = true
	c.startHandshake(ctx)
	if err = c.waitReady(c.writeReady); err != nil {
		_ = c.Close()
		return nil, err
	}
	return c, nil
}

func (c *Conn) initLifecycle() {
	c.lifecycleOnce.Do(func() {
		c.lifetime, c.cancelLifetime = context.WithCancelCause(context.Background())
		c.handshakeComplete = make(chan struct{})
		c.handshakeDone = make(chan struct{})
		c.writeReady = make(chan struct{})
		c.applicationReady = make(chan struct{})
	})
}

// Context is canceled when the connection is closed or fails. Its cause is
// available through context.Cause. Calling Context does not start a handshake.
// A peer close_notify closes only the receive side and does not cancel Context.
func (c *Conn) Context() context.Context {
	c.initLifecycle()
	return c.lifetime
}

// HandshakeComplete is closed after a successful handshake. It remains open
// on failure: also select on Context().Done() to observe failure or closure.
// Calling HandshakeComplete does not start a handshake.
func (c *Conn) HandshakeComplete() <-chan struct{} {
	c.initLifecycle()
	return c.handshakeComplete
}

func (c *Conn) startHandshake(ctx context.Context) {
	if ctx == nil {
		panic("dtls13: nil handshake context")
	}
	c.initLifecycle()
	c.initInput()
	c.handshakeOnce.Do(func() {
		c.readerMu.Lock()
		if c.readerClosed {
			c.readerMu.Unlock()
			c.handshakeErr = net.ErrClosed
			c.cancelLifetime(c.handshakeErr)
			close(c.handshakeDone)
			return
		}
		c.handshaking = true
		c.readerMu.Unlock()
		go c.completeHandshake(ctx)
	})
}

func (c *Conn) completeHandshake(ctx context.Context) {
	var err error
	if ctx.Err() != nil {
		err = ctx.Err()
	} else if c.conn == nil {
		err = &ConfigError{"nil underlying connection"}
	} else {
		var config *Config
		if config, err = c.config.normalized(); err == nil {
			c.readerMu.Lock()
			c.config = config
			c.readerMu.Unlock()
			err = c.runHandshake(ctx)
		}
	}
	c.readerMu.Lock()
	if c.readerClosed {
		err = net.ErrClosed
	} else if cause := context.Cause(c.lifetime); cause != nil {
		err = cause
	}
	if err != nil {
		c.clearTrafficSecrets(err)
		c.mu.Lock()
		c.state.HandshakeComplete = false
		c.mu.Unlock()
		c.cancelLifetime(err)
	} else {
		c.mu.Lock()
		c.state.HandshakeComplete = true
		c.mu.Unlock()
		close(c.handshakeComplete)
		c.signalApplicationReady()
		c.signalWriteReady()
	}
	c.handshakeErr = err
	c.handshaking = false
	close(c.handshakeDone)
	c.readerMu.Unlock()
	c.notifyRead()
	if err == nil {
		c.startRecordReader()
	}
}

// Only the handshake goroutine publishes readiness.
func (c *Conn) signalWriteReady() {
	select {
	case <-c.writeReady:
	default:
		close(c.writeReady)
	}
}

func (c *Conn) signalApplicationReady() {
	select {
	case <-c.applicationReady:
	default:
		close(c.applicationReady)
	}
}

func (c *Conn) waitReady(ready <-chan struct{}) error {
	select {
	case <-ready:
	case <-c.lifetime.Done():
	}
	return context.Cause(c.lifetime)
}

func (c *Conn) stopEarlyWrites() {
	c.writeMu.Lock()
	c.earlyWriteCipher = nil
	c.earlyWriteRemaining = 0
	c.writeMu.Unlock()
}

func (c *Conn) writeEarlyDatagramLocked(p []byte) (int, error) {
	if err := context.Cause(c.lifetime); err != nil {
		return 0, err
	}
	for {
		if len(p) > c.maxApplicationDatagramForCipher(c.earlyWriteCipher) {
			return 0, datagramTooLargeError(c.RemoteAddr())
		}
		wire, err := c.earlyWriteCipher.seal(recordTypeApplicationData, p)
		if err != nil {
			return 0, err
		}
		if err = c.writeRecord(wire); err != nil {
			if !c.config.IgnorePathMTU && isMessageTooLong(err) {
				if _, reduced := c.reducePathMTU(); reduced {
					continue
				}
			}
			return 0, err
		}
		c.earlyWriteRemaining -= uint32(len(p))
		return len(p), nil
	}
}
