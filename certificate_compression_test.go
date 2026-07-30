package dtls13

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"testing"
	"time"
)

func testCertificateCompressionAlgorithms(t testing.TB, algorithms ...uint16) *certificateCompressionAlgorithms {
	t.Helper()
	raw := make([]byte, 1+2*len(algorithms))
	raw[0] = byte(2 * len(algorithms))
	for index, algorithm := range algorithms {
		binary.BigEndian.PutUint16(raw[1+2*index:], algorithm)
	}
	offer, err := parseCertificateCompressionAlgorithms(raw)
	if err != nil {
		t.Fatal(err)
	}
	return offer
}

func TestCertificateCompressionAlgorithms(t *testing.T) {
	offer := testCertificateCompressionAlgorithms(t, certificateCompressionZlib, 0xfafa)
	wire, err := marshalCertificateCompressionAlgorithms(offer)
	if err != nil {
		t.Fatal(err)
	}
	algorithms, err := parseCertificateCompressionAlgorithms(wire)
	if err != nil || !supportsCertificateCompression(algorithms, certificateCompressionZlib) || !supportsCertificateCompression(algorithms, 0xfafa) {
		t.Fatalf("algorithms=%x err=%v", algorithms, err)
	}
	for _, malformed := range [][]byte{nil, {0}, {1, 0}, {2, 0}, {3, 0, 1, 0}} {
		if _, err = parseCertificateCompressionAlgorithms(malformed); err == nil {
			t.Fatalf("accepted malformed algorithm vector %x", malformed)
		} else if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
			t.Fatalf("vector %x alert=%d ok=%v err=%v", malformed, description, ok, err)
		}
	}
	if _, err = marshalCertificateCompressionAlgorithms(nil); err == nil {
		t.Fatal("marshaled an empty algorithm vector")
	}
}

func TestCompressedCertificateRoundTripAndFallback(t *testing.T) {
	want := &certificateMessage{requestContext: []byte{1}, certificates: []certificateEntry{{data: bytes.Repeat([]byte("certificate"), 200)}}}
	certificate, err := want.marshal()
	if err != nil {
		t.Fatal(err)
	}
	typ, body, err := certificateHandshakeMessage(certificate, &certificateCompressionZlibOffer, true, nil)
	if err != nil {
		t.Fatal(err)
	}
	if typ != handshakeTypeCompressedCertificate || len(body) >= len(certificate) {
		t.Fatalf("type=%d compressed=%d certificate=%d", typ, len(body), len(certificate))
	}
	got, err := parseCertificateHandshakeMessage(typ, body, testCertificateCompressionAlgorithms(t, 0xfafa, certificateCompressionZlib), 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if !equalBytes(got.requestContext, want.requestContext) || len(got.certificates) != 1 || !equalBytes(got.certificates[0].data, want.certificates[0].data) {
		t.Fatalf("round trip mismatch: %#v", got)
	}

	tiny, err := (&certificateMessage{}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	typ, body, err = certificateHandshakeMessage(tiny, &certificateCompressionZlibOffer, true, nil)
	if err != nil || typ != handshakeTypeCertificate || !equalBytes(body, tiny) {
		t.Fatalf("tiny Certificate type=%d body=%x err=%v", typ, body, err)
	}
	typ, body, err = certificateHandshakeMessage(certificate, testCertificateCompressionAlgorithms(t, 2), true, nil)
	if err != nil || typ != handshakeTypeCertificate || !equalBytes(body, certificate) {
		t.Fatalf("unoffered zlib type=%d err=%v", typ, err)
	}
}

func TestCompressedCertificateRejectsInvalidInput(t *testing.T) {
	certificate, err := (&certificateMessage{certificates: []certificateEntry{{data: bytes.Repeat([]byte{0x42}, 512)}}}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	valid, err := marshalCompressedCertificate(certificate)
	if err != nil {
		t.Fatal(err)
	}
	badCertificate := func(t *testing.T, body []byte, offered *certificateCompressionAlgorithms, maxSize int) {
		t.Helper()
		_, parseErr := parseCertificateHandshakeMessage(handshakeTypeCompressedCertificate, body, offered, maxSize)
		if description, ok := protocolAlert(parseErr); !ok || description != alertBadCertificate {
			t.Fatalf("alert=%d ok=%v err=%v", description, ok, parseErr)
		}
	}

	badCertificate(t, valid, nil, 1<<20)
	unknown := append([]byte(nil), valid...)
	binary.BigEndian.PutUint16(unknown, 2)
	badCertificate(t, unknown, testCertificateCompressionAlgorithms(t, 2), 1<<20)
	tooLarge := append([]byte(nil), valid...)
	tooLarge[2], tooLarge[3], tooLarge[4] = 0, 4, 1
	badCertificate(t, tooLarge, &certificateCompressionZlibOffer, 1024)
	short := append([]byte(nil), valid...)
	declared := len(certificate) + 1
	short[2], short[3], short[4] = byte(declared>>16), byte(declared>>8), byte(declared)
	badCertificate(t, short, &certificateCompressionZlibOffer, 1<<20)
	long := append([]byte(nil), valid...)
	declared = len(certificate) - 1
	long[2], long[3], long[4] = byte(declared>>16), byte(declared>>8), byte(declared)
	badCertificate(t, long, &certificateCompressionZlibOffer, 1<<20)
	corrupt := append([]byte(nil), valid...)
	corrupt[len(corrupt)-1] ^= 1
	badCertificate(t, corrupt, &certificateCompressionZlibOffer, 1<<20)
	trailing := append(append([]byte(nil), valid...), 0)
	compressedLength := len(trailing) - 8
	trailing[5], trailing[6], trailing[7] = byte(compressedLength>>16), byte(compressedLength>>8), byte(compressedLength)
	badCertificate(t, trailing, &certificateCompressionZlibOffer, 1<<20)

	for _, malformed := range [][]byte{nil, make([]byte, 7), make([]byte, 8)} {
		_, parseErr := parseCertificateHandshakeMessage(handshakeTypeCompressedCertificate, malformed, &certificateCompressionZlibOffer, 1<<20)
		if description, ok := protocolAlert(parseErr); !ok || description != alertDecodeError {
			t.Fatalf("body %x alert=%d ok=%v err=%v", malformed, description, ok, parseErr)
		}
	}
	invalidCertificate, err := marshalCompressedCertificate([]byte{0})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = parseCertificateHandshakeMessage(handshakeTypeCompressedCertificate, invalidCertificate, &certificateCompressionZlibOffer, 1<<20); err == nil {
		t.Fatal("accepted malformed Certificate after decompression")
	}
}

func TestCertificateCompressionExtensionRoundTrips(t *testing.T) {
	hello := &clientHello{
		cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, signatureSchemes: defaultSignatureSchemes(),
		supportedGroups: []tls.CurveID{tls.X25519}, keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
		certificateCompressionOffered: true,
		unknownExtensions:             map[uint16][]byte{extCompressCertificate: {4, 0, 1, 0xfa, 0xfa}},
	}
	body, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseClientHello(body)
	if err != nil || !supportsCertificateCompression(parsed.certificateCompressionAlgorithms(), certificateCompressionZlib) || !supportsCertificateCompression(parsed.certificateCompressionAlgorithms(), 0xfafa) {
		t.Fatalf("ClientHello algorithms=%x err=%v", parsed.certificateCompressionAlgorithms(), err)
	}
	request := &certificateRequestMessage{
		signatureSchemes: []tls.SignatureScheme{tls.Ed25519},
	}
	body, err = request.marshalWithCertificateCompression(&certificateCompressionZlibOffer)
	if err != nil {
		t.Fatal(err)
	}
	_, parsedAlgorithms, err := parseCertificateRequestWithCompression(body)
	if err != nil || !supportsCertificateCompression(parsedAlgorithms, certificateCompressionZlib) {
		t.Fatalf("CertificateRequest algorithms=%x err=%v", parsedAlgorithms, err)
	}

	second := *hello
	second.cookie = []byte("cookie")
	if !equalClientHelloAfterHRR(hello, &second, 0) {
		t.Fatal("unchanged certificate compression offer failed HRR invariant")
	}
	second.certificateCompressionOffered = true
	second.unknownExtensions = map[uint16][]byte{extCompressCertificate: {2, 0, 2}}
	if equalClientHelloAfterHRR(hello, &second, 0) {
		t.Fatal("accepted changed certificate compression offer after HRR")
	}
}

func completeCertificateCompressionHandshake(t *testing.T, clientConfig, serverConfig *Config) (*Conn, *Conn) {
	t.Helper()
	left, right := memoryDatagramPair()
	t.Cleanup(func() { _ = left.Close(); _ = right.Close() })
	client := Client(left, clientConfig)
	server := Server(right, serverConfig)
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	serverErr := <-serverDone
	if clientErr != nil || serverErr != nil {
		t.Fatalf("handshake failed: client=%v server=%v", clientErr, serverErr)
	}
	return client, server
}

func compressibleTestCertificate(certificate tls.Certificate) tls.Certificate {
	chain := append([][]byte(nil), certificate.Certificate...)
	for range 3 {
		chain = append(chain, certificate.Certificate[0])
	}
	certificate.Certificate = chain
	return certificate
}

func TestCertificateCompressionEndToEnd(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	certificate = compressibleTestCertificate(certificate)
	message := &certificateMessage{}
	for _, der := range certificate.Certificate {
		message.certificates = append(message.certificates, certificateEntry{data: der})
	}
	certificateBody, err := message.marshal()
	if err != nil {
		t.Fatal(err)
	}
	if typ, _, err := certificateHandshakeMessage(certificateBody, &certificateCompressionZlibOffer, true, nil); err != nil || typ != handshakeTypeCompressedCertificate {
		t.Fatalf("test certificate was not compressed: type=%d err=%v", typ, err)
	}
	for _, test := range []struct {
		name                         string
		clientEnabled, serverEnabled bool
	}{
		{"both", true, true},
		{"client-only", true, false},
		{"server-only", false, true},
		{"disabled", false, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, server := completeCertificateCompressionHandshake(t,
				&Config{RootCAs: roots, ServerName: "server.test", EnableCertificateCompression: test.clientEnabled, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond},
				&Config{Certificates: []tls.Certificate{certificate}, EnableCertificateCompression: test.serverEnabled, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
			if len(client.ConnectionState().PeerCertificates) == 0 || !client.ConnectionState().HandshakeComplete || !server.ConnectionState().HandshakeComplete {
				t.Fatal("certificate compression handshake lost authenticated state")
			}
		})
	}
}

func TestCertificateCompressionSessionResumption(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	certificate = compressibleTestCertificate(certificate)
	cache := NewLRUClientSessionCache(2)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x83}, len(ticketKey)))
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache, EnableCertificateCompression: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour, EnableCertificateCompression: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	firstClient, firstServer := completeCertificateCompressionHandshake(t, clientConfig, serverConfig)
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if _, ok := cache.Get("server.test"); ok {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if _, ok := cache.Get("server.test"); !ok {
		t.Fatal("compressed full handshake did not produce a session ticket")
	}
	_ = firstClient.Close()
	_ = firstServer.Close()
	client, server := completeCertificateCompressionHandshake(t, clientConfig, serverConfig)
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume {
		t.Fatal("certificate compression configuration prevented session resumption")
	}
}

func TestCertificateCompressionMutualTLSAndRecordLimit(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	serverCertificate = compressibleTestCertificate(serverCertificate)
	clientCertificate = compressibleTestCertificate(clientCertificate)
	client, server := completeCertificateCompressionHandshake(t,
		&Config{RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate}, EnableCertificateCompression: true, RecordSizeLimit: 64, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond},
		&Config{Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots, EnableCertificateCompression: true, RecordSizeLimit: 64, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	if len(client.ConnectionState().PeerCertificates) == 0 || len(server.ConnectionState().PeerCertificates) == 0 {
		t.Fatal("compressed mutual TLS lost peer certificates")
	}
}

func TestCertificateCompressionWithHelloRetryRequest(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	certificate = compressibleTestCertificate(certificate)
	completeCertificateCompressionHandshake(t,
		&Config{RootCAs: roots, ServerName: "server.test", EnableCertificateCompression: true, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond},
		&Config{Certificates: []tls.Certificate{certificate}, CurvePreferences: []tls.CurveID{tls.CurveP256}, EnableCertificateCompression: true, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
}

func TestCertificateCompressionPostHandshakeAuthentication(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	serverCertificate = compressibleTestCertificate(serverCertificate)
	clientCertificate = compressibleTestCertificate(clientCertificate)
	client, server := completeCertificateCompressionHandshake(t,
		&Config{RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate}, PostHandshakeAuth: true, EnableCertificateCompression: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond},
		&Config{Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots, EnableCertificateCompression: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.RequestClientCertificate(ctx); err != nil {
		t.Fatal(err)
	}
	if len(server.ConnectionState().PeerCertificates) == 0 {
		t.Fatal("compressed post-handshake authentication lost client certificate")
	}
	_ = client.Close()
	_ = server.Close()
}

func TestCertificateCompressionWeakNetwork(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	serverCertificate = compressibleTestCertificate(serverCertificate)
	clientCertificate = compressibleTestCertificate(clientCertificate)
	left, right := memoryDatagramPair()
	clientWire := &weakNetworkConn{Conn: left, enabled: true}
	serverWire := &weakNetworkConn{Conn: right, enabled: true}
	client := Client(clientWire, &Config{RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate}, EnableCertificateCompression: true, SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(serverWire, &Config{Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots, EnableCertificateCompression: true, SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	serverErr := <-serverDone
	if clientErr != nil || serverErr != nil {
		t.Fatalf("weak-network handshake failed: client=%v server=%v", clientErr, serverErr)
	}
	_ = left.Close()
	_ = right.Close()
}

func TestCompressedCertificateUsesWireTranscript(t *testing.T) {
	certificate, err := (&certificateMessage{certificates: []certificateEntry{{data: bytes.Repeat([]byte{1}, 512)}}}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	typ, compressed, err := certificateHandshakeMessage(certificate, &certificateCompressionZlibOffer, true, nil)
	if err != nil || typ != handshakeTypeCompressedCertificate {
		t.Fatalf("type=%d err=%v", typ, err)
	}
	suite, err := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	if err != nil {
		t.Fatal(err)
	}
	wireTranscript := newTranscriptHash(suite.hash.New())
	plainTranscript := newTranscriptHash(suite.hash.New())
	if err = wireTranscript.add(typ, 3, compressed); err != nil {
		t.Fatal(err)
	}
	if err = plainTranscript.add(handshakeTypeCertificate, 3, certificate); err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(wireTranscript.sum(), plainTranscript.sum()) {
		t.Fatal("CompressedCertificate transcript used decompressed Certificate wire")
	}
}

func BenchmarkCertificateCompression(b *testing.B) {
	certificate, err := (&certificateMessage{certificates: []certificateEntry{{data: bytes.Repeat([]byte("certificate-der"), 256)}}}).marshal()
	if err != nil {
		b.Fatal(err)
	}
	compressed, err := marshalCompressedCertificate(certificate)
	if err != nil {
		b.Fatal(err)
	}
	b.ReportMetric(float64(len(compressed))/float64(len(certificate)), "wire-ratio")
	b.Run("Compress", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := marshalCompressedCertificate(certificate); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("Decompress", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := parseCertificateHandshakeMessage(handshakeTypeCompressedCertificate, compressed, &certificateCompressionZlibOffer, 1<<20); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkCertificateCompressionHandshakeLifecycle(b *testing.B) {
	serverCertificate, roots := testServerCertificate(b)
	clientCertificate, clientRoots := testClientCertificate(b)
	serverCertificate = compressibleTestCertificate(serverCertificate)
	clientCertificate = compressibleTestCertificate(clientCertificate)
	profiles := []struct {
		name   string
		client *Config
		server *Config
	}{
		{
			name:   "ServerCertificate",
			client: &Config{RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true, HandshakeTimeout: time.Second},
			server: &Config{Certificates: []tls.Certificate{serverCertificate}, SessionTicketsDisabled: true, HandshakeTimeout: time.Second},
		},
		{
			name:   "MutualTLS",
			client: &Config{RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate}, SessionTicketsDisabled: true, HandshakeTimeout: time.Second},
			server: &Config{Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots, SessionTicketsDisabled: true, HandshakeTimeout: time.Second},
		},
	}
	for _, profile := range profiles {
		for _, mode := range []struct {
			name    string
			enabled bool
		}{{"Plain", false}, {"Zlib", true}} {
			clientConfig, serverConfig := profile.client.Clone(), profile.server.Clone()
			clientConfig.EnableCertificateCompression = mode.enabled
			serverConfig.EnableCertificateCompression = mode.enabled
			b.Run(profile.name+"/"+mode.name, func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
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
						b.Fatalf("handshake failed: client=%v server=%v", clientErr, serverErr)
					}
				}
			})
		}
	}
}
