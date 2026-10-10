package dtls13

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"runtime"
	"testing"
	"time"
)

type recordBenchmarkConfig struct {
	name     string
	suiteID  uint16
	cidBytes int
}

var recordBenchmarkSuites = []recordBenchmarkConfig{
	{name: "AES128GCM", suiteID: TLS_AES_128_GCM_SHA256},
	{name: "AES256GCM", suiteID: TLS_AES_256_GCM_SHA384},
	{name: "ChaCha20Poly1305", suiteID: TLS_CHACHA20_POLY1305_SHA256},
	{name: "AES128CCM", suiteID: TLS_AES_128_CCM_SHA256},
}

func BenchmarkConnectionHandshakeLifecycle(b *testing.B) {
	certificate, roots := testServerCertificate(b)
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true,
		HandshakeTimeout: time.Second,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true,
		HandshakeTimeout: time.Second,
	}
	b.ReportAllocs()

	for b.Loop() {
		left, right := memoryDatagramPair()
		client := Client(left, clientConfig)
		server := Server(right, serverConfig)
		serverDone := make(chan error, 1)
		go func() { serverDone <- server.Handshake() }()
		clientErr := client.Handshake()
		serverErr := <-serverDone
		_ = left.Close()
		_ = right.Close()
		if clientErr != nil || serverErr != nil {
			b.Fatalf("handshake failed: client=%v server=%v", clientErr, serverErr)
		}
	}
}

func BenchmarkSessionTicketRequestHandshakeLifecycle(b *testing.B) {
	certificate, roots := testServerCertificate(b)
	var ticketKey [32]byte
	ticketKey[0] = 1
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", SessionTicketRequest: SessionTicketRequest{Enabled: true, NewSessionCount: 4, ResumptionCount: 1},
		HandshakeTimeout: time.Second,
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey, MaxSessionTickets: 4,
		HandshakeTimeout: time.Second,
	}
	b.ReportAllocs()
	for b.Loop() {
		cache := NewLRUClientSessionCache(1).(*lruSessionCache)
		clientConfig.ClientSessionCache = cache
		left, right := memoryDatagramPair()
		client := Client(left, clientConfig)
		server := Server(right, serverConfig)
		serverDone := make(chan error, 1)
		go func() { serverDone <- server.Handshake() }()
		clientErr := client.Handshake()
		serverErr := <-serverDone
		if clientErr == nil && serverErr == nil {
			deadline := time.Now().Add(time.Second)
			for {
				cache.mu.Lock()
				element := cache.entries["server.test"]
				count := 0
				if element != nil {
					count = len(element.Value.(*sessionCacheEntry).states)
				}
				cache.mu.Unlock()
				server.writeMu.Lock()
				acknowledged := server.ticketFlight == nil
				server.writeMu.Unlock()
				if count == 4 && acknowledged {
					break
				}
				if time.Now().After(deadline) {
					b.Fatalf("ticket lifecycle timed out: cached=%d acknowledged=%v", count, acknowledged)
				}
				runtime.Gosched()
			}
		}
		_ = client.Close()
		_ = server.Close()
		if clientErr != nil || serverErr != nil {
			b.Fatalf("ticket_request handshake failed: client=%v server=%v", clientErr, serverErr)
		}
	}
}

func BenchmarkGREASEHandshakeLifecycle(b *testing.B) {
	serverCertificate, roots := testServerCertificate(b)
	clientCertificate, clientRoots := testClientCertificate(b)
	for _, test := range []struct {
		name    string
		enabled bool
	}{
		{name: "Disabled"},
		{name: "Enabled", enabled: true},
	} {
		b.Run(test.name, func(b *testing.B) {
			clientConfig := &Config{
				RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate},
				EnableGREASE: test.enabled, HandshakeTimeout: time.Second,
			}
			serverConfig := &Config{
				Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert,
				ClientCAs: clientRoots, EnableGREASE: test.enabled, HandshakeTimeout: time.Second,
			}
			b.ReportAllocs()
			for b.Loop() {
				cache := NewLRUClientSessionCache(1).(*lruSessionCache)
				clientConfig.ClientSessionCache = cache
				left, right := memoryDatagramPair()
				client := Client(left, clientConfig)
				server := Server(right, serverConfig)
				serverDone := make(chan error, 1)
				go func() { serverDone <- server.Handshake() }()
				clientErr := client.Handshake()
				serverErr := <-serverDone
				if clientErr == nil && serverErr == nil {
					deadline := time.Now().Add(time.Second)
					for {
						cache.mu.Lock()
						element := cache.entries["server.test"]
						cached := element != nil && len(element.Value.(*sessionCacheEntry).states) == 1
						cache.mu.Unlock()
						server.writeMu.Lock()
						acknowledged := server.ticketFlight == nil
						server.writeMu.Unlock()
						if cached && acknowledged {
							break
						}
						if time.Now().After(deadline) {
							b.Fatalf("GREASE session ticket timed out: cached=%v acknowledged=%v", cached, acknowledged)
						}
						runtime.Gosched()
					}
				}
				_ = client.Close()
				_ = server.Close()
				if clientErr != nil || serverErr != nil {
					b.Fatalf("GREASE handshake failed: client=%v server=%v", clientErr, serverErr)
				}
			}
		})
	}
}

func BenchmarkECHHandshakeLifecycle(b *testing.B) {
	certificate, roots := testServerCertificate(b)
	configList, key := testECHConfig(b, "public.test", 1)
	for _, test := range []struct {
		name   string
		curves []tls.CurveID
	}{
		{name: "Direct"},
		{name: "HRR", curves: []tls.CurveID{tls.CurveP256}},
	} {
		b.Run(test.name, func(b *testing.B) {
			clientConfig := &Config{
				RootCAs: roots, ServerName: "server.test", EncryptedClientHelloConfigList: configList,
				SessionTicketsDisabled: true, HandshakeTimeout: time.Second,
			}
			serverConfig := &Config{
				Certificates: []tls.Certificate{certificate}, EncryptedClientHelloKeys: []EncryptedClientHelloKey{key},
				CurvePreferences: test.curves, SessionTicketsDisabled: true, HandshakeTimeout: time.Second,
			}
			b.ReportAllocs()
			for b.Loop() {
				left, right := memoryDatagramPair()
				client := Client(left, clientConfig)
				server := Server(right, serverConfig)
				serverDone := make(chan error, 1)
				go func() { serverDone <- server.Handshake() }()
				clientErr := client.Handshake()
				serverErr := <-serverDone
				_ = left.Close()
				_ = right.Close()
				if clientErr != nil || serverErr != nil {
					b.Fatalf("ECH handshake failed: client=%v server=%v", clientErr, serverErr)
				}
			}
		})
	}
}

func BenchmarkExternalPSKHandshakeLifecycle(b *testing.B) {
	psk, err := ImportExternalPSK([]byte("benchmark-device"), bytes.Repeat([]byte{0x5d}, 32), []byte("client=benchmark;server=benchmark"), 0)
	if err != nil {
		b.Fatal(err)
	}
	config := &Config{ExternalPSKs: []*ExternalPSK{psk}, SessionTicketsDisabled: true, HandshakeTimeout: time.Second}
	b.ReportAllocs()
	for b.Loop() {
		left, right := memoryDatagramPair()
		client := Client(left, config)
		server := Server(right, config)
		serverDone := make(chan error, 1)
		go func() { serverDone <- server.Handshake() }()
		clientErr := client.Handshake()
		serverErr := <-serverDone
		_ = left.Close()
		_ = right.Close()
		if clientErr != nil || serverErr != nil {
			b.Fatalf("external PSK handshake failed: client=%v server=%v", clientErr, serverErr)
		}
	}
}

func BenchmarkMutualTLSHandshakeLifecycle(b *testing.B) {
	serverCertificate, roots := testServerCertificate(b)
	clientCertificate, clientRoots := testClientCertificate(b)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x6e}, 32))
	serverConfig := &Config{
		Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: clientRoots, SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour,
		HandshakeTimeout: time.Second,
	}

	b.Run("Full", func(b *testing.B) {
		fullServerConfig := serverConfig.Clone()
		fullServerConfig.SessionTicketsDisabled = true
		clientConfig := &Config{
			RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate},
			SessionTicketsDisabled: true, HandshakeTimeout: time.Second,
		}
		b.ReportAllocs()
		for b.Loop() {
			left, right := memoryDatagramPair()
			client := Client(left, clientConfig)
			server := Server(right, fullServerConfig)
			serverDone := make(chan error, 1)
			go func() { serverDone <- server.Handshake() }()
			clientErr := client.Handshake()
			serverErr := <-serverDone
			_ = left.Close()
			_ = right.Close()
			if clientErr != nil || serverErr != nil {
				b.Fatalf("full mutual TLS handshake failed: client=%v server=%v", clientErr, serverErr)
			}
		}
	})

	cache := NewLRUClientSessionCache(1)
	clientConfig := &Config{
		RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{clientCertificate},
		ClientSessionCache: cache, HandshakeTimeout: time.Second,
	}
	left, right := memoryDatagramPair()
	client := Client(left, clientConfig)
	server := Server(right, serverConfig)
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	if err := client.Handshake(); err != nil {
		b.Fatal(err)
	}
	if err := <-serverDone; err != nil {
		b.Fatal(err)
	}
	var session *ClientSessionState
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if state, ok := cache.Get("server.test"); ok {
			session = state
			cache.Put("server.test", nil)
			break
		}
		time.Sleep(time.Millisecond)
	}
	_ = client.Close()
	_ = server.Close()
	if session == nil {
		b.Fatal("initial mutual TLS handshake did not produce a session ticket")
	}
	clientConfig.Certificates = nil

	b.Run("Resumed", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			cache.Put("server.test", session)
			left, right := memoryDatagramPair()
			client := Client(left, clientConfig)
			server := Server(right, serverConfig)
			serverDone := make(chan error, 1)
			go func() { serverDone <- server.Handshake() }()
			clientErr := client.Handshake()
			serverErr := <-serverDone
			resumed := client.ConnectionState().DidResume && server.ConnectionState().DidResume
			_ = left.Close()
			_ = right.Close()
			if clientErr != nil || serverErr != nil || !resumed {
				b.Fatalf("resumed mutual TLS handshake failed: client=%v server=%v resumed=%v", clientErr, serverErr, resumed)
			}
		}
	})
}

func BenchmarkCertificateSelectionHandshakeLifecycle(b *testing.B) {
	serverCertificate, roots := testServerCertificate(b)
	wrongServer := testSelectionCertificate(b, 1000, "wrong-server", "server.test", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageServerAuth)
	wrongClient := testSelectionCertificate(b, 1001, "wrong-client", "", x509.KeyUsageDigitalSignature, x509.ExtKeyUsageClientAuth)
	rightClient := testSelectionCertificate(b, 1002, "right-client", "", x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment, x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageCodeSigning)
	filters := []CertificateOIDFilter{
		keyUsageFilter(b, x509.KeyUsageDigitalSignature|x509.KeyUsageKeyEncipherment),
		extendedKeyUsageFilter(b, oidExtKeyUsageCodeSigning),
	}
	serverConfig := &Config{
		Certificates: []tls.Certificate{wrongServer, serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certificatePool(rightClient), ClientCertificateOIDFilters: filters,
		SessionTicketsDisabled: true, HandshakeTimeout: time.Second,
	}
	run := func(b *testing.B, postHandshake bool) {
		clientConfig := &Config{
			RootCAs: roots, ServerName: "server.test", Certificates: []tls.Certificate{wrongClient, rightClient},
			ServerCertificateAuthorities: [][]byte{serverCertificate.Leaf.RawSubject},
			PostHandshakeAuth:            postHandshake, SessionTicketsDisabled: true, HandshakeTimeout: time.Second,
		}
		b.ReportAllocs()
		for b.Loop() {
			left, right := memoryDatagramPair()
			client := Client(left, clientConfig)
			server := Server(right, serverConfig)
			serverDone := make(chan error, 1)
			go func() { serverDone <- server.Handshake() }()
			clientErr := client.Handshake()
			serverErr := <-serverDone
			if clientErr == nil && serverErr == nil && postHandshake {
				serverErr = server.RequestClientCertificate(context.Background())
			}
			_ = left.Close()
			_ = right.Close()
			if clientErr != nil || serverErr != nil {
				b.Fatalf("certificate selection handshake failed: client=%v server=%v", clientErr, serverErr)
			}
		}
	}
	b.Run("InitialMutualTLS", func(b *testing.B) { run(b, false) })
	b.Run("PostHandshakeAuthentication", func(b *testing.B) { run(b, true) })
}

func newBenchmarkRecordCipher(b *testing.B, suiteID uint16, secretByte byte, cidBytes int) *recordCipher {
	b.Helper()
	suite, err := cipherSuiteForID(suiteID)
	if err != nil {
		b.Fatal(err)
	}
	secret := bytes.Repeat([]byte{secretByte}, suite.hash.Size())
	cipher, err := newRecordCipher(suite, secret, 3, 64)
	if err != nil {
		b.Fatal(err)
	}
	if cidBytes > 0 {
		if err = cipher.setConnectionID(bytes.Repeat([]byte{0xa5}, cidBytes)); err != nil {
			b.Fatal(err)
		}
	}
	return cipher
}

func benchmarkProtectedRecordSeal(b *testing.B, config recordBenchmarkConfig) {
	cipher := newBenchmarkRecordCipher(b, config.suiteID, 1, config.cidBytes)
	payload := bytes.Repeat([]byte{2}, 1200)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := cipher.seal(recordTypeApplicationData, payload); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkProtectedRecordRoundTrip(b *testing.B, config recordBenchmarkConfig) {
	sender := newBenchmarkRecordCipher(b, config.suiteID, 3, config.cidBytes)
	receiver := newBenchmarkRecordCipher(b, config.suiteID, 3, config.cidBytes)
	payload := bytes.Repeat([]byte{4}, 1200)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wire, sealErr := sender.seal(recordTypeApplicationData, payload)
		if sealErr != nil {
			b.Fatal(sealErr)
		}
		if _, _, _, openErr := receiver.open(wire); openErr != nil {
			b.Fatal(openErr)
		}
	}
}

func BenchmarkProtectedRecordSeal(b *testing.B) {
	benchmarkProtectedRecordSeal(b, recordBenchmarkSuites[0])
}

func BenchmarkProtectedRecordRoundTrip(b *testing.B) {
	benchmarkProtectedRecordRoundTrip(b, recordBenchmarkSuites[0])
}

func BenchmarkProtectedRecordRoundTripInPlace(b *testing.B) {
	sender := newBenchmarkRecordCipher(b, TLS_AES_128_GCM_SHA256, 3, 0)
	receiver := newBenchmarkRecordCipher(b, TLS_AES_128_GCM_SHA256, 3, 0)
	payload := bytes.Repeat([]byte{4}, 1200)
	b.ReportAllocs()
	b.SetBytes(int64(len(payload)))

	for b.Loop() {
		wire, err := sender.seal(recordTypeApplicationData, payload)
		if err != nil {
			b.Fatal(err)
		}
		if _, _, _, err = receiver.openInPlace(wire); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProtectedRecordReceiveErrorUnauthenticated(b *testing.B) {
	err := error(&ProtocolError{"protected record authentication failed"})
	b.ReportAllocs()

	for b.Loop() {
		if fatalErr := protectedRecordReceiveError(err); fatalErr != nil {
			b.Fatal(fatalErr)
		}
	}
}

func BenchmarkProtectedRecordSealSuites(b *testing.B) {
	for _, config := range recordBenchmarkSuites {
		b.Run(config.name, func(b *testing.B) {
			benchmarkProtectedRecordSeal(b, config)
		})
	}
}

func BenchmarkProtectedRecordRoundTripSuites(b *testing.B) {
	for _, config := range recordBenchmarkSuites {
		b.Run(config.name, func(b *testing.B) {
			benchmarkProtectedRecordRoundTrip(b, config)
		})
	}
}

func BenchmarkProtectedRecordCID(b *testing.B) {
	config := recordBenchmarkConfig{name: "AES128GCM-CID16", suiteID: TLS_AES_128_GCM_SHA256, cidBytes: 16}
	b.Run("Seal", func(b *testing.B) {
		benchmarkProtectedRecordSeal(b, config)
	})
	b.Run("RoundTrip", func(b *testing.B) {
		benchmarkProtectedRecordRoundTrip(b, config)
	})
}

func BenchmarkBuildProtectedACKRecords(b *testing.B) {
	cases := []struct {
		name    string
		numbers []recordNumber
	}{
		{name: "Single", numbers: []recordNumber{{epoch: 3, sequence: 1}}},
		{name: "Sorted64", numbers: make([]recordNumber, 64)},
		{name: "Reversed64", numbers: make([]recordNumber, 64)},
	}
	for i := range cases[1].numbers {
		cases[1].numbers[i] = recordNumber{epoch: 3, sequence: uint64(i + 1)}
		cases[2].numbers[i] = recordNumber{epoch: 3, sequence: uint64(len(cases[2].numbers) - i)}
	}
	for _, benchmark := range cases {
		b.Run(benchmark.name, func(b *testing.B) {
			cipher := newBenchmarkRecordCipher(b, TLS_AES_128_GCM_SHA256, 7, 0)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, _, err := buildACKRecords(benchmark.numbers, 1200, 0, cipher); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
	b.Run("SingleReuse", func(b *testing.B) {
		cipher := newBenchmarkRecordCipher(b, TLS_AES_128_GCM_SHA256, 7, 0)
		numbers := []recordNumber{{epoch: 3, sequence: 1}}
		var storage [1][]byte
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			if _, _, err := buildACKRecordsInto(storage[:0], numbers, 1200, 0, cipher); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkBuildPlainACKRecords(b *testing.B) {
	cases := []struct {
		name    string
		numbers []recordNumber
	}{
		{name: "Empty"},
		{name: "Single", numbers: []recordNumber{{epoch: 3, sequence: 1}}},
		{name: "Sorted64", numbers: make([]recordNumber, 64)},
	}
	for i := range cases[2].numbers {
		cases[2].numbers[i] = recordNumber{epoch: 3, sequence: uint64(i + 1)}
	}
	for _, benchmark := range cases {
		b.Run(benchmark.name, func(b *testing.B) {
			var storage [1][]byte
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, _, err := buildACKRecordsInto(storage[:0], benchmark.numbers, 1200, 0, nil); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkParseACK(b *testing.B) {
	wire, err := marshalCanonicalACK([]recordNumber{{epoch: 3, sequence: 1}})
	if err != nil {
		b.Fatal(err)
	}
	b.Run("Owned", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := parseACK(wire); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("ReuseSingle", func(b *testing.B) {
		var storage [1]recordNumber
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := parseACKInto(wire, storage[:0]); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkFlightWireWindow(b *testing.B) {
	f := &flight{records: make([]flightRecord, 10)}
	for i := range f.records {
		f.records[i].wire = make([]byte, 1200)
		f.records[i].number = recordNumber{epoch: 2, sequence: uint64(i)}
		f.records[i].sent = true
	}
	var storage [10][]byte
	b.Run("Retransmit", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if records := f.retransmitWire(10, storage[:0]); len(records) != 10 {
				b.Fatal("missing retransmit wire")
			}
		}
	})
}

func BenchmarkFlightFirstRefresh(b *testing.B) {
	f := &flight{records: make([]flightRecord, 10)}
	var sequence uint64
	f.rebuildPending = func(indices flightRecordIndices) error {
		for i := 0; i < indices.count; i++ {
			index := indices.at(i)
			sequence++
			f.records[index].replaceNumber(recordNumber{epoch: 2, sequence: sequence})
		}
		return nil
	}
	b.ReportAllocs()
	for b.Loop() {
		for index := range f.records {
			f.records[index].number = recordNumber{epoch: 2, sequence: uint64(index)}
			f.records[index].priorNumber = recordNumber{}
			f.records[index].earlierNumbers = nil
			f.records[index].hasPrior = false
		}
		if err := f.refreshPending(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFlightInitialHistoryBatch(b *testing.B) {
	records := make([]flightRecord, 10)
	var indices flightRecordIndices
	for i := range records {
		indices.add(i, len(records))
	}
	b.ReportAllocs()
	for b.Loop() {
		for index := range records {
			records[index].number = recordNumber{epoch: 2, sequence: uint64(index)}
			records[index].priorNumber = recordNumber{epoch: 2, sequence: uint64(index + 10)}
			records[index].earlierNumbers = nil
			records[index].hasPrior = true
		}
		reserveInitialRecordHistory(records, indices)
		for index := range records {
			records[index].replaceNumber(recordNumber{epoch: 2, sequence: uint64(index + 20)})
			records[index].replaceNumber(recordNumber{epoch: 2, sequence: uint64(index + 30)})
			records[index].replaceNumber(recordNumber{epoch: 2, sequence: uint64(index + 40)})
		}
	}
}

func BenchmarkBuildProtectedFlight(b *testing.B) {
	cipher := newBenchmarkRecordCipher(b, TLS_AES_128_GCM_SHA256, 2, 0)
	messages := []handshakeMessage{{
		typ: handshakeTypeCertificate, body: make([]byte, 4096),
	}}
	b.ReportAllocs()
	b.SetBytes(int64(len(messages[0].body)))

	for b.Loop() {
		if _, err := buildProtectedFlight(messages, 1200, cipher); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBuildPlainFlight(b *testing.B) {
	messages := []handshakeMessage{{
		typ: handshakeTypeCertificate, body: make([]byte, 4096),
	}}
	b.ReportAllocs()
	b.SetBytes(int64(len(messages[0].body)))

	for b.Loop() {
		if _, _, err := buildPlainFlight(messages, 1200, 0, 0); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCombineFlights(b *testing.B) {
	first := &flight{records: make([]flightRecord, 1)}
	second := &flight{records: make([]flightRecord, 4)}
	b.ReportAllocs()
	for b.Loop() {
		if combined := combineFlights(first, second); len(combined.records) != 5 {
			b.Fatal("combined flight has wrong record count")
		}
	}
}
