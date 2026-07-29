package dtls13

import (
	"testing"
	"time"
)

func TestAmplificationGuard(t *testing.T) {
	var g amplificationGuard
	if g.allowSend(1) {
		t.Fatal("sent before receiving data")
	}
	g.recordReceived(100)
	if !g.allowSend(250) {
		t.Fatal("rejected within budget")
	}
	if g.allowSend(51) {
		t.Fatal("exceeded three-times budget")
	}
	if !g.allowSend(50) {
		t.Fatal("rejected exact remaining budget")
	}
	g.validate()
	if !g.allowSend(1 << 20) {
		t.Fatal("limited validated address")
	}
}

func TestAmplificationConnAccountsActualDatagrams(t *testing.T) {
	left, right := memoryDatagramPair()
	var guard amplificationGuard
	conn := &amplificationConn{Conn: left, guard: &guard}
	if n, err := conn.Write([]byte("suppressed")); err != nil || n != len("suppressed") {
		t.Fatalf("suppressed write n=%d err=%v", n, err)
	}
	_ = right.SetReadDeadline(time.Now().Add(5 * time.Millisecond))
	if _, err := right.Read(make([]byte, 32)); err == nil {
		t.Fatal("write before receiving data reached the network")
	}
	_ = right.SetReadDeadline(time.Time{})
	inbound := make([]byte, 100)
	if _, err := right.Write(inbound); err != nil {
		t.Fatal(err)
	}
	if n, err := conn.Read(make([]byte, 100)); err != nil || n != 100 {
		t.Fatalf("accounted read n=%d err=%v", n, err)
	}
	if _, err := conn.Write(make([]byte, 300)); err != nil {
		t.Fatal(err)
	}
	_ = right.SetReadDeadline(time.Now().Add(time.Second))
	if n, err := right.Read(make([]byte, 300)); err != nil || n != 300 {
		t.Fatalf("budgeted write n=%d err=%v", n, err)
	}
	if _, err := conn.Write([]byte{1}); err != nil {
		t.Fatal(err)
	}
	_ = right.SetReadDeadline(time.Now().Add(5 * time.Millisecond))
	if _, err := right.Read(make([]byte, 1)); err == nil {
		t.Fatal("write beyond amplification budget reached the network")
	}
	guard.validate()
	_ = right.SetReadDeadline(time.Time{})
	if _, err := conn.Write([]byte{2}); err != nil {
		t.Fatal(err)
	}
	if n, err := right.Read(make([]byte, 1)); err != nil || n != 1 {
		t.Fatalf("validated write n=%d err=%v", n, err)
	}
}
