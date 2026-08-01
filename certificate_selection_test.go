package dtls13

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"math/big"
	"testing"
	"time"
)

var oidExtKeyUsageClientAuth = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 2}
var oidExtKeyUsageCodeSigning = asn1.ObjectIdentifier{1, 3, 6, 1, 5, 5, 7, 3, 3}

func testSelectionCertificate(t testing.TB, serial int64, commonName, dnsName string, usage x509.KeyUsage, extended ...x509.ExtKeyUsage) tls.Certificate {
	t.Helper()
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber:          big.NewInt(serial),
		Subject:               pkix.Name{CommonName: commonName},
		DNSNames:              nil,
		NotBefore:             now.Add(-time.Hour),
		NotAfter:              now.Add(time.Hour),
		KeyUsage:              usage,
		ExtKeyUsage:           extended,
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	if dnsName != "" {
		template.DNSNames = []string{dnsName}
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, public, private)
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: private, Leaf: leaf}
}

func certificatePool(certificates ...tls.Certificate) *x509.CertPool {
	pool := x509.NewCertPool()
	for i := range certificates {
		pool.AddCert(certificates[i].Leaf)
	}
	return pool
}

func completeSelectionHandshake(t *testing.T, clientConfig, serverConfig *Config) (*Conn, *Conn) {
	t.Helper()
	left, right := memoryDatagramPair()
	t.Cleanup(func() {
		_ = left.Close()
		_ = right.Close()
	})
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

func keyUsageFilter(t testing.TB, usage x509.KeyUsage) CertificateOIDFilter {
	t.Helper()
	bitLength := 0
	for value := usage; value != 0; value >>= 1 {
		bitLength++
	}
	bits := make([]byte, (bitLength+7)/8)
	for bit := 0; bit < bitLength; bit++ {
		if usage&(1<<bit) != 0 {
			bits[bit/8] |= 1 << (7 - bit%8)
		}
	}
	values, err := asn1.Marshal(asn1.BitString{Bytes: bits, BitLength: bitLength})
	if err != nil {
		t.Fatal(err)
	}
	return CertificateOIDFilter{OID: append(asn1.ObjectIdentifier(nil), oidExtensionKeyUsage...), Values: values}
}

func extendedKeyUsageFilter(t testing.TB, usages ...asn1.ObjectIdentifier) CertificateOIDFilter {
	t.Helper()
	values, err := asn1.Marshal(usages)
	if err != nil {
		t.Fatal(err)
	}
	return CertificateOIDFilter{OID: append(asn1.ObjectIdentifier(nil), oidExtensionExtendedKeyUsage...), Values: values}
}

func TestCertificateSelectionExtensionsRoundTripAndValidation(t *testing.T) {
	certificate := testSelectionCertificate(t, 1, "client-ca", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	authorities := [][]byte{certificate.Leaf.RawSubject}
	filters := []CertificateOIDFilter{
		keyUsageFilter(t, x509.KeyUsageDigitalSignature),
		extendedKeyUsageFilter(t, oidExtKeyUsageClientAuth),
	}
	request := &certificateRequestMessage{
		requestContext:         []byte{1, 2},
		signatureSchemes:       []tls.SignatureScheme{tls.Ed25519},
		certificateAuthorities: authorities,
		oidFilters:             filters,
	}
	body, err := request.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseCertificateRequest(body)
	if err != nil {
		t.Fatal(err)
	}
	if !equalByteSlices(parsed.certificateAuthorities, authorities) || len(parsed.oidFilters) != 2 || !parsed.oidFilters[1].OID.Equal(oidExtensionExtendedKeyUsage) {
		t.Fatalf("parsed request = %#v", parsed)
	}

	hello := &clientHello{
		cipherSuites:     []uint16{TLS_AES_128_GCM_SHA256},
		keyShares:        []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
		supportedGroups:  []tls.CurveID{tls.X25519},
		signatureSchemes: []tls.SignatureScheme{tls.Ed25519},
	}
	if err = hello.setCertificateAuthorities(authorities); err != nil {
		t.Fatal(err)
	}
	body, err = hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsedHello, err := parseClientHello(body)
	if err != nil || !equalByteSlices(parsedHello.certificateAuthorityNames(), authorities) {
		t.Fatalf("ClientHello authorities=%x err=%v", parsedHello.certificateAuthorityNames(), err)
	}
	if _, err = parseClientHello(appendClientHelloExtension(t, body, extOIDFilters, []byte{0, 0})); err == nil {
		t.Fatal("accepted oid_filters in ClientHello")
	}

	for _, malformed := range [][]byte{{0, 0}, {0, 2, 0, 0}, {0, 3, 0, 1, 0}} {
		if _, err = parseCertificateAuthorities(malformed); err == nil {
			t.Fatalf("accepted malformed certificate_authorities %x", malformed)
		}
	}
	request.certificateAuthorities = [][]byte{{1, 2, 3}}
	if _, err = request.marshal(); err == nil {
		t.Fatal("marshaled malformed certificate authority")
	}

	oidDER, err := asn1.Marshal(oidExtensionKeyUsage)
	if err != nil {
		t.Fatal(err)
	}
	duplicate := newWireBuilder(32)
	start := duplicate.startVector16()
	for range 2 {
		duplicate.bytes8(oidDER)
		duplicate.bytes16(filters[0].Values)
	}
	duplicate.endVector16(start)
	if _, err = parseOIDFilters(duplicate.b); err == nil {
		t.Fatal("accepted duplicate oid_filters OID")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("duplicate OID alert=%d ok=%v err=%v", description, ok, err)
	}

	anyUsage, err := asn1.Marshal([]asn1.ObjectIdentifier{oidAnyExtendedKeyUsage})
	if err != nil {
		t.Fatal(err)
	}
	anyFilter := newWireBuilder(32)
	start = anyFilter.startVector16()
	anyOID, _ := asn1.Marshal(oidExtensionExtendedKeyUsage)
	anyFilter.bytes8(anyOID)
	anyFilter.bytes16(anyUsage)
	anyFilter.endVector16(start)
	if _, err = parseOIDFilters(anyFilter.b); err == nil {
		t.Fatal("accepted anyExtendedKeyUsage in oid_filters")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("anyExtendedKeyUsage alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestCertificateRequestInfoMatchesCAKeyUsageAndExtendedKeyUsage(t *testing.T) {
	first := testSelectionCertificate(t, 2, "first-client", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	second := testSelectionCertificate(t, 3, "second-client", "", x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment, x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageCodeSigning)
	unknownValue := []byte{0xff}
	info := &CertificateRequestInfo{
		AcceptableCAs:    [][]byte{second.Leaf.RawSubject},
		SignatureSchemes: []tls.SignatureScheme{tls.Ed25519},
		OIDFilters: []CertificateOIDFilter{
			keyUsageFilter(t, x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment),
			extendedKeyUsageFilter(t, oidExtKeyUsageCodeSigning),
			{OID: asn1.ObjectIdentifier{1, 2, 3, 4}, Values: unknownValue},
		},
		Version: VersionDTLS13,
	}
	if err := info.SupportsCertificate(&first); err == nil {
		t.Fatal("accepted certificate that does not match CA and OID filters")
	}
	if err := info.SupportsCertificate(&second); err != nil {
		t.Fatal(err)
	}
	request := &certificateRequestMessage{signatureSchemes: []tls.SignatureScheme{tls.Ed25519}, oidFilters: info.OIDFilters[2:]}
	body, err := request.marshal()
	if err != nil {
		t.Fatalf("marshal unknown OID filter: %v", err)
	}
	parsed, err := parseCertificateRequest(body)
	if err != nil || len(parsed.oidFilters) != 1 || !bytes.Equal(parsed.oidFilters[0].Values, unknownValue) {
		t.Fatalf("unknown OID filter round trip: parsed=%#v err=%v", parsed, err)
	}
}

func TestDynamicCertificateSelectionEndToEnd(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	wrongClient := testSelectionCertificate(t, 4, "wrong-client", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	rightClient := testSelectionCertificate(t, 5, "right-client", "", x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment, x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageCodeSigning)
	filters := []CertificateOIDFilter{
		keyUsageFilter(t, x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment),
		extendedKeyUsageFilter(t, oidExtKeyUsageCodeSigning),
	}
	_, server := completeSelectionHandshake(t, &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{wrongClient, rightClient},
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}, &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certificatePool(rightClient), ClientCertificateOIDFilters: filters,
		SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if got := server.ConnectionState().PeerCertificates[0].Subject.CommonName; got != "right-client" {
		t.Fatalf("selected client certificate %q", got)
	}

	wrongServer := testSelectionCertificate(t, 6, "wrong-server", "server.test", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageServerAuth)
	rightServer := testSelectionCertificate(t, 7, "right-server", "server.test", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageServerAuth)
	client, _ := completeSelectionHandshake(t, &Config{
		RootCAs: certificatePool(rightServer), ServerName: "server.test",
		ServerCertificateAuthorities: [][]byte{rightServer.Leaf.RawSubject},
		SessionTicketsDisabled:       true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}, &Config{
		Certificates: []tls.Certificate{wrongServer, rightServer}, SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if got := client.ConnectionState().PeerCertificates[0].Subject.CommonName; got != "right-server" {
		t.Fatalf("selected server certificate %q", got)
	}
}

func TestGetClientCertificateAndPostHandshakeSelection(t *testing.T) {
	clientCertificate := testSelectionCertificate(t, 8, "callback-client", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	request := &certificateRequestMessage{signatureSchemes: []tls.SignatureScheme{tls.Ed25519}}
	callbackError := errors.New("callback failure")
	for _, test := range []struct {
		name     string
		callback func(*CertificateRequestInfo) (*tls.Certificate, error)
		wantNil  bool
		wantErr  error
	}{
		{"certificate", func(info *CertificateRequestInfo) (*tls.Certificate, error) {
			if info.Version != VersionDTLS13 || info.Conn == nil {
				t.Fatalf("callback info = %#v", info)
			}
			return &clientCertificate, nil
		}, false, nil},
		{"empty", func(*CertificateRequestInfo) (*tls.Certificate, error) { return new(tls.Certificate), nil }, false, nil},
		{"nil", func(*CertificateRequestInfo) (*tls.Certificate, error) { return nil, nil }, true, nil},
		{"error", func(*CertificateRequestInfo) (*tls.Certificate, error) { return nil, callbackError }, true, callbackError},
	} {
		t.Run(test.name, func(t *testing.T) {
			conn := &Conn{config: &Config{GetClientCertificate: test.callback}}
			certificate, err := conn.selectClientCertificate(request)
			if test.wantErr != nil && !errors.Is(err, test.wantErr) {
				t.Fatalf("error=%v want=%v", err, test.wantErr)
			}
			if test.name == "nil" && err == nil {
				t.Fatal("accepted nil GetClientCertificate result")
			}
			if test.wantErr == nil && test.name != "nil" && err != nil {
				t.Fatal(err)
			}
			if test.wantNil && certificate != nil {
				t.Fatalf("certificate=%#v", certificate)
			}
		})
	}

	serverCertificate, roots := testServerCertificate(t)
	wrongClient := testSelectionCertificate(t, 9, "wrong-client", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	rightClient := testSelectionCertificate(t, 10, "pha-client", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	calls := 0
	client, server := completeSelectionHandshake(t, &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{wrongClient, rightClient}, PostHandshakeAuth: true,
		GetClientCertificate: func(info *CertificateRequestInfo) (*tls.Certificate, error) {
			calls++
			if !equalByteSlices(info.AcceptableCAs, [][]byte{rightClient.Leaf.RawSubject}) {
				t.Fatalf("acceptable CAs = %x", info.AcceptableCAs)
			}
			return &rightClient, nil
		},
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}, &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certificatePool(rightClient), HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if err := server.RequestClientCertificate(context.Background()); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("GetClientCertificate calls=%d want=2", calls)
	}
	if got := server.ConnectionState().PeerCertificates[0].Subject.CommonName; got != "pha-client" {
		t.Fatalf("post-handshake certificate=%q", got)
	}
	_ = client.Close()
}

func TestResumedMutualTLSDoesNotSelectClientCertificateAgain(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate := testSelectionCertificate(t, 11, "resumed-client", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	cache := NewLRUClientSessionCache(2)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x75}, 32))
	calls := 0
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache,
		GetClientCertificate: func(*CertificateRequestInfo) (*tls.Certificate, error) {
			calls++
			return &clientCertificate, nil
		},
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certificatePool(clientCertificate), SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	_ = issueEarlyDataTicket(t, clientConfig, serverConfig)
	if calls != 1 {
		t.Fatalf("initial GetClientCertificate calls=%d", calls)
	}
	calls = 0
	client, server := completeSelectionHandshake(t, clientConfig, serverConfig)
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume {
		t.Fatal("mutual TLS connection did not resume")
	}
	if calls != 0 {
		t.Fatalf("resumed handshake selected a fresh client certificate %d times", calls)
	}
}

func TestCertificateAuthoritiesHRRAndECHInvariants(t *testing.T) {
	first := testSelectionCertificate(t, 12, "first", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	second := testSelectionCertificate(t, 13, "second", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	initial := new(clientHello)
	if err := initial.setCertificateAuthorities([][]byte{first.Leaf.RawSubject}); err != nil {
		t.Fatal(err)
	}
	retry := cloneClientHello(initial)
	retry.cookie = []byte("cookie")
	if !equalClientHelloAfterHRR(initial, retry, 0) {
		t.Fatal("unchanged certificate_authorities failed HRR invariant")
	}
	if err := retry.setCertificateAuthorities([][]byte{second.Leaf.RawSubject}); err != nil {
		t.Fatal(err)
	}
	if equalClientHelloAfterHRR(initial, retry, 0) {
		t.Fatal("changed certificate_authorities passed HRR invariant")
	}
	outer, err := makeECHOuter(initial, &echConfig{publicName: "public.test"}, zeroReader{})
	if err != nil {
		t.Fatal(err)
	}
	if len(outer.certificateAuthorityNames()) != 0 {
		t.Fatal("ECH outer ClientHello exposed certificate_authorities")
	}
	if !equalByteSlices(cloneClientHello(initial).certificateAuthorityNames(), initial.certificateAuthorityNames()) {
		t.Fatal("ECH clone lost certificate_authorities")
	}
}

func TestCertificateSelectionConfigValidationAndClone(t *testing.T) {
	filter := keyUsageFilter(t, x509.KeyUsageDigitalSignature)
	config := &Config{
		ServerCertificateAuthorities: [][]byte{{0x30, 0x00}},
		ClientCertificateOIDFilters:  []CertificateOIDFilter{filter},
	}
	clone := config.Clone()
	clone.ServerCertificateAuthorities[0][0] = 0
	clone.ClientCertificateOIDFilters[0].OID[0] = 9
	clone.ClientCertificateOIDFilters[0].Values[0] = 0
	if config.ServerCertificateAuthorities[0][0] != 0x30 || config.ClientCertificateOIDFilters[0].OID[0] != 2 || config.ClientCertificateOIDFilters[0].Values[0] == 0 {
		t.Fatal("Config.Clone shared certificate selection slices")
	}

	any := extendedKeyUsageFilter(t, oidAnyExtendedKeyUsage)
	if _, err := (&Config{ClientCertificateOIDFilters: []CertificateOIDFilter{any}}).normalized(); err == nil {
		t.Fatal("accepted anyExtendedKeyUsage in configured oid_filters")
	}
	if _, err := (&Config{ServerCertificateAuthorities: [][]byte{{1, 2, 3}}}).normalized(); err == nil {
		t.Fatal("accepted malformed configured certificate authority")
	}
}
