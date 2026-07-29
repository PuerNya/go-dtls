package dtls13

import (
	"bytes"
	"crypto/hkdf"
	"encoding/hex"
	"testing"
)

func TestCipherSuiteDescriptorsAreShared(t *testing.T) {
	for _, id := range defaultCipherSuites() {
		first, err := cipherSuiteForID(id)
		if err != nil {
			t.Fatal(err)
		}
		second, err := cipherSuiteForID(id)
		if err != nil {
			t.Fatal(err)
		}
		if first != second {
			t.Fatalf("cipher suite %04x descriptor was reallocated", id)
		}
	}
}

func TestHKDFRFC5869Case1(t *testing.T) {
	ikm := bytes.Repeat([]byte{0x0b}, 22)
	salt, _ := hex.DecodeString("000102030405060708090a0b0c")
	info, _ := hex.DecodeString("f0f1f2f3f4f5f6f7f8f9")
	prkWant, _ := hex.DecodeString("077709362c2e32df0ddc3f0dc47bba6390b6c73bb50f9c3122ec844ad7c2b3e5")
	okmWant, _ := hex.DecodeString("3cb25f25faacd57a90434f64d0362f2a2d2d0a90cf1a5a4c5db02d56ecc4c5bf34007208d5b887185865")
	s, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	prk := hkdfExtract(s.hash.New, ikm, salt)
	if !bytes.Equal(prk, prkWant) {
		t.Fatalf("PRK %x", prk)
	}
	if got := hkdfExpand(s.hash.New, prk, info, 42); !bytes.Equal(got, okmWant) {
		t.Fatalf("OKM %x", got)
	}
}

func TestHKDFExtractIntoMatchesStandardLibrary(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		for _, secret := range [][]byte{nil, {}, bytes.Repeat([]byte{0x0b}, 22), bytes.Repeat([]byte{0xa5}, 2*suite.hash.Size()+3)} {
			for _, salt := range [][]byte{nil, {}, bytes.Repeat([]byte{0x11}, suite.hash.Size()), bytes.Repeat([]byte{0x22}, suite.hash.Size()+7)} {
				want, err := hkdf.Extract(suite.hash.New, secret, salt)
				if err != nil {
					t.Fatal(err)
				}
				got := bytes.Repeat([]byte{0xff}, suite.hash.Size())
				hkdfExtractInto(suite.hash.New, secret, salt, got)
				if !bytes.Equal(got, want) {
					t.Fatalf("suite %04x secret=%d salt=%d differs from crypto/hkdf", suiteID, len(secret), len(salt))
				}
			}
		}
		saltAndOutput := bytes.Repeat([]byte{0x33}, suite.hash.Size())
		want, err := hkdf.Extract(suite.hash.New, bytes.Repeat([]byte{0x44}, suite.hash.Size()+5), append([]byte(nil), saltAndOutput...))
		if err != nil {
			t.Fatal(err)
		}
		hkdfExtractInto(suite.hash.New, bytes.Repeat([]byte{0x44}, suite.hash.Size()+5), saltAndOutput, saltAndOutput)
		if !bytes.Equal(saltAndOutput, want) {
			t.Fatalf("suite %04x overlapping salt and output differs from crypto/hkdf", suiteID)
		}
	}
}

func TestHKDFExpandMatchesStandardLibrary(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		info := bytes.Repeat([]byte{0xa5}, 97)
		lengths := []int{0, 1, suite.hash.Size() - 1, suite.hash.Size(), suite.hash.Size() + 1, 2*suite.hash.Size() + 7, 255 * suite.hash.Size()}
		for _, length := range lengths {
			want, err := hkdf.Expand(suite.hash.New, secret, string(info), length)
			if err != nil {
				t.Fatal(err)
			}
			if got := hkdfExpand(suite.hash.New, secret, info, length); !bytes.Equal(got, want) {
				t.Fatalf("suite %04x length %d differs from crypto/hkdf", suiteID, length)
			}
		}

		expander := newHKDFExpander(suite.hash.New, secret)
		for _, label := range [][]byte{[]byte("first"), []byte("second"), []byte("first")} {
			got := expander.expand(label, suite.hash.Size()+3)
			want, err := hkdf.Expand(suite.hash.New, secret, string(label), suite.hash.Size()+3)
			if err != nil || !bytes.Equal(got, want) {
				t.Fatalf("suite %04x sequential expand label %q differs: err=%v", suiteID, label, err)
			}
		}
	}
}

func TestDTLSLabelAndTrafficKeySizes(t *testing.T) {
	s, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	secret := bytes.Repeat([]byte{0x42}, s.hash.Size())
	keys := deriveTrafficKeys(s, secret)
	if len(keys.key) != 16 || len(keys.iv) != 12 || len(keys.sn) != 16 {
		t.Fatalf("unexpected key sizes: %d/%d/%d", len(keys.key), len(keys.iv), len(keys.sn))
	}
	if bytes.Equal(expandLabel(s, secret, "key", nil, 16), hkdfExpand(s.hash.New, secret, []byte("key"), 16)) {
		t.Fatal("label encoding was not applied")
	}
	context := []byte{1, 2, 3}
	info := []byte{0, 16, byte(len("dtls13key"))}
	info = append(info, "dtls13key"...)
	info = append(info, byte(len(context)))
	info = append(info, context...)
	if got, want := expandLabel(s, secret, "key", context, 16), hkdfExpand(s.hash.New, secret, info, 16); !bytes.Equal(got, want) {
		t.Fatalf("HKDF label includes an incorrect DTLS prefix: got %x want %x", got, want)
	}
}

func TestSingleBlockHKDFLabelsMatchStandardAndRemainImmutable(t *testing.T) {
	labels := []*singleBlockHKDFLabel{
		labelClientEarlyTraffic,
		labelDerived,
		labelFinished,
		labelResumptionMaster,
		labelResumptionBinder,
		labelTrafficUpdate,
		labelResumption,
		newSingleBlockHKDFLabel("dynamic exporter"),
	}
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		context := bytes.Repeat([]byte{0xa5}, suite.hash.Size())
		for _, label := range labels {
			fullBefore := append([]byte(nil), label.full...)
			sha256HeaderBefore := label.sha256Header
			sha384HeaderBefore := label.sha384Header
			got := make([]byte, suite.hash.Size())
			expandLabelHashInto(suite, secret, label, context, got)
			name := string(label.full[len("dtls13"):])
			want := expandLabel(suite, secret, name, context, suite.hash.Size())
			if !bytes.Equal(got, want) {
				t.Fatalf("suite %04x label %q differs from standard HKDF path", suiteID, name)
			}
			inPlace := append([]byte(nil), secret...)
			expandLabelHashInto(suite, inPlace, label, context, inPlace)
			if !bytes.Equal(inPlace, want) {
				t.Fatalf("suite %04x in-place label %q differs from standard HKDF path", suiteID, name)
			}
			if !bytes.Equal(label.full, fullBefore) || label.sha256Header != sha256HeaderBefore || label.sha384Header != sha384HeaderBefore {
				t.Fatalf("suite %04x label %q descriptor was modified", suiteID, name)
			}
		}
	}
}

func TestDeriveTrafficKeysIntoMatchesOwnedResult(t *testing.T) {
	for _, suiteID := range defaultCipherSuites() {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		secret := bytes.Repeat([]byte{byte(suiteID)}, suite.hash.Size())
		want := deriveTrafficKeys(suite, secret)
		key := make([]byte, suite.keyLen)
		iv := make([]byte, suite.ivLen)
		sn := make([]byte, suite.keyLen)
		deriveTrafficKeysInto(suite, secret, key, iv, sn)
		if !bytes.Equal(key, want.key) || !bytes.Equal(iv, want.iv) || !bytes.Equal(sn, want.sn) {
			t.Fatalf("suite %04x caller-destination traffic keys differ", suiteID)
		}
		record, err := newRecordCipher(suite, secret, 3, 64)
		if err != nil {
			t.Fatal(err)
		}
		if record.ivLen != suite.ivLen || !bytes.Equal(record.iv[:record.ivLen], want.iv) {
			t.Fatalf("suite %04x record IV differs after temporary material cleanup", suiteID)
		}
	}
}

func TestCCMUsageLimits(t *testing.T) {
	suite, err := cipherSuiteForID(TLS_AES_128_CCM_SHA256)
	if err != nil {
		t.Fatal(err)
	}
	if suite.recordLimit != 1<<23 || suite.authFailureLimit != 1<<23 {
		t.Fatalf("CCM limits = %d/%d", suite.recordLimit, suite.authFailureLimit)
	}
	if _, err = cipherSuiteForID(TLS_AES_128_CCM_8_SHA256); err == nil {
		t.Fatal("enabled CCM_8 without deployment-specific forgery safeguards")
	}
}
