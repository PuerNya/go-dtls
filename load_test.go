package dtls13

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"runtime"
	"slices"
	"strconv"
	"sync"
	"testing"
	"time"
)

type loadPair struct{ client, server *Conn }

func newMemoryLoadPairs(t testing.TB, count int, cid bool) []loadPair {
	t.Helper()
	certificate, roots := testServerCertificate(t)
	clientConfig := &Config{RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true, MTU: 1500}
	serverConfig := &Config{Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true, MTU: 1500}
	if cid {
		clientConfig.ConnectionID, serverConfig.ConnectionID = []byte("client"), []byte("server")
	}
	pairs := make([]loadPair, count)
	for i := range pairs {
		left, right := memoryDatagramPair()
		pairs[i] = loadPair{Client(left, clientConfig), Server(right, serverConfig)}
		handshakePair(t, pairs[i].client, pairs[i].server)
	}
	return pairs
}

func BenchmarkMemoryDatagramRoundTrip(b *testing.B) {
	for _, size := range []int{64, 1200} {
		b.Run(fmt.Sprintf("Bytes%d", size), func(b *testing.B) {
			pair := newMemoryLoadPairs(b, 1, false)[0]
			payload, buffer := make([]byte, size), make([]byte, size)
			b.ReportAllocs()
			b.SetBytes(int64(2 * size))
			for b.Loop() {
				if err := pair.roundTrip(payload, buffer); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// Use real UDP and the public API so routing, record readers, queues and
// transport costs are included alongside cryptography.
func newLoadPairs(t testing.TB, count int, cid bool) []loadPair {
	t.Helper()
	certificate, roots := testServerCertificate(t)
	config := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true,
		MaxPendingConnections: count, MTU: 1500, HandshakeTimeout: 10 * time.Second,
	}
	if cid {
		var next uint64
		config.GetConnectionID = func() ([]byte, error) {
			next++
			return binary.BigEndian.AppendUint64(nil, next), nil
		}
	}
	listener, err := Listen("udp4", "127.0.0.1:0", config)
	if err != nil {
		t.Fatal(err)
	}
	if os.Getenv("GO_DTLS_STRESS") == "1" {
		if err = listener.conn.(*net.UDPConn).SetReadBuffer(4 << 20); err != nil {
			_ = listener.Close()
			t.Fatal(err)
		}
	}
	pairs := make([]loadPair, 0, count)
	t.Cleanup(func() {
		for _, pair := range pairs {
			_ = pair.client.Close()
			_ = pair.server.Close()
		}
		_ = listener.Close()
	})
	type accepted struct {
		conn *Conn
		err  error
	}
	serverReady := make(chan accepted, 1)
	for i := range count {
		go func() {
			conn, acceptErr := listener.Accept()
			if acceptErr == nil {
				acceptErr = conn.Handshake()
			}
			serverReady <- accepted{conn, acceptErr}
		}()
		clientConfig := &Config{
			RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true,
			MTU: 1500, HandshakeTimeout: 10 * time.Second,
		}
		if cid {
			clientConfig.ConnectionID = binary.BigEndian.AppendUint64(nil, uint64(i+1))
		}
		client, dialErr := Dial("udp4", listener.Addr().String(), clientConfig)
		if dialErr != nil {
			_ = listener.Close()
		}
		server := <-serverReady
		if dialErr != nil || server.err != nil {
			if client != nil {
				_ = client.Close()
			}
			if server.conn != nil {
				_ = server.conn.Close()
			}
			t.Fatalf("pair %d: client=%v server=%v", i, dialErr, server.err)
		}
		pairs = append(pairs, loadPair{client, server.conn})
		deadline := time.Now().Add(5 * time.Minute)
		if err = client.SetDeadline(deadline); err != nil {
			t.Fatal(err)
		}
		if err = server.conn.SetDeadline(deadline); err != nil {
			t.Fatal(err)
		}
	}
	return pairs
}

func (p loadPair) roundTrip(payload, buffer []byte) error {
	if _, err := p.client.WriteDatagram(payload); err != nil {
		return err
	}
	n, info, err := p.server.ReadDatagram(buffer)
	if err != nil {
		return err
	}
	if info.Truncated || !bytes.Equal(buffer[:n], payload) {
		return fmt.Errorf("server datagram mismatch: length=%d truncated=%v", n, info.Truncated)
	}
	if _, err = p.server.WriteDatagram(buffer[:n]); err != nil {
		return err
	}
	n, info, err = p.client.ReadDatagram(buffer)
	if err != nil {
		return err
	}
	if info.Truncated || !bytes.Equal(buffer[:n], payload) {
		return fmt.Errorf("client datagram mismatch: length=%d truncated=%v", n, info.Truncated)
	}
	return nil
}

func BenchmarkListenerDatagramRoundTrip(b *testing.B) {
	for _, count := range []int{1, 64, 1024} {
		for _, cid := range []bool{false, true} {
			cidName := "NoCID"
			if cid {
				cidName = "CID"
			}
			for _, size := range []int{64, 1200} {
				b.Run(fmt.Sprintf("Connections%d/%s/Bytes%d", count, cidName, size), func(b *testing.B) {
					pairs := newLoadPairs(b, count, cid)
					payload, buffer := make([]byte, size), make([]byte, size)
					for _, pair := range pairs {
						if err := pair.roundTrip(payload, buffer); err != nil {
							b.Fatal(err)
						}
					}
					b.ReportAllocs()
					b.SetBytes(int64(2 * size))
					var sequence uint64
					for b.Loop() {
						binary.BigEndian.PutUint64(payload, sequence)
						if err := pairs[sequence%uint64(count)].roundTrip(payload, buffer); err != nil {
							b.Fatal(err)
						}
						sequence++
					}
				})
			}
		}
	}
}

func TestListenerDatagramLoad(t *testing.T) {
	testDatagramLoad(t, newLoadPairs)
}

func TestMemoryDatagramLoad(t *testing.T) {
	testDatagramLoad(t, newMemoryLoadPairs)
}

func testDatagramLoad(t *testing.T, makePairs func(testing.TB, int, bool) []loadPair) {
	t.Helper()
	if os.Getenv("GO_DTLS_STRESS") != "1" {
		t.Skip("set GO_DTLS_STRESS=1 for the large datagram workload")
	}
	const connections, exchanges, payloadSize = 1024, 1024, 1200
	concurrency := 16
	if value := os.Getenv("GO_DTLS_STRESS_WORKERS"); value != "" {
		var err error
		concurrency, err = strconv.Atoi(value)
		if err != nil || concurrency < 1 || concurrency > connections {
			t.Fatal("GO_DTLS_STRESS_WORKERS must be between 1 and 1024")
		}
	}
	runtime.GC()
	var before, live, after runtime.MemStats
	runtime.ReadMemStats(&before)
	beforeGoroutines := runtime.NumGoroutine()
	started := time.Now()
	t.Run("exchange", func(t *testing.T) {
		pairs := makePairs(t, connections, true)
		runtime.GC()
		runtime.ReadMemStats(&live)
		var workers sync.WaitGroup
		failures := make(chan error, concurrency)
		latencies := make([][]time.Duration, concurrency)
		for worker := range concurrency {
			workers.Go(func() {
				payload, buffer := make([]byte, payloadSize), make([]byte, payloadSize)
				for sequence := range exchanges {
					for i := worker; i < len(pairs); i += concurrency {
						binary.BigEndian.PutUint64(payload, uint64(i)*exchanges+uint64(sequence))
						var start time.Time
						if sequence%16 == 0 {
							start = time.Now()
						}
						if err := pairs[i].roundTrip(payload, buffer); err != nil {
							failures <- fmt.Errorf("connection %d exchange %d: %w", i, sequence, err)
							return
						}
						if !start.IsZero() {
							latencies[worker] = append(latencies[worker], time.Since(start))
						}
					}
				}
			})
		}
		workers.Wait()
		close(failures)
		for err := range failures {
			t.Error(err)
		}
		samples := slices.Concat(latencies...)
		slices.Sort(samples)
		if len(samples) > 0 {
			t.Logf("workers=%d round_trip_samples=%d p50=%s p95=%s p99=%s max=%s", concurrency, len(samples),
				samples[len(samples)/2], samples[len(samples)*95/100], samples[len(samples)*99/100], samples[len(samples)-1])
		}
	})
	deadline := time.Now().Add(5 * time.Second)
	for runtime.NumGoroutine() > beforeGoroutines+8 && time.Now().Before(deadline) {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
	// A second collection also drains unused sync.Pool buffers.
	runtime.GC()
	runtime.GC()
	runtime.ReadMemStats(&after)
	if got := runtime.NumGoroutine(); got > beforeGoroutines+8 {
		t.Errorf("goroutines did not converge: %d -> %d", beforeGoroutines, got)
	}
	if retained := int64(after.HeapAlloc) - int64(before.HeapAlloc); retained > 32<<20 {
		t.Errorf("retained heap grew by %d bytes after closing connections", retained)
	}
	t.Logf("connections=%d exchanges=%d payload_bytes=%d elapsed=%s live_heap=%d retained_heap=%d allocated=%d gc_cycles=%d gc_pause=%s goroutines=%d->%d",
		connections, connections*exchanges, 2*connections*exchanges*payloadSize, time.Since(started),
		int64(live.HeapAlloc)-int64(before.HeapAlloc), int64(after.HeapAlloc)-int64(before.HeapAlloc),
		after.TotalAlloc-before.TotalAlloc, after.NumGC-before.NumGC, time.Duration(after.PauseTotalNs-before.PauseTotalNs), beforeGoroutines, runtime.NumGoroutine())
}

func BenchmarkEarlyDataHandshakeLifecycle(b *testing.B) {
	for _, entries := range []int{64, 4096} {
		b.Run(fmt.Sprintf("Entries%d", entries), func(b *testing.B) {
			certificate, roots := testServerCertificate(b)
			clientConfig := &Config{
				RootCAs: roots, ServerName: "server.test", ClientSessionCache: NewLRUClientSessionCache(1),
				HandshakeTimeout: 5 * time.Second,
			}
			serverConfig := &Config{
				Certificates: []tls.Certificate{certificate}, MaxEarlyData: 1024, AllowEarlyDataWithoutCookie: true,
				EarlyDataReplayCache: NewLRUEarlyDataReplayCache(entries + b.N), HandshakeTimeout: 5 * time.Second,
			}
			serverConfig.SessionTicketKey[0] = 1
			expires := time.Now().Add(time.Hour)
			for i := range entries {
				if !serverConfig.EarlyDataReplayCache.CheckAndStore(fmt.Sprint(i), expires) {
					b.Fatal("could not seed replay cache")
				}
			}
			exchange := func(early bool) {
				left, right := memoryDatagramPair()
				server := Server(right, serverConfig)
				var client *Conn
				if early {
					client = ClientEarly(left, clientConfig)
				} else {
					client = Client(left, clientConfig)
				}
				defer client.Close()
				defer server.Close()
				serverDone := make(chan error, 1)
				if !early {
					go func() { serverDone <- server.Handshake() }()
				}
				var payload, buffer [32]byte
				if _, err := client.WriteDatagram(payload[:]); err != nil {
					b.Fatal(err)
				}
				if early {
					go func() { serverDone <- server.Handshake() }()
				}
				if err := client.Handshake(); err != nil {
					b.Fatal(err)
				}
				if err := <-serverDone; err != nil {
					b.Fatal(err)
				}
				if early && (client.ConnectionState().EarlyData != EarlyDataAccepted || server.ConnectionState().EarlyData != EarlyDataAccepted) {
					b.Fatal("resumed connection did not accept early data")
				}
				n, info, err := server.ReadDatagram(buffer[:])
				if err != nil || n != len(payload) || info.EarlyData != early || buffer != payload {
					b.Fatalf("early=%v: payload=%x info=%+v error=%v", early, buffer, info, err)
				}
				// Reading the response also drains the preceding session ticket.
				if _, err = server.WriteDatagram(payload[:]); err != nil {
					b.Fatal(err)
				}
				if n, info, err = client.ReadDatagram(buffer[:]); err != nil || n != len(payload) || info.Truncated || buffer != payload {
					b.Fatalf("response: payload=%x info=%+v error=%v", buffer, info, err)
				}
				if _, ok := clientConfig.ClientSessionCache.Get("server.test"); !ok {
					b.Fatal("missing ticket for the next connection")
				}
			}
			exchange(false)
			b.ReportAllocs()
			b.ResetTimer()
			for range b.N {
				exchange(true)
			}
		})
	}
}
