package dtls13

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"testing"
	"time"
)

func makeTestRSACertificate(t testing.TB, bits int, algorithm x509.SignatureAlgorithm) (tls.Certificate, *x509.Certificate, *x509.CertPool) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	template := &x509.Certificate{
		SerialNumber: big.NewInt(100), Subject: pkix.Name{CommonName: "example.test"}, DNSNames: []string{"example.test"},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), SignatureAlgorithm: algorithm,
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		IsCA:        true, BasicConstraintsValid: true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(parsed)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key}, parsed, roots
}

func makeTestCertificateWithSHA1Root(t testing.TB) (tls.Certificate, *x509.Certificate, *x509.Certificate, *x509.CertPool) {
	t.Helper()
	rootKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	rootTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(101), Subject: pkix.Name{CommonName: "SHA-1 root"},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), SignatureAlgorithm: x509.SHA1WithRSA,
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		IsCA:     true, BasicConstraintsValid: true,
	}
	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, &rootKey.PublicKey, rootKey)
	if err != nil {
		t.Fatal(err)
	}
	root, err := x509.ParseCertificate(rootDER)
	if err != nil {
		t.Fatal(err)
	}
	leafPublic, leafKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(102), Subject: pkix.Name{CommonName: "example.test"}, DNSNames: []string{"example.test"},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, root, leafPublic, rootKey)
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := x509.ParseCertificate(leafDER)
	if err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AddCert(root)
	return tls.Certificate{Certificate: [][]byte{leafDER, rootDER}, PrivateKey: leafKey}, leaf, root, roots
}

func makeTestCertificate(t *testing.T, dns string, usage x509.ExtKeyUsage) ([]byte, *x509.CertPool) {
	t.Helper()
	pub, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(1), Subject: pkix.Name{CommonName: dns}, DNSNames: []string{dns}, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, ExtKeyUsage: []x509.ExtKeyUsage{usage}, IsCA: true, BasicConstraintsValid: true}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(cert)
	return der, pool
}

func TestVerifyServerCertificate(t *testing.T) {
	der, roots := makeTestCertificate(t, "example.test", x509.ExtKeyUsageServerAuth)
	cfg, err := (&Config{RootCAs: roots, ServerName: "example.test"}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	msg := &certificateMessage{certificates: []certificateEntry{{data: der}}}
	certs, chains, err := verifyCertificateChain(cfg, msg, true, defaultSignatureSchemes())
	if err != nil {
		t.Fatal(err)
	}
	if len(certs) != 1 || len(chains) == 0 {
		t.Fatal("missing verified chain")
	}
	cfg.ServerName = "wrong.test"
	if _, _, err = verifyCertificateChain(cfg, msg, true, defaultSignatureSchemes()); err == nil {
		t.Fatal("accepted wrong DNS name")
	}
}
func TestVerifyCertificateCallbackWithInsecureSkip(t *testing.T) {
	der, _ := makeTestCertificate(t, "example.test", x509.ExtKeyUsageServerAuth)
	called := false
	cfg, err := (&Config{InsecureSkipVerify: true, VerifyPeerCertificate: func(raw [][]byte, chains [][]*x509.Certificate) error {
		called = true
		if len(raw) != 1 || chains != nil {
			t.Fatal("unexpected callback arguments")
		}
		return nil
	}}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = verifyCertificateChain(cfg, &certificateMessage{certificates: []certificateEntry{{data: der}}}, true, defaultSignatureSchemes()); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("callback was not called")
	}
}

func TestVerifyCertificateTLS13SigningRequirements(t *testing.T) {
	cfg, err := (&Config{InsecureSkipVerify: true}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = verifyCertificateChain(cfg, &certificateMessage{}, true, defaultSignatureSchemes())
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertDecodeError {
		t.Fatalf("empty server Certificate returned %v", err)
	}

	pub, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(2), Subject: pkix.Name{CommonName: "no-signing"},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour),
		KeyUsage: x509.KeyUsageKeyEncipherment,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, key)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = verifyCertificateChain(cfg, &certificateMessage{certificates: []certificateEntry{{data: der}}}, true, defaultSignatureSchemes())
	local = nil
	if !errors.As(err, &local) || local.description != alertBadCertificate {
		t.Fatalf("non-signing server certificate returned %v", err)
	}
}

func TestSignatureAlgorithmsCertIsApplied(t *testing.T) {
	hello := &clientHello{
		cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, supportedGroups: []tls.CurveID{tls.X25519},
		keyShares:                   []keyShareEntry{{group: tls.X25519, data: make([]byte, 32)}},
		certificateSignatureSchemes: []tls.SignatureScheme{tls.Ed25519},
	}
	wire, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsedHello, err := parseClientHello(wire)
	if err != nil || len(parsedHello.certificateSignatureSchemes) != 1 || parsedHello.certificateSignatureSchemes[0] != tls.Ed25519 {
		t.Fatalf("ClientHello signature_algorithms_cert=%v err=%v", parsedHello.certificateSignatureSchemes, err)
	}
	request := &certificateRequestMessage{
		signatureSchemes:            []tls.SignatureScheme{tls.Ed25519},
		certificateSignatureSchemes: []tls.SignatureScheme{tls.Ed25519},
	}
	wire, err = request.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsedRequest, err := parseCertificateRequest(wire)
	if err != nil || len(parsedRequest.certificateSignatureSchemes) != 1 || parsedRequest.certificateSignatureSchemes[0] != tls.Ed25519 {
		t.Fatalf("CertificateRequest signature_algorithms_cert=%v err=%v", parsedRequest.certificateSignatureSchemes, err)
	}

	rootPublic, rootKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	leafPublic, leafKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	rootTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(10), Subject: pkix.Name{CommonName: "root"},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), IsCA: true,
		BasicConstraintsValid: true, KeyUsage: x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}
	rootDER, err := x509.CreateCertificate(rand.Reader, rootTemplate, rootTemplate, rootPublic, rootKey)
	if err != nil {
		t.Fatal(err)
	}
	root, err := x509.ParseCertificate(rootDER)
	if err != nil {
		t.Fatal(err)
	}
	leafTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(11), Subject: pkix.Name{CommonName: "leaf"},
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature,
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTemplate, root, leafPublic, rootKey)
	if err != nil {
		t.Fatal(err)
	}
	configured := &tls.Certificate{Certificate: [][]byte{leafDER, rootDER}, PrivateKey: leafKey}
	if err = validateConfiguredCertificate(configured, []tls.SignatureScheme{tls.Ed25519}, true); err != nil {
		t.Fatal(err)
	}
	if err = validateConfiguredCertificate(configured, []tls.SignatureScheme{tls.ECDSAWithP256AndSHA256}, true); err == nil {
		t.Fatal("accepted a certificate chain signature excluded by signature_algorithms_cert")
	}
}

func TestRFC9325RejectsSmallRSAServerCertificate(t *testing.T) {
	certificate, parsed, roots := makeTestRSACertificate(t, 1024, x509.SHA256WithRSA)
	if _, err := parsed.Verify(x509.VerifyOptions{Roots: roots, DNSName: "example.test"}); err != nil {
		t.Fatalf("crypto/x509 rejected the differential fixture: %v", err)
	}
	if err := validateConfiguredCertificate(&certificate, defaultSignatureSchemes(), true); err == nil {
		t.Fatal("accepted a configured 1024-bit RSA server certificate")
	}
	if err := validateConfiguredCertificate(&certificate, defaultSignatureSchemes(), false); err != nil {
		t.Fatalf("applied the server-only RSA modulus requirement to a client certificate: %v", err)
	}

	config, err := (&Config{InsecureSkipVerify: true}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = verifyCertificateChain(config, &certificateMessage{certificates: []certificateEntry{{data: parsed.Raw}}}, true, defaultSignatureSchemes())
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertBadCertificate {
		t.Fatalf("weak peer RSA certificate returned %v", err)
	}

	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	client := Client(left, &Config{RootCAs: roots, ServerName: "example.test", HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	serverErr := <-serverDone
	local = nil
	if !errors.As(serverErr, &local) || local.description != alertHandshakeFailure {
		t.Fatalf("server returned %v", serverErr)
	}
	if !errors.Is(clientErr, AlertError(alertHandshakeFailure)) {
		t.Fatalf("client returned %v", clientErr)
	}

	strong, _, _ := makeTestRSACertificate(t, 2048, x509.SHA256WithRSA)
	if err = validateConfiguredCertificate(&strong, defaultSignatureSchemes(), true); err != nil {
		t.Fatalf("rejected a 2048-bit RSA server certificate: %v", err)
	}
}

func TestRFC9325RejectsSHA1AndMD5CertificateSignatures(t *testing.T) {
	certificate, parsed, roots := makeTestRSACertificate(t, 2048, x509.SHA1WithRSA)
	if err := parsed.CheckSignature(parsed.SignatureAlgorithm, parsed.RawTBSCertificate, parsed.Signature); err != nil {
		t.Fatalf("crypto/x509 no longer accepts the SHA-1 differential fixture: %v", err)
	}
	if _, err := parsed.Verify(x509.VerifyOptions{Roots: roots, DNSName: "example.test"}); err != nil {
		t.Fatalf("crypto/x509 rejected the SHA-1 trust anchor fixture: %v", err)
	}
	if err := validateConfiguredCertificate(&certificate, append(defaultSignatureSchemes(), tls.PKCS1WithSHA1), true); err == nil {
		t.Fatal("accepted a configured SHA-1 certificate")
	}

	config, err := (&Config{RootCAs: roots, ServerName: "example.test"}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = verifyCertificateChain(config, &certificateMessage{certificates: []certificateEntry{{data: parsed.Raw}}}, true, append(defaultSignatureSchemes(), tls.PKCS1WithSHA1))
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertBadCertificate {
		t.Fatalf("SHA-1 trust anchor returned %v", err)
	}

	md5 := *parsed
	md5.SignatureAlgorithm = x509.MD5WithRSA
	if err = validateCertificateSecurityPolicy([]*x509.Certificate{&md5}, false); err == nil {
		t.Fatal("accepted an MD5 certificate signature")
	}
}

func TestRFC9325CertificatePolicyAppliesToResumption(t *testing.T) {
	configured, leaf, root, roots := makeTestCertificateWithSHA1Root(t)
	if _, err := leaf.Verify(x509.VerifyOptions{
		Roots: roots, DNSName: "example.test", KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}); err != nil {
		t.Fatalf("crypto/x509 rejected the SHA-1 trust anchor differential fixture: %v", err)
	}
	if err := validateConfiguredCertificate(&configured, defaultSignatureSchemes(), true); err == nil {
		t.Fatal("accepted a configured chain with a SHA-1 trust anchor")
	}
	peerConfig, err := (&Config{RootCAs: roots, ServerName: "example.test"}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = verifyCertificateChain(peerConfig, &certificateMessage{certificates: []certificateEntry{{data: leaf.Raw}}}, true, defaultSignatureSchemes())
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertBadCertificate {
		t.Fatalf("omitted SHA-1 trust anchor returned %v", err)
	}

	now := time.Now()
	serverConfig, err := (&Config{
		ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: roots,
		SessionTicketLifetime: time.Hour, Time: func() time.Time { return now },
	}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	ticket := &sessionTicketState{
		clientAuthAt: now.Unix(), peerCertificates: []*x509.Certificate{leaf},
		verifiedChains: [][]*x509.Certificate{{leaf, root}},
	}
	if validClientAuthenticationTicket(serverConfig, ticket) {
		t.Fatal("accepted a resumed mTLS identity with a SHA-1 trust anchor")
	}

	suite, err := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	if err != nil {
		t.Fatal(err)
	}
	cache := NewLRUClientSessionCache(1)
	cache.Put("example.test", &ClientSessionState{
		ticket: []byte{1}, psk: make([]byte, suite.hash.Size()), suite: suite.id,
		receivedAt: now, lifetime: 60, serverName: "example.test",
		peerCertificates: []*x509.Certificate{leaf}, verifiedChains: [][]*x509.Certificate{{leaf, root}},
	})
	clientConfig, err := (&Config{ServerName: "example.test", ClientSessionCache: cache, Time: func() time.Time { return now }}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	if state, _ := usableClientSession(clientConfig, left); state != nil {
		t.Fatal("accepted a client session with a SHA-1 server trust anchor")
	}
}
