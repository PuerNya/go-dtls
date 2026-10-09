package dtls13

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Holding epoch 2 also holds ACKs, without needing access to either peer's keys.
type heldHandshakeConn struct {
	net.Conn
	mu      sync.Mutex
	read    bool
	release bool
	held    [][]byte
}

func (c *heldHandshakeConn) Write(p []byte) (int, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.read && !c.release && isUnifiedRecord(p) && p[0]&unifiedHeaderEpochMask == 2 {
		c.held = append(c.held, bytes.Clone(p))
		return len(p), nil
	}
	return c.Conn.Write(p)
}

func (c *heldHandshakeConn) Read(p []byte) (int, error) {
	for {
		c.mu.Lock()
		if c.read && c.release && len(c.held) != 0 {
			wire := c.held[0]
			c.held = c.held[1:]
			c.mu.Unlock()
			return copy(p, wire), nil
		}
		c.mu.Unlock()
		n, err := c.Conn.Read(p)
		if err != nil {
			return n, err
		}
		c.mu.Lock()
		hold := c.read && !c.release && isUnifiedRecord(p[:n]) && p[0]&unifiedHeaderEpochMask == 2
		if hold {
			c.held = append(c.held, bytes.Clone(p[:n]))
		}
		c.mu.Unlock()
		if !hold {
			return n, nil
		}
	}
}

func (c *heldHandshakeConn) flush() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.release = true
	if c.read {
		return c.SetReadDeadline(time.Now())
	}
	for _, p := range c.held {
		if _, err := c.Conn.Write(p); err != nil {
			return err
		}
	}
	c.held = nil
	return nil
}

func earlyIOConfigs(t *testing.T) (*Config, *Config) {
	t.Helper()
	certificate, roots := testServerCertificate(t)
	client := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: NewLRUClientSessionCache(4),
		HandshakeTimeout: 2 * time.Second, FlightInterval: 10 * time.Millisecond,
	}
	server := &Config{
		Certificates: []tls.Certificate{certificate}, MaxEarlyData: 1024,
		EnableEarlyDataIO: true, AllowEarlyDataWithoutCookie: true,
		EarlyDataReplayCache: NewLRUEarlyDataReplayCache(16),
		HandshakeTimeout:     2 * time.Second, FlightInterval: 10 * time.Millisecond,
	}
	copy(server.SessionTicketKey[:], bytes.Repeat([]byte{0xec}, 32))
	return client, server
}

func expectEarlyIODatagram(t *testing.T, conn *Conn, want string, early bool) {
	t.Helper()
	result := make(chan error, 1)
	go func() {
		var buffer [2048]byte
		n, info, err := conn.ReadDatagram(buffer[:])
		if err == nil && (string(buffer[:n]) != want || info.EarlyData != early || info.FullLength != len(want)) {
			err = fmt.Errorf("received %q %+v; want %q early=%t", buffer[:n], info, want, early)
		}
		result <- err
	}()
	if err := waitCloseHandshakeResult(t, result); err != nil {
		t.Fatal(err)
	}
}

func writeEarlyIODatagram(t *testing.T, conn *Conn, payload string) {
	t.Helper()
	if n, err := conn.WriteDatagram([]byte(payload)); err != nil || n != len(payload) {
		t.Fatalf("WriteDatagram(%q) = %d, %v", payload, n, err)
	}
}

func exchangeEarlyBeforeFinished(t *testing.T, conn *Conn, server bool) {
	t.Helper()
	gate := &heldHandshakeConn{Conn: conn.conn, read: server}
	conn.conn = gate
	conn.earlyIO = true
	if server {
		expectEarlyIODatagram(t, conn, "early", true)
		writeEarlyIODatagram(t, conn, "early reply")
	} else {
		writeEarlyIODatagram(t, conn, "early")
		expectEarlyIODatagram(t, conn, "early reply", false)
	}
	select {
	case <-conn.HandshakeComplete():
		t.Fatal("early exchange waited for client Finished")
	default:
	}
	if err := gate.flush(); err != nil {
		t.Fatal(err)
	}
}

func TestEarlyIOExchangeBeforeFinished(t *testing.T) {
	for _, suite := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384, TLS_CHACHA20_POLY1305_SHA256, TLS_AES_128_CCM_SHA256} {
		t.Run(fmt.Sprintf("%x", suite), func(t *testing.T) {
			clientConfig, serverConfig := earlyIOConfigs(t)
			clientConfig.CipherSuites, serverConfig.CipherSuites = []uint16{suite}, []uint16{suite}
			issueEarlyDataTicket(t, clientConfig, serverConfig)
			left, right := memoryDatagramPair()
			gate := &heldHandshakeConn{Conn: left}
			capture := &captureWritesConn{Conn: right}
			client, server := ClientEarly(gate, clientConfig), Server(capture, serverConfig)
			t.Cleanup(func() { _ = client.Close(); _ = server.Close() })

			// No server is reading yet: every write must use the same epoch-1 cipher.
			for _, payload := range []string{"one", "", "three"} {
				writeEarlyIODatagram(t, client, payload)
			}
			if client.ConnectionState().EarlyData != EarlyDataPending {
				t.Fatal("early offer was not pending")
			}
			for _, payload := range []string{"one", "", "three"} {
				expectEarlyIODatagram(t, server, payload, true)
			}
			writeEarlyIODatagram(t, server, "before Finished")
			expectEarlyIODatagram(t, client, "before Finished", false)
			for _, conn := range []*Conn{client, server} {
				select {
				case <-conn.HandshakeComplete():
					t.Fatal("handshake completed through a held Finished")
				default:
				}
				if conn.ConnectionState().EarlyData != EarlyDataAccepted {
					t.Fatal("early offer was not accepted")
				}
			}
			var replay []byte
			capture.mu.Lock()
			for _, wire := range capture.writes {
				if isUnifiedRecord(wire) && wire[0]&unifiedHeaderEpochMask == 3 {
					replay = bytes.Clone(wire)
					break
				}
			}
			capture.mu.Unlock()
			if replay == nil {
				t.Fatal("response did not use epoch 3")
			}
			if err := gate.flush(); err != nil {
				t.Fatal(err)
			}
			handshakePair(t, client, server)
			// A replay spanning handshake completion must not be delivered again.
			if _, err := right.Write(replay); err != nil {
				t.Fatal(err)
			}
			writeEarlyIODatagram(t, server, "after Finished")
			expectEarlyIODatagram(t, client, "after Finished", false)
			writeEarlyIODatagram(t, client, "ordinary")
			expectEarlyIODatagram(t, server, "ordinary", false)
		})
	}
}

func TestEarlyIOPendingFinalACK(t *testing.T) {
	for _, mode := range []string{"ticket", "external-psk", "background-failure", "handshake-application"} {
		t.Run(mode, func(t *testing.T) {
			clientConfig, serverConfig := earlyIOConfigs(t)
			if mode == "external-psk" {
				psk, err := ImportExternalPSK([]byte("device"), bytes.Repeat([]byte{0x91}, 32), []byte("roles"), 0)
				if err != nil {
					t.Fatal(err)
				}
				clientConfig.ExternalPSKs, serverConfig.ExternalPSKs = []*ExternalPSK{psk}, []*ExternalPSK{psk}
			} else {
				issueEarlyDataTicket(t, clientConfig, serverConfig)
			}
			left, right := memoryDatagramPair()
			clientGate := &heldHandshakeConn{Conn: left}
			serverGate := &heldHandshakeConn{Conn: right, release: true}
			client, server := ClientEarly(clientGate, clientConfig), Server(serverGate, serverConfig)
			t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
			var dropACK atomic.Bool
			dropACK.Store(true)
			serverGate.Conn = &observeRecordsConn{Conn: right, owner: server, observe: func(r record) (bool, error) {
				if r.typ == recordTypeACK && r.epoch >= 3 && dropACK.Load() {
					numbers, err := parseACK(r.payload)
					for _, number := range numbers {
						if number.epoch == 2 {
							return true, err
						}
					}
					return false, err
				}
				return false, nil
			}}
			client.startHandshake(context.Background())
			server.startHandshake(context.Background())
			if err := client.waitReady(client.applicationReady); err != nil {
				t.Fatal(err)
			}
			// Hold the final ACK, letting tickets and KeyUpdate precede the response.
			serverGate.mu.Lock()
			serverGate.release = false
			serverGate.mu.Unlock()
			if err := clientGate.flush(); err != nil {
				t.Fatal(err)
			}
			if err := server.Handshake(); err != nil {
				t.Fatal(err)
			}
			if err := server.SendKeyUpdate(true); err != nil {
				t.Fatal(err)
			}
			writeEarlyIODatagram(t, server, "pending ACK response")
			expectEarlyIODatagram(t, client, "pending ACK response", false)
			select {
			case <-client.HandshakeComplete():
				t.Fatal("client completed without its final ACK")
			default:
			}
			if _, err := client.ConnectionState().ExportKeyingMaterial("test", nil, 32); err == nil {
				t.Fatal("exporter available before handshake completion")
			}
			ticket, ok := clientConfig.ClientSessionCache.Get("server.test")
			if !ok {
				t.Fatal("ticket before the final ACK was not processed")
			}
			if mode == "external-psk" && (ticket.externalPSK == nil || string(client.ConnectionState().ExternalPSKIdentity()) != "device") {
				t.Fatal("early ticket lost its external PSK authentication origin")
			}
			if mode == "background-failure" {
				failure := errors.New("post-handshake retransmission failed")
				client.failConnection(failure)
				if err := client.Handshake(); !errors.Is(err, failure) {
					t.Fatalf("handshake failure = %v", err)
				}
				select {
				case <-client.HandshakeComplete():
					t.Fatal("failed handshake reported completion")
				default:
				}
				return
			}
			if mode == "handshake-application" {
				server.writeMu.Lock()
				wire, err := server.finishedACKCipher.seal(recordTypeApplicationData, []byte("wrong epoch"))
				server.writeMu.Unlock()
				if err != nil {
					t.Fatal(err)
				}
				if _, err = right.Write(wire); err != nil {
					t.Fatal(err)
				}
				if description, ok := protocolAlert(client.Handshake()); !ok || description != alertUnexpectedMessage {
					t.Fatalf("handshake-key application alert = %d, %v", description, ok)
				}
				return
			}
			dropACK.Store(false)
			if err := serverGate.flush(); err != nil {
				t.Fatal(err)
			}
			handshakePair(t, client, server)
			writeEarlyIODatagram(t, client, "after ACK")
			expectEarlyIODatagram(t, server, "after ACK", false)
		})
	}
}

func TestEarlyIOLimitAndRejection(t *testing.T) {
	for _, mode := range []string{"quota", "HRR", "replay", "ordinary-server"} {
		t.Run(mode, func(t *testing.T) {
			clientConfig, serverConfig := earlyIOConfigs(t)
			serverConfig.MaxEarlyData = 3
			ticket := issueEarlyDataTicket(t, clientConfig, serverConfig)
			switch mode {
			case "HRR":
				serverConfig.AllowEarlyDataWithoutCookie = false
			case "replay":
				serverConfig.EarlyDataReplayCache.CheckAndStore(string(ticket.ticket), time.Now().Add(time.Hour))
			case "ordinary-server":
				serverConfig.EnableEarlyDataIO = false
			}
			left, right := memoryDatagramPair()
			client, server := ClientEarly(left, clientConfig), Server(right, serverConfig)
			t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
			writeEarlyIODatagram(t, client, "123")
			if n, err := client.WriteDatagram(make([]byte, 65536)); n != 0 || !errors.Is(err, ErrDatagramTooLarge) {
				t.Fatalf("oversized early write = %d, %v", n, err)
			}
			// The next datagram cannot fit the early quota; it must wait and use epoch 3.
			done := make(chan error, 1)
			go func() { _, err := client.WriteDatagram([]byte("next")); done <- err }()
			server.startHandshake(context.Background())
			if mode == "quota" || mode == "ordinary-server" {
				expectEarlyIODatagram(t, server, "123", true)
			}
			if err := waitCloseHandshakeResult(t, done); err != nil {
				t.Fatal(err)
			}
			expectEarlyIODatagram(t, server, "next", false)
			handshakePair(t, client, server)
			want := EarlyDataAccepted
			if mode == "HRR" || mode == "replay" {
				want = EarlyDataRejected
			}
			if client.ConnectionState().EarlyData != want {
				t.Fatalf("early status = %v, want %v", client.ConnectionState().EarlyData, want)
			}
		})
	}
}

func TestEarlyIOLifecycle(t *testing.T) {
	for _, mode := range []string{"close", "cancel", "record-failure", "bad-server", "no-ticket"} {
		t.Run(mode, func(t *testing.T) {
			clientConfig, serverConfig := earlyIOConfigs(t)
			if mode != "no-ticket" {
				issueEarlyDataTicket(t, clientConfig, serverConfig)
			}
			left, right := memoryDatagramPair()
			client, server := ClientEarly(left, clientConfig), Server(right, serverConfig)
			t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			// Notifications and construction must not consume a cached ticket or do I/O.
			select {
			case <-client.HandshakeComplete():
				t.Fatal("idle handshake completed")
			case <-client.Context().Done():
				t.Fatal("idle connection closed")
			default:
			}
			if len(right.(*memoryDatagramConn).in) != 0 {
				t.Fatal("lazy client performed I/O")
			}
			client.startHandshake(ctx)
			if mode == "no-ticket" {
				handshakePair(t, client, server)
				cancel()
				writeEarlyIODatagram(t, client, "fallback")
				expectEarlyIODatagram(t, server, "fallback", false)
				if client.ConnectionState().EarlyData != EarlyDataNotAttempted {
					t.Fatal("uncached handshake attempted early data")
				}
				return
			}
			if err := client.waitReady(client.writeReady); err != nil {
				t.Fatal(err)
			}
			want := error(net.ErrClosed)
			switch mode {
			case "close":
				_ = client.Close()
			case "cancel":
				cancel()
				want = context.Canceled
			case "record-failure":
				want = errors.New("record write failed")
				client.failConnection(want)
			case "bad-server":
				serverConfig.SessionTicketsDisabled = true
				certificate, _ := testServerCertificate(t)
				serverConfig.Certificates = []tls.Certificate{certificate}
				server.startHandshake(context.Background())
				want = nil
			}
			select {
			case <-client.Context().Done():
			case <-time.After(3 * time.Second):
				t.Fatal("connection failure was not observable")
			}
			if err := context.Cause(client.Context()); err == nil || (want != nil && !errors.Is(err, want)) {
				t.Fatalf("connection cause = %v, want %v", err, want)
			}
			if _, _, err := client.ReadDatagram(make([]byte, 1)); err == nil {
				t.Fatal("failed connection remained readable")
			}
			if err := client.Handshake(); err == nil {
				t.Fatal("failed early connection completed its handshake")
			}
			select {
			case <-client.HandshakeComplete():
				t.Fatal("failed handshake reported success")
			default:
			}
		})
	}
}

func TestEarlyIOClientAuthentication(t *testing.T) {
	for _, required := range []bool{false, true} {
		t.Run(fmt.Sprint(required), func(t *testing.T) {
			clientConfig, serverConfig := earlyIOConfigs(t)
			serverConfig.ClientAuth = tls.RequestClientCert
			if required {
				serverConfig.ClientAuth = tls.RequireAnyClientCert
			}
			left, right := memoryDatagramPair()
			gate := &heldHandshakeConn{Conn: left}
			client, server := ClientEarly(gate, clientConfig), Server(right, serverConfig)
			t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
			client.startHandshake(context.Background())
			done := make(chan error, 1)
			go func() { _, err := server.WriteDatagram([]byte("authenticated")); done <- err }()
			if err := client.waitReady(client.applicationReady); err != nil {
				t.Fatal(err)
			}
			select {
			case <-server.writeReady:
				t.Fatal("server allowed writing before client authentication")
			default:
			}
			if err := gate.flush(); err != nil {
				t.Fatal(err)
			}
			err := waitCloseHandshakeResult(t, done)
			if (err != nil) != required {
				t.Fatalf("server write error = %v, client certificate required = %v", err, required)
			}
			if !required {
				expectEarlyIODatagram(t, client, "authenticated", false)
			}
		})
	}
}

func TestEarlyIOConcurrentWrites(t *testing.T) {
	clientConfig, serverConfig := earlyIOConfigs(t)
	issueEarlyDataTicket(t, clientConfig, serverConfig)
	left, right := memoryDatagramPair()
	client, server := ClientEarly(left, clientConfig), Server(right, serverConfig)
	t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
	const count = 16
	done := make(chan error, count)
	for i := range count {
		go func() { _, err := client.WriteDatagram([]byte{byte(i), 1}); done <- err }()
	}
	for range count {
		if err := waitCloseHandshakeResult(t, done); err != nil {
			t.Fatal(err)
		}
	}
	seen := make(map[byte]bool)
	for range count {
		var buffer [1]byte
		n, info, err := server.ReadDatagram(buffer[:])
		if err != nil || n != 1 || !info.EarlyData || !info.Truncated || info.FullLength != 2 || seen[buffer[0]] {
			t.Fatalf("concurrent early record = %x, %+v, %v", buffer, info, err)
		}
		seen[buffer[0]] = true
	}
	handshakePair(t, client, server)
}

func TestDialEarly(t *testing.T) {
	clientConfig, serverConfig := earlyIOConfigs(t)
	listener, err := Listen("udp4", "127.0.0.1:0", serverConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	accepted := make(chan *Conn, 1)
	go func() { conn, _ := listener.Accept(); accepted <- conn }()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// No cached ticket: DialEarly must complete a full handshake before returning.
	dialed := make(chan *Conn, 1)
	go func() { conn, _ := DialEarly(ctx, "udp4", listener.Addr().String(), clientConfig); dialed <- conn }()
	server := <-accepted
	if server == nil {
		t.Fatal("accept failed")
	}
	defer server.Close()
	if err = server.HandshakeContext(ctx); err != nil {
		t.Fatal(err)
	}
	client := <-dialed
	if client == nil || !client.ConnectionState().HandshakeComplete {
		t.Fatal("uncached DialEarly did not complete the handshake")
	}
	defer client.Close()
	// An application response follows the ticket, making cache readiness observable.
	writeEarlyIODatagram(t, server, "ticket ready")
	expectEarlyIODatagram(t, client, "ticket ready", false)
	if _, ok := clientConfig.ClientSessionCache.Get("server.test"); !ok {
		t.Fatal("session ticket was not cached")
	}
	go func() { conn, _ := listener.Accept(); accepted <- conn }()
	resumed, err := DialEarly(ctx, "udp4", listener.Addr().String(), clientConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer resumed.Close()
	if resumed.ConnectionState().HandshakeComplete {
		t.Fatal("resumed DialEarly waited for the idle server")
	}
	writeEarlyIODatagram(t, resumed, "early UDP")
	server2 := <-accepted
	if server2 == nil {
		t.Fatal("second accept failed")
	}
	defer server2.Close()
	expectEarlyIODatagram(t, server2, "early UDP", true)
	writeEarlyIODatagram(t, server2, "UDP reply")
	expectEarlyIODatagram(t, resumed, "UDP reply", false)
	handshakePair(t, resumed, server2)
	if _, err = DialEarly(ctx, "tcp", listener.Addr().String(), clientConfig); err == nil {
		t.Fatal("DialEarly accepted a stream transport")
	}
}

func TestEarlyIOServerAmplificationLimit(t *testing.T) {
	clientConfig, serverConfig := earlyIOConfigs(t)
	clientConfig.IgnorePathMTU, serverConfig.IgnorePathMTU = true, true
	issueEarlyDataTicket(t, clientConfig, serverConfig)
	left, right := memoryDatagramPair()
	gate := &heldHandshakeConn{Conn: left}
	client, server := ClientEarly(gate, clientConfig), Server(right, serverConfig)
	t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
	writeEarlyIODatagram(t, client, "request")
	expectEarlyIODatagram(t, server, "request", true)
	done := make(chan error, 1)
	payload := bytes.Repeat([]byte("x"), 8000)
	go func() { _, err := server.WriteDatagram(payload); done <- err }()
	// The receive budget is much smaller than the response; it must not be sent.
	select {
	case err := <-done:
		t.Fatalf("amplification-limited write returned before validation: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	if err := gate.flush(); err != nil {
		t.Fatal(err)
	}
	if err := waitCloseHandshakeResult(t, done); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, len(payload))
	n, info, err := client.ReadDatagram(buffer)
	if err != nil || info.Truncated || !bytes.Equal(buffer[:n], payload) {
		t.Fatalf("post-validation response: n=%d info=%+v error=%v", n, info, err)
	}
}
