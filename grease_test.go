package dtls13

import (
	"bytes"
	"crypto/tls"
	"testing"
	"time"
)

type fixedByteReader byte

func (r fixedByteReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(r)
	}
	return len(p), nil
}

func TestGREASEExtensionWireAndParsing(t *testing.T) {
	for seed := range 256 {
		want := uint16(0x0a0a + (seed&15)*0x1010)
		if got := greaseValue(byte(seed)); got != want {
			t.Fatalf("greaseValue(%d) = 0x%04x, want 0x%04x", seed, got, want)
		}
	}

	const grease = 0x4a4a
	find := func(t *testing.T, extensions []orderedExtension) int {
		t.Helper()
		for i, extension := range extensions {
			if extension.typ == grease {
				if len(extension.value) != 0 {
					t.Fatalf("GREASE extension has %d bytes", len(extension.value))
				}
				return i
			}
		}
		t.Fatal("GREASE extension is missing")
		return -1
	}

	t.Run("ClientHello", func(t *testing.T) {
		hello := &clientHello{
			random:          [32]byte{4},
			cipherSuites:    []uint16{TLS_AES_128_GCM_SHA256},
			keyShares:       []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
			supportedGroups: []tls.CurveID{tls.X25519}, signatureSchemes: defaultSignatureSchemes(),
			pskIdentity: []byte("ticket"), pskBinder: make([]byte, 32), pskDHE: true,
			grease: true,
		}
		wire, err := hello.marshal()
		if err != nil {
			t.Fatal(err)
		}
		p := wireParser{b: wire}
		_ = p.u16()
		_ = p.take(32)
		_ = p.bytes8()
		_ = p.bytes8()
		_ = p.bytes16()
		_ = p.bytes8()
		extensions, err := parseOrderedExtensionsView(p.take(len(p.b)-p.off), nil)
		if err != nil {
			t.Fatal(err)
		}
		if find(t, extensions) >= len(extensions)-1 || extensions[len(extensions)-1].typ != extPreSharedKey {
			t.Fatal("GREASE displaced the final pre_shared_key extension")
		}
		parsed, err := parseClientHello(wire)
		if err != nil {
			t.Fatal(err)
		}
		if value, ok := parsed.greaseExtension(); !ok || value != grease {
			t.Fatalf("parsed ClientHello GREASE = %x, present=%v", value, ok)
		}
		wire[2] = 0
		parsed, err = parseClientHello(wire)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := parsed.unknownExtensions[grease]; !ok || parsed.grease {
			t.Fatalf("non-derived GREASE was not preserved in the fallback map: %#v", parsed)
		}
	})

	t.Run("CertificateRequest", func(t *testing.T) {
		message := &certificateRequestMessage{signatureSchemes: []tls.SignatureScheme{tls.Ed25519}}
		wire, err := message.marshalWithCertificateCompression(nil, grease)
		if err != nil {
			t.Fatal(err)
		}
		p := wireParser{b: wire}
		_ = p.bytes8()
		extensions, err := parseOrderedExtensionsView(p.take(len(p.b)-p.off), nil)
		if err != nil {
			t.Fatal(err)
		}
		find(t, extensions)
		parsed, err := parseCertificateRequest(wire)
		if err != nil {
			t.Fatal(err)
		}
		if len(parsed.signatureSchemes) != 1 || parsed.signatureSchemes[0] != tls.Ed25519 {
			t.Fatalf("parsed GREASE entered CertificateRequest state: %#v", parsed)
		}
	})

	t.Run("NewSessionTicket", func(t *testing.T) {
		message := &newSessionTicketMessage{lifetime: 60, ageAdd: 1, nonce: []byte{1}, ticket: []byte{2}}
		wire, err := message.marshalWithGREASE(grease)
		if err != nil {
			t.Fatal(err)
		}
		p := wireParser{b: wire}
		_ = p.u32()
		_ = p.u32()
		_ = p.bytes8()
		_ = p.bytes16()
		extensions, err := parseOrderedExtensionsView(p.take(len(p.b)-p.off), nil)
		if err != nil {
			t.Fatal(err)
		}
		find(t, extensions)
		parsed, err := parseNewSessionTicket(wire)
		if err != nil {
			t.Fatal(err)
		}
		if parsed.lifetime != 60 || parsed.ageAdd != 1 {
			t.Fatalf("parsed GREASE entered NewSessionTicket state: %#v", parsed)
		}
	})
}

func TestGREASEPreservedAcrossHRRWithFixedRandom(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	capture := &captureWritesConn{Conn: left}
	completeHybridHandshake(t, capture, right,
		&Config{Rand: fixedByteReader(0x2f), RootCAs: roots, ServerName: "server.test", CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}, EnableGREASE: true, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second},
		&Config{Rand: fixedByteReader(0x37), Certificates: []tls.Certificate{certificate}, CurvePreferences: []tls.CurveID{tls.CurveP256}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second})
	hellos := capturedClientHellos(t, capture)
	if len(hellos) != 2 {
		t.Fatalf("captured %d ClientHellos, want 2", len(hellos))
	}
	want := greaseValue(0x2f)
	for i, hello := range hellos {
		value, ok := hello.greaseExtension()
		if !ok || value != want {
			t.Fatalf("ClientHello %d GREASE = 0x%04x, present=%v, want 0x%04x", i, value, ok, want)
		}
	}
}

func TestGREASEWithECH(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	list, key := testECHConfig(t, "public.test", 17)
	context, err := newECHClientContext(list)
	if err != nil {
		t.Fatal(err)
	}
	inner := testECHClientHello()
	inner.random[0] = 0x12
	inner.grease = true
	inner.unknownExtensions = map[uint16][]byte{0xfefe: {1}}
	outer, err := makeECHOuter(inner, context.config, fixedByteReader(0x37))
	if err != nil {
		t.Fatal(err)
	}
	want := greaseValue(inner.random[0])
	if value, ok := outer.greaseExtension(); !ok || value != want {
		t.Fatalf("ECH outer GREASE = 0x%04x, present=%v, want 0x%04x", value, ok, want)
	}
	if !bytes.Equal(outer.unknownExtensions[0xfefe], []byte{1}) {
		t.Fatalf("ECH outer lost existing unknown extensions: %#v", outer.unknownExtensions)
	}

	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	client, server := completeHybridHandshake(t, left, right,
		&Config{RootCAs: roots, ServerName: "server.test", EncryptedClientHelloConfigList: list, EnableGREASE: true, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second},
		&Config{Certificates: []tls.Certificate{certificate}, EncryptedClientHelloKeys: []EncryptedClientHelloKey{key}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second})
	if !client.ConnectionState().ECHAccepted || !server.ConnectionState().ECHAccepted {
		t.Fatalf("ECH state: client=%#v server=%#v", client.ConnectionState(), server.ConnectionState())
	}
}

func TestGREASEEndToEndMutualTLSAndTicket(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	cache := NewLRUClientSessionCache(1).(*lruSessionCache)
	left, right := memoryDatagramPair()
	client, server := runTicketRequestHandshake(t,
		&Config{RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate}, ClientSessionCache: cache, EnableGREASE: true, HandshakeTimeout: 2 * time.Second},
		&Config{Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots, EnableGREASE: true, HandshakeTimeout: 2 * time.Second},
		left, right)
	waitForTicketCount(t, cache, 1)
	_ = client.Close()
	_ = server.Close()
}
