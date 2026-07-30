package dtls13

import (
	"bytes"
	"errors"
	"io"
	"net"
	"reflect"
	"testing"
)

type failingCIDRegistrarConn struct {
	net.Conn
	registered [][]byte
	removed    [][]byte
	failWrite  bool
}

func (c *failingCIDRegistrarConn) registerConnectionIDs(cids [][]byte) error {
	c.registered = cloneConnectionIDs(cids)
	return nil
}

func (c *failingCIDRegistrarConn) unregisterConnectionIDs(cids [][]byte) {
	c.removed = cloneConnectionIDs(cids)
	c.registered = nil
}

func (c *failingCIDRegistrarConn) Write(p []byte) (int, error) {
	if c.failWrite {
		return 0, io.ErrClosedPipe
	}
	return c.Conn.Write(p)
}

func TestNewConnectionIDRoundTrip(t *testing.T) {
	for _, want := range []newConnectionIDMessage{
		{connectionIDs: [][]byte{{1, 2}, {}, {3, 4, 5}}, usage: connectionIDImmediate},
		{usage: connectionIDSpare},
	} {
		b, err := want.marshal()
		if err != nil {
			t.Fatal(err)
		}
		got, err := parseNewConnectionID(b)
		if err != nil {
			t.Fatal(err)
		}
		if got.usage != want.usage || !reflect.DeepEqual(got.connectionIDs, want.connectionIDs) {
			t.Fatalf("got %#v want %#v", got, want)
		}
	}
}

func TestNewConnectionIDRejectsOversizedCID(t *testing.T) {
	message := newConnectionIDMessage{connectionIDs: [][]byte{make([]byte, 256)}, usage: connectionIDSpare}
	if wire, err := message.marshal(); err == nil {
		t.Fatalf("marshaled 256-byte CID as %x", wire)
	}
}

func TestNewConnectionIDRejectsOversizedVector(t *testing.T) {
	connectionIDs := make([][]byte, 256)
	for i := range connectionIDs {
		connectionIDs[i] = make([]byte, 255)
	}
	if wire, err := (newConnectionIDMessage{connectionIDs: connectionIDs, usage: connectionIDSpare}).marshal(); err == nil {
		t.Fatalf("marshaled oversized CID vector as %d bytes", len(wire))
	}
}

func TestConnectionIDUpdateMessagesRejectMalformed(t *testing.T) {
	if _, err := (newConnectionIDMessage{usage: connectionIDImmediate}).marshal(); err == nil {
		t.Fatal("marshaled an empty immediate NewConnectionId")
	}
	for _, b := range [][]byte{{}, {0, 0}, {0, 0, 2}, {0, 1, 2, 1}} {
		if _, err := parseNewConnectionID(b); err == nil {
			t.Fatalf("accepted malformed NewConnectionId %x", b)
		}
	}
	if _, err := parseRequestConnectionID(nil); err == nil {
		t.Fatal("accepted a truncated RequestConnectionId")
	}
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	client.connectionIDNegotiated = true
	client.peerCIDUpdatesAllowed = true
	client.sendConnectionID = []byte{1}
	if err := client.processNewConnectionID(1, nil); err == nil {
		t.Fatal("accepted truncated NewConnectionId")
	} else if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
		t.Fatalf("truncated NewConnectionId alert=%d ok=%v err=%v", description, ok, err)
	}
	if _, _, err := client.processRequestConnectionID(1, nil); err == nil {
		t.Fatal("accepted truncated RequestConnectionId")
	} else if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
		t.Fatalf("truncated RequestConnectionId alert=%d ok=%v err=%v", description, ok, err)
	}
	for _, count := range []uint8{0, 1, 255} {
		got, err := parseRequestConnectionID((requestConnectionIDMessage{count: count}).marshal())
		if err != nil || got.count != count || !bytes.Equal(got.marshal(), []byte{count}) {
			t.Fatalf("count=%d got=%#v err=%v", count, got, err)
		}
	}
}

func TestConnectionIDUpdateRequiresNegotiation(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	if err := client.SendNewConnectionIDs([][]byte{{1}}, false); err == nil {
		t.Fatal("sent NewConnectionId without CID negotiation")
	}
	if err := client.RequestConnectionIDs(1); err == nil {
		t.Fatal("sent RequestConnectionId without CID negotiation")
	}
	body, err := (newConnectionIDMessage{connectionIDs: [][]byte{{1}}, usage: connectionIDSpare}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err = server.processNewConnectionID(1, body); err == nil {
		t.Fatal("accepted NewConnectionId without CID negotiation")
	}
	if _, _, err = server.processRequestConnectionID(1, []byte{1}); err == nil {
		t.Fatal("accepted RequestConnectionId without CID negotiation")
	}
}

func TestConnectionIDResourceLimits(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	client.connectionIDNegotiated = true
	client.localCIDUpdatesAllowed = true
	client.receiveConnectionID = []byte{1, 2}
	client.config.MaxConnectionIDs = 2
	if err := client.receivingTraffic.setConnectionID(client.receiveConnectionID); err != nil {
		t.Fatal(err)
	}
	if err := client.SendNewConnectionIDs([][]byte{{3, 4}, {5, 6}}, false); err == nil {
		t.Fatal("local connection ID set exceeded MaxConnectionIDs")
	}

	server.connectionIDNegotiated = true
	server.peerCIDUpdatesAllowed = true
	server.sendConnectionID = []byte{7, 8}
	server.config.MaxConnectionIDs = 2
	body, err := (newConnectionIDMessage{connectionIDs: [][]byte{{9, 10}, {11, 12}, {13, 14}}, usage: connectionIDSpare}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err = server.processNewConnectionID(1, body); err != nil {
		t.Fatal(err)
	}
	if len(server.peerSpareConnectionIDs) != 1 {
		t.Fatalf("retained %d spare peer CIDs, want 1", len(server.peerSpareConnectionIDs))
	}
}

func TestSendNewConnectionIDsRollsBackOnWriteFailure(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	client.connectionIDNegotiated = true
	client.localCIDUpdatesAllowed = true
	client.receiveConnectionID = []byte{1, 2}
	if err := client.receivingTraffic.setConnectionID(client.receiveConnectionID); err != nil {
		t.Fatal(err)
	}
	sequence := client.sendingTraffic.messageSequence
	wrapped := &failingCIDRegistrarConn{Conn: client.conn, failWrite: true}
	client.conn = wrapped
	if err := client.SendNewConnectionIDs([][]byte{{1, 2}, {3, 4}}, false); !errors.Is(err, io.ErrClosedPipe) {
		t.Fatalf("SendNewConnectionIDs error = %v", err)
	}
	if got := client.receivingTraffic.acceptedCIDs; !reflect.DeepEqual(got, [][]byte{{1, 2}}) {
		t.Fatalf("accepted CIDs after rollback = %x", got)
	}
	if len(wrapped.registered) != 0 {
		t.Fatalf("registered CIDs after rollback = %x", wrapped.registered)
	}
	if !reflect.DeepEqual(wrapped.removed, [][]byte{{3, 4}}) {
		t.Fatalf("rollback removed CIDs = %x, want only the newly registered CID", wrapped.removed)
	}
	if client.sendingTraffic.messageSequence != sequence {
		t.Fatalf("message sequence advanced from %d to %d", sequence, client.sendingTraffic.messageSequence)
	}
	if client.newConnectionIDFlight != nil {
		t.Fatal("failed NewConnectionId left an outstanding flight")
	}
}

func TestPeerConnectionIDSetRejectsPrefixesAndBoundsStorage(t *testing.T) {
	client, server := establishedConnPair(t)
	client.connectionIDNegotiated = true
	client.peerCIDUpdatesAllowed = true
	client.sendConnectionID = []byte{1, 2}
	client.config.MaxConnectionIDs = 2

	prefixBody, err := (newConnectionIDMessage{usage: connectionIDSpare, connectionIDs: [][]byte{{1, 2, 3}}}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	err = client.processNewConnectionID(1, prefixBody)
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertIllegalParameter {
		t.Fatalf("prefix collision returned %v", err)
	}

	overflowBody, err := (newConnectionIDMessage{usage: connectionIDSpare, connectionIDs: [][]byte{{3}, {4}}}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	err = client.processNewConnectionID(2, overflowBody)
	if err != nil {
		t.Fatalf("bounded CID advertisement returned %v", err)
	}
	if len(client.peerSpareConnectionIDs) != 1 {
		t.Fatalf("retained %d spare CIDs, want 1", len(client.peerSpareConnectionIDs))
	}
	_ = client.conn.Close()
	_ = server.conn.Close()
}
