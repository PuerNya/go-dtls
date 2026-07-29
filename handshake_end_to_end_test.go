package dtls13

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

type dropWritesConn struct {
	net.Conn
	mu        sync.Mutex
	remaining int
}
type dropNthWriteConn struct {
	net.Conn
	mu          sync.Mutex
	count, drop int
}
type captureWritesConn struct {
	net.Conn
	mu     sync.Mutex
	writes [][]byte
}

type corruptSecondClientHelloCookieConn struct {
	net.Conn
	mu     sync.Mutex
	hellos int
}

type addInitialClientHelloLegacyCookieConn struct {
	net.Conn
	mu      sync.Mutex
	mutated bool
}

func (c *addInitialClientHelloLegacyCookieConn) Write(p []byte) (int, error) {
	c.mu.Lock()
	if c.mutated {
		c.mu.Unlock()
		return c.Conn.Write(p)
	}
	c.mu.Unlock()
	records, err := parsePlainRecords(p)
	if err != nil || len(records) != 1 || records[0].typ != recordTypeHandshake {
		return c.Conn.Write(p)
	}
	fragments, err := parseHandshakeFragments(records[0].payload)
	if err != nil || len(fragments) != 1 || fragments[0].typ != handshakeTypeClientHello || fragments[0].offset != 0 || int(fragments[0].length) != len(fragments[0].body) {
		return c.Conn.Write(p)
	}
	body := fragments[0].body
	if len(body) < 36 {
		return c.Conn.Write(p)
	}
	sessionEnd := 35 + int(body[34])
	if sessionEnd >= len(body) || body[sessionEnd] != 0 {
		return c.Conn.Write(p)
	}
	mutated := make([]byte, 0, len(body)+1)
	mutated = append(mutated, body[:sessionEnd]...)
	mutated = append(mutated, 1, 0x42)
	mutated = append(mutated, body[sessionEnd+1:]...)
	fragment := fragments[0]
	fragment.length = uint32(len(mutated))
	fragment.body = mutated
	payload, err := marshalHandshakeFragment(fragment)
	if err != nil {
		return 0, err
	}
	wire, err := marshalPlainRecord(record{typ: recordTypeHandshake, sequence: records[0].sequence, payload: payload})
	if err != nil {
		return 0, err
	}
	c.mu.Lock()
	c.mutated = true
	c.mu.Unlock()
	if _, err = c.Conn.Write(wire); err != nil {
		return 0, err
	}
	return len(p), nil
}

type replaceServerHelloWithSecondHRRConn struct {
	net.Conn
	mu      sync.Mutex
	hrrBody []byte
}

func (c *replaceServerHelloWithSecondHRRConn) Write(p []byte) (int, error) {
	records, err := parsePlainRecords(p)
	if err != nil || len(records) != 1 || records[0].typ != recordTypeHandshake {
		return c.Conn.Write(p)
	}
	fragments, err := parseHandshakeFragments(records[0].payload)
	if err != nil || len(fragments) != 1 || fragments[0].typ != handshakeTypeServerHello || fragments[0].offset != 0 || int(fragments[0].length) != len(fragments[0].body) {
		return c.Conn.Write(p)
	}
	body := fragments[0].body
	isHRR := len(body) >= 34 && string(body[2:34]) == string(helloRetryRequestRandom[:])
	c.mu.Lock()
	if isHRR {
		c.hrrBody = append([]byte(nil), body...)
		c.mu.Unlock()
		return c.Conn.Write(p)
	}
	hrrBody := append([]byte(nil), c.hrrBody...)
	c.mu.Unlock()
	if len(hrrBody) == 0 {
		return c.Conn.Write(p)
	}
	fragment := fragments[0]
	fragment.length = uint32(len(hrrBody))
	fragment.body = hrrBody
	payload, err := marshalHandshakeFragment(fragment)
	if err != nil {
		return 0, err
	}
	wire, err := marshalPlainRecord(record{typ: recordTypeHandshake, sequence: records[0].sequence, payload: payload})
	if err != nil {
		return 0, err
	}
	if _, err = c.Conn.Write(wire); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *corruptSecondClientHelloCookieConn) Write(p []byte) (int, error) {
	records, err := parsePlainRecords(p)
	if err != nil || len(records) != 1 || records[0].typ != recordTypeHandshake {
		return c.Conn.Write(p)
	}
	fragments, err := parseHandshakeFragments(records[0].payload)
	if err != nil || len(fragments) != 1 || fragments[0].typ != handshakeTypeClientHello || fragments[0].offset != 0 || int(fragments[0].length) != len(fragments[0].body) {
		return c.Conn.Write(p)
	}
	c.mu.Lock()
	c.hellos++
	second := c.hellos == 2
	c.mu.Unlock()
	if !second {
		return c.Conn.Write(p)
	}
	hello, err := parseClientHello(fragments[0].body)
	if err != nil || len(hello.cookie) == 0 {
		return c.Conn.Write(p)
	}
	hello.cookie[0] ^= 1
	body, err := hello.marshal()
	if err != nil {
		return 0, err
	}
	fragment := fragments[0]
	fragment.length = uint32(len(body))
	fragment.body = body
	payload, err := marshalHandshakeFragment(fragment)
	if err != nil {
		return 0, err
	}
	wire, err := marshalPlainRecord(record{typ: recordTypeHandshake, sequence: records[0].sequence, payload: payload})
	if err != nil {
		return 0, err
	}
	if _, err = c.Conn.Write(wire); err != nil {
		return 0, err
	}
	return len(p), nil
}

type weakNetworkConn struct {
	net.Conn
	mu      sync.Mutex
	enabled bool
	count   int
	held    []byte
}

func (c *weakNetworkConn) setEnabled(enabled bool) {
	c.mu.Lock()
	held := c.held
	c.held = nil
	c.enabled = enabled
	c.mu.Unlock()
	if len(held) > 0 {
		_, _ = c.Conn.Write(held)
	}
}

func (c *weakNetworkConn) Write(p []byte) (int, error) {
	c.mu.Lock()
	if !c.enabled {
		c.mu.Unlock()
		return c.Conn.Write(p)
	}
	c.count++
	count := c.count
	if count%13 == 0 || count%17 == 0 {
		c.mu.Unlock()
		return len(p), nil
	}
	current := append([]byte(nil), p...)
	if count%5 == 0 && c.held == nil {
		c.held = current
		c.mu.Unlock()
		return len(p), nil
	}
	held := c.held
	c.held = nil
	c.mu.Unlock()
	time.Sleep(time.Duration(count%3) * 200 * time.Microsecond)
	if _, err := c.Conn.Write(current); err != nil {
		return 0, err
	}
	if len(held) > 0 {
		if _, err := c.Conn.Write(held); err != nil {
			return 0, err
		}
	}
	if count%7 == 0 {
		_, _ = c.Conn.Write(current)
	}
	return len(p), nil
}

type memoryAddr string

func (a memoryAddr) Network() string { return "memory-datagram" }
func (a memoryAddr) String() string  { return string(a) }

type memoryDatagramConn struct {
	in, out                     chan []byte
	mu                          sync.Mutex
	readDeadline, writeDeadline time.Time
	closed                      chan struct{}
	closeOnce                   sync.Once
}

func memoryDatagramPair() (net.Conn, net.Conn) {
	aToB := make(chan []byte, 256)
	bToA := make(chan []byte, 256)
	aClosed, bClosed := make(chan struct{}), make(chan struct{})
	return &memoryDatagramConn{in: bToA, out: aToB, closed: aClosed},
		&memoryDatagramConn{in: aToB, out: bToA, closed: bClosed}
}
func (c *memoryDatagramConn) Read(p []byte) (int, error) {
	c.mu.Lock()
	deadline := c.readDeadline
	c.mu.Unlock()
	var timer <-chan time.Time
	if !deadline.IsZero() {
		duration := time.Until(deadline)
		if duration < 0 {
			duration = 0
		}
		timer = time.After(duration)
	}
	select {
	case b := <-c.in:
		return copy(p, b), nil
	case <-c.closed:
		return 0, net.ErrClosed
	case <-timer:
		return 0, &timeoutError{}
	}
}
func (c *memoryDatagramConn) Write(p []byte) (int, error) {
	b := append([]byte(nil), p...)
	c.mu.Lock()
	deadline := c.writeDeadline
	c.mu.Unlock()
	var timer <-chan time.Time
	if !deadline.IsZero() {
		duration := time.Until(deadline)
		if duration < 0 {
			duration = 0
		}
		timer = time.After(duration)
	}
	select {
	case c.out <- b:
		return len(p), nil
	case <-c.closed:
		return 0, net.ErrClosed
	case <-timer:
		return 0, &timeoutError{}
	}
}
func (c *memoryDatagramConn) Close() error {
	c.closeOnce.Do(func() { close(c.closed) })
	return nil
}
func (c *memoryDatagramConn) LocalAddr() net.Addr  { return memoryAddr("local") }
func (c *memoryDatagramConn) RemoteAddr() net.Addr { return memoryAddr("remote") }
func (c *memoryDatagramConn) SetDeadline(t time.Time) error {
	c.mu.Lock()
	c.readDeadline = t
	c.writeDeadline = t
	c.mu.Unlock()
	return nil
}
func (c *memoryDatagramConn) SetReadDeadline(t time.Time) error {
	c.mu.Lock()
	c.readDeadline = t
	c.mu.Unlock()
	return nil
}
func (c *memoryDatagramConn) SetWriteDeadline(t time.Time) error {
	c.mu.Lock()
	c.writeDeadline = t
	c.mu.Unlock()
	return nil
}

type timeoutError struct{}

func (*timeoutError) Error() string   { return "i/o timeout" }
func (*timeoutError) Timeout() bool   { return true }
func (*timeoutError) Temporary() bool { return true }

func (c *dropNthWriteConn) Write(p []byte) (int, error) {
	c.mu.Lock()
	c.count++
	drop := c.count == c.drop
	c.mu.Unlock()
	if drop {
		return len(p), nil
	}
	return c.Conn.Write(p)
}

func (c *captureWritesConn) Write(p []byte) (int, error) {
	c.mu.Lock()
	c.writes = append(c.writes, append([]byte(nil), p...))
	c.mu.Unlock()
	return c.Conn.Write(p)
}

func (c *captureWritesConn) reset() {
	c.mu.Lock()
	c.writes = nil
	c.mu.Unlock()
}

func (c *captureWritesConn) writeCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.writes)
}

func (c *captureWritesConn) lastWrite(t *testing.T) []byte {
	t.Helper()
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.writes) == 0 {
		t.Fatal("connection captured no writes")
	}
	return append([]byte(nil), c.writes[len(c.writes)-1]...)
}

func (c *dropWritesConn) Write(p []byte) (int, error) {
	c.mu.Lock()
	if c.remaining > 0 {
		c.remaining--
		c.mu.Unlock()
		return len(p), nil
	}
	c.mu.Unlock()
	return c.Conn.Write(p)
}

func waitForSendEpoch(t *testing.T, conn *Conn, epoch uint64) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		conn.writeMu.Lock()
		current := conn.sendCipher.epoch
		conn.writeMu.Unlock()
		if current == epoch {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("send epoch did not advance to %d", epoch)
}

func testServerCertificate(t testing.TB) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	pub, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(99), Subject: pkix.Name{CommonName: "server.test"}, DNSNames: []string{"server.test"}, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}, IsCA: true, BasicConstraintsValid: true}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, key)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(parsed)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key, Leaf: parsed}, pool
}

func testClientCertificate(t *testing.T) (tls.Certificate, *x509.CertPool) {
	t.Helper()
	pub, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now()
	tmpl := &x509.Certificate{SerialNumber: big.NewInt(100), Subject: pkix.Name{CommonName: "client"}, NotBefore: now.Add(-time.Hour), NotAfter: now.Add(time.Hour), KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}, IsCA: true, BasicConstraintsValid: true}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, pub, key)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	pool := x509.NewCertPool()
	pool.AddCert(parsed)
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: key, Leaf: parsed}, pool
}

func TestInvalidHRRCookieReturnsPlaintextIllegalParameter(t *testing.T) {
	left, right := memoryDatagramPair()
	client := Client(&corruptSecondClientHelloCookieConn{Conn: left}, &Config{InsecureSkipVerify: true, HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	defer left.Close()
	defer right.Close()

	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	clientErr := client.Handshake()
	if !errors.Is(clientErr, AlertError(alertIllegalParameter)) {
		t.Fatalf("client received %v", clientErr)
	}
	err := <-serverErr
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertIllegalParameter {
		t.Fatalf("server returned %v", err)
	}
}

func TestNonemptyLegacyCookieReturnsPlaintextIllegalParameter(t *testing.T) {
	left, right := memoryDatagramPair()
	client := Client(&addInitialClientHelloLegacyCookieConn{Conn: left}, &Config{InsecureSkipVerify: true, HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	defer left.Close()
	defer right.Close()

	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); !errors.Is(err, AlertError(alertIllegalParameter)) {
		t.Fatalf("client received %v", err)
	}
	err := <-serverErr
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertIllegalParameter {
		t.Fatalf("server returned %v", err)
	}
}

func TestClientRejectsSecondHRRWithPlaintextUnexpectedMessage(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(&replaceServerHelloWithSecondHRRConn{Conn: right}, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	defer left.Close()
	defer right.Close()

	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	err := client.Handshake()
	var local *localAlertError
	if !errors.As(err, &local) || local.description != alertUnexpectedMessage {
		t.Fatalf("client returned %v", err)
	}
	_ = left.Close()
	_ = right.Close()
	<-serverErr
}

func TestEndToEndCertificateHandshake(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", NextProtos: []string{"coap"}, HandshakeTimeout: 5 * time.Second})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, NextProtos: []string{"coap"}, HandshakeTimeout: 5 * time.Second})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if !client.ConnectionState().HandshakeComplete || client.ConnectionState().NegotiatedProtocol != "coap" {
		t.Fatalf("client state %#v", client.ConnectionState())
	}
	payload := []byte("protected application data")
	writeErr := make(chan error, 1)
	go func() { _, err := client.WriteDatagram(payload); writeErr <- err }()
	buf := make([]byte, 100)
	n, _, err := server.ReadDatagram(buf)
	if err != nil {
		t.Fatal(err)
	}
	if err = <-writeErr; err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf[:n], payload) {
		t.Fatalf("got %q", buf[:n])
	}
	_ = left.Close()
	_ = right.Close()
}

func TestCertificateVerificationFailureSendsUnknownCA(t *testing.T) {
	certificate, _ := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: x509.NewCertPool(), ServerName: "server.test", HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
	defer left.Close()
	defer right.Close()
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	clientErr := client.Handshake()
	var local *localAlertError
	if !errors.As(clientErr, &local) || local.description != alertUnknownCA {
		t.Fatalf("client returned %v", clientErr)
	}
	if err := <-serverErr; !errors.Is(err, AlertError(alertUnknownCA)) {
		t.Fatalf("server received %v", err)
	}
}

func TestHelloRetryRequestSelectsAdvertisedP256Share(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, CurvePreferences: []tls.CurveID{tls.CurveP256}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	defer left.Close()
	defer right.Close()
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
}

func TestEndToEndExportKeyingMaterial(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	defer left.Close()
	defer right.Close()
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	clientState, serverState := client.ConnectionState(), server.ConnectionState()
	clientKey, err := clientState.ExportKeyingMaterial("EXPORTER-test", []byte("context"), 64)
	if err != nil {
		t.Fatal(err)
	}
	serverKey, err := serverState.ExportKeyingMaterial("EXPORTER-test", []byte("context"), 64)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(clientKey, serverKey) || len(clientKey) != 64 {
		t.Fatal("client and server exporter output differ")
	}
	other, err := clientState.ExportKeyingMaterial("EXPORTER-test", []byte("other"), 64)
	if err != nil || bytes.Equal(clientKey, other) {
		t.Fatal("exporter context did not separate output")
	}
	if _, err = clientState.ExportKeyingMaterial(strings.Repeat("x", 250), nil, 1); err == nil {
		t.Fatal("accepted an exporter label longer than the HKDF label field")
	}
	if _, err = clientState.ExportKeyingMaterial("test", nil, 65536); err == nil {
		t.Fatal("accepted an exporter output longer than the HKDF length field")
	}
}

func TestEndToEndAES128CCM(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", CipherSuites: []uint16{TLS_AES_128_CCM_SHA256}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, CipherSuites: []uint16{TLS_AES_128_CCM_SHA256}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverDone; err != nil {
		t.Fatal(err)
	}
	if client.ConnectionState().CipherSuite != TLS_AES_128_CCM_SHA256 || server.ConnectionState().CipherSuite != TLS_AES_128_CCM_SHA256 {
		t.Fatal("AES-128-CCM was not negotiated")
	}
	if _, err := client.WriteDatagram([]byte("ccm application data")); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 64)
	n, _, err := server.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if string(buffer[:n]) != "ccm application data" {
		t.Fatalf("CCM payload = %q", buffer[:n])
	}
}

func TestPostHandshakeClientAuthentication(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate},
		PostHandshakeAuth: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	server := Server(right, &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.RequestClientCertificate(ctx); err != nil {
		t.Fatal(err)
	}
	state := server.ConnectionState()
	if len(state.PeerCertificates) == 0 || state.PeerCertificates[0].Subject.CommonName != "client" {
		t.Fatalf("post-handshake peer certificates %#v", state.PeerCertificates)
	}
	_ = client.Close()
	_ = server.Close()
}

func TestPostHandshakeClientAuthenticationRetransmitsDroppedACK(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	left, right := memoryDatagramPair()
	serverWire := &dropNthWriteConn{Conn: right}
	client := Client(left, &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate},
		PostHandshakeAuth: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	server := Server(serverWire, &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	serverWire.mu.Lock()
	serverWire.drop = serverWire.count + 2 // Request, then the first response ACK.
	serverWire.mu.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.RequestClientCertificate(ctx); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		client.writeMu.Lock()
		complete := client.clientAuthResponseFlight == nil
		client.writeMu.Unlock()
		if complete {
			_ = client.Close()
			_ = server.Close()
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("client authentication response did not recover from a lost ACK")
}

func TestPostHandshakeAuthenticationWhileClientKeyUpdateAwaitingACK(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	left, right := memoryDatagramPair()
	serverWire := &dropWritesConn{Conn: right}
	client := Client(left, &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate},
		PostHandshakeAuth: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 500 * time.Millisecond, MaxFlightInterval: 500 * time.Millisecond,
	})
	server := Server(serverWire, &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	ticketDeadline := time.Now().Add(time.Second)
	for {
		server.writeMu.Lock()
		ticketComplete := server.ticketFlight == nil
		server.writeMu.Unlock()
		if ticketComplete {
			break
		}
		if time.Now().After(ticketDeadline) {
			t.Fatal("initial NewSessionTicket was not acknowledged")
		}
		time.Sleep(time.Millisecond)
	}
	serverWire.mu.Lock()
	serverWire.remaining = 1
	serverWire.mu.Unlock()
	if err := client.SendKeyUpdate(false); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for {
		server.writeMu.Lock()
		receivedUpdate := server.receivingTraffic.current == 4
		server.writeMu.Unlock()
		client.writeMu.Lock()
		stillPending := !client.sendingTraffic.update.canUseNewKeys()
		client.writeMu.Unlock()
		if receivedUpdate && stillPending {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("client KeyUpdate did not remain pending after the dropped ACK")
		}
		time.Sleep(time.Millisecond)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.RequestClientCertificate(ctx); err != nil {
		t.Fatal(err)
	}
	state := server.ConnectionState()
	if len(state.PeerCertificates) == 0 || state.PeerCertificates[0].Subject.CommonName != "client" {
		t.Fatalf("post-handshake peer certificates %#v", state.PeerCertificates)
	}
	_ = client.Close()
	_ = server.Close()
}

func TestHandshakeRetransmitsDroppedClientHello(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(&dropWritesConn{Conn: left, remaining: 1}, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	_ = left.Close()
	_ = right.Close()
}
func TestHandshakeRetransmitsDroppedServerHello(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(&dropWritesConn{Conn: right, remaining: 1}, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	_ = left.Close()
	_ = right.Close()
}

func TestHandshakeSurvivesBurstLossBothDirections(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(&dropWritesConn{Conn: left, remaining: 3}, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 3 * time.Second, FlightInterval: 2 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond})
	server := Server(&dropWritesConn{Conn: right, remaining: 3}, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 3 * time.Second, FlightInterval: 2 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
}

func TestHandshakeStableUnderLossDelayReorderingAndDuplication(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	clientWire := &weakNetworkConn{Conn: left, enabled: true}
	serverWire := &weakNetworkConn{Conn: right, enabled: true}
	client := Client(clientWire, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 5 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(serverWire, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 5 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	clientWire.setEnabled(false)
	serverWire.setEnabled(false)
	payload := bytes.Repeat([]byte("weak-network"), 90)
	writeErr := make(chan error, 1)
	go func() { _, err := client.WriteDatagram(payload); writeErr <- err }()
	buffer := make([]byte, len(payload))
	n, info, err := server.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(payload) || info.Truncated {
		t.Fatalf("ReadDatagram=%d, %+v", n, info)
	}
	if err := <-writeErr; err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer, payload) {
		t.Fatal("application data mismatch after weak-network handshake")
	}
	_ = client.Close()
	_ = server.Close()
}

func TestHandshakeRetransmitsDroppedClientFinished(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(&dropNthWriteConn{Conn: left, drop: 3}, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 100 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	_ = left.Close()
	_ = right.Close()
}

func TestHandshakeRetransmitsWhenFinalACKIsDropped(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(&dropNthWriteConn{Conn: right, drop: 7}, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	_ = left.Close()
	_ = right.Close()
}

func TestClientAuthFlightRetransmitsWhenFinalACKIsDropped(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(&dropNthWriteConn{Conn: right, drop: 8}, &Config{Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	_ = left.Close()
	_ = right.Close()
}

func TestEndToEndRequiredClientCertificate(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	clientCertificate, clientRoots := testClientCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientRoots, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if len(server.ConnectionState().PeerCertificates) != 1 || len(server.ConnectionState().VerifiedChains) == 0 {
		t.Fatalf("server state %#v", server.ConnectionState())
	}
}

func TestClientAuthPolicyMatrixWithEmptyCertificate(t *testing.T) {
	serverCertificate, roots := testServerCertificate(t)
	tests := []struct {
		name    string
		policy  tls.ClientAuthType
		wantErr bool
	}{
		{"NoClientCert", tls.NoClientCert, false},
		{"RequestClientCert", tls.RequestClientCert, false},
		{"RequireAnyClientCert", tls.RequireAnyClientCert, true},
		{"VerifyClientCertIfGiven", tls.VerifyClientCertIfGiven, false},
		{"RequireAndVerifyClientCert", tls.RequireAndVerifyClientCert, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			left, right := memoryDatagramPair()
			defer left.Close()
			defer right.Close()
			client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
			server := Server(right, &Config{Certificates: []tls.Certificate{serverCertificate}, ClientAuth: test.policy, HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond})
			serverDone := make(chan error, 1)
			go func() { serverDone <- server.Handshake() }()
			clientErr := client.Handshake()
			serverErr := <-serverDone
			if !test.wantErr {
				if clientErr != nil || serverErr != nil {
					t.Fatalf("client error=%v, server error=%v", clientErr, serverErr)
				}
				return
			}
			if serverErr == nil {
				t.Fatalf("server accepted an empty required certificate; client=%v", clientErr)
			}
			if clientErr == nil {
				deadline := time.Now().Add(time.Second)
				for clientErr == nil && time.Now().Before(deadline) {
					client.inputMu.Lock()
					clientErr = client.readErr
					client.inputMu.Unlock()
					if clientErr == nil {
						time.Sleep(time.Millisecond)
					}
				}
			}
			var received AlertError
			if !errors.As(clientErr, &received) || uint8(received) != alertCertificateRequired {
				t.Fatalf("client error=%v, want certificate_required", clientErr)
			}
			var local *localAlertError
			if !errors.As(serverErr, &local) || local.description != alertCertificateRequired {
				t.Fatalf("server error=%v, want local certificate_required", serverErr)
			}
		})
	}
}

func TestEndToEndKeyUpdate(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if err := client.SendKeyUpdate(false); err != nil {
		t.Fatal(err)
	}
	waitForSendEpoch(t, client, 4)
	if _, err := client.WriteDatagram([]byte("new")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 32)
	n, _, err := server.ReadDatagram(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "new" {
		t.Fatalf("got %q", buf[:n])
	}
	if err = server.SendKeyUpdate(false); err != nil {
		t.Fatal(err)
	}
	waitForSendEpoch(t, server, 4)
	if _, err = server.WriteDatagram([]byte("response")); err != nil {
		t.Fatal(err)
	}
	n, _, err = client.ReadDatagram(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "response" {
		t.Fatalf("got %q", buf[:n])
	}
}

func TestKeyUpdateRetransmitsDroppedACK(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond})
	droppingServerConn := &dropWritesConn{Conn: right}
	server := Server(droppingServerConn, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	droppingServerConn.mu.Lock()
	droppingServerConn.remaining = 1
	droppingServerConn.mu.Unlock()
	if err := client.SendKeyUpdate(false); err != nil {
		t.Fatal(err)
	}
	waitForSendEpoch(t, client, 4)
	server.writeMu.Lock()
	receiveEpoch := server.receivingTraffic.current
	server.writeMu.Unlock()
	if receiveEpoch != 4 {
		t.Fatalf("duplicate KeyUpdate advanced receiver to epoch %d", receiveEpoch)
	}
	if _, err := client.WriteDatagram([]byte("after retransmit")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 32)
	n, _, err := server.ReadDatagram(buf)
	if err != nil {
		t.Fatal(err)
	}
	if string(buf[:n]) != "after retransmit" {
		t.Fatalf("got %q", buf[:n])
	}
}

func TestApplicationWriteAutomaticallyStartsKeyUpdateNearAEADLimit(t *testing.T) {
	client, server := establishedConnPair(t)
	client.startRecordReader()
	server.startRecordReader()
	client.writeMu.Lock()
	client.sendCipher.recordLimit = 4
	client.sendCipher.nextSequence = 2
	client.writeMu.Unlock()
	if _, err := client.WriteDatagram([]byte("trigger")); err != nil {
		t.Fatal(err)
	}
	waitForSendEpoch(t, client, 4)
	if client.ConnectionState().HandshakeComplete != server.ConnectionState().HandshakeComplete {
		t.Fatal("connection state changed during automatic KeyUpdate")
	}
	_ = client.Close()
	_ = server.Close()
}

func TestRepeatedBidirectionalKeyUpdatesRemainBounded(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	clientEpoch, serverEpoch := uint64(3), uint64(3)
	for i := 0; i < 64; i++ {
		if i%2 == 0 {
			if err := client.SendKeyUpdate(false); err != nil {
				t.Fatalf("client update %d: %v", i, err)
			}
			clientEpoch++
			waitForSendEpoch(t, client, clientEpoch)
		} else {
			if err := server.SendKeyUpdate(false); err != nil {
				t.Fatalf("server update %d: %v", i, err)
			}
			serverEpoch++
			waitForSendEpoch(t, server, serverEpoch)
		}
		for _, conn := range []*Conn{client, server} {
			conn.receiveEpochs.mu.RLock()
			retained := len(conn.receiveEpochs.ciphers)
			conn.receiveEpochs.mu.RUnlock()
			if retained > 2 {
				t.Fatalf("update %d retained %d receive epochs", i, retained)
			}
		}
	}
	_ = client.Close()
	_ = server.Close()
}

func TestApplicationBufferLimit(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, MaxBufferedApplicationData: 4, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if _, err := client.WriteDatagram([]byte("12345")); err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 8)
	if _, _, err := server.ReadDatagram(buf); err == nil || err.Error() != "dtls13: protocol error: buffered application data limit exceeded" {
		t.Fatalf("unexpected buffer limit error: %v", err)
	}
}

func TestCloseWakesBlockedRead(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	readErr := make(chan error, 1)
	go func() {
		buf := make([]byte, 1)
		_, _, err := server.ReadDatagram(buf)
		readErr <- err
	}()
	if err := server.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-readErr:
		if !errors.Is(err, net.ErrClosed) {
			t.Fatalf("blocked Read returned %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("blocked Read was not released by Close")
	}
}

func TestFatalAlertIsReported(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	body, err := (alertMessage{level: alertLevelFatal, description: 80}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	client.writeMu.Lock()
	wire, err := client.sendCipher.seal(recordTypeAlert, body)
	if err == nil {
		_, err = client.conn.Write(wire)
	}
	client.writeMu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1)
	_, _, err = server.ReadDatagram(buf)
	var alertErr AlertError
	if !errors.As(err, &alertErr) || alertErr != 80 {
		t.Fatalf("fatal alert returned %v", err)
	}
}

func TestEndToEndNewSessionTicket(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	cache := NewLRUClientSessionCache(4)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x6a}, 32))
	client := Client(left, &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	server := Server(right, &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey,
		SessionTicketLifetime: time.Hour, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	var cached *ClientSessionState
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if state, ok := cache.Get("server.test"); ok {
			cached = state
			break
		}
		time.Sleep(time.Millisecond)
	}
	if cached == nil {
		t.Fatal("client did not cache NewSessionTicket")
	}
	protector, err := newSessionTicketProtector(ticketKey, nil, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	ticketState, err := protector.open(cached.ticket)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(ticketState.psk, cached.psk) {
		t.Fatal("client and server derived different resumption PSKs")
	}
	if ticketState.suite != cached.suite || cached.serverName != "server.test" || cached.lifetime != 3600 {
		t.Fatalf("unexpected cached state: %#v", cached)
	}
	if client.ConnectionState().DidResume {
		t.Fatal("initial handshake was marked resumed")
	}
}

func TestSessionTicketsDisabled(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	cache := NewLRUClientSessionCache(1)
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	if _, ok := cache.Get("server.test"); ok {
		t.Fatal("client cached a ticket while server tickets were disabled")
	}
}

func TestEndToEndSessionResumption(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	cache := NewLRUClientSessionCache(4)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x5c}, 32))
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", NextProtos: []string{"coap"}, ClientSessionCache: cache,
		CipherSuites:     []uint16{TLS_CHACHA20_POLY1305_SHA256},
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, NextProtos: []string{"coap"}, SessionTicketKey: ticketKey,
		CipherSuites:          []uint16{TLS_CHACHA20_POLY1305_SHA256},
		SessionTicketLifetime: time.Hour, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	left, right := memoryDatagramPair()
	firstClient := Client(left, clientConfig)
	firstServer := Server(right, serverConfig)
	serverErr := make(chan error, 1)
	go func() { serverErr <- firstServer.Handshake() }()
	if err := firstClient.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if _, ok := cache.Get("server.test"); ok {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if _, ok := cache.Get("server.test"); !ok {
		t.Fatal("initial handshake did not produce a session ticket")
	}
	secondLeft, secondRight := memoryDatagramPair()
	client := Client(secondLeft, clientConfig)
	server := Server(secondRight, serverConfig)
	serverErr = make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume {
		t.Fatalf("resumption state client=%v server=%v", client.ConnectionState().DidResume, server.ConnectionState().DidResume)
	}
	if client.ConnectionState().CipherSuite != TLS_CHACHA20_POLY1305_SHA256 || server.ConnectionState().CipherSuite != TLS_CHACHA20_POLY1305_SHA256 {
		t.Fatal("resumed connection did not negotiate ChaCha20-Poly1305")
	}
	if len(client.ConnectionState().PeerCertificates) == 0 || len(client.ConnectionState().VerifiedChains) == 0 {
		t.Fatal("resumed client lost authenticated peer state")
	}
	if client.ConnectionState().NegotiatedProtocol != "coap" {
		t.Fatalf("resumed ALPN %q", client.ConnectionState().NegotiatedProtocol)
	}
	if _, err := client.WriteDatagram([]byte("resumed")); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 16)
	n, _, err := server.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if string(buffer[:n]) != "resumed" {
		t.Fatalf("got %q", buffer[:n])
	}
}

func TestSessionResumptionAllowsDifferentCipherSuiteWithSameHash(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	cache := NewLRUClientSessionCache(2)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x6d}, 32))
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache,
		CipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey,
		CipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, SessionTicketLifetime: time.Hour, MaxEarlyData: 1024,
		AllowEarlyDataWithoutCookie: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	_ = issueEarlyDataTicket(t, clientConfig, serverConfig)

	clientConfig.CipherSuites = []uint16{TLS_CHACHA20_POLY1305_SHA256}
	serverConfig.CipherSuites = []uint16{TLS_CHACHA20_POLY1305_SHA256}
	left, right := memoryDatagramPair()
	client := Client(left, clientConfig)
	server := Server(right, serverConfig)
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if n, err := client.WriteEarlyData([]byte("bound to the original cipher suite")); n != 0 || !errors.Is(err, ErrEarlyDataRejected) {
		t.Fatalf("cross-suite WriteEarlyData = %d, %v; want 0, ErrEarlyDataRejected", n, err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	clientState, serverState := client.ConnectionState(), server.ConnectionState()
	if !clientState.DidResume || !serverState.DidResume {
		t.Fatalf("cross-suite resumption state client=%v server=%v", clientState.DidResume, serverState.DidResume)
	}
	if clientState.CipherSuite != TLS_CHACHA20_POLY1305_SHA256 || serverState.CipherSuite != TLS_CHACHA20_POLY1305_SHA256 {
		t.Fatalf("cross-suite negotiation client=%x server=%x", clientState.CipherSuite, serverState.CipherSuite)
	}
}

func TestEndToEndEarlyData(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	cache := NewLRUClientSessionCache(4)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x7d}, 32))
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache,
		MTU: 256, IgnorePathMTU: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey,
		SessionTicketLifetime: time.Hour, MaxEarlyData: 4096,
		AllowEarlyDataWithoutCookie: true,
		HandshakeTimeout:            2 * time.Second,
		FlightInterval:              5 * time.Millisecond,
	}

	firstLeft, firstRight := memoryDatagramPair()
	firstClient := Client(firstLeft, clientConfig)
	firstServer := Server(firstRight, serverConfig)
	serverErr := make(chan error, 1)
	go func() { serverErr <- firstServer.Handshake() }()
	if err := firstClient.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if state, ok := cache.Get("server.test"); ok && state.maxEarlyData == 4096 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	state, ok := cache.Get("server.test")
	if !ok || state.maxEarlyData != 4096 {
		t.Fatal("initial handshake did not cache an early-data ticket")
	}

	secondLeft, secondRight := memoryDatagramPair()
	client := Client(secondLeft, clientConfig)
	server := Server(secondRight, serverConfig)
	serverErr = make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	payload := bytes.Repeat([]byte("e"), 700)
	n, err := client.WriteEarlyData(payload)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(payload) {
		t.Fatalf("WriteEarlyData wrote %d, want %d", n, len(payload))
	}
	if err = <-serverErr; err != nil {
		t.Fatal(err)
	}
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume {
		t.Fatal("0-RTT connection did not resume")
	}
	buffer := make([]byte, len(payload))
	n, _, err = server.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buffer[:n], payload) {
		t.Fatalf("early data %q, want %q", buffer[:n], payload)
	}
}

func issueEarlyDataTicket(t *testing.T, clientConfig, serverConfig *Config) *ClientSessionState {
	t.Helper()
	left, right := memoryDatagramPair()
	client := Client(left, clientConfig)
	server := Server(right, serverConfig)
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	key := clientSessionCacheKey(client.config, client.conn)
	for time.Now().Before(deadline) {
		if state, ok := clientConfig.ClientSessionCache.Get(key); ok {
			return state
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("initial handshake did not produce an early-data ticket")
	return nil
}

func TestEarlyDataRejectedAfterHelloRetryRequest(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	cache := NewLRUClientSessionCache(2)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x3e}, 32))
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey,
		SessionTicketLifetime: time.Hour, MaxEarlyData: 1024,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	_ = issueEarlyDataTicket(t, clientConfig, serverConfig)

	left, right := memoryDatagramPair()
	client := Client(left, clientConfig)
	server := Server(right, serverConfig)
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if n, err := client.WriteEarlyData([]byte("must be rejected after HRR")); n != 0 || !errors.Is(err, ErrEarlyDataRejected) {
		t.Fatalf("WriteEarlyData = %d, %v", n, err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume {
		t.Fatal("HRR fallback did not retain PSK resumption")
	}
	if len(server.earlyReadDatagrams) != 0 || server.earlyAccepted {
		t.Fatal("server accepted early data after HRR")
	}
}

func TestEarlyDataReplayIsRejected(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	initialCache := NewLRUClientSessionCache(1)
	replayCache := NewLRUEarlyDataReplayCache(8)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x4f}, 32))
	baseClient := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: initialCache,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey,
		SessionTicketLifetime: time.Hour, MaxEarlyData: 1024,
		AllowEarlyDataWithoutCookie: true, EarlyDataReplayCache: replayCache,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	ticket := issueEarlyDataTicket(t, baseClient, serverConfig)

	run := func(payload string) (int, *Conn, error) {
		cache := NewLRUClientSessionCache(1)
		cache.Put("server.test", ticket)
		clientConfig := baseClient.Clone()
		clientConfig.ClientSessionCache = cache
		left, right := memoryDatagramPair()
		client := Client(left, clientConfig)
		server := Server(right, serverConfig)
		serverErr := make(chan error, 1)
		go func() { serverErr <- server.Handshake() }()
		n, err := client.WriteEarlyData([]byte(payload))
		if serverHandshakeErr := <-serverErr; serverHandshakeErr != nil {
			t.Fatal(serverHandshakeErr)
		}
		if !server.ConnectionState().DidResume {
			t.Fatal("replayed connection did not retain 1-RTT resumption")
		}
		return n, server, err
	}

	if n, server, err := run("accepted once"); err != nil || n != len("accepted once") || !server.earlyAccepted {
		t.Fatalf("first early use = %d, %v, accepted=%v", n, err, server.earlyAccepted)
	}
	if n, server, err := run("replayed"); n != 0 || !errors.Is(err, ErrEarlyDataRejected) || server.earlyAccepted {
		t.Fatalf("replay = %d, %v, accepted=%v", n, err, server.earlyAccepted)
	}
}

func TestEarlyDataTicketLimit(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	cache := NewLRUClientSessionCache(1)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x2d}, 32))
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey,
		SessionTicketLifetime: time.Hour, MaxEarlyData: 8,
		AllowEarlyDataWithoutCookie: true,
		HandshakeTimeout:            2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	_ = issueEarlyDataTicket(t, clientConfig, serverConfig)

	left, right := memoryDatagramPair()
	client := Client(left, clientConfig)
	server := Server(right, serverConfig)
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if n, err := client.WriteEarlyData([]byte("ninebytes")); n != 0 || !errors.Is(err, ErrEarlyDataUnavailable) {
		t.Fatalf("oversized WriteEarlyData = %d, %v", n, err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if len(server.earlyReadDatagrams) != 0 || bufferedApplicationDatagrams(server) != 0 {
		t.Fatal("oversized early data reached the server application buffer")
	}
}

func TestWriteEarlyDataRejectsOversizedDatagram(t *testing.T) {
	cache := NewLRUClientSessionCache(1)
	cache.Put("server.test", &ClientSessionState{
		ticket: []byte("ticket"), psk: bytes.Repeat([]byte{1}, 32),
		suite: TLS_AES_128_GCM_SHA256, receivedAt: time.Now(), lifetime: 3600,
		serverName: "server.test", maxEarlyData: 4096,
	})
	left, right := memoryDatagramPair()
	defer left.Close()
	defer right.Close()
	client := Client(left, &Config{
		ServerName: "server.test", ClientSessionCache: cache, MTU: 256,
		HandshakeTimeout: time.Second,
	})
	n, err := client.WriteEarlyData(bytes.Repeat([]byte{1}, 256))
	if n != 0 || !errors.Is(err, ErrDatagramTooLarge) {
		t.Fatalf("WriteEarlyData oversized=%d, %v", n, err)
	}
	if client.earlySent {
		t.Fatal("oversized early data was sent")
	}
}

func TestSessionResumptionTicketAndBinderFailures(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	baseCache := NewLRUClientSessionCache(2)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x3d}, 32))
	baseClientConfig := &Config{RootCAs: roots, ServerName: "server.test", ClientSessionCache: baseCache, HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond}
	serverConfig := &Config{Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour, HandshakeTimeout: time.Second, FlightInterval: 5 * time.Millisecond}
	left, right := memoryDatagramPair()
	client := Client(left, baseClientConfig)
	server := Server(right, serverConfig)
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	var baseState *ClientSessionState
	for time.Now().Before(deadline) {
		if state, ok := baseCache.Get("server.test"); ok {
			baseState = state
			break
		}
		time.Sleep(time.Millisecond)
	}
	if baseState == nil {
		t.Fatal("initial handshake did not produce a ticket")
	}

	t.Run("tampered ticket falls back", func(t *testing.T) {
		state := cloneClientSessionState(baseState)
		state.ticket[len(state.ticket)-1] ^= 1
		cache := NewLRUClientSessionCache(1)
		cache.Put("server.test", state)
		clientConfig := baseClientConfig.Clone()
		clientConfig.ClientSessionCache = cache
		clientSide, serverSide := memoryDatagramPair()
		fallbackClient := Client(clientSide, clientConfig)
		fallbackServer := Server(serverSide, serverConfig)
		errors := make(chan error, 1)
		go func() { errors <- fallbackServer.Handshake() }()
		if err := fallbackClient.Handshake(); err != nil {
			t.Fatal(err)
		}
		if err := <-errors; err != nil {
			t.Fatal(err)
		}
		if fallbackClient.ConnectionState().DidResume || fallbackServer.ConnectionState().DidResume {
			t.Fatal("tampered ticket resumed a session")
		}
	})

	t.Run("valid ticket with wrong binder aborts", func(t *testing.T) {
		state := cloneClientSessionState(baseState)
		state.psk[0] ^= 1
		cache := NewLRUClientSessionCache(1)
		cache.Put("server.test", state)
		clientConfig := baseClientConfig.Clone()
		clientConfig.ClientSessionCache = cache
		clientConfig.HandshakeTimeout = 300 * time.Millisecond
		serverCopy := serverConfig.Clone()
		serverCopy.HandshakeTimeout = 300 * time.Millisecond
		clientSide, serverSide := memoryDatagramPair()
		badClient := Client(clientSide, clientConfig)
		badServer := Server(serverSide, serverCopy)
		errors := make(chan error, 1)
		go func() { errors <- badServer.Handshake() }()
		clientErr := badClient.Handshake()
		serverErr := <-errors
		if clientErr == nil || serverErr == nil || !strings.Contains(serverErr.Error(), "invalid PSK binder") {
			t.Fatalf("client error=%v server error=%v", clientErr, serverErr)
		}
	})
}

func TestNewSessionTicketLossRecovery(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x71}, 32))
	for _, test := range []struct {
		name       string
		dropServer int
		dropClient int
	}{
		{name: "ticket lost", dropServer: 8},
		{name: "ticket ACK lost", dropClient: 4},
	} {
		t.Run(test.name, func(t *testing.T) {
			cache := NewLRUClientSessionCache(1)
			left, right := memoryDatagramPair()
			clientConn := left
			serverConn := right
			if test.dropClient != 0 {
				clientConn = &dropNthWriteConn{Conn: left, drop: test.dropClient}
			}
			if test.dropServer != 0 {
				serverConn = &dropNthWriteConn{Conn: right, drop: test.dropServer}
			}
			client := Client(clientConn, &Config{RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond})
			server := Server(serverConn, &Config{Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond})
			serverErr := make(chan error, 1)
			go func() { serverErr <- server.Handshake() }()
			if err := client.Handshake(); err != nil {
				t.Fatal(err)
			}
			if err := <-serverErr; err != nil {
				t.Fatal(err)
			}
			deadline := time.Now().Add(time.Second)
			for time.Now().Before(deadline) {
				_, cached := cache.Get("server.test")
				server.writeMu.Lock()
				acknowledged := server.ticketFlight == nil
				server.writeMu.Unlock()
				if cached && acknowledged {
					return
				}
				time.Sleep(time.Millisecond)
			}
			t.Fatal("NewSessionTicket flight did not recover")
		})
	}
}

func TestFragmentedNewSessionTicket(t *testing.T) {
	certificate, _ := testServerCertificate(t)
	serverName := strings.Repeat("a", 400)
	cache := NewLRUClientSessionCache(1)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{InsecureSkipVerify: true, ServerName: serverName, ClientSessionCache: cache, MTU: 256, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, MTU: 256, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		_, cached := cache.Get(serverName)
		server.writeMu.Lock()
		acknowledged := server.ticketFlight == nil
		server.writeMu.Unlock()
		if cached && acknowledged {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("fragmented NewSessionTicket was not reassembled and acknowledged")
}

func TestEndToEndConnectionIDAndKeyUpdate(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	clientWire := &captureWritesConn{Conn: left}
	serverWire := &captureWritesConn{Conn: right}
	clientCID := []byte{0xc1, 0xc2}
	serverCID := []byte{0xd1, 0xd2, 0xd3}
	client := Client(clientWire, &Config{RootCAs: roots, ServerName: "server.test", CipherSuites: []uint16{TLS_CHACHA20_POLY1305_SHA256}, ConnectionID: clientCID, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(serverWire, &Config{Certificates: []tls.Certificate{certificate}, CipherSuites: []uint16{TLS_CHACHA20_POLY1305_SHA256}, ConnectionID: serverCID, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	clientState, serverState := client.ConnectionState(), server.ConnectionState()
	if clientState.CipherSuite != TLS_CHACHA20_POLY1305_SHA256 || serverState.CipherSuite != TLS_CHACHA20_POLY1305_SHA256 {
		t.Fatal("connection did not negotiate ChaCha20-Poly1305")
	}
	if !bytes.Equal(clientState.LocalConnectionID, clientCID) || !bytes.Equal(clientState.PeerConnectionID, serverCID) ||
		!bytes.Equal(serverState.LocalConnectionID, serverCID) || !bytes.Equal(serverState.PeerConnectionID, clientCID) {
		t.Fatalf("CID states client=%#v server=%#v", clientState, serverState)
	}
	clientWire.reset()
	if _, err := client.WriteDatagram([]byte("client CID")); err != nil {
		t.Fatal(err)
	}
	wire := clientWire.lastWrite(t)
	if wire[0]&unifiedHeaderCID == 0 || !bytes.Equal(wire[1:1+len(serverCID)], serverCID) {
		t.Fatalf("client used wrong peer CID: %x", wire)
	}
	buffer := make([]byte, 32)
	if _, _, err := server.ReadDatagram(buffer); err != nil {
		t.Fatal(err)
	}
	serverWire.reset()
	if _, err := server.WriteDatagram([]byte("server CID")); err != nil {
		t.Fatal(err)
	}
	wire = serverWire.lastWrite(t)
	if wire[0]&unifiedHeaderCID == 0 || !bytes.Equal(wire[1:1+len(clientCID)], clientCID) {
		t.Fatalf("server used wrong peer CID: %x", wire)
	}
	if _, _, err := client.ReadDatagram(buffer); err != nil {
		t.Fatal(err)
	}
	if err := client.SendKeyUpdate(false); err != nil {
		t.Fatal(err)
	}
	waitForSendEpoch(t, client, 4)
	client.writeMu.Lock()
	updatedCID := append([]byte(nil), client.sendCipher.connectionID...)
	client.writeMu.Unlock()
	if !bytes.Equal(updatedCID, serverCID) {
		t.Fatalf("KeyUpdate lost CID: %x", updatedCID)
	}
	if _, err := client.WriteDatagram([]byte("epoch 4 CID")); err != nil {
		t.Fatal(err)
	}
	n, _, err := server.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if string(buffer[:n]) != "epoch 4 CID" {
		t.Fatalf("got %q after CID KeyUpdate", buffer[:n])
	}
}

func TestConnectionIDRequiresBothPeers(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{1}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if len(client.ConnectionState().LocalConnectionID) != 0 || len(server.ConnectionState().PeerConnectionID) != 0 || len(client.sendCipher.connectionID) != 0 {
		t.Fatal("CID was enabled without server negotiation")
	}
}

func TestEndToEndEmptyConnectionID(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	clientWire := &captureWritesConn{Conn: left}
	client := Client(clientWire, &Config{RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, ConnectionID: []byte{}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	clientWire.reset()
	if _, err := client.WriteDatagram([]byte("empty CID")); err != nil {
		t.Fatal(err)
	}
	wire := clientWire.lastWrite(t)
	if wire[0]&unifiedHeaderCID == 0 {
		t.Fatalf("negotiated empty CID did not set the C bit: %x", wire)
	}
	buffer := make([]byte, 16)
	n, _, err := server.ReadDatagram(buffer)
	if err != nil || string(buffer[:n]) != "empty CID" {
		t.Fatalf("empty CID application record failed: %q %v", buffer[:n], err)
	}
}

func TestEndToEndImmediateConnectionIDUpdate(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	serverWire := &captureWritesConn{Conn: right}
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{1, 2}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(serverWire, &Config{Certificates: []tls.Certificate{certificate}, ConnectionID: []byte{3, 4}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	newClientCID := []byte{5, 6, 7}
	serverWire.reset()
	if err := client.SendNewConnectionIDs([][]byte{newClientCID}, true); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if bytes.Equal(server.ConnectionState().PeerConnectionID, newClientCID) && serverWire.writeCount() > 0 {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if !bytes.Equal(server.ConnectionState().PeerConnectionID, newClientCID) {
		t.Fatal("server did not activate the immediate peer CID")
	}
	wire := serverWire.lastWrite(t)
	if wire[0]&unifiedHeaderCID == 0 || !bytes.Equal(wire[1:1+len(newClientCID)], newClientCID) {
		t.Fatalf("server ACK did not use the immediate CID: %x", wire)
	}
	if _, err := server.WriteDatagram([]byte("new client CID")); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 32)
	n, _, err := client.ReadDatagram(buffer)
	if err != nil || string(buffer[:n]) != "new client CID" {
		t.Fatalf("application data after immediate CID update: %q %v", buffer[:n], err)
	}
}

func TestEndToEndRequestAndUseSpareConnectionID(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	clientWire := &captureWritesConn{Conn: left}
	spareServerCID := []byte{9, 10, 11}
	client := Client(clientWire, &Config{RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{1, 2}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, ConnectionID: []byte{3, 4}, GetConnectionID: func() ([]byte, error) { return append([]byte(nil), spareServerCID...), nil }, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if err := client.RequestConnectionIDs(1); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		client.writeMu.Lock()
		ready := !client.connectionIDRequestOpen && len(client.peerSpareConnectionIDs) > 0
		client.writeMu.Unlock()
		if ready {
			break
		}
		time.Sleep(time.Millisecond)
	}
	clientWire.reset()
	if err := client.UseNextConnectionID(); err != nil {
		t.Fatal(err)
	}
	if _, err := client.WriteDatagram([]byte("spare server CID")); err != nil {
		t.Fatal(err)
	}
	wire := clientWire.lastWrite(t)
	if wire[0]&unifiedHeaderCID == 0 || !bytes.Equal(wire[1:1+len(spareServerCID)], spareServerCID) {
		t.Fatalf("client did not use the spare server CID: %x", wire)
	}
	buffer := make([]byte, 32)
	n, _, err := server.ReadDatagram(buffer)
	if err != nil || string(buffer[:n]) != "spare server CID" {
		t.Fatalf("application data with spare CID: %q %v", buffer[:n], err)
	}
}

func TestConnectionIDUpdateRetransmitsDroppedACK(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	clientWire := &captureWritesConn{Conn: left}
	droppingServer := &dropWritesConn{Conn: right}
	client := Client(clientWire, &Config{RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{1, 2}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond})
	server := Server(droppingServer, &Config{Certificates: []tls.Certificate{certificate}, ConnectionID: []byte{3, 4}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	droppingServer.mu.Lock()
	droppingServer.remaining = 1
	droppingServer.mu.Unlock()
	clientWire.reset()
	if err := client.SendNewConnectionIDs([][]byte{{5, 6, 7}}, true); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		client.writeMu.Lock()
		complete := client.newConnectionIDFlight == nil
		client.writeMu.Unlock()
		if complete {
			clientWire.mu.Lock()
			writes := len(clientWire.writes)
			clientWire.mu.Unlock()
			if writes < 2 {
				t.Fatalf("CID update completed without retransmission: writes=%d", writes)
			}
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("CID update did not recover from a lost ACK")
}

func TestFragmentedNewConnectionID(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{1, 2}, MaxConnectionIDs: 128, MTU: 256, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, ConnectionID: []byte{3, 4}, MaxConnectionIDs: 128, MTU: 256, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	connectionIDs := make([][]byte, 80)
	for i := range connectionIDs {
		connectionIDs[i] = []byte{0x80, byte(i), 1, 2, 3, 4, 5, 6}
	}
	if err := client.SendNewConnectionIDs(connectionIDs, false); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		server.writeMu.Lock()
		received := len(server.peerSpareConnectionIDs)
		server.writeMu.Unlock()
		if received == len(connectionIDs) {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("fragmented NewConnectionId was not reassembled")
}

func TestConnectionIDRequestRemainsOpenUntilFulfilled(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	started := make(chan struct{})
	release := make(chan struct{})
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{1, 2}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, ConnectionID: []byte{3, 4}, GetConnectionID: func() ([]byte, error) {
		close(started)
		<-release
		return []byte{7, 8, 9}, nil
	}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	if err := client.RequestConnectionIDs(1); err != nil {
		t.Fatal(err)
	}
	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("peer did not process RequestConnectionId")
	}
	if err := client.RequestConnectionIDs(1); err == nil {
		t.Fatal("sent another RequestConnectionId before the first was fulfilled")
	}
	close(release)
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		client.writeMu.Lock()
		fulfilled := !client.connectionIDRequestOpen
		client.writeMu.Unlock()
		if fulfilled {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("NewConnectionId did not fulfill the outstanding request")
}

func TestConnectionIDRequestGeneratorFailureReturnsPartialResponse(t *testing.T) {
	for _, successful := range []int{0, 1} {
		t.Run(fmt.Sprintf("successful=%d", successful), func(t *testing.T) {
			certificate, roots := testServerCertificate(t)
			left, right := memoryDatagramPair()
			defer left.Close()
			defer right.Close()
			calls := 0
			client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{1, 2}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
			server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, ConnectionID: []byte{3, 4}, GetConnectionID: func() ([]byte, error) {
				calls++
				if calls > successful {
					return nil, errors.New("temporary CID generation failure")
				}
				return []byte{7, byte(calls), 9}, nil
			}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
			serverErr := make(chan error, 1)
			go func() { serverErr <- server.Handshake() }()
			if err := client.Handshake(); err != nil {
				t.Fatal(err)
			}
			if err := <-serverErr; err != nil {
				t.Fatal(err)
			}
			if err := client.RequestConnectionIDs(2); err != nil {
				t.Fatal(err)
			}
			deadline := time.Now().Add(time.Second)
			for time.Now().Before(deadline) {
				client.writeMu.Lock()
				fulfilled := !client.connectionIDRequestOpen
				spares := len(client.peerSpareConnectionIDs)
				client.writeMu.Unlock()
				if fulfilled {
					if spares != successful {
						t.Fatalf("received %d spare CIDs, want %d", spares, successful)
					}
					return
				}
				time.Sleep(time.Millisecond)
			}
			t.Fatal("partial NewConnectionId did not fulfill the request")
		})
	}
}

func TestUnexpectedConnectionIDMessageSendsFatalAlert(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	server := Server(right, &Config{Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond})
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	body := (requestConnectionIDMessage{count: 1}).marshal()
	fragment, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeRequestConnectionID, messageSequence: client.sendingTraffic.messageSequence, length: uint32(len(body)), body: body})
	if err != nil {
		t.Fatal(err)
	}
	client.writeMu.Lock()
	wire, err := client.sendCipher.seal(recordTypeHandshake, fragment)
	if err == nil {
		_, err = client.conn.Write(wire)
	}
	client.writeMu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		client.inputMu.Lock()
		readErr := client.readErr
		client.inputMu.Unlock()
		var alertErr AlertError
		if errors.As(readErr, &alertErr) {
			if uint8(alertErr) != alertUnexpectedMessage {
				t.Fatalf("received alert %d, want unexpected_message", alertErr)
			}
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("peer did not send a fatal unexpected_message alert")
}
