package dtls13

import (
	"bytes"
	"crypto/ecdh"
	"crypto/hpke"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net"
	"testing"
	"time"
)

func testECHConfig(t testing.TB, publicName string, configID uint8) ([]byte, EncryptedClientHelloKey) {
	t.Helper()
	kem := hpke.DHKEM(ecdh.X25519())
	privateKey, err := kem.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	privateKeyBytes, err := privateKey.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	contents := newWireBuilder(128)
	contents.u8(int(configID))
	contents.u16(int(kem.ID()))
	contents.bytes16(privateKey.PublicKey().Bytes())
	cipherSuites := contents.startVector16()
	contents.u16(int(hpke.HKDFSHA256().ID()))
	contents.u16(int(hpke.AES128GCM().ID()))
	contents.endVector16(cipherSuites)
	contents.u8(64)
	contents.bytes8([]byte(publicName))
	contents.bytes16(nil)
	if contents.err != nil {
		t.Fatal(contents.err)
	}
	config := newWireBuilder(4 + len(contents.b))
	config.u16(int(extECH))
	config.u16(len(contents.b))
	config.b = append(config.b, contents.b...)
	if config.err != nil {
		t.Fatal(config.err)
	}
	list := newWireBuilder(2 + len(config.b))
	list.bytes16(config.b)
	if list.err != nil {
		t.Fatal(list.err)
	}
	key := EncryptedClientHelloKey{Config: config.b, PrivateKey: privateKeyBytes, SendAsRetry: true}
	return list.b, key
}

func testECHClientHello() *clientHello {
	return &clientHello{
		cipherSuites:                  []uint16{TLS_AES_128_GCM_SHA256},
		supportedGroups:               []tls.CurveID{tls.X25519},
		keyShares:                     []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
		signatureSchemes:              []tls.SignatureScheme{tls.Ed25519},
		serverName:                    "server.test",
		alpn:                          []string{"coap"},
		encryptedClientHelloExtension: []byte{echInnerType},
	}
}

func TestECHConfigVectorsAndValidation(t *testing.T) {
	presentation := "AEj+DQBEAQAgACAdd+scUi0IYFsXnUIU7ko2Nd9+F8M26pAGZVpz/KrWPgAEAAEAAWQVZWNoLXNpdGVzLmV4YW1wbGUubmV0AAA="
	wire, err := base64.StdEncoding.DecodeString(presentation)
	if err != nil {
		t.Fatal(err)
	}
	configs, err := parseECHConfigList(wire)
	if err != nil || len(configs) != 1 || configs[0].publicName != "ech-sites.example.net" {
		t.Fatalf("RFC 9848 ECHConfigList = %#v, %v", configs, err)
	}
	if config, _, _, _ := pickECHConfig(configs); config == nil {
		t.Fatal("RFC 9848 ECHConfigList contains no supported configuration")
	}

	cloudflare, err := hex.DecodeString("0045fe0d0041590020002092a01233db2218518ccbbbbc24df20686af417b37388de6460e94011974777090004000100010012636c6f7564666c6172652d6563682e636f6d0000")
	if err != nil {
		t.Fatal(err)
	}
	configs, err = parseECHConfigList(cloudflare)
	if err != nil || len(configs) != 1 || configs[0].publicName != "cloudflare-ech.com" {
		t.Fatalf("public ECHConfigList = %#v, %v", configs, err)
	}

	list, key := testECHConfig(t, "public.test", 7)
	if _, err = (&Config{EncryptedClientHelloConfigList: list}).normalized(); err != nil {
		t.Fatalf("valid client configuration: %v", err)
	}
	if _, err = (&Config{EncryptedClientHelloKeys: []EncryptedClientHelloKey{key}}).normalized(); err != nil {
		t.Fatalf("valid server configuration: %v", err)
	}
	badList := append([]byte(nil), list...)
	badList[1]++
	if _, err = (&Config{EncryptedClientHelloConfigList: badList}).normalized(); err == nil {
		t.Fatal("accepted malformed ECHConfigList")
	}
	badKey := key
	badKey.PrivateKey = append([]byte(nil), key.PrivateKey...)
	badKey.PrivateKey[len(badKey.PrivateKey)/2] ^= 1
	if _, err = (&Config{EncryptedClientHelloKeys: []EncryptedClientHelloKey{badKey}}).normalized(); err == nil {
		t.Fatal("accepted an ECH private key that does not match its configuration")
	}
	badKey = key
	badKey.SendAsRetry = false
	if _, err = (&Config{EncryptedClientHelloKeys: []EncryptedClientHelloKey{badKey}}).normalized(); err == nil {
		t.Fatal("accepted ECH keys with no retry configuration")
	}
}

func TestECHEncodingAuthenticatesOuterAndPadding(t *testing.T) {
	list, key := testECHConfig(t, "public.test", 9)
	context, err := newECHClientContext(list)
	if err != nil {
		t.Fatal(err)
	}
	inner := testECHClientHello()
	outer, err := makeECHOuter(inner, context.config, rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	body, err := computeOuterECH(outer, inner, context, true)
	if err != nil {
		t.Fatal(err)
	}
	parsedOuter, err := parseClientHello(body)
	if err != nil {
		t.Fatal(err)
	}
	decoded, _, serverContext, err := processECHClientHello(parsedOuter, body, []EncryptedClientHelloKey{key})
	if err != nil || serverContext == nil || decoded.serverName != inner.serverName || !bytes.Equal(decoded.encryptedClientHello(), []byte{echInnerType}) {
		t.Fatalf("ECH decode = %#v, context=%v, err=%v", decoded, serverContext, err)
	}

	encoded, err := encodeInnerClientHello(inner, int(context.config.maxNameLen))
	if err != nil {
		t.Fatal(err)
	}
	if len(encoded)%32 != 0 {
		t.Fatalf("EncodedClientHelloInner length = %d, want a multiple of 32", len(encoded))
	}
	badPadding := append([]byte(nil), encoded...)
	badPadding[len(badPadding)-1] = 1
	if _, _, err = decodeInnerClientHello(parsedOuter, body, badPadding); err == nil {
		t.Fatal("accepted non-zero EncodedClientHelloInner padding")
	}

	serverName, err := marshalServerName("qublic.test")
	if err != nil {
		t.Fatal(err)
	}
	tamperedBody := replaceClientHelloExtension(t, body, extServerName, serverName)
	tampered, err := parseClientHello(tamperedBody)
	if err != nil {
		t.Fatal(err)
	}
	effective, _, accepted, err := processECHClientHello(tampered, tamperedBody, []EncryptedClientHelloKey{key})
	if err != nil || accepted != nil || effective.serverName != "qublic.test" {
		t.Fatalf("tampered outer was not cleanly rejected: effective=%#v context=%v err=%v", effective, accepted, err)
	}
}

func encodedInnerWithOuterReferences(t *testing.T, encoded []byte, remove, references []uint16) []byte {
	t.Helper()
	p := wireParser{b: encoded}
	p.take(2 + 32)
	p.bytes8()
	p.bytes8()
	p.bytes16()
	p.bytes8()
	extensionsOffset := p.off
	p.bytes16()
	if p.err != nil {
		t.Fatal(p.err)
	}
	padding := encoded[p.off:]
	extensions, err := parseOrderedExtensionsView(encoded[extensionsOffset:p.off], nil)
	if err != nil {
		t.Fatal(err)
	}
	removed := make(map[uint16]bool, len(remove))
	for _, typ := range remove {
		removed[typ] = true
	}
	refs := newWireBuilder(1 + 2*len(references))
	values := make([]byte, 0, 2*len(references))
	for _, typ := range references {
		values = append(values, byte(typ>>8), byte(typ))
	}
	refs.bytes8(values)
	if refs.err != nil {
		t.Fatal(refs.err)
	}
	w := newWireBuilder(len(encoded))
	w.b = append(w.b, encoded[:extensionsOffset]...)
	start := w.startVector16()
	inserted := false
	for _, extension := range extensions {
		if removed[extension.typ] {
			if !inserted {
				w.u16(int(extECHOuterExtensions))
				w.bytes16(refs.b)
				inserted = true
			}
			continue
		}
		w.u16(int(extension.typ))
		w.bytes16(extension.value)
	}
	w.endVector16(start)
	w.b = append(w.b, padding...)
	if w.err != nil {
		t.Fatal(w.err)
	}
	return w.b
}

func TestECHOuterExtensionsReconstruction(t *testing.T) {
	inner := testECHClientHello()
	encoded, err := encodeInnerClientHello(inner, 0)
	if err != nil {
		t.Fatal(err)
	}
	outer := cloneClientHello(inner)
	outer.serverName = "public.test"
	outerECH, err := generateOuterECHExt(1, echCipher{kdfID: 1, aeadID: 1}, nil, []byte{1})
	if err != nil {
		t.Fatal(err)
	}
	outer.setEncryptedClientHello(outerECH)
	outerBody, err := outer.marshal()
	if err != nil {
		t.Fatal(err)
	}
	outer, err = parseClientHello(outerBody)
	if err != nil {
		t.Fatal(err)
	}
	valid := encodedInnerWithOuterReferences(t, encoded,
		[]uint16{extSupportedGroups, extSignatureAlgorithms},
		[]uint16{extSupportedGroups, extSignatureAlgorithms})
	decoded, _, err := decodeInnerClientHello(outer, outerBody, valid)
	if err != nil || len(decoded.supportedGroups) != 1 || len(decoded.signatureSchemes) != 1 {
		t.Fatalf("valid OuterExtensions decode = %#v, %v", decoded, err)
	}
	for _, test := range []struct {
		name string
		refs []uint16
	}{
		{"missing", []uint16{0xbeef}},
		{"duplicate", []uint16{extSupportedGroups, extSupportedGroups}},
		{"encrypted-client-hello", []uint16{extECH}},
		{"out-of-order", []uint16{extSignatureAlgorithms, extSupportedGroups}},
	} {
		t.Run(test.name, func(t *testing.T) {
			invalid := encodedInnerWithOuterReferences(t, encoded,
				[]uint16{extSupportedGroups, extSignatureAlgorithms}, test.refs)
			if _, _, err := decodeInnerClientHello(outer, outerBody, invalid); err == nil {
				t.Fatal("accepted invalid OuterExtensions")
			}
		})
	}
}

func TestECHEndToEndHRRFragmentationAndWeakNetwork(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	list, key := testECHConfig(t, "public.test", 11)
	left, right := memoryDatagramPair()
	clientWire := &weakNetworkConn{Conn: left, enabled: true}
	serverWire := &weakNetworkConn{Conn: right, enabled: true}
	var outerName, innerName string
	var outerALPN, innerALPN []string
	client := Client(clientWire, &Config{
		RootCAs: roots, ServerName: "server.test", NextProtos: []string{"coap"},
		EncryptedClientHelloConfigList: list, SessionTicketsDisabled: true, MTU: 256,
		HandshakeTimeout: 5 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	server := Server(serverWire, &Config{
		GetEncryptedClientHelloKeys: func(info *ClientHelloInfo) ([]EncryptedClientHelloKey, error) {
			outerName, outerALPN = info.ServerName, append([]string(nil), info.SupportedProtos...)
			return []EncryptedClientHelloKey{key}, nil
		},
		GetCertificate: func(info *ClientHelloInfo) (*tls.Certificate, error) {
			innerName, innerALPN = info.ServerName, append([]string(nil), info.SupportedProtos...)
			return &certificate, nil
		},
		NextProtos: []string{"coap"}, CurvePreferences: []tls.CurveID{tls.CurveP256},
		SessionTicketsDisabled: true, MTU: 256,
		HandshakeTimeout: 5 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	defer left.Close()
	defer right.Close()
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	serverErr := <-serverDone
	if clientErr != nil || serverErr != nil {
		t.Fatalf("ECH weak-network handshake failed: client=%v server=%v", clientErr, serverErr)
	}
	if !client.ConnectionState().ECHAccepted || !server.ConnectionState().ECHAccepted {
		t.Fatalf("ECH state: client=%#v server=%#v", client.ConnectionState(), server.ConnectionState())
	}
	if outerName != "public.test" || len(outerALPN) != 0 || innerName != "server.test" || len(innerALPN) != 1 || innerALPN[0] != "coap" {
		t.Fatalf("outer=(%q,%v), inner=(%q,%v)", outerName, outerALPN, innerName, innerALPN)
	}
}

type corruptFinalECHConfirmationConn struct {
	net.Conn
	sawHRR bool
}

func (c *corruptFinalECHConfirmationConn) Write(p []byte) (int, error) {
	records, err := parsePlainRecords(p)
	if err != nil || len(records) != 1 || records[0].typ != recordTypeHandshake {
		return c.Conn.Write(p)
	}
	fragments, err := parseHandshakeFragments(records[0].payload)
	if err != nil || len(fragments) != 1 || fragments[0].typ != handshakeTypeServerHello ||
		fragments[0].offset != 0 || int(fragments[0].length) != len(fragments[0].body) {
		return c.Conn.Write(p)
	}
	body := fragments[0].body
	if len(body) < 34 {
		return c.Conn.Write(p)
	}
	if bytes.Equal(body[2:34], helloRetryRequestRandom[:]) {
		c.sawHRR = true
		return c.Conn.Write(p)
	}
	if !c.sawHRR {
		return c.Conn.Write(p)
	}
	body = append([]byte(nil), body...)
	body[33] ^= 1
	fragment := fragments[0]
	fragment.body = body
	payload, err := marshalHandshakeFragment(fragment)
	if err != nil {
		return 0, err
	}
	wire, err := marshalPlainRecord(record{typ: recordTypeHandshake, sequence: records[0].sequence, payload: payload})
	if err != nil {
		return 0, err
	}
	if _, err = c.Conn.Write(wire); err != nil {
		return 0, err
	}
	return len(p), nil
}

func TestECHRejectsFinalConfirmationDowngradeAfterHRR(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	list, key := testECHConfig(t, "public.test", 16)
	left, right := memoryDatagramPair()
	capture := &captureWritesConn{Conn: left}
	client := Client(capture, &Config{
		RootCAs: roots, ServerName: "server.test", EncryptedClientHelloConfigList: list,
		SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	server := Server(&corruptFinalECHConfirmationConn{Conn: right}, &Config{
		Certificates: []tls.Certificate{certificate}, EncryptedClientHelloKeys: []EncryptedClientHelloKey{key},
		CurvePreferences: []tls.CurveID{tls.CurveP256}, SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	defer left.Close()
	defer right.Close()
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	var local *localAlertError
	if !errors.As(clientErr, &local) || local.description != alertIllegalParameter {
		t.Fatalf("client downgrade error = %v", clientErr)
	}
	capture.mu.Lock()
	writes := append([][]byte(nil), capture.writes...)
	capture.mu.Unlock()
	sentAlert := false
	for _, datagram := range writes {
		records, parseErr := parsePlainRecords(datagram)
		if parseErr == nil {
			for _, record := range records {
				sentAlert = sentAlert || record.typ == recordTypeAlert && bytes.Equal(record.payload, []byte{alertLevelFatal, alertIllegalParameter})
			}
		}
	}
	if !sentAlert {
		t.Fatal("client did not send fatal illegal_parameter after ECH downgrade")
	}
	_ = left.Close()
	_ = right.Close()
	<-serverDone
}

func TestECHRejectionAuthenticatesRetryAndSuppressesClientCertificate(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	clientList, _ := testECHConfig(t, "server.test", 12)
	retryList, retryKey := testECHConfig(t, "server.test", 13)
	left, right := memoryDatagramPair()
	verifyPeerCalled := false
	client := Client(left, &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate},
		EncryptedClientHelloConfigList: clientList, SessionTicketsDisabled: true,
		VerifyPeerCertificate: func([][]byte, [][]*x509.Certificate) error {
			verifyPeerCalled = true
			return nil
		},
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	server := Server(right, &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequestClientCert, ClientCAs: clientRoots,
		EncryptedClientHelloKeys: []EncryptedClientHelloKey{retryKey}, SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	defer left.Close()
	defer right.Close()
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	serverErr := <-serverDone
	var rejection *ECHRejectionError
	if !errors.As(clientErr, &rejection) {
		t.Fatalf("client error = %v, want ECHRejectionError", clientErr)
	}
	if !bytes.Equal(rejection.RetryConfigList, retryList) {
		t.Fatalf("retry configs = %x, want %x", rejection.RetryConfigList, retryList)
	}
	if serverErr != nil {
		t.Fatalf("server rejection handshake: %v", serverErr)
	}
	if server.ConnectionState().ECHAccepted || len(server.ConnectionState().PeerCertificates) != 0 {
		t.Fatalf("server exposed rejected ECH/client identity: %#v", server.ConnectionState())
	}
	if verifyPeerCalled {
		t.Fatal("ordinary VerifyPeerCertificate ran for an ECH rejection")
	}
}

func TestECHGREASEPreservedAcrossHRR(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	capture := &captureWritesConn{Conn: left}
	client := Client(capture, &Config{
		RootCAs: roots, ServerName: "server.test", EncryptedClientHelloGrease: true,
		SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	server := Server(right, &Config{
		Certificates: []tls.Certificate{certificate}, CurvePreferences: []tls.CurveID{tls.CurveP256},
		SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	defer left.Close()
	defer right.Close()
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	serverErr := <-serverDone
	if clientErr != nil || serverErr != nil {
		t.Fatalf("GREASE handshake failed: client=%v server=%v", clientErr, serverErr)
	}
	capture.mu.Lock()
	writes := append([][]byte(nil), capture.writes...)
	capture.mu.Unlock()
	hellos := make(map[uint16]*clientHello)
	for _, datagram := range writes {
		if len(datagram) == 0 || datagram[0] != recordTypeHandshake {
			continue
		}
		records, err := parsePlainRecords(datagram)
		if err != nil {
			t.Fatal(err)
		}
		for _, record := range records {
			fragments, err := parseHandshakeFragments(record.payload)
			if err != nil {
				t.Fatal(err)
			}
			for _, fragment := range fragments {
				if fragment.typ != handshakeTypeClientHello || fragment.offset != 0 || int(fragment.length) != len(fragment.body) {
					continue
				}
				hello, err := parseClientHello(fragment.body)
				if err != nil {
					t.Fatal(err)
				}
				hellos[fragment.messageSequence] = hello
			}
		}
	}
	if hellos[0] == nil || hellos[1] == nil || !bytes.Equal(hellos[0].encryptedClientHello(), hellos[1].encryptedClientHello()) {
		t.Fatalf("GREASE ECH changed across HRR: first=%x second=%x", echExtension(hellos[0]), echExtension(hellos[1]))
	}
	if client.ConnectionState().ECHAccepted || server.ConnectionState().ECHAccepted {
		t.Fatal("GREASE ECH was reported as accepted")
	}
}

func echExtension(hello *clientHello) []byte {
	if hello == nil {
		return nil
	}
	return hello.encryptedClientHello()
}

func TestECHResumptionAndEarlyData(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	list, key := testECHConfig(t, "public.test", 14)
	cache := NewLRUClientSessionCache(2)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x8e}, len(ticketKey)))
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", EncryptedClientHelloConfigList: list,
		ClientSessionCache: cache, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, EncryptedClientHelloKeys: []EncryptedClientHelloKey{key},
		SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour, MaxEarlyData: 1024,
		AllowEarlyDataWithoutCookie: true, EarlyDataReplayCache: NewLRUEarlyDataReplayCache(4),
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	_ = issueEarlyDataTicket(t, clientConfig, serverConfig)
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	client := Client(left, clientConfig)
	server := Server(right, serverConfig)
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	payload := []byte("ECH early data")
	if n, err := client.WriteEarlyData(payload); err != nil || n != len(payload) {
		t.Fatalf("WriteEarlyData = %d, %v", n, err)
	}
	if err := <-serverDone; err != nil {
		t.Fatal(err)
	}
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume ||
		!client.ConnectionState().ECHAccepted || !server.ConnectionState().ECHAccepted {
		t.Fatalf("ECH resumption state: client=%#v server=%#v", client.ConnectionState(), server.ConnectionState())
	}
	buffer := make([]byte, len(payload))
	n, _, err := server.ReadDatagram(buffer)
	if err != nil || !bytes.Equal(buffer[:n], payload) {
		t.Fatalf("early data = %q, %v", buffer[:n], err)
	}
}

func TestECHRealUDP(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	list, key := testECHConfig(t, "public.test", 15)
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, EncryptedClientHelloKeys: []EncryptedClientHelloKey{key},
		SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second,
		FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
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
		if handshakeErr := server.Handshake(); handshakeErr != nil {
			serverDone <- handshakeErr
			return
		}
		if !server.ConnectionState().ECHAccepted {
			serverDone <- errors.New("server did not accept ECH")
			return
		}
		buffer := make([]byte, 16)
		n, _, readErr := server.ReadDatagram(buffer)
		if readErr == nil && string(buffer[:n]) != "ping" {
			readErr = errors.New("unexpected ECH UDP payload")
		}
		if readErr == nil {
			_, readErr = server.WriteDatagram([]byte("pong"))
		}
		serverDone <- readErr
	}()
	client, err := DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "udp4", listener.Addr().String(), &Config{
		RootCAs: roots, ServerName: "server.test", EncryptedClientHelloConfigList: list,
		SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second,
		FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	if !client.ConnectionState().ECHAccepted {
		t.Fatal("client did not confirm ECH")
	}
	_ = client.SetDeadline(time.Now().Add(2 * time.Second))
	if _, err = client.WriteDatagram([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 16)
	n, _, err := client.ReadDatagram(buffer)
	if err != nil || string(buffer[:n]) != "pong" {
		t.Fatalf("ECH UDP response = %q, %v", buffer[:n], err)
	}
	if err = <-serverDone; err != nil {
		t.Fatal(err)
	}
}
