package dtls13

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"net"
	"sync"
	"testing"
	"time"
)

func replaceClientHelloTicketRequest(t *testing.T, wire, value []byte) []byte {
	t.Helper()
	p := wireParser{b: wire}
	p.u16()
	p.take(32)
	p.bytes8()
	p.bytes8()
	p.bytes16()
	p.bytes8()
	extensionsOffset := p.off
	extensions := p.bytes16()
	if err := p.done(); err != nil {
		t.Fatal(err)
	}
	for offset := 0; offset < len(extensions); {
		if len(extensions)-offset < 4 {
			t.Fatal("truncated test extension")
		}
		typ := binary.BigEndian.Uint16(extensions[offset:])
		length := int(binary.BigEndian.Uint16(extensions[offset+2:]))
		if length > len(extensions)-offset-4 {
			t.Fatal("invalid test extension length")
		}
		if typ == extTicketRequest {
			absolute := extensionsOffset + 2 + offset
			out := make([]byte, 0, len(wire)-length+len(value))
			out = append(out, wire[:absolute+4]...)
			out = append(out, value...)
			out = append(out, wire[absolute+4+length:]...)
			binary.BigEndian.PutUint16(out[absolute+2:absolute+4], uint16(len(value)))
			binary.BigEndian.PutUint16(out[extensionsOffset:extensionsOffset+2], uint16(len(extensions)-length+len(value)))
			return out
		}
		offset += 4 + length
	}
	t.Fatal("test ClientHello has no ticket_request extension")
	return nil
}

func appendServerHelloExtension(t *testing.T, body []byte, typ uint16, value []byte) []byte {
	t.Helper()
	if len(body) < 40 || body[34] != 0 {
		t.Fatal("test ServerHello has unexpected legacy fields")
	}
	extensionsOffset := 38
	length := int(binary.BigEndian.Uint16(body[extensionsOffset:]))
	if extensionsOffset+2+length != len(body) {
		t.Fatal("test ServerHello has invalid extensions")
	}
	out := append([]byte(nil), body...)
	binary.BigEndian.PutUint16(out[extensionsOffset:], uint16(length+4+len(value)))
	out = binary.BigEndian.AppendUint16(out, typ)
	out = binary.BigEndian.AppendUint16(out, uint16(len(value)))
	out = append(out, value...)
	return out
}

func TestTicketRequestWireFormatAndMessagePlacement(t *testing.T) {
	hello := &clientHello{
		cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, supportedGroups: []tls.CurveID{tls.X25519},
		keyShares:     []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
		ticketRequest: SessionTicketRequest{Enabled: true, NewSessionCount: 4, ResumptionCount: 1},
	}
	wire, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseClientHello(wire)
	if err != nil {
		t.Fatal(err)
	}
	if !parsed.ticketRequest.Enabled || parsed.ticketRequest != hello.ticketRequest {
		t.Fatalf("ticket_request = %#v, present=%v", parsed.ticketRequest, parsed.ticketRequest.Enabled)
	}
	for _, value := range [][]byte{nil, {1}, {1, 2, 3}} {
		_, err = parseClientHello(replaceClientHelloTicketRequest(t, wire, value))
		if description, ok := protocolAlert(err); !ok || description != alertDecodeError {
			t.Fatalf("length %d alert=%d ok=%v err=%v", len(value), description, ok, err)
		}
	}

	serverHello, err := (&serverHello{cipherSuite: TLS_AES_128_GCM_SHA256, keyShare: keyShareEntry{group: tls.X25519, data: bytes.Repeat([]byte{2}, 32)}}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = parseServerHello(appendServerHelloExtension(t, serverHello, extTicketRequest, []byte{1})); err == nil {
		t.Fatal("accepted ticket_request in ServerHello")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("ServerHello alert=%d ok=%v err=%v", description, ok, err)
	}
	hrr, err := (&helloRetryRequest{cipherSuite: TLS_AES_128_GCM_SHA256, selectedGroup: tls.CurveP256}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	if _, err = parseHelloRetryRequest(appendServerHelloExtension(t, hrr, extTicketRequest, []byte{1})); err == nil {
		t.Fatal("accepted ticket_request in HelloRetryRequest")
	} else if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("HelloRetryRequest alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestTicketRequestEncryptedExtensionsAndHRRInvariants(t *testing.T) {
	hello := &clientHello{ticketRequest: SessionTicketRequest{Enabled: true, NewSessionCount: 4, ResumptionCount: 1}}
	message := &encryptedExtensions{extensions: map[uint16][]byte{extTicketRequest: {3}}}
	if _, _, _, err := validateEncryptedExtensions(hello, message); err != nil {
		t.Fatal(err)
	}
	if !message.hasTicketRequest || message.expectedTicketCount != 3 {
		t.Fatalf("expected_count=%d present=%v", message.expectedTicketCount, message.hasTicketRequest)
	}
	for _, test := range []struct {
		name      string
		hello     *clientHello
		value     []byte
		wantAlert uint8
	}{
		{name: "Unsolicited", hello: &clientHello{}, value: []byte{1}, wantAlert: alertUnsupportedExtension},
		{name: "Empty", hello: hello, wantAlert: alertDecodeError},
		{name: "Long", hello: hello, value: []byte{1, 2}, wantAlert: alertDecodeError},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, _, _, err := validateEncryptedExtensions(test.hello, &encryptedExtensions{extensions: map[uint16][]byte{extTicketRequest: test.value}})
			if description, ok := protocolAlert(err); !ok || description != test.wantAlert {
				t.Fatalf("alert=%d ok=%v err=%v", description, ok, err)
			}
		})
	}

	second := *hello
	if !equalClientHelloAfterHRR(hello, &second, 0) {
		t.Fatal("unchanged ticket_request failed the HRR invariant")
	}
	second.ticketRequest.ResumptionCount++
	if equalClientHelloAfterHRR(hello, &second, 0) {
		t.Fatal("changed ticket_request value passed the HRR invariant")
	}
	second = *hello
	second.ticketRequest.Enabled = false
	if equalClientHelloAfterHRR(hello, &second, 0) {
		t.Fatal("removed ticket_request passed the HRR invariant")
	}
}

func ticketCacheSnapshot(cache *lruSessionCache, key string) []*ClientSessionState {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	element := cache.entries[key]
	if element == nil {
		return nil
	}
	states := element.Value.(*sessionCacheEntry).states
	out := make([]*ClientSessionState, len(states))
	for i, state := range states {
		out[i] = cloneClientSessionState(state)
	}
	return out
}

func waitForTicketCount(t *testing.T, cache *lruSessionCache, want int) []*ClientSessionState {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		states := ticketCacheSnapshot(cache, "server.test")
		if len(states) == want {
			return states
		}
		time.Sleep(time.Millisecond)
	}
	states := ticketCacheSnapshot(cache, "server.test")
	t.Fatalf("cached %d tickets, want %d", len(states), want)
	return nil
}

func runTicketRequestHandshake(t *testing.T, clientConfig, serverConfig *Config, clientConn, serverConn net.Conn) (*Conn, *Conn) {
	t.Helper()
	client := Client(clientConn, clientConfig)
	server := Server(serverConn, serverConfig)
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		t.Fatal(err)
	}
	if err := <-serverErr; err != nil {
		t.Fatal(err)
	}
	return client, server
}

func TestTicketRequestEndToEndCountsAndResumption(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	cache := NewLRUClientSessionCache(2).(*lruSessionCache)
	request := SessionTicketRequest{Enabled: true, NewSessionCount: 4, ResumptionCount: 1}
	var ticketKey [32]byte
	ticketKey[0] = 1
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache, SessionTicketRequest: request,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, MaxSessionTickets: 2, SessionTicketKey: ticketKey,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	}
	left, right := memoryDatagramPair()
	client, server := runTicketRequestHandshake(t, clientConfig, serverConfig, left, right)
	states := waitForTicketCount(t, cache, 2)
	if client.ConnectionState().DidResume || server.ConnectionState().DidResume {
		t.Fatal("initial ticket_request handshake resumed")
	}
	tickets, nonces, psks := make(map[string]bool), make(map[string]bool), make(map[string]bool)
	group := states[0].ticketGroup
	for _, state := range states {
		tickets[string(state.ticket)] = true
		nonces[string(state.nonce)] = true
		psks[string(state.psk)] = true
		if state.ticketGroup != group || group == ([32]byte{}) {
			t.Fatal("tickets from one connection did not share a nonzero group")
		}
	}
	if len(tickets) != 2 || len(nonces) != 2 || len(psks) != 2 {
		t.Fatal("server did not issue distinct ticket, nonce, and PSK values")
	}
	_ = client.Close()
	_ = server.Close()

	left, right = memoryDatagramPair()
	client, server = runTicketRequestHandshake(t, clientConfig, serverConfig, left, right)
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume {
		t.Fatal("ticket_request resumption did not resume")
	}
	waitForTicketCount(t, cache, 2)
	_ = client.Close()
	_ = server.Close()
}

func TestTicketRequestZeroLegacyAndWeakNetwork(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, MaxSessionTickets: 4,
		HandshakeTimeout: 3 * time.Second, FlightInterval: 2 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
	}
	for _, test := range []struct {
		name    string
		request SessionTicketRequest
		want    int
		weak    bool
	}{
		{name: "ExplicitZero", request: SessionTicketRequest{Enabled: true}, want: 0},
		{name: "Legacy", want: 1},
		{name: "WeakNetwork", request: SessionTicketRequest{Enabled: true, NewSessionCount: 4, ResumptionCount: 1}, want: 4, weak: true},
	} {
		t.Run(test.name, func(t *testing.T) {
			cache := NewLRUClientSessionCache(1).(*lruSessionCache)
			clientConfig := &Config{
				RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache, SessionTicketRequest: test.request,
				HandshakeTimeout: 3 * time.Second, FlightInterval: 2 * time.Millisecond, MaxFlightInterval: 20 * time.Millisecond,
			}
			left, right := memoryDatagramPair()
			clientConn := net.Conn(left)
			serverConn := net.Conn(right)
			var clientWire, serverWire *weakNetworkConn
			if test.weak {
				clientWire = &weakNetworkConn{Conn: left, enabled: true}
				serverWire = &weakNetworkConn{Conn: right, enabled: true}
				clientConn, serverConn = clientWire, serverWire
			}
			client, server := runTicketRequestHandshake(t, clientConfig, serverConfig, clientConn, serverConn)
			waitForTicketCount(t, cache, test.want)
			if test.weak {
				clientWire.disable()
				serverWire.disable()
			}
			_ = client.Close()
			_ = server.Close()
		})
	}
}

func TestTicketRequestRealUDP(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	cache := NewLRUClientSessionCache(1).(*lruSessionCache)
	var ticketKey [32]byte
	ticketKey[0] = 1
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey, MaxSessionTickets: 3,
		HandshakeTimeout: 3 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	listener, err := Listen("udp4", "127.0.0.1:0", serverConfig)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	serverResults := make(chan bool, 2)
	serverErr := make(chan error, 1)
	go func() {
		for range 2 {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				serverErr <- acceptErr
				return
			}
			if acceptErr = conn.SetDeadline(time.Now().Add(3 * time.Second)); acceptErr == nil {
				acceptErr = conn.Handshake()
			}
			if acceptErr == nil {
				var request [8]byte
				_, _, acceptErr = conn.ReadDatagram(request[:])
			}
			if acceptErr == nil {
				_, acceptErr = conn.WriteDatagram([]byte("ok"))
			}
			serverResults <- conn.ConnectionState().DidResume
			_ = conn.Close()
			if acceptErr != nil {
				serverErr <- acceptErr
				return
			}
		}
		serverErr <- nil
	}()
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache,
		SessionTicketRequest: SessionTicketRequest{Enabled: true, NewSessionCount: 3, ResumptionCount: 1},
		HandshakeTimeout:     3 * time.Second, FlightInterval: 5 * time.Millisecond,
	}
	for connection := range 2 {
		client, dialErr := Dial("udp4", listener.Addr().String(), clientConfig)
		if dialErr != nil {
			t.Fatal(dialErr)
		}
		_ = client.SetDeadline(time.Now().Add(3 * time.Second))
		if _, dialErr = client.WriteDatagram([]byte("ping")); dialErr != nil {
			t.Fatal(dialErr)
		}
		var reply [8]byte
		n, _, dialErr := client.ReadDatagram(reply[:])
		if dialErr != nil || string(reply[:n]) != "ok" {
			t.Fatalf("real UDP reply=%q err=%v", reply[:n], dialErr)
		}
		if client.ConnectionState().DidResume != (connection == 1) {
			t.Fatalf("connection %d resumed=%v", connection+1, client.ConnectionState().DidResume)
		}
		waitForTicketCount(t, cache, 3)
		_ = client.Close()
	}
	if first, second := <-serverResults, <-serverResults; first || !second {
		t.Fatalf("server resumption states = %v, %v", first, second)
	}
	if err = <-serverErr; err != nil {
		t.Fatal(err)
	}
}

func TestTicketRequestPoolConcurrentConsumptionAndRejectedResumption(t *testing.T) {
	t.Run("ConcurrentConsumption", func(t *testing.T) {
		cache := NewLRUClientSessionCache(1).(*lruSessionCache)
		suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
		now := time.Now()
		group := [32]byte{1}
		for i := range 8 {
			cache.putTicket("server.test", &ClientSessionState{
				ticket: []byte{byte(i + 1)}, psk: bytes.Repeat([]byte{byte(i + 1)}, suite.hash.Size()),
				suite: suite.id, receivedAt: now, lifetime: 60, serverName: "server.test", ticketGroup: group,
			})
		}
		config, err := (&Config{ServerName: "server.test", ClientSessionCache: cache, Time: func() time.Time { return now }}).normalized()
		if err != nil {
			t.Fatal(err)
		}
		left, right := net.Pipe()
		defer left.Close()
		defer right.Close()
		results := make(chan string, 8)
		var wg sync.WaitGroup
		for range 8 {
			wg.Go(func() {
				state, _ := usableClientSession(config, left)
				if state != nil {
					results <- string(state.ticket)
				}
			})
		}
		wg.Wait()
		close(results)
		unique := make(map[string]bool)
		for ticket := range results {
			unique[ticket] = true
		}
		if len(unique) != 8 {
			t.Fatalf("concurrent clients consumed %d distinct tickets", len(unique))
		}
	})

	t.Run("RejectedResumption", func(t *testing.T) {
		certificate, roots := testServerCertificate(t)
		cache := NewLRUClientSessionCache(1).(*lruSessionCache)
		clientConfig := &Config{
			RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache,
			SessionTicketRequest: SessionTicketRequest{Enabled: true, NewSessionCount: 3, ResumptionCount: 1},
			HandshakeTimeout:     2 * time.Second, FlightInterval: 5 * time.Millisecond,
		}
		var firstKey [32]byte
		firstKey[0] = 1
		serverConfig := &Config{
			Certificates: []tls.Certificate{certificate}, SessionTicketKey: firstKey, MaxSessionTickets: 3,
			HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond,
		}
		left, right := memoryDatagramPair()
		client, server := runTicketRequestHandshake(t, clientConfig, serverConfig, left, right)
		waitForTicketCount(t, cache, 3)
		_ = client.Close()
		_ = server.Close()

		resumingClient := clientConfig.Clone()
		resumingClient.SessionTicketRequest = SessionTicketRequest{Enabled: true, NewSessionCount: 0, ResumptionCount: 1}
		rotatedServer := serverConfig.Clone()
		rotatedServer.SessionTicketKey[0] = 2
		left, right = memoryDatagramPair()
		client, server = runTicketRequestHandshake(t, resumingClient, rotatedServer, left, right)
		if client.ConnectionState().DidResume || server.ConnectionState().DidResume {
			t.Fatal("server accepted a ticket encrypted with the old key")
		}
		waitForTicketCount(t, cache, 0)
		_ = client.Close()
		_ = server.Close()
	})
}

func TestTicketRequestRetransmissionIsDeduplicated(t *testing.T) {
	cache := NewLRUClientSessionCache(1).(*lruSessionCache)
	config, err := (&Config{ClientSessionCache: cache, ServerName: "server.test"}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	c := &Conn{
		config: config, resumptionSuite: suite, resumptionMasterSecret: bytes.Repeat([]byte{1}, suite.hash.Size()),
		sessionTicketRequest: &sessionTicketRequestState{limit: 2, group: [32]byte{1}},
	}
	body, err := (&newSessionTicketMessage{lifetime: 60, ageAdd: 1, nonce: []byte{1}, ticket: []byte{2}}).marshal()
	if err != nil {
		t.Fatal(err)
	}
	if err = c.processNewSessionTicket(7, body); err != nil {
		t.Fatal(err)
	}
	if err = c.processNewSessionTicket(7, body); err != nil {
		t.Fatal(err)
	}
	if states := ticketCacheSnapshot(cache, "server.test"); len(states) != 1 {
		t.Fatalf("cached a retransmitted ticket %d times", len(states))
	}
}
