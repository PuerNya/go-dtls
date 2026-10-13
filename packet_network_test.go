package dtls13

import (
	"bytes"
	"container/heap"
	"errors"
	"math/rand/v2"
	"net"
	"os"
	"sync"
	"testing"
	"time"
)

type packetNetworkProfile struct {
	delay, jitter            time.Duration
	minRate, maxRate         int // Bytes per second; a new rate is drawn every 16 writes.
	rateChoices              []int
	loss, maxLoss, duplicate int // Probabilities in thousandths; maxLoss enables variable-length loss bursts.
	queuePackets, queueBytes int
}

type packetNetworkStats struct {
	writes, lost, overflow, corrupted, duplicated, sent, failed, canceled uint64
	sentBytes, peakBytes, peakPackets, minRate, maxRate, minLoss, maxLoss int
}

type scheduledPacket struct {
	at    time.Time
	data  []byte
	to    net.Addr
	write uint64
	copy  int
}

type scheduledPackets []scheduledPacket

func (q *scheduledPackets) Len() int           { return len(*q) }
func (q *scheduledPackets) Less(i, j int) bool { return (*q)[i].at.Before((*q)[j].at) }
func (q *scheduledPackets) Swap(i, j int)      { (*q)[i], (*q)[j] = (*q)[j], (*q)[i] }
func (q *scheduledPackets) Push(x any)         { *q = append(*q, x.(scheduledPacket)) }
func (q *scheduledPackets) Pop() any {
	i := len(*q) - 1
	p := (*q)[i]
	(*q)[i] = scheduledPacket{}
	*q = (*q)[:i]
	return p
}

// Each socket has one bounded outgoing link and one scheduler, not one goroutine
// per packet. Seeds reproduce draws; concurrent write ordering remains scheduler-dependent.
type impairedPacketConn struct {
	net.PacketConn
	mu                                                  sync.Mutex
	profile                                             packetNetworkProfile
	random                                              *rand.Rand
	queue                                               scheduledPackets
	queuedBytes, rate, corrupt, lossRate, lossRemaining int
	next, writeDeadline                                 time.Time
	stats                                               packetNetworkStats
	writeErr                                            error
	wake, stop, done                                    chan struct{}
	closeOnce                                           sync.Once
	trace                                               *packetNetworkTrace
}

func newImpairedPacketConn(t testing.TB, seed uint64, profile packetNetworkProfile) *impairedPacketConn {
	t.Helper()
	c, err := openImpairedPacketConn(seed, profile)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func openImpairedPacketConn(seed uint64, profile packetNetworkProfile) (*impairedPacketConn, error) {
	if profile.minRate <= 0 || profile.maxRate < profile.minRate || profile.queuePackets < 1 || profile.queueBytes < 1 || profile.delay < 0 || profile.jitter < 0 || profile.loss < 0 || profile.loss > 1000 || profile.maxLoss > 1000 || profile.maxLoss != 0 && profile.maxLoss < profile.loss || profile.duplicate < 0 || profile.duplicate > 1000 {
		return nil, errors.New("invalid packet network profile")
	}
	for _, rate := range profile.rateChoices {
		if rate < profile.minRate || rate > profile.maxRate {
			return nil, errors.New("invalid packet network rate")
		}
	}
	raw, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	//nolint:gosec // Reproducible test impairments, never cryptographic randomness.
	c := &impairedPacketConn{PacketConn: raw, profile: profile, random: rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15)),
		wake: make(chan struct{}, 1), stop: make(chan struct{}), done: make(chan struct{})}
	go c.run()
	return c, nil
}

func (c *impairedPacketConn) WriteTo(p []byte, to net.Addr) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.stop:
		return 0, net.ErrClosed
	default:
	}
	if c.writeErr != nil {
		return 0, c.writeErr
	}
	now := time.Now()
	if !c.writeDeadline.IsZero() && !now.Before(c.writeDeadline) {
		return 0, os.ErrDeadlineExceeded
	}
	decision, err := c.nextDecision(len(p))
	if err != nil {
		c.writeErr = err
		return 0, err
	}
	c.stats.writes++
	if c.stats.minRate == 0 || decision.Rate < c.stats.minRate {
		c.stats.minRate = decision.Rate
	}
	c.stats.maxRate = max(c.stats.maxRate, decision.Rate)
	if c.stats.writes == 1 || decision.Loss < c.stats.minLoss {
		c.stats.minLoss = decision.Loss
	}
	c.stats.maxLoss = max(c.stats.maxLoss, decision.Loss)
	if decision.Lost {
		c.stats.lost++
		if err := c.recordDecision(decision); err != nil {
			c.writeErr = err
			return 0, err
		}
		return len(p), nil
	}
	copies := 1
	if decision.Duplicate {
		copies++
		c.stats.duplicated++
	}
	expected := decision.Queued
	decision.Queued = [2]bool{}
	for i := range copies {
		if len(c.queue) >= c.profile.queuePackets || len(p) > c.profile.queueBytes-c.queuedBytes {
			c.stats.overflow++
			continue
		}
		data := bytes.Clone(p)
		if len(data) != 0 && decision.Corrupt[i] != 0 {
			data[len(data)-1] ^= decision.Corrupt[i]
			c.stats.corrupted++
		}
		if c.next.Before(now) {
			c.next = now
		}
		c.next = c.next.Add(time.Duration(int64(len(data)) * int64(time.Second) / int64(decision.Rate)))
		at := c.next.Add(decision.Delay[i])
		address := *to.(*net.UDPAddr)
		address.IP = bytes.Clone(address.IP)
		heap.Push(&c.queue, scheduledPacket{at: at, data: data, to: &address, write: decision.Write, copy: i})
		decision.Queued[i] = true
		c.queuedBytes += len(data)
		c.stats.peakPackets = max(c.stats.peakPackets, len(c.queue))
		c.stats.peakBytes = max(c.stats.peakBytes, c.queuedBytes)
	}
	if c.trace != nil && c.trace.decode != nil && decision.Queued != expected {
		c.writeErr = errors.New("replay queue occupancy diverged")
		return 0, c.writeErr
	}
	if err := c.recordDecision(decision); err != nil {
		c.writeErr = err
		return 0, err
	}
	select {
	case c.wake <- struct{}{}:
	default:
	}
	return len(p), nil
}

func (c *impairedPacketConn) run() {
	defer close(c.done)
	timer := time.NewTimer(time.Hour)
	defer timer.Stop()
	for {
		c.mu.Lock()
		select {
		case <-c.stop:
			for _, p := range c.queue {
				if err := c.recordDecision(packetNetworkDecision{Kind: "canceled", Write: p.write, Copy: p.copy, Size: len(p.data)}); err != nil && c.writeErr == nil {
					c.writeErr = err
				}
			}
			c.stats.canceled += uint64(len(c.queue))
			c.queue = nil
			c.queuedBytes = 0
			c.mu.Unlock()
			return
		default:
		}
		delay := time.Hour
		if len(c.queue) != 0 {
			delay = time.Until(c.queue[0].at)
		}
		if delay > 0 {
			c.mu.Unlock()
			timer.Reset(delay)
			select {
			case <-timer.C:
			case <-c.wake:
			case <-c.stop:
			}
			continue
		}
		p := heap.Pop(&c.queue).(scheduledPacket)
		c.queuedBytes -= len(p.data)
		c.mu.Unlock()
		n, err := c.PacketConn.WriteTo(p.data, p.to)
		c.mu.Lock()
		outcome := "sent"
		if err != nil {
			outcome = "failed"
		}
		if traceErr := c.recordDecision(packetNetworkDecision{Kind: outcome, Write: p.write, Copy: p.copy, Size: len(p.data)}); traceErr != nil && c.writeErr == nil {
			c.writeErr = traceErr
		}
		if err != nil {
			c.stats.failed++
			c.writeErr = err
		} else {
			c.stats.sent++
			c.stats.sentBytes += n
		}
		c.mu.Unlock()
	}
}

func (c *impairedPacketConn) Close() error {
	c.closeOnce.Do(func() {
		c.mu.Lock()
		close(c.stop)
		c.mu.Unlock()
		_ = c.PacketConn.Close()
	})
	<-c.done
	return nil
}

func (c *impairedPacketConn) SetWriteDeadline(deadline time.Time) error {
	c.mu.Lock()
	c.writeDeadline = deadline
	c.mu.Unlock()
	return nil
}

func (c *impairedPacketConn) SetDeadline(deadline time.Time) error {
	_ = c.SetWriteDeadline(deadline)
	return c.SetReadDeadline(deadline)
}

func (c *impairedPacketConn) snapshot() packetNetworkStats {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.stats
}

func TestImpairedPacketConn(t *testing.T) {
	profile := packetNetworkProfile{minRate: 1 << 20, maxRate: 1 << 20, queuePackets: 2, queueBytes: 8, delay: time.Hour}
	receiver := newImpairedPacketConn(t, 1, profile)
	sender := newImpairedPacketConn(t, 2, profile)
	if err := receiver.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	payload := []byte("four")
	for range 3 {
		if _, err := sender.WriteTo(payload, receiver.LocalAddr()); err != nil {
			t.Fatal(err)
		}
	}
	clear(payload)
	// Release only after all enqueues, so preemption cannot change overflow.
	sender.mu.Lock()
	for i := range sender.queue {
		sender.queue[i].at = time.Now()
	}
	sender.profile.delay = 0
	sender.mu.Unlock()
	select {
	case sender.wake <- struct{}{}:
	default:
	}
	var buffer [4]byte
	for range 2 {
		n, from, err := receiver.ReadFrom(buffer[:])
		if err != nil || n != 4 || string(buffer[:]) != "four" || from.String() != sender.LocalAddr().String() {
			t.Fatalf("read=%q from=%v error=%v", buffer, from, err)
		}
	}
	stats := sender.snapshot()
	if stats.overflow != 1 || stats.peakBytes != 8 || stats.peakPackets != 2 {
		t.Fatalf("queue bounds: %+v", stats)
	}
	if err := sender.SetWriteDeadline(time.Now()); err != nil {
		t.Fatal(err)
	}
	if _, err := sender.WriteTo(payload, receiver.LocalAddr()); !errors.Is(err, os.ErrDeadlineExceeded) {
		t.Fatalf("write deadline: %v", err)
	}
	if err := sender.SetWriteDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	sender.mu.Lock()
	sender.profile.loss = 1000
	sender.mu.Unlock()
	if _, err := sender.WriteTo(payload, receiver.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	if sender.snapshot().lost != 1 {
		t.Fatal("loss not recorded")
	}
	sender.mu.Lock()
	sender.profile.loss, sender.profile.duplicate, sender.corrupt = 0, 1000, 1000
	sender.mu.Unlock()
	if _, err := sender.WriteTo([]byte("four"), receiver.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		n, _, err := receiver.ReadFrom(buffer[:])
		if err != nil || n != 4 || string(buffer[:3]) != "fou" || buffer[3] == 'r' {
			t.Fatalf("corrupted copy: %q error=%v", buffer, err)
		}
	}
	stats = sender.snapshot()
	if stats.duplicated != 1 || stats.corrupted != 2 {
		t.Fatalf("impairment counters: %+v", stats)
	}
	sender.mu.Lock()
	sender.profile.duplicate = 0
	sender.profile.delay = time.Hour
	sender.mu.Unlock()
	if _, err := sender.WriteTo(payload, receiver.LocalAddr()); err != nil {
		t.Fatal(err)
	}
	if err := sender.Close(); err != nil {
		t.Fatal(err)
	}
	if sender.snapshot().canceled != 1 {
		t.Fatal("close did not cancel queued packet")
	}
	if _, err := sender.WriteTo(payload, receiver.LocalAddr()); !errors.Is(err, net.ErrClosed) {
		t.Fatalf("closed write: %v", err)
	}
	if err := receiver.SetReadDeadline(time.Now()); err != nil {
		t.Fatal(err)
	}
	if _, _, err := receiver.ReadFrom(buffer[:]); !errors.Is(err, os.ErrDeadlineExceeded) {
		t.Fatalf("read deadline: %v", err)
	}
	if err := receiver.SetReadDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	readDone := make(chan error, 1)
	go func() { _, _, err := receiver.ReadFrom(buffer[:]); readDone <- err }()
	if err := receiver.Close(); err != nil {
		t.Fatal(err)
	}
	if err := <-readDone; !errors.Is(err, net.ErrClosed) {
		t.Fatalf("closed read: %v", err)
	}
}

func TestImpairedPacketConnVariableBursts(t *testing.T) {
	profile := packetNetworkProfile{minRate: 128 << 10, maxRate: 64 << 20,
		rateChoices: []int{128 << 10, 64 << 20}, maxLoss: 950,
		queuePackets: 1, queueBytes: 1, delay: time.Hour}
	sender := newImpairedPacketConn(t, 99, profile)
	to := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9}
	for range 2048 {
		if _, err := sender.WriteTo([]byte{1}, to); err != nil {
			t.Fatal(err)
		}
	}
	if err := sender.Close(); err != nil {
		t.Fatal(err)
	}
	stats := sender.snapshot()
	if stats.minRate != profile.minRate || stats.maxRate != profile.maxRate || stats.minLoss > 100 || stats.maxLoss < 850 ||
		stats.writes != stats.lost+stats.overflow+stats.sent+stats.failed+stats.canceled {
		t.Fatalf("variable bursts were not exercised: %+v", stats)
	}
}
