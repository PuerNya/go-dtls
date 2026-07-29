package dtls13

import (
	"bytes"
	"crypto/hmac"
	"testing"
)

func TestKeyScheduleFinishedAndApplicationSecrets(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	client := newKeySchedule(suite, nil)
	server := newKeySchedule(suite, nil)
	shared := bytes.Repeat([]byte{0x33}, 32)
	helloHash := bytes.Repeat([]byte{0x44}, suite.hash.Size())
	for _, k := range []*keySchedule{client, server} {
		if err := k.deriveHandshake(shared, helloHash); err != nil {
			t.Fatal(err)
		}
	}
	if !bytes.Equal(client.clientHandshakeTraffic, server.clientHandshakeTraffic) || !bytes.Equal(client.serverHandshakeTraffic, server.serverHandshakeTraffic) {
		t.Fatal("peers derived different handshake traffic secrets")
	}
	transcript := bytes.Repeat([]byte{0x55}, suite.hash.Size())
	verify := client.finishedVerifyData(client.clientHandshakeTraffic, transcript)
	if !server.verifyFinished(server.clientHandshakeTraffic, transcript, verify) {
		t.Fatal("valid Finished was rejected")
	}
	verify[0] ^= 1
	if server.verifyFinished(server.clientHandshakeTraffic, transcript, verify) {
		t.Fatal("tampered Finished was accepted")
	}
	derived := deriveSecret(suite, client.handshakeSecret, labelDerived, emptyTranscriptHash(suite))
	wantMaster := hkdfExtract(suite.hash.New, make([]byte, suite.hash.Size()), derived)
	wrongMaster := hkdfExtract(suite.hash.New, nil, derived)
	for _, k := range []*keySchedule{client, server} {
		if err := k.deriveApplication(transcript); err != nil {
			t.Fatal(err)
		}
	}
	if !bytes.Equal(client.masterSecret, wantMaster) {
		t.Fatal("master secret did not use a Hash.length zero input secret")
	}
	if bytes.Equal(client.masterSecret, wrongMaster) {
		t.Fatal("master secret incorrectly used an empty input secret")
	}
	if !bytes.Equal(client.clientApplicationTraffic, server.clientApplicationTraffic) || !bytes.Equal(client.serverApplicationTraffic, server.serverApplicationTraffic) {
		t.Fatal("peers derived different application traffic secrets")
	}
	if bytes.Equal(client.clientApplicationTraffic, client.serverApplicationTraffic) {
		t.Fatal("client and server traffic secrets are equal")
	}
}

func TestKeyScheduleTrafficSecretWindowsAreCapacityIsolated(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		schedule := newKeySchedule(suite, nil)
		transcript := bytes.Repeat([]byte{0x44}, suite.hash.Size())
		root := &schedule.earlySecret[0]
		if len(schedule.secretStorage) != keyScheduleSecretSlots*suite.hash.Size() || cap(schedule.earlySecret) != suite.hash.Size() {
			t.Fatalf("suite %04x key schedule storage is not exact or isolated", suiteID)
		}
		earlyBefore := append([]byte(nil), schedule.earlySecret...)
		earlyTraffic := schedule.earlyTrafficSecret(transcript)
		if cap(earlyTraffic) != suite.hash.Size() || &earlyTraffic[0] != &schedule.secretStorage[suite.hash.Size()] {
			t.Fatalf("suite %04x early traffic secret did not use its isolated storage window", suiteID)
		}
		if !bytes.Equal(schedule.earlySecret, earlyBefore) {
			t.Fatalf("suite %04x early traffic derivation modified the early secret", suiteID)
		}
		if err = schedule.deriveHandshake(bytes.Repeat([]byte{0x33}, suite.hash.Size()), transcript); err != nil {
			t.Fatal(err)
		}
		if schedule.earlySecret != nil || &schedule.handshakeSecret[0] != root {
			t.Fatalf("suite %04x did not reuse and retire the early-secret window", suiteID)
		}
		if cap(schedule.clientHandshakeTraffic) != suite.hash.Size() || cap(schedule.serverHandshakeTraffic) != suite.hash.Size() {
			t.Fatalf("suite %04x handshake secret capacities are not isolated", suiteID)
		}
		clientHandshakeBeforeApplication := append([]byte(nil), schedule.clientHandshakeTraffic...)
		serverHandshake := append([]byte(nil), schedule.serverHandshakeTraffic...)
		clientHandshake := append(schedule.clientHandshakeTraffic, 0xff)
		clientHandshake[0] ^= 0xff
		if !bytes.Equal(schedule.serverHandshakeTraffic, serverHandshake) {
			t.Fatalf("suite %04x appending client handshake secret modified server secret", suiteID)
		}

		if err = schedule.deriveApplication(transcript); err != nil {
			t.Fatal(err)
		}
		if schedule.handshakeSecret != nil || &schedule.masterSecret[0] != root {
			t.Fatalf("suite %04x did not reuse and retire the handshake-secret window", suiteID)
		}
		if !bytes.Equal(schedule.clientHandshakeTraffic, clientHandshakeBeforeApplication) || !bytes.Equal(schedule.serverHandshakeTraffic, serverHandshake) {
			t.Fatalf("suite %04x application derivation modified handshake traffic secrets", suiteID)
		}
		windows := [][]byte{schedule.clientApplicationTraffic, schedule.serverApplicationTraffic, schedule.exporterMasterSecret}
		for i, window := range windows {
			if cap(window) != suite.hash.Size() {
				t.Fatalf("suite %04x application secret %d capacity=%d", suiteID, i, cap(window))
			}
		}
		serverApplication := append([]byte(nil), schedule.serverApplicationTraffic...)
		exporter := append([]byte(nil), schedule.exporterMasterSecret...)
		clientApplication := append(schedule.clientApplicationTraffic, 0xff)
		clientApplication[0] ^= 0xff
		if !bytes.Equal(schedule.serverApplicationTraffic, serverApplication) || !bytes.Equal(schedule.exporterMasterSecret, exporter) {
			t.Fatalf("suite %04x appending client application secret modified adjacent secrets", suiteID)
		}
		clientApplicationBeforeResumption := append([]byte(nil), schedule.clientApplicationTraffic...)
		if err = schedule.deriveResumption(transcript); err != nil {
			t.Fatal(err)
		}
		if schedule.masterSecret != nil || &schedule.resumptionMasterSecret[0] != root {
			t.Fatalf("suite %04x did not reuse and retire the master-secret window", suiteID)
		}
		if !bytes.Equal(schedule.clientApplicationTraffic, clientApplicationBeforeResumption) || !bytes.Equal(schedule.serverApplicationTraffic, serverApplication) || !bytes.Equal(schedule.exporterMasterSecret, exporter) {
			t.Fatalf("suite %04x resumption derivation modified application secrets", suiteID)
		}
	}
}

func TestEmptyTranscriptHashMatchesSuiteAndRemainsImmutable(t *testing.T) {
	for _, suiteID := range defaultCipherSuites() {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		h := suite.hash.New()
		want := h.Sum(nil)
		first := emptyTranscriptHash(suite)
		second := emptyTranscriptHash(suite)
		if !bytes.Equal(first, want) || len(first) == 0 || &first[0] != &second[0] {
			t.Fatalf("suite %04x empty transcript digest is not the shared standard digest", suiteID)
		}
		before := append([]byte(nil), first...)
		_ = deriveSecret(suite, bytes.Repeat([]byte{0x5a}, suite.hash.Size()), labelDerived, first)
		if !bytes.Equal(first, before) {
			t.Fatalf("suite %04x derivation modified the shared empty transcript digest", suiteID)
		}
		zeros := zeroHashSecret(suite)
		secondZeros := zeroHashSecret(suite)
		if len(zeros) != suite.hash.Size() || !bytes.Equal(zeros, make([]byte, suite.hash.Size())) || &zeros[0] != &secondZeros[0] {
			t.Fatalf("suite %04x zero secret is not shared Hash.length zero storage", suiteID)
		}
		_ = hkdfExtract(suite.hash.New, zeros, first)
		if !bytes.Equal(zeros, make([]byte, suite.hash.Size())) {
			t.Fatalf("suite %04x HKDF-Extract modified the shared zero secret", suiteID)
		}
	}
}

func TestFinishedVerifyDataMatchesIndependentDerivation(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		transcript := bytes.Repeat([]byte{0xa5}, suite.hash.Size())
		key := expandLabel(suite, secret, "finished", nil, suite.hash.Size())
		mac := hmac.New(suite.hash.New, key)
		_, _ = mac.Write(transcript)
		want := mac.Sum(nil)
		if got := computeFinishedVerifyData(suite, secret, transcript); !bytes.Equal(got, want) {
			t.Fatalf("suite %04x Finished verify_data differs", suiteID)
		}
	}
}

func TestKeyScheduleRejectsInvalidTransitions(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	k := newKeySchedule(suite, nil)
	if err := k.deriveApplication(make([]byte, suite.hash.Size())); err == nil {
		t.Fatal("derived application keys before handshake keys")
	}
	if err := k.deriveHandshake(make([]byte, 32), make([]byte, suite.hash.Size())); err != nil {
		t.Fatal(err)
	}
	if err := k.deriveHandshake(make([]byte, 32), make([]byte, suite.hash.Size())); err == nil {
		t.Fatal("derived handshake keys twice")
	}
}

func TestKeyScheduleResumptionPSK(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	shared := bytes.Repeat([]byte{0x11}, 32)
	helloHash := bytes.Repeat([]byte{0x22}, suite.hash.Size())
	finishedHash := bytes.Repeat([]byte{0x33}, suite.hash.Size())
	peers := []*keySchedule{newKeySchedule(suite, nil), newKeySchedule(suite, nil)}
	var first []byte
	for _, schedule := range peers {
		if err := schedule.deriveHandshake(shared, helloHash); err != nil {
			t.Fatal(err)
		}
		if err := schedule.deriveApplication(finishedHash); err != nil {
			t.Fatal(err)
		}
		if err := schedule.deriveResumption(finishedHash); err != nil {
			t.Fatal(err)
		}
		psk, err := schedule.resumptionPSK([]byte{1, 2, 3})
		if err != nil {
			t.Fatal(err)
		}
		if first == nil {
			first = psk
		} else if !bytes.Equal(first, psk) {
			t.Fatal("peers derived different resumption PSKs")
		}
		other, err := schedule.resumptionPSK([]byte{1, 2, 4})
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(psk, other) {
			t.Fatal("ticket nonce did not separate resumption PSKs")
		}
	}
}

func TestExporterEnforcesHKDFOutputLimit(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		exporter := newExporter(suite, bytes.Repeat([]byte{0x5a}, suite.hash.Size()))
		limit := 255 * suite.hash.Size()
		if output, err := exporter.export("limit", nil, limit); err != nil || len(output) != limit {
			t.Fatalf("suite %04x maximum output length=%d err=%v", suiteID, len(output), err)
		}
		if _, err := exporter.export("limit", nil, limit+1); err == nil {
			t.Fatalf("suite %04x accepted output beyond HKDF limit", suiteID)
		}
		maxLabel := string(bytes.Repeat([]byte{'x'}, 255-len("dtls13")))
		if output, err := exporter.export(maxLabel, nil, suite.hash.Size()); err != nil || len(output) != suite.hash.Size() {
			t.Fatalf("suite %04x maximum label length=%d err=%v", suiteID, len(output), err)
		}
		if _, err := exporter.export(maxLabel+"x", nil, suite.hash.Size()); err == nil {
			t.Fatalf("suite %04x accepted label beyond uint8 length", suiteID)
		}
	}
}

func TestExporterMatchesStandardDerivation(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		context := bytes.Repeat([]byte{0xa5}, suite.hash.Size()+7)
		exporter := newExporter(suite, secret)
		for _, length := range []int{0, 1, suite.hash.Size() - 1, suite.hash.Size(), 2*suite.hash.Size() + 3, 255 * suite.hash.Size()} {
			labelSecret := expandLabel(suite, secret, "differential", emptyTranscriptHash(suite), suite.hash.Size())
			h := suite.hash.New()
			_, _ = h.Write(context)
			want := expandLabel(suite, labelSecret, "exporter", h.Sum(nil), length)
			got, err := exporter.export("differential", context, length)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("suite %04x length %d exporter differs from standard path", suiteID, length)
			}
		}
	}
}
