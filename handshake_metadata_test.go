package dtls13

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func metadataConfigs(t *testing.T) (*Config, *Config) {
	t.Helper()
	certificate, roots := testServerCertificate(t)
	list, key := testECHConfig(t, "public.test", 21)
	return &Config{
		RootCAs: roots, ServerName: "server.test", EncryptedClientHelloConfigList: list,
		HandshakeMetadata: []byte("private client metadata"), SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}, &Config{
		Certificates: []tls.Certificate{certificate}, EncryptedClientHelloKeys: []EncryptedClientHelloKey{key},
		SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
		AcceptHandshakeMetadata: func(_ *ClientHelloInfo, _ []byte) ([]byte, error) { return []byte("private server metadata"), nil },
	}
}

func TestHandshakeMetadataExchange(t *testing.T) {
	for _, mode := range []string{"full", "sha384", "chacha", "ccm", "empty", "maximum", "hrr-loss", "resumed", "resumed-hrr", "resumed-no-cookie", "ticket-rejected", "udp"} {
		t.Run(mode, func(t *testing.T) {
			cc, sc := metadataConfigs(t)
			response := []byte("private server metadata")
			if suite := map[string]uint16{"sha384": TLS_AES_256_GCM_SHA384, "chacha": TLS_CHACHA20_POLY1305_SHA256, "ccm": TLS_AES_128_CCM_SHA256}[mode]; suite != 0 {
				cc.CipherSuites, sc.CipherSuites = []uint16{suite}, []uint16{suite}
			}
			if mode == "empty" {
				cc.HandshakeMetadata, response = []byte{}, nil
			}
			if mode == "maximum" {
				cc.HandshakeMetadata = bytes.Repeat([]byte{0x93}, maxHandshakeMetadata)
				response = bytes.Repeat([]byte{0xb7}, maxHandshakeMetadata)
			}
			resumed := mode == "resumed" || mode == "resumed-hrr" || mode == "resumed-no-cookie"
			if mode == "resumed-no-cookie" {
				sc.MaxEarlyData, sc.AllowEarlyDataWithoutCookie = 64, true
				sc.EarlyDataReplayCache = NewLRUEarlyDataReplayCache(4)
			}
			if resumed || mode == "ticket-rejected" {
				cc.SessionTicketsDisabled, sc.SessionTicketsDisabled = false, false
				cc.ClientSessionCache, sc.SessionTicketKey = NewLRUClientSessionCache(2), [32]byte{21}
				_ = issueEarlyDataTicket(t, cc, sc)
				cc.HandshakeMetadata = []byte("fresh resumption metadata")
				response = []byte("fresh resumption response")
			}
			if mode == "ticket-rejected" {
				sc.SessionTicketKey = [32]byte{22}
			}
			if mode == "resumed-hrr" {
				sc.CurvePreferences = []tls.CurveID{tls.CurveP256}
			}
			var calls atomic.Int32
			sc.AcceptHandshakeMetadata = func(info *ClientHelloInfo, request []byte) ([]byte, error) {
				calls.Add(1)
				if info.ServerName != "server.test" || !bytes.Equal(request, cc.HandshakeMetadata) || request == nil {
					return nil, errors.New("wrong metadata or inner SNI")
				}
				if info.Conn.ConnectionState().HandshakeMetadata != nil {
					return nil, errors.New("published metadata before handshake completion")
				}
				clear(request) // Callback ownership must not change the handshake transcript.
				return response, nil
			}
			if mode == "udp" {
				listener, err := Listen("udp4", "127.0.0.1:0", sc)
				if err != nil {
					t.Fatal(err)
				}
				defer listener.Close()
				accepted := make(chan *Conn, 1)
				serverErr := make(chan error, 1)
				go func() {
					s, err := listener.Accept()
					if err != nil {
						serverErr <- err
						return
					}
					accepted <- s
					serverErr <- s.Handshake()
				}()
				client, err := Dial("udp4", listener.Addr().String(), cc)
				if err != nil {
					t.Fatal(err)
				}
				defer client.Close()
				if err := <-serverErr; err != nil {
					t.Fatal(err)
				}
				server := <-accepted
				defer server.Close()
				checkMetadataState(t, client, server, cc.HandshakeMetadata, response, false)
			} else {
				left, right := memoryDatagramPair()
				clientCapture, serverCapture := &captureWritesConn{Conn: left}, &captureWritesConn{Conn: right}
				client, server := Client(clientCapture, cc), Server(serverCapture, sc)
				if mode == "resumed-no-cookie" {
					client = ClientEarly(clientCapture, cc)
				}
				loss := &handshakeLoss{typ: handshakeTypeEncryptedExtensions, epoch: 2}
				if mode == "hrr-loss" {
					cc.MTU, sc.MTU = 256, 256
					sc.CurvePreferences = []tls.CurveID{tls.CurveP256}
					server.conn = &observeRecordsConn{Conn: serverCapture, owner: server, observe: loss.observeMessage}
					client.conn = &observeRecordsConn{Conn: clientCapture, owner: client, observe: loss.observeACK}
				}
				handshakePair(t, client, server)
				checkMetadataState(t, client, server, cc.HandshakeMetadata, response, resumed)
				if mode == "hrr-loss" {
					loss.requireRecovery(t)
				}
				for _, capture := range []*captureWritesConn{clientCapture, serverCapture} {
					capture.mu.Lock()
					for _, wire := range capture.writes {
						if len(cc.HandshakeMetadata) > 0 && bytes.Contains(wire, cc.HandshakeMetadata) || len(response) > 0 && bytes.Contains(wire, response) {
							t.Error("metadata leaked onto the wire")
						}
					}
					capture.mu.Unlock()
				}
				// The default server policy adds a cookie HRR with or without metadata.
				wantHellos := 2
				if mode == "resumed-no-cookie" {
					wantHellos = 1
				}
				hellos := capturedClientHellos(t, clientCapture)
				if len(hellos) != wantHellos {
					t.Fatalf("ClientHello count=%d, want %d", len(hellos), wantHellos)
				}
				for _, hello := range hellos {
					if hello.handshakeMetadata != nil {
						t.Fatal("metadata in ClientHelloOuter")
					}
				}
			}
			if calls.Load() != 1 {
				t.Fatalf("callback calls = %d", calls.Load())
			}
		})
	}
}

func TestHandshakeMetadataWaitsForCompletion(t *testing.T) {
	for _, cancelHandshake := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "canceled"}[cancelHandshake], func(t *testing.T) {
			cc, sc := metadataConfigs(t)
			left, right := memoryDatagramPair()
			gate := &heldHandshakeConn{Conn: left}
			client, server := ClientEarly(gate, cc), Server(right, sc)
			defer client.Close()
			defer server.Close()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			server.startHandshake(context.Background())
			client.startHandshake(ctx)
			if err := client.waitReady(client.applicationReady); err != nil {
				t.Fatal(err)
			}
			if client.ConnectionState().HandshakeMetadata != nil || server.ConnectionState().HandshakeMetadata != nil {
				t.Fatal("metadata exposed before handshake completion")
			}
			if cancelHandshake {
				cancel()
				if err := client.Handshake(); !errors.Is(err, context.Canceled) {
					t.Fatalf("cancel: %v", err)
				}
				if client.ConnectionState().HandshakeMetadata != nil {
					t.Fatal("metadata exposed after cancellation")
				}
				return
			}
			if err := gate.flush(); err != nil {
				t.Fatal(err)
			}
			handshakePair(t, client, server)
			checkMetadataState(t, client, server, cc.HandshakeMetadata, []byte("private server metadata"), false)
		})
	}
}

func checkMetadataState(t *testing.T, client, server *Conn, request, response []byte, resumed bool) {
	t.Helper()
	for _, side := range []struct {
		conn *Conn
		want []byte
	}{{client, response}, {server, request}} {
		state := side.conn.ConnectionState()
		if !state.HandshakeComplete || !state.ECHAccepted || state.DidResume != resumed || state.HandshakeMetadata == nil || !bytes.Equal(state.HandshakeMetadata, side.want) {
			t.Fatalf("metadata=%x, complete=%v, ECH=%v, resumed=%v", state.HandshakeMetadata, state.HandshakeComplete, state.ECHAccepted, state.DidResume)
		}
		clear(state.HandshakeMetadata)
		if !bytes.Equal(side.conn.ConnectionState().HandshakeMetadata, side.want) {
			t.Fatal("state snapshot aliases connection metadata")
		}
	}
}

func TestHandshakeMetadataRejection(t *testing.T) {
	for _, mode := range []string{"no-ech", "grease-only", "oversized-request", "unsupported", "callback-error", "oversized-response", "ech-rejected", "close-in-callback"} {
		t.Run(mode, func(t *testing.T) {
			cc, sc := metadataConfigs(t)
			want := errors.New("metadata policy rejection")
			switch mode {
			case "no-ech", "grease-only":
				cc.EncryptedClientHelloConfigList = nil
				cc.EncryptedClientHelloGrease = mode == "grease-only"
			case "oversized-request":
				cc.HandshakeMetadata = make([]byte, maxHandshakeMetadata+1)
			case "unsupported":
				sc.AcceptHandshakeMetadata = nil
			case "callback-error":
				sc.AcceptHandshakeMetadata = func(*ClientHelloInfo, []byte) ([]byte, error) { return nil, want }
			case "oversized-response":
				sc.AcceptHandshakeMetadata = func(*ClientHelloInfo, []byte) ([]byte, error) { return make([]byte, maxHandshakeMetadata+1), nil }
			case "ech-rejected":
				cc.EncryptedClientHelloConfigList, _ = testECHConfig(t, "server.test", 22)
				sc.AcceptHandshakeMetadata = func(*ClientHelloInfo, []byte) ([]byte, error) { t.Error("callback on ECH rejection"); return nil, nil }
			case "close-in-callback":
				sc.AcceptHandshakeMetadata = func(info *ClientHelloInfo, _ []byte) ([]byte, error) { _ = info.Conn.Close(); return nil, nil }
			}
			left, right := memoryDatagramPair()
			client, server := Client(left, cc), Server(right, sc)
			defer client.Close()
			defer server.Close()
			done := make(chan error, 1)
			go func() { done <- server.Handshake() }()
			err := client.Handshake()
			if mode == "no-ech" || mode == "grease-only" || mode == "oversized-request" {
				_ = server.Close()
			}
			serverErr := <-done
			if err == nil {
				t.Fatal("client accepted missing/rejected metadata")
			}
			if mode == "no-ech" || mode == "grease-only" || mode == "oversized-request" {
				if _, ok := errors.AsType[*ConfigError](err); !ok {
					t.Fatalf("client configuration error: %v", err)
				}
			}
			if mode == "oversized-response" {
				if _, ok := errors.AsType[*ConfigError](serverErr); !ok {
					t.Fatalf("server response error: %v", serverErr)
				}
			}
			if mode == "callback-error" && !errors.Is(serverErr, want) {
				t.Fatalf("callback error = %v", serverErr)
			}
			if mode == "unsupported" {
				var alert *localAlertError
				if !errors.As(err, &alert) || alert.description != alertMissingExtension {
					t.Fatalf("missing acceptance: %v", err)
				}
			}
			if mode == "ech-rejected" {
				if _, ok := errors.AsType[*ECHRejectionError](err); !ok {
					t.Fatalf("ECH rejection lost: %v", err)
				}
			}
			if client.ConnectionState().HandshakeMetadata != nil || server.ConnectionState().HandshakeMetadata != nil {
				t.Fatal("metadata published after failed exchange")
			}
		})
	}
}

func TestHandshakeMetadataWireValidation(t *testing.T) {
	for _, test := range []struct {
		name      string
		raw       []byte
		offered   bool
		wantAlert uint8
	}{
		{"empty", []byte{1}, true, 0},
		{"unsolicited", []byte{1}, false, alertUnsupportedExtension},
		{"truncated", nil, true, alertDecodeError},
		{"version", []byte{2}, true, alertIllegalParameter},
		{"oversized", make([]byte, maxHandshakeMetadata+2), true, alertDecodeError},
	} {
		t.Run(test.name, func(t *testing.T) {
			hello := &clientHello{
				cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, signatureSchemes: defaultSignatureSchemes(),
				supportedGroups: []tls.CurveID{tls.X25519},
				keyShares:       []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
			}
			if test.offered {
				hello.handshakeMetadata = []byte{}
			}
			wire, err := (&encryptedExtensions{extensions: map[uint16][]byte{extHandshakeMetadata: test.raw}}).marshal()
			if err != nil {
				t.Fatal(err)
			}
			ee, err := parseEncryptedExtensions(wire)
			if err != nil {
				t.Fatal(err)
			}
			_, _, _, err = validateEncryptedExtensions(hello, &ee)
			errs := []error{err}
			if test.offered {
				body, marshalErr := hello.marshal()
				if marshalErr != nil {
					t.Fatal(marshalErr)
				}
				_, parseErr := parseClientHello(replaceClientHelloExtension(t, body, extHandshakeMetadata, test.raw))
				errs = append(errs, parseErr)
			}
			for _, err := range errs {
				if test.wantAlert == 0 {
					if err != nil || ee.handshakeMetadata == nil {
						t.Fatalf("empty exchange: %v", err)
					}
				} else if description, ok := protocolAlert(err); !ok || description != test.wantAlert {
					t.Fatalf("alert=%d, want %d: %v", description, test.wantAlert, err)
				}
			}
		})
	}
	cc, _ := metadataConfigs(t)
	cloned := cc.Clone()
	cloned.HandshakeMetadata[0] ^= 0xff
	if bytes.Equal(cloned.HandshakeMetadata, cc.HandshakeMetadata) {
		t.Fatal("Config.Clone aliases metadata")
	}
	cc.HandshakeMetadata = []byte{}
	if cc.Clone().HandshakeMetadata == nil {
		t.Fatal("Config.Clone loses empty offer")
	}
	initial := &clientHello{handshakeMetadata: []byte{}}
	if !equalClientHelloAfterHRR(initial, initial, 0) {
		t.Fatal("unchanged metadata rejected after HRR")
	}
	for _, data := range [][]byte{nil, []byte("b")} {
		second := *initial
		second.handshakeMetadata = data
		if equalClientHelloAfterHRR(initial, &second, 0) {
			t.Fatal("HRR accepted changed metadata")
		}
	}
}

// A server must not accept this extension from an unencrypted ClientHello.
func TestHandshakeMetadataRejectsPlaintext(t *testing.T) {
	cc, sc := metadataConfigs(t)
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	cc.EncryptedClientHelloConfigList = nil
	cc.HandshakeMetadata = nil
	client := Client(left, cc)
	client.config, _ = cc.normalized()
	state := &clientHandshakeState{}
	if err := client.clientPrepareHello(state); err != nil {
		t.Fatal(err)
	}
	state.hello.handshakeMetadata = []byte("plaintext metadata")
	body, err := state.hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	flight, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeClientHello, body: body}}, 1200, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	server := Server(right, sc)
	defer server.Close()
	done := make(chan error, 1)
	go func() { done <- server.Handshake() }()
	if err := client.writeFlight(left, flight); err != nil {
		t.Fatal(err)
	}
	err = <-done
	var alert *localAlertError
	if !errors.As(err, &alert) || alert.description != alertIllegalParameter {
		t.Fatalf("plaintext metadata: %v", err)
	}
}
