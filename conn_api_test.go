package dtls13

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"testing"
)

type localAddrConn struct {
	net.Conn
	local net.Addr
}

func (c *localAddrConn) LocalAddr() net.Addr { return c.local }

func TestDialRejectsStreamNetwork(t *testing.T) {
	if _, err := Dial("tcp", "127.0.0.1:4433", &Config{}); err == nil {
		t.Fatal("accepted TCP network")
	}
}
func TestNilUnderlyingConnection(t *testing.T) {
	if err := Client(nil, &Config{}).Handshake(); err == nil {
		t.Fatal("accepted nil connection")
	}
}

func TestConnUsesOnlyNativeDatagramAPI(t *testing.T) {
	var conn any = (*Conn)(nil)
	if _, ok := conn.(net.Conn); ok {
		t.Fatal("Conn unexpectedly implements net.Conn")
	}
	if _, ok := conn.(net.PacketConn); ok {
		t.Fatal("Conn unexpectedly implements net.PacketConn")
	}
}

func TestPathMTUAndRecordOverhead(t *testing.T) {
	client, _ := establishedConnPair(t)
	if client.PathMTU() != 1200 {
		t.Fatalf("PathMTU=%d", client.PathMTU())
	}
	want := client.sendCipher.headerLen16() + client.sendCipher.aead.Overhead() + 1
	if got := client.RecordOverhead(); got != want {
		t.Fatalf("RecordOverhead=%d, want %d", got, want)
	}
}

func TestPathMTUFloorUsesUDPPayloadSize(t *testing.T) {
	for _, test := range []struct {
		ip   net.IP
		want int
	}{{net.IPv4(127, 0, 0, 1), 548}, {net.ParseIP("::1"), 1232}} {
		conn := &Conn{conn: &localAddrConn{local: &net.UDPAddr{IP: test.ip}}}
		if got := conn.pathMTUFloor(); got != test.want {
			t.Fatalf("IP=%s payload floor=%d, want %d", test.ip, got, test.want)
		}
	}
}

func TestConnDatagramBoundariesAndSources(t *testing.T) {
	client, server := establishedConnPair(t)

	cases := []struct {
		payload  []byte
		readSize int
	}{
		{[]byte("first"), 1024},
		{bytes.Repeat([]byte("x"), 700), 1024},
		{nil, 0},
		{[]byte("abcdef"), 3},
		{[]byte("next"), 1024},
		{bytes.Repeat([]byte("y"), 1000), 1024},
	}
	writes := make(chan error, 1)
	go func() {
		for _, test := range cases {
			payload := test.payload
			n, err := client.WriteDatagram(payload)
			if err == nil && n != len(payload) {
				err = fmt.Errorf("WriteDatagram count=%d, want %d", n, len(payload))
			}
			if err != nil {
				writes <- err
				return
			}
		}
		writes <- nil
	}()
	buffer := make([]byte, 1024)
	for _, test := range cases {
		want := test.payload
		destination := buffer[:test.readSize]
		if test.readSize == 0 {
			destination = nil
		}
		n, info, err := server.ReadDatagram(destination)
		if err != nil {
			t.Fatal(err)
		}
		if n != min(len(destination), len(want)) || !bytes.Equal(buffer[:n], want[:n]) {
			t.Fatalf("ReadDatagram=%q, want %q", buffer[:n], want)
		}
		if info.Source == nil || info.Source.String() != server.RemoteAddr().String() || info.FullLength != len(want) || info.Truncated != (len(want) > len(destination)) {
			t.Fatalf("info=%+v, want source %v and length %d", info, server.RemoteAddr(), len(want))
		}
	}
	if err := <-writes; err != nil {
		t.Fatal(err)
	}
}

func TestConnWriteDatagramRejectsOversizedDatagram(t *testing.T) {
	client, _ := establishedConnPair(t)
	sink := &recordSinkConn{Conn: client.conn}
	client.conn = sink
	maximum := client.maxApplicationDatagramLocked()
	n, err := client.WriteDatagram(bytes.Repeat([]byte{1}, maximum+1))
	if n != 0 || !errors.Is(err, ErrDatagramTooLarge) {
		t.Fatalf("WriteDatagram oversized=%d, %v", n, err)
	}
	if len(sink.writes) != 0 {
		t.Fatalf("oversized datagram generated %d records", len(sink.writes))
	}
}
