package dtls13

import (
	"bytes"
	"testing"
)

func TestEpochSetSelectsCurrentAndPastKeys(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	oldSecret := bytes.Repeat([]byte{1}, suite.hash.Size())
	newSecret := bytes.Repeat([]byte{2}, suite.hash.Size())
	oldSend, _ := newRecordCipher(suite, oldSecret, 2, 64)
	oldRecv, _ := newRecordCipher(suite, oldSecret, 2, 64)
	newSend, _ := newRecordCipher(suite, newSecret, 3, 64)
	newRecv, _ := newRecordCipher(suite, newSecret, 3, 64)
	epochs := newEpochSet()
	if err := epochs.install(oldRecv); err != nil {
		t.Fatal(err)
	}
	if err := epochs.install(newRecv); err != nil {
		t.Fatal(err)
	}
	if err := epochs.setCurrent(3); err != nil {
		t.Fatal(err)
	}
	oldWire, _ := oldSend.seal(recordTypeApplicationData, []byte("old"))
	content, _, epoch, _, err := epochs.open(oldWire)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "old" || epoch != 2 {
		t.Fatalf("content=%q epoch=%d", content, epoch)
	}
	newWire, _ := newSend.seal(recordTypeApplicationData, []byte("new"))
	content, _, epoch, _, err = epochs.open(newWire)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "new" || epoch != 3 {
		t.Fatalf("content=%q epoch=%d", content, epoch)
	}
}
func TestEpochSetDiscardAndMonotonicCurrent(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	secret := bytes.Repeat([]byte{1}, suite.hash.Size())
	two, _ := newRecordCipher(suite, secret, 2, 64)
	three, _ := newRecordCipher(suite, secret, 3, 64)
	epochs := newEpochSet()
	_ = epochs.install(two)
	_ = epochs.install(three)
	if err := epochs.setCurrent(3); err != nil {
		t.Fatal(err)
	}
	if err := epochs.setCurrent(2); err == nil {
		t.Fatal("moved epoch backwards")
	}
	epochs.discardBefore(3)
	if _, err := epochs.selectCipher(byte(unifiedHeaderFixed | 2)); err == nil {
		t.Fatal("selected discarded epoch")
	}
}

func TestEpochLowBitsCollisionSelectsCurrentEpoch(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	secret := bytes.Repeat([]byte{3}, suite.hash.Size())
	oldSender, _ := newRecordCipher(suite, secret, 3, 64)
	oldReceiver, _ := newRecordCipher(suite, secret, 3, 64)
	currentSender, _ := newRecordCipher(suite, secret, 7, 64)
	currentReceiver, _ := newRecordCipher(suite, secret, 7, 64)
	epochs := newEpochSet()
	_ = epochs.install(oldReceiver)
	_ = epochs.install(currentReceiver)
	if err := epochs.setCurrent(7); err != nil {
		t.Fatal(err)
	}
	currentWire, _ := currentSender.seal(recordTypeApplicationData, []byte("current"))
	content, _, epoch, _, err := epochs.open(currentWire)
	if err != nil || epoch != 7 || string(content) != "current" {
		t.Fatalf("current collision epoch=%d content=%q err=%v", epoch, content, err)
	}
	oldWire, _ := oldSender.seal(recordTypeApplicationData, []byte("old"))
	if _, _, _, _, err = epochs.open(oldWire); err == nil {
		t.Fatal("accepted an epoch outside the retained reconstruction range")
	}
}
