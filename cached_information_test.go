package dtls13

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"net"
	"sync"
	"testing"
	"time"
)

func cachedInformationConfigs(t *testing.T) (*Config, *Config) {
	t.Helper()
	cert, roots := testServerCertificate(t)
	clientCert, clientRoots := testClientCertificate(t)
	return &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCert},
		CachedInformationCache: NewCachedInformationCache(2), SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}, &Config{
		Certificates: []tls.Certificate{cert}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots,
		EnableCachedInformation: true, SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
}

// Observe authenticated wire messages so successful handshakes cannot hide a
// failure to negotiate caching. Reuse the existing record observer and loss gate.
func cachedInformationPair(t *testing.T, cc, sc *Config, want uint8, loss *handshakeLoss) (*Conn, *Conn) {
	t.Helper()
	left, right := memoryDatagramPair()
	client, server := Client(left, cc), Server(right, sc)
	if loss != nil {
		client.conn = &observeRecordsConn{Conn: client.conn, owner: client, observe: loss.observeACK}
	}
	check := observeCachedInformation(t, server, want, loss)
	handshakePair(t, client, server)
	check()
	return client, server
}

func observeCachedInformation(t *testing.T, server *Conn, want uint8, loss *handshakeLoss) func() {
	t.Helper()
	var mu sync.Mutex
	seen := make(map[uint8]int)
	server.conn = &observeRecordsConn{Conn: server.conn, owner: server, observe: func(r record) (bool, error) {
		if r.typ == recordTypeHandshake && r.epoch == 2 {
			fragments, err := parseHandshakeFragments(r.payload)
			if err != nil {
				return false, err
			}
			mu.Lock()
			for _, f := range fragments {
				seen[f.typ] = int(f.length)
			}
			mu.Unlock()
		}
		if loss != nil {
			return loss.observeMessage(r)
		}
		return false, nil
	}}
	return func() {
		t.Helper()
		if loss != nil {
			loss.requireRecovery(t)
		}
		mu.Lock()
		defer mu.Unlock()
		for i, typ := range []uint8{handshakeTypeCertificate, handshakeTypeCertificateRequest} {
			if got := seen[typ] == 33; got != (want&uint8(i+1) != 0) {
				t.Fatalf("message %d length=%d, cached selection=%d", typ, seen[typ], want)
			}
		}
	}
}

func TestCachedInformationHandshake(t *testing.T) {
	for _, mode := range []string{"mutual", "server-only", "compression", "RPK", "SHA384", "HRR-loss", "ECH", "server-disabled", "client-disabled"} {
		t.Run(mode, func(t *testing.T) {
			cc, sc := cachedInformationConfigs(t)
			want := cachedCert | cachedCertRequest
			var loss *handshakeLoss
			switch mode {
			case "server-only":
				sc.ClientAuth = tls.NoClientCert
				want = cachedCert
			case "compression":
				cc.EnableCertificateCompression, sc.EnableCertificateCompression = true, true
				sc.Certificates[0] = compressibleTestCertificate(sc.Certificates[0])
			case "RPK":
				cc, sc = rawPublicKeyConfigs(t)
				cc.CachedInformationCache, sc.EnableCachedInformation = NewCachedInformationCache(2), true
			case "SHA384":
				cc.CipherSuites, sc.CipherSuites = []uint16{TLS_AES_256_GCM_SHA384}, []uint16{TLS_AES_256_GCM_SHA384}
			case "HRR-loss":
				sc.CurvePreferences = []tls.CurveID{tls.CurveP256}
				cc.MTU, sc.MTU = 256, 256
				loss = &handshakeLoss{typ: handshakeTypeCertificate, epoch: 2}
			case "ECH":
				list, key := testECHConfig(t, "public.test", 7)
				cc.EncryptedClientHelloConfigList, sc.EncryptedClientHelloKeys = list, []EncryptedClientHelloKey{key}
			case "server-disabled":
				sc.EnableCachedInformation = false
				want = 0
			case "client-disabled":
				cc.CachedInformationCache = nil
				want = 0
			}
			cachedInformationPair(t, cc, sc, 0, nil)
			client, server := cachedInformationPair(t, cc, sc, want, loss)
			if mode == "ECH" && !client.ConnectionState().ECHAccepted {
				t.Fatal("ECH not accepted")
			}
			for _, pair := range [][2]*Conn{{client, server}, {server, client}} {
				_ = pair[1].SetReadDeadline(time.Now().Add(time.Second))
				if _, err := pair[0].WriteDatagram([]byte("cached")); err != nil {
					t.Fatal(err)
				}
				var data [16]byte
				if n, _, err := pair[1].ReadDatagram(data[:]); err != nil || string(data[:n]) != "cached" {
					t.Fatalf("application data: %d %v", n, err)
				}
			}
		})
	}
}

func TestCachedInformationRotationAndTrust(t *testing.T) {
	cc, sc := cachedInformationConfigs(t)
	cc.EnableOCSPStapling, sc.EnableOCSPStapling = true, true
	sc.Certificates[0].OCSPStaple = []byte("old")
	cachedInformationPair(t, cc, sc, 0, nil)
	rotated := sc.Clone()
	rotated.Certificates[0].OCSPStaple = []byte("new")
	client, _ := cachedInformationPair(t, cc, rotated, cachedCertRequest, nil)
	if !bytes.Equal(client.ConnectionState().OCSPResponse, []byte("new")) {
		t.Fatal("stale cached OCSP response")
	}
	cachedInformationPair(t, cc, rotated, cachedCert|cachedCertRequest, nil)
	rotated.EnableOCSPStapling = false
	cachedInformationPair(t, cc, rotated, cachedCert, nil)
	for _, mode := range []string{"roots", "callback", "expiry", "identity"} {
		t.Run(mode, func(t *testing.T) {
			changed := cc.Clone()
			switch mode {
			case "roots":
				changed.RootCAs = x509.NewCertPool()
			case "callback":
				changed.VerifyPeerCertificate = func([][]byte, [][]*x509.Certificate) error { return errors.New("revoked") }
			case "expiry":
				changed.Time = func() time.Time { return time.Date(2200, 1, 1, 0, 0, 0, 0, time.UTC) }
			case "identity":
				changed.ServerName = "wrong.test"
				changed.CachedInformationCache = NewCachedInformationCache(1)
				changed.CachedInformationCache.put("wrong.test", cc.CachedInformationCache.get("server.test"))
			}
			left, right := memoryDatagramPair()
			client, server := Client(left, changed), Server(right, rotated)
			defer client.Close()
			defer server.Close()
			done := make(chan error, 1)
			go func() { done <- server.Handshake() }()
			if err := client.Handshake(); err == nil {
				t.Fatal("cached identity bypassed current trust")
			}
			if err := <-done; err == nil {
				t.Fatal("server accepted rejected client handshake")
			}
		})
	}
	// A newly authenticated credential replaces the old entry only on success.
	rotated.Certificates[0], cc.RootCAs = testServerCertificate(t)
	before := cc.CachedInformationCache.get("server.test")
	rejected := cc.Clone()
	rejected.VerifyPeerCertificate = func([][]byte, [][]*x509.Certificate) error { return errors.New("reject rotated identity") }
	left, right := memoryDatagramPair()
	client, server := Client(left, rejected), Server(right, rotated)
	defer client.Close()
	defer server.Close()
	done := make(chan error, 1)
	go func() { done <- server.Handshake() }()
	if err := client.Handshake(); err == nil {
		t.Fatal("accepted rejected identity")
	}
	if err := <-done; err == nil {
		t.Fatal("peer accepted failed authentication")
	}
	if got := cc.CachedInformationCache.get("server.test"); !bytes.Equal(got[0], before[0]) || !bytes.Equal(got[1], before[1]) {
		t.Fatal("failed handshake updated cache")
	}
	cachedInformationPair(t, cc, rotated, cachedCertRequest, nil)
	cachedInformationPair(t, cc, rotated, cachedCert|cachedCertRequest, nil)

	cc, sc = rawPublicKeyConfigs(t)
	cc.CachedInformationCache, sc.EnableCachedInformation = NewCachedInformationCache(1), true
	cachedInformationPair(t, cc, sc, 0, nil)
	cc.VerifyPeerRawPublicKey = pinRawPublicKey(rawPublicKeyDER(t, ed25519.NewKeyFromSeed(make([]byte, ed25519.SeedSize))))
	left, right = memoryDatagramPair()
	client, server = Client(left, cc), Server(right, sc)
	defer client.Close()
	defer server.Close()
	go func() { done <- server.Handshake() }()
	if err := client.Handshake(); err == nil {
		t.Fatal("cached RPK bypassed changed pin")
	}
	if err := <-done; err == nil {
		t.Fatal("peer accepted failed RPK authentication")
	}
}

func TestCachedInformationResumptionAndPHA(t *testing.T) {
	cc, sc := cachedInformationConfigs(t)
	cc.SessionTicketsDisabled, sc.SessionTicketsDisabled = false, false
	cc.ClientSessionCache, cc.PostHandshakeAuth = NewLRUClientSessionCache(2), true
	sc.SessionTicketKey = [32]byte{3}
	sc.MaxEarlyData, sc.AllowEarlyDataWithoutCookie = 1024, true
	issueEarlyDataTicket(t, cc, sc)
	left, right := memoryDatagramPair()
	client, server := ClientEarly(left, cc), Server(right, sc)
	t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
	check := observeCachedInformation(t, server, 0, nil)
	done := make(chan error, 1)
	go func() { done <- server.Handshake() }()
	writeEarlyTestDatagram(t, client, []byte("early"))
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	check()
	_ = server.SetReadDeadline(time.Now().Add(time.Second))
	var data [16]byte
	if n, info, err := server.ReadDatagram(data[:]); err != nil || !info.EarlyData || string(data[:n]) != "early" {
		t.Fatalf("early data: %d %v", n, err)
	}
	if !client.ConnectionState().DidResume {
		t.Fatal("did not resume")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.RequestClientCertificate(ctx); err != nil {
		t.Fatal(err)
	}
	full := cc.Clone()
	full.SessionTicketsDisabled = true
	_, server = cachedInformationPair(t, full, sc, cachedCert|cachedCertRequest, nil)
	if err := server.RequestClientCertificate(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestCachedInformationCache(t *testing.T) {
	cache := NewCachedInformationCache(2)
	body := []byte("one")
	cache.put("a", [2][]byte{body})
	body[0] = 'x'
	if string(cache.get("a")[0]) != "one" {
		t.Fatal("cache retained caller memory")
	}
	cache.put("b", [2][]byte{[]byte("two")})
	snapshot := cache.get("a")
	cache.put("c", [2][]byte{[]byte("three")})
	if cache.get("b")[0] != nil {
		t.Fatal("LRU eviction failed")
	}
	cache.put("a", [2][]byte{bytes.Repeat([]byte{1}, maxCachedInformationSize+1)})
	if cache.get("a")[0] != nil || string(snapshot[0]) != "one" {
		t.Fatal("size limit or snapshot lifetime failed")
	}
	var wg sync.WaitGroup
	for range 4 {
		wg.Go(func() {
			for range 20 {
				cache.put("a", [2][]byte{[]byte("new")})
				_ = cache.get("a")
			}
		})
	}
	wg.Wait()
	var disabled CachedInformationCache
	disabled.put("a", [2][]byte{body})
	if disabled.get("a")[0] != nil {
		t.Fatal("zero-value cache enabled")
	}
}

func TestCachedInformationAuthentication(t *testing.T) {
	cc, sc := cachedInformationConfigs(t)
	cachedInformationPair(t, cc, sc, 0, nil)
	bodies := cc.CachedInformationCache.get("server.test")
	config, err := cc.normalized()
	if err != nil {
		t.Fatal(err)
	}
	conn := &Conn{config: config}
	state := &clientHandshakeState{hello: &clientHello{signatureSchemes: defaultSignatureSchemes()},
		cachedInfo: &clientCachedInformation{offered: bodies, selected: cachedCert | cachedCertRequest},
		transcript: newTranscriptHash(sha256.New())}
	want := sha256.New()
	for i, typ := range []uint8{handshakeTypeCertificateRequest, handshakeTypeCertificate} {
		body := bodies[1-i]
		hash := cachedInformationHash(typ, body)
		reference := append([]byte{32}, hash[:]...)
		message := completedHandshake{typ: typ, sequence: uint16(i), body: reference}
		if typ == handshakeTypeCertificateRequest {
			err = conn.clientCertificateRequest(state, message)
		} else {
			err = conn.clientServerCertificate(state, message)
		}
		if err != nil {
			t.Fatal(err)
		}
		_, _ = want.Write([]byte{typ, 0, 0, 33})
		_, _ = want.Write(reference)
	}
	if !bytes.Equal(state.transcript.sum(), want.Sum(nil)) {
		t.Fatal("transcript used restored body instead of transmitted fingerprint")
	}
	for _, test := range []struct {
		name  string
		alert uint8
		run   func(*clientHandshakeState) error
	}{
		{"missing CertificateRequest", alertIllegalParameter, func(state *clientHandshakeState) error {
			state.certificateRequest = nil
			return conn.clientServerCertificate(state, completedHandshake{typ: handshakeTypeCertificate})
		}},
		{"compressed certificate after selection", alertUnexpectedMessage, func(state *clientHandshakeState) error {
			_, err := conn.parseCachedServerCertificate(state, completedHandshake{typ: handshakeTypeCompressedCertificate})
			return err
		}},
		{"full certificate after selection", alertDecodeError, func(state *clientHandshakeState) error {
			_, err := conn.parseCachedServerCertificate(state, completedHandshake{typ: handshakeTypeCertificate, body: bodies[0]})
			return err
		}},
		{"selection in PSK", alertIllegalParameter, func(state *clientHandshakeState) error {
			conn.offerCachedInformation(state, state.hello)
			state.usingPSK = true
			body, err := (&encryptedExtensions{extensions: map[uint16][]byte{extCachedInfo: {0, 1, cachedCert}}}).marshal()
			if err != nil {
				return err
			}
			return conn.clientEncryptedExtensions(state, completedHandshake{typ: handshakeTypeEncryptedExtensions, body: body})
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			copy := *state
			if err := test.run(&copy); err == nil {
				t.Fatal("accepted invalid cached authentication")
			} else if alert, ok := protocolAlert(err); !ok || alert != test.alert {
				t.Fatalf("alert=%d error=%v", alert, err)
			}
		})
	}
}

func TestCachedInformationWire(t *testing.T) {
	requireAlert := func(err error, want uint8) {
		t.Helper()
		if got, ok := protocolAlert(err); !ok || got != want {
			t.Fatalf("alert=%d, want=%d: %v", got, want, err)
		}
	}
	// RFC 7924 Appendix A, independently specified TLS 1.2 Certificate vector.
	der, err := hex.DecodeString("3082022c308201b2a00302010202010d300a06082a8648ce3d040302303e310b3009060355040613024e4c3111300f060355040a1308506f6c617253534c311c301a06035504031313506f6c617273736c2054657374204543204341301e170d3133303932343135353230345a170d3233303932323135353230345a3041310b3009060355040613024e4c3111300f060355040a1308506f6c617253534c311f301d06035504031316506f6c617253534c205465737420436c69656e7420323059301306072a8648ce3d020106082a8648ce3d0301070342000457e5aeb173dfd3acbb93b881ff12aeeee653acce5553f6340ecc2ee363250bdf98e2f35c603696c0d5181470e57f9fd54b4518e5b06cd55cf8968f8770a3e4c7a3819d30819a30090603551d1304023000301d0603551d0e041604147a005f8664fce05de511103bb2e63bc4263fcfe2306e0603551d230467306580149d6d202449013f2bcb78b519bc7e24c9dbfb367ca142a440303e310b3009060355040613024e4c3111300f060355040a1308506f6c617253534c311c301a06035504031313506f6c617273736c2054657374204543204341820900c143e27e6243cce8300a06082a8648ce3d040302036800306502304a650d7b2083a299b9a80ffc8dee8f3dbb704c9603ac8e7870ddf20ea0b216cb658e1ac93f2c617ef83cefad1cee36200231009df227a6d574b824aee16a3f31a1ca542f08d08dee4f0c61df77787db4fdfc4249eee5b26ac2cd2677628e287c9e5745")
	if err != nil {
		t.Fatal(err)
	}
	vector := append([]byte{0, 2, 0x33, 0, 2, 0x30}, der...)
	vectorHash := cachedInformationHash(handshakeTypeCertificate, vector)
	if hex.EncodeToString(vectorHash[:]) != "086eefb4859adfe977defac494fff6b73033b4ce1f86b8f2a9fc0c6bf98605af" {
		t.Fatalf("RFC 7924 fingerprint=%x", vectorHash)
	}
	body := []byte("full certificate")
	hash := cachedInformationHash(handshakeTypeCertificate, body)
	offer := append([]byte{0, 34, cachedCert, 32}, hash[:]...)
	hello := &clientHello{unknownExtensions: map[uint16][]byte{extCachedInfo: offer}}
	for _, raw := range [][]byte{nil, {0, 0}, {0, 2, 1, 0}, {0, 3, 1, 1, 7}, append(bytes.Clone(offer), 0)} {
		_, err := parseCachedInformationOffer(raw)
		requireAlert(err, alertDecodeError)
	}
	for _, test := range []struct {
		raw   []byte
		alert uint8
	}{
		{nil, alertDecodeError}, {[]byte{0, 0}, alertDecodeError},
		{[]byte{0, 1, 2}, alertIllegalParameter}, {[]byte{0, 1, 255}, alertIllegalParameter},
	} {
		_, err := validateCachedInformationSelection(hello, test.raw)
		requireAlert(err, test.alert)
	}
	// RFC 7924 forbids unoffered types, but not repeated supported selections.
	for _, raw := range [][]byte{{0, 1, 1}, {0, 2, 1, 1}} {
		if got, err := validateCachedInformationSelection(hello, raw); err != nil || got != cachedCert {
			t.Fatalf("selection=%d %v", got, err)
		}
	}
	_, err = validateCachedInformationSelection(&clientHello{}, []byte{0, 1, 1})
	requireAlert(err, alertUnsupportedExtension)
	// Multiple fingerprints of one type are legal; unknown types are ignored.
	multiple := append(append([]byte{0, 71, 99, 1, 0}, offer[2:]...), offer[2:]...)
	objects, err := parseCachedInformationOffer(multiple)
	if err != nil || len(objects) != 2 {
		t.Fatalf("multiple offers: %v %v", objects, err)
	}
	encoded, selected := cachedHandshakeBody(handshakeTypeCertificate, body, objects)
	state := &clientCachedInformation{offered: [2][]byte{body}, selected: selected}
	if decoded, err := state.resolve(handshakeTypeCertificate, encoded, 1024); err != nil || !bytes.Equal(decoded, body) {
		t.Fatalf("fingerprint resolution: %x %v", decoded, err)
	}
	for _, test := range []struct {
		raw   []byte
		alert uint8
	}{
		{encoded[:32], alertDecodeError}, {append(bytes.Clone(encoded), 0), alertDecodeError},
		{append([]byte{32}, make([]byte, 32)...), alertIllegalParameter},
	} {
		_, err := state.resolve(handshakeTypeCertificate, test.raw, 1024)
		requireAlert(err, test.alert)
	}
	_, err = state.resolve(handshakeTypeCertificate, encoded, 1)
	requireAlert(err, alertIllegalParameter)
	outer, err := makeECHOuter(hello, &echConfig{publicName: "public.test"}, bytes.NewReader(make([]byte, 32)))
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := outer.unknownExtensions[extCachedInfo]; ok {
		t.Fatal("ECH outer exposes cache fingerprints")
	}
}

func TestInteropWolfSSLCachedInformationFallback(t *testing.T) {
	t.Run("server", func(t *testing.T) {
		var offered bool
		testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
			args:                    []string{"-C", "2"},
			connections:             2,
			verifyServerCertificate: true,
			configure: func(t *testing.T, root string, config *Config) {
				config.RootCAs = wolfSSLRootCAs(t, root, "certs/ca-cert.pem")
				config.ServerName = "example.com"
				config.CachedInformationCache = NewCachedInformationCache(1)
				config.SessionTicketsDisabled = true
			},
			wrapClientConn: func(conn *Conn) {
				offered = false
				conn.conn = &observeRecordsConn{Conn: conn.conn, owner: conn, observe: func(r record) (bool, error) {
					if r.typ != recordTypeHandshake || r.epoch != 0 {
						return false, nil
					}
					fragments, err := parseHandshakeFragments(r.payload)
					for _, f := range fragments {
						if f.typ == handshakeTypeClientHello && f.offset == 0 && len(f.body) == int(f.length) {
							hello, parseErr := parseClientHello(f.body)
							if parseErr != nil {
								return false, parseErr
							}
							offered = hello.unknownExtensions[extCachedInfo] != nil
						}
					}
					return false, err
				}}
			},
			connected: func(t *testing.T, conn *Conn, index int) {
				state := conn.ConnectionState()
				if state.DidResume || len(state.VerifiedChains) == 0 {
					t.Fatal("expected verified full handshake")
				}
				if offered != (index == 1) {
					t.Fatalf("connection %d cached_info offered=%v", index, offered)
				}
			},
		})
	})
	t.Run("client", func(t *testing.T) {
		testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
			args:                    []string{"-A", "certs/ca-cert.pem"},
			verifyServerCertificate: true,
			configure:               func(_ *testing.T, _ string, config *Config) { config.EnableCachedInformation = true },
		})
	})
}

func TestCachedInformationRealUDP(t *testing.T) {
	cc, sc := cachedInformationConfigs(t)
	listener, err := Listen("udp4", "127.0.0.1:0", sc)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
	defer stop()
	for round := range 2 {
		wire, err := net.Dial("udp4", listener.Addr().String())
		if err != nil {
			t.Fatal(err)
		}
		client := Client(wire, cc)
		defer client.Close()
		done := make(chan error, 1)
		go func() { done <- client.HandshakeContext(ctx) }()
		server, err := listener.Accept()
		if err != nil {
			t.Fatal(err)
		}
		defer server.Close()
		var want uint8
		if round > 0 {
			want = cachedCert | cachedCertRequest
		}
		check := observeCachedInformation(t, server, want, nil)
		serverErr := server.HandshakeContext(ctx)
		if clientErr := <-done; clientErr != nil || serverErr != nil {
			t.Fatalf("UDP handshake: %v / %v", clientErr, serverErr)
		}
		check()
		_ = server.SetDeadline(time.Now().Add(time.Second))
		_ = client.SetDeadline(time.Now().Add(time.Second))
		if _, err := client.WriteDatagram([]byte("cached UDP")); err != nil {
			t.Fatal(err)
		}
		var data [32]byte
		if n, _, err := server.ReadDatagram(data[:]); err != nil || string(data[:n]) != "cached UDP" {
			t.Fatalf("UDP data: %d %v", n, err)
		}
		_ = client.Close()
		_ = server.Close()
	}
}
