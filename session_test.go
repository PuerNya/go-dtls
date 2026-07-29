package dtls13

import (
	"bytes"
	"crypto/hkdf"
	"crypto/hmac"
	"crypto/tls"
	"crypto/x509"
	"net"
	"sync"
	"testing"
	"time"
)

func TestCalculatePSKBinderIntoMatchesIndependentDerivation(t *testing.T) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			t.Fatal(err)
		}
		psk := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		transcript := bytes.Repeat([]byte{0xa5}, suite.hash.Size())
		earlySecret, err := hkdf.Extract(suite.hash.New, psk, nil)
		if err != nil {
			t.Fatal(err)
		}
		binderKey := expandLabel(suite, earlySecret, "res binder", emptyTranscriptHash(suite), suite.hash.Size())
		finishedKey := expandLabel(suite, binderKey, "finished", nil, suite.hash.Size())
		mac := hmac.New(suite.hash.New, finishedKey)
		_, _ = mac.Write(transcript)
		want := mac.Sum(nil)

		out := append([]byte(nil), psk...)
		calculatePSKBinderInto(suite, out, transcript, out)
		if !bytes.Equal(out, want) {
			t.Fatalf("suite %04x in-place PSK binder differs from independent derivation", suiteID)
		}
	}
}

func TestNewSessionTicketRoundTrip(t *testing.T) {
	want := &newSessionTicketMessage{
		lifetime: 3600, ageAdd: 0x10203040, nonce: []byte{1, 2, 3},
		ticket: []byte("opaque ticket"), maxEarlyData: 4096,
	}
	wire, err := want.marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseNewSessionTicket(wire)
	if err != nil {
		t.Fatal(err)
	}
	if got.lifetime != want.lifetime || got.ageAdd != want.ageAdd ||
		!bytes.Equal(got.nonce, want.nonce) || !bytes.Equal(got.ticket, want.ticket) ||
		got.maxEarlyData != want.maxEarlyData {
		t.Fatalf("got %#v", got)
	}
	if _, err = (&newSessionTicketMessage{lifetime: 8 * 24 * 60 * 60, ticket: []byte{1}}).marshal(); err == nil {
		t.Fatal("accepted a ticket lifetime over seven days")
	}
	if _, err = parseNewSessionTicket(wire[:len(wire)-1]); err == nil {
		t.Fatal("accepted a truncated NewSessionTicket")
	}
}

func TestNewSessionTicketIgnoresUnknownExtensions(t *testing.T) {
	exts, err := marshalExtensions(map[uint16][]byte{0xffa5: {1, 2, 3}}, []uint16{0xffa5})
	if err != nil {
		t.Fatal(err)
	}
	var w wireBuilder
	w.u32(60)
	w.u32(7)
	w.bytes8([]byte{1})
	w.bytes16([]byte{2})
	w.b = append(w.b, exts...)
	message, err := parseNewSessionTicket(w.b)
	if err != nil {
		t.Fatal(err)
	}
	if message.lifetime != 60 || message.ageAdd != 7 || message.maxEarlyData != 0 {
		t.Fatalf("unexpected ticket: %#v", message)
	}
}

func TestNewSessionTicketRejectsRecognizedExtensionInWrongMessage(t *testing.T) {
	exts, err := marshalExtensions(map[uint16][]byte{extALPN: nil}, []uint16{extALPN})
	if err != nil {
		t.Fatal(err)
	}
	var w wireBuilder
	w.u32(60)
	w.u32(7)
	w.bytes8([]byte{1})
	w.bytes16([]byte{2})
	w.b = append(w.b, exts...)
	_, err = parseNewSessionTicket(w.b)
	if description, ok := protocolAlert(err); !ok || description != alertIllegalParameter {
		t.Fatalf("wrong-message extension alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestNewSessionTicketAcceptsZeroLifetimeForImmediateDiscard(t *testing.T) {
	var w wireBuilder
	w.u32(0)
	w.u32(7)
	w.bytes8([]byte{1})
	w.bytes16([]byte{2})
	exts, err := marshalExtensions(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	w.b = append(w.b, exts...)
	message, err := parseNewSessionTicket(w.b)
	if err != nil {
		t.Fatal(err)
	}
	if message.lifetime != 0 {
		t.Fatalf("lifetime=%d", message.lifetime)
	}
}

func TestNewSessionTicketOverlongLifetimeIsParsedAndDiscarded(t *testing.T) {
	var w wireBuilder
	w.u32(8 * 24 * 60 * 60)
	w.u32(7)
	w.bytes8([]byte{1})
	w.bytes16([]byte{2})
	exts, err := marshalExtensions(nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	w.b = append(w.b, exts...)
	message, err := parseNewSessionTicket(w.b)
	if err != nil || message.lifetime != 8*24*60*60 {
		t.Fatalf("parsed ticket=%#v err=%v", message, err)
	}
	cache := NewLRUClientSessionCache(1)
	config, err := (&Config{ClientSessionCache: cache, ServerName: "server.test"}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	c := &Conn{config: config, resumptionSuite: suite, resumptionMasterSecret: make([]byte, suite.hash.Size())}
	if err = c.processNewSessionTicket(w.b); err != nil {
		t.Fatal(err)
	}
	if _, ok := cache.Get("server.test"); ok {
		t.Fatal("cached a ticket with lifetime over seven days")
	}
}

func TestSessionTicketProtectionTamperAndExpiry(t *testing.T) {
	now := time.Unix(1700000000, 0)
	var key [32]byte
	copy(key[:], bytes.Repeat([]byte{9}, 32))
	protector, err := newSessionTicketProtector(key, bytes.NewReader(bytes.Repeat([]byte{5}, 64)), func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	want := &sessionTicketState{
		createdAt: now.Unix(), lifetime: 60, suite: TLS_AES_128_GCM_SHA256,
		psk: bytes.Repeat([]byte{7}, 32), serverName: "server.test", protocol: "coap",
	}
	ticket, err := protector.seal(want)
	if err != nil {
		t.Fatal(err)
	}
	got, err := protector.open(ticket)
	if err != nil {
		t.Fatal(err)
	}
	if got.suite != want.suite || got.serverName != want.serverName || got.protocol != want.protocol || !bytes.Equal(got.psk, want.psk) {
		t.Fatalf("got %#v", got)
	}
	tampered := append([]byte(nil), ticket...)
	tampered[len(tampered)-1] ^= 1
	if _, err = protector.open(tampered); err == nil {
		t.Fatal("accepted a tampered ticket")
	}
	now = now.Add(61 * time.Second)
	if _, err = protector.open(ticket); err == nil {
		t.Fatal("accepted an expired ticket")
	}
}

func TestSessionTicketStateClientAuthenticationRoundTrip(t *testing.T) {
	certificate, _ := testClientCertificate(t)
	leaf, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	want := &sessionTicketState{
		createdAt: 1700000000, lifetime: 3600, suite: TLS_AES_128_GCM_SHA256,
		psk: bytes.Repeat([]byte{0x42}, 32), serverName: "server.test", protocol: "coap", ageAdd: 7,
		maxEarlyData: 1024, recordSizeLimit: 256, clientAuthAt: 1699999990,
		peerCertificates: []*x509.Certificate{leaf}, verifiedChains: [][]*x509.Certificate{{leaf}},
	}
	wire, err := want.marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := parseSessionTicketState(wire)
	if err != nil {
		t.Fatal(err)
	}
	if got.createdAt != want.createdAt || got.clientAuthAt != want.clientAuthAt || got.maxEarlyData != want.maxEarlyData || got.recordSizeLimit != want.recordSizeLimit || len(got.peerCertificates) != 1 || len(got.verifiedChains) != 1 || len(got.verifiedChains[0]) != 1 || !got.peerCertificates[0].Equal(leaf) || !got.verifiedChains[0][0].Equal(leaf) {
		t.Fatalf("unexpected client authentication ticket state: %#v", got)
	}
	for length := range len(wire) {
		if _, err = parseSessionTicketState(wire[:length]); err == nil {
			t.Fatalf("accepted client authentication ticket truncated to %d bytes", length)
		}
	}
	if _, err = parseSessionTicketState(append(append([]byte(nil), wire...), 0)); err == nil {
		t.Fatal("accepted client authentication ticket with trailing data")
	}
}

func TestClientAuthenticationTicketPolicy(t *testing.T) {
	certificate, roots := testClientCertificate(t)
	leaf, err := x509.ParseCertificate(certificate.Certificate[0])
	if err != nil {
		t.Fatal(err)
	}
	now := leaf.NotBefore.Add(time.Minute)
	state := &sessionTicketState{
		clientAuthAt: now.Add(-time.Minute).Unix(), peerCertificates: []*x509.Certificate{leaf},
		verifiedChains: [][]*x509.Certificate{{leaf}},
	}
	config := func(policy tls.ClientAuthType, clientRoots *x509.CertPool) *Config {
		cfg, err := (&Config{
			ClientAuth: policy, ClientCAs: clientRoots, SessionTicketLifetime: time.Hour,
			Time: func() time.Time { return now },
		}).normalized()
		if err != nil {
			t.Fatal(err)
		}
		return cfg
	}
	tests := []struct {
		name   string
		policy tls.ClientAuthType
		state  *sessionTicketState
		roots  *x509.CertPool
		want   bool
	}{
		{"NoClientCertWithoutCertificate", tls.NoClientCert, &sessionTicketState{}, roots, true},
		{"NoClientCertWithCertificate", tls.NoClientCert, state, roots, false},
		{"RequestWithoutCertificate", tls.RequestClientCert, &sessionTicketState{}, roots, false},
		{"RequireAnyWithoutCertificate", tls.RequireAnyClientCert, &sessionTicketState{}, roots, false},
		{"RequireAnyWithCertificate", tls.RequireAnyClientCert, &sessionTicketState{clientAuthAt: state.clientAuthAt, peerCertificates: state.peerCertificates}, roots, true},
		{"VerifyIfGivenWithoutCertificate", tls.VerifyClientCertIfGiven, &sessionTicketState{}, roots, false},
		{"VerifyIfGivenWithoutChain", tls.VerifyClientCertIfGiven, &sessionTicketState{clientAuthAt: state.clientAuthAt, peerCertificates: state.peerCertificates}, roots, false},
		{"VerifyIfGivenWithChain", tls.VerifyClientCertIfGiven, state, roots, true},
		{"RequireAndVerifyWithChain", tls.RequireAndVerifyClientCert, state, roots, true},
		{"RequireAndVerifyWithDifferentCA", tls.RequireAndVerifyClientCert, state, x509.NewCertPool(), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := validClientAuthenticationTicket(config(test.policy, test.roots), test.state); got != test.want {
				t.Fatalf("validClientAuthenticationTicket = %v, want %v", got, test.want)
			}
		})
	}

	expiredAuthentication := *state
	expiredAuthentication.clientAuthAt = now.Add(-2 * time.Hour).Unix()
	if validClientAuthenticationTicket(config(tls.RequireAndVerifyClientCert, roots), &expiredAuthentication) {
		t.Fatal("accepted client authentication older than the configured total ticket lifetime")
	}
	expiredCertificate := *state
	expiredCertificate.clientAuthAt = leaf.NotAfter.Add(-time.Minute).Unix()
	now = leaf.NotAfter.Add(time.Minute)
	if validClientAuthenticationTicket(config(tls.RequireAndVerifyClientCert, roots), &expiredCertificate) {
		t.Fatal("accepted an expired client certificate")
	}
}

func TestLRUClientSessionCacheIsBoundedAndClones(t *testing.T) {
	cache := NewLRUClientSessionCache(2).(*lruSessionCache)
	one := &ClientSessionState{ticket: []byte{1}, psk: []byte{2}}
	cache.Put("one", one)
	one.ticket[0] = 9
	got, ok := cache.Get("one")
	if !ok || got.ticket[0] != 1 {
		t.Fatal("cache did not clone inserted session")
	}
	got.ticket[0] = 8
	again, _ := cache.Get("one")
	if again.ticket[0] != 1 {
		t.Fatal("cache exposed mutable stored session")
	}
	cache.Put("two", &ClientSessionState{ticket: []byte{2}})
	cache.Put("three", &ClientSessionState{ticket: []byte{3}})
	if _, ok = cache.Get("one"); ok {
		t.Fatal("least recently used session was not evicted")
	}
	if cache.order.Len() != 2 || len(cache.entries) != 2 {
		t.Fatalf("cache grew beyond capacity: list=%d map=%d", cache.order.Len(), len(cache.entries))
	}
	cache.Put("two", nil)
	if _, ok = cache.Get("two"); ok {
		t.Fatal("nil Put did not remove session")
	}
}

func TestClientSessionTicketIsConsumedOnUse(t *testing.T) {
	cache := NewLRUClientSessionCache(1)
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	now := time.Now()
	cache.Put("server.test", &ClientSessionState{
		ticket: []byte("ticket"), psk: bytes.Repeat([]byte{1}, suite.hash.Size()),
		suite: suite.id, receivedAt: now, lifetime: 60, serverName: "server.test",
	})
	config, err := (&Config{ServerName: "server.test", ClientSessionCache: cache, Time: func() time.Time { return now }}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	left, right := net.Pipe()
	defer left.Close()
	defer right.Close()
	if state, selected := usableClientSession(config, left); state == nil || selected == nil || selected.id != suite.id {
		t.Fatal("valid client ticket was not selected")
	}
	if _, ok := cache.Get("server.test"); ok {
		t.Fatal("selected client ticket remained in the cache")
	}
	if state, selected := usableClientSession(config, left); state != nil || selected != nil {
		t.Fatal("consumed client ticket was selected twice")
	}

	cache.Put("server.test", &ClientSessionState{
		ticket: []byte("ticket-2"), psk: bytes.Repeat([]byte{2}, suite.hash.Size()),
		suite: suite.id, receivedAt: now, lifetime: 60, serverName: "server.test",
	})
	const contenders = 32
	results := make(chan bool, contenders)
	var wg sync.WaitGroup
	for range contenders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			state, _ := usableClientSession(config, left)
			results <- state != nil
		}()
	}
	wg.Wait()
	close(results)
	selectedCount := 0
	for selected := range results {
		if selected {
			selectedCount++
		}
	}
	if selectedCount != 1 {
		t.Fatalf("concurrent connections selected one ticket %d times", selectedCount)
	}
}

func TestEarlyDataReplayCacheIsConcurrentAndBounded(t *testing.T) {
	cacheInterface := NewLRUEarlyDataReplayCache(2)
	cache, ok := cacheInterface.(*earlyDataReplayCache)
	if !ok {
		t.Fatal("unexpected replay cache implementation")
	}
	expires := time.Now().Add(time.Hour)
	const contenders = 32
	var wg sync.WaitGroup
	results := make(chan bool, contenders)
	for range contenders {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results <- cache.CheckAndStore("same-ticket", expires)
		}()
	}
	wg.Wait()
	close(results)
	accepted := 0
	for result := range results {
		if result {
			accepted++
		}
	}
	if accepted != 1 {
		t.Fatalf("accepted identical early-data ticket %d times", accepted)
	}
	if !cache.CheckAndStore("ticket-two", expires) {
		t.Fatal("replay cache rejected a fresh identity within capacity")
	}
	if cache.CheckAndStore("ticket-three", expires) {
		t.Fatal("replay cache accepted new early data after reaching capacity")
	}
	if cache.CheckAndStore("same-ticket", expires) {
		t.Fatal("capacity pressure made an unexpired identity replayable")
	}
	cache.mu.Lock()
	entries := len(cache.entries)
	cache.mu.Unlock()
	if entries != 2 {
		t.Fatalf("replay cache retained %d entries, want 2", entries)
	}
	expiryCache := NewLRUEarlyDataReplayCache(1)
	if !expiryCache.CheckAndStore("already-expired", time.Now().Add(-time.Second)) {
		t.Fatal("replay cache rejected a fresh identity with an expired retention deadline")
	}
	if !expiryCache.CheckAndStore("already-expired", expires) {
		t.Fatal("expired identity was not reusable")
	}
}

func TestPSKBinderInitialAndHelloRetryRequest(t *testing.T) {
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	psk := bytes.Repeat([]byte{0x42}, suite.hash.Size())
	hello := &clientHello{
		cipherSuites: []uint16{suite.id}, supportedGroups: []tls.CurveID{tls.X25519},
		keyShares:        []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
		signatureSchemes: defaultSignatureSchemes(), pskIdentity: []byte("ticket"), obfuscatedAge: 1234,
	}
	initialBody, err := marshalClientHelloWithPSKBinder(hello, suite, psk, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = verifyClientHelloPSKBinder(initialBody, suite, psk, nil, nil); err != nil {
		t.Fatal(err)
	}
	tampered := append([]byte(nil), initialBody...)
	tampered[len(tampered)-1] ^= 1
	if err = verifyClientHelloPSKBinder(tampered, suite, psk, nil, nil); err == nil {
		t.Fatal("accepted a tampered initial PSK binder")
	} else if description, ok := protocolAlert(err); !ok || description != alertDecryptError {
		t.Fatalf("binder alert=%d ok=%v err=%v", description, ok, err)
	}
	initialTranscript := newTranscriptHash(suite.hash.New())
	_ = initialTranscript.add(handshakeTypeClientHello, 0, initialBody)
	hrrBody := []byte("synthetic HRR body")
	hello.cookie = []byte("cookie")
	secondBody, err := marshalClientHelloWithPSKBinder(hello, suite, psk, initialTranscript.sum(), hrrBody)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(parseBinderForTest(t, initialBody), parseBinderForTest(t, secondBody)) {
		t.Fatal("HRR did not change the PSK binder")
	}
	if err = verifyClientHelloPSKBinder(secondBody, suite, psk, initialTranscript.sum(), hrrBody); err != nil {
		t.Fatal(err)
	}
}

func parseBinderForTest(t *testing.T, body []byte) []byte {
	t.Helper()
	hello, err := parseClientHello(body)
	if err != nil {
		t.Fatal(err)
	}
	return hello.pskBinder
}

func TestServerSelectsSecondCompatiblePSKIdentity(t *testing.T) {
	now := time.Unix(1700000000, 0)
	var ticketKey [32]byte
	for i := range ticketKey {
		ticketKey[i] = byte(i + 1)
	}
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	psk := bytes.Repeat([]byte{0x77}, suite.hash.Size())
	protector, err := newSessionTicketProtector(ticketKey, bytes.NewReader(bytes.Repeat([]byte{5}, 64)), func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	ticket, err := protector.seal(&sessionTicketState{createdAt: now.Unix(), lifetime: 60, suite: suite.id, psk: psk, serverName: "server.test", ageAdd: 9})
	if err != nil {
		t.Fatal(err)
	}
	hello := &clientHello{
		cipherSuites: []uint16{suite.id}, supportedGroups: []tls.CurveID{tls.X25519}, serverName: "server.test",
		keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}, signatureSchemes: defaultSignatureSchemes(),
		pskIdentities: []pskIdentityEntry{{identity: []byte("unknown")}, {identity: ticket, obfuscatedAge: 9}},
		pskBinders:    [][]byte{make([]byte, suite.hash.Size()), make([]byte, suite.hash.Size())},
	}
	body, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	entriesLength := 0
	for _, binder := range hello.pskBinders {
		entriesLength += 1 + len(binder)
	}
	truncated, err := truncateClientHelloForBinderEntries(body, entriesLength)
	if err != nil {
		t.Fatal(err)
	}
	h := suite.hash.New()
	writeBinderTranscriptMessage(h, handshakeTypeClientHello, truncated, len(body))
	transcriptHash := h.Sum(nil)
	hello.pskBinders[0] = calculatePSKBinder(suite, bytes.Repeat([]byte{0x11}, suite.hash.Size()), transcriptHash)
	hello.pskBinders[1] = calculatePSKBinder(suite, psk, transcriptHash)
	body, err = hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseClientHello(body)
	if err != nil {
		t.Fatal(err)
	}
	config, err := (&Config{SessionTicketKey: ticketKey, Time: func() time.Time { return now }}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	c := &Conn{config: config}
	selectedSession, selected, err := c.acceptSessionTicket(parsed, body, nil, nil, suite, "")
	if err != nil {
		t.Fatal(err)
	}
	if selectedSession == nil || selected != 1 || !bytes.Equal(selectedSession.psk, psk) {
		t.Fatalf("resumed=%v selected=%d psk_match=%v", selectedSession != nil, selected, selectedSession != nil && bytes.Equal(selectedSession.psk, psk))
	}
	body = replaceClientHelloExtension(t, body, extPSKKeyExchangeModes, []byte{1, 0})
	parsed, err = parseClientHello(body)
	if err != nil {
		t.Fatal(err)
	}
	selectedSession, _, err = c.acceptSessionTicket(parsed, body, nil, nil, suite, "")
	if err != nil || selectedSession != nil {
		t.Fatalf("unsupported PSK mode resumed=%v err=%v", selectedSession != nil, err)
	}
}

func TestSecondClientHelloCannotRemoveCompatiblePSK(t *testing.T) {
	now := time.Unix(1700000000, 0)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x4a}, 32))
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	protector, err := newSessionTicketProtector(ticketKey, bytes.NewReader(bytes.Repeat([]byte{5}, 64)), func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	ticket, err := protector.seal(&sessionTicketState{
		createdAt: now.Unix(), lifetime: 60, suite: suite.id, psk: bytes.Repeat([]byte{7}, suite.hash.Size()), serverName: "server.test",
	})
	if err != nil {
		t.Fatal(err)
	}
	initial := &clientHello{pskIdentities: []pskIdentityEntry{{identity: ticket}, {identity: []byte("unknown")}}}
	second := &clientHello{pskIdentities: []pskIdentityEntry{{identity: []byte("unknown")}}}
	config, err := (&Config{SessionTicketKey: ticketKey, Time: func() time.Time { return now }}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	ok, err := removedPSKsAreIncompatible(config, initial, second, suite)
	if err != nil || ok {
		t.Fatalf("compatible removal accepted=%v err=%v", ok, err)
	}
	second.pskIdentities = nil
	initial.pskIdentities = []pskIdentityEntry{{identity: []byte("unknown")}}
	ok, err = removedPSKsAreIncompatible(config, initial, second, suite)
	if err != nil || !ok {
		t.Fatalf("unknown removal accepted=%v err=%v", ok, err)
	}
}

func TestTicketAgeMismatchRejectsOnlyEarlyData(t *testing.T) {
	now := time.Unix(1700000000, 0)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{3}, 32))
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	psk := bytes.Repeat([]byte{0x44}, suite.hash.Size())
	protector, err := newSessionTicketProtector(ticketKey, bytes.NewReader(bytes.Repeat([]byte{5}, 64)), func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	ticket, err := protector.seal(&sessionTicketState{
		createdAt: now.Unix(), lifetime: 3600, suite: suite.id, psk: psk,
		serverName: "server.test", ageAdd: 7, maxEarlyData: 1024,
	})
	if err != nil {
		t.Fatal(err)
	}
	hello := &clientHello{
		cipherSuites: []uint16{suite.id}, supportedGroups: []tls.CurveID{tls.X25519}, serverName: "server.test",
		keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}, signatureSchemes: defaultSignatureSchemes(),
		pskIdentity: ticket, obfuscatedAge: 7 + uint32(20*time.Minute/time.Millisecond),
		pskBinder: make([]byte, suite.hash.Size()), earlyData: true,
	}
	body, err := hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	truncated, err := truncateClientHelloForBinder(body, suite.hash.Size())
	if err != nil {
		t.Fatal(err)
	}
	h := suite.hash.New()
	writeBinderTranscriptMessage(h, handshakeTypeClientHello, truncated, len(body))
	hello.pskBinder = calculatePSKBinder(suite, psk, h.Sum(nil))
	body, err = hello.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseClientHello(body)
	if err != nil {
		t.Fatal(err)
	}
	config, err := (&Config{
		SessionTicketKey: ticketKey, Time: func() time.Time { return now }, MaxEarlyData: 1024,
		AllowEarlyDataWithoutCookie: true,
	}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	c := &Conn{config: config}
	selectedSession, selected, err := c.acceptSessionTicket(parsed, body, nil, nil, suite, "")
	if err != nil {
		t.Fatal(err)
	}
	if selectedSession == nil || selected != 0 || !bytes.Equal(selectedSession.psk, psk) {
		t.Fatal("ticket age mismatch incorrectly disabled 1-RTT resumption")
	}
	if c.earlyAccepted || c.earlyDataLimit != 0 {
		t.Fatal("ticket age mismatch accepted early data")
	}
}
