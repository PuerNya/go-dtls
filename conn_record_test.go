package dtls13

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"testing"
	"time"
)

type mtuLimitedConn struct {
	net.Conn
	limit  int
	writes int
}

type recordSinkConn struct {
	net.Conn
	writes [][]byte
}

type failNthRecordWriteConn struct {
	net.Conn
	writes int
	failAt int
}

func bufferedApplicationPayload(c *Conn) []byte {
	c.inputMu.Lock()
	defer c.inputMu.Unlock()
	var payload []byte
	for _, datagram := range c.applicationDatagrams {
		payload = append(payload, datagram.payload...)
	}
	return payload
}

func bufferedApplicationDatagrams(c *Conn) int {
	c.inputMu.Lock()
	defer c.inputMu.Unlock()
	return len(c.applicationDatagrams)
}

func (c *failNthRecordWriteConn) Write(p []byte) (int, error) {
	c.writes++
	if c.writes == c.failAt {
		return 0, io.ErrClosedPipe
	}
	return len(p), nil
}

func (c *recordSinkConn) Write(p []byte) (int, error) {
	c.writes = append(c.writes, append([]byte(nil), p...))
	return len(p), nil
}

func (c *mtuLimitedConn) Write(p []byte) (int, error) {
	c.writes++
	if len(p) > c.limit {
		return 0, syscall.Errno(10040)
	}
	return c.Conn.Write(p)
}

func establishedConnPair(t *testing.T) (*Conn, *Conn) {
	t.Helper()
	left, right := net.Pipe()
	client := Client(left, &Config{})
	server := Server(right, &Config{})
	var err error
	client.config, err = client.config.normalized()
	if err != nil {
		t.Fatal(err)
	}
	server.config, err = server.config.normalized()
	if err != nil {
		t.Fatal(err)
	}
	client.handshakeOnce.Do(func() {})
	server.handshakeOnce.Do(func() {})
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	clientSecret := bytes.Repeat([]byte{1}, suite.hash.Size())
	serverSecret := bytes.Repeat([]byte{2}, suite.hash.Size())
	if err = client.installApplicationKeys(suite, clientSecret, serverSecret); err != nil {
		t.Fatal(err)
	}
	if err = server.installApplicationKeys(suite, clientSecret, serverSecret); err != nil {
		t.Fatal(err)
	}
	return client, server
}

func TestConnPropagatesAEADAuthenticationFailureLimit(t *testing.T) {
	client, server := establishedConnPair(t)
	server.receiveEpochs.mu.Lock()
	server.receiveEpochs.ciphers[3].authFailureLimit = 1
	server.receiveEpochs.mu.Unlock()
	wire, err := client.sendCipher.seal(recordTypeApplicationData, []byte("authenticated"))
	if err != nil {
		t.Fatal(err)
	}
	wire[len(wire)-1] ^= 1
	if err = server.dispatchDatagram(wire); !errors.Is(err, errAEADAuthenticationFailureLimit) {
		t.Fatalf("authentication failure limit returned %v", err)
	}
	_ = client.conn.Close()
	_ = server.conn.Close()
}

func TestConnRequestsKeyUpdateBeforeAuthenticationFailureLimit(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	updateWire := make(chan []byte, 1)
	go func() {
		buf := make([]byte, 2048)
		n, _ := client.conn.Read(buf)
		updateWire <- append([]byte(nil), buf[:n]...)
	}()

	receiver := server.receiveEpochs.ciphers[3]
	receiver.authFailureLimit = 4
	for i := 0; i < 3; i++ {
		wire, err := client.sendCipher.seal(recordTypeApplicationData, []byte("forged"))
		if err != nil {
			t.Fatal(err)
		}
		wire[len(wire)-1] ^= 1
		if err = server.dispatchDatagram(wire); err != nil {
			t.Fatalf("failure %d returned %v", i+1, err)
		}
	}
	if wire := <-updateWire; len(wire) == 0 {
		t.Fatal("peer did not receive KeyUpdate request")
	}
	if server.sendingTraffic.update.canUseNewKeys() {
		t.Fatal("authentication failure threshold did not start KeyUpdate")
	}
	if len(server.sendingTraffic.pendingFragment) == 0 {
		t.Fatal("KeyUpdate request fragment was not retained for reliable retransmission")
	}
	fragments, err := parseHandshakeFragments(server.sendingTraffic.pendingFragment)
	if err != nil || len(fragments) != 1 {
		t.Fatalf("parse KeyUpdate: fragments=%d err=%v", len(fragments), err)
	}
	message, err := parseKeyUpdate(fragments[0].body)
	if err != nil || !message.requestUpdate {
		t.Fatalf("KeyUpdate request=%v err=%v", message.requestUpdate, err)
	}
}

func TestEarlyDataOverLimitReturnsUnexpectedMessage(t *testing.T) {
	config, err := (&Config{}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	c := &Conn{config: config, earlyAccepted: true, earlyDataLimit: 4}
	if err := c.queueEarlyApplicationData([]byte("1234")); err != nil {
		t.Fatal(err)
	}
	err = c.queueEarlyApplicationData([]byte("5"))
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertUnexpectedMessage {
		t.Fatalf("over-limit early data returned %v", err)
	}
}

func TestApplicationDatagramQueueIsCountBounded(t *testing.T) {
	config, err := (&Config{MaxBufferedApplicationDatagrams: 2}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	c := &Conn{config: config}
	if err = c.queueApplicationData(nil, memoryAddr("first")); err != nil {
		t.Fatal(err)
	}
	if err = c.queueApplicationData(nil, memoryAddr("second")); err != nil {
		t.Fatal(err)
	}
	if err = c.queueApplicationData(nil, memoryAddr("third")); err == nil {
		t.Fatal("zero-length datagrams bypassed the queue count limit")
	}
}

func TestConnEncryptedDatagramRoundTrip(t *testing.T) {
	client, server := establishedConnPair(t)
	payload := bytes.Repeat([]byte("x"), 700)
	errCh := make(chan error, 1)
	go func() { _, err := client.WriteDatagram(payload); errCh <- err }()
	buf := make([]byte, len(payload))
	n, info, err := server.ReadDatagram(buf)
	if err != nil {
		t.Fatal(err)
	}
	if err := <-errCh; err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf[:n], payload) || info.Truncated || info.FullLength != len(payload) {
		t.Fatal("application plaintext mismatch")
	}
	_ = client.conn.Close()
	_ = server.conn.Close()
}

func TestConnWriteDatagramIgnoresPathMTU(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	client.config.IgnorePathMTU = true
	client.pathMTU.Store(256)
	payload := bytes.Repeat([]byte("x"), 700)
	written := make(chan error, 1)
	go func() { _, err := client.WriteDatagram(payload); written <- err }()
	buf := make([]byte, len(payload))
	n, _, err := server.ReadDatagram(buf)
	if err != nil {
		t.Fatal(err)
	}
	if err = <-written; err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf[:n], payload) {
		t.Fatal("application plaintext mismatch")
	}
}

func TestConnWriteDatagramIgnorePathMTUDoesNotRetry(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	limited := &mtuLimitedConn{Conn: client.conn, limit: 500}
	client.conn = limited
	client.config.IgnorePathMTU = true
	wantMTU := client.currentMTU()
	n, err := client.WriteDatagram(bytes.Repeat([]byte("p"), 700))
	if n != 0 || !errors.Is(err, ErrDatagramTooLarge) {
		t.Fatalf("WriteDatagram=%d, %v", n, err)
	}
	if limited.writes != 1 {
		t.Fatalf("transport writes=%d, want 1", limited.writes)
	}
	if got := client.currentMTU(); got != wantMTU {
		t.Fatalf("path MTU=%d, want %d", got, wantMTU)
	}
}

func TestConnDatagramWriteReducesPathMTUWithoutPartialSend(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	client.conn = &mtuLimitedConn{Conn: client.conn, limit: 500}
	payload := bytes.Repeat([]byte("p"), 700)
	n, err := client.WriteDatagram(payload)
	if n != 0 || !errors.Is(err, ErrDatagramTooLarge) {
		t.Fatalf("WriteDatagram=%d, %v", n, err)
	}
	if client.currentMTU() >= client.config.MTU {
		t.Fatalf("path MTU=%d was not reduced", client.currentMTU())
	}
}

func TestConnDatagramWriteRespectsRecordLimitWithLargeMTU(t *testing.T) {
	client, server := establishedConnPair(t)
	originalClientConn := client.conn
	defer originalClientConn.Close()
	defer server.conn.Close()
	sink := &recordSinkConn{Conn: originalClientConn}
	client.conn = sink
	client.config.IgnorePathMTU = true
	client.pathMTU.Store(65535)
	want := bytes.Repeat([]byte{0x5a}, 2*maxRecordContent)
	n, err := client.WriteDatagram(want)
	if n != 0 || !errors.Is(err, ErrDatagramTooLarge) {
		t.Fatalf("WriteDatagram = %d, %v", n, err)
	}
	if len(sink.writes) != 0 {
		t.Fatalf("large application datagram generated %d records", len(sink.writes))
	}
}

func TestConnDropsDelayedEarlyDataAfterApplicationKeys(t *testing.T) {
	client, server := establishedConnPair(t)
	suite, err := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	if err != nil {
		t.Fatal(err)
	}
	early, err := newRecordCipher(suite, bytes.Repeat([]byte{9}, suite.hash.Size()), 1, 64)
	if err != nil {
		t.Fatal(err)
	}
	wire, err := early.seal(recordTypeApplicationData, []byte("late early data"))
	if err != nil {
		t.Fatal(err)
	}
	if err = server.dispatchDatagram(wire); err != nil {
		t.Fatalf("delayed epoch-1 record was not silently discarded: %v", err)
	}
	buffered := bufferedApplicationDatagrams(server)
	if buffered != 0 {
		t.Fatal("delayed epoch-1 data reached the application")
	}
	_ = client.conn.Close()
	_ = server.conn.Close()
}

func TestConnRejectsUnexpectedPostHandshakeMessages(t *testing.T) {
	for _, typ := range []uint8{handshakeTypeEndOfEarlyData, 0xff} {
		t.Run(fmt.Sprintf("type-%d", typ), func(t *testing.T) {
			client, server := establishedConnPair(t)
			fragment, err := marshalHandshakeFragment(handshakeFragment{typ: typ, messageSequence: 0, length: 0})
			if err != nil {
				t.Fatal(err)
			}
			wire, err := client.sendCipher.seal(recordTypeHandshake, fragment)
			if err != nil {
				t.Fatal(err)
			}
			err = server.dispatchDatagram(wire)
			var alertErr *localAlertError
			if !errors.As(err, &alertErr) || alertErr.description != alertUnexpectedMessage {
				t.Fatalf("unexpected post-handshake type %d returned %v", typ, err)
			}
			_ = client.conn.Close()
			_ = server.conn.Close()
		})
	}
}

func TestConnRejectsAuthenticatedChangeCipherSpec(t *testing.T) {
	client, server := establishedConnPair(t)
	wire := sealInvalidInnerType(t, client.sendCipher, recordTypeChangeCipherSpec)
	err := server.dispatchDatagram(wire)
	var alertErr *localAlertError
	if !errors.As(err, &alertErr) || alertErr.description != alertUnexpectedMessage {
		t.Fatalf("authenticated ChangeCipherSpec returned %v", err)
	}
	_ = client.conn.Close()
	_ = server.conn.Close()
}

func TestConnRejectsMalformedAuthenticatedACK(t *testing.T) {
	for _, test := range []struct {
		name string
		body []byte
		want uint8
	}{
		{name: "malformed-vector", body: []byte{0, 1, 0}, want: alertDecodeError},
		{name: "non-increasing", body: []byte{
			0, 32,
			0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 2,
			0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 1,
		}, want: alertIllegalParameter},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, server := establishedConnPair(t)
			wire, err := client.sendCipher.seal(recordTypeACK, test.body)
			if err != nil {
				t.Fatal(err)
			}
			err = server.dispatchDatagram(wire)
			var local *localAlertError
			if !errors.As(err, &local) || local.description != test.want {
				t.Fatalf("authenticated ACK returned %v, want alert %d", err, test.want)
			}
			_ = client.conn.Close()
			_ = server.conn.Close()
		})
	}
}

func TestRecordReaderAlertsOnAuthenticatedHandshakeDecodeError(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	server.mu.Lock()
	server.state.exporter = newExporter(suite, bytes.Repeat([]byte{7}, suite.hash.Size()))
	server.mu.Unlock()
	stateBeforeFatal := server.ConnectionState()
	server.startRecordReader()
	wire, err := client.sendCipher.seal(recordTypeHandshake, []byte{1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = client.conn.Write(wire); err != nil {
		t.Fatal(err)
	}
	response := make([]byte, 2048)
	n, err := client.conn.Read(response)
	if err != nil {
		t.Fatal(err)
	}
	content, typ, _, _, err := client.receiveEpochs.open(response[:n])
	if err != nil {
		t.Fatal(err)
	}
	alert, err := parseAlert(content)
	if err != nil || typ != recordTypeAlert || alert.description != alertDecodeError {
		t.Fatalf("type=%d alert=%v err=%v", typ, alert, err)
	}
	deadline := time.Now().Add(time.Second)
	cleared := false
	for time.Now().Before(deadline) {
		server.writeMu.Lock()
		cleared = server.sendCipher == nil && server.sendingTraffic == nil && server.receivingTraffic == nil
		server.writeMu.Unlock()
		cleared = cleared && server.ConnectionState().exporter == nil
		if cleared {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if !cleared {
		t.Fatal("fatal record error retained traffic or exporter secrets")
	}
	if _, err = stateBeforeFatal.ExportKeyingMaterial("test", nil, 16); err == nil {
		t.Fatal("ConnectionState copied before fatal retained exporter access")
	}
}

func TestConnUsesRecordOverflowForOversizedAuthenticatedInnerPlaintext(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	plain := make([]byte, (1<<14)+2)
	plain[len(plain)-1] = recordTypeApplicationData
	wire := sealRawInnerPlaintext(t, client.sendCipher, plain)
	err := server.dispatchDatagram(wire)
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertRecordOverflow {
		t.Fatalf("oversized authenticated inner plaintext returned %v", err)
	}
}

func TestConnRejectsZeroLengthAuthenticatedHandshake(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	wire := sealRawInnerPlaintext(t, client.sendCipher, []byte{recordTypeHandshake})
	err := server.dispatchDatagram(wire)
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertUnexpectedMessage {
		t.Fatalf("zero-length authenticated Handshake returned %v", err)
	}
}

func TestConnTreatsUnknownAlertLevelAsError(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	wire := sealRawInnerPlaintext(t, client.sendCipher, []byte{255, alertCloseNotify, recordTypeAlert})
	err := server.dispatchDatagram(wire)
	var peerAlert AlertError
	if !errors.As(err, &peerAlert) || uint8(peerAlert) != alertCloseNotify {
		t.Fatalf("unknown-level alert returned %v", err)
	}
	if server.peerReadClosed {
		t.Fatal("unknown-level close_notify closed the read side cleanly")
	}
}

func TestConnCloseClearsTrafficSecrets(t *testing.T) {
	client, server := establishedConnPair(t)
	originalClientConn := client.conn
	client.conn = &recordSinkConn{Conn: originalClientConn}
	defer originalClientConn.Close()
	defer server.conn.Close()

	client.resumptionMasterSecret = []byte("resumption secret")
	client.resumptionSuite = client.sendingTraffic.suite
	client.sendingTraffic.nextSecret = bytes.Repeat([]byte{0x41}, client.sendingTraffic.suite.hash.Size())
	staleReceiveSecret := bytes.Repeat([]byte{0x42}, client.receivingTraffic.suite.hash.Size())
	client.receivingTraffic.secrets[client.receivingTraffic.current-1] = staleReceiveSecret
	sendingSecret := client.sendingTraffic.secret
	nextSendingSecret := client.sendingTraffic.nextSecret
	receivingSecret := client.receivingTraffic.secret
	resumptionSecret := client.resumptionMasterSecret
	client.mu.Lock()
	client.state.exporter = newExporter(client.sendingTraffic.suite, bytes.Repeat([]byte{0x43}, client.sendingTraffic.suite.hash.Size()))
	exporterSecret := client.state.exporter.secret
	client.mu.Unlock()
	client.postHandshakeReassembly = newReassembler()
	client.protectedHandshakeRanges = []protectedHandshakeRecordRange{{first: recordNumber{epoch: 3}, last: recordNumber{epoch: 3, sequence: 2}}}
	client.recentApplicationRecords = []recordNumber{{epoch: 3, sequence: 1}}
	client.pendingHandshakeApplications = []pendingPostAuthApplication{{content: []byte("pending")}}
	if err := client.Close(); err != nil {
		t.Fatal(err)
	}
	if client.sendCipher != nil || client.sendingTraffic != nil || client.receivingTraffic != nil || client.resumptionMasterSecret != nil || client.resumptionSuite != nil {
		t.Fatal("Close retained traffic or resumption secrets")
	}
	for name, secret := range map[string][]byte{
		"sending": sendingSecret, "next sending": nextSendingSecret, "receiving": receivingSecret,
		"stale receiving": staleReceiveSecret, "resumption": resumptionSecret, "exporter": exporterSecret,
	} {
		if !bytes.Equal(secret, make([]byte, len(secret))) {
			t.Fatalf("Close did not clear %s secret backing", name)
		}
	}
	client.receiveEpochs.mu.RLock()
	retainedEpochs := len(client.receiveEpochs.ciphers)
	client.receiveEpochs.mu.RUnlock()
	if retainedEpochs != 0 {
		t.Fatalf("Close retained %d receive epochs", retainedEpochs)
	}
	if client.postHandshakeReassembly != nil || client.protectedHandshakeRanges != nil || client.recentApplicationRecords != nil || client.pendingHandshakeApplications != nil {
		t.Fatal("Close retained post-handshake ordering or reassembly state")
	}
}

func TestPostHandshakeRetransmissionFailureTerminatesConnection(t *testing.T) {
	client, server := establishedConnPair(t)
	originalClientConn := client.conn
	client.conn = &recordSinkConn{Conn: originalClientConn}
	defer originalClientConn.Close()
	defer server.conn.Close()

	flight, err := buildProtectedFlight([]handshakeMessage{{
		typ: handshakeTypeNewSessionTicket, sequence: client.sendingTraffic.messageSequence, body: []byte{0},
	}}, client.currentMTU(), client.sendCipher)
	if err != nil {
		t.Fatal(err)
	}
	flight.setIntervals(time.Millisecond, time.Millisecond)
	client.writeMu.Lock()
	client.config.FlightInterval = time.Millisecond
	client.config.MaxFlightInterval = time.Millisecond
	client.ticketFlight = flight
	client.sendCipher.nextSequence = client.sendCipher.recordLimit
	client.writeMu.Unlock()
	client.startTicketRetransmission()

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		client.inputMu.Lock()
		readErr := client.readErr
		client.inputMu.Unlock()
		if readErr != nil {
			var protocol *ProtocolError
			if !errors.As(readErr, &protocol) {
				t.Fatalf("retransmission failure = %v", readErr)
			}
			if client.sendCipher != nil || client.sendingTraffic != nil || client.receivingTraffic != nil {
				t.Fatal("failed retransmission retained traffic secrets")
			}
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("retransmission failure was not surfaced to the connection")
}

func TestConnRequiresKeyUpdateRecordBoundaryAlignment(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	sequence := client.sendingTraffic.messageSequence
	one, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeKeyUpdate, messageSequence: sequence, length: 1, body: []byte{0}})
	if err != nil {
		t.Fatal(err)
	}
	two, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeKeyUpdate, messageSequence: sequence + 1, length: 1, body: []byte{0}})
	if err != nil {
		t.Fatal(err)
	}
	wire, err := client.sendCipher.seal(recordTypeHandshake, append(one, two...))
	if err != nil {
		t.Fatal(err)
	}
	err = server.dispatchDatagram(wire)
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertUnexpectedMessage {
		t.Fatalf("coalesced KeyUpdate returned %v", err)
	}
	if server.receivingTraffic.current != 3 {
		t.Fatalf("coalesced KeyUpdate advanced receive epoch to %d", server.receivingTraffic.current)
	}
}

func TestRequestedKeyUpdateResponseWaitsForOutstandingUpdateACK(t *testing.T) {
	client, server := establishedConnPair(t)
	originalClientConn := client.conn
	defer originalClientConn.Close()
	defer server.conn.Close()
	sink := &recordSinkConn{Conn: originalClientConn}
	client.conn = sink

	_, ownNumber, err := client.sendingTraffic.beginKeyUpdate(false)
	if err != nil {
		t.Fatal(err)
	}
	peerWire, _, err := server.sendingTraffic.beginKeyUpdate(true)
	if err != nil {
		t.Fatal(err)
	}
	if err = client.dispatchDatagram(peerWire); err != nil {
		t.Fatal(err)
	}
	if !client.keyUpdateResponsePending {
		t.Fatal("update_requested was not deferred while a local KeyUpdate was outstanding")
	}
	ackBody, err := marshalACK([]recordNumber{ownNumber})
	if err != nil {
		t.Fatal(err)
	}
	ackWire, err := server.sendCipher.seal(recordTypeACK, ackBody)
	if err != nil {
		t.Fatal(err)
	}
	if err = client.dispatchDatagram(ackWire); err != nil {
		t.Fatal(err)
	}
	if client.keyUpdateResponsePending {
		t.Fatal("deferred KeyUpdate response remained pending after the local ACK")
	}
	if client.sendingTraffic.cipher.epoch != 4 || client.sendingTraffic.update.canUseNewKeys() {
		t.Fatalf("deferred response state epoch=%d available=%v", client.sendingTraffic.cipher.epoch, client.sendingTraffic.update.canUseNewKeys())
	}
	fragments, err := parseHandshakeFragments(client.sendingTraffic.pendingFragment)
	if err != nil || len(fragments) != 1 {
		t.Fatalf("deferred response fragments=%d err=%v", len(fragments), err)
	}
	message, err := parseKeyUpdate(fragments[0].body)
	if err != nil || message.requestUpdate {
		t.Fatalf("deferred response message=%#v err=%v", message, err)
	}
	client.readerMu.Lock()
	client.readerClosed = true
	client.readerMu.Unlock()
}

func TestIncomingKeyUpdatePropagatesACKAndResponseWriteErrors(t *testing.T) {
	for _, test := range []struct {
		name   string
		failAt int
	}{
		{name: "ACK", failAt: 1},
		{name: "automatic response", failAt: 2},
	} {
		t.Run(test.name, func(t *testing.T) {
			client, server := establishedConnPair(t)
			originalClientConn := client.conn
			defer originalClientConn.Close()
			defer server.conn.Close()
			client.conn = &failNthRecordWriteConn{Conn: originalClientConn, failAt: test.failAt}

			wire, _, err := server.sendingTraffic.beginKeyUpdate(true)
			if err != nil {
				t.Fatal(err)
			}
			if err = client.dispatchDatagram(wire); !errors.Is(err, io.ErrClosedPipe) {
				t.Fatalf("dispatch returned %v, want write error", err)
			}
			if client.receivingTraffic.current != 4 {
				t.Fatalf("authenticated KeyUpdate did not install receive epoch: %d", client.receivingTraffic.current)
			}
		})
	}
}

func TestSendKeyUpdateIsIndependentOfNewSessionTicketFlight(t *testing.T) {
	client, server := establishedConnPair(t)
	originalClientConn := client.conn
	defer originalClientConn.Close()
	defer server.conn.Close()
	client.conn = &recordSinkConn{Conn: originalClientConn}
	client.ticketFlight = &flight{records: []flightRecord{{number: recordNumber{epoch: 3, sequence: 100}}}}
	if err := client.SendKeyUpdate(false); err != nil {
		t.Fatalf("KeyUpdate was blocked by independent NewSessionTicket flight: %v", err)
	}
	client.readerMu.Lock()
	client.readerClosed = true
	client.readerMu.Unlock()
}

func TestDeferredKeyUpdateResponseIsClearedAtSendingEpochLimit(t *testing.T) {
	suite, err := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	if err != nil {
		t.Fatal(err)
	}
	traffic, err := newSendingTraffic(suite, bytes.Repeat([]byte{1}, suite.hash.Size()), maxSendingEpoch-1, 0, 64)
	if err != nil {
		t.Fatal(err)
	}
	_, number, err := traffic.beginKeyUpdate(false)
	if err != nil {
		t.Fatal(err)
	}
	client, server := establishedConnPair(t)
	originalClientConn := client.conn
	defer originalClientConn.Close()
	defer server.conn.Close()
	client.conn = &recordSinkConn{Conn: originalClientConn}
	client.sendingTraffic = traffic
	client.sendCipher = traffic.cipher
	client.keyUpdateResponsePending = true
	ackSender, ackReceiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, maxSendingEpoch)
	client.receiveEpochs = newEpochSet()
	if err = client.receiveEpochs.install(ackReceiver); err != nil {
		t.Fatal(err)
	}
	if err = client.receiveEpochs.setCurrent(maxSendingEpoch); err != nil {
		t.Fatal(err)
	}
	ackBody, err := marshalACK([]recordNumber{number})
	if err != nil {
		t.Fatal(err)
	}
	ackWire, err := ackSender.seal(recordTypeACK, ackBody)
	if err != nil {
		t.Fatal(err)
	}
	if err = client.dispatchDatagram(ackWire); err != nil {
		t.Fatal(err)
	}
	if client.keyUpdateResponsePending || client.sendingTraffic.cipher.epoch != maxSendingEpoch {
		t.Fatalf("pending=%v epoch=%d", client.keyUpdateResponsePending, client.sendingTraffic.cipher.epoch)
	}
}

func TestConnRejectsKeyUpdateAfterIncompleteHandshakeMessage(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	sequence := client.sendingTraffic.messageSequence
	fragment, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeNewConnectionID, messageSequence: sequence, length: 2, body: []byte{0}})
	if err != nil {
		t.Fatal(err)
	}
	wire, err := client.sendCipher.seal(recordTypeHandshake, fragment)
	if err != nil {
		t.Fatal(err)
	}
	ackRead := make(chan error, 1)
	go func() {
		buffer := make([]byte, 2048)
		_, readErr := client.conn.Read(buffer)
		ackRead <- readErr
	}()
	if err = server.dispatchDatagram(wire); err != nil {
		t.Fatalf("buffering incomplete message: %v", err)
	}
	if err = <-ackRead; err != nil {
		t.Fatalf("reading fragment ACK: %v", err)
	}
	update, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeKeyUpdate, messageSequence: sequence + 1, length: 1, body: []byte{0}})
	if err != nil {
		t.Fatal(err)
	}
	wire, err = client.sendCipher.seal(recordTypeHandshake, update)
	if err != nil {
		t.Fatal(err)
	}
	err = server.dispatchDatagram(wire)
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertUnexpectedMessage {
		t.Fatalf("KeyUpdate after incomplete message returned %v", err)
	}
	if server.receivingTraffic.current != 3 {
		t.Fatalf("KeyUpdate after incomplete message advanced receive epoch to %d", server.receivingTraffic.current)
	}
}

func TestConnRejectsApplicationDataInsidePostHandshakeAuthResponse(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	if _, err := client.sendCipher.seal(recordTypeHandshake, []byte{1}); err != nil {
		t.Fatal(err)
	}
	nextRecord := client.sendCipher.nextSequence
	server.postHandshakeAuthState = &postHandshakeAuthState{
		stage: postAuthExpectCertificateVerify, hasResponseEpoch: true, responseEpoch: 3,
		firstResponseRecord: recordNumber{epoch: 3, sequence: nextRecord - 1},
		lastResponseRecord:  recordNumber{epoch: 3, sequence: nextRecord - 1},
	}
	wire, err := client.sendCipher.seal(recordTypeApplicationData, []byte("interleaved"))
	if err != nil {
		t.Fatal(err)
	}
	if err = server.dispatchDatagram(wire); err != nil {
		t.Fatalf("out-of-order application data was not buffered: %v", err)
	}
	server.postHandshakeAuthState.lastResponseRecord = recordNumber{epoch: 3, sequence: nextRecord + 1}
	err = validatePostHandshakeAuthApplicationOrder(server.postHandshakeAuthState)
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertUnexpectedMessage {
		t.Fatalf("interleaved application data returned %v", err)
	}
	server.postHandshakeAuthState.pendingApplication = []pendingPostAuthApplication{{
		number: recordNumber{epoch: 3, sequence: nextRecord + 2}, content: []byte("after response"),
	}}
	if err = validatePostHandshakeAuthApplicationOrder(server.postHandshakeAuthState); err != nil {
		t.Fatalf("reordered application data after Finished was rejected: %v", err)
	}
	server.postHandshakeAuthState = &postHandshakeAuthState{stage: postAuthExpectCertificate}
	wire, err = client.sendCipher.seal(recordTypeApplicationData, []byte("before response"))
	if err != nil || server.dispatchDatagram(wire) != nil {
		t.Fatalf("application data before PHA response was rejected: %v", err)
	}
}

func TestConnRejectsApplicationDataInsideFragmentedPostHandshakeMessage(t *testing.T) {
	client, server := establishedConnPair(t)
	originalClientConn := client.conn
	client.conn = &recordSinkConn{Conn: originalClientConn}
	defer originalClientConn.Close()
	defer server.conn.Close()

	sequence := server.sendingTraffic.messageSequence
	first, err := marshalHandshakeFragment(handshakeFragment{
		typ: handshakeTypeNewSessionTicket, messageSequence: sequence, length: 2, body: []byte{1},
	})
	if err != nil {
		t.Fatal(err)
	}
	wire, err := server.sendCipher.seal(recordTypeHandshake, first)
	if err != nil || client.dispatchDatagram(wire) != nil {
		t.Fatalf("first fragment: %v", err)
	}
	application, err := server.sendCipher.seal(recordTypeApplicationData, []byte("interleaved"))
	if err != nil || client.dispatchDatagram(application) != nil {
		t.Fatalf("buffering interleaved application data: %v", err)
	}
	exposed := bufferedApplicationDatagrams(client)
	if exposed != 0 {
		t.Fatal("application data was exposed before fragmented handshake completion")
	}
	second, err := marshalHandshakeFragment(handshakeFragment{
		typ: handshakeTypeNewSessionTicket, messageSequence: sequence, length: 2, offset: 1, body: []byte{2},
	})
	if err != nil {
		t.Fatal(err)
	}
	wire, err = server.sendCipher.seal(recordTypeHandshake, second)
	if err != nil {
		t.Fatal(err)
	}
	err = client.dispatchDatagram(wire)
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertUnexpectedMessage {
		t.Fatalf("fragment interleaving returned %v", err)
	}
}

func TestConnDeliversReorderedApplicationDataAfterFragmentedHandshakeMessage(t *testing.T) {
	client, server := establishedConnPair(t)
	originalClientConn := client.conn
	client.conn = &recordSinkConn{Conn: originalClientConn}
	defer originalClientConn.Close()
	defer server.conn.Close()

	sequence := server.sendingTraffic.messageSequence
	body, err := (&newSessionTicketMessage{lifetime: 1, nonce: []byte{1}, ticket: []byte{1}}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	split := len(body) / 2
	first, err := marshalHandshakeFragment(handshakeFragment{
		typ: handshakeTypeNewSessionTicket, messageSequence: sequence, length: uint32(len(body)), body: body[:split],
	})
	if err != nil {
		t.Fatal(err)
	}
	firstWire, err := server.sendCipher.seal(recordTypeHandshake, first)
	if err != nil {
		t.Fatal(err)
	}
	second, err := marshalHandshakeFragment(handshakeFragment{
		typ: handshakeTypeNewSessionTicket, messageSequence: sequence, length: uint32(len(body)), offset: uint32(split), body: body[split:],
	})
	if err != nil {
		t.Fatal(err)
	}
	secondWire, err := server.sendCipher.seal(recordTypeHandshake, second)
	if err != nil {
		t.Fatal(err)
	}
	applicationWire, err := server.sendCipher.seal(recordTypeApplicationData, []byte("after fragments"))
	if err != nil {
		t.Fatal(err)
	}
	if err = client.dispatchDatagram(firstWire); err != nil {
		t.Fatal(err)
	}
	if err = client.dispatchDatagram(applicationWire); err != nil {
		t.Fatalf("early arrival after final fragment: %v", err)
	}
	exposed := bufferedApplicationDatagrams(client)
	if exposed != 0 {
		t.Fatal("reordered application data was exposed before reassembly completed")
	}
	if err = client.dispatchDatagram(secondWire); err != nil {
		t.Fatalf("completing fragmented message: %v", err)
	}
	got := bufferedApplicationPayload(client)
	if string(got) != "after fragments" {
		t.Fatalf("delivered application data %q", got)
	}
}

func TestProtectedHandshakeApplicationRecordBufferIsCountBounded(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	server.config.ReplayWindow = 1
	server.config.MaxBufferedHandshakeMessages = 1
	server.postHandshakeReassembly = newReassembler()
	fragment := handshakeFragment{typ: handshakeTypeNewSessionTicket, length: 2, body: []byte{1}}
	if _, complete, err := server.postHandshakeReassembly.addProtected(fragment, 3); err != nil || complete {
		t.Fatalf("incomplete fragment: complete=%v err=%v", complete, err)
	}
	for sequence := uint64(0); sequence < 8; sequence++ {
		buffered, err := server.bufferIncompleteHandshakeApplicationLocked(nil, recordNumber{epoch: 3, sequence: sequence})
		if err != nil || !buffered {
			t.Fatalf("buffer %d: buffered=%v err=%v", sequence, buffered, err)
		}
	}
	if _, err := server.bufferIncompleteHandshakeApplicationLocked(nil, recordNumber{epoch: 3, sequence: 8}); err == nil {
		t.Fatal("zero-length application records exceeded the ordering buffer count limit")
	}

	server.postHandshakeAuthState = &postHandshakeAuthState{hasResponseEpoch: true, responseEpoch: 3}
	for sequence := uint64(0); sequence < 8; sequence++ {
		buffered, err := server.bufferPostHandshakeAuthApplicationLocked(nil, recordNumber{epoch: 3, sequence: sequence})
		if err != nil || !buffered {
			t.Fatalf("PHA buffer %d: buffered=%v err=%v", sequence, buffered, err)
		}
	}
	if _, err := server.bufferPostHandshakeAuthApplicationLocked(nil, recordNumber{epoch: 3, sequence: 8}); err == nil {
		t.Fatal("zero-length PHA application records exceeded the ordering buffer count limit")
	}
}

func TestCloseNotifyOnlyClosesPeerReadSide(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	body, err := (alertMessage{level: alertLevelWarning, description: alertCloseNotify}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	wire, err := client.sendCipher.seal(recordTypeAlert, body)
	if err != nil {
		t.Fatal(err)
	}
	readResult := make(chan error, 1)
	go func() {
		buffer := make([]byte, 1)
		_, _, readErr := server.ReadDatagram(buffer)
		readResult <- readErr
	}()
	if _, err = client.conn.Write(wire); err != nil {
		t.Fatal(err)
	}
	if err = <-readResult; !errors.Is(err, io.EOF) {
		t.Fatalf("Read returned %v after close_notify", err)
	}
	writeResult := make(chan error, 1)
	go func() {
		_, writeErr := server.WriteDatagram([]byte("still writable"))
		writeResult <- writeErr
	}()
	response := make([]byte, 2048)
	n, err := client.conn.Read(response)
	if err != nil {
		t.Fatal(err)
	}
	content, typ, _, _, err := client.receiveEpochs.open(response[:n])
	if err != nil || typ != recordTypeApplicationData || string(content) != "still writable" {
		t.Fatalf("type=%d content=%q err=%v", typ, content, err)
	}
	if err = <-writeResult; err != nil {
		t.Fatal(err)
	}
}

func TestCloseNotifyUsesRecordNumberOrdering(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	preClose, err := client.sendCipher.seal(recordTypeApplicationData, []byte("before"))
	if err != nil {
		t.Fatal(err)
	}
	body, err := (alertMessage{level: alertLevelWarning, description: alertCloseNotify}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	closeWire, err := client.sendCipher.seal(recordTypeAlert, body)
	if err != nil {
		t.Fatal(err)
	}
	postClose, err := client.sendCipher.seal(recordTypeApplicationData, []byte("after"))
	if err != nil {
		t.Fatal(err)
	}
	if err = server.dispatchDatagram(closeWire); err != nil {
		t.Fatal(err)
	}
	if err = server.dispatchDatagram(postClose); err != nil {
		t.Fatal(err)
	}
	if err = server.dispatchDatagram(preClose); err != nil {
		t.Fatal(err)
	}
	server.inputMu.Lock()
	closed := server.peerReadClosed
	server.inputMu.Unlock()
	buffered := string(bufferedApplicationPayload(server))
	if !closed || buffered != "before" {
		t.Fatalf("closed=%v buffered=%q; want only pre-close data", closed, buffered)
	}
}
