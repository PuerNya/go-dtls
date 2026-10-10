package dtls13

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"runtime"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func newImpairedPairs(t *testing.T, count int, seed uint64, profile packetNetworkProfile, clientConfig, serverConfig *Config, early bool) ([]loadPair, []*impairedPacketConn) {
	t.Helper()
	serverSocket := newImpairedPacketConn(t, seed, profile)
	listener := NewListener(serverSocket, serverConfig)
	t.Cleanup(func() { _ = listener.Close() })
	// Bound Accept as well as Handshake when every ClientHello is lost.
	timer := time.AfterFunc(2*serverConfig.HandshakeTimeout, func() { _ = listener.Close() })
	defer timer.Stop()
	pairs := make([]loadPair, count)
	links := []*impairedPacketConn{serverSocket}
	byAddress := make(map[string]int, count)
	for i := range pairs {
		socket := newImpairedPacketConn(t, seed+uint64(i)+1, profile)
		links = append(links, socket)
		byAddress[socket.LocalAddr().String()] = i
		transport := &connectedPacketConn{PacketConn: socket, remote: serverSocket.LocalAddr()}
		if early {
			pairs[i].client = ClientEarly(transport, clientConfig)
		} else {
			pairs[i].client = Client(transport, clientConfig)
		}
	}
	t.Cleanup(func() {
		for _, pair := range pairs {
			_ = pair.client.Close()
			if pair.server != nil {
				_ = pair.server.Close()
			}
		}
	})
	type acceptedConn struct {
		conn *Conn
		err  error
	}
	servers := make(chan acceptedConn, count+1)
	go func() {
		defer close(servers)
		var handshakes sync.WaitGroup
		defer handshakes.Wait()
		for range count {
			server, err := listener.Accept()
			if err != nil {
				servers <- acceptedConn{err: err}
				return
			}
			handshakes.Go(func() { servers <- acceptedConn{server, server.Handshake()} })
		}
	}()
	var clients sync.WaitGroup
	for _, pair := range pairs {
		clients.Go(func() {
			var err error
			if early {
				_, err = pair.client.WriteDatagram([]byte("early"))
			}
			if err == nil {
				err = pair.client.Handshake()
			}
			if err != nil {
				t.Errorf("client %s handshake: %v", pair.client.LocalAddr(), err)
				_ = listener.Close()
			}
		})
	}
	clients.Wait()
	for result := range servers {
		if result.err != nil {
			t.Error("server handshake:", result.err)
		}
		server := result.conn
		if server == nil {
			continue
		}
		i, ok := byAddress[server.RemoteAddr().String()]
		if !ok || pairs[i].server != nil {
			_ = server.Close()
			t.Error("unexpected accepted connection", server.RemoteAddr())
			continue
		}
		pairs[i].server = server
	}
	if t.Failed() {
		t.FailNow()
	}
	for _, pair := range pairs {
		if pair.server == nil {
			t.Fatal("missing accepted connection")
		}
		if early && (!pair.client.ConnectionState().DidResume || pair.client.ConnectionState().EarlyData != EarlyDataAccepted) {
			t.Fatal("resumption did not accept early data")
		}
	}
	return pairs, links
}

func TestPacketNetworkScenarios(t *testing.T) {
	beforeGoroutines := runtime.NumGoroutine()
	t.Cleanup(func() {
		deadline := time.Now().Add(3 * time.Second)
		for runtime.NumGoroutine() > beforeGoroutines+8 && time.Now().Before(deadline) {
			time.Sleep(10 * time.Millisecond)
		}
		if got := runtime.NumGoroutine(); got > beforeGoroutines+8 {
			t.Errorf("goroutines did not recover: %d -> %d", beforeGoroutines, got)
		}
	})
	count, messages, seeds := 4, 128, 2
	stress := os.Getenv("GO_DTLS_STRESS") == "1"
	if stress {
		count, messages, seeds = 64, 2048, 5
	}
	seedBase := uint64(1)
	if value := os.Getenv("GO_DTLS_NETWORK_SEED"); value != "" {
		var err error
		seedBase, err = strconv.ParseUint(value, 10, 64)
		if err != nil {
			t.Fatal("GO_DTLS_NETWORK_SEED must be a uint64")
		}
	} else if stress {
		var entropy [8]byte
		if _, err := rand.Read(entropy[:]); err != nil {
			t.Fatal(err)
		}
		seedBase = binary.LittleEndian.Uint64(entropy[:])
	}
	t.Logf("network_seed_base=%d", seedBase)
	variable := packetNetworkProfile{delay: time.Millisecond, jitter: 4 * time.Millisecond, minRate: 2 << 20, maxRate: 20 << 20,
		loss: 30, duplicate: 20, queuePackets: 256, queueBytes: 256 << 10}
	extreme := packetNetworkProfile{delay: time.Millisecond, jitter: 30 * time.Millisecond, minRate: 128 << 10, maxRate: 64 << 20,
		rateChoices: []int{128 << 10, 512 << 10, 2 << 20, 8 << 20, 64 << 20}, loss: 0, maxLoss: 950,
		duplicate: 50, queuePackets: 256, queueBytes: 256 << 10}
	lossStorm := extreme
	lossStorm.minRate, lossStorm.maxRate, lossStorm.rateChoices = 4<<20, 20<<20, nil
	rateCollapse := extreme
	rateCollapse.loss, rateCollapse.maxLoss = 30, 0
	networks := []struct {
		name    string
		profile packetNetworkProfile
		traffic *packetNetworkProfile
	}{
		{"Variable", variable, nil},
		{"Constrained", packetNetworkProfile{delay: 3 * time.Millisecond, jitter: 12 * time.Millisecond, minRate: 256 << 10, maxRate: 4 << 20, loss: 100, duplicate: 40, queuePackets: 128, queueBytes: 128 << 10}, nil},
		{"ExtremeBursts", variable, &extreme},
		{"LossStorm", variable, &lossStorm},
		{"RateCollapse", variable, &rateCollapse},
	}
	if !stress {
		networks = networks[:3]
	}
	for networkIndex, network := range networks {
		for cidIndex, cid := range []bool{false, true} {
			for seed := 1; seed <= seeds; seed++ {
				t.Run(fmt.Sprintf("%s/CID%t/Seed%d", network.name, cid, seed), func(t *testing.T) {
					scenarioCount, pause := count, 2*time.Millisecond
					if network.traffic != nil {
						scenarioCount, pause = min(count, 8), 10*time.Millisecond
					}
					scenarioSeed := seedBase + uint64(networkIndex*10000+cidIndex*1000+seed*100)
					t.Logf("scenario_seed=%d connections=%d pause_per_16=%s", scenarioSeed, scenarioCount, pause)
					certificate, roots := testServerCertificate(t)
					client := &Config{RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true}
					server := &Config{Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true, MaxPendingConnections: scenarioCount}
					for _, config := range []*Config{client, server} {
						config.HandshakeTimeout = 15 * time.Second
						config.FlightInterval = 40 * time.Millisecond
						config.MaxFlightInterval = 500 * time.Millisecond
						config.MTU = 1280
						config.MaxBufferedApplicationDatagrams = messages + 64
						config.MaxBufferedApplicationData = (messages + 64) * 1200
					}
					if cid {
						client.ConnectionID = []byte("client")
						var next uint64
						server.GetConnectionID = func() ([]byte, error) { next++; return binary.BigEndian.AppendUint64(nil, next), nil }
					}
					runtime.GC()
					var before, after runtime.MemStats
					runtime.ReadMemStats(&before)
					started := time.Now()
					pairs, links := newImpairedPairs(t, scenarioCount, scenarioSeed, network.profile, client, server, false)
					t.Logf("handshakes=%d elapsed=%s", scenarioCount, time.Since(started))
					// Corrupt authenticated application records, not plaintext handshake messages.
					for _, link := range links {
						link.mu.Lock()
						if network.traffic != nil {
							link.profile = *network.traffic
							link.lossRemaining = 0
							link.next = time.Now()
						}
						link.corrupt = 10
						link.mu.Unlock()
					}
					var traffic sync.WaitGroup
					var received atomic.Uint64
					for i, pair := range pairs {
						traffic.Go(func() {
							received.Add(uint64(exerciseImpairedTraffic(t, pair.client, pair.server, uint64(2*i), messages, false, pause)))
						})
						traffic.Go(func() {
							received.Add(uint64(exerciseImpairedTraffic(t, pair.server, pair.client, uint64(2*i+1), messages, i%2 == 0, pause)))
						})
					}
					traffic.Wait()
					if received.Load() == 0 {
						t.Error("no application traffic crossed the network")
					}
					if network.traffic != nil {
						for _, link := range links {
							link.mu.Lock()
							link.profile = network.profile
							link.lossRemaining = 0
							link.rate = network.profile.maxRate
							link.next = time.Now()
							link.mu.Unlock()
						}
					}
					// Collection deadlines stop record readers; resume them to allow
					// control retransmissions to recover after the application burst.
					for _, pair := range pairs {
						for _, conn := range []*Conn{pair.client, pair.server} {
							if err := conn.SetReadDeadline(time.Time{}); err != nil {
								t.Fatal(err)
							}
							conn.startRecordReader()
						}
					}
					recoveryDeadline := time.Now().Add(3 * time.Second)
					for _, pair := range pairs {
						for _, conn := range []*Conn{pair.client, pair.server} {
							for {
								conn.writeMu.Lock()
								updated := conn.sendCipher != nil && conn.sendCipher.epoch >= 4
								conn.writeMu.Unlock()
								if updated {
									break
								}
								if time.Now().After(recoveryDeadline) {
									t.Error("KeyUpdate was not acknowledged")
									break
								}
								time.Sleep(time.Millisecond)
							}
						}
					}
					for _, pair := range pairs {
						_ = pair.client.Close()
						_ = pair.server.Close()
					}
					for _, link := range links {
						_ = link.Close()
						stats := link.snapshot()
						if stats.peakPackets > network.profile.queuePackets || stats.peakBytes > network.profile.queueBytes || stats.writes+stats.duplicated != stats.lost+stats.overflow+stats.sent+stats.failed+stats.canceled {
							t.Errorf("unbalanced link: %+v", stats)
						}
						if link == links[0] {
							t.Logf("server_link=%+v", stats)
						}
					}
					runtime.GC()
					runtime.ReadMemStats(&after)
					t.Logf("offered_datagrams=%d received=%d elapsed=%s allocated=%d heap_delta=%d gc_cycles=%d gc_pause=%s", 2*scenarioCount*messages, received.Load(), time.Since(started), after.TotalAlloc-before.TotalAlloc, int64(after.HeapAlloc)-int64(before.HeapAlloc), after.NumGC-before.NumGC, time.Duration(after.PauseTotalNs-before.PauseTotalNs))
				})
			}
		}
	}
}

// Application loss is expected: send once and validate every received sequence,
// payload, source and boundary. A deadline ends collection, not a reliable retry.
func exerciseImpairedTraffic(t *testing.T, sender, receiver *Conn, stream uint64, messages int, slow bool, pause time.Duration) int {
	t.Helper()
	size := min(1200, sender.PathMTU()-sender.RecordOverhead())
	if size < 24 {
		t.Error("insufficient application MTU")
		return 0
	}
	deadline := time.Now().Add(2*time.Second + time.Duration(messages/16)*3*time.Millisecond)
	if err := receiver.SetReadDeadline(deadline); err != nil {
		t.Error(err)
		return 0
	}
	written := make(chan error, 1)
	go func() {
		payload := make([]byte, size)
		for sequence := range messages {
			if sequence == messages/2 {
				if err := sender.SendKeyUpdate(false); err != nil {
					written <- err
					return
				}
			}
			for i := range payload {
				payload[i] = byte(sequence)
			}
			binary.BigEndian.PutUint64(payload[:8], stream)
			binary.BigEndian.PutUint64(payload[8:16], uint64(sequence))
			binary.BigEndian.PutUint64(payload[16:24], uint64(time.Now().UnixNano()))
			if _, err := sender.WriteDatagram(payload); err != nil {
				written <- err
				return
			}
			if sequence%16 == 15 {
				time.Sleep(pause)
			}
		}
		written <- nil
	}()
	seen := make([]bool, messages)
	received, reordered, highest := 0, 0, -1
	latencies := make([]time.Duration, 0, messages)
	buffer := make([]byte, size)
	for {
		n, info, err := receiver.ReadDatagram(buffer[:])
		if err != nil {
			if networkErr, ok := errors.AsType[net.Error](err); !ok || !networkErr.Timeout() {
				t.Errorf("stream %d: %v", stream, err)
			}
			break
		}
		sequence := binary.BigEndian.Uint64(buffer[8:16])
		if n != len(buffer) || info.Truncated || info.FullLength != len(buffer) || info.EarlyData || binary.BigEndian.Uint64(buffer[:8]) != stream || sequence >= uint64(messages) || info.Source == nil || info.Source.String() != sender.LocalAddr().String() {
			t.Errorf("stream %d invalid datagram: n=%d info=%+v sequence=%d", stream, n, info, sequence)
			break
		}
		if seen[sequence] {
			t.Errorf("stream %d duplicated/corrupted sequence %d", stream, sequence)
			break
		}
		for _, value := range buffer[24:] {
			if value != byte(sequence) {
				t.Errorf("stream %d corrupted sequence %d", stream, sequence)
				break
			}
		}
		seen[sequence] = true
		latencies = append(latencies, time.Since(time.Unix(0, int64(binary.BigEndian.Uint64(buffer[16:24])))))
		received++
		if int(sequence) < highest {
			reordered++
		}
		highest = max(highest, int(sequence))
		if slow {
			time.Sleep(time.Millisecond)
		}
	}
	if err := <-written; err != nil {
		t.Errorf("stream %d write: %v", stream, err)
	}
	if len(latencies) > 0 {
		slices.Sort(latencies)
		t.Logf("stream=%d bytes=%d slow=%t received=%d/%d reordered=%d p50=%s p95=%s p99=%s", stream, size, slow, received, messages, reordered, latencies[len(latencies)/2], latencies[len(latencies)*95/100], latencies[len(latencies)*99/100])
	}
	return received
}

func TestPacketNetworkResumption(t *testing.T) {
	for _, early := range []bool{false, true} {
		t.Run(fmt.Sprintf("Early%t", early), func(t *testing.T) {
			client, server := earlyIOConfigs(t)
			for _, config := range []*Config{client, server} {
				config.HandshakeTimeout = 10 * time.Second
				config.FlightInterval = 40 * time.Millisecond
				config.MaxFlightInterval = 500 * time.Millisecond
			}
			profile := packetNetworkProfile{delay: time.Millisecond, jitter: 4 * time.Millisecond, minRate: 256 << 10, maxRate: 2 << 20, loss: 100, duplicate: 30, queuePackets: 128, queueBytes: 128 << 10}
			for round := range 3 {
				t.Run(fmt.Sprintf("Connection%d", round), func(t *testing.T) {
					pairs, _ := newImpairedPairs(t, 1, uint64(90+round), profile, client, server, early && round > 0)
					pair := pairs[0]
					if pair.client.ConnectionState().DidResume != (round > 0) {
						t.Fatal("unexpected resumption state")
					}
					deadline := time.Now().Add(3 * time.Second)
					for {
						if _, ok := client.ClientSessionCache.Get("server.test"); ok {
							break
						}
						if time.Now().After(deadline) {
							t.Fatal("session ticket was not recovered under loss")
						}
						time.Sleep(time.Millisecond)
					}
					if early && round > 0 {
						if err := pair.server.SetReadDeadline(time.Now().Add(50 * time.Millisecond)); err != nil {
							t.Fatal(err)
						}
						var payload [16]byte
						n, info, err := pair.server.ReadDatagram(payload[:])
						if err == nil {
							if !info.EarlyData || !bytes.Equal(payload[:n], []byte("early")) {
								t.Fatalf("bad early datagram: %q %+v", payload[:n], info)
							}
						} else if networkErr, ok := errors.AsType[net.Error](err); !ok || !networkErr.Timeout() {
							t.Fatal(err)
						}
					}
				})
				if t.Failed() {
					break
				}
			}
		})
	}
}

func TestPacketNetworkSlowReaderBound(t *testing.T) {
	client, server := earlyIOConfigs(t)
	client.SessionTicketsDisabled, server.SessionTicketsDisabled = true, true
	server.MaxBufferedApplicationDatagrams, server.MaxBufferedApplicationData = 4, 64
	profile := packetNetworkProfile{minRate: 1 << 20, maxRate: 2 << 20, queuePackets: 64, queueBytes: 64 << 10, jitter: time.Millisecond}
	pairs, _ := newImpairedPairs(t, 1, 71, profile, client, server, false)
	pair := pairs[0]
	for range 32 {
		if _, err := pair.client.WriteDatagram([]byte("unread")); err != nil {
			t.Fatal(err)
		}
	}
	deadline := time.Now().Add(2 * time.Second)
	for {
		pair.server.inputMu.Lock()
		err, queued, queuedBytes := pair.server.readErr, len(pair.server.applicationDatagrams), pair.server.bufferedApplicationBytes
		pair.server.inputMu.Unlock()
		if queued > 4 || queuedBytes > 64 {
			t.Fatalf("queue exceeded bound: %d/%d", queued, queuedBytes)
		}
		if err != nil {
			if _, ok := errors.AsType[*ProtocolError](err); !ok {
				t.Fatalf("unexpected overflow error: %v", err)
			}
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("slow reader did not reach its configured bound")
		}
		time.Sleep(time.Millisecond)
	}
}
