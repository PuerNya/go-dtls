package dtls13

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"testing"
)

func TestEphemeralKeyAgreement(t *testing.T) {
	for _, group := range []tls.CurveID{tls.X25519, tls.CurveP256, tls.CurveP384} {
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
		if !bytes.Equal(as, bs) {
			t.Fatalf("%v shared secrets differ", group)
		}
	}
}

func TestHybridKeyAgreementWireFormatAndFallback(t *testing.T) {
	for _, group := range []tls.CurveID{tls.X25519MLKEM768, tls.SecP256r1MLKEM768, tls.SecP384r1MLKEM1024} {
		t.Run(group.String(), func(t *testing.T) {
			hybrid, _ := hybridKeyExchangeForID(group)
			client, err := generateEphemeralKey(group, nil)
			if err != nil {
				t.Fatal(err)
			}
			hybridClient := client.(*hybridEphemeralKey)
			clientShare := client.publicBytes()
			if len(clientShare) != hybrid.ecdhElementSize+hybrid.mlkemPublicKeySize {
				t.Fatalf("client share length = %d", len(clientShare))
			}
			ecdhShare, mlkemShare := splitHybridShare(hybrid, clientShare, hybrid.mlkemPublicKeySize)
			if !bytes.Equal(ecdhShare, hybridClient.private.PublicKey().Bytes()) || !bytes.Equal(mlkemShare, hybridClient.mlkem.Encapsulator().Bytes()) {
				t.Fatal("client share component order is wrong")
			}

			serverShare, serverSecret, err := generateServerKeyShare(group, clientShare, nil)
			if err != nil {
				t.Fatal(err)
			}
			if len(serverShare) != hybrid.ecdhElementSize+hybrid.mlkemCiphertextSize {
				t.Fatalf("server share length = %d", len(serverShare))
			}
			serverECDH, ciphertext := splitHybridShare(hybrid, serverShare, hybrid.mlkemCiphertextSize)
			if _, err = hybrid.ecdhCurve.NewPublicKey(serverECDH); err != nil || len(ciphertext) != hybrid.mlkemCiphertextSize {
				t.Fatalf("invalid server share components: %v", err)
			}
			clientSecret, err := client.sharedSecret(group, serverShare)
			if err != nil {
				t.Fatal(err)
			}
			wantSecretLen := 64
			if group == tls.SecP384r1MLKEM1024 {
				wantSecretLen = 80
			}
			if !bytes.Equal(clientSecret, serverSecret) || len(clientSecret) != wantSecretLen {
				t.Fatalf("hybrid shared secret length = %d or values differ", len(clientSecret))
			}

			fallbackGroup, fallbackShare, ok := client.fallbackPublicBytes()
			if !ok || fallbackGroup != hybrid.ecdhGroup || !bytes.Equal(fallbackShare, ecdhShare) {
				t.Fatal("hybrid key did not reuse its ECDH component for fallback")
			}
			fallbackServerShare, fallbackServerSecret, err := generateServerKeyShare(fallbackGroup, fallbackShare, nil)
			if err != nil {
				t.Fatal(err)
			}
			fallbackClientSecret, err := client.sharedSecret(fallbackGroup, fallbackServerShare)
			if err != nil || !bytes.Equal(fallbackClientSecret, fallbackServerSecret) {
				t.Fatalf("fallback key agreement failed: %v", err)
			}
		})
	}
}

func TestHybridKeyExchangeRejectsInvalidPeerInputs(t *testing.T) {
	for _, group := range []tls.CurveID{tls.X25519MLKEM768, tls.SecP256r1MLKEM768, tls.SecP384r1MLKEM1024} {
		t.Run(group.String(), func(t *testing.T) {
			hybrid, _ := hybridKeyExchangeForID(group)
			client, err := generateEphemeralKey(group, nil)
			if err != nil {
				t.Fatal(err)
			}
			clientShare := client.publicBytes()
			for _, bad := range [][]byte{clientShare[:len(clientShare)-1], append(append([]byte(nil), clientShare...), 0)} {
				_, _, err = generateServerKeyShare(group, bad, nil)
				requireIllegalParameter(t, err)
			}
			badPublic := append([]byte(nil), clientShare...)
			if hybrid.mlkemFirst {
				for i := range hybrid.mlkemPublicKeySize {
					badPublic[i] = 0xff
				}
			} else {
				for i := hybrid.ecdhElementSize; i < len(badPublic); i++ {
					badPublic[i] = 0xff
				}
			}
			_, _, err = generateServerKeyShare(group, badPublic, nil)
			requireIllegalParameter(t, err)

			serverShare, _, err := generateServerKeyShare(group, clientShare, nil)
			if err != nil {
				t.Fatal(err)
			}
			_, err = client.sharedSecret(group, serverShare[:len(serverShare)-1])
			requireIllegalParameter(t, err)
			badECDH := append([]byte(nil), serverShare...)
			serverECDH, _ := splitHybridShare(hybrid, badECDH, hybrid.mlkemCiphertextSize)
			clear(serverECDH)
			_, err = client.sharedSecret(group, badECDH)
			requireIllegalParameter(t, err)
		})
	}
}

func requireIllegalParameter(t *testing.T, err error) {
	t.Helper()
	if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("peer key error = %v; alert = %d, %v", err, description, ok)
	}
}

type failingEntropyReader struct{}

func (failingEntropyReader) Read([]byte) (int, error) { return 0, io.ErrUnexpectedEOF }

func TestServerKeyShareClassifiesEntropyFailureAsInternalError(t *testing.T) {
	t.Setenv("GODEBUG", "cryptocustomrand=1")
	client, err := generateEphemeralKey(tls.X25519MLKEM768, nil)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = generateServerKeyShare(tls.X25519MLKEM768, client.publicBytes(), failingEntropyReader{})
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertInternalError {
		t.Fatalf("entropy failure = %v", err)
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
