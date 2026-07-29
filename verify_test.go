package dtls13

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"math/big"
	"testing"
	"time"
)

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
	if err = validateConfiguredCertificate(configured, []tls.SignatureScheme{tls.Ed25519}); err != nil {
		t.Fatal(err)
	}
	if err = validateConfiguredCertificate(configured, []tls.SignatureScheme{tls.ECDSAWithP256AndSHA256}); err == nil {
		t.Fatal("accepted a certificate chain signature excluded by signature_algorithms_cert")
	}
}
