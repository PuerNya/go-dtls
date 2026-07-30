package dtls13

import (
	"bytes"
	"crypto/tls"
	"errors"
	"net"
	"testing"
	"time"
)

func completeHybridHandshake(t testing.TB, left, right net.Conn, clientConfig, serverConfig *Config) (*Conn, *Conn) {
	t.Helper()
	client := Client(left, clientConfig)
	server := Server(right, serverConfig)
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	serverErr := <-serverDone
	if clientErr != nil || serverErr != nil {
		t.Fatalf("hybrid handshake failed: client=%v server=%v", clientErr, serverErr)
	}
	return client, server
}

func TestHybridKeyExchangeCertificateHandshake(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	for _, group := range []tls.CurveID{tls.X25519MLKEM768, tls.SecP256r1MLKEM768, tls.SecP384r1MLKEM1024} {
		t.Run(group.String(), func(t *testing.T) {
			left, right := memoryDatagramPair()
			defer left.Close()
			defer right.Close()
			client, server := completeHybridHandshake(t, left, right,
				&Config{RootCAs: roots, ServerName: "server.test", CurvePreferences: []tls.CurveID{group}, SessionTicketsDisabled: true, MTU: 256, HandshakeTimeout: 3 * time.Second, FlightInterval: 5 * time.Millisecond},
				&Config{Certificates: []tls.Certificate{certificate}, CurvePreferences: []tls.CurveID{group}, SessionTicketsDisabled: true, MTU: 256, HandshakeTimeout: 3 * time.Second, FlightInterval: 5 * time.Millisecond})
			payload := []byte("hybrid application data")
			writeDone := make(chan error, 1)
			go func() { _, err := client.WriteDatagram(payload); writeDone <- err }()
			buffer := make([]byte, len(payload))
			n, _, err := server.ReadDatagram(buffer)
			if err != nil || !bytes.Equal(buffer[:n], payload) || <-writeDone != nil {
				t.Fatalf("hybrid application data = %q, %v", buffer[:n], err)
			}
		})
	}
}

func TestHybridKeyExchangeFallbackAndHelloRetryRequestShares(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	tests := []struct {
		name         string
		clientGroups []tls.CurveID
		serverGroups []tls.CurveID
		firstShares  []tls.CurveID
		secondShares []tls.CurveID
	}{
		{
			name: "ClassicalFallback", clientGroups: []tls.CurveID{tls.X25519MLKEM768, tls.X25519}, serverGroups: []tls.CurveID{tls.X25519},
			firstShares: []tls.CurveID{tls.X25519MLKEM768, tls.X25519}, secondShares: []tls.CurveID{tls.X25519MLKEM768, tls.X25519},
		},
		{
			name: "HybridGroupRetry", clientGroups: []tls.CurveID{tls.X25519, tls.X25519MLKEM768}, serverGroups: []tls.CurveID{tls.X25519MLKEM768},
			firstShares: []tls.CurveID{tls.X25519}, secondShares: []tls.CurveID{tls.X25519MLKEM768},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			left, right := memoryDatagramPair()
			defer left.Close()
			defer right.Close()
			capture := &captureWritesConn{Conn: left}
			completeHybridHandshake(t, capture, right,
				&Config{RootCAs: roots, ServerName: "server.test", CurvePreferences: test.clientGroups, SessionTicketsDisabled: true, HandshakeTimeout: 3 * time.Second, FlightInterval: 5 * time.Millisecond},
				&Config{Certificates: []tls.Certificate{certificate}, CurvePreferences: test.serverGroups, SessionTicketsDisabled: true, HandshakeTimeout: 3 * time.Second, FlightInterval: 5 * time.Millisecond})
			hellos := capturedClientHellos(t, capture)
			if len(hellos) != 2 {
				t.Fatalf("captured %d ClientHellos", len(hellos))
			}
			for i, want := range [][]tls.CurveID{test.firstShares, test.secondShares} {
				got := make([]tls.CurveID, len(hellos[i].keyShares))
				for j := range hellos[i].keyShares {
					got[j] = hellos[i].keyShares[j].group
				}
				if !bytes.Equal(curveIDsForTest(got), curveIDsForTest(want)) {
					t.Fatalf("ClientHello %d shares = %v, want %v", i, got, want)
				}
			}
		})
	}
}

func capturedClientHellos(t *testing.T, capture *captureWritesConn) []*clientHello {
	t.Helper()
	capture.mu.Lock()
	writes := append([][]byte(nil), capture.writes...)
	capture.mu.Unlock()
	reassembler := newReassembler()
	bySequence := make(map[uint16]*clientHello)
	for _, datagram := range writes {
		records, err := parsePlainRecords(datagram)
		if err != nil {
			continue
		}
		for _, record := range records {
			if record.typ != recordTypeHandshake {
				continue
			}
			fragments, err := parseHandshakeFragments(record.payload)
			if err != nil {
				t.Fatal(err)
			}
			for _, fragment := range fragments {
				if fragment.typ != handshakeTypeClientHello {
					continue
				}
				body, done, err := reassembler.add(fragment)
				if err != nil {
					t.Fatal(err)
				}
				if done {
					hello, err := parseClientHello(body)
					if err != nil {
						t.Fatal(err)
					}
					bySequence[fragment.messageSequence] = hello
				}
			}
		}
	}
	var hellos []*clientHello
	for sequence := uint16(0); ; sequence++ {
		hello := bySequence[sequence]
		if hello == nil {
			break
		}
		hellos = append(hellos, hello)
	}
	return hellos
}

func curveIDsForTest(groups []tls.CurveID) []byte {
	wire := make([]byte, 2*len(groups))
	for i, group := range groups {
		wire[2*i], wire[2*i+1] = byte(group>>8), byte(group)
	}
	return wire
}

func TestHybridKeyExchangeMutualTLSResumptionAndEarlyData(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	cache := NewLRUClientSessionCache(2)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x91}, len(ticketKey)))
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate}, ClientSessionCache: cache,
		CurvePreferences: []tls.CurveID{tls.X25519MLKEM768}, HandshakeTimeout: 3 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots,
		CurvePreferences: []tls.CurveID{tls.X25519MLKEM768}, SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour,
		MaxEarlyData: 1024, AllowEarlyDataWithoutCookie: true, EarlyDataReplayCache: NewLRUEarlyDataReplayCache(4),
		HandshakeTimeout: 3 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	_ = issueEarlyDataTicket(t, clientConfig, serverConfig)
	resumingConfig := clientConfig.Clone()
	resumingConfig.Certificates = nil
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	client := Client(left, resumingConfig)
	server := Server(right, serverConfig)
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	payload := []byte("hybrid mTLS early data")
	if n, err := client.WriteEarlyData(payload); err != nil || n != len(payload) {
		t.Fatalf("WriteEarlyData = %d, %v", n, err)
	}
	if err := <-serverDone; err != nil {
		t.Fatal(err)
	}
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume || len(server.ConnectionState().PeerCertificates) != 1 {
		t.Fatal("hybrid mTLS connection did not restore the authenticated session")
	}
	buffer := make([]byte, len(payload))
	n, _, err := server.ReadDatagram(buffer)
	if err != nil || !bytes.Equal(buffer[:n], payload) {
		t.Fatalf("early data = %q, %v", buffer[:n], err)
	}
}

func TestHybridKeyExchangeECHFragmentationAndWeakNetwork(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	configList, echKey := testECHConfig(t, "public.test", 21)
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	clientWire := &weakNetworkConn{Conn: left, enabled: true}
	serverWire := &weakNetworkConn{Conn: right, enabled: true}
	client, server := completeHybridHandshake(t, clientWire, serverWire,
		&Config{RootCAs: roots, ServerName: "server.test", EncryptedClientHelloConfigList: configList, CurvePreferences: []tls.CurveID{tls.X25519MLKEM768}, SessionTicketsDisabled: true, MTU: 600, HandshakeTimeout: 5 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond},
		&Config{Certificates: []tls.Certificate{certificate}, EncryptedClientHelloKeys: []EncryptedClientHelloKey{echKey}, CurvePreferences: []tls.CurveID{tls.X25519MLKEM768}, SessionTicketsDisabled: true, MTU: 600, HandshakeTimeout: 5 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond})
	if !client.ConnectionState().ECHAccepted || !server.ConnectionState().ECHAccepted {
		t.Fatal("hybrid weak-network handshake did not accept ECH")
	}
}

func TestHybridKeyExchangeRealUDP(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, CurvePreferences: []tls.CurveID{tls.X25519MLKEM768}, SessionTicketsDisabled: true,
		HandshakeTimeout: 3 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	serverDone := make(chan error, 1)
	go func() {
		server, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverDone <- acceptErr
			return
		}
		defer server.Close()
		buffer := make([]byte, 8)
		n, _, readErr := server.ReadDatagram(buffer)
		if readErr == nil && string(buffer[:n]) != "ping" {
			readErr = errors.New("unexpected hybrid UDP payload")
		}
		if readErr == nil {
			_, readErr = server.WriteDatagram([]byte("pong"))
		}
		serverDone <- readErr
	}()
	client, err := DialWithDialer(&net.Dialer{Timeout: 3 * time.Second}, "udp4", listener.Addr().String(), &Config{
		RootCAs: roots, ServerName: "server.test", CurvePreferences: []tls.CurveID{tls.X25519MLKEM768}, SessionTicketsDisabled: true,
		HandshakeTimeout: 3 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	_ = client.SetDeadline(time.Now().Add(3 * time.Second))
	if _, err = client.WriteDatagram([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 8)
	n, _, err := client.ReadDatagram(buffer)
	if err != nil || string(buffer[:n]) != "pong" {
		t.Fatalf("hybrid UDP response = %q, %v", buffer[:n], err)
	}
	if err = <-serverDone; err != nil {
		t.Fatal(err)
	}
}

func BenchmarkHybridKeyExchangeHandshakeLifecycle(b *testing.B) {
	certificate, roots := testServerCertificate(b)
	for _, group := range []tls.CurveID{tls.X25519MLKEM768, tls.SecP256r1MLKEM768, tls.SecP384r1MLKEM1024} {
		b.Run(group.String(), func(b *testing.B) {
			clientConfig := &Config{RootCAs: roots, ServerName: "server.test", CurvePreferences: []tls.CurveID{group}, SessionTicketsDisabled: true, HandshakeTimeout: time.Second}
			serverConfig := &Config{Certificates: []tls.Certificate{certificate}, CurvePreferences: []tls.CurveID{group}, SessionTicketsDisabled: true, HandshakeTimeout: time.Second}
			b.ReportAllocs()
			for b.Loop() {
				left, right := memoryDatagramPair()
				completeHybridHandshake(b, left, right, clientConfig, serverConfig)
				_ = left.Close()
				_ = right.Close()
			}
		})
	}
}
