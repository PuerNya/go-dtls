package dtls13

import (
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

// Signal ticket delivery without adding polling delays to every connection.
type packetLifecycleCache struct {
	*lruSessionCache
	ready chan struct{}
}

func (c *packetLifecycleCache) Put(key string, state *ClientSessionState) {
	c.lruSessionCache.Put(key, state)
	if state != nil {
		select {
		case c.ready <- struct{}{}:
		default:
		}
	}
}

type packetProcessLifecycle struct {
	configs []*Config
	locals  []*net.UDPAddr
}

func newPacketProcessLifecycle(config *Config, lanes int) packetProcessLifecycle {
	l := packetProcessLifecycle{configs: make([]*Config, lanes), locals: make([]*net.UDPAddr, lanes)}
	for i := range l.configs {
		l.configs[i] = config.Clone()
		l.configs[i].ClientSessionCache = &packetLifecycleCache{
			lruSessionCache: NewLRUClientSessionCache(1).(*lruSessionCache), ready: make(chan struct{}, 1),
		}
	}
	return l
}

func (l *packetProcessLifecycle) run(t *testing.T, spec packetProcessSpec, phase string, listener *Listener, wire *packetProcessWire) packetProcessStats {
	t.Helper()
	operations := spec.Operations
	if phase == "prime" {
		operations = spec.Connections
	}
	resumed := spec.Resume && phase != "prime"
	remote, err := net.ResolveUDPAddr("udp4", spec.Address)
	if listener == nil && err != nil {
		t.Fatal(err)
	}
	completed := make([]int, min(operations, spec.Connections))
	var closeBarrier *sync.WaitGroup
	if phase == "prime" && listener == nil {
		// Reserve every client port until all lanes and relay sockets exist.
		closeBarrier = new(sync.WaitGroup)
		closeBarrier.Add(len(completed))
	}
	var workers sync.WaitGroup
	for lane := range completed {
		workers.Go(func() {
			for i := lane; i < operations; i += len(completed) {
				if err := l.connection(spec, lane, resumed, remote, listener, wire, closeBarrier); err != nil {
					t.Errorf("%s lifecycle %d: %v", spec.Role, i, err)
					return
				}
				completed[lane]++
			}
		})
	}
	workers.Wait()
	if t.Failed() {
		t.FailNow()
	}
	if listener != nil {
		listener.mu.Lock()
		remaining := len(listener.sessions) + len(listener.pending) + len(listener.cidSessions)
		listener.mu.Unlock()
		if remaining != 0 {
			t.Fatalf("lifecycle left %d Listener routes", remaining)
		}
	}
	var result packetProcessStats
	for _, count := range completed {
		result.Handshakes += uint64(count)
	}
	if resumed {
		result.Resumed = result.Handshakes
	}
	return result
}

func (l *packetProcessLifecycle) connection(spec packetProcessSpec, lane int, resumed bool, remote *net.UDPAddr, listener *Listener, wire *packetProcessWire, closeBarrier *sync.WaitGroup) (err error) {
	var conn *Conn
	var cache *packetLifecycleCache
	defer func() {
		if closeBarrier != nil {
			closeBarrier.Done()
			closeBarrier.Wait()
		}
		if conn != nil {
			// A validated successor can already have retired the old Listener route.
			if closeErr := conn.Close(); err == nil && !errors.Is(closeErr, net.ErrClosed) {
				err = closeErr
			}
		}
	}()
	if listener != nil {
		conn, err = listener.Accept()
		if err != nil {
			return err
		}
	} else {
		config := l.configs[lane]
		cache = config.ClientSessionCache.(*packetLifecycleCache)
		if !resumed {
			cache.Put("server.test", nil)
		}
		// Rebind each lane so the relay's socket and route count stays bounded.
		socket, err := net.DialUDP("udp4", l.locals[lane], remote)
		if err != nil {
			return err
		}
		l.locals[lane] = socket.LocalAddr().(*net.UDPAddr)
		var transport net.Conn = socket
		if spec.Observe == "wire" {
			transport = &countedPacketProcessConn{Conn: socket, wire: wire}
		}
		conn = Client(transport, config)
	}
	if err := conn.Handshake(); err != nil {
		return err
	}
	if state := conn.ConnectionState(); !state.HandshakeComplete || state.DidResume != resumed {
		return errors.New("unexpected lifecycle handshake state")
	}
	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return err
	}
	var marker [1]byte
	if cache != nil {
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()
		select {
		case <-cache.ready:
		case <-timer.C:
			return errors.New("lifecycle ticket was not delivered")
		}
		if _, err := conn.WriteDatagram(marker[:]); err != nil {
			return err
		}
		if _, _, err := conn.ReadDatagram(marker[:]); !errors.Is(err, io.EOF) {
			return errors.Join(errors.New("lifecycle peer did not close after ticket delivery"), err)
		}
	} else {
		if n, info, err := conn.ReadDatagram(marker[:]); err != nil {
			return err
		} else if n != 1 || marker[0] != 0 || info.Truncated || info.EarlyData {
			return errors.New("invalid lifecycle completion marker")
		}
	}
	return nil
}

func TestPacketProcessLifecycle(t *testing.T) {
	const connections, operations = 64, 131
	for _, resumed := range []bool{false, true} {
		name := "Full"
		if resumed {
			name = "Resumed"
		}
		t.Run(name, func(t *testing.T) {
			run := newPacketProcessRun(t, packetProcessSpec{Lifecycle: true, Resume: resumed, Connections: connections, Operations: operations})
			for range 2 {
				stats := run.stage(t, "lifecycle")
				for _, role := range []string{"client", "server"} {
					wantResumed := uint64(0)
					if resumed {
						wantResumed = operations
					}
					if got := stats[role]; got.Handshakes != operations || got.Resumed != wantResumed {
						t.Fatalf("%s completed %d handshakes, %d resumed", role, got.Handshakes, got.Resumed)
					}
				}
			}
			run.close(t)
		})
	}
}
