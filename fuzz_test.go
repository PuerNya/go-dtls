package dtls13

import (
	"bytes"
	"reflect"
	"testing"
)

func FuzzPlainRecordParsers(f *testing.F) {
	seed, err := marshalPlainRecord(record{typ: recordTypeHandshake, sequence: 7, payload: []byte("seed")})
	if err != nil {
		f.Fatal(err)
	}
	f.Add(seed)
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, data []byte) {
		owned, ownedErr := parsePlainRecords(data)
		view, viewErr := parsePlainRecordsView(data)
		var storage [2]record
		reused, reusedErr := parsePlainRecordsViewInto(data, storage[:0])
		if (ownedErr == nil) != (viewErr == nil) || (ownedErr == nil) != (reusedErr == nil) {
			t.Fatalf("ownership modes disagree: owned=%v view=%v reused=%v", ownedErr, viewErr, reusedErr)
		}
		if ownedErr != nil {
			return
		}
		if !reflect.DeepEqual(owned, view) || !reflect.DeepEqual(owned, reused) {
			t.Fatalf("ownership modes parsed different records: owned=%#v view=%#v reused=%#v", owned, view, reused)
		}
	})
}

func FuzzHandshakeFragmentParsers(f *testing.F) {
	seed, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeClientHello, messageSequence: 1, length: 4, body: []byte("seed")})
	if err != nil {
		f.Fatal(err)
	}
	f.Add(seed)
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, data []byte) {
		owned, ownedErr := parseHandshakeFragments(data)
		view, viewErr := parseHandshakeFragmentsView(data)
		var storage [2]handshakeFragment
		reused, reusedErr := parseHandshakeFragmentsViewInto(data, storage[:0])
		if (ownedErr == nil) != (viewErr == nil) || (ownedErr == nil) != (reusedErr == nil) {
			t.Fatalf("ownership modes disagree: owned=%v view=%v reused=%v", ownedErr, viewErr, reusedErr)
		}
		if ownedErr != nil {
			return
		}
		if !reflect.DeepEqual(owned, view) || !reflect.DeepEqual(owned, reused) {
			t.Fatalf("ownership modes parsed different fragments: owned=%#v view=%#v reused=%#v", owned, view, reused)
		}
	})
}

func FuzzACKParser(f *testing.F) {
	seed, err := marshalACK([]recordNumber{{epoch: 2, sequence: 1}, {epoch: 3, sequence: 4}})
	if err != nil {
		f.Fatal(err)
	}
	f.Add(seed)
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, data []byte) {
		owned, ownedErr := parseACK(data)
		var storage [4]recordNumber
		reused, reusedErr := parseACKInto(data, storage[:0])
		if (ownedErr == nil) != (reusedErr == nil) {
			t.Fatalf("ownership modes disagree: owned=%v reused=%v", ownedErr, reusedErr)
		}
		if ownedErr != nil {
			return
		}
		if !reflect.DeepEqual(owned, reused) {
			t.Fatalf("ownership modes parsed different ACKs: owned=%#v reused=%#v", owned, reused)
		}
		wire, err := marshalACK(owned)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(wire, data) {
			t.Fatalf("accepted non-canonical ACK: input=%x canonical=%x", data, wire)
		}
	})
}

func FuzzHandshakeMessageParsers(f *testing.F) {
	f.Add([]byte{})
	f.Add([]byte{0, 0, 0, 0})
	f.Fuzz(func(t *testing.T, data []byte) {
		_, _ = parseClientHello(data)
		_, _ = parseServerHello(data)
		_, _ = parseHelloRetryRequest(data)
		_, _ = parseEncryptedExtensions(data)
		_, _ = parseCertificateMessage(data, 1<<20)
		_, _ = parseCertificateHandshakeMessage(handshakeTypeCompressedCertificate, data, &certificateCompressionZlibOffer, 4096)
		_, _ = parseCertificateCompressionAlgorithms(data)
		_, _ = parseCertificateVerify(data)
		_, _ = parseFinished(data, 32)
		_, _ = parseCertificateRequest(data)
		_, _ = parseKeyUpdate(data)
		_, _ = parseNewSessionTicket(data)
		_, _ = parseNewConnectionID(data)
		_, _ = parseRequestConnectionID(data)
		_, _ = parseAlert(data)
	})
}

func FuzzProtectedRecordOpenModes(f *testing.F) {
	suiteIDs := []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384, TLS_CHACHA20_POLY1305_SHA256, TLS_AES_128_CCM_SHA256}
	for mode := uint8(0); mode < 8; mode++ {
		suite, err := cipherSuiteForID(suiteIDs[int(mode)&3])
		if err != nil {
			f.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		sender, err := newRecordCipher(suite, secret, 3, 64)
		if err != nil {
			f.Fatal(err)
		}
		if mode&4 != 0 {
			if err = sender.setConnectionID([]byte{1, 2, 3, 4}); err != nil {
				f.Fatal(err)
			}
		}
		wire, err := sender.seal(recordTypeApplicationData, []byte("fuzz record"))
		if err != nil {
			f.Fatal(err)
		}
		f.Add(mode, wire)
	}
	f.Add(uint8(0), []byte{})
	f.Fuzz(func(t *testing.T, mode uint8, data []byte) {
		suite, err := cipherSuiteForID(suiteIDs[int(mode)&3])
		if err != nil {
			t.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		ownedCipher, err := newRecordCipher(suite, secret, 3, 64)
		if err != nil {
			t.Fatal(err)
		}
		viewCipher, err := newRecordCipher(suite, secret, 3, 64)
		if err != nil {
			t.Fatal(err)
		}
		if mode&4 != 0 {
			cid := []byte{1, 2, 3, 4}
			if err = ownedCipher.setConnectionID(cid); err != nil {
				t.Fatal(err)
			}
			if err = viewCipher.setConnectionID(cid); err != nil {
				t.Fatal(err)
			}
		}
		ownedInput := append([]byte(nil), data...)
		viewInput := append([]byte(nil), data...)
		owned, ownedType, ownedConsumed, ownedErr := ownedCipher.open(ownedInput)
		view, viewType, viewConsumed, viewErr := viewCipher.openInPlace(viewInput)
		if (ownedErr == nil) != (viewErr == nil) {
			t.Fatalf("open modes disagree: owned=%v view=%v", ownedErr, viewErr)
		}
		if ownedErr != nil {
			if ownedErr.Error() != viewErr.Error() {
				t.Fatalf("open modes returned different errors: owned=%v view=%v", ownedErr, viewErr)
			}
		} else if !bytes.Equal(owned, view) || ownedType != viewType || ownedConsumed != viewConsumed {
			t.Fatalf("open modes returned different records: owned=%x/%d/%d view=%x/%d/%d", owned, ownedType, ownedConsumed, view, viewType, viewConsumed)
		}
		if !reflect.DeepEqual(ownedCipher.replay, viewCipher.replay) || ownedCipher.lastOpened != viewCipher.lastOpened || ownedCipher.authFailures != viewCipher.authFailures {
			t.Fatalf("open modes committed different state: owned=%+v/%d/%d view=%+v/%d/%d", ownedCipher.replay, ownedCipher.lastOpened, ownedCipher.authFailures, viewCipher.replay, viewCipher.lastOpened, viewCipher.authFailures)
		}
	})
}
