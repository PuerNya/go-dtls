package dtls13

import (
	"crypto/tls"
	"testing"
)

func TestEphemeralKeyAgreement(t *testing.T) {
	for _, group := range []tls.CurveID{tls.X25519, tls.CurveP256} {
		a, err := generateEphemeralKey(group, nil)
		if err != nil {
			t.Fatal(err)
		}
		b, err := generateEphemeralKey(group, nil)
		if err != nil {
			t.Fatal(err)
		}
		as, err := a.sharedSecret(group, b.publicBytes())
		if err != nil {
			t.Fatal(err)
		}
		bs, err := b.sharedSecret(group, a.publicBytes())
		if err != nil {
			t.Fatal(err)
		}
		if string(as) != string(bs) {
			t.Fatalf("%v shared secrets differ", group)
		}
	}
}

func TestEphemeralKeyRejectsBadPeer(t *testing.T) {
	k, err := generateEphemeralKey(tls.X25519, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = k.sharedSecret(tls.CurveP256, k.publicBytes()); err == nil {
		t.Fatal("accepted wrong group")
	}
	if _, err = k.sharedSecret(tls.X25519, []byte{1, 2, 3}); err == nil {
		t.Fatal("accepted invalid public key")
	}
}
