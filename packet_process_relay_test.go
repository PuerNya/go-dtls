package dtls13

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"testing"
	"time"
)

type packetRelayStats struct {
	Direction                                                               string
	Peer                                                                    int
	Received, Sent, Lost, Overflow, Duplicated, Corrupted, Failed, Canceled uint64
	SentBytes, PeakBytes, PeakPackets, MinRate, MaxRate, MinLoss, MaxLoss   int
}

type packetRelayPeer struct {
	conn           net.PacketConn
	client         net.Addr
	link           *impairedPacketConn
	received, sent uint64
}

type packetProcessRelay struct {
	front          net.PacketConn
	down           *impairedPacketConn
	server         *net.UDPAddr
	spec           packetProcessSpec
	mu             sync.Mutex
	peers          []*packetRelayPeer
	byAddress      map[string]*packetRelayPeer
	received, sent uint64
	err            error
	wg             sync.WaitGroup
	closing        bool
	traceFiles     []*os.File
	traceBuffers   []*bufio.Writer
}

func packetProcessNetwork(name string) (packetNetworkProfile, *packetNetworkProfile) {
	variable := packetNetworkProfile{delay: time.Millisecond, jitter: 4 * time.Millisecond, minRate: 2 << 20, maxRate: 20 << 20, loss: 30, duplicate: 20, queuePackets: 256, queueBytes: 256 << 10}
	extreme := packetNetworkProfile{delay: time.Millisecond, jitter: 30 * time.Millisecond, minRate: 128 << 10, maxRate: 64 << 20, rateChoices: []int{128 << 10, 512 << 10, 2 << 20, 8 << 20, 64 << 20}, maxLoss: 950, duplicate: 50, queuePackets: 256, queueBytes: 256 << 10}
	switch name {
	case "Variable":
		return variable, nil
	case "Constrained":
		return packetNetworkProfile{delay: 3 * time.Millisecond, jitter: 12 * time.Millisecond, minRate: 256 << 10, maxRate: 4 << 20, loss: 100, duplicate: 40, queuePackets: 128, queueBytes: 128 << 10}, nil
	case "ExtremeBursts":
		return variable, &extreme
	case "LossStorm":
		extreme.minRate, extreme.maxRate, extreme.rateChoices = 4<<20, 20<<20, nil
		return variable, &extreme
	case "RateCollapse":
		extreme.loss, extreme.maxLoss = 30, 0
		return variable, &extreme
	}
	return packetNetworkProfile{}, nil
}

func newPacketProcessRelay(t *testing.T, spec packetProcessSpec) *packetProcessRelay {
	t.Helper()
	server, err := net.ResolveUDPAddr("udp4", spec.Address)
	if err != nil {
		t.Fatal(err)
	}
	r := &packetProcessRelay{server: server, spec: spec, byAddress: make(map[string]*packetRelayPeer)}
	if spec.Network != "" {
		profile, _ := packetProcessNetwork(spec.Network)
		r.down = newImpairedPacketConn(t, spec.Seed, profile)
		if err := r.attachTrace(r.down, "down"); err != nil {
			t.Fatal(err)
		}
		r.front = r.down
	} else {
		r.front, err = net.ListenPacket("udp4", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
	}
	t.Cleanup(func() { r.close(t) })
	if err := r.raw(r.front).SetReadBuffer(4 << 20); err != nil {
		t.Fatal(err)
	}
	r.wg.Go(func() {
		buffer := make([]byte, 65535)
		for {
			n, from, err := r.front.ReadFrom(buffer)
			if err != nil {
				r.fail(err)
				return
			}
			r.mu.Lock()
			if r.closing {
				r.mu.Unlock()
				return
			}
			peer := r.byAddress[from.String()]
			if peer == nil && len(r.peers) < spec.Connections {
				peer = &packetRelayPeer{client: from}
				if spec.Network != "" {
					profile, _ := packetProcessNetwork(spec.Network)
					peer.link, err = openImpairedPacketConn(spec.Seed+uint64(len(r.peers))+1, profile)
					if err == nil {
						peer.conn = peer.link
						err = r.attachTrace(peer.link, fmt.Sprintf("up-%d", len(r.peers)))
					}
				} else {
					peer.conn, err = net.ListenPacket("udp4", "127.0.0.1:0")
				}
				if err != nil {
					if peer.conn != nil {
						_ = peer.conn.Close()
					}
					r.mu.Unlock()
					r.fail(err)
					return
				}
				if err = r.raw(peer.conn).SetReadBuffer(1 << 20); err != nil {
					_ = peer.conn.Close()
					r.mu.Unlock()
					r.fail(err)
					return
				}
				r.byAddress[from.String()] = peer
				r.peers = append(r.peers, peer)
				r.wg.Go(func() { r.readServer(peer) })
			}
			if peer != nil {
				peer.received++
			}
			r.mu.Unlock()
			if peer == nil {
				continue
			}
			if _, err := peer.conn.WriteTo(buffer[:n], server); err != nil {
				r.fail(err)
				return
			}
			r.mu.Lock()
			peer.sent++
			r.mu.Unlock()
		}
	})
	return r
}

func (r *packetProcessRelay) raw(conn net.PacketConn) *net.UDPConn {
	if link, ok := conn.(*impairedPacketConn); ok {
		return link.PacketConn.(*net.UDPConn)
	}
	return conn.(*net.UDPConn)
}

func (r *packetProcessRelay) readServer(peer *packetRelayPeer) {
	buffer := make([]byte, 65535)
	for {
		n, from, err := peer.conn.ReadFrom(buffer)
		if err != nil {
			r.fail(err)
			return
		}
		if from.String() != r.server.String() {
			r.fail(errors.New("relay received an unexpected server source"))
			return
		}
		r.mu.Lock()
		r.received++
		r.mu.Unlock()
		if _, err := r.front.WriteTo(buffer[:n], peer.client); err != nil {
			r.fail(err)
			return
		}
		r.mu.Lock()
		r.sent++
		r.mu.Unlock()
	}
}

func (r *packetProcessRelay) fail(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.closing && r.err == nil {
		r.err = err
	}
}

func (r *packetProcessRelay) phase(t *testing.T, phase string) {
	t.Helper()
	if r.spec.Network == "" {
		return
	}
	profile, traffic := packetProcessNetwork(r.spec.Network)
	if phase == "warmup" {
		profile = packetNetworkProfile{minRate: 64 << 20, maxRate: 64 << 20, queuePackets: 256, queueBytes: 256 << 10}
	}
	corrupt := 0
	if phase == "traffic" {
		if traffic != nil {
			profile = *traffic
		}
		corrupt = 10
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	links := []*impairedPacketConn{r.down}
	for _, peer := range r.peers {
		links = append(links, peer.link)
	}
	for _, link := range links {
		link.mu.Lock()
		link.profile, link.corrupt, link.lossRemaining, link.rate = profile, corrupt, 0, profile.maxRate
		if link.trace != nil {
			link.trace.phase = phase
		}
		link.next = time.Now()
		link.mu.Unlock()
	}
}

func (r *packetProcessRelay) snapshot(t *testing.T) []packetRelayStats {
	t.Helper()
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		t.Fatal(r.err)
	}
	stats := make([]packetRelayStats, 0, len(r.peers)+1)
	appendStats := func(direction string, index int, link *impairedPacketConn, received, sent uint64) {
		item := packetRelayStats{Direction: direction, Peer: index, Received: received, Sent: sent}
		if link != nil {
			s := link.snapshot()
			item.Received, item.Sent, item.Lost, item.Overflow, item.Duplicated, item.Corrupted, item.Failed, item.Canceled = s.writes, s.sent, s.lost, s.overflow, s.duplicated, s.corrupted, s.failed, s.canceled
			item.SentBytes, item.PeakBytes, item.PeakPackets, item.MinRate, item.MaxRate, item.MinLoss, item.MaxLoss = s.sentBytes, s.peakBytes, s.peakPackets, s.minRate, s.maxRate, s.minLoss, s.maxLoss
		}
		stats = append(stats, item)
	}
	appendStats("down", 0, r.down, r.received, r.sent)
	for i, peer := range r.peers {
		appendStats("up", i, peer.link, peer.received, peer.sent)
	}
	return stats
}

func (r *packetProcessRelay) close(t *testing.T) {
	t.Helper()
	r.mu.Lock()
	if r.closing {
		r.mu.Unlock()
		return
	}
	r.closing = true
	peers := slices.Clone(r.peers)
	r.mu.Unlock()
	_ = r.front.Close()
	for _, peer := range peers {
		_ = peer.conn.Close()
	}
	r.wg.Wait()
	for _, buffer := range r.traceBuffers {
		if err := buffer.Flush(); err != nil {
			t.Fatal(err)
		}
	}
	for _, file := range r.traceFiles {
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
	for _, stats := range r.snapshot(t) {
		if stats.Received+stats.Duplicated != stats.Sent+stats.Lost+stats.Overflow+stats.Failed+stats.Canceled {
			t.Errorf("unbalanced relay: %+v", stats)
		}
		if stats.Failed != 0 {
			t.Errorf("relay send failures: %+v", stats)
		}
	}
}

func (r *packetProcessRelay) attachTrace(link *impairedPacketConn, name string) error {
	if r.spec.Observe != "decisions" {
		return nil
	}
	//nolint:gosec // The local artifact root and internally assigned link names determine this path.
	file, err := os.Create(filepath.Join(r.spec.Output, "relay-"+name+".decisions.jsonl"))
	if err != nil {
		return err
	}
	buffer := bufio.NewWriterSize(file, 64<<10)
	r.traceFiles = append(r.traceFiles, file)
	r.traceBuffers = append(r.traceBuffers, buffer)
	link.mu.Lock()
	link.trace = &packetNetworkTrace{encode: json.NewEncoder(buffer), phase: "handshake"}
	link.mu.Unlock()
	return nil
}

func TestPacketProcessRelayTraceError(t *testing.T) {
	output := t.TempDir()
	path := filepath.Join(output, "relay-up-0.decisions.jsonl")
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	relay := newPacketProcessRelay(t, packetProcessSpec{Address: "127.0.0.1:9", Connections: 1, Network: "Variable", Observe: "decisions", Output: output})
	client, err := net.Dial("udp4", relay.front.LocalAddr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	if _, err := client.Write([]byte("trigger peer creation")); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for {
		relay.mu.Lock()
		err = relay.err
		if err != nil {
			relay.err = nil // Consume the expected setup error before cleanup.
		}
		relay.mu.Unlock()
		if err != nil {
			pathErr, ok := errors.AsType[*os.PathError](err)
			if !ok || pathErr.Path != path {
				t.Fatalf("relay setup error: %v", err)
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("relay did not report its trace setup error")
		}
		time.Sleep(time.Millisecond)
	}
	relay.close(t)
}
