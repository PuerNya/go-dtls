package dtls13

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

func TestCompletedHandshakeBatchInlineAndOverflow(t *testing.T) {
	first := completedHandshake{typ: handshakeTypeClientHello, sequence: 0, body: []byte("first")}
	second := completedHandshake{typ: handshakeTypeServerHello, sequence: 1, body: []byte("second")}

	var batch completedHandshakeBatch
	batch.add(first)
	if batch.len() != 1 || batch.values != nil || batch.at(0).typ != first.typ || string(batch.at(0).body) != "first" {
		t.Fatalf("unexpected inline batch: %#v", batch)
	}

	batch.add(second)
	if batch.len() != 2 || len(batch.values) != 2 {
		t.Fatalf("unexpected overflow batch: %#v", batch)
	}
	if batch.at(0).typ != first.typ || string(batch.at(0).body) != "first" || batch.at(1).typ != second.typ || string(batch.at(1).body) != "second" {
		t.Fatalf("overflow changed message order: %#v", batch.values)
	}
}

func TestHandshakeInboxDeliversInOrder(t *testing.T) {
	inbox := newHandshakeInbox(0, 1024, 8, 4096)
	future := handshakeFragment{typ: 2, messageSequence: 1, length: 1, body: []byte("b")}
	messages, err := inbox.add(future)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 0 {
		t.Fatal("delivered future message")
	}
	messages, err = inbox.add(handshakeFragment{typ: 1, messageSequence: 0, length: 1, body: []byte("a")})
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 2 || string(messages[0].body) != "a" || string(messages[1].body) != "b" {
		t.Fatalf("unexpected delivery %#v", messages)
	}
	if inbox.expected != 2 {
		t.Fatalf("expected sequence %d", inbox.expected)
	}
}

func TestHandshakeInboxReusesDeliveryDestination(t *testing.T) {
	inbox := newHandshakeInbox(0, 1024, 8, 4096)
	if messages, err := inbox.add(handshakeFragment{typ: 2, messageSequence: 1, length: 1, body: []byte("b")}); err != nil || len(messages) != 0 {
		t.Fatalf("future messages=%v err=%v", messages, err)
	}
	var storage [2]completedHandshake
	messages, err := inbox.addInto(storage[:0], handshakeFragment{typ: 1, messageSequence: 0, length: 1, body: []byte("a")})
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 2 || &messages[0] != &storage[0] || string(messages[0].body) != "a" || string(messages[1].body) != "b" {
		t.Fatalf("unexpected reused delivery %#v", messages)
	}
	if len(inbox.ready) != 0 || inbox.readyBytes != 0 || inbox.expected != 2 {
		t.Fatalf("ready=%d bytes=%d expected=%d", len(inbox.ready), inbox.readyBytes, inbox.expected)
	}
}

func TestHandshakeInboxBatchDeliversBufferedMessagesInOrder(t *testing.T) {
	inbox := newHandshakeInbox(0, 1024, 8, 4096)
	var batch completedHandshakeBatch
	if err := inbox.addBatch(&batch, handshakeFragment{typ: 2, messageSequence: 1, length: 1, body: []byte("b")}); err != nil {
		t.Fatal(err)
	}
	if batch.len() != 0 {
		t.Fatalf("delivered future message: %#v", batch)
	}
	if err := inbox.addBatch(&batch, handshakeFragment{typ: 1, messageSequence: 0, length: 1, body: []byte("a")}); err != nil {
		t.Fatal(err)
	}
	if batch.len() != 2 || string(batch.at(0).body) != "a" || string(batch.at(1).body) != "b" {
		t.Fatalf("unexpected batch delivery: %#v", batch)
	}
	if len(inbox.ready) != 0 || inbox.readyBytes != 0 || inbox.expected != 2 {
		t.Fatalf("ready=%d bytes=%d expected=%d", len(inbox.ready), inbox.readyBytes, inbox.expected)
	}
}

func TestHandshakeInboxDiscardsOldMessage(t *testing.T) {
	inbox := newHandshakeInbox(3, 1024, 8, 4096)
	messages, err := inbox.add(handshakeFragment{typ: 1, messageSequence: 2, length: 1, body: []byte("x")})
	if err != nil || len(messages) != 0 {
		t.Fatalf("messages=%v err=%v", messages, err)
	}
	if len(inbox.reassembler.messages) != 0 {
		t.Fatal("buffered old message")
	}
}

func TestProtectedHandshakeACKsBufferedOutOfOrderFragment(t *testing.T) {
	recordSender, recordReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	ackSender, ackReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	flight, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeCertificate, sequence: 0, body: bytes.Repeat([]byte{0x42}, 100)}}, 100, recordSender)
	if err != nil {
		t.Fatal(err)
	}
	if len(flight.records) != 2 {
		t.Fatalf("flight records=%d, want 2", len(flight.records))
	}
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	_ = left.SetDeadline(time.Now().Add(time.Second))
	_ = right.SetDeadline(time.Now().Add(time.Second))
	result := make(chan []completedHandshake, 1)
	errCh := make(chan error, 1)
	go func() {
		messages, receiveErr := receiveHandshakeMessageWithEarly(right, newHandshakeInbox(0, 1024, 8, 4096), recordReceiver, nil, nil, nil, ackSender, 100, nil)
		if receiveErr != nil {
			errCh <- receiveErr
			return
		}
		result <- messages
	}()

	if _, err = left.Write(flight.records[1].wire); err != nil {
		t.Fatal(err)
	}
	ackWire := make([]byte, 256)
	n, err := left.Read(ackWire)
	if err != nil {
		t.Fatal(err)
	}
	content, typ, _, err := ackReceiver.open(ackWire[:n])
	if err != nil || typ != recordTypeACK {
		t.Fatalf("open partial ACK: type=%d err=%v", typ, err)
	}
	numbers, err := parseACK(content)
	if err != nil || len(numbers) != 1 || numbers[0] != flight.records[1].number {
		t.Fatalf("partial ACK numbers=%v err=%v", numbers, err)
	}

	if _, err = left.Write(flight.records[0].wire); err != nil {
		t.Fatal(err)
	}
	select {
	case messages := <-result:
		if len(messages) != 1 || !bytes.Equal(messages[0].body, bytes.Repeat([]byte{0x42}, 100)) {
			t.Fatalf("unexpected delivered messages %#v", messages)
		}
	case err = <-errCh:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("complete message was not delivered")
	}
}

func TestHandshakeReceivePropagatesAEADAuthenticationFailureLimit(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	receiver.authFailureLimit = 1
	wire, err := sender.seal(recordTypeHandshake, []byte{1})
	if err != nil {
		t.Fatal(err)
	}
	wire[len(wire)-1] ^= 1
	go func() { _, _ = left.Write(wire) }()
	_, err = receiveHandshakeMessageWithEarly(right, newHandshakeInbox(0, 1024, 8, 4096), receiver, nil, nil, nil, nil, 1200, nil)
	if !errors.Is(err, errAEADAuthenticationFailureLimit) {
		t.Fatalf("handshake receive returned %v", err)
	}
}

func TestHandshakeReceiveProcessesPlaintextFatalAlert(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	body, err := (alertMessage{level: alertLevelFatal, description: alertIllegalParameter}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	wire, err := marshalPlainRecord(record{typ: recordTypeAlert, payload: body})
	if err != nil {
		t.Fatal(err)
	}
	go func() { _, _ = left.Write(wire) }()
	_, err = receiveHandshakeMessage(right, newHandshakeInbox(0, 1024, 8, 4096), nil)
	if !errors.Is(err, AlertError(alertIllegalParameter)) {
		t.Fatalf("receive returned %v", err)
	}
}

func TestHandshakeReceiveTreatsWarningAlertAsFatal(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	body, err := (alertMessage{level: alertLevelWarning, description: alertIllegalParameter}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	wire, err := marshalPlainRecord(record{typ: recordTypeAlert, payload: body})
	if err != nil {
		t.Fatal(err)
	}
	go func() { _, _ = left.Write(wire) }()
	_, err = receiveHandshakeMessage(right, newHandshakeInbox(0, 1024, 8, 4096), nil)
	if !errors.Is(err, AlertError(alertIllegalParameter)) {
		t.Fatalf("receive returned %v", err)
	}
}

func TestHelloRetryRequestRetransmitsOnlyForRepeatedClientHello(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	config, err := (&Config{MTU: 1200}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	conn := &Conn{conn: right, config: config}
	hrr, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeServerHello, sequence: 0, body: []byte("hrr")}}, 1200, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err = conn.writeFlight(right, hrr); err != nil {
		t.Fatal(err)
	}
	// Consume the initial HRR and then retransmit the first ClientHello.
	buffer := make([]byte, 2048)
	if _, err = left.Read(buffer); err != nil {
		t.Fatal(err)
	}
	first, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeClientHello, sequence: 0, body: []byte("first")}}, 1200, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := buildPlainFlight([]handshakeMessage{{typ: handshakeTypeClientHello, sequence: 1, body: []byte("second")}}, 1200, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	result := make(chan completedHandshakeBatch, 1)
	errCh := make(chan error, 1)
	go func() {
		messages, receiveErr := conn.receiveSecondClientHello(right, newHandshakeInbox(1, 1024, 8, 4096), hrr)
		if receiveErr != nil {
			errCh <- receiveErr
			return
		}
		result <- messages
	}()
	if _, err = left.Write(first.records[0].wire); err != nil {
		t.Fatal(err)
	}
	_ = left.SetReadDeadline(time.Now().Add(time.Second))
	if _, err = left.Read(buffer); err != nil {
		t.Fatal("repeated ClientHello did not trigger HRR retransmission")
	}
	if _, err = left.Write(second.records[0].wire); err != nil {
		t.Fatal(err)
	}
	select {
	case messages := <-result:
		if messages.len() != 1 || string(messages.at(0).body) != "second" {
			t.Fatalf("messages=%#v", messages)
		}
	case err = <-errCh:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("second ClientHello was not delivered")
	}
}

func TestUndecryptableHandshakeRecordSendsEmptyACK(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	_ = left.SetDeadline(time.Now().Add(time.Second))
	_ = right.SetDeadline(time.Now().Add(50 * time.Millisecond))
	config, err := (&Config{}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	owner := &Conn{conn: right, config: config}
	owner.plainSendSequence.Store(5)
	sender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	wire, err := sender.seal(recordTypeHandshake, []byte{1})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		_, _ = receiveHandshakeMessageWithEarly(right, newHandshakeInbox(0, 1024, 8, 4096), nil, nil, nil, nil, nil, owner.currentMTU(), owner)
	}()
	if _, err = left.Write(wire); err != nil {
		t.Fatal(err)
	}
	ackWire := make([]byte, 256)
	n, err := left.Read(ackWire)
	if err != nil {
		t.Fatal(err)
	}
	records, err := parsePlainRecords(ackWire[:n])
	if err != nil || len(records) != 1 || records[0].typ != recordTypeACK || records[0].sequence != 5 {
		t.Fatalf("empty ACK record=%v err=%v", records, err)
	}
	numbers, err := parseACK(records[0].payload)
	if err != nil || len(numbers) != 0 {
		t.Fatalf("empty ACK numbers=%v err=%v", numbers, err)
	}
}

func TestDiscardedOldHandshakeRecordIsNotAcknowledged(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	_ = left.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	handshakeSender, handshakeReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	ackSender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	result := make(chan error, 1)
	go func() {
		_, err := receiveHandshakeMessageWithEarly(right, newHandshakeInbox(1, 1024, 8, 4096), handshakeReceiver, nil, nil, nil, ackSender, 1200, nil)
		result <- err
	}()
	fragment, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeFinished, messageSequence: 0, length: 1, body: []byte{1}})
	if err != nil {
		t.Fatal(err)
	}
	wire, err := handshakeSender.seal(recordTypeHandshake, fragment)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = left.Write(wire); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 256)
	if _, err = left.Read(buffer); err == nil {
		t.Fatal("receiver acknowledged a discarded old handshake record")
	}
	_ = right.Close()
	<-result
}

func TestPeerFlightRetransmissionImmediatelyRetransmitsUnacknowledgedResponse(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	_ = left.SetReadDeadline(time.Now().Add(time.Second))
	config, err := (&Config{}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	owner := &Conn{conn: right, config: config}
	responseSender, responseReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	response, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeFinished, sequence: 1, body: []byte("response")}}, 1200, responseSender)
	if err != nil {
		t.Fatal(err)
	}
	if err = owner.writeFlight(right, response); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 2048)
	if _, err = left.Read(buffer); err != nil {
		t.Fatal(err)
	}

	requestSender, requestReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	done := make(chan error, 1)
	go func() {
		_, receiveErr := receiveHandshakeMessageWithEarly(right, newHandshakeInbox(1, 1024, 8, 4096), requestReceiver, nil, nil, response, nil, 1200, owner)
		done <- receiveErr
	}()
	fragment, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeCertificate, messageSequence: 0, length: 1, body: []byte{1}})
	if err != nil {
		t.Fatal(err)
	}
	wire, err := requestSender.seal(recordTypeHandshake, fragment)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = left.Write(wire); err != nil {
		t.Fatal(err)
	}
	n, err := left.Read(buffer)
	if err != nil {
		t.Fatal("peer flight retransmission did not trigger an immediate response retransmission")
	}
	content, typ, _, err := responseReceiver.open(buffer[:n])
	if err != nil || typ != recordTypeHandshake || len(content) == 0 {
		t.Fatalf("immediate retransmission type=%d content=%x err=%v", typ, content, err)
	}
	_ = right.Close()
	<-done
}

func TestEmptyPlaintextACKRetransmitsProtectedFlight(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	_ = left.SetDeadline(time.Now().Add(time.Second))
	_ = right.SetDeadline(time.Now().Add(50 * time.Millisecond))
	config, err := (&Config{}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	owner := &Conn{conn: right, config: config}
	flightSender, flightReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	outgoing, err := buildProtectedFlight([]handshakeMessage{{typ: handshakeTypeEncryptedExtensions, sequence: 1, body: []byte("flight")}}, 1200, flightSender)
	if err != nil {
		t.Fatal(err)
	}
	if err = owner.writeFlight(right, outgoing); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 2048)
	if _, err = left.Read(buffer); err != nil {
		t.Fatal(err)
	}
	_, incomingCipher := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	go func() {
		_, _ = receiveHandshakeMessageWithEarly(right, newHandshakeInbox(0, 1024, 8, 4096), incomingCipher, nil, nil, outgoing, nil, owner.currentMTU(), owner)
	}()
	empty, _, err := buildACKRecords(nil, 1200, 0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = left.Write(empty[0]); err != nil {
		t.Fatal(err)
	}
	n, err := left.Read(buffer)
	if err != nil {
		t.Fatal("empty ACK did not trigger retransmission")
	}
	content, typ, _, err := flightReceiver.open(buffer[:n])
	if err != nil || typ != recordTypeHandshake || len(content) == 0 {
		t.Fatalf("retransmitted flight type=%d content=%x err=%v", typ, content, err)
	}
}

func TestHandshakeDropsPrematureApplicationEpoch(t *testing.T) {
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	_ = left.SetDeadline(time.Now().Add(time.Second))
	_ = right.SetDeadline(time.Now().Add(time.Second))
	handshakeSender, handshakeReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 2)
	futureSender, _ := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	result := make(chan []completedHandshake, 1)
	errCh := make(chan error, 1)
	go func() {
		messages, err := receiveHandshakeMessageWithEarly(right, newHandshakeInbox(0, 1024, 8, 4096), handshakeReceiver, nil, nil, nil, nil, 1200, nil)
		if err != nil {
			errCh <- err
			return
		}
		result <- messages
	}()
	premature, err := futureSender.seal(recordTypeApplicationData, []byte("future"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err = left.Write(premature); err != nil {
		t.Fatal(err)
	}
	fragment, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeFinished, messageSequence: 0, length: 1, body: []byte{1}})
	if err != nil {
		t.Fatal(err)
	}
	valid, err := handshakeSender.seal(recordTypeHandshake, fragment)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = left.Write(valid); err != nil {
		t.Fatal(err)
	}
	select {
	case messages := <-result:
		if len(messages) != 1 || messages[0].typ != handshakeTypeFinished {
			t.Fatalf("messages=%#v", messages)
		}
	case err = <-errCh:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("valid handshake record was blocked by premature epoch")
	}
}
