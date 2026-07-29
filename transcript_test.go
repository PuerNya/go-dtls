package dtls13

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"testing"
)

func TestTranscriptCanonicalHeader(t *testing.T) {
	tr := newTranscriptHash(sha256.New())
	if err := tr.add(1, 9, []byte("abc")); err != nil {
		t.Fatal(err)
	}
	wantInput := []byte{1, 0, 0, 3, 'a', 'b', 'c'}
	want := sha256.Sum256(wantInput)
	got := tr.sum()
	if string(got) != string(want[:]) {
		t.Fatalf("got %x, want %x", got, want)
	}
}

func TestTranscriptCloneIsIndependent(t *testing.T) {
	for _, newHash := range []func() hash.Hash{sha256.New, sha512.New384} {
		original := newTranscriptHash(newHash())
		_ = original.add(handshakeTypeClientHello, 0, bytes.Repeat([]byte("hello"), 17))
		clone := original.clone()
		before := append([]byte(nil), original.sum()...)
		_ = clone.add(handshakeTypeCertificateRequest, 1, []byte("request"))
		if !bytes.Equal(original.sum(), before) {
			t.Fatal("mutating clone changed original transcript")
		}
		if bytes.Equal(clone.sum(), original.sum()) {
			t.Fatal("clone did not advance independently")
		}
	}
}

func TestTranscriptSumIntoMatchesOwnedResult(t *testing.T) {
	for _, newHash := range []func() hash.Hash{sha256.New, sha512.New384} {
		transcript := newTranscriptHash(newHash())
		_ = transcript.add(handshakeTypeCertificate, 1, bytes.Repeat([]byte{0x5a}, 257))
		var scratch [maxSupportedHashSize]byte
		got := transcript.sumInto(scratch[:0])
		if !bytes.Equal(got, transcript.sum()) || &got[0] != &scratch[0] {
			t.Fatal("caller-storage transcript hash differs from owned result")
		}
	}
}

func TestTranscriptExcludesMessageSequence(t *testing.T) {
	a := newTranscriptHash(sha256.New())
	b := newTranscriptHash(sha256.New())
	if err := a.add(1, 1, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	if err := b.add(1, 999, []byte("hello")); err != nil {
		t.Fatal(err)
	}
	if string(a.sum()) != string(b.sum()) {
		t.Fatal("message_seq affected transcript")
	}
}
