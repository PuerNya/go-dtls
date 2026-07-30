package dtls13

import (
	"crypto/tls"
	"runtime"
	"testing"
	"time"
)

func TestConnectionLifecycleResourceStability(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	const connections = 32
	runtime.GC()
	baselineGoroutines := runtime.NumGoroutine()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)

	for i := 0; i < connections; i++ {
		left, right := memoryDatagramPair()
		client := Client(left, &Config{
			RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true,
			HandshakeTimeout: time.Second, FlightInterval: 2 * time.Millisecond,
		})
		server := Server(right, &Config{
			Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true,
			HandshakeTimeout: time.Second, FlightInterval: 2 * time.Millisecond,
		})
		serverDone := make(chan error, 1)
		go func() { serverDone <- server.Handshake() }()
		if err := client.Handshake(); err != nil {
			t.Fatalf("connection %d client handshake: %v", i, err)
		}
		if err := <-serverDone; err != nil {
			t.Fatalf("connection %d server handshake: %v", i, err)
		}
		if _, err := client.WriteDatagram([]byte{byte(i)}); err != nil {
			t.Fatalf("connection %d write: %v", i, err)
		}
		payload := make([]byte, 1)
		if _, _, err := server.ReadDatagram(payload); err != nil || payload[0] != byte(i) {
			t.Fatalf("connection %d read payload=%x err=%v", i, payload, err)
		}
		_ = left.Close()
		_ = right.Close()
	}

	deadline := time.Now().Add(2 * time.Second)
	for runtime.NumGoroutine() > baselineGoroutines+8 && time.Now().Before(deadline) {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	if got := runtime.NumGoroutine(); got > baselineGoroutines+8 {
		stack := make([]byte, 1<<20)
		n := runtime.Stack(stack, true)
		t.Logf("goroutine dump:\n%s", stack[:n])
		t.Fatalf("goroutines did not converge: before=%d after=%d", baselineGoroutines, got)
	}
	if retained := int64(after.HeapAlloc) - int64(before.HeapAlloc); retained > 8<<20 {
		t.Fatalf("retained heap grew by %d bytes across %d connections", retained, connections)
	}
	averageAllocated := (after.TotalAlloc - before.TotalAlloc) / connections
	if averageAllocated > 512<<10 {
		t.Fatalf("average allocation is %d bytes per connection", averageAllocated)
	}
	t.Logf("connections=%d average_alloc=%dB retained_heap_delta=%dB goroutines=%d->%d",
		connections, averageAllocated, int64(after.HeapAlloc)-int64(before.HeapAlloc), baselineGoroutines, runtime.NumGoroutine())
}

func TestCertificateCompressionResourceStability(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	serverCertificate = compressibleTestCertificate(serverCertificate)
	clientCertificate = compressibleTestCertificate(clientCertificate)
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate},
		EnableCertificateCompression: true, SessionTicketsDisabled: true, HandshakeTimeout: time.Second,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots,
		EnableCertificateCompression: true, SessionTicketsDisabled: true, HandshakeTimeout: time.Second,
	}
	runtime.GC()
	baselineGoroutines := runtime.NumGoroutine()
	var before runtime.MemStats
	runtime.ReadMemStats(&before)
	for range 16 {
		left, right := memoryDatagramPair()
		client := Client(left, clientConfig)
		server := Server(right, serverConfig)
		serverDone := make(chan error, 1)
		go func() { serverDone <- server.Handshake() }()
		clientErr := client.Handshake()
		serverErr := <-serverDone
		_ = left.Close()
		_ = right.Close()
		if clientErr != nil || serverErr != nil {
			t.Fatalf("handshake failed: client=%v server=%v", clientErr, serverErr)
		}
	}
	deadline := time.Now().Add(2 * time.Second)
	for runtime.NumGoroutine() > baselineGoroutines+8 && time.Now().Before(deadline) {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
	runtime.GC()
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	if got := runtime.NumGoroutine(); got > baselineGoroutines+8 {
		t.Fatalf("goroutines did not converge: before=%d after=%d", baselineGoroutines, got)
	}
	if retained := int64(after.HeapAlloc) - int64(before.HeapAlloc); retained > 8<<20 {
		t.Fatalf("retained heap grew by %d bytes", retained)
	}
}
