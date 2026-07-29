package dtls13

import (
	"bytes"
	"reflect"
	"testing"
)

func TestACKGatedTrafficUpdate(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	secret := bytes.Repeat([]byte{4}, suite.hash.Size())
	sender, err := newSendingTraffic(suite, secret, 3, 7, 64)
	if err != nil {
		t.Fatal(err)
	}
	receiver, err := newReceivingTraffic(suite, secret, 3, 64)
	if err != nil {
		t.Fatal(err)
	}
	sender.cipher.setPlaintextLimit(minRecordSizeLimit)
	receiver.setPlaintextLimit(minRecordSizeLimit)
	wire, number, err := sender.beginKeyUpdate(true)
	if err != nil {
		t.Fatal(err)
	}
	if sender.cipher.epoch != 3 {
		t.Fatal("sender changed epoch before ACK")
	}
	if sender.pendingCipher.plaintextLimit != minRecordSizeLimit {
		t.Fatal("sending KeyUpdate did not preserve record_size_limit")
	}
	content, typ, epoch, _, err := receiver.epochs.open(wire)
	if err != nil {
		t.Fatal(err)
	}
	if typ != recordTypeHandshake || epoch != 3 {
		t.Fatalf("type=%d epoch=%d", typ, epoch)
	}
	fragments, err := parseHandshakeFragments(content)
	if err != nil {
		t.Fatal(err)
	}
	message, updated, err := receiver.processKeyUpdate(fragments[0].messageSequence, fragments[0].body)
	if err != nil || !updated || !message.requestUpdate {
		t.Fatalf("message=%v err=%v", message, err)
	}
	if receiver.current != 4 {
		t.Fatal("receiver did not prepare epoch 4")
	}
	if receiver.epochs.ciphers[4].plaintextLimit != minRecordSizeLimit {
		t.Fatal("receiving KeyUpdate did not preserve record_size_limit")
	}
	retransmitted, retransmittedNumber, err := sender.retransmitKeyUpdate()
	if err != nil {
		t.Fatal(err)
	}
	if retransmittedNumber == number {
		t.Fatal("KeyUpdate retransmission reused its record number")
	}
	content, _, _, _, err = receiver.epochs.open(retransmitted)
	if err != nil {
		t.Fatal(err)
	}
	fragments, err = parseHandshakeFragments(content)
	if err != nil {
		t.Fatal(err)
	}
	_, updated, err = receiver.processKeyUpdate(fragments[0].messageSequence, fragments[0].body)
	if err != nil || updated || receiver.current != 4 {
		t.Fatalf("duplicate KeyUpdate updated=%v epoch=%d err=%v", updated, receiver.current, err)
	}
	if sender.processACK([]recordNumber{{epoch: 3, sequence: retransmittedNumber.sequence + 100}}) {
		t.Fatal("accepted unrelated ACK")
	}
	if !sender.processACK([]recordNumber{retransmittedNumber}) || sender.cipher.epoch != 4 {
		t.Fatal("sender did not activate epoch 4")
	}
	app, err := sender.cipher.seal(recordTypeApplicationData, []byte("new epoch"))
	if err != nil {
		t.Fatal(err)
	}
	plain, _, epoch, _, err := receiver.epochs.open(app)
	if err != nil {
		t.Fatal(err)
	}
	if epoch != 4 || string(plain) != "new epoch" {
		t.Fatalf("epoch=%d plain=%q", epoch, plain)
	}
}

func TestSendingTrafficReusesAndClearsKeyUpdateStorage(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		initial := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		sender, err := newSendingTraffic(suite, initial, 3, 0, 64)
		if err != nil {
			t.Fatal(err)
		}
		firstSecretSlot := sender.secret
		wire, number, err := sender.beginKeyUpdate(false)
		if err != nil || len(wire) == 0 {
			t.Fatalf("suite %04x first KeyUpdate: wire=%d err=%v", suiteID, len(wire), err)
		}
		secondSecretSlot := sender.nextSecret
		firstFragment := &sender.pendingFragment[0]
		firstRecordSlot := &sender.update.records[0]
		if !sender.processACK([]recordNumber{number}) {
			t.Fatalf("suite %04x did not ACK first KeyUpdate", suiteID)
		}
		if &sender.secret[0] != &secondSecretSlot[0] || &sender.nextSecret[0] != &firstSecretSlot[0] {
			t.Fatalf("suite %04x did not rotate secret slots", suiteID)
		}
		if !bytes.Equal(firstSecretSlot, make([]byte, suite.hash.Size())) {
			t.Fatalf("suite %04x did not clear prior sending secret", suiteID)
		}

		wire, number, err = sender.beginKeyUpdate(true)
		if err != nil || len(wire) == 0 {
			t.Fatalf("suite %04x second KeyUpdate: wire=%d err=%v", suiteID, len(wire), err)
		}
		if &sender.nextSecret[0] != &firstSecretSlot[0] {
			t.Fatalf("suite %04x did not reuse prior secret storage", suiteID)
		}
		if &sender.pendingFragment[0] != firstFragment {
			t.Fatalf("suite %04x did not reuse KeyUpdate fragment storage", suiteID)
		}
		if &sender.update.records[0] != firstRecordSlot {
			t.Fatalf("suite %04x did not reuse ACK record storage", suiteID)
		}
		if !sender.processACK([]recordNumber{number}) {
			t.Fatalf("suite %04x did not ACK second KeyUpdate", suiteID)
		}
		if !bytes.Equal(secondSecretSlot, make([]byte, suite.hash.Size())) {
			t.Fatalf("suite %04x did not clear second prior sending secret", suiteID)
		}
		sender.clearSecrets()
		if !bytes.Equal(firstSecretSlot, make([]byte, suite.hash.Size())) {
			t.Fatalf("suite %04x did not clear current sending secret", suiteID)
		}
	}
}

func TestApplicationTrafficSecretsShareCapacityIsolatedStorage(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		for _, isClient := range []bool{false, true} {
			clientSecret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
			serverSecret := bytes.Repeat([]byte{0xa5}, suite.hash.Size())
			expectedSend := clientSecret
			expectedReceive := serverSecret
			if !isClient {
				expectedSend, expectedReceive = serverSecret, clientSecret
			}
			expectedSend = bytes.Clone(expectedSend)
			expectedReceive = bytes.Clone(expectedReceive)
			conn := &Conn{config: &Config{ReplayWindow: 64}, isClient: isClient}
			if err = conn.installApplicationKeys(suite, clientSecret, serverSecret); err != nil {
				t.Fatal(err)
			}
			sendingSecret := conn.sendingTraffic.secret
			receivingSecret := conn.receivingTraffic.secret
			clientSecret[0] ^= 0xff
			serverSecret[0] ^= 0xff
			if !bytes.Equal(sendingSecret, expectedSend) || !bytes.Equal(receivingSecret, expectedReceive) {
				t.Fatalf("suite %04x client=%v traffic secrets alias caller storage", suiteID, isClient)
			}
			if cap(sendingSecret) != suite.hash.Size() || cap(receivingSecret) != suite.hash.Size() {
				t.Fatalf("suite %04x client=%v traffic secret capacity is not isolated", suiteID, isClient)
			}
			sendingAddress := reflect.ValueOf(&sendingSecret[0]).Pointer()
			receivingAddress := reflect.ValueOf(&receivingSecret[0]).Pointer()
			if receivingAddress-sendingAddress != uintptr(suite.hash.Size()) {
				t.Fatalf("suite %04x client=%v traffic secrets do not share adjacent storage", suiteID, isClient)
			}
			conn.sendingTraffic.clearSecrets()
			if !bytes.Equal(receivingSecret, expectedReceive) {
				t.Fatalf("suite %04x client=%v clearing sending secret changed receive window", suiteID, isClient)
			}
			conn.receivingTraffic.clearSecrets()
			if !bytes.Equal(receivingSecret, make([]byte, suite.hash.Size())) {
				t.Fatalf("suite %04x client=%v did not clear receive window", suiteID, isClient)
			}
		}
	}
}

func TestReceivingTrafficBoundsRetainedEpochs(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	secret := bytes.Repeat([]byte{7}, suite.hash.Size())
	receiver, err := newReceivingTraffic(suite, secret, 3, 64)
	if err != nil {
		t.Fatal(err)
	}
	for sequence := uint16(10); sequence < 30; sequence++ {
		if _, updated, updateErr := receiver.processKeyUpdate(sequence, []byte{0}); updateErr != nil || !updated {
			t.Fatalf("sequence %d updated=%v err=%v", sequence, updated, updateErr)
		}
		receiver.epochs.mu.RLock()
		retained := len(receiver.epochs.ciphers)
		receiver.epochs.mu.RUnlock()
		if retained > 2 {
			t.Fatalf("retained %d epoch ciphers after sequence %d", retained, sequence)
		}
		if len(receiver.secrets) > 2 {
			t.Fatalf("retained %d epoch secrets after sequence %d", len(receiver.secrets), sequence)
		}
	}
}

func TestReceivingTrafficSharesAndClearsEpochSecretStorage(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		initial := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		receiver, err := newReceivingTraffic(suite, initial, 3, 64)
		if err != nil {
			t.Fatal(err)
		}
		storedInitial := receiver.secrets[3]
		if cap(receiver.secret) != suite.hash.Size() || &receiver.secret[0] != &storedInitial[0] {
			t.Fatalf("suite %04x current and epoch secrets do not share isolated storage", suiteID)
		}
		if _, updated, err := receiver.processKeyUpdate(1, []byte{0}); err != nil || !updated {
			t.Fatalf("suite %04x first update: updated=%v err=%v", suiteID, updated, err)
		}
		if current := receiver.secrets[4]; cap(receiver.secret) != suite.hash.Size() || &receiver.secret[0] != &current[0] {
			t.Fatalf("suite %04x updated current and epoch secrets do not share isolated storage", suiteID)
		}
		if _, updated, err := receiver.processKeyUpdate(2, []byte{0}); err != nil || !updated {
			t.Fatalf("suite %04x second update: updated=%v err=%v", suiteID, updated, err)
		}
		if _, retained := receiver.secrets[3]; retained {
			t.Fatalf("suite %04x retained an expired epoch secret", suiteID)
		}
		if !bytes.Equal(storedInitial, make([]byte, suite.hash.Size())) {
			t.Fatalf("suite %04x did not clear the expired epoch secret", suiteID)
		}
	}
}

func TestKeyUpdateSequenceMaySkipOtherPostHandshakeMessages(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	receiver, err := newReceivingTraffic(suite, bytes.Repeat([]byte{8}, suite.hash.Size()), 3, 64)
	if err != nil {
		t.Fatal(err)
	}
	if _, updated, updateErr := receiver.processKeyUpdate(10, []byte{0}); updateErr != nil || !updated {
		t.Fatalf("first update: updated=%v err=%v", updated, updateErr)
	}
	if _, updated, updateErr := receiver.processKeyUpdate(12, []byte{0}); updateErr != nil || !updated {
		t.Fatalf("update after another post-handshake message: updated=%v err=%v", updated, updateErr)
	}
	if _, updated, updateErr := receiver.processKeyUpdate(11, []byte{0}); updateErr != nil || updated {
		t.Fatalf("old update was not treated as a duplicate: updated=%v err=%v", updated, updateErr)
	}
}

func TestSendingEpochLimitDoesNotConstrainReceiver(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	secret := bytes.Repeat([]byte{9}, suite.hash.Size())
	sender, err := newSendingTraffic(suite, secret, maxSendingEpoch, 0, 64)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = sender.beginKeyUpdate(false); err == nil {
		t.Fatal("sender exceeded the 2^48-1 epoch limit")
	}
	if sender.canBeginKeyUpdate() {
		t.Fatal("sender advertised a KeyUpdate response at the 2^48-1 epoch limit")
	}
	receiver, err := newReceivingTraffic(suite, secret, maxSendingEpoch, 64)
	if err != nil {
		t.Fatal(err)
	}
	if _, updated, updateErr := receiver.processKeyUpdate(0, []byte{0}); updateErr != nil || !updated || receiver.current != maxSendingEpoch+1 {
		t.Fatalf("receiver incorrectly enforced sender epoch limit: updated=%v epoch=%d err=%v", updated, receiver.current, updateErr)
	}
}

func TestSendingTrafficRejectsHandshakeMessageSequenceWrap(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	sender, err := newSendingTraffic(suite, bytes.Repeat([]byte{10}, suite.hash.Size()), 3, ^uint16(0), 64)
	if err != nil {
		t.Fatal(err)
	}
	_, number, err := sender.beginKeyUpdate(false)
	if err != nil {
		t.Fatalf("last message sequence was rejected: %v", err)
	}
	if !sender.processACK([]recordNumber{number}) {
		t.Fatal("last-sequence KeyUpdate was not acknowledged")
	}
	if _, _, err = sender.beginKeyUpdate(false); err == nil {
		t.Fatal("handshake message sequence wrapped to zero")
	}
	if sender.messageSequence != ^uint16(0) || !sender.sequenceExhausted {
		t.Fatalf("sequence state = %d, exhausted=%v", sender.messageSequence, sender.sequenceExhausted)
	}
}
