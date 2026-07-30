package dtls13

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"testing"
)

func TestCertificateVerifySignatures(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	transcript := make([]byte, suite.hash.Size())
	transcript[0] = 42
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	_, edKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		key    crypto.Signer
		scheme tls.SignatureScheme
	}{{rsaKey, tls.PSSWithSHA256}, {ecKey, tls.ECDSAWithP256AndSHA256}, {edKey, tls.Ed25519}}
	for _, tc := range cases {
		sig, err := signCertificateVerify(rand.Reader, tc.key, tc.scheme, transcript, true)
		if err != nil {
			t.Fatal(err)
		}
		if err = verifyCertificateVerify(tc.key.Public(), tc.scheme, transcript, sig, true); err != nil {
			t.Fatalf("scheme %v: %v", tc.scheme, err)
		}
		sig[0] ^= 1
		if err = verifyCertificateVerify(tc.key.Public(), tc.scheme, transcript, sig, true); err == nil {
			t.Fatalf("scheme %v accepted tampered signature", tc.scheme)
		}
	}
}
