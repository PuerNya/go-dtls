package dtls13

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"testing"
)

func TestProtectedRecordReceiveErrorClassification(t *testing.T) {
	unauthenticated := &ProtocolError{"protected record authentication failed"}
	if got := protectedRecordReceiveError(unauthenticated); got != nil {
		t.Fatalf("unauthenticated error became fatal: %v", got)
	}

	limit := fmt.Errorf("wrapped limit: %w", errAEADAuthenticationFailureLimit)
	if got := protectedRecordReceiveError(limit); !errors.Is(got, errAEADAuthenticationFailureLimit) {
		t.Fatalf("authentication failure limit returned %v", got)
	}

	authenticated := fmt.Errorf("wrapped authenticated error: %w", authenticatedRecordAlert(alertRecordOverflow, &ProtocolError{"oversized record"}))
	got := protectedRecordReceiveError(authenticated)
	description, ok := protocolAlert(got)
	if !ok || description != alertRecordOverflow {
		t.Fatalf("authenticated error returned description=%d ok=%v err=%v", description, ok, got)
	}
}

func sealInvalidInnerType(t *testing.T, cipher *recordCipher, innerType uint8) []byte {
	t.Helper()
	return sealRawInnerPlaintext(t, cipher, []byte{innerType})
}

func sealRawInnerPlaintext(t *testing.T, cipher *recordCipher, plain []byte) []byte {
	t.Helper()
	sequence := cipher.nextSequence
	header := make([]byte, 5)
	header[0] = unifiedHeaderFixed | unifiedHeaderSequence16 | unifiedHeaderLength | byte(cipher.epoch&unifiedHeaderEpochMask)
	binary.BigEndian.PutUint16(header[1:3], uint16(sequence))
	binary.BigEndian.PutUint16(header[3:5], uint16(len(plain)+cipher.aead.Overhead()))
	ciphertext := cipher.aead.Seal(nil, cipher.nonce(sequence), plain, header)
	mask, err := cipher.sequenceMask.mask(ciphertext[:sequenceProtectionSampleSize])
	if err != nil {
		t.Fatal(err)
	}
	header[1] ^= mask[0]
	header[2] ^= mask[1]
	cipher.nextSequence++
	return append(header, ciphertext...)
}

func TestProtectedRecordRejectsOversizedAuthenticatedInnerPlaintext(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	plain := make([]byte, (1<<14)+2)
	plain[len(plain)-1] = recordTypeApplicationData
	wire := sealRawInnerPlaintext(t, sender, plain)
	_, _, _, err := receiver.open(wire)
	var authenticatedErr *authenticatedRecordError
	if !errors.As(err, &authenticatedErr) || authenticatedErr.description != alertRecordOverflow {
		t.Fatalf("oversized authenticated inner plaintext returned %v", err)
	}
	if receiver.replay.nextExpected() != 0 {
		t.Fatal("oversized authenticated inner plaintext committed replay state")
	}
}

func TestProtectedRecordEnforcesNegotiatedPlaintextLimit(t *testing.T) {
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	sender.setPlaintextLimit(minRecordSizeLimit)
	if _, err := sender.seal(recordTypeApplicationData, make([]byte, minRecordSizeLimit-1)); err != nil {
		t.Fatal(err)
	}
	if _, err := sender.seal(recordTypeApplicationData, make([]byte, minRecordSizeLimit)); err == nil {
		t.Fatal("sent content whose inner plaintext exceeds record_size_limit")
	}

	oversizedSender, limitedReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	limitedReceiver.setPlaintextLimit(minRecordSizeLimit)
	wire, err := oversizedSender.seal(recordTypeApplicationData, make([]byte, minRecordSizeLimit))
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, err = limitedReceiver.open(wire)
	var authenticatedErr *authenticatedRecordError
	if !errors.As(err, &authenticatedErr) || authenticatedErr.description != alertRecordOverflow {
		t.Fatalf("oversized authenticated record returned %v", err)
	}
	if limitedReceiver.replay.nextExpected() != 0 {
		t.Fatal("record_size_limit failure committed replay state")
	}
}

func TestProtectedRecordRejectsZeroLengthHandshakeAndAlert(t *testing.T) {
	for _, contentType := range []uint8{recordTypeHandshake, recordTypeAlert} {
		sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
		if _, err := sender.seal(contentType, nil); err == nil {
			t.Fatalf("sent zero-length protected content type %d", contentType)
		}
		wire := sealRawInnerPlaintext(t, sender, []byte{contentType})
		_, _, _, err := receiver.open(wire)
		var authenticatedErr *authenticatedRecordError
		if !errors.As(err, &authenticatedErr) || authenticatedErr.description != alertUnexpectedMessage {
			t.Fatalf("zero-length protected content type %d returned %v", contentType, err)
		}
	}
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	wire, err := sender.seal(recordTypeApplicationData, nil)
	if err != nil {
		t.Fatalf("zero-length application data was rejected: %v", err)
	}
	content, contentType, _, err := receiver.open(wire)
	if err != nil || contentType != recordTypeApplicationData || len(content) != 0 {
		t.Fatalf("zero-length application data content=%x type=%d err=%v", content, contentType, err)
	}
}

func recordCipherPair(t *testing.T, suiteID uint16, epoch uint64) (*recordCipher, *recordCipher) {
	t.Helper()
	suite, err := cipherSuiteForID(suiteID)
	if err != nil {
		t.Fatal(err)
	}
	secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
	sender, err := newRecordCipher(suite, secret, epoch, 64)
	if err != nil {
		t.Fatal(err)
	}
	receiver, err := newRecordCipher(suite, secret, epoch, 64)
	if err != nil {
		t.Fatal(err)
	}
	return sender, receiver
}

func TestProtectedACKRecord(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	wire, err := sender.seal(recordTypeACK, []byte{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	content, typ, _, err := receiver.open(wire)
	if err != nil {
		t.Fatal(err)
	}
	if typ != recordTypeACK || !bytes.Equal(content, []byte{0, 0}) {
		t.Fatalf("got type=%d content=%x", typ, content)
	}
}

func TestSealHandshakeFragmentIntoMatchesAllocated(t *testing.T) {
	allocated, into := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	fragment := handshakeFragment{typ: handshakeTypeCertificate, messageSequence: 7, length: 4, body: []byte("body")}
	payload, err := marshalHandshakeFragment(fragment)
	if err != nil {
		t.Fatal(err)
	}
	want, err := allocated.seal(recordTypeHandshake, payload)
	if err != nil {
		t.Fatal(err)
	}
	storage := make([]byte, 0, len(want))
	got, err := into.sealHandshakeFragmentInto(storage, fragment)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) || len(got) != cap(got) || &got[0] != &storage[:cap(storage)][0] {
		t.Fatalf("into=%x allocated=%x len=%d cap=%d", got, want, len(got), cap(got))
	}
}

func TestProtectedRecordReportsAuthenticatedInvalidInnerType(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	wire := sealInvalidInnerType(t, sender, recordTypeChangeCipherSpec)
	_, _, _, err := receiver.open(wire)
	var authenticatedErr *authenticatedRecordError
	if !errors.As(err, &authenticatedErr) {
		t.Fatalf("invalid authenticated inner type returned %v", err)
	}
}

func TestProtectedRecordRoundTrip(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384, TLS_CHACHA20_POLY1305_SHA256, TLS_AES_128_CCM_SHA256} {
		sender, receiver := recordCipherPair(t, suiteID, 2)
		wire, err := sender.seal(recordTypeApplicationData, []byte("hello"))
		if err != nil {
			t.Fatal(err)
		}
		if wire[0]&0xe0 != unifiedHeaderFixed {
			t.Fatalf("invalid unified header %02x", wire[0])
		}
		content, typ, n, err := receiver.open(wire)
		if err != nil {
			t.Fatal(err)
		}
		if typ != recordTypeApplicationData || !bytes.Equal(content, []byte("hello")) || n != len(wire) {
			t.Fatalf("unexpected open result %q/%d/%d", content, typ, n)
		}
	}
}

func TestProtectedRecordOpenOwnershipModes(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	wire, err := sender.seal(recordTypeApplicationData, []byte("independent"))
	if err != nil {
		t.Fatal(err)
	}
	wireBefore := append([]byte(nil), wire...)
	content, _, _, err := receiver.open(wire)
	if err != nil {
		t.Fatal(err)
	}
	content[0] ^= 0xff
	if !bytes.Equal(wire, wireBefore) {
		t.Fatal("open returned content that aliases the input datagram")
	}

	sender, receiver = recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	wire, err = sender.seal(recordTypeApplicationData, []byte("in-place"))
	if err != nil {
		t.Fatal(err)
	}
	content, _, _, err = receiver.openInPlace(wire)
	if err != nil {
		t.Fatal(err)
	}
	content[0] ^= 0xff
	if wire[unifiedHeaderLen16] != content[0] {
		t.Fatal("openInPlace did not reuse the input datagram")
	}
}

func TestProtectedRecordMultipleAndSequenceRollover(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	sender.nextSequence = 65535
	one, err := sender.seal(recordTypeHandshake, []byte("one"))
	if err != nil {
		t.Fatal(err)
	}
	two, err := sender.seal(recordTypeHandshake, []byte("two"))
	if err != nil {
		t.Fatal(err)
	}
	datagram := append(one, two...)
	content, _, n, err := receiver.open(datagram)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "one" {
		t.Fatalf("got %q", content)
	}
	content, _, used, err := receiver.open(datagram[n:])
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "two" || n+used != len(datagram) {
		t.Fatalf("got %q, consumed %d", content, n+used)
	}
}

func TestProtectedRecordRejectsTamperWithoutCommittingReplay(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	wire, err := sender.seal(recordTypeApplicationData, []byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	bad := append([]byte(nil), wire...)
	bad[len(bad)-1] ^= 1
	if _, _, _, err = receiver.open(bad); err == nil {
		t.Fatal("accepted tampered ciphertext")
	}
	if content, _, _, openErr := receiver.open(wire); openErr != nil || string(content) != "secret" {
		t.Fatalf("valid record after tamper: %q, %v", content, openErr)
	}
	if _, _, _, err = receiver.open(wire); err == nil {
		t.Fatal("accepted replay")
	}
}

func TestSequenceReconstruction(t *testing.T) {
	cases := []struct{ expected, truncated, want uint64 }{{0, 7, 7}, {65536, 0, 65536}, {65535, 0, 65536}, {65536, 65535, 65535}, {131071, 0, 131072}}
	for _, tc := range cases {
		if got := reconstructSequence(tc.expected, tc.truncated, 16); got != tc.want {
			t.Fatalf("expected=%d truncated=%d: got %d want %d", tc.expected, tc.truncated, got, tc.want)
		}
	}
}

func TestRecordNonceUsesSequenceWithoutEpoch(t *testing.T) {
	a, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	b, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	if !bytes.Equal(a.nonce(99), b.nonce(99)) {
		t.Fatal("epoch was mixed into the AEAD nonce")
	}
}

func TestProtectedRecordHeaderVariants(t *testing.T) {
	for _, sequence16 := range []bool{false, true} {
		for _, length := range []bool{false, true} {
			sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
			wire, err := sender.sealWithHeader(recordTypeApplicationData, []byte("variant"), sequence16, length)
			if err != nil {
				t.Fatal(err)
			}
			content, typ, consumed, err := receiver.open(wire)
			if err != nil {
				t.Fatalf("S=%v L=%v: %v", sequence16, length, err)
			}
			if string(content) != "variant" || typ != recordTypeApplicationData || consumed != len(wire) {
				t.Fatalf("S=%v L=%v bad result", sequence16, length)
			}
		}
	}
}

func TestAEADUsageLimit(t *testing.T) {
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	sender.nextSequence = sender.recordLimit
	if _, err := sender.seal(recordTypeApplicationData, []byte{1}); err == nil {
		t.Fatal("sealed beyond AEAD usage limit")
	}
}

func TestAuthFailureKeyUpdateThreshold(t *testing.T) {
	_, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	receiver.authFailureLimit = 8
	for failures := uint64(0); failures <= 8; failures++ {
		receiver.authFailures = failures
		want := failures >= 6
		if got := receiver.shouldRequestKeyUpdateForAuthFailures(); got != want {
			t.Fatalf("failures=%d: got %v want %v", failures, got, want)
		}
	}
}

func TestProtectedRecordConnectionIDHeaderVariants(t *testing.T) {
	connectionID := []byte{0xaa, 0xbb, 0xcc, 0xdd}
	for _, sequence16 := range []bool{false, true} {
		for _, length := range []bool{false, true} {
			sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
			if err := sender.setConnectionID(connectionID); err != nil {
				t.Fatal(err)
			}
			if err := receiver.setConnectionID(connectionID); err != nil {
				t.Fatal(err)
			}
			wire, err := sender.sealWithHeader(recordTypeApplicationData, []byte("cid"), sequence16, length)
			if err != nil {
				t.Fatal(err)
			}
			if wire[0]&unifiedHeaderCID == 0 || !bytes.Equal(wire[1:1+len(connectionID)], connectionID) {
				t.Fatalf("CID missing from unified header: %x", wire)
			}
			content, typ, consumed, err := receiver.open(wire)
			if err != nil {
				t.Fatalf("S=%v L=%v: %v", sequence16, length, err)
			}
			if string(content) != "cid" || typ != recordTypeApplicationData || consumed != len(wire) {
				t.Fatalf("S=%v L=%v bad CID result", sequence16, length)
			}
		}
	}
}

func TestProtectedRecordEmptyConnectionIDIsPresent(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	if err := sender.setConnectionID([]byte{}); err != nil {
		t.Fatal(err)
	}
	if err := receiver.setConnectionID([]byte{}); err != nil {
		t.Fatal(err)
	}
	wire, err := sender.seal(recordTypeApplicationData, []byte("empty CID"))
	if err != nil {
		t.Fatal(err)
	}
	if wire[0]&unifiedHeaderCID == 0 {
		t.Fatalf("empty negotiated CID did not set the C bit: %x", wire)
	}
	content, _, _, err := receiver.open(wire)
	if err != nil || string(content) != "empty CID" {
		t.Fatalf("empty CID record failed: %q %v", content, err)
	}
}

func TestProtectedRecordAcceptsSpareConnectionID(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	if err := sender.setConnectionID([]byte{3, 4, 5}); err != nil {
		t.Fatal(err)
	}
	if err := receiver.setConnectionID([]byte{1, 2}); err != nil {
		t.Fatal(err)
	}
	if err := receiver.addAcceptedConnectionIDs([][]byte{{3, 4, 5}}); err != nil {
		t.Fatal(err)
	}
	wire, err := sender.seal(recordTypeApplicationData, []byte("spare"))
	if err != nil {
		t.Fatal(err)
	}
	content, _, _, err := receiver.open(wire)
	if err != nil || string(content) != "spare" {
		t.Fatalf("spare CID record failed: %q %v", content, err)
	}
	if err = receiver.addAcceptedConnectionIDs([][]byte{{3, 4}}); err == nil {
		t.Fatal("accepted ambiguous CID prefixes")
	}
}

func TestProtectedRecordRetainsOwnedMatchedConnectionID(t *testing.T) {
	cid := []byte{1, 2, 3, 4}
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	if err := sender.setConnectionID(cid); err != nil {
		t.Fatal(err)
	}
	if err := receiver.setConnectionID(cid); err != nil {
		t.Fatal(err)
	}
	wire, err := sender.seal(recordTypeApplicationData, []byte("owned CID"))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, _, err = receiver.open(wire); err != nil {
		t.Fatal(err)
	}
	wire[1] ^= 0xff
	if !bytes.Equal(receiver.lastConnectionID, cid) {
		t.Fatalf("matched CID aliases input datagram: %x", receiver.lastConnectionID)
	}
}

func TestProtectedRecordRejectsConnectionIDMismatchAndTamper(t *testing.T) {
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	if err := sender.setConnectionID([]byte{1, 2, 3}); err != nil {
		t.Fatal(err)
	}
	wire, err := sender.seal(recordTypeApplicationData, []byte("authenticated CID"))
	if err != nil {
		t.Fatal(err)
	}
	_, receiverWithoutCID := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	if _, _, _, err = receiverWithoutCID.open(wire); err == nil {
		t.Fatal("accepted a CID record in an epoch without CID")
	}
	_, wrongReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	_ = wrongReceiver.setConnectionID([]byte{1, 2, 4})
	if _, _, _, err = wrongReceiver.open(wire); err == nil {
		t.Fatal("accepted the wrong connection ID")
	}
	tampered := append([]byte(nil), wire...)
	tampered[3] ^= 1
	_, tamperedReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	_ = tamperedReceiver.setConnectionID([]byte{1, 2, 2})
	if _, _, _, err = tamperedReceiver.open(tampered); err == nil {
		t.Fatal("CID tamper was not detected by record authentication")
	}
	_, validReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	_ = validReceiver.setConnectionID([]byte{1, 2, 3})
	if content, _, _, err := validReceiver.open(wire); err != nil || string(content) != "authenticated CID" {
		t.Fatalf("valid CID record failed after attacks: %q %v", content, err)
	}
	if err = sender.setConnectionID(make([]byte, 256)); err == nil {
		t.Fatal("accepted a connection ID over 255 bytes")
	}
}
