package dtls13

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"testing"
)

func TestNegotiation(t *testing.T) {
	suite, err := selectCipherSuite(defaultCipherSuites(), []uint16{0xffff, TLS_AES_256_GCM_SHA384})
	if err != nil || suite.id != TLS_AES_256_GCM_SHA384 {
		t.Fatalf("suite=%v err=%v", suite, err)
	}
	share, err := selectKeyShare([]tls.CurveID{tls.X25519}, []keyShareEntry{{group: tls.CurveP256}, {group: tls.X25519, data: []byte{1}}})
	if err != nil || share.group != tls.X25519 {
		t.Fatalf("share=%v err=%v", share, err)
	}
	proto, err := negotiateALPN([]string{"h3", "coap"}, []string{"coap", "h3"})
	if err != nil || proto != "h3" {
		t.Fatalf("proto=%q err=%v", proto, err)
	}
}

func TestCipherSuiteUsesServerPreference(t *testing.T) {
	preferences := []uint16{TLS_CHACHA20_POLY1305_SHA256, TLS_AES_128_GCM_SHA256}
	offered := []uint16{TLS_AES_128_GCM_SHA256, TLS_CHACHA20_POLY1305_SHA256}
	suite, err := selectCipherSuite(preferences, offered)
	if err != nil || suite.id != TLS_CHACHA20_POLY1305_SHA256 {
		t.Fatalf("suite=%v err=%v", suite, err)
	}
	if _, err = selectCipherSuite(preferences, []uint16{TLS_AES_256_GCM_SHA384}); err == nil {
		t.Fatal("selected a cipher suite outside the configured intersection")
	} else if description, ok := protocolAlert(err); !ok || description != alertHandshakeFailure {
		t.Fatalf("cipher alert=%d ok=%v", description, ok)
	}
}
func TestSignatureNegotiation(t *testing.T) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	scheme, err := selectSignatureScheme(key, []tls.SignatureScheme{tls.PSSWithSHA256, tls.Ed25519})
	if err != nil || scheme != tls.Ed25519 {
		t.Fatalf("scheme=%v err=%v", scheme, err)
	}
	if _, err = selectSignatureScheme(key, []tls.SignatureScheme{tls.PSSWithSHA256}); err == nil {
		t.Fatal("selected incompatible scheme")
	} else if description, ok := protocolAlert(err); !ok || description != alertHandshakeFailure {
		t.Fatalf("signature alert=%d ok=%v", description, ok)
	}
}

func TestALPNFailureUsesNoApplicationProtocol(t *testing.T) {
	if _, err := negotiateALPN([]string{"h3"}, []string{"coap"}); err == nil {
		t.Fatal("negotiated disjoint ALPN lists")
	} else if description, ok := protocolAlert(err); !ok || description != alertNoApplicationProtocol {
		t.Fatalf("ALPN alert=%d ok=%v", description, ok)
	}
}
