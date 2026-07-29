package dtls13

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"testing"
)

func TestHelloRetryRequestRoundTrip(t *testing.T) {
	h := &helloRetryRequest{cipherSuite: TLS_AES_128_GCM_SHA256, selectedGroup: tls.X25519, cookie: []byte("cookie")}
	b, err := h.marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseHelloRetryRequest(b)
	if err != nil {
		t.Fatal(err)
	}
	if got.cipherSuite != h.cipherSuite || got.selectedGroup != h.selectedGroup || !bytes.Equal(got.cookie, h.cookie) || !bytes.Equal(got.sessionID, h.sessionID) {
		t.Fatalf("round trip mismatch: %#v", got)
	}
}

func TestHelloRetryRequestAllowsGroupOnly(t *testing.T) {
	body, err := (&helloRetryRequest{cipherSuite: TLS_AES_128_GCM_SHA256, selectedGroup: tls.CurveP256}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseHelloRetryRequest(body)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.selectedGroup != tls.CurveP256 || len(parsed.cookie) != 0 {
		t.Fatalf("parsed %#v", parsed)
	}
}
func TestHelloRetryRequestRequiresAChange(t *testing.T) {
	if _, err := (&helloRetryRequest{cipherSuite: TLS_AES_128_GCM_SHA256}).marshal(); err == nil {
		t.Fatal("marshaled HRR without cookie or selected group")
	}
}
func TestHelloRetryRequestRejectsOversizedExtensionVector(t *testing.T) {
	if _, err := (&helloRetryRequest{cipherSuite: TLS_AES_128_GCM_SHA256, cookie: make([]byte, 65535)}).marshal(); err == nil {
		t.Fatal("marshaled an oversized HelloRetryRequest extension vector")
	}
}
func TestHelloRetryTranscriptRewrite(t *testing.T) {
	clientHash := sha256.Sum256([]byte("client hello"))
	tr := newTranscriptHash(sha256.New())
	if err := tr.addHelloRetryRequest(clientHash[:], []byte("hrr")); err != nil {
		t.Fatal(err)
	}
	manual := sha256.New()
	manual.Write([]byte{254, 0, 0, 32})
	manual.Write(clientHash[:])
	manual.Write([]byte{2, 0, 0, 3, 'h', 'r', 'r'})
	if !bytes.Equal(tr.sum(), manual.Sum(nil)) {
		t.Fatal("incorrect HRR transcript rewrite")
	}
}
