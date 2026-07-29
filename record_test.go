package dtls13

import (
	"bytes"
	"testing"
)

func TestPlainRecordRoundTrip(t *testing.T) {
	want := record{typ: recordTypeHandshake, epoch: 0, sequence: 0x010203040506, payload: []byte("hello")}
	b, err := marshalPlainRecord(want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := parsePlainRecords(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].typ != want.typ || got[0].epoch != want.epoch || got[0].sequence != want.sequence || !bytes.Equal(got[0].payload, want.payload) {
		t.Fatalf("round trip mismatch: %#v", got)
	}
}

func TestPlainRecordAllocatorMatchesMarshal(t *testing.T) {
	for _, want := range []record{
		{typ: recordTypeAlert},
		{typ: recordTypeHandshake, sequence: 0x010203040506, payload: []byte("handshake")},
		{typ: recordTypeACK, sequence: 9, payload: []byte("ack")},
	} {
		allocated, err := allocatePlainRecordWire(want.typ, want.epoch, want.sequence, len(want.payload))
		if err != nil {
			t.Fatal(err)
		}
		copy(allocated[plainRecordHeaderLen:], want.payload)
		marshaled, err := marshalPlainRecord(want)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(allocated, marshaled) {
			t.Fatalf("allocated=%x marshaled=%x", allocated, marshaled)
		}
	}
}

func TestPlainRecordAllocatorMatchesMarshalErrors(t *testing.T) {
	tests := []struct {
		record     record
		payloadLen int
	}{
		{record: record{typ: recordTypeApplicationData}},
		{record: record{typ: recordTypeHandshake, epoch: 1}},
		{record: record{typ: recordTypeHandshake, sequence: 1 << 48}},
		{record: record{typ: recordTypeHandshake, payload: make([]byte, maxRecordContent+1)}, payloadLen: maxRecordContent + 1},
	}
	for _, test := range tests {
		if test.payloadLen == 0 {
			test.payloadLen = len(test.record.payload)
		}
		_, marshalErr := marshalPlainRecord(test.record)
		_, allocateErr := allocatePlainRecordWire(test.record.typ, test.record.epoch, test.record.sequence, test.payloadLen)
		if marshalErr == nil || allocateErr == nil || marshalErr.Error() != allocateErr.Error() {
			t.Fatalf("record=%#v marshal=%v allocate=%v", test.record, marshalErr, allocateErr)
		}
	}
	if wire, err := allocatePlainRecordWire(recordTypeHandshake, 0, 0, -1); err == nil || wire != nil {
		t.Fatalf("negative payload length wire=%x err=%v", wire, err)
	}
}

func TestPlainRecordParserOwnershipModes(t *testing.T) {
	wire, err := marshalPlainRecord(record{typ: recordTypeHandshake, payload: []byte("payload")})
	if err != nil {
		t.Fatal(err)
	}
	records, err := parsePlainRecords(wire)
	if err != nil {
		t.Fatal(err)
	}
	records[0].payload[0] ^= 0xff
	if wire[plainRecordHeaderLen] != 'p' {
		t.Fatal("parsePlainRecords payload aliases input")
	}
	records, err = parsePlainRecordsView(wire)
	if err != nil {
		t.Fatal(err)
	}
	records[0].payload[0] ^= 0xff
	if wire[plainRecordHeaderLen] != records[0].payload[0] {
		t.Fatal("parsePlainRecordsView did not reuse input")
	}
	var storage [1]record
	records, err = parsePlainRecordsViewInto(wire, storage[:0])
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || &records[0] != &storage[0] || &records[0].payload[0] != &wire[plainRecordHeaderLen] {
		t.Fatal("parsePlainRecordsViewInto did not reuse destination and input")
	}
}
func TestPlainRecordRejectsTruncation(t *testing.T) {
	b, _ := marshalPlainRecord(record{typ: recordTypeHandshake, payload: []byte("hello")})
	for i := 1; i < len(b); i++ {
		if _, err := parsePlainRecords(b[:i]); err == nil {
			t.Fatalf("accepted truncation at %d", i)
		}
	}
}

func TestPlainRecordRequiresEpochZero(t *testing.T) {
	if _, err := marshalPlainRecord(record{typ: recordTypeHandshake, epoch: 1}); err == nil {
		t.Fatal("marshaled nonzero plaintext epoch")
	}
	b, _ := marshalPlainRecord(record{typ: recordTypeHandshake})
	b[3], b[4] = 0, 1
	if _, err := parsePlainRecords(b); err == nil {
		t.Fatal("parsed nonzero plaintext epoch")
	}
}

func TestPlainRecordIgnoresLegacyVersion(t *testing.T) {
	for _, version := range []uint16{0xfeff, 0xfefd, 0x1234} {
		b, err := marshalPlainRecord(record{typ: recordTypeHandshake, payload: []byte("hello")})
		if err != nil {
			t.Fatal(err)
		}
		b[1], b[2] = byte(version>>8), byte(version)
		records, err := parsePlainRecords(b)
		if err != nil || len(records) != 1 || string(records[0].payload) != "hello" {
			t.Fatalf("legacy version %04x: records=%v err=%v", version, records, err)
		}
	}
}

func TestPlainRecordRejectsProtectedOnlyContentTypes(t *testing.T) {
	for _, contentType := range []uint8{recordTypeChangeCipherSpec, recordTypeApplicationData, recordTypeHeartbeat, 0xff} {
		if _, err := marshalPlainRecord(record{typ: contentType}); err == nil {
			t.Fatalf("marshaled plaintext content type %d", contentType)
		}
		wire := make([]byte, plainRecordHeaderLen)
		wire[0] = contentType
		if _, err := parsePlainRecords(wire); err == nil {
			t.Fatalf("parsed plaintext content type %d", contentType)
		}
	}
}

func TestReplayWindow(t *testing.T) {
	w := newReplayWindow(64)
	for _, seq := range []uint64{5, 4, 70} {
		if !w.check(seq) {
			t.Fatalf("unexpected rejection of %d", seq)
		}
		w.accept(seq)
	}
	for _, seq := range []uint64{5, 4, 0} {
		if w.check(seq) {
			t.Fatalf("accepted duplicate or stale %d", seq)
		}
	}
	if !w.check(69) {
		t.Fatal("rejected unseen in-window record")
	}
}
