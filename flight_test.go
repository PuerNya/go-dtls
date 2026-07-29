package dtls13

import (
	"bytes"
	"context"
	"syscall"
	"testing"
	"time"
)

type mtuFlightWriter struct {
	limit   int
	records [][]byte
}

func (w *mtuFlightWriter) Write(p []byte) (int, error) {
	if len(p) > w.limit {
		return 0, syscall.Errno(10040)
	}
	w.records = append(w.records, append([]byte(nil), p...))
	return len(p), nil
}

func TestBuildPlainFlightFragmentsToMTU(t *testing.T) {
	body := bytes.Repeat([]byte{0x42}, 100)
	f, next, err := buildPlainFlight([]handshakeMessage{{typ: 1, sequence: 7, body: body}}, 50, 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if next != 14 || len(f.records) != 4 {
		t.Fatalf("next=%d records=%d", next, len(f.records))
	}
	r := newReassembler()
	for i, fr := range f.records {
		if len(fr.wire) > 50 {
			t.Fatalf("record %d exceeds MTU", i)
		}
		records, err := parsePlainRecords(fr.wire)
		if err != nil {
			t.Fatal(err)
		}
		fragments, err := parseHandshakeFragments(records[0].payload)
		if err != nil {
			t.Fatal(err)
		}
		got, done, err := r.add(fragments[0])
		if err != nil {
			t.Fatal(err)
		}
		if i == len(f.records)-1 {
			if !done || !bytes.Equal(got, body) {
				t.Fatal("message did not reassemble")
			}
		} else if done {
			t.Fatal("message completed early")
		}
	}
}

func TestHandshakeFragmenterReusesDestinationViews(t *testing.T) {
	body := []byte("abcdefgh")
	var storage [2]handshakeFragment
	fragments := fragmentHandshakeMessageInto(storage[:0], handshakeMessage{typ: 1, sequence: 7, body: body}, 4)
	if len(fragments) != 2 || &fragments[0] != &storage[0] {
		t.Fatalf("fragment destination not reused: %#v", fragments)
	}
	if &fragments[0].body[0] != &body[0] || &fragments[1].body[0] != &body[4] {
		t.Fatal("fragment bodies do not view the input message")
	}
	var shortStorage [1]handshakeFragment
	fragments = fragmentHandshakeMessageInto(shortStorage[:0], handshakeMessage{typ: 1, body: body}, 4)
	if len(fragments) != 2 || &fragments[0] == &shortStorage[0] {
		t.Fatalf("fragment allocation fallback failed: %#v", fragments)
	}
}

func TestCountHandshakeFragments(t *testing.T) {
	messages := []handshakeMessage{
		{typ: 1},
		{typ: 2, body: make([]byte, 1200)},
		{typ: 3, body: make([]byte, 1201)},
	}
	count, err := countHandshakeFragments(messages, 1200)
	if err != nil || count != 4 {
		t.Fatalf("fragment count=%d err=%v", count, err)
	}
	if _, err = countHandshakeFragments([]handshakeMessage{{typ: 1, body: make([]byte, 1<<24)}}, 1200); err == nil {
		t.Fatal("accepted a handshake message exceeding the 24-bit length field")
	}
}

func TestBuildProtectedFlightFragmentsToMTU(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	body := bytes.Repeat([]byte{9}, 100)
	f, err := buildProtectedFlight([]handshakeMessage{{typ: 11, sequence: 2, body: body}}, 60, sender)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.records) < 2 {
		t.Fatal("message was not fragmented")
	}
	r := newReassembler()
	for index, item := range f.records {
		if len(item.wire) > 60 {
			t.Fatal("record exceeds MTU")
		}
		content, typ, _, err := receiver.open(item.wire)
		if err != nil {
			t.Fatal(err)
		}
		if typ != recordTypeHandshake {
			t.Fatalf("type %d", typ)
		}
		fragments, err := parseHandshakeFragments(content)
		if err != nil {
			t.Fatal(err)
		}
		got, done, err := r.add(fragments[0])
		if err != nil {
			t.Fatal(err)
		}
		if index == len(f.records)-1 && (!done || !bytes.Equal(got, body)) {
			t.Fatal("protected message did not reassemble")
		}
	}
}

func TestProtectedFlightRetainsMessageHeadersAndPayload(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	body := []byte("original")
	messages := []handshakeMessage{{typ: handshakeTypeCertificate, body: body}}
	f, err := buildProtectedFlight(messages, 1200, sender)
	if err != nil {
		t.Fatal(err)
	}
	messages[0] = handshakeMessage{typ: handshakeTypeFinished, body: []byte("modified")}
	if err = f.refreshPending(); err != nil {
		t.Fatal(err)
	}
	content, typ, _, err := receiver.open(f.records[0].wire)
	if err != nil || typ != recordTypeHandshake {
		t.Fatalf("open retransmission: type=%d err=%v", typ, err)
	}
	fragments, err := parseHandshakeFragments(content)
	if err != nil || len(fragments) != 1 {
		t.Fatalf("parse retransmission: fragments=%d err=%v", len(fragments), err)
	}
	if string(fragments[0].body) != "original" {
		t.Fatalf("retransmission payload=%q", fragments[0].body)
	}
	if err = f.resize(1000); err != nil {
		t.Fatal(err)
	}
	content, typ, _, err = receiver.open(f.records[0].wire)
	if err != nil || typ != recordTypeHandshake {
		t.Fatalf("open resized record: type=%d err=%v", typ, err)
	}
	fragments, err = parseHandshakeFragments(content)
	if err != nil || len(fragments) != 1 || fragments[0].typ != handshakeTypeCertificate || string(fragments[0].body) != "original" {
		t.Fatalf("resized fragments=%#v err=%v", fragments, err)
	}
}

func TestPlainFlightRetainsSingleMessageHeader(t *testing.T) {
	messages := []handshakeMessage{{typ: handshakeTypeCertificate, body: []byte("original")}}
	f, _, err := buildPlainFlight(messages, 1200, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	messages[0] = handshakeMessage{typ: handshakeTypeFinished, body: []byte("modified")}
	if err = f.resize(1000); err != nil {
		t.Fatal(err)
	}
	records, err := parsePlainRecords(f.records[0].wire)
	if err != nil || len(records) != 1 {
		t.Fatalf("records=%d err=%v", len(records), err)
	}
	fragments, err := parseHandshakeFragments(records[0].payload)
	if err != nil || len(fragments) != 1 || fragments[0].typ != handshakeTypeCertificate || string(fragments[0].body) != "original" {
		t.Fatalf("fragments=%#v err=%v", fragments, err)
	}
}

func TestProtectedFlightRetainsMultipleMessageHeaders(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	messages := []handshakeMessage{
		{typ: handshakeTypeCertificate, body: []byte("first")},
		{typ: handshakeTypeFinished, sequence: 1, body: []byte("second")},
	}
	f, err := buildProtectedFlight(messages, 1200, sender)
	if err != nil {
		t.Fatal(err)
	}
	messages[0], messages[1] = messages[1], messages[0]
	if err = f.resize(1000); err != nil {
		t.Fatal(err)
	}
	for i, want := range []struct {
		typ  uint8
		body string
	}{{handshakeTypeCertificate, "first"}, {handshakeTypeFinished, "second"}} {
		content, typ, _, openErr := receiver.open(f.records[i].wire)
		if openErr != nil || typ != recordTypeHandshake {
			t.Fatalf("record %d type=%d err=%v", i, typ, openErr)
		}
		fragments, parseErr := parseHandshakeFragments(content)
		if parseErr != nil || len(fragments) != 1 || fragments[0].typ != want.typ || string(fragments[0].body) != want.body {
			t.Fatalf("record %d fragments=%#v err=%v", i, fragments, parseErr)
		}
	}
}

func TestCombinedFlightRefreshAndResizePreservePartOrder(t *testing.T) {
	first := &flight{records: []flightRecord{{wire: []byte("first")}}}
	second := &flight{records: []flightRecord{{wire: []byte("second")}}}
	first.rebuildPending = func(flightRecordIndices) error {
		first.records[0].wire = []byte("first refreshed")
		return nil
	}
	second.rebuildPending = func(flightRecordIndices) error {
		second.records[0].wire = []byte("second refreshed")
		return nil
	}
	first.resizeForMTU = func(int) error {
		first.records = []flightRecord{{wire: []byte("first resized 1")}, {wire: []byte("first resized 2")}}
		return nil
	}
	second.resizeForMTU = func(int) error {
		second.records = []flightRecord{{wire: []byte("second resized")}}
		return nil
	}
	combined := combineFlights(first, second)
	if got := []string{string(combined.records[0].wire), string(combined.records[1].wire)}; got[0] != "first" || got[1] != "second" {
		t.Fatalf("initial records=%q", got)
	}
	if err := combined.refreshPending(); err != nil {
		t.Fatal(err)
	}
	if got := []string{string(combined.records[0].wire), string(combined.records[1].wire)}; got[0] != "first refreshed" || got[1] != "second refreshed" {
		t.Fatalf("refreshed records=%q", got)
	}
	if err := combined.resize(900); err != nil {
		t.Fatal(err)
	}
	want := []string{"first resized 1", "first resized 2", "second resized"}
	if len(combined.records) != len(want) {
		t.Fatalf("resized records=%d, want %d", len(combined.records), len(want))
	}
	for i := range want {
		if got := string(combined.records[i].wire); got != want[i] {
			t.Fatalf("resized record %d=%q, want %q", i, got, want[i])
		}
	}
}

func TestFlightsRespectRecordLimitWithLargeMTU(t *testing.T) {
	message := handshakeMessage{typ: handshakeTypeCertificate, sequence: 1, body: make([]byte, 2*maxRecordContent)}
	plain, _, err := buildPlainFlight([]handshakeMessage{message}, 65535, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, record := range plain.records {
		parsed, parseErr := parsePlainRecords(record.wire)
		if parseErr != nil || len(parsed) != 1 || len(parsed[0].payload) > maxRecordContent {
			t.Fatalf("plaintext payload length=%d err=%v", len(parsed[0].payload), parseErr)
		}
	}
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	protected, err := buildProtectedFlight([]handshakeMessage{message}, 65535, sender)
	if err != nil {
		t.Fatal(err)
	}
	if len(protected.records) < 2 {
		t.Fatalf("large protected handshake used %d records", len(protected.records))
	}
}

func TestProtectedFlightRefragmentsAfterMTUError(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	body := bytes.Repeat([]byte{0x5a}, 2000)
	flight, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeCertificate, sequence: 4, body: body}}, 1200, sender)
	if err != nil {
		t.Fatal(err)
	}
	config, err := (&Config{MTU: 1200}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	conn := &Conn{config: config}
	writer := &mtuFlightWriter{limit: 500}
	if err = conn.writeFlight(writer, flight); err != nil {
		t.Fatal(err)
	}
	if conn.currentMTU() > writer.limit || len(writer.records) < 2 {
		t.Fatalf("path MTU=%d records=%d", conn.currentMTU(), len(writer.records))
	}
	reassembler := newReassembler()
	var got []byte
	for _, wire := range writer.records {
		if len(wire) > writer.limit {
			t.Fatalf("record length %d exceeds limit", len(wire))
		}
		content, typ, _, openErr := receiver.open(wire)
		if openErr != nil || typ != recordTypeHandshake {
			t.Fatalf("open refragmented record: type=%d err=%v", typ, openErr)
		}
		fragments, parseErr := parseHandshakeFragments(content)
		if parseErr != nil || len(fragments) != 1 {
			t.Fatalf("parse fragment: count=%d err=%v", len(fragments), parseErr)
		}
		if complete, done, addErr := reassembler.add(fragments[0]); addErr != nil {
			t.Fatal(addErr)
		} else if done {
			got = complete
		}
	}
	if !bytes.Equal(got, body) {
		t.Fatal("refragmented flight did not reassemble")
	}
}

func TestFlightBlackHoleTimeoutReducesMTU(t *testing.T) {
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	flight, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeCertificate, sequence: 1, body: bytes.Repeat([]byte{1}, 2000)}}, 1200, sender)
	if err != nil {
		t.Fatal(err)
	}
	config, err := (&Config{MTU: 1200}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	conn := &Conn{config: config}
	for timeout := 1; timeout <= 3; timeout++ {
		if _, err = conn.prepareFlightRetransmission(flight, timeout); err != nil {
			t.Fatal(err)
		}
	}
	if conn.currentMTU() != 900 {
		t.Fatalf("path MTU=%d, want 900", conn.currentMTU())
	}
	for _, record := range flight.records {
		if len(record.wire) > 900 {
			t.Fatalf("refragmented record length=%d", len(record.wire))
		}
	}
}

func TestProtectedFlightUsesTenRecordSlidingWindow(t *testing.T) {
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	flight, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeCertificate, sequence: 1, body: bytes.Repeat([]byte{1}, 6000)}}, 256, sender)
	if err != nil {
		t.Fatal(err)
	}
	if len(flight.records) <= 20 {
		t.Fatalf("test flight has only %d records", len(flight.records))
	}
	config, err := (&Config{MTU: 256}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	conn := &Conn{config: config}
	writer := &mtuFlightWriter{limit: 256}
	if err = conn.writeFlight(writer, flight); err != nil {
		t.Fatal(err)
	}
	if len(writer.records) != 10 {
		t.Fatalf("initial burst=%d, want 10", len(writer.records))
	}
	var firstFive []recordNumber
	for i := 0; i < 5; i++ {
		firstFive = append(firstFive, flight.records[i].number)
	}
	flight.ack(firstFive)
	if err = conn.writeFlight(writer, flight); err != nil {
		t.Fatal(err)
	}
	if len(writer.records) != 15 {
		t.Fatalf("burst after five ACKs=%d, want 15 total", len(writer.records))
	}
}

func TestFlightWindowMethodsReuseDestination(t *testing.T) {
	f := &flight{records: []flightRecord{{wire: []byte("record"), number: recordNumber{epoch: 2}}}}
	var indexStorage [10]int
	indices := f.pendingIndices(indexStorage[:0])
	if len(indices) != 1 || &indices[0] != &indexStorage[0] {
		t.Fatalf("pending indices did not reuse destination: %v", indices)
	}
	var wireStorage [10][]byte
	wires := f.pendingWire(wireStorage[:0])
	if len(wires) != 1 || &wires[0] != &wireStorage[0] || string(wires[0]) != "record" {
		t.Fatalf("pending wire did not reuse destination: %q", wires)
	}
	if &wires[0][0] != &f.records[0].wire[0] {
		t.Fatal("pending wire copied immutable flight storage")
	}
}

func TestFlightWireSnapshotSurvivesRefresh(t *testing.T) {
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	f, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeCertificate, body: bytes.Repeat([]byte("certificate"), 300)}}, 1200, sender)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.records) < 2 {
		t.Fatal("test message was not fragmented")
	}
	old := make([][]byte, len(f.records))
	want := make([][]byte, len(f.records))
	for i := range f.records {
		if cap(f.records[i].wire) != len(f.records[i].wire) {
			t.Fatalf("record %d len=%d cap=%d", i, len(f.records[i].wire), cap(f.records[i].wire))
		}
		old[i] = f.records[i].wire
		want[i] = append([]byte(nil), old[i]...)
	}
	if err = f.refreshPending(); err != nil {
		t.Fatal(err)
	}
	for i := range f.records {
		if !bytes.Equal(old[i], want[i]) {
			t.Fatalf("refresh modified published wire snapshot %d", i)
		}
		if &old[i][0] == &f.records[i].wire[0] {
			t.Fatalf("refresh reused published wire backing %d", i)
		}
		if cap(f.records[i].wire) != len(f.records[i].wire) {
			t.Fatalf("refreshed record %d len=%d cap=%d", i, len(f.records[i].wire), cap(f.records[i].wire))
		}
	}
}

func TestFlightRTTAdjustsAndResetsRetransmitTimer(t *testing.T) {
	now := time.Unix(1000, 0)
	config, err := (&Config{FlightInterval: time.Second, MaxFlightInterval: 60 * time.Second, Time: func() time.Time { return now }}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	conn := &Conn{config: config}
	flight, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeClientHello, body: []byte{1}}}, 1200, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	flight.noteSend(now, false)
	now = now.Add(100 * time.Millisecond)
	conn.observeFlightRTT(flight)
	if got := conn.flightInterval(); got != 150*time.Millisecond {
		t.Fatalf("adjusted interval=%v", got)
	}
	now = now.Add(1500 * time.Millisecond)
	if got := conn.flightInterval(); got != time.Second {
		t.Fatalf("idle reset interval=%v", got)
	}

	retransmitted, _, _ := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeClientHello, body: []byte{1}}}, 1200, 0, 0)
	retransmitted.noteSend(now, false)
	retransmitted.noteSend(now.Add(time.Millisecond), true)
	now = now.Add(50 * time.Millisecond)
	conn.observeFlightRTT(retransmitted)
	if got := conn.flightInterval(); got != time.Second {
		t.Fatalf("retransmitted flight changed interval to %v", got)
	}
}

func TestFlightRetransmitsOnlyUnackedRecords(t *testing.T) {
	f, _, err := buildPlainFlight([]handshakeMessage{{typ: 1, sequence: 0, body: bytes.Repeat([]byte{1}, 40)}}, 40, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.records) < 2 {
		t.Fatal("test requires multiple records")
	}
	f.setIntervals(5*time.Millisecond, 10*time.Millisecond)
	ackEvent := make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	counts := make([]int, len(f.records))
	oldRemaining := make([]recordNumber, len(f.records)-1)
	for i := 1; i < len(f.records); i++ {
		oldRemaining[i-1] = f.records[i].number
	}
	completed := false
	err = f.transmit(ctx, func(wire []byte) error {
		for i := range f.records {
			if bytes.Equal(f.records[i].wire, wire) {
				counts[i]++
				if i == 0 && counts[i] == 1 {
					f.ack([]recordNumber{f.records[i].number})
					ackEvent <- struct{}{}
				}
				break
			}
		}
		allTwice := true
		for i := 1; i < len(counts); i++ {
			if counts[i] < 2 {
				allTwice = false
			}
		}
		if allTwice && !completed {
			completed = true
			f.ack(oldRemaining)
			ackEvent <- struct{}{}
		}
		return nil
	}, ackEvent)
	if err != nil {
		t.Fatal(err)
	}
	if counts[0] != 1 {
		t.Fatalf("acked record sent %d times", counts[0])
	}
	for i := 1; i < len(f.records); i++ {
		if counts[i] < 2 || !f.records[i].hasPrior || f.records[i].priorNumber == f.records[i].number {
			t.Fatalf("record %d was not retransmitted", i)
		}
	}
}

func TestFlightDelayedACKCoversEveryRetransmissionGeneration(t *testing.T) {
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	f, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeCertificate, body: []byte{1}}}, 1200, sender)
	if err != nil {
		t.Fatal(err)
	}
	initial := f.records[0].number
	if err = f.refreshPending(); err != nil {
		t.Fatal(err)
	}
	firstRetransmission := f.records[0].number
	if err = f.refreshPending(); err != nil {
		t.Fatal(err)
	}
	if f.records[0].historyLength() != 2 || !f.records[0].hasNumber(initial) || !f.records[0].hasNumber(firstRetransmission) {
		t.Fatalf("record number history=%v prior=%v current=%v", f.records[0].earlierNumbers, f.records[0].priorNumber, f.records[0].number)
	}
	if changed := f.ack([]recordNumber{initial}); changed != 1 || !f.complete() {
		t.Fatalf("delayed initial ACK changed=%d complete=%v", changed, f.complete())
	}
}

func TestReserveInitialRecordHistoryUsesIndependentWindows(t *testing.T) {
	records := make([]flightRecord, 11)
	var indices flightRecordIndices
	for i := range records {
		records[i].number = recordNumber{epoch: 2, sequence: uint64(i)}
		records[i].priorNumber = recordNumber{epoch: 2, sequence: uint64(i + 100)}
		records[i].hasPrior = true
		indices.add(i, len(records))
	}
	reserveInitialRecordHistory(records, indices)
	for i := range records {
		if len(records[i].earlierNumbers) != 0 || cap(records[i].earlierNumbers) != initialRecordHistoryCapacity {
			t.Fatalf("record %d history window len=%d cap=%d", i, len(records[i].earlierNumbers), cap(records[i].earlierNumbers))
		}
		records[i].replaceNumber(recordNumber{epoch: 2, sequence: uint64(i + 200)})
	}
	for i := range records {
		want := recordNumber{epoch: 2, sequence: uint64(i + 100)}
		if len(records[i].earlierNumbers) != 1 || records[i].earlierNumbers[0] != want {
			t.Fatalf("record %d history=%v, want %v", i, records[i].earlierNumbers, want)
		}
	}
}

func TestFlightRefreshSupportsMoreThanInlineWindow(t *testing.T) {
	f := &flight{records: make([]flightRecord, 11)}
	f.rebuildPending = func(indices flightRecordIndices) error {
		if indices.count != len(f.records) {
			t.Fatalf("pending indices=%d, want %d", indices.count, len(f.records))
		}
		for i := 0; i < indices.count; i++ {
			if index := indices.at(i); index != i {
				t.Fatalf("pending index[%d]=%d", i, index)
			}
			f.records[i].replaceNumber(recordNumber{epoch: 2, sequence: uint64(i + 100)})
		}
		return nil
	}
	if err := f.refreshPending(); err != nil {
		t.Fatal(err)
	}
	for i := range f.records {
		if !f.records[i].hasPrior || f.records[i].priorNumber != (recordNumber{}) {
			t.Fatalf("record %d history not updated", i)
		}
	}
}

func TestFlightContextCancellation(t *testing.T) {
	f, _, err := buildPlainFlight([]handshakeMessage{{typ: 1, body: []byte{1}}}, 1200, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = f.transmit(ctx, func([]byte) error { return nil }, make(chan struct{}))
	if err != context.Canceled {
		t.Fatalf("got %v", err)
	}
}

func TestFlightImplicitAcknowledgement(t *testing.T) {
	f, _, err := buildPlainFlight([]handshakeMessage{{typ: 1, body: []byte{1}}}, 1200, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	events := make(chan flightEvent, 1)
	events <- flightEvent{kind: flightEventNextFlight}
	sends := 0
	if err = f.runStateMachine(context.Background(), func([]byte) error { sends++; return nil }, events, false); err != nil {
		t.Fatal(err)
	}
	if sends != 1 || f.currentState() != flightFinished {
		t.Fatalf("sends=%d state=%d", sends, f.currentState())
	}
}

func TestFlightExplicitACKAndPeerRetransmit(t *testing.T) {
	f, _, err := buildPlainFlight([]handshakeMessage{{typ: 1, body: []byte{1}}}, 1200, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	events := make(chan flightEvent, 3)
	events <- flightEvent{kind: flightEventNextFlight} // Must not finish this flight.
	events <- flightEvent{kind: flightEventPeerRetransmit}
	events <- flightEvent{kind: flightEventACK, numbers: []recordNumber{f.records[0].number}}
	sends := 0
	if err = f.runStateMachine(context.Background(), func([]byte) error { sends++; return nil }, events, true); err != nil {
		t.Fatal(err)
	}
	if sends != 2 || f.currentState() != flightFinished {
		t.Fatalf("sends=%d state=%d", sends, f.currentState())
	}
}

func TestFlightPartialACKImmediatelyRetransmitsOnlyPendingRecords(t *testing.T) {
	f, _, err := buildPlainFlight([]handshakeMessage{{typ: 1, body: bytes.Repeat([]byte{1}, 40)}}, 40, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(f.records) < 2 {
		t.Fatal("test requires multiple records")
	}
	events := make(chan flightEvent, 2)
	events <- flightEvent{kind: flightEventACK, numbers: []recordNumber{f.records[0].number}}
	remaining := make([]recordNumber, 0, len(f.records)-1)
	for i := 1; i < len(f.records); i++ {
		remaining = append(remaining, f.records[i].number)
	}
	events <- flightEvent{kind: flightEventACK, numbers: remaining}

	counts := make([]int, len(f.records))
	err = f.runStateMachine(context.Background(), func(wire []byte) error {
		for i := range f.records {
			if bytes.Equal(f.records[i].wire, wire) {
				counts[i]++
				break
			}
		}
		return nil
	}, events, true)
	if err != nil {
		t.Fatal(err)
	}
	if counts[0] != 1 {
		t.Fatalf("acked record sent %d times", counts[0])
	}
	for i := 1; i < len(counts); i++ {
		if counts[i] != 2 {
			t.Fatalf("pending record %d sent %d times, want immediate retransmission", i, counts[i])
		}
	}
}

func TestActivePartialACKRetransmitsAndFillsProtectedWindow(t *testing.T) {
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	messages := make([]handshakeMessage, 11)
	for i := range messages {
		messages[i] = handshakeMessage{typ: handshakeTypeCertificate, sequence: uint16(i), body: []byte{byte(i)}}
	}
	f, err := buildProtectedFlight(messages, 1200, sender)
	if err != nil {
		t.Fatal(err)
	}
	initial := f.nextUnsentWire(10, nil)
	if len(initial) != 10 {
		t.Fatalf("initial window=%d", len(initial))
	}
	oldNumbers := make([]recordNumber, len(f.records))
	for i := range f.records {
		oldNumbers[i] = f.records[i].number
	}
	f.ack([]recordNumber{oldNumbers[0]})
	var sent bytes.Buffer
	c := &Conn{config: &Config{Time: time.Now}}
	if err = c.retransmitPartialFlight(&sent, f); err != nil {
		t.Fatal(err)
	}
	for i := 1; i < 10; i++ {
		if f.records[i].number == oldNumbers[i] || f.records[i].historyLength() != 1 {
			t.Fatalf("record %d did not receive a fresh retransmission number", i)
		}
	}
	if f.records[0].number != oldNumbers[0] || f.records[0].historyLength() != 0 {
		t.Fatal("acknowledged record was rebuilt")
	}
	if !f.records[10].sent {
		t.Fatal("partial ACK did not fill the protected-record window")
	}
}
