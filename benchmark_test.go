package dtls13

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"fmt"
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

func BenchmarkTranscriptClone(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		transcript := newTranscriptHash(suite.hash.New())
		if err = transcript.add(handshakeTypeCertificate, 1, bytes.Repeat([]byte{0x5a}, 64<<10)); err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				clone := transcript.clone()
				if len(clone.sum()) != suite.hash.Size() {
					b.Fatal("invalid cloned transcript hash")
				}
			}
		})
	}
}

func BenchmarkTranscriptSum(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		transcript := newTranscriptHash(suite.hash.New())
		if err = transcript.add(handshakeTypeCertificate, 1, bytes.Repeat([]byte{0x5a}, 64<<10)); err != nil {
			b.Fatal(err)
		}
		prefix := fmt.Sprintf("%04x", suiteID)
		b.Run(prefix+"/Owned", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if len(transcript.sum()) != suite.hash.Size() {
					b.Fatal("invalid transcript hash")
				}
			}
		})
		b.Run(prefix+"/Reuse", func(b *testing.B) {
			var scratch [maxSupportedHashSize]byte
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if len(transcript.sumInto(scratch[:0])) != suite.hash.Size() {
					b.Fatal("invalid transcript hash")
				}
			}
		})
	}
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

func BenchmarkParseHandshakeFragment(b *testing.B) {
	body := bytes.Repeat([]byte{0x5a}, 1200)
	wire, err := marshalHandshakeFragment(handshakeFragment{typ: handshakeTypeCertificate, messageSequence: 1, length: uint32(len(body)), body: body})
	if err != nil {
		b.Fatal(err)
	}
	b.Run("View", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := parseHandshakeFragmentsView(wire); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("ReuseSingle", func(b *testing.B) {
		var storage [1]handshakeFragment
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := parseHandshakeFragmentsViewInto(wire, storage[:0]); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkParsePlainRecord(b *testing.B) {
	payload := bytes.Repeat([]byte{0x5a}, 1200)
	wire, err := marshalPlainRecord(record{typ: recordTypeHandshake, sequence: 1, payload: payload})
	if err != nil {
		b.Fatal(err)
	}
	b.Run("View", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := parsePlainRecordsView(wire); err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("ReuseSingle", func(b *testing.B) {
		var storage [1]record
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if _, err := parsePlainRecordsViewInto(wire, storage[:0]); err != nil {
				b.Fatal(err)
			}
		}
	})
}

func BenchmarkFlightPendingIndices(b *testing.B) {
	f := &flight{records: make([]flightRecord, 10)}
	b.Run("Allocated", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if indices := f.pendingIndices(nil); len(indices) != len(f.records) {
				b.Fatal("missing pending record")
			}
		}
	})
	b.Run("ReuseWindow", func(b *testing.B) {
		var storage [10]int
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if indices := f.pendingIndices(storage[:0]); len(indices) != len(f.records) {
				b.Fatal("missing pending record")
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
	b.Run("Pending", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if records := f.pendingWire(storage[:0]); len(records) != 10 {
				b.Fatal("missing pending wire")
			}
		}
	})
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

func BenchmarkParseExtensions(b *testing.B) {
	items := map[uint16][]byte{
		extSupportedVersions:   {byte(VersionDTLS13 >> 8), byte(VersionDTLS13 & 0xff)},
		extSupportedGroups:     make([]byte, 16),
		extSignatureAlgorithms: make([]byte, 20),
		extKeyShare:            make([]byte, 40),
		extServerName:          make([]byte, 24),
		extALPN:                make([]byte, 12),
	}
	order := []uint16{extSupportedVersions, extSupportedGroups, extSignatureAlgorithms, extKeyShare, extServerName, extALPN}
	wire, err := marshalExtensions(items, order)
	if err != nil {
		b.Fatal(err)
	}
	for _, benchmark := range []struct {
		name  string
		parse func([]byte) (map[uint16][]byte, error)
	}{{"Owned", parseExtensions}, {"View", parseExtensionsView}} {
		b.Run(benchmark.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if _, err := benchmark.parse(wire); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
	b.Run("OrderedView", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var storage [8]orderedExtension
			items, parseErr := parseOrderedExtensionsView(wire, storage[:0])
			if parseErr != nil || len(items) != len(order) {
				b.Fatal("failed to parse ordered extensions")
			}
		}
	})
}

func BenchmarkParseKeyShares(b *testing.B) {
	for _, count := range []int{1, 4, 9} {
		shares := make([]keyShareEntry, count)
		for i := range shares {
			shares[i] = keyShareEntry{group: tls.CurveID(0x100 + i), data: bytes.Repeat([]byte{byte(i + 1)}, 32)}
		}
		wire, err := marshalKeyShares(shares, true)
		if err != nil {
			b.Fatal(err)
		}
		for _, benchmark := range []struct {
			name  string
			parse func([]byte, bool) ([]keyShareEntry, error)
		}{{"Owned", parseKeyShares}, {"View", parseKeySharesView}} {
			b.Run(fmt.Sprintf("Client%d/%s", count, benchmark.name), func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					parsed, parseErr := benchmark.parse(wire, true)
					if parseErr != nil || len(parsed) != count {
						b.Fatal("failed to parse key shares")
					}
				}
			})
		}
		b.Run(fmt.Sprintf("Client%d/ViewInto", count), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				var storage [8]keyShareEntry
				parsed, parseErr := parseKeySharesViewInto(wire, true, storage[:0])
				if parseErr != nil || len(parsed) != count {
					b.Fatal("failed to parse key shares into caller storage")
				}
			}
		})
	}
}

func BenchmarkDeriveTrafficKeys(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				keys := deriveTrafficKeys(suite, secret)
				if len(keys.key) != suite.keyLen || len(keys.iv) != suite.ivLen || len(keys.sn) != suite.keyLen {
					b.Fatal("invalid traffic keys")
				}
			}
		})
	}
}

func BenchmarkDeriveTrafficKeysInto(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		key := make([]byte, suite.keyLen)
		iv := make([]byte, suite.ivLen)
		sn := make([]byte, suite.keyLen)
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				deriveTrafficKeysInto(suite, secret, key, iv, sn)
			}
		})
	}
}

func BenchmarkNewRecordCipher(b *testing.B) {
	for _, config := range recordBenchmarkSuites {
		suite, err := cipherSuiteForID(config.suiteID)
		if err != nil {
			b.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		b.Run(config.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				record, err := newRecordCipher(suite, secret, 3, 64)
				if err != nil || record.aead == nil || record.sequenceMask == nil {
					b.Fatal("failed to construct record cipher")
				}
			}
		})
	}
}

func BenchmarkKeyScheduleDerivation(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		shared := bytes.Repeat([]byte{0x33}, suite.hash.Size())
		transcript := bytes.Repeat([]byte{0x44}, suite.hash.Size())
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				schedule := newKeySchedule(suite, nil)
				if err := schedule.deriveHandshake(shared, transcript); err != nil {
					b.Fatal(err)
				}
				if err := schedule.deriveApplication(transcript); err != nil {
					b.Fatal(err)
				}
				if err := schedule.deriveResumption(transcript); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkEmptyTranscriptHash(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if len(emptyTranscriptHash(suite)) != suite.hash.Size() {
					b.Fatal("invalid empty transcript hash")
				}
			}
		})
	}
}

func BenchmarkFinishedVerifyData(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		transcript := bytes.Repeat([]byte{0xa5}, suite.hash.Size())
		schedule := &keySchedule{suite: suite}
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if verify := schedule.finishedVerifyData(secret, transcript); len(verify) != suite.hash.Size() {
					b.Fatal("invalid Finished verify_data")
				}
			}
		})
	}
}

func BenchmarkCalculatePSKBinder(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		psk := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		transcript := bytes.Repeat([]byte{0xa5}, suite.hash.Size())
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if binder := calculatePSKBinder(suite, psk, transcript); len(binder) != suite.hash.Size() {
					b.Fatal("invalid PSK binder")
				}
			}
		})
	}
}

func BenchmarkKeyScheduleSideDerivations(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		context := bytes.Repeat([]byte{0xa5}, suite.hash.Size())
		prefix := fmt.Sprintf("%04x", suiteID)

		b.Run(prefix+"/EarlyTraffic", func(b *testing.B) {
			schedule := newKeySchedule(suite, secret)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if got := schedule.earlyTrafficSecret(context); len(got) != suite.hash.Size() {
					b.Fatal("invalid early traffic secret")
				}
			}
		})
		b.Run(prefix+"/TrafficUpdate", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if got := scheduleTrafficUpdate(suite, secret); len(got) != suite.hash.Size() {
					b.Fatal("invalid traffic update secret")
				}
			}
		})
		b.Run(prefix+"/ResumptionPSK", func(b *testing.B) {
			schedule := &keySchedule{suite: suite, resumptionMasterSecret: secret}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if got, err := schedule.resumptionPSK(context); err != nil || len(got) != suite.hash.Size() {
					b.Fatal("invalid resumption PSK")
				}
			}
		})
		b.Run(prefix+"/Exporter", func(b *testing.B) {
			exporter := newExporter(suite, secret)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if got, err := exporter.export("benchmark", context, suite.hash.Size()); err != nil || len(got) != suite.hash.Size() {
					b.Fatal("invalid exporter output")
				}
			}
		})
		b.Run(prefix+"/ExporterZero", func(b *testing.B) {
			exporter := newExporter(suite, secret)
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if got, err := exporter.export("benchmark", context, 0); err != nil || len(got) != 0 {
					b.Fatal("invalid zero-length exporter output")
				}
			}
		})
	}
}

func BenchmarkReceivingTrafficKeyUpdate(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		newReceiver := func() *receivingTraffic {
			receiver, newErr := newReceivingTraffic(suite, secret, 3)
			if newErr != nil {
				b.Fatal(newErr)
			}
			return receiver
		}
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			receiver := newReceiver()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if i > 0 && i%60000 == 0 {
					b.StopTimer()
					receiver = newReceiver()
					b.StartTimer()
				}
				if _, updated, updateErr := receiver.processKeyUpdate(uint16(i%60000), []byte{0}); updateErr != nil || !updated {
					b.Fatalf("key update %d: updated=%v err=%v", i, updated, updateErr)
				}
			}
		})
	}
}

func BenchmarkSendingTrafficKeyUpdate(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		secret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		newSender := func() *sendingTraffic {
			sender, newErr := newSendingTraffic(suite, secret, 3, 0)
			if newErr != nil {
				b.Fatal(newErr)
			}
			wire, number, newErr := sender.beginKeyUpdate(false)
			if newErr != nil || len(wire) == 0 || !sender.processACK([]recordNumber{number}) {
				b.Fatalf("warm-up KeyUpdate: wire=%d err=%v", len(wire), newErr)
			}
			return sender
		}
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			sender := newSender()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if i > 0 && i%60000 == 0 {
					b.StopTimer()
					sender = newSender()
					b.StartTimer()
				}
				wire, number, updateErr := sender.beginKeyUpdate(false)
				if updateErr != nil || len(wire) == 0 || !sender.processACK([]recordNumber{number}) {
					b.Fatalf("KeyUpdate %d: wire=%d err=%v", i, len(wire), updateErr)
				}
			}
		})
	}
}

func BenchmarkInstallApplicationKeys(b *testing.B) {
	for _, suiteID := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384} {
		suite, err := cipherSuiteForID(suiteID)
		if err != nil {
			b.Fatal(err)
		}
		clientSecret := bytes.Repeat([]byte{0x5a}, suite.hash.Size())
		serverSecret := bytes.Repeat([]byte{0xa5}, suite.hash.Size())
		b.Run(fmt.Sprintf("%04x", suiteID), func(b *testing.B) {
			conn := &Conn{config: &Config{ReplayWindow: 64}, isClient: true}
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				if installErr := conn.installApplicationKeysAt(suite, clientSecret, serverSecret, uint16(i)); installErr != nil {
					b.Fatal(installErr)
				}
				conn.sendingTraffic.clearSecrets()
				conn.receivingTraffic.clearSecrets()
			}
		})
	}
}

func BenchmarkMarshalExtensions(b *testing.B) {
	items := map[uint16][]byte{
		extSupportedVersions:   {2, 0xfe, 0xfc},
		extSupportedGroups:     {0, 4, 0, 29, 0, 23},
		extSignatureAlgorithms: bytes.Repeat([]byte{0x5a}, 20),
		extKeyShare:            bytes.Repeat([]byte{0xa5}, 72),
		extConnectionID:        {4, 1, 2, 3, 4},
	}
	order := []uint16{extSupportedVersions, extSupportedGroups, extSignatureAlgorithms, extKeyShare, extConnectionID}
	b.ReportAllocs()
	for b.Loop() {
		if wire, err := marshalExtensions(items, order); err != nil || len(wire) == 0 {
			b.Fatal("failed to marshal extensions")
		}
	}
}

func BenchmarkMarshalHandshakeMessages(b *testing.B) {
	client := &clientHello{
		cipherSuites:     defaultCipherSuites(),
		supportedGroups:  []tls.CurveID{tls.X25519, tls.CurveP256},
		keyShares:        []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{0x11}, 32)}},
		signatureSchemes: defaultSignatureSchemes(),
		serverName:       "server.test",
		alpn:             []string{"coap"},
		cookie:           bytes.Repeat([]byte{0x22}, 64),
	}
	resumptionClient := *client
	resumptionClient.pskIdentities = []pskIdentityEntry{{
		identity:      bytes.Repeat([]byte{0xaa}, 96),
		obfuscatedAge: 0x10203040,
	}}
	resumptionClient.pskBinders = [][]byte{bytes.Repeat([]byte{0xbb}, sha256.Size)}
	server := &serverHello{
		cipherSuite:     TLS_AES_128_GCM_SHA256,
		keyShare:        keyShareEntry{group: tls.X25519, data: bytes.Repeat([]byte{0x33}, 32)},
		connectionID:    bytes.Repeat([]byte{0x44}, 8),
		hasConnectionID: true,
	}
	retry := &helloRetryRequest{
		cipherSuite:   TLS_AES_128_GCM_SHA256,
		selectedGroup: tls.X25519,
		cookie:        bytes.Repeat([]byte{0x55}, 64),
	}
	certificateVerify := &certificateVerifyMessage{
		algorithm: tls.Ed25519,
		signature: bytes.Repeat([]byte{0x66}, 64),
	}
	certificate := &certificateMessage{certificates: []certificateEntry{{
		data:       bytes.Repeat([]byte{0x77}, 1024),
		extensions: map[uint16][]byte{0xff01: {0x01, 0x02}},
	}}}
	ticket := &newSessionTicketMessage{
		lifetime:     3600,
		ageAdd:       0x12345678,
		nonce:        bytes.Repeat([]byte{0x88}, 8),
		ticket:       bytes.Repeat([]byte{0x99}, 64),
		maxEarlyData: 4096,
	}
	ticketState := &sessionTicketState{
		createdAt:    123456789,
		lifetime:     3600,
		suite:        TLS_AES_128_GCM_SHA256,
		psk:          bytes.Repeat([]byte{0xcc}, sha256.Size),
		serverName:   "server.test",
		protocol:     "coap",
		ageAdd:       0x50607080,
		maxEarlyData: 4096,
	}
	connectionIDUpdate := &newConnectionIDMessage{
		connectionIDs: [][]byte{bytes.Repeat([]byte{0xdd}, 8), bytes.Repeat([]byte{0xee}, 16)},
		usage:         connectionIDSpare,
	}
	benchmarks := []struct {
		name    string
		marshal func() ([]byte, error)
	}{
		{name: "ClientHello", marshal: client.marshal},
		{name: "ResumptionClientHello", marshal: resumptionClient.marshal},
		{name: "ServerHello", marshal: server.marshal},
		{name: "HelloRetryRequest", marshal: retry.marshal},
		{name: "CertificateVerify", marshal: certificateVerify.marshal},
		{name: "Certificate", marshal: certificate.marshal},
		{name: "NewSessionTicket", marshal: ticket.marshal},
		{name: "SessionTicketState", marshal: ticketState.marshal},
		{name: "NewConnectionID", marshal: connectionIDUpdate.marshal},
	}
	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				wire, err := benchmark.marshal()
				if err != nil || len(wire) == 0 {
					b.Fatalf("marshal returned %d bytes, %v", len(wire), err)
				}
			}
		})
	}
}

func BenchmarkHandshakeReassembly(b *testing.B) {
	const messageSize = 64 << 10
	body := bytes.Repeat([]byte{0x5a}, messageSize)
	const fragmentSize = 1200
	b.ReportAllocs()
	b.SetBytes(messageSize)

	for b.Loop() {
		r := newReassemblerWithLimit(messageSize)
		for offset := 0; offset < len(body); offset += fragmentSize {
			end := offset + fragmentSize
			if end > len(body) {
				end = len(body)
			}
			fragment := handshakeFragment{typ: handshakeTypeCertificate, messageSequence: 1, length: messageSize, offset: uint32(offset), body: body[offset:end]}
			if _, _, err := r.add(fragment); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkHandshakeReassemblySingleFragment(b *testing.B) {
	body := bytes.Repeat([]byte{0x5a}, 1200)
	fragment := handshakeFragment{typ: handshakeTypeCertificate, messageSequence: 1, length: uint32(len(body)), body: body}
	b.ReportAllocs()
	b.SetBytes(int64(len(body)))

	for b.Loop() {
		r := newReassemblerWithLimit(len(body))
		got, complete, err := r.add(fragment)
		if err != nil || !complete || len(got) != len(body) {
			b.Fatal("failed to reassemble complete fragment")
		}
	}
}

func BenchmarkHandshakeInboxSequentialSingleFragment(b *testing.B) {
	body := bytes.Repeat([]byte{0x5a}, 1200)
	fragment := handshakeFragment{typ: handshakeTypeCertificate, length: uint32(len(body)), body: body}
	b.ReportAllocs()
	b.SetBytes(int64(len(body)))

	for b.Loop() {
		inbox := newHandshakeInbox(0, len(body), 1, len(body))
		messages, err := inbox.add(fragment)
		if err != nil || len(messages) != 1 || len(messages[0].body) != len(body) {
			b.Fatal("failed to deliver sequential fragment")
		}
	}
}

func BenchmarkHandshakeInboxSequentialSingleFragmentReuse(b *testing.B) {
	body := bytes.Repeat([]byte{0x5a}, 1200)
	fragment := handshakeFragment{typ: handshakeTypeCertificate, length: uint32(len(body)), body: body}
	var storage [1]completedHandshake
	b.ReportAllocs()
	b.SetBytes(int64(len(body)))

	for b.Loop() {
		inbox := newHandshakeInbox(0, len(body), 1, len(body))
		messages, err := inbox.addInto(storage[:0], fragment)
		if err != nil || len(messages) != 1 || &messages[0] != &storage[0] || len(messages[0].body) != len(body) {
			b.Fatal("failed to reuse sequential delivery destination")
		}
	}
}

func BenchmarkHandshakeInboxSequentialSingleFragmentBatch(b *testing.B) {
	body := bytes.Repeat([]byte{0x5a}, 1200)
	fragment := handshakeFragment{typ: handshakeTypeCertificate, length: uint32(len(body)), body: body}
	b.ReportAllocs()
	b.SetBytes(int64(len(body)))

	for b.Loop() {
		inbox := newHandshakeInbox(0, len(body), 1, len(body))
		var messages completedHandshakeBatch
		if err := inbox.addBatch(&messages, fragment); err != nil || messages.len() != 1 || len(messages.at(0).body) != len(body) {
			b.Fatal("failed to deliver sequential fragment into batch")
		}
	}
}
