package dtls13

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
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
	run := func(b *testing.B, spec packetProcessSpec) {
		b.StopTimer()
		spec.Operations = b.N
		process := newPacketProcessRun(b, spec)
		b.StartTimer()
		stats := process.stage(b, "traffic")
		b.StopTimer()
		process.close(b)
		client, server, relay := stats["client"], stats["server"], stats["relay"]
		if client.Received != uint64(b.N) || server.Received != uint64(b.N) {
			b.Fatal("incomplete UDP round trips")
		}
		b.SetBytes(int64(2 * client.Payload))
		b.ReportMetric(float64(2*client.Payload*b.N)/client.Elapsed.Seconds()/1e6, "MB/s")
		reportPacketProcessBenchmark(b, client, server, relay)
	}
	for _, count := range []int{1, 64, 1024} {
		for _, cid := range []bool{false, true} {
			cidName := "NoCID"
			if cid {
				cidName = "CID"
			}
			for _, size := range []int{64, 1200} {
				name := fmt.Sprintf("Connections%d/%s/Bytes%d", count, cidName, size)
				spec := packetProcessSpec{Connections: count, Payload: size, CID: cid}
				b.Run(name, func(b *testing.B) {
					run(b, spec)
				})
				if count == 1024 && cid && size == 1200 {
					spec.NoServerReadDeadline = true
					b.Run("NoServerReadDeadline/"+name, func(b *testing.B) {
						run(b, spec)
					})
				}
			}
		}
	}
}

func reportPacketProcessBenchmark(b *testing.B, client, server, relay packetProcessStats) {
	b.ReportMetric(float64(client.Elapsed)/float64(b.N), "ns/op")
	b.ReportMetric(float64(client.Allocated+server.Allocated)/float64(b.N), "B/op")
	b.ReportMetric(float64(client.Allocs+server.Allocs)/float64(b.N), "allocs/op")
	for role, resources := range map[string]packetProcessStats{"client": client, "server": server, "relay": relay} {
		b.ReportMetric(float64(resources.Allocated)/float64(b.N), role+"_B/op")
		b.ReportMetric(float64(resources.Resources.User+resources.Resources.System)/float64(b.N), role+"_cpu_ns/op")
	}
}

func BenchmarkListenerHandshakeLifecycle(b *testing.B) {
	for _, resumed := range []bool{false, true} {
		name := "Full"
		if resumed {
			name = "Resumed"
		}
		for _, concurrency := range []int{1, 64} {
			b.Run(fmt.Sprintf("%s/Concurrency%d", name, concurrency), func(b *testing.B) {
				b.StopTimer()
				process := newPacketProcessRun(b, packetProcessSpec{Lifecycle: true, Resume: resumed, Connections: min(concurrency, b.N), Operations: b.N})
				b.StartTimer()
				stats := process.stage(b, "lifecycle")
				b.StopTimer()
				process.close(b)
				for _, role := range []string{"client", "server"} {
					if stats[role].Handshakes != uint64(b.N) {
						b.Fatal("incomplete UDP connection lifecycles")
					}
				}
				reportPacketProcessBenchmark(b, stats["client"], stats["server"], stats["relay"])
			})
		}
	}
}

func TestListenerDatagramLoad(t *testing.T) {
	if os.Getenv("GO_DTLS_STRESS") != "1" {
		t.Skip("set GO_DTLS_STRESS=1 for the large datagram workload")
	}
	connections, operations, workers := 1024, 1024*1024, 16
	for name, target := range map[string]*int{"GO_DTLS_PROCESS_CONNECTIONS": &connections, "GO_DTLS_PROCESS_OPERATIONS": &operations, "GO_DTLS_STRESS_WORKERS": &workers} {
		if value := os.Getenv(name); value != "" {
			parsed, err := strconv.Atoi(value)
			if err != nil || parsed < 1 {
				t.Fatalf("invalid %s", name)
			}
			*target = parsed
		}
	}
	run := newPacketProcessRun(t, packetProcessSpec{Connections: connections, Operations: operations, Workers: workers, Payload: 1200, CID: true})
	stats := run.stage(t, "traffic")
	if stats["client"].Received != uint64(operations) || stats["server"].Received != uint64(operations) {
		t.Fatal("incomplete round trips")
	}
	run.close(t)
}

func TestMemoryDatagramLoad(t *testing.T) {
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
		pairs := newMemoryLoadPairs(t, connections, true)
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
