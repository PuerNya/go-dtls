package dtls13

import (
	"crypto/ed25519"
	"crypto/tls"
	"crypto/x509"
	"testing"
)

type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	clear(p)
	return len(p), nil
}

func TestPostHandshakeAuthContextIsUniqueWithDeterministicRandom(t *testing.T) {
	c := &Conn{config: &Config{Rand: zeroReader{}}}
	first, err := c.newPostHandshakeAuthContext()
	if err != nil {
		t.Fatal(err)
	}
	second, err := c.newPostHandshakeAuthContext()
	if err != nil {
		t.Fatal(err)
	}
	if equalBytes(first, second) {
		t.Fatal("post-handshake authentication contexts repeated")
	}
}

func TestCertificateRequestRoundTrip(t *testing.T) {
	want := &certificateRequestMessage{requestContext: []byte{1, 2}, signatureSchemes: []tls.SignatureScheme{tls.Ed25519, tls.PSSWithSHA256}}
	b, err := want.marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseCertificateRequest(b)
	if err != nil {
		t.Fatal(err)
	}
	if string(got.requestContext) != string(want.requestContext) || len(got.signatureSchemes) != 2 || got.signatureSchemes[0] != tls.Ed25519 {
		t.Fatalf("got %#v", got)
	}
}

func TestCertificateRequestRequiresSignatureAlgorithms(t *testing.T) {
	_, err := parseCertificateRequest([]byte{0, 0, 0})
	if description, ok := protocolAlert(err); !ok || description != alertMissingExtension {
		t.Fatalf("parseCertificateRequest alert = %d, %v; want missing_extension", description, err)
	}
}

func TestCertificateRequestRejectsRecognizedExtensionInWrongMessage(t *testing.T) {
	signatures, err := marshalSignatureSchemes([]tls.SignatureScheme{tls.Ed25519})
	if err != nil {
		t.Fatal(err)
	}
	exts, err := marshalExtensions(map[uint16][]byte{
		extSignatureAlgorithms: signatures,
		extEarlyData:           nil,
	}, []uint16{extSignatureAlgorithms, extEarlyData})
	if err != nil {
		t.Fatal(err)
	}
	body := append([]byte{0}, exts...)
	_, err = parseCertificateRequest(body)
	if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("wrong-message extension alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestPostHandshakeCertificateVerifyParseErrorDoesNotPanic(t *testing.T) {
	suite, err := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	if err != nil {
		t.Fatal(err)
	}
	c := &Conn{receivingTraffic: &receivingTraffic{suite: suite}}
	state := &postHandshakeAuthState{
		stage:            postAuthExpectCertificateVerify,
		peerCertificates: []*x509.Certificate{{PublicKey: ed25519.PublicKey(make([]byte, ed25519.PublicKeySize))}},
		signatureSchemes: []tls.SignatureScheme{tls.Ed25519},
	}
	err = c.processPostHandshakeAuthMessageLocked(state, completedHandshake{typ: handshakeTypeCertificateVerify, body: []byte{0}})
	if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
		t.Fatalf("malformed post-handshake CertificateVerify alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestKeyUpdateMessageAndState(t *testing.T) {
	for _, request := range []bool{false, true} {
		got, err := parseKeyUpdate((keyUpdateMessage{requestUpdate: request}).marshal())
		if err != nil || got.requestUpdate != request {
			t.Fatalf("got %#v err=%v", got, err)
		}
	}
	if _, err := parseKeyUpdate([]byte{2}); err == nil {
		t.Fatal("accepted invalid KeyUpdate")
	}
	var state keyUpdateState
	number := recordNumber{epoch: 3, sequence: 9}
	if err := state.begin(number); err != nil {
		t.Fatal(err)
	}
	if state.canUseNewKeys() {
		t.Fatal("allowed new keys before ACK")
	}
	if err := state.begin(recordNumber{}); err == nil {
		t.Fatal("allowed concurrent KeyUpdate")
	}
	if state.ack([]recordNumber{{epoch: 3, sequence: 8}}) {
		t.Fatal("accepted unrelated ACK")
	}
	if !state.ack([]recordNumber{number}) || !state.canUseNewKeys() {
		t.Fatal("did not release keys after ACK")
	}
}
