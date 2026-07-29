package dtls13

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

type reusableUDPClientConn struct {
	conn   *net.UDPConn
	remote *net.UDPAddr
}

func TestListenerConnectionIDBatchRegistrationIsAtomic(t *testing.T) {
	owner := &packetSession{}
	other := &packetSession{}
	listener := &packetListener{cidSessions: map[string]*packetSession{string([]byte{3, 4}): other}}
	if err := listener.registerSessionCIDs(owner, [][]byte{{1, 2}, {3, 4}}); err == nil {
		t.Fatal("registered a batch containing another session's CID")
	}
	if listener.cidSessions[string([]byte{1, 2})] != nil {
		t.Fatal("part of a rejected CID batch was registered")
	}
	if listener.cidSessions[string([]byte{3, 4})] != other {
		t.Fatal("rejected CID batch changed the existing owner")
	}
}

func (c *reusableUDPClientConn) Read(p []byte) (int, error) {
	for {
		n, from, err := c.conn.ReadFromUDP(p)
		if err != nil {
			return 0, err
		}
		if from.String() == c.remote.String() {
			return n, nil
		}
	}
}
func (c *reusableUDPClientConn) Write(p []byte) (int, error) {
	return c.conn.WriteToUDP(p, c.remote)
}
func (*reusableUDPClientConn) Close() error                    { return nil }
func (c *reusableUDPClientConn) LocalAddr() net.Addr           { return c.conn.LocalAddr() }
func (c *reusableUDPClientConn) RemoteAddr() net.Addr          { return c.remote }
func (c *reusableUDPClientConn) SetDeadline(t time.Time) error { return c.conn.SetDeadline(t) }
func (c *reusableUDPClientConn) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *reusableUDPClientConn) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

func TestListenAcceptRealUDP(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second,
		FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	serverErr := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverErr <- acceptErr
			return
		}
		defer conn.Close()
		_ = conn.SetDeadline(time.Now().Add(2 * time.Second))
		buffer := make([]byte, 32)
		n, _, readErr := conn.ReadDatagram(buffer)
		if readErr == nil && string(buffer[:n]) != "ping" {
			readErr = &ProtocolError{"unexpected listener payload"}
		}
		if readErr == nil {
			_, readErr = conn.WriteDatagram([]byte("pong"))
		}
		serverErr <- readErr
	}()
	dialer := &net.Dialer{Timeout: 2 * time.Second}
	client, err := DialWithDialer(dialer, "udp4", listener.Addr().String(), &Config{
		RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second,
		FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	_ = client.SetDeadline(time.Now().Add(2 * time.Second))
	if _, err = client.WriteDatagram([]byte("ping")); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 32)
	n, _, err := client.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if string(buffer[:n]) != "pong" {
		t.Fatalf("got %q", buffer[:n])
	}
	if err = <-serverErr; err != nil {
		t.Fatal(err)
	}
}

func TestListenerBoundsAndReclaimsSessions(t *testing.T) {
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{MaxPendingConnections: 1, MaxSessionQueueDatagrams: 1})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	internal := listener
	first, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer first.Close()
	second, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	address := listener.Addr()
	fragment, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeClientHello, messageSequence: 0, length: 1, body: []byte{0}})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 4; i++ {
		initial, marshalErr := marshalPlainRecord(record{typ: recordTypeHandshake, sequence: uint64(i), payload: fragment})
		if marshalErr != nil {
			t.Fatal(marshalErr)
		}
		_, _ = first.WriteTo(initial, address)
		_, _ = second.WriteTo(initial, address)
	}
	accepted, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	internal.mu.Lock()
	sessions := len(internal.sessions)
	internal.mu.Unlock()
	if sessions != 1 {
		t.Fatalf("listener retained %d sessions", sessions)
	}
	session := accepted.conn.(*packetSession)
	if queued := len(session.in); queued > 1 {
		t.Fatalf("session queued %d datagrams", queued)
	}
	if err = accepted.Close(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		internal.mu.Lock()
		reclaimed := true
		for _, candidate := range internal.sessions {
			if candidate == session {
				reclaimed = false
				break
			}
		}
		internal.mu.Unlock()
		if reclaimed {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("closed session was not reclaimed")
}

func TestListenerRejectsNonUDPNetwork(t *testing.T) {
	if _, err := Listen("tcp", "127.0.0.1:0", &Config{}); err == nil {
		t.Fatal("Listen accepted a non-UDP network")
	}
}

func TestInitialClientHelloDatagramFilter(t *testing.T) {
	fragment, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeClientHello, messageSequence: 0, length: 1, body: []byte{0}})
	if err != nil {
		t.Fatal(err)
	}
	datagram, err := marshalPlainRecord(record{typ: recordTypeHandshake, payload: fragment})
	if err != nil {
		t.Fatal(err)
	}
	if !isInitialClientHelloDatagram(datagram, 1024) {
		t.Fatal("rejected a structurally valid initial ClientHello fragment")
	}
	for i := 0; i < len(datagram); i++ {
		if isInitialClientHelloDatagram(datagram[:i], 1024) {
			t.Fatalf("accepted truncation at %d", i)
		}
	}
	for _, invalid := range [][]byte{{1}, append([]byte(nil), datagram...)} {
		if len(invalid) == len(datagram) {
			invalid[3], invalid[4] = 0, 1
		}
		if isInitialClientHelloDatagram(invalid, 1024) {
			t.Fatalf("accepted invalid initial datagram: %x", invalid)
		}
	}
	if isInitialClientHelloDatagram(datagram, 0) {
		t.Fatal("accepted a ClientHello over the configured message limit")
	}
}

func TestListenerConcurrentClients(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 3 * time.Second,
		FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	const clients = 8
	serverErrors := make(chan error, clients)
	go func() {
		for i := 0; i < clients; i++ {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				serverErrors <- acceptErr
				continue
			}
			go func(conn *Conn) {
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
				buffer := make([]byte, 32)
				n, _, exchangeErr := conn.ReadDatagram(buffer)
				if exchangeErr == nil {
					_, exchangeErr = conn.WriteDatagram(append([]byte("echo:"), buffer[:n]...))
				}
				serverErrors <- exchangeErr
			}(conn)
		}
	}()
	var wait sync.WaitGroup
	clientErrors := make(chan error, clients)
	for i := 0; i < clients; i++ {
		wait.Add(1)
		go func(index int) {
			defer wait.Done()
			dialer := &net.Dialer{Timeout: 3 * time.Second}
			conn, dialErr := DialWithDialer(dialer, "udp4", listener.Addr().String(), &Config{
				RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 3 * time.Second,
				FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
			})
			if dialErr != nil {
				clientErrors <- dialErr
				return
			}
			defer conn.Close()
			_ = conn.SetDeadline(time.Now().Add(3 * time.Second))
			message := fmt.Sprintf("client-%d", index)
			if _, dialErr = conn.WriteDatagram([]byte(message)); dialErr == nil {
				buffer := make([]byte, 32)
				var n int
				n, _, dialErr = conn.ReadDatagram(buffer)
				if dialErr == nil && string(buffer[:n]) != "echo:"+message {
					dialErr = &ProtocolError{"listener mixed client datagrams"}
				}
			}
			clientErrors <- dialErr
		}(i)
	}
	wait.Wait()
	for i := 0; i < clients; i++ {
		if err = <-clientErrors; err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < clients; i++ {
		if err = <-serverErrors; err != nil {
			t.Fatal(err)
		}
	}
}

func TestListenerConnReadDeadlineCanRecover(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, HandshakeTimeout: 2 * time.Second,
		FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	accepted := make(chan *Conn, 1)
	serverErr := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr == nil {
			acceptErr = conn.Handshake()
		}
		if acceptErr != nil {
			serverErr <- acceptErr
			return
		}
		accepted <- conn
	}()
	client, err := DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "udp4", listener.Addr().String(), &Config{
		RootCAs: roots, ServerName: "server.test", HandshakeTimeout: 2 * time.Second,
		FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	var server *Conn
	select {
	case server = <-accepted:
	case err = <-serverErr:
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("server handshake did not complete")
	}
	defer server.Close()
	if err = server.SetReadDeadline(time.Now().Add(10 * time.Millisecond)); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 16)
	if _, _, err = server.ReadDatagram(buffer); err == nil {
		t.Fatal("Read did not time out")
	} else if networkErr, ok := err.(net.Error); !ok || !networkErr.Timeout() {
		t.Fatalf("Read returned non-timeout error: %v", err)
	}
	if err = server.SetReadDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	if _, err = client.WriteDatagram([]byte("recovered")); err != nil {
		t.Fatal(err)
	}
	n, _, err := server.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if string(buffer[:n]) != "recovered" {
		t.Fatalf("got %q", buffer[:n])
	}
}

func TestListenerConnectionIDAuthenticatedMigration(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	serverCID := []byte{0x91, 0x92, 0x93, 0x94}
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true,
		GetConnectionID:  func() ([]byte, error) { return append([]byte(nil), serverCID...), nil },
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	accepted := make(chan *Conn, 1)
	serverErr := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverErr <- acceptErr
			return
		}
		server := conn
		if acceptErr = server.Handshake(); acceptErr != nil {
			serverErr <- acceptErr
			return
		}
		accepted <- server
	}()
	client, err := DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "udp4", listener.Addr().String(), &Config{
		RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{0x81, 0x82}, SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	var server *Conn
	select {
	case server = <-accepted:
	case err = <-serverErr:
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("server handshake timed out")
	}
	defer server.Close()
	runtimeServerCID := []byte{0xa1, 0xa2, 0xa3, 0xa4, 0xa5}
	if err = server.SendNewConnectionIDs([][]byte{runtimeServerCID}, false); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		client.writeMu.Lock()
		ready := len(client.peerSpareConnectionIDs) > 0
		client.writeMu.Unlock()
		if ready {
			break
		}
		time.Sleep(time.Millisecond)
	}
	if err = client.UseNextConnectionID(); err != nil {
		t.Fatal(err)
	}
	originalRemote := server.RemoteAddr().String()
	migrated, err := net.ListenPacket("udp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer migrated.Close()
	client.writeMu.Lock()
	wire, err := client.sendCipher.seal(recordTypeApplicationData, []byte("migrated"))
	client.writeMu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	tampered := append([]byte(nil), wire...)
	tampered[len(tampered)-1] ^= 1
	if _, err = migrated.WriteTo(tampered, listener.Addr()); err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	if server.RemoteAddr().String() != originalRemote {
		t.Fatal("unauthenticated CID record changed the peer address")
	}
	if _, err = migrated.WriteTo(wire, listener.Addr()); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 32)
	n, _, err := server.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if string(buffer[:n]) != "migrated" {
		t.Fatalf("got %q", buffer[:n])
	}
	if server.RemoteAddr().String() != originalRemote {
		t.Fatalf("server changed peer address without path validation: got %s want %s", server.RemoteAddr(), originalRemote)
	}
	if _, err = server.WriteDatagram([]byte("new path")); err != nil {
		t.Fatal(err)
	}
	_ = client.SetReadDeadline(time.Now().Add(time.Second))
	n, _, err = client.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if string(buffer[:n]) != "new path" {
		t.Fatalf("response content=%q", buffer[:n])
	}
}

func TestListenerReturnRoutabilityNATRebinding(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	serverCID := []byte{0x91, 0x92, 0x93, 0x94}
	clientCID := []byte{0x81, 0x82, 0x83}
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true,
		GetConnectionID:  func() ([]byte, error) { return append([]byte(nil), serverCID...), nil },
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	accepted := make(chan *Conn, 1)
	serverErr := make(chan error, 1)
	go func() {
		server, acceptErr := listener.Accept()
		if acceptErr == nil {
			acceptErr = server.Handshake()
		}
		if acceptErr != nil {
			serverErr <- acceptErr
			return
		}
		accepted <- server
	}()
	client, err := DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "udp4", listener.Addr().String(), &Config{
		RootCAs: roots, ServerName: "server.test", ConnectionID: clientCID, SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	var server *Conn
	select {
	case server = <-accepted:
	case err = <-serverErr:
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("server handshake timed out")
	}
	defer server.Close()
	if !client.ConnectionState().ReturnRoutabilityCheck || !server.ConnectionState().ReturnRoutabilityCheck {
		t.Fatal("RRC was not negotiated")
	}

	client.writeMu.Lock()
	wire, err := client.sendCipher.seal(recordTypeApplicationData, []byte("rebinding"))
	client.writeMu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	oldTransport, ok := client.conn.(*net.UDPConn)
	if !ok {
		t.Fatalf("client transport is %T", client.conn)
	}
	oldAddress := oldTransport.LocalAddr()
	if err = oldTransport.SetReadDeadline(time.Now()); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		client.readerMu.Lock()
		running := client.readerRunning
		client.readerMu.Unlock()
		if !running {
			break
		}
		time.Sleep(time.Millisecond)
	}
	client.readerMu.Lock()
	readerRunning := client.readerRunning
	client.readerMu.Unlock()
	if readerRunning {
		t.Fatal("old path reader did not stop")
	}
	if err = oldTransport.SetReadDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}
	migrated, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer migrated.Close()
	server.retransmitNanos.Store(int64(5 * time.Millisecond))
	server.lastRTTSampleUnixNano.Store(time.Now().UnixNano())
	if _, err = migrated.WriteTo(wire, listener.Addr()); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 2048)
	if n, _, readErr := server.ReadDatagram(buffer); readErr != nil || string(buffer[:n]) != "rebinding" {
		t.Fatalf("server read %q, %v", buffer[:n], readErr)
	}

	if err = migrated.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	n, _, err := migrated.ReadFromUDP(buffer)
	if err != nil {
		t.Fatal(err)
	}
	content, contentType, _, _, err := client.receiveEpochs.openInPlace(buffer[:n])
	if err != nil || contentType != recordTypeReturnRoutability {
		t.Fatalf("basic challenge type=%d err=%v", contentType, err)
	}
	challenge, known, err := parseReturnRoutabilityMessage(content)
	if err != nil || !known || challenge.typ != returnRoutabilityPathChallenge {
		t.Fatalf("challenge=%#v known=%v err=%v", challenge, known, err)
	}
	response, err := (returnRoutabilityMessage{typ: returnRoutabilityPathResponse, cookie: challenge.cookie}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	client.writeMu.Lock()
	responseWire, err := client.sendCipher.seal(recordTypeReturnRoutability, response[:])
	client.writeMu.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = migrated.WriteTo(responseWire, listener.Addr()); err != nil {
		t.Fatal(err)
	}
	deadline = time.Now().Add(time.Second)
	for time.Now().Before(deadline) && !sameNetworkAddress(server.RemoteAddr(), migrated.LocalAddr()) {
		time.Sleep(time.Millisecond)
	}
	if !sameNetworkAddress(server.RemoteAddr(), migrated.LocalAddr()) {
		t.Fatalf("server did not rebind: got %v want %v", server.RemoteAddr(), migrated.LocalAddr())
	}

	if _, err = server.WriteDatagram([]byte("new path")); err != nil {
		t.Fatal(err)
	}
	n, _, err = migrated.ReadFromUDP(buffer)
	if err != nil {
		t.Fatal(err)
	}
	content, contentType, _, _, err = client.receiveEpochs.openInPlace(buffer[:n])
	if err != nil || contentType != recordTypeApplicationData || string(content) != "new path" {
		t.Fatalf("rebound response type=%d content=%q err=%v", contentType, content, err)
	}

	internal := listener
	internal.mu.Lock()
	bound := internal.sessions[sessionKey(migrated.LocalAddr())]
	oldBound := internal.sessions[sessionKey(oldAddress)]
	internal.mu.Unlock()
	if bound != server.conn.(*packetSession) || oldBound != nil {
		t.Fatal("Listener tuple map was not updated with the validated path")
	}
}

func TestListenerRoutesNegotiatedEmptyConnectionID(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, ConnectionID: []byte{}, SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	accepted := make(chan *Conn, 1)
	serverErr := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			serverErr <- acceptErr
			return
		}
		server := conn
		if acceptErr = server.Handshake(); acceptErr != nil {
			serverErr <- acceptErr
			return
		}
		accepted <- server
	}()
	client, err := DialWithDialer(&net.Dialer{Timeout: 2 * time.Second}, "udp4", listener.Addr().String(), &Config{
		RootCAs: roots, ServerName: "server.test", ConnectionID: []byte{}, SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	var server *Conn
	select {
	case server = <-accepted:
	case err = <-serverErr:
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("server handshake timed out")
	}
	defer server.Close()
	if !client.connectionIDNegotiated || !server.connectionIDNegotiated {
		t.Fatal("empty connection ID was not negotiated")
	}
	if _, err = client.WriteDatagram([]byte("empty cid")); err != nil {
		t.Fatal(err)
	}
	buffer := make([]byte, 32)
	n, _, err := server.ReadDatagram(buffer)
	if err != nil || string(buffer[:n]) != "empty cid" {
		t.Fatalf("server read %q, %v", buffer[:n], err)
	}
}

func TestListenerReplacesAssociationOnlyAfterNewFinished(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	serverAddress := listener.Addr().(*net.UDPAddr)
	udpClient, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	defer udpClient.Close()
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
	}

	firstServerRaw := make(chan *Conn, 1)
	firstServerErr := make(chan error, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr == nil {
			acceptErr = conn.Handshake()
		}
		if acceptErr != nil {
			firstServerErr <- acceptErr
			return
		}
		firstServerRaw <- conn
	}()
	firstClient := Client(&reusableUDPClientConn{conn: udpClient, remote: serverAddress}, clientConfig)
	if err = firstClient.Handshake(); err != nil {
		t.Fatal(err)
	}
	var firstServer *Conn
	select {
	case conn := <-firstServerRaw:
		firstServer = conn
	case err = <-firstServerErr:
		t.Fatal(err)
	case <-time.After(2 * time.Second):
		t.Fatal("first server handshake timed out")
	}
	defer firstServer.Close()
	oldSession := firstServer.conn.(*packetSession)

	// Stop the first client's background reader without closing the UDP
	// socket, simulating a reboot/abandonment followed by tuple reuse.
	if err = udpClient.SetReadDeadline(time.Now()); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	readerStopped := false
	for time.Now().Before(deadline) {
		firstClient.readerMu.Lock()
		running := firstClient.readerRunning
		firstClient.readerMu.Unlock()
		if !running {
			readerStopped = true
			break
		}
		time.Sleep(time.Millisecond)
	}
	if !readerStopped {
		t.Fatal("first client reader did not stop before reusing its UDP socket")
	}
	if err = udpClient.SetDeadline(time.Time{}); err != nil {
		t.Fatal(err)
	}

	secondClient := Client(&reusableUDPClientConn{conn: udpClient, remote: serverAddress}, clientConfig)
	secondClientErr := make(chan error, 1)
	go func() { secondClientErr <- secondClient.Handshake() }()
	secondRaw, err := listener.Accept()
	if err != nil {
		t.Fatal(err)
	}
	secondServer := secondRaw
	defer secondServer.Close()
	select {
	case <-oldSession.done:
		t.Fatal("old association was destroyed before the new Finished was verified")
	default:
	}
	oldPayload := []byte("old association still valid")
	// WriteDatagram would restart the reader stopped above. Seal directly so
	// the old association stays usable without competing for the reused socket.
	firstClient.writeMu.Lock()
	oldWire, sealErr := firstClient.sendCipher.seal(recordTypeApplicationData, oldPayload)
	firstClient.writeMu.Unlock()
	if sealErr != nil {
		t.Fatal(sealErr)
	}
	if _, err = udpClient.WriteToUDP(oldWire, serverAddress); err != nil {
		t.Fatal(err)
	}
	oldBuffer := make([]byte, 32)
	oldN, _, err := firstServer.ReadDatagram(oldBuffer)
	if err != nil {
		t.Fatal(err)
	}
	if string(oldBuffer[:oldN]) != string(oldPayload) {
		t.Fatalf("old association payload %q", oldBuffer[:oldN])
	}
	secondServerErr := make(chan error, 1)
	go func() { secondServerErr <- secondServer.Handshake() }()
	if err = <-secondClientErr; err != nil {
		select {
		case serverErr := <-secondServerErr:
			t.Fatalf("second client handshake: %v; second server handshake: %v", err, serverErr)
		case <-time.After(3 * time.Second):
			t.Fatalf("second client handshake: %v; second server handshake did not return", err)
		}
	}
	if err = <-secondServerErr; err != nil {
		t.Fatal(err)
	}
	select {
	case <-oldSession.done:
	case <-time.After(time.Second):
		t.Fatal("old association remained active after the new Finished")
	}

	internal := listener
	key := sessionKey(udpClient.LocalAddr())
	newSession := secondServer.conn.(*packetSession)
	internal.mu.Lock()
	active, pending := internal.sessions[key], internal.pending[key]
	internal.mu.Unlock()
	if active != newSession || pending != nil {
		t.Fatal("listener did not atomically promote the validated association")
	}

	payload := []byte("new association")
	writeErr := make(chan error, 1)
	go func() { _, writeErrValue := secondClient.WriteDatagram(payload); writeErr <- writeErrValue }()
	buffer := make([]byte, 32)
	n, _, err := secondServer.ReadDatagram(buffer)
	if err != nil {
		t.Fatal(err)
	}
	if err = <-writeErr; err != nil {
		t.Fatal(err)
	}
	if string(buffer[:n]) != string(payload) {
		t.Fatalf("new association payload %q", buffer[:n])
	}
}

func TestListenerRejectsAmbiguousConnectionIDs(t *testing.T) {
	l := &packetListener{cidSessions: make(map[string]*packetSession), sessions: make(map[string]*packetSession)}
	first := &packetSession{}
	second := &packetSession{}
	if err := l.registerSessionCID(first, []byte{1, 2}); err != nil {
		t.Fatal(err)
	}
	if err := l.registerSessionCID(second, []byte{1, 2, 3}); err == nil {
		t.Fatal("Listener accepted prefix-ambiguous connection IDs")
	}
}

func TestListenerMixedCIDDatagramDoesNotCrossAssociations(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	secretA := bytes.Repeat([]byte{0x41}, suite.hash.Size())
	secretB := bytes.Repeat([]byte{0x42}, suite.hash.Size())
	cidA := []byte{0xa1, 0xa2}
	cidB := []byte{0xb1, 0xb2, 0xb3}

	senderA, err := newRecordCipher(suite, secretA, 3, 64)
	if err != nil {
		t.Fatal(err)
	}
	senderB, err := newRecordCipher(suite, secretB, 3, 64)
	if err != nil {
		t.Fatal(err)
	}
	if err = senderA.setConnectionID(cidA); err != nil {
		t.Fatal(err)
	}
	if err = senderB.setConnectionID(cidB); err != nil {
		t.Fatal(err)
	}
	recordA, err := senderA.seal(recordTypeApplicationData, []byte("association-a"))
	if err != nil {
		t.Fatal(err)
	}
	recordB, err := senderB.seal(recordTypeApplicationData, []byte("association-b"))
	if err != nil {
		t.Fatal(err)
	}
	mixed := append(append([]byte(nil), recordA...), recordB...)

	sessionA, sessionB := &packetSession{}, &packetSession{}
	listener := &packetListener{cidSessions: map[string]*packetSession{string(cidA): sessionA, string(cidB): sessionB}}
	if got := listener.sessionForCIDLocked(mixed); got != sessionA {
		t.Fatalf("mixed datagram routed to %p, want association A %p", got, sessionA)
	}

	client, peer := establishedConnPair(t)
	defer client.conn.Close()
	defer peer.conn.Close()
	receiving, err := newReceivingTraffic(suite, secretA, 3, 64)
	if err != nil {
		t.Fatal(err)
	}
	if err = receiving.setConnectionID(cidA); err != nil {
		t.Fatal(err)
	}
	client.receivingTraffic = receiving
	client.receiveEpochs = receiving.epochs
	if err = client.dispatchDatagram(mixed); err != nil {
		t.Fatal(err)
	}
	got := string(bufferedApplicationPayload(client))
	if got != "association-a" {
		t.Fatalf("mixed datagram delivered %q, want only association A", got)
	}
}
