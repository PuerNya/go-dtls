package dtls13

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

type returnRoutabilityTestAddr string

func (a returnRoutabilityTestAddr) Network() string { return "rrc-test" }
func (a returnRoutabilityTestAddr) String() string  { return string(a) }

type returnRoutabilityTestWrite struct {
	wire []byte
	to   net.Addr
}

type returnRoutabilityTestTransport struct {
	mu     sync.Mutex
	remote net.Addr
	writes []returnRoutabilityTestWrite
}

func (*returnRoutabilityTestTransport) Read([]byte) (int, error) { return 0, io.EOF }
func (t *returnRoutabilityTestTransport) Write(wire []byte) (int, error) {
	return t.WriteTo(wire, t.RemoteAddr())
}
func (t *returnRoutabilityTestTransport) WriteTo(wire []byte, address net.Addr) (int, error) {
	t.mu.Lock()
	t.writes = append(t.writes, returnRoutabilityTestWrite{wire: append([]byte(nil), wire...), to: address})
	t.mu.Unlock()
	return len(wire), nil
}
func (*returnRoutabilityTestTransport) Close() error { return nil }
func (*returnRoutabilityTestTransport) LocalAddr() net.Addr {
	return returnRoutabilityTestAddr("local")
}
func (*returnRoutabilityTestTransport) SetDeadline(time.Time) error      { return nil }
func (*returnRoutabilityTestTransport) SetReadDeadline(time.Time) error  { return nil }
func (*returnRoutabilityTestTransport) SetWriteDeadline(time.Time) error { return nil }
func (t *returnRoutabilityTestTransport) RemoteAddr() net.Addr {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.remote
}
func (t *returnRoutabilityTestTransport) rebindRemote(address net.Addr) {
	t.mu.Lock()
	t.remote = address
	t.mu.Unlock()
}
func (t *returnRoutabilityTestTransport) writeCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.writes)
}
func (t *returnRoutabilityTestTransport) write(index int) returnRoutabilityTestWrite {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.writes[index]
}

func newReturnRoutabilityTestConn(t *testing.T, random []byte) (*Conn, *returnRoutabilityTestTransport, *recordCipher) {
	t.Helper()
	oldAddress := returnRoutabilityTestAddr("old")
	transport := &returnRoutabilityTestTransport{remote: oldAddress}
	config, err := (&Config{Rand: bytes.NewReader(random)}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	suite, err := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	if err != nil {
		t.Fatal(err)
	}
	secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
	sender, err := newRecordCipher(suite, secret, 3, 64)
	if err != nil {
		t.Fatal(err)
	}
	receiver, err := newRecordCipher(suite, secret, 3, 64)
	if err != nil {
		t.Fatal(err)
	}
	for _, cipher := range []*recordCipher{sender, receiver} {
		if err = cipher.setConnectionID([]byte{0xc1}); err != nil {
			t.Fatal(err)
		}
	}
	conn := &Conn{conn: transport, config: config, sendCipher: sender, returnRoutabilityCheckNegotiated: true}
	return conn, transport, receiver
}

func openReturnRoutabilityTestWrite(t *testing.T, receiver *recordCipher, write returnRoutabilityTestWrite) returnRoutabilityMessage {
	t.Helper()
	content, contentType, _, err := receiver.open(write.wire)
	if err != nil {
		t.Fatal(err)
	}
	if contentType != recordTypeReturnRoutability {
		t.Fatalf("content type = %d", contentType)
	}
	message, known, err := parseReturnRoutabilityMessage(content)
	if err != nil || !known {
		t.Fatalf("message known=%v err=%v", known, err)
	}
	return message
}

func returnRoutabilityPhaseOf(conn *Conn) returnRoutabilityPhase {
	conn.writeMu.Lock()
	defer conn.writeMu.Unlock()
	if conn.returnRoutability == nil {
		return returnRoutabilityIdle
	}
	return conn.returnRoutability.phase
}

func TestReturnRoutabilityMessageCodec(t *testing.T) {
	for _, typ := range []uint8{returnRoutabilityPathChallenge, returnRoutabilityPathResponse, returnRoutabilityPathDrop} {
		want := returnRoutabilityMessage{typ: typ, cookie: [8]byte{1, 2, 3, 4, 5, 6, 7, 8}}
		wire, err := want.marshal()
		if err != nil {
			t.Fatal(err)
		}
		got, known, err := parseReturnRoutabilityMessage(wire[:])
		if err != nil || !known || got != want {
			t.Fatalf("type=%d got=%#v known=%v err=%v", typ, got, known, err)
		}
	}
	for _, wire := range [][]byte{nil, {returnRoutabilityPathChallenge}, {returnRoutabilityPathResponse, 1, 2}} {
		if _, _, err := parseReturnRoutabilityMessage(wire); err == nil {
			t.Fatalf("accepted malformed message %x", wire)
		}
	}
	for _, wire := range [][]byte{{0xff}, {3, 1, 2, 3}} {
		if _, known, err := parseReturnRoutabilityMessage(wire); err != nil || known {
			t.Fatalf("unknown message %x known=%v err=%v", wire, known, err)
		}
	}
	if _, err := (returnRoutabilityMessage{typ: 3}).marshal(); err == nil {
		t.Fatal("marshaled an unknown local message type")
	}
}

func TestReturnRoutabilityResponderRepliesExactlyOnce(t *testing.T) {
	conn, transport, receiver := newReturnRoutabilityTestConn(t, bytes.Repeat([]byte{1}, 32))
	cookie := [8]byte{8, 7, 6, 5, 4, 3, 2, 1}
	challenge, _ := (returnRoutabilityMessage{typ: returnRoutabilityPathChallenge, cookie: cookie}).marshal()
	oldAddress := transport.RemoteAddr()
	if err := conn.handleReturnRoutability(challenge[:], oldAddress); err != nil {
		t.Fatal(err)
	}
	first := transport.write(0)
	if !sameNetworkAddress(first.to, oldAddress) {
		t.Fatalf("response sent to %v", first.to)
	}
	if got := openReturnRoutabilityTestWrite(t, receiver, first); got.typ != returnRoutabilityPathResponse || got.cookie != cookie {
		t.Fatalf("old-path response = %#v", got)
	}

	newAddress := returnRoutabilityTestAddr("new")
	if err := conn.handleReturnRoutability(challenge[:], newAddress); err != nil {
		t.Fatal(err)
	}
	second := transport.write(1)
	if !sameNetworkAddress(second.to, newAddress) {
		t.Fatalf("drop sent to %v", second.to)
	}
	if got := openReturnRoutabilityTestWrite(t, receiver, second); got.typ != returnRoutabilityPathDrop || got.cookie != cookie {
		t.Fatalf("non-preferred-path response = %#v", got)
	}

	before := transport.writeCount()
	for _, content := range [][]byte{{0xff}, {returnRoutabilityPathChallenge}} {
		if err := conn.handleReturnRoutability(content, oldAddress); err != nil {
			t.Fatal(err)
		}
	}
	if transport.writeCount() != before {
		t.Fatal("invalid or unknown message elicited a response")
	}
	conn.returnRoutabilityCheckNegotiated = false
	if err := conn.handleReturnRoutability(challenge[:], oldAddress); err != nil || transport.writeCount() != before {
		t.Fatal("RRC was used without negotiation")
	}
}

func TestEnhancedReturnRoutabilityKeepsReachableOldPath(t *testing.T) {
	random := append(bytes.Repeat([]byte{1}, 8), bytes.Repeat([]byte{2}, 8)...)
	conn, transport, _ := newReturnRoutabilityTestConn(t, random)
	conn.retransmitNanos.Store(int64(time.Hour))
	conn.lastRTTSampleUnixNano.Store(time.Now().UnixNano())
	defer conn.clearReturnRoutability()
	newAddress := returnRoutabilityTestAddr("new")
	if err := conn.observeReturnRoutabilityRecord(newAddress, 64); err != nil {
		t.Fatal(err)
	}
	if transport.writeCount() != 1 || !sameNetworkAddress(transport.write(0).to, returnRoutabilityTestAddr("old")) {
		t.Fatal("enhanced check did not challenge the old path first")
	}
	cookie := conn.returnRoutability.cookie
	wrong := cookie
	wrong[0] ^= 1
	wrongResponse, _ := (returnRoutabilityMessage{typ: returnRoutabilityPathResponse, cookie: wrong}).marshal()
	if err := conn.handleReturnRoutability(wrongResponse[:], transport.RemoteAddr()); err != nil {
		t.Fatal(err)
	}
	if conn.returnRoutability.phase != returnRoutabilityCheckingOldPath {
		t.Fatal("invalid old-path response changed validation state")
	}
	response, _ := (returnRoutabilityMessage{typ: returnRoutabilityPathResponse, cookie: cookie}).marshal()
	if err := conn.handleReturnRoutability(response[:], transport.RemoteAddr()); err != nil {
		t.Fatal(err)
	}
	if conn.returnRoutability.phase != returnRoutabilityIdle || !sameNetworkAddress(transport.RemoteAddr(), returnRoutabilityTestAddr("old")) {
		t.Fatal("reachable old path was not retained")
	}
}

func TestReturnRoutabilityRequiresFreshRandomCookie(t *testing.T) {
	conn, transport, _ := newReturnRoutabilityTestConn(t, nil)
	if err := conn.observeReturnRoutabilityRecord(returnRoutabilityTestAddr("new"), 64); err == nil {
		t.Fatal("started path validation without CSPRNG output")
	}
	if returnRoutabilityPhaseOf(conn) != returnRoutabilityIdle || transport.writeCount() != 0 {
		t.Fatal("randomness failure left an active or transmitted challenge")
	}
}

func TestEnhancedReturnRoutabilityDropFallsBackToBasic(t *testing.T) {
	random := append(bytes.Repeat([]byte{1}, 8), bytes.Repeat([]byte{2}, 8)...)
	conn, transport, receiver := newReturnRoutabilityTestConn(t, random)
	conn.retransmitNanos.Store(int64(time.Hour))
	conn.lastRTTSampleUnixNano.Store(time.Now().UnixNano())
	defer conn.clearReturnRoutability()
	oldAddress := transport.RemoteAddr()
	newAddress := returnRoutabilityTestAddr("new")
	spareConnectionID := []byte{0xd1, 0xd2}
	conn.sendingTraffic = &sendingTraffic{cipher: conn.sendCipher}
	conn.peerSpareConnectionIDs = [][]byte{spareConnectionID}
	if err := receiver.addAcceptedConnectionIDs([][]byte{spareConnectionID}); err != nil {
		t.Fatal(err)
	}
	if err := conn.observeReturnRoutabilityRecord(newAddress, 64); err != nil {
		t.Fatal(err)
	}
	oldCookie := conn.returnRoutability.cookie
	drop, _ := (returnRoutabilityMessage{typ: returnRoutabilityPathDrop, cookie: oldCookie}).marshal()
	if err := conn.handleReturnRoutability(drop[:], oldAddress); err != nil {
		t.Fatal(err)
	}
	if conn.returnRoutability.phase != returnRoutabilityCheckingNewPath || conn.returnRoutability.cookie == oldCookie {
		t.Fatal("path_drop did not start a basic check with a fresh cookie")
	}
	if transport.writeCount() != 2 || !sameNetworkAddress(transport.write(1).to, newAddress) {
		t.Fatal("basic challenge was not sent to the candidate path")
	}
	if got := openReturnRoutabilityTestWrite(t, receiver, transport.write(1)); got.typ != returnRoutabilityPathChallenge || !equalBytes(receiver.lastConnectionID, spareConnectionID) {
		t.Fatal("candidate path challenge did not use the available spare CID")
	}
	if !equalBytes(conn.sendCipher.connectionID, []byte{0xc1}) {
		t.Fatal("candidate probe changed the CID used on the old path before validation")
	}
	if conn.returnRoutability.sent > 3*conn.returnRoutability.received {
		t.Fatal("candidate path exceeded the anti-amplification limit")
	}
	cookie := conn.returnRoutability.cookie
	delayedOldResponse, _ := (returnRoutabilityMessage{typ: returnRoutabilityPathResponse, cookie: oldCookie}).marshal()
	if err := conn.handleReturnRoutability(delayedOldResponse[:], oldAddress); err != nil {
		t.Fatal(err)
	}
	if conn.returnRoutability.phase != returnRoutabilityCheckingNewPath {
		t.Fatal("delayed old-path response replaced the active basic check")
	}
	response, _ := (returnRoutabilityMessage{typ: returnRoutabilityPathResponse, cookie: cookie}).marshal()
	if err := conn.handleReturnRoutability(response[:], returnRoutabilityTestAddr("other")); err != nil {
		t.Fatal(err)
	}
	if !sameNetworkAddress(transport.RemoteAddr(), oldAddress) {
		t.Fatal("response from the wrong address rebound the connection")
	}
	if err := conn.handleReturnRoutability(response[:], newAddress); err != nil {
		t.Fatal(err)
	}
	if !sameNetworkAddress(transport.RemoteAddr(), newAddress) || conn.returnRoutability.phase != returnRoutabilityIdle {
		t.Fatal("valid basic response did not rebind the connection")
	}
	if !equalBytes(conn.sendCipher.connectionID, spareConnectionID) || len(conn.peerSpareConnectionIDs) != 0 {
		t.Fatal("validated new path did not activate the available spare CID")
	}
	if err := conn.handleReturnRoutability(response[:], newAddress); err != nil || conn.returnRoutability.phase != returnRoutabilityIdle {
		t.Fatal("duplicate response changed completed validation")
	}
}

func TestReturnRoutabilityTimeoutAndNestedRebinding(t *testing.T) {
	conn, transport, _ := newReturnRoutabilityTestConn(t, bytes.Repeat([]byte{3}, 32))
	conn.retransmitNanos.Store(int64(2 * time.Millisecond))
	conn.lastRTTSampleUnixNano.Store(time.Now().UnixNano())
	firstAddress := returnRoutabilityTestAddr("first")
	if err := conn.observeReturnRoutabilityRecord(firstAddress, 64); err != nil {
		t.Fatal(err)
	}
	if err := conn.observeReturnRoutabilityRecord(returnRoutabilityTestAddr("nested"), 64); err != nil {
		t.Fatal(err)
	}
	if !sameNetworkAddress(conn.returnRoutability.newAddress, firstAddress) {
		t.Fatal("nested rebinding replaced the active candidate")
	}
	deadline := time.Now().Add(200 * time.Millisecond)
	for time.Now().Before(deadline) && returnRoutabilityPhaseOf(conn) != returnRoutabilityIdle {
		time.Sleep(time.Millisecond)
	}
	phase := returnRoutabilityPhaseOf(conn)
	if phase != returnRoutabilityIdle || transport.writeCount() != 2 {
		t.Fatalf("timed-out validation phase=%d writes=%d", phase, transport.writeCount())
	}
	if !sameNetworkAddress(transport.RemoteAddr(), returnRoutabilityTestAddr("old")) {
		t.Fatal("timeout rebound the connection")
	}
	if err := conn.observeReturnRoutabilityRecord(firstAddress, 64); err != nil {
		t.Fatal(err)
	}
	if returnRoutabilityPhaseOf(conn) != returnRoutabilityCheckingOldPath || transport.writeCount() != 3 {
		t.Fatal("new candidate data did not restart validation after packet loss")
	}
	conn.clearReturnRoutability()
}

func TestReturnRoutabilityTimerAndAmplificationLimits(t *testing.T) {
	conn, _, _ := newReturnRoutabilityTestConn(t, bytes.Repeat([]byte{1}, 16))
	if got := conn.returnRoutabilityTimeout(); got != time.Second {
		t.Fatalf("fallback timeout = %v", got)
	}
	conn.retransmitNanos.Store(int64(4 * time.Millisecond))
	conn.lastRTTSampleUnixNano.Store(time.Now().UnixNano())
	if got := conn.returnRoutabilityTimeout(); got != 8*time.Millisecond {
		t.Fatalf("3xRTT timeout = %v", got)
	}
	conn.returnRoutability = &returnRoutabilityState{phase: returnRoutabilityCheckingNewPath, newAddress: returnRoutabilityTestAddr("new"), received: 10}
	if !conn.allowReturnRoutabilitySendLocked(returnRoutabilityTestAddr("new"), 30) || conn.allowReturnRoutabilitySendLocked(returnRoutabilityTestAddr("new"), 31) {
		t.Fatal("anti-amplification boundary is not exactly three times received bytes")
	}
	conn.returnRoutability.sent = 20
	if !conn.allowReturnRoutabilitySendLocked(returnRoutabilityTestAddr("new"), 10) || conn.allowReturnRoutabilitySendLocked(returnRoutabilityTestAddr("new"), 11) {
		t.Fatal("anti-amplification accounting ignored previously sent bytes")
	}
}

func TestProtectedReturnRoutabilityRecord(t *testing.T) {
	sender, receiver := recordCipherPair(t, TLS_AES_128_GCM_SHA256, 3)
	payload := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
	wire, err := sender.seal(recordTypeReturnRoutability, payload)
	if err != nil {
		t.Fatal(err)
	}
	content, contentType, _, err := receiver.open(wire)
	if err != nil || contentType != recordTypeReturnRoutability || !bytes.Equal(content, payload) {
		t.Fatalf("content=%x type=%d err=%v", content, contentType, err)
	}
}
