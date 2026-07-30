package dtls13

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"encoding/hex"
	"slices"
	"testing"
)

func TestHandshakeMarshalRFCWireVectors(t *testing.T) {
	decode := func(encoded string) []byte {
		t.Helper()
		wire, err := hex.DecodeString(encoded)
		if err != nil {
			t.Fatal(err)
		}
		return wire
	}
	tests := []struct {
		name    string
		marshal func() ([]byte, error)
		want    string
	}{
		{
			name: "ClientHello",
			marshal: (&clientHello{
				cipherSuites:     []uint16{TLS_AES_128_GCM_SHA256},
				supportedGroups:  []tls.CurveID{tls.X25519},
				signatureSchemes: []tls.SignatureScheme{tls.Ed25519},
				keyShares:        []keyShareEntry{{group: tls.X25519, data: []byte{0xaa}}},
			}).marshal,
			want: "fefd000000000000000000000000000000000000000000000000000000000000000000000002130101000022000a00040002001d000d000400020807002b000302fefc003300070005001d0001aa",
		},
		{
			name:    "ServerHello",
			marshal: (&serverHello{cipherSuite: TLS_AES_128_GCM_SHA256, keyShare: keyShareEntry{group: tls.X25519, data: []byte{0xaa}}}).marshal,
			want:    "fefd000000000000000000000000000000000000000000000000000000000000000000130100000f002b0002fefc00330005001d0001aa",
		},
		{
			name:    "HelloRetryRequest",
			marshal: (&helloRetryRequest{cipherSuite: TLS_AES_128_GCM_SHA256, selectedGroup: tls.X25519}).marshal,
			want:    "fefdcf21ad74e59a6111be1d8c021e65b891c2a211167abb8c5e079e09e2c8a8339c00130100000c002b0002fefc00330002001d",
		},
		{
			name: "Certificate",
			marshal: (&certificateMessage{requestContext: []byte{0xaa}, certificates: []certificateEntry{{
				data: []byte{0x01, 0x02}, extensions: map[uint16][]byte{0x1234: {0x05}},
			}}}).marshal,
			want: "01aa00000c000002010200051234000105",
		},
		{
			name:    "CertificateVerify",
			marshal: (&certificateVerifyMessage{algorithm: tls.Ed25519, signature: []byte{0x01, 0x02, 0x03}}).marshal,
			want:    "08070003010203",
		},
		{
			name: "NewSessionTicket",
			marshal: (&newSessionTicketMessage{
				lifetime: 1, ageAdd: 0x01020304, nonce: []byte{0xaa, 0xbb}, ticket: []byte{0xcc, 0xdd}, maxEarlyData: 4096,
			}).marshal,
			want: "000000010102030402aabb0002ccdd0008002a000400001000",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.marshal()
			if err != nil {
				t.Fatal(err)
			}
			if want := decode(test.want); !bytes.Equal(got, want) {
				t.Fatalf("wire = %x, want %x", got, want)
			}
		})
	}
}

func TestClientHelloRoundTrip(t *testing.T) {
	h := &clientHello{sessionID: []byte{1, 2, 3}, cookie: []byte{4, 5}, cipherSuites: []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384}, keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{7}, 32)}}, serverName: "example.test", alpn: []string{"h3", "coap"}, connectionID: []byte{1}, hasConnectionID: true, returnRoutability: true}
	for i := range h.random {
		h.random[i] = byte(i)
	}
	b, err := h.marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseClientHello(b)
	if err != nil {
		t.Fatal(err)
	}
	if got.random != h.random || !bytes.Equal(got.sessionID, h.sessionID) || !bytes.Equal(got.cookie, h.cookie) || len(got.cipherSuites) != 2 || len(got.keyShares) != 1 || got.keyShares[0].group != tls.X25519 || !bytes.Equal(got.keyShares[0].data, h.keyShares[0].data) || got.serverName != "example.test" || len(got.alpn) != 2 || got.alpn[1] != "coap" || !got.returnRoutability {
		t.Fatalf("round trip mismatch: %#v", got)
	}
}

func TestClientHelloPostHandshakeAuthExtension(t *testing.T) {
	h := &clientHello{cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, supportedGroups: []tls.CurveID{tls.X25519}, keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}, postHandshakeAuth: true}
	wire, err := h.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseClientHello(wire)
	if err != nil {
		t.Fatal(err)
	}
	if !parsed.postHandshakeAuth {
		t.Fatal("post_handshake_auth was not preserved")
	}
}

func TestClientHelloRejectsUnadvertisedKeyShare(t *testing.T) {
	h := &clientHello{cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, supportedGroups: []tls.CurveID{tls.CurveP256}, keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}}
	b, err := h.marshal()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = parseClientHello(b); err == nil {
		t.Fatal("accepted unadvertised key share group")
	}
}

func TestClientHelloRejectsLegacyCookie(t *testing.T) {
	valid := &clientHello{cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, supportedGroups: []tls.CurveID{tls.X25519}, keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}}
	wire, err := valid.marshal()
	if err != nil {
		t.Fatal(err)
	}
	legacy := make([]byte, 0, len(wire)+1)
	legacy = append(legacy, wire[:35]...)
	legacy = append(legacy, 1, 0x42)
	legacy = append(legacy, wire[36:]...)
	_, err = parseClientHello(legacy)
	if err == nil {
		t.Fatal("parsed nonempty legacy_cookie")
	}
	if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("legacy_cookie error classified as alert %d, ok=%v: %v", description, ok, err)
	}
}

func TestClientHelloRejectsEarlyDataWithoutPSK(t *testing.T) {
	h := &clientHello{cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, supportedGroups: []tls.CurveID{tls.X25519}, keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}, earlyData: true}
	wire, err := h.marshal()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = parseClientHello(wire); err == nil {
		t.Fatal("accepted early_data without pre_shared_key")
	}
}

func TestServerHelloRoundTrip(t *testing.T) {
	h := &serverHello{cipherSuite: TLS_AES_128_GCM_SHA256, keyShare: keyShareEntry{group: tls.X25519, data: bytes.Repeat([]byte{9}, 32)}, connectionID: []byte{7, 8}, hasConnectionID: true, returnRoutability: true}
	h.random[0] = 42
	b, err := h.marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseServerHello(b)
	if err != nil {
		t.Fatal(err)
	}
	if got.random != h.random || got.cipherSuite != h.cipherSuite || got.keyShare.group != h.keyShare.group || !bytes.Equal(got.keyShare.data, h.keyShare.data) || !got.hasConnectionID || !bytes.Equal(got.connectionID, h.connectionID) || !got.returnRoutability {
		t.Fatalf("round trip mismatch: %#v", got)
	}
}

func TestServerHelloAllowsPreSharedKeyInAnyExtensionPosition(t *testing.T) {
	shares, err := marshalKeyShares([]keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{9}, 32)}}, false)
	if err != nil {
		t.Fatal(err)
	}
	exts, err := marshalExtensions(map[uint16][]byte{
		extSupportedVersions: {byte(VersionDTLS13 >> 8), byte(VersionDTLS13 & 0xff)},
		extPreSharedKey:      {0, 0},
		extKeyShare:          shares,
	}, []uint16{extSupportedVersions, extPreSharedKey, extKeyShare})
	if err != nil {
		t.Fatal(err)
	}
	var w wireBuilder
	w.u16(int(dtlsLegacyVersion))
	w.b = append(w.b, make([]byte, 32)...)
	w.bytes8(nil)
	w.u16(int(TLS_AES_128_GCM_SHA256))
	w.u8(0)
	w.b = append(w.b, exts...)
	parsed, err := parseServerHello(w.b)
	if err != nil || parsed.selectedIdentity == nil || *parsed.selectedIdentity != 0 {
		t.Fatalf("parsed=%#v err=%v", parsed, err)
	}
}

func TestServerHelloRejectsLegacySessionIDAndUnknownExtension(t *testing.T) {
	withSessionID := &serverHello{sessionID: []byte{1}, cipherSuite: TLS_AES_128_GCM_SHA256, keyShare: keyShareEntry{group: tls.X25519, data: bytes.Repeat([]byte{9}, 32)}}
	if _, err := withSessionID.marshal(); err == nil {
		t.Fatal("marshaled a DTLS 1.3 ServerHello with legacy_session_id")
	}
	valid := &serverHello{cipherSuite: TLS_AES_128_GCM_SHA256, keyShare: keyShareEntry{group: tls.X25519, data: bytes.Repeat([]byte{9}, 32)}}
	wire, err := valid.marshal()
	if err != nil {
		t.Fatal(err)
	}
	// Insert one legacy_session_id byte immediately after its length.
	legacy := make([]byte, 0, len(wire)+1)
	legacy = append(legacy, wire[:35]...)
	legacy[34] = 1
	legacy = append(legacy, 0x42)
	legacy = append(legacy, wire[35:]...)
	if _, err = parseServerHello(legacy); err == nil {
		t.Fatal("accepted a DTLS 1.3 ServerHello with legacy_session_id")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("legacy session ID alert=%d ok=%v", description, ok)
	}

	// The extensions vector starts at offset 38 for an empty session ID.
	unknown := append([]byte(nil), wire...)
	extensionLength := int(binary.BigEndian.Uint16(unknown[38:40]))
	binary.BigEndian.PutUint16(unknown[38:40], uint16(extensionLength+4))
	unknown = append(unknown, 0xff, 0xff, 0, 0)
	if _, err = parseServerHello(unknown); err == nil {
		t.Fatal("accepted an unknown ServerHello extension")
	} else if description, ok := protocolAlert(err); !ok || description != alertUnsupportedExtension {
		t.Fatalf("unknown extension alert=%d ok=%v", description, ok)
	}
}

func TestServerHelloPreservesEmptyConnectionID(t *testing.T) {
	h := &serverHello{cipherSuite: TLS_AES_128_GCM_SHA256, keyShare: keyShareEntry{group: tls.X25519, data: bytes.Repeat([]byte{9}, 32)}, hasConnectionID: true}
	b, err := h.marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseServerHello(b)
	if err != nil || !got.hasConnectionID || len(got.connectionID) != 0 {
		t.Fatalf("empty connection ID was not preserved: %#v, %v", got, err)
	}
}

func TestServerHelloRejectsUnsolicitedConnectionIDAsUnsupportedExtension(t *testing.T) {
	err := validateServerHelloConnectionID(&clientHello{}, &serverHello{hasConnectionID: true})
	if description, ok := protocolAlert(err); !ok || description != alertUnsupportedExtension {
		t.Fatalf("unsolicited connection_id alert=%d ok=%v err=%v", description, ok, err)
	}
	if err = validateServerHelloConnectionID(&clientHello{hasConnectionID: true}, &serverHello{hasConnectionID: true}); err != nil {
		t.Fatalf("requested connection_id was rejected: %v", err)
	}
}

func TestReturnRoutabilityHelloRequirements(t *testing.T) {
	if err := validateServerHelloReturnRoutability(&clientHello{}, &serverHello{returnRoutability: true}); err == nil {
		t.Fatal("accepted an unsolicited rrc extension")
	} else if description, ok := protocolAlert(err); !ok || description != alertUnsupportedExtension {
		t.Fatalf("unsolicited rrc alert=%d ok=%v", description, ok)
	}
	offer := &clientHello{returnRoutability: true, hasConnectionID: true}
	if err := validateServerHelloReturnRoutability(offer, &serverHello{returnRoutability: true}); err == nil {
		t.Fatal("accepted rrc without a negotiated connection ID")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("rrc without CID alert=%d ok=%v", description, ok)
	}
	if err := validateServerHelloReturnRoutability(offer, &serverHello{returnRoutability: true, hasConnectionID: true}); err != nil {
		t.Fatal(err)
	}

	client := &clientHello{cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, signatureSchemes: defaultSignatureSchemes(), supportedGroups: []tls.CurveID{tls.X25519}, keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}}
	clientWire, err := client.marshal()
	if err != nil {
		t.Fatal(err)
	}
	p := wireParser{b: clientWire}
	_ = p.u16()
	_ = p.take(32)
	_ = p.bytes8()
	_ = p.bytes8()
	_ = p.bytes16()
	_ = p.bytes8()
	extensionLengthOffset := p.off
	extensionLength := int(binary.BigEndian.Uint16(clientWire[extensionLengthOffset:]))
	clientWire = append(clientWire, byte(extReturnRoutability>>8), byte(extReturnRoutability), 0, 1, 1)
	binary.BigEndian.PutUint16(clientWire[extensionLengthOffset:], uint16(extensionLength+5))
	if _, err = parseClientHello(clientWire); err == nil {
		t.Fatal("accepted a non-empty ClientHello rrc extension")
	} else if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
		t.Fatalf("malformed ClientHello rrc alert=%d ok=%v", description, ok)
	}

	server := &serverHello{cipherSuite: TLS_AES_128_GCM_SHA256, keyShare: keyShareEntry{group: tls.X25519, data: bytes.Repeat([]byte{9}, 32)}}
	serverWire, err := server.marshal()
	if err != nil {
		t.Fatal(err)
	}
	serverExtensionLength := int(binary.BigEndian.Uint16(serverWire[38:40]))
	serverWire = append(serverWire, byte(extReturnRoutability>>8), byte(extReturnRoutability), 0, 1, 1)
	binary.BigEndian.PutUint16(serverWire[38:40], uint16(serverExtensionLength+5))
	if _, err = parseServerHello(serverWire); err == nil {
		t.Fatal("accepted a non-empty ServerHello rrc extension")
	} else if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
		t.Fatalf("malformed ServerHello rrc alert=%d ok=%v", description, ok)
	}
}

func TestHelloRejectsTruncation(t *testing.T) {
	h := &clientHello{cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}}
	b, err := h.marshal()
	if err != nil {
		t.Fatal(err)
	}
	for i := range b {
		if _, err := parseClientHello(b[:i]); err == nil {
			t.Fatalf("accepted truncation at %d", i)
		}
	}
}

func TestExtensionsRejectDuplicate(t *testing.T) {
	// vector length 8, followed by two empty extensions with type 43.
	duplicate := []byte{0, 8, 0, 43, 0, 0, 0, 43, 0, 0}
	for _, parse := range []func([]byte) (map[uint16][]byte, error){parseExtensions, parseExtensionsView} {
		if _, err := parse(duplicate); err == nil {
			t.Fatal("accepted duplicate extension")
		} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
			t.Fatalf("duplicate extension alert=%d ok=%v err=%v", description, ok, err)
		}
	}
	var storage [2]orderedExtension
	if _, err := parseOrderedExtensionsView(duplicate, storage[:0]); err == nil {
		t.Fatal("ordered parser accepted duplicate extension")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("ordered parser duplicate extension alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestExtensionParserOwnershipModes(t *testing.T) {
	wire, err := marshalExtensions(map[uint16][]byte{43: {1, 2, 3}}, []uint16{43})
	if err != nil {
		t.Fatal(err)
	}
	owned, err := parseExtensions(wire)
	if err != nil {
		t.Fatal(err)
	}
	view, err := parseExtensionsView(wire)
	if err != nil {
		t.Fatal(err)
	}
	var storage [1]orderedExtension
	ordered, err := parseOrderedExtensionsView(wire, storage[:0])
	if err != nil {
		t.Fatal(err)
	}
	orderedValue, ok := orderedExtensionValue(ordered, 43)
	if !ok || &owned[43][0] == &wire[len(wire)-3] || &view[43][0] != &wire[len(wire)-3] || &orderedValue[0] != &wire[len(wire)-3] {
		t.Fatal("extension parser ownership does not match its mode")
	}
	wire[len(wire)-1] = 9
	if owned[43][2] != 3 || view[43][2] != 9 || orderedValue[2] != 9 {
		t.Fatalf("owned=%v view=%v ordered=%v", owned[43], view[43], orderedValue)
	}
}

func TestOrderedExtensionParserInlineAndOverflow(t *testing.T) {
	items := make(map[uint16][]byte, 12)
	order := make([]uint16, 12)
	for i := range order {
		order[i] = uint16(0xff00 + i)
		items[order[i]] = []byte{byte(i)}
	}
	wire, err := marshalExtensions(items, order)
	if err != nil {
		t.Fatal(err)
	}
	var storage [2]orderedExtension
	parsed, err := parseOrderedExtensionsView(wire, storage[:0])
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != len(order) || cap(parsed) <= len(storage) {
		t.Fatalf("parsed len=%d cap=%d", len(parsed), cap(parsed))
	}
	for i, extension := range parsed {
		if extension.typ != order[i] || len(extension.value) != 1 || extension.value[0] != byte(i) {
			t.Fatalf("extension %d = {%04x %x}", i, extension.typ, extension.value)
		}
	}
	wire[len(wire)-1] = 0xaa
	if parsed[len(parsed)-1].value[0] != 0xaa {
		t.Fatal("overflow parser did not retain value view")
	}
}

func TestOfferedVersionsAcceptsListContainingDTLS13(t *testing.T) {
	versions := []byte{4, 0xfe, 0xfd, byte(VersionDTLS13 >> 8), byte(VersionDTLS13 & 0xff)}
	if !offeredDTLS13(versions) {
		t.Fatal("rejected version list containing DTLS 1.3")
	}
	if offeredDTLS13([]byte{2, 0xfe, 0xfd}) {
		t.Fatal("accepted list without DTLS 1.3")
	}
	if offeredDTLS13([]byte{3, 0xfe, 0xfc}) {
		t.Fatal("accepted malformed version vector")
	}
}

func TestSupportedVersionsAlertClassification(t *testing.T) {
	if _, err := parseClientSupportedVersions([]byte{3, 0xfe, 0xfc}); err == nil {
		t.Fatal("accepted malformed ClientHello supported_versions")
	} else if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
		t.Fatalf("malformed versions alert=%d ok=%v err=%v", description, ok, err)
	}
	hello := &clientHello{
		cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, supportedGroups: []tls.CurveID{tls.X25519},
		keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}, signatureSchemes: defaultSignatureSchemes(),
	}
	body, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	body = replaceClientHelloExtension(t, body, extSupportedVersions, []byte{2, 0xfe, 0xfd})
	_, err = parseClientHello(body)
	if description, ok := protocolAlert(err); !ok || description != alertProtocolVersion {
		t.Fatalf("unavailable version alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestServerHelloAndHRRKeyShareAlertClassification(t *testing.T) {
	extensions, err := marshalExtensions(map[uint16][]byte{
		extSupportedVersions: {byte(VersionDTLS13 >> 8), byte(VersionDTLS13 & 0xff)},
	}, []uint16{extSupportedVersions})
	if err != nil {
		t.Fatal(err)
	}
	var sh wireBuilder
	sh.u16(int(dtlsLegacyVersion))
	sh.b = append(sh.b, make([]byte, 32)...)
	sh.bytes8(nil)
	sh.u16(int(TLS_AES_128_GCM_SHA256))
	sh.u8(0)
	sh.b = append(sh.b, extensions...)
	if _, err = parseServerHello(sh.b); err == nil {
		t.Fatal("accepted ServerHello without key_share")
	} else if description, ok := protocolAlert(err); !ok || description != alertMissingExtension {
		t.Fatalf("ServerHello missing key_share alert=%d ok=%v err=%v", description, ok, err)
	}

	extensions, err = marshalExtensions(map[uint16][]byte{
		extSupportedVersions: {byte(VersionDTLS13 >> 8), byte(VersionDTLS13 & 0xff)},
		extKeyShare:          {0},
	}, []uint16{extSupportedVersions, extKeyShare})
	if err != nil {
		t.Fatal(err)
	}
	var hrr wireBuilder
	hrr.u16(int(dtlsLegacyVersion))
	hrr.b = append(hrr.b, helloRetryRequestRandom[:]...)
	hrr.bytes8(nil)
	hrr.u16(int(TLS_AES_128_GCM_SHA256))
	hrr.u8(0)
	hrr.b = append(hrr.b, extensions...)
	if _, err = parseHelloRetryRequest(hrr.b); err == nil {
		t.Fatal("accepted malformed HelloRetryRequest key_share")
	} else if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
		t.Fatalf("HelloRetryRequest key_share alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestSignatureSchemesRoundTrip(t *testing.T) {
	want := []tls.SignatureScheme{tls.PSSWithSHA256, tls.Ed25519}
	b, err := marshalSignatureSchemes(want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseSignatureSchemes(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %v", got)
	}
}

func TestNestedExtensionVectorsRejectOversizedEntries(t *testing.T) {
	if _, err := marshalServerName(string(make([]byte, 1<<16))); err == nil {
		t.Fatal("marshaled a 65536-byte server name")
	}
	if _, err := marshalALPN([]string{string(make([]byte, 256))}); err == nil {
		t.Fatal("marshaled a 256-byte ALPN protocol")
	}
	identity := []pskIdentityEntry{{identity: make([]byte, 1<<16)}}
	if _, err := marshalClientPSKs(identity, [][]byte{make([]byte, 32)}); err == nil {
		t.Fatal("marshaled a 65536-byte PSK identity")
	}
	identity[0].identity = []byte{1}
	if _, err := marshalClientPSKs(identity, [][]byte{make([]byte, 256)}); err == nil {
		t.Fatal("marshaled a 256-byte PSK binder")
	}
	if _, err := (&serverHello{keyShare: keyShareEntry{group: tls.X25519, data: make([]byte, 65535)}}).marshal(); err == nil {
		t.Fatal("marshaled an oversized ServerHello extension vector")
	}
}

func TestWireBuilderVectorScopesRollbackOnError(t *testing.T) {
	for _, test := range []struct {
		name  string
		build func(*wireBuilder)
	}{
		{name: "vector16-overflow", build: func(w *wireBuilder) {
			start := w.startVector16()
			w.b = append(w.b, make([]byte, 1<<16)...)
			w.endVector16(start)
		}},
		{name: "vector24-overflow", build: func(w *wireBuilder) {
			start := w.startVector24()
			w.b = append(w.b, make([]byte, 1<<24)...)
			w.endVector24(start)
		}},
		{name: "nested-error", build: func(w *wireBuilder) {
			start := w.startVector16()
			w.bytes8(make([]byte, 256))
			w.endVector16(start)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			w := wireBuilder{b: []byte{0xaa}}
			test.build(&w)
			if w.err == nil || !bytes.Equal(w.b, []byte{0xaa}) {
				t.Fatalf("scope did not roll back: wire=%x err=%v", w.b, w.err)
			}
		})
	}
}

func TestDefaultSignatureSchemesCoverImplementedAlgorithms(t *testing.T) {
	want := []tls.SignatureScheme{
		tls.Ed25519, tls.ECDSAWithP256AndSHA256, tls.ECDSAWithP384AndSHA384,
		tls.ECDSAWithP521AndSHA512, tls.PSSWithSHA256, tls.PSSWithSHA384,
		tls.PSSWithSHA512, tls.PKCS1WithSHA256, tls.PKCS1WithSHA384, tls.PKCS1WithSHA512,
	}
	got := defaultSignatureSchemes()
	second := defaultSignatureSchemes()
	if len(got) == 0 || &got[0] != &second[0] {
		t.Fatal("default signature schemes do not share immutable storage")
	}
	before := append([]tls.SignatureScheme(nil), got...)
	if _, err := marshalSignatureSchemes(got); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(got, before) {
		t.Fatal("signature scheme marshal modified shared defaults")
	}
	for _, scheme := range want {
		found := false
		for _, candidate := range got {
			found = found || candidate == scheme
		}
		if !found {
			t.Fatalf("default signature_algorithms omitted 0x%04x", uint16(scheme))
		}
	}
}

func TestPSKKeyExchangeModesUseImmutableSharedStorage(t *testing.T) {
	first := marshalPSKKeyExchangeModes()
	second := marshalPSKKeyExchangeModes()
	if len(first) != 2 || first[0] != 1 || first[1] != 1 || &first[0] != &second[0] {
		t.Fatalf("PSK key exchange modes = %v", first)
	}
	items := map[uint16][]byte{extPSKKeyExchangeModes: first}
	if _, err := marshalExtensions(items, []uint16{extPSKKeyExchangeModes}); err != nil {
		t.Fatal(err)
	}
	if first[0] != 1 || first[1] != 1 {
		t.Fatal("extension marshal modified shared PSK mode storage")
	}
}

func TestSecondClientHelloAllowsOnlyHRRChanges(t *testing.T) {
	initial := &clientHello{
		cipherSuites:  []uint16{TLS_AES_128_GCM_SHA256},
		keyShares:     []keyShareEntry{{group: tls.X25519, data: []byte{1}}},
		pskIdentities: []pskIdentityEntry{{identity: []byte("first"), obfuscatedAge: 1}, {identity: []byte("second"), obfuscatedAge: 2}},
		pskBinders:    [][]byte{{1}, {2}}, earlyData: true,
		unknownExtensions: map[uint16][]byte{0xfefe: {1}},
	}
	second := *initial
	second.cookie = []byte("cookie")
	second.earlyData = false
	second.pskIdentities = []pskIdentityEntry{{identity: []byte("second"), obfuscatedAge: 99}}
	second.pskBinders = [][]byte{{9}}
	if !equalClientHelloAfterHRR(initial, &second, 0) {
		t.Fatal("rejected permitted PSK age, binder, and identity removal changes")
	}
	second.unknownExtensions = map[uint16][]byte{0xfefe: {2}}
	if equalClientHelloAfterHRR(initial, &second, 0) {
		t.Fatal("accepted a changed unknown extension after HRR")
	}
	second.unknownExtensions = map[uint16][]byte{0xfefe: {1}}
	second.pskIdentities = []pskIdentityEntry{{identity: []byte("added")}}
	if equalClientHelloAfterHRR(initial, &second, 0) {
		t.Fatal("accepted an added PSK identity after HRR")
	}
}

func TestSecondClientHelloRejectsChangedInvariantFields(t *testing.T) {
	initial := &clientHello{
		sessionID:                     []byte{1},
		cipherSuites:                  []uint16{TLS_AES_128_GCM_SHA256},
		keyShares:                     []keyShareEntry{{group: tls.X25519, data: []byte{2}}},
		signatureSchemes:              []tls.SignatureScheme{tls.Ed25519},
		certificateSignatureSchemes:   []tls.SignatureScheme{tls.PSSWithSHA256},
		supportedGroups:               []tls.CurveID{tls.X25519},
		serverName:                    "server.test",
		alpn:                          []string{"coap"},
		pskDHE:                        true,
		connectionID:                  []byte{3},
		hasConnectionID:               true,
		returnRoutability:             true,
		postHandshakeAuth:             true,
		certificateCompressionOffered: true,
		recordSizeLimit:               512,
		hasRecordSizeLimit:            true,
		unknownExtensions:             map[uint16][]byte{0xffa5: {4}},
	}
	initial.random[0] = 5
	for _, test := range []struct {
		name   string
		mutate func(*clientHello)
	}{
		{"random", func(h *clientHello) { h.random[0]++ }},
		{"session-id", func(h *clientHello) { h.sessionID = []byte{9} }},
		{"cipher-suites", func(h *clientHello) { h.cipherSuites = []uint16{TLS_AES_256_GCM_SHA384} }},
		{"key-share", func(h *clientHello) { h.keyShares = []keyShareEntry{{group: tls.X25519, data: []byte{9}}} }},
		{"signature-schemes", func(h *clientHello) { h.signatureSchemes = []tls.SignatureScheme{tls.PSSWithSHA256} }},
		{"certificate-signature-schemes", func(h *clientHello) { h.certificateSignatureSchemes = []tls.SignatureScheme{tls.Ed25519} }},
		{"supported-groups", func(h *clientHello) { h.supportedGroups = []tls.CurveID{tls.CurveP256} }},
		{"server-name", func(h *clientHello) { h.serverName = "other.test" }},
		{"alpn", func(h *clientHello) { h.alpn = []string{"h3"} }},
		{"psk-mode", func(h *clientHello) { h.pskDHE = false }},
		{"connection-id", func(h *clientHello) { h.connectionID = []byte{9} }},
		{"connection-id-presence", func(h *clientHello) { h.hasConnectionID = false }},
		{"return-routability", func(h *clientHello) { h.returnRoutability = false }},
		{"post-handshake-auth", func(h *clientHello) { h.postHandshakeAuth = false }},
		{"certificate-compression", func(h *clientHello) {
			h.certificateCompressionOffered = true
			h.unknownExtensions = map[uint16][]byte{extCompressCertificate: {2, 0, 2}}
		}},
		{"record-size-limit", func(h *clientHello) { h.recordSizeLimit++ }},
		{"record-size-limit-presence", func(h *clientHello) { h.hasRecordSizeLimit = false }},
		{"unknown-extension", func(h *clientHello) { h.unknownExtensions = map[uint16][]byte{0xffa5: {9}} }},
	} {
		t.Run(test.name, func(t *testing.T) {
			second := *initial
			second.cookie = []byte("cookie")
			test.mutate(&second)
			if equalClientHelloAfterHRR(initial, &second, 0) {
				t.Fatal("accepted changed invariant field")
			}
		})
	}
}

func TestClientHelloRecordSizeLimitValidation(t *testing.T) {
	hello := &clientHello{
		cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, signatureSchemes: defaultSignatureSchemes(),
		supportedGroups: []tls.CurveID{tls.X25519},
		keyShares:       []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
		recordSizeLimit: minRecordSizeLimit, hasRecordSizeLimit: true,
	}
	body, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseClientHello(body)
	if err != nil || !parsed.hasRecordSizeLimit || parsed.recordSizeLimit != minRecordSizeLimit {
		t.Fatalf("record_size_limit round trip = %d, %v, %v", parsed.recordSizeLimit, parsed.hasRecordSizeLimit, err)
	}
	for _, test := range []struct {
		name      string
		value     []byte
		wantAlert uint8
	}{
		{"bad length", []byte{64}, alertDecodeError},
		{"below minimum", []byte{0, 63}, alertIllegalParameter},
	} {
		t.Run(test.name, func(t *testing.T) {
			wire := replaceClientHelloExtension(t, body, extRecordSizeLimit, test.value)
			_, parseErr := parseClientHello(wire)
			description, ok := protocolAlert(parseErr)
			if !ok || description != test.wantAlert {
				t.Fatalf("alert=%d ok=%v err=%v", description, ok, parseErr)
			}
		})
	}
	high := replaceClientHelloExtension(t, body, extRecordSizeLimit, []byte{0xff, 0xff})
	parsed, err = parseClientHello(high)
	if err != nil || parsed.recordSizeLimit != 0xffff {
		t.Fatalf("server rejected future client limit: %#v, %v", parsed, err)
	}
	hello.recordSizeLimit = minRecordSizeLimit - 1
	if _, err = hello.marshal(); err == nil {
		t.Fatal("marshaled record_size_limit below 64")
	}
	hello.recordSizeLimit = defaultRecordSizeLimit + 1
	if _, err = hello.marshal(); err == nil {
		t.Fatal("marshaled record_size_limit above the DTLS 1.3 maximum")
	}
}

func TestClientHelloRejectsNonZeroPadding(t *testing.T) {
	hello := &clientHello{
		cipherSuites:    []uint16{TLS_AES_128_GCM_SHA256},
		supportedGroups: []tls.CurveID{tls.X25519},
		keyShares:       []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
	}
	body, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	// Extensions follow the fixed fields and their vectors. Locate the final
	// extension vector structurally so the test is independent of CH contents.
	p := wireParser{b: body}
	_ = p.u16()
	_ = p.take(32)
	_ = p.bytes8()
	_ = p.bytes8()
	_ = p.bytes16()
	_ = p.bytes8()
	extensionOffset := p.off
	existing := append([]byte(nil), p.bytes16()...)
	if err = p.done(); err != nil {
		t.Fatal(err)
	}
	var extension wireBuilder
	extension.u16(int(extPadding))
	extension.bytes16([]byte{0, 1})
	all := append(existing, extension.b...)
	modified := append([]byte(nil), body[:extensionOffset]...)
	var vector wireBuilder
	vector.bytes16(all)
	modified = append(modified, vector.b...)
	if _, err = parseClientHello(modified); err == nil {
		t.Fatal("accepted non-zero ClientHello padding")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("padding alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestPSKClientHelloMayOmitSignatureAlgorithms(t *testing.T) {
	hello := &clientHello{
		cipherSuites:    []uint16{TLS_AES_128_GCM_SHA256},
		supportedGroups: []tls.CurveID{tls.X25519},
		keyShares:       []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
		pskIdentity:     []byte("ticket"), pskBinder: bytes.Repeat([]byte{2}, 32),
	}
	body, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	body = removeClientHelloExtension(t, body, extSignatureAlgorithms)
	parsed, err := parseClientHello(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.signatureSchemes) != 0 || len(parsed.pskIdentity) == 0 {
		t.Fatalf("parsed signatures=%v PSK=%x", parsed.signatureSchemes, parsed.pskIdentity)
	}
	if err = requireCertificateSignatureAlgorithms(parsed, true); err != nil {
		t.Fatal(err)
	}
	err = requireCertificateSignatureAlgorithms(parsed, false)
	if description, ok := protocolAlert(err); !ok || description != alertMissingExtension {
		t.Fatalf("fallback alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestPSKKeyExchangeModesWithoutPSKAndUnsupportedMode(t *testing.T) {
	base := &clientHello{
		cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, supportedGroups: []tls.CurveID{tls.X25519},
		keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}, signatureSchemes: defaultSignatureSchemes(),
	}
	body, err := base.marshal()
	if err != nil {
		t.Fatal(err)
	}
	body = appendClientHelloExtension(t, body, extPSKKeyExchangeModes, []byte{1, 1})
	parsed, err := parseClientHello(body)
	if err != nil || !parsed.pskDHE || len(parsed.pskIdentity) != 0 {
		t.Fatalf("standalone modes parsed=%#v err=%v", parsed, err)
	}

	base.pskIdentity = []byte("ticket")
	base.pskBinder = bytes.Repeat([]byte{2}, 32)
	body, err = base.marshal()
	if err != nil {
		t.Fatal(err)
	}
	body = replaceClientHelloExtension(t, body, extPSKKeyExchangeModes, []byte{1, 0})
	parsed, err = parseClientHello(body)
	if err != nil || parsed.pskDHE || len(parsed.pskIdentity) == 0 {
		t.Fatalf("unsupported mode parsed=%#v err=%v", parsed, err)
	}
	body = removeClientHelloExtension(t, body, extPSKKeyExchangeModes)
	_, err = parseClientHello(body)
	if description, ok := protocolAlert(err); !ok || description != alertMissingExtension {
		t.Fatalf("missing PSK modes alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestClientHelloAllowsEmptyKeyShareAndRequiresExtensionPair(t *testing.T) {
	hello := &clientHello{
		cipherSuites:    []uint16{TLS_AES_128_GCM_SHA256},
		supportedGroups: []tls.CurveID{tls.CurveP256},
	}
	body, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseClientHello(body)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed.keyShares) != 0 || len(parsed.supportedGroups) != 1 {
		t.Fatalf("shares=%v groups=%v", parsed.keyShares, parsed.supportedGroups)
	}
	withoutShare := removeClientHelloExtension(t, body, extKeyShare)
	_, err = parseClientHello(withoutShare)
	if description, ok := protocolAlert(err); !ok || description != alertMissingExtension {
		t.Fatalf("missing key_share alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestKeyShareRejectsDuplicateGroups(t *testing.T) {
	wire, err := marshalKeyShares([]keyShareEntry{
		{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)},
		{group: tls.X25519, data: bytes.Repeat([]byte{2}, 32)},
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, parse := range []func([]byte, bool) ([]keyShareEntry, error){parseKeyShares, parseKeySharesView} {
		_, err = parse(wire, true)
		if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
			t.Fatalf("duplicate key share alert=%d ok=%v err=%v", description, ok, err)
		}
	}
}

func TestKeyShareParserOwnershipModes(t *testing.T) {
	wire, err := marshalKeyShares([]keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{7}, 32)}}, true)
	if err != nil {
		t.Fatal(err)
	}
	owned, err := parseKeyShares(wire, true)
	if err != nil {
		t.Fatal(err)
	}
	view, err := parseKeySharesView(wire, true)
	if err != nil {
		t.Fatal(err)
	}
	var storage [1]keyShareEntry
	viewInto, err := parseKeySharesViewInto(wire, true, storage[:0])
	if err != nil {
		t.Fatal(err)
	}
	wire[len(wire)-1] = 9
	if owned[0].data[len(owned[0].data)-1] != 7 || view[0].data[len(view[0].data)-1] != 9 || viewInto[0].data[len(viewInto[0].data)-1] != 9 {
		t.Fatalf("owned=%x view=%x viewInto=%x", owned[0].data, view[0].data, viewInto[0].data)
	}
}

func TestClientKeySharesFollowSupportedGroupOrder(t *testing.T) {
	valid := &clientHello{
		cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, signatureSchemes: defaultSignatureSchemes(),
		supportedGroups: []tls.CurveID{tls.CurveP256, tls.CurveP384, tls.X25519},
		keyShares: []keyShareEntry{
			{group: tls.CurveP256, data: bytes.Repeat([]byte{1}, 65)},
			{group: tls.X25519, data: bytes.Repeat([]byte{2}, 32)},
		},
	}
	body, err := valid.marshal()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = parseClientHello(body); err != nil {
		t.Fatalf("rejected ordered non-contiguous key shares: %v", err)
	}
	valid.keyShares[0], valid.keyShares[1] = valid.keyShares[1], valid.keyShares[0]
	body, err = valid.marshal()
	if err != nil {
		t.Fatal(err)
	}
	_, err = parseClientHello(body)
	if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("out-of-order shares alert=%d ok=%v err=%v", description, ok, err)
	}
}

func removeClientHelloExtension(t *testing.T, body []byte, remove uint16) []byte {
	t.Helper()
	p := wireParser{b: body}
	_ = p.u16()
	_ = p.take(32)
	_ = p.bytes8()
	_ = p.bytes8()
	_ = p.bytes16()
	_ = p.bytes8()
	extensionOffset := p.off
	raw := wireParser{b: p.bytes16()}
	if err := p.done(); err != nil {
		t.Fatal(err)
	}
	var kept wireBuilder
	for raw.off < len(raw.b) {
		typ := uint16(raw.u16())
		value := raw.bytes16()
		if raw.err != nil {
			t.Fatal(raw.err)
		}
		if typ != remove {
			kept.u16(int(typ))
			kept.bytes16(value)
		}
	}
	if err := raw.done(); err != nil {
		t.Fatal(err)
	}
	var vector wireBuilder
	vector.bytes16(kept.b)
	return append(append([]byte(nil), body[:extensionOffset]...), vector.b...)
}

func appendClientHelloExtension(t *testing.T, body []byte, typ uint16, value []byte) []byte {
	t.Helper()
	p := wireParser{b: body}
	_ = p.u16()
	_ = p.take(32)
	_ = p.bytes8()
	_ = p.bytes8()
	_ = p.bytes16()
	_ = p.bytes8()
	extensionOffset := p.off
	existing := append([]byte(nil), p.bytes16()...)
	if err := p.done(); err != nil {
		t.Fatal(err)
	}
	var extension wireBuilder
	extension.u16(int(typ))
	extension.bytes16(value)
	var vector wireBuilder
	vector.bytes16(append(existing, extension.b...))
	return append(append([]byte(nil), body[:extensionOffset]...), vector.b...)
}

func replaceClientHelloExtension(t *testing.T, body []byte, typ uint16, value []byte) []byte {
	t.Helper()
	p := wireParser{b: body}
	_ = p.u16()
	_ = p.take(32)
	_ = p.bytes8()
	_ = p.bytes8()
	_ = p.bytes16()
	_ = p.bytes8()
	extensionOffset := p.off
	extensions := wireParser{b: p.bytes16()}
	var rebuilt wireBuilder
	found := false
	for extensions.off < len(extensions.b) {
		candidate := uint16(extensions.u16())
		raw := extensions.bytes16()
		if candidate == typ {
			raw = value
			found = true
		}
		rebuilt.u16(int(candidate))
		rebuilt.bytes16(raw)
	}
	if err := extensions.done(); err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatalf("ClientHello extension %d not found", typ)
	}
	var vector wireBuilder
	vector.bytes16(rebuilt.b)
	return append(append([]byte(nil), body[:extensionOffset]...), vector.b...)
}
