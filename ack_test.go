package dtls13

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"testing"
)

func TestACKRoundTrip(t *testing.T) {
	want := []recordNumber{{epoch: 2, sequence: 7}, {epoch: 1<<63 + 9, sequence: 1<<63 + 11}}
	b, err := marshalACK(want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseACK(b)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func TestPlainACKRecordBuilderMatchesComposedMarshal(t *testing.T) {
	for _, numbers := range [][]recordNumber{
		nil,
		{{epoch: 2, sequence: 7}},
		{{epoch: 2, sequence: 7}, {epoch: 3, sequence: 9}},
	} {
		body, err := marshalCanonicalACK(numbers)
		if err != nil {
			t.Fatal(err)
		}
		want, err := marshalPlainRecord(record{typ: recordTypeACK, sequence: 11, payload: body})
		if err != nil {
			t.Fatal(err)
		}
		got, err := marshalPlainACKRecord(numbers, 11)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("direct=%x composed=%x", got, want)
		}
	}
	tooMany := make([]recordNumber, 4096)
	if wire, err := marshalPlainACKRecord(tooMany, 0); err == nil || wire != nil {
		t.Fatalf("oversized ACK wire=%x err=%v", wire, err)
	}
}

func TestACKParserReusesDestination(t *testing.T) {
	wire, err := marshalACK([]recordNumber{{epoch: 2, sequence: 7}})
	if err != nil {
		t.Fatal(err)
	}
	var storage [1]recordNumber
	got, err := parseACKInto(wire, storage[:0])
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != (recordNumber{epoch: 2, sequence: 7}) || &got[0] != &storage[0] {
		t.Fatalf("parser did not reuse destination: got=%#v", got)
	}

	largeWire, err := marshalACK([]recordNumber{{epoch: 2, sequence: 7}, {epoch: 2, sequence: 8}})
	if err != nil {
		t.Fatal(err)
	}
	got, err = parseACKInto(largeWire, storage[:0])
	if err != nil || len(got) != 2 || &got[0] == &storage[0] {
		t.Fatalf("parser fallback got=%#v err=%v", got, err)
	}
}

func TestACKRejectsMalformedVector(t *testing.T) {
	for _, b := range [][]byte{{}, {0}, {0, 8}, {0, 1, 0}} {
		if _, err := parseACK(b); err == nil {
			t.Fatalf("accepted malformed ACK %x", b)
		}
	}
}

func TestACKRejectsUnsortedOrDuplicateRecordNumbers(t *testing.T) {
	for _, numbers := range [][]recordNumber{
		{{epoch: 2, sequence: 1}, {epoch: 1, sequence: 2}},
		{{epoch: 2, sequence: 2}, {epoch: 2, sequence: 1}},
		{{epoch: 2, sequence: 1}, {epoch: 2, sequence: 1}},
	} {
		body := make([]byte, 2+16*len(numbers))
		binary.BigEndian.PutUint16(body, uint16(16*len(numbers)))
		for i, number := range numbers {
			off := 2 + 16*i
			binary.BigEndian.PutUint64(body[off:off+8], number.epoch)
			binary.BigEndian.PutUint64(body[off+8:off+16], number.sequence)
		}
		if _, err := parseACK(body); err == nil {
			t.Fatalf("accepted non-increasing ACK record numbers: %#v", numbers)
		}
	}
}

func TestACKRejectsAcknowledgementFromLowerEpoch(t *testing.T) {
	if err := validateACKEpoch([]recordNumber{{epoch: 3, sequence: 1}}, 2); err == nil {
		t.Fatal("accepted an ACK from an epoch below the acknowledged record")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("alert=%d ok=%v err=%v", description, ok, err)
	}
	if err := validateACKEpoch([]recordNumber{{epoch: 2}, {epoch: 3}}, 3); err != nil {
		t.Fatal(err)
	}
}

func TestACKCanonicalSortAndDeduplicate(t *testing.T) {
	input := []recordNumber{{epoch: 2, sequence: 3}, {epoch: 1, sequence: 9}, {epoch: 2, sequence: 3}}
	wantInput := append([]recordNumber(nil), input...)
	b, err := marshalACK(input)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input, wantInput) {
		t.Fatalf("marshalACK modified input: got %#v, want %#v", input, wantInput)
	}
	got, err := parseACK(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != (recordNumber{epoch: 1, sequence: 9}) || got[1] != (recordNumber{epoch: 2, sequence: 3}) {
		t.Fatalf("got %#v", got)
	}
}
func TestBuildACKRecordsPlainAndProtected(t *testing.T) {
	numbers := []recordNumber{{epoch: 2, sequence: 2}, {epoch: 2, sequence: 1}, {epoch: 2, sequence: 3}}
	plain, next, err := buildACKRecords(numbers, 50, 7, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plain) != 2 || next != 9 {
		t.Fatalf("records=%d next=%d", len(plain), next)
	}
	for _, wire := range plain {
		parsed, err := parsePlainRecords(wire)
		if err != nil {
			t.Fatal(err)
		}
		if parsed[0].typ != recordTypeACK {
			t.Fatal("wrong plaintext ACK type")
		}
	}
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	protected, _, err := buildACKRecords(numbers, 60, 0, sender)
	if err != nil {
		t.Fatal(err)
	}
	for _, wire := range protected {
		content, typ, _, err := receiver.open(wire)
		if err != nil {
			t.Fatal(err)
		}
		if typ != recordTypeACK {
			t.Fatal("wrong protected ACK type")
		}
		if _, err = parseACK(content); err != nil {
			t.Fatal(err)
		}
	}
}

func TestACKBuilderReusesDestination(t *testing.T) {
	var storage [1][]byte
	records, next, err := buildACKRecordsInto(storage[:0], []recordNumber{{epoch: 2, sequence: 1}}, 50, 7, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || next != 8 || &records[0] != &storage[0] {
		t.Fatalf("builder did not reuse destination: records=%d next=%d", len(records), next)
	}

	numbers := []recordNumber{{epoch: 2, sequence: 1}, {epoch: 2, sequence: 2}, {epoch: 2, sequence: 3}}
	records, next, err = buildACKRecordsInto(storage[:0], numbers, 50, 9, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 || next != 11 || &records[0] == &storage[0] {
		t.Fatalf("builder fallback records=%d next=%d", len(records), next)
	}
	for _, wire := range records {
		parsed, parseErr := parsePlainRecords(wire)
		if parseErr != nil || len(parsed) != 1 || parsed[0].typ != recordTypeACK {
			t.Fatalf("parse fallback ACK: records=%v err=%v", parsed, parseErr)
		}
	}
}

func TestBuildEmptyACKRecordsPlainAndProtected(t *testing.T) {
	plain, next, err := buildACKRecords(nil, 50, 7, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(plain) != 1 || next != 8 {
		t.Fatalf("plain empty ACK records=%d next=%d", len(plain), next)
	}
	records, err := parsePlainRecords(plain[0])
	if err != nil || len(records) != 1 || records[0].typ != recordTypeACK {
		t.Fatalf("parse plain empty ACK: records=%v err=%v", records, err)
	}
	if numbers, err := parseACK(records[0].payload); err != nil || len(numbers) != 0 {
		t.Fatalf("plain empty ACK payload: numbers=%v err=%v", numbers, err)
	}

	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	protected, _, err := buildACKRecords(nil, 60, 0, sender)
	if err != nil {
		t.Fatal(err)
	}
	if len(protected) != 1 {
		t.Fatalf("protected empty ACK records=%d", len(protected))
	}
	content, typ, _, err := receiver.open(protected[0])
	if err != nil || typ != recordTypeACK {
		t.Fatalf("open protected empty ACK: type=%d err=%v", typ, err)
	}
	if numbers, err := parseACK(content); err != nil || len(numbers) != 0 {
		t.Fatalf("protected empty ACK payload: numbers=%v err=%v", numbers, err)
	}
}

func TestACKRecordsRespectRecordLimitWithLargeMTU(t *testing.T) {
	numbers := make([]recordNumber, 2048)
	for i := range numbers {
		numbers[i] = recordNumber{epoch: 3, sequence: uint64(i)}
	}
	records, _, err := buildACKRecords(numbers, 65535, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) < 2 {
		t.Fatalf("large ACK used %d records", len(records))
	}
	for _, wire := range records {
		parsed, parseErr := parsePlainRecords(wire)
		if parseErr != nil || len(parsed) != 1 || len(parsed[0].payload) > maxRecordContent {
			t.Fatalf("ACK payload length=%d err=%v", len(parsed[0].payload), parseErr)
		}
	}
}

func TestReceiveFinalFlightACKAtApplicationEpoch(t *testing.T) {
	sender3, receiver3 := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	_, receiver2 := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	wire, _, err := buildACKRecords([]recordNumber{{epoch: 2, sequence: 5}}, 1200, 0, sender3)
	if err != nil {
		t.Fatal(err)
	}
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	go func() { _, _ = left.Write(wire[0]) }()
	numbers, err := receiveACKRecord(right, nil, receiver2, receiver3)
	if err != nil {
		t.Fatal(err)
	}
	want := recordNumber{epoch: 2, sequence: 5}
	if len(numbers) != 1 || numbers[0] != want {
		t.Fatalf("ACK numbers = %#v, want %#v", numbers, []recordNumber{want})
	}
}

func TestReceiveACKRecordAcceptsEarlierCipherCandidate(t *testing.T) {
	sender2, receiver2 := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	_, receiver3 := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	wire, _, err := buildACKRecords([]recordNumber{{epoch: 2, sequence: 5}}, 1200, 0, sender2)
	if err != nil {
		t.Fatal(err)
	}
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	go func() { _, _ = left.Write(wire[0]) }()
	var storage [1]recordNumber
	numbers, err := receiveACKRecord(right, storage[:0], receiver2, receiver3)
	if err != nil {
		t.Fatal(err)
	}
	want := recordNumber{epoch: 2, sequence: 5}
	if len(numbers) != 1 || numbers[0] != want {
		t.Fatalf("ACK numbers = %#v, want %#v", numbers, []recordNumber{want})
	}
}

func TestReceiveACKRecordSkipsNonUnifiedDatagram(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	want := recordNumber{epoch: 2, sequence: 5}
	plain, _, err := buildACKRecords([]recordNumber{want}, 1200, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	protected, _, err := buildACKRecords([]recordNumber{want}, 1200, 0, sender)
	if err != nil {
		t.Fatal(err)
	}
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	go func() {
		_, _ = left.Write(plain[0])
		_, _ = left.Write(protected[0])
	}()
	numbers, err := receiveACKRecord(right, nil, receiver)
	if err != nil {
		t.Fatal(err)
	}
	if len(numbers) != 1 || numbers[0] != want {
		t.Fatalf("ACK numbers = %#v, want %#v", numbers, []recordNumber{want})
	}
}

func TestReceiveACKRecordTriesCollidingEpochBitsNonDestructively(t *testing.T) {
	sender6, receiver6 := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 6)
	suite, err := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	if err != nil {
		t.Fatal(err)
	}
	receiver2, err := newRecordCipher(suite, bytes.Repeat([]byte{0x4a}, suite.hash.Size()), 2, 64)
	if err != nil {
		t.Fatal(err)
	}
	want := recordNumber{epoch: 6, sequence: 9}
	wire, _, err := buildACKRecords([]recordNumber{want}, 1200, 0, sender6)
	if err != nil {
		t.Fatal(err)
	}
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	go func() { _, _ = left.Write(wire[0]) }()
	numbers, err := receiveACKRecord(right, nil, receiver2, receiver6)
	if err != nil {
		t.Fatal(err)
	}
	if len(numbers) != 1 || numbers[0] != want {
		t.Fatalf("ACK numbers = %#v, want %#v", numbers, []recordNumber{want})
	}
	if receiver2.authFailures != 1 {
		t.Fatalf("colliding old epoch attempts = %d, want 1", receiver2.authFailures)
	}
}

func TestReceiveACKPropagatesAEADAuthenticationFailureLimit(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	receiver.authFailureLimit = 1
	wire, err := sender.seal(recordTypeACK, []byte{0, 0})
	if err != nil {
		t.Fatal(err)
	}
	wire[len(wire)-1] ^= 1
	go func() { _, _ = left.Write(wire) }()
	_, err = receiveACKRecord(right, nil, receiver)
	if !errors.Is(err, errAEADAuthenticationFailureLimit) {
		t.Fatalf("ACK receive returned %v", err)
	}
}
