package dtls13

import (
	"bytes"
	"errors"
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
	client, server := establishedConnPair(t)
	if client.PathMTU() != 1200 {
		t.Fatalf("PathMTU=%d", client.PathMTU())
	}
	want := client.sendCipher.headerLen16() + client.sendCipher.aead.Overhead() + 1
	if got := client.RecordOverhead(); got != want {
		t.Fatalf("RecordOverhead=%d, want %d", got, want)
	}
	_ = client.conn.Close()
	_ = server.conn.Close()
}

func TestPathMTUFloorUsesUDPPayloadSize(t *testing.T) {
	left, right := net.Pipe()
	defer left.Close()
	defer right.Close()
	for _, test := range []struct {
		ip   net.IP
		want int
	}{{net.IPv4(127, 0, 0, 1), 548}, {net.ParseIP("::1"), 1232}} {
		conn := &Conn{conn: &localAddrConn{Conn: left, local: &net.UDPAddr{IP: test.ip}}}
		if got := conn.pathMTUFloor(); got != test.want {
			t.Fatalf("IP=%s payload floor=%d, want %d", test.ip, got, test.want)
		}
	}
}

func TestConnDatagramBoundariesAndSources(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()

	writes := make(chan error, 1)
	go func() {
		if _, err := client.WriteDatagram([]byte("first")); err != nil {
			writes <- err
			return
		}
		_, err := client.WriteDatagram([]byte("second"))
		writes <- err
	}()
	buffer := make([]byte, 32)
	for _, want := range []string{"first", "second"} {
		n, info, err := server.ReadDatagram(buffer)
		if err != nil {
			t.Fatal(err)
		}
		if string(buffer[:n]) != want {
			t.Fatalf("ReadDatagram=%q, want %q", buffer[:n], want)
		}
		if info.Source == nil || info.Source.String() != server.RemoteAddr().String() || info.FullLength != len(want) || info.Truncated {
			t.Fatalf("info=%+v, want source %v and length %d", info, server.RemoteAddr(), len(want))
		}
	}
	if err := <-writes; err != nil {
		t.Fatal(err)
	}
}

func TestConnReadDatagramReportsTruncation(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	writes := make(chan error, 1)
	go func() {
		if _, err := client.WriteDatagram([]byte("abcdef")); err != nil {
			writes <- err
			return
		}
		_, err := client.WriteDatagram([]byte("next"))
		writes <- err
	}()
	buffer := make([]byte, 3)
	n, info, err := server.ReadDatagram(buffer)
	if err != nil || string(buffer[:n]) != "abc" || !info.Truncated || info.FullLength != 6 {
		t.Fatalf("first ReadDatagram=%q, %+v, %v", buffer[:n], info, err)
	}
	buffer = make([]byte, 8)
	n, info, err = server.ReadDatagram(buffer)
	if err != nil || string(buffer[:n]) != "next" || info.Truncated || info.FullLength != 4 {
		t.Fatalf("second ReadDatagram=%q, %+v, %v", buffer[:n], info, err)
	}
	if err := <-writes; err != nil {
		t.Fatal(err)
	}
}

func TestConnZeroLengthDatagram(t *testing.T) {
	client, server := establishedConnPair(t)
	defer client.conn.Close()
	defer server.conn.Close()
	written := make(chan error, 1)
	go func() {
		n, err := client.WriteDatagram(nil)
		if err == nil && n != 0 {
			err = errors.New("zero-length WriteDatagram returned a non-zero count")
		}
		written <- err
	}()
	n, info, err := server.ReadDatagram(nil)
	if err != nil || n != 0 || info.Source == nil || info.FullLength != 0 || info.Truncated {
		t.Fatalf("ReadDatagram(nil)=%d, %+v, %v", n, info, err)
	}
	if err := <-written; err != nil {
		t.Fatal(err)
	}
}

func TestConnWriteDatagramRejectsOversizedDatagram(t *testing.T) {
	client, server := establishedConnPair(t)
	original := client.conn
	defer original.Close()
	defer server.conn.Close()
	sink := &recordSinkConn{Conn: original}
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
