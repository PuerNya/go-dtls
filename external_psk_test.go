package dtls13

import (
	"bytes"
	"crypto"
	"crypto/hkdf"
	"crypto/sha256"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"testing"
	"time"
)

func TestImportExternalPSKDerivation(t *testing.T) {
	identity := []byte("device-17")
	context := []byte("client=A;server=B")
	key := bytes.Repeat([]byte{0x5a}, 32)
	psk, err := ImportExternalPSK(identity, key, context, crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	if len(psk.keys) != 2 {
		t.Fatalf("derived %d target keys, want 2", len(psk.keys))
	}
	for i, target := range []struct {
		hash crypto.Hash
		kdf  uint16
	}{{crypto.SHA256, tlsKDFHKDFSHA256}, {crypto.SHA384, tlsKDFHKDFSHA384}} {
		wantIdentity := make([]byte, 0, 8+len(identity)+len(context))
		wantIdentity = binary.BigEndian.AppendUint16(wantIdentity, uint16(len(identity)))
		wantIdentity = append(wantIdentity, identity...)
		wantIdentity = binary.BigEndian.AppendUint16(wantIdentity, uint16(len(context)))
		wantIdentity = append(wantIdentity, context...)
		wantIdentity = binary.BigEndian.AppendUint16(wantIdentity, VersionDTLS13)
		wantIdentity = binary.BigEndian.AppendUint16(wantIdentity, target.kdf)
		if !bytes.Equal(psk.keys[i].wireIdentity, wantIdentity) {
			t.Fatalf("target %v identity=%x, want %x", target.hash, psk.keys[i].wireIdentity, wantIdentity)
		}

		h := sha256.New()
		_, _ = h.Write(wantIdentity)
		extracted, extractErr := hkdf.Extract(sha256.New, key, nil)
		if extractErr != nil {
			t.Fatal(extractErr)
		}
		const label = "dtls13derived psk"
		info := make([]byte, 2+1+len(label)+1+sha256.Size)
		binary.BigEndian.PutUint16(info, uint16(target.hash.Size()))
		info[2] = byte(len(label))
		copy(info[3:], label)
		at := 3 + len(label)
		info[at] = sha256.Size
		copy(info[at+1:], h.Sum(nil))
		wantKey, expandErr := hkdf.Expand(sha256.New, extracted, string(info), target.hash.Size())
		if expandErr != nil {
			t.Fatal(expandErr)
		}
		if !bytes.Equal(psk.keys[i].key, wantKey) {
			t.Fatalf("target %v key=%x, want %x", target.hash, psk.keys[i].key, wantKey)
		}
	}
	if bytes.Equal(psk.keys[0].key, psk.keys[1].key) {
		t.Fatal("SHA-256 and SHA-384 target keys are not separated")
	}
}

func TestExternalPSKValidationAndClone(t *testing.T) {
	if _, err := NewDirectExternalPSK(nil, bytes.Repeat([]byte{1}, 16), crypto.SHA256); err == nil {
		t.Fatal("accepted an empty identity")
	}
	if _, err := NewDirectExternalPSK([]byte("id"), bytes.Repeat([]byte{1}, 15), crypto.SHA256); err == nil {
		t.Fatal("accepted a PSK shorter than 128 bits")
	}
	if _, err := NewDirectExternalPSK([]byte("id"), bytes.Repeat([]byte{1}, 16), crypto.SHA512); err == nil {
		t.Fatal("accepted an unsupported direct PSK hash")
	}
	if _, err := ImportExternalPSK(bytes.Repeat([]byte{1}, 65535), bytes.Repeat([]byte{2}, 16), nil, 0); err == nil {
		t.Fatal("accepted an oversized ImportedIdentity")
	}
	psk, err := NewDirectExternalPSK([]byte("id"), bytes.Repeat([]byte{3}, 16), 0)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = (&Config{ExternalPSKs: []*ExternalPSK{psk, psk}}).normalized(); err == nil {
		t.Fatal("accepted duplicate external PSK wire identities")
	}
	if _, err = (&Config{ExternalPSKs: []*ExternalPSK{psk}, ClientAuth: tls.RequestClientCert}).normalized(); err == nil {
		t.Fatal("accepted external PSK with certificate client authentication")
	}
	config := &Config{ExternalPSKs: []*ExternalPSK{psk}}
	clone := config.Clone()
	clone.ExternalPSKs[0] = nil
	if config.ExternalPSKs[0] != psk {
		t.Fatal("Config.Clone shared the ExternalPSKs slice backing array")
	}
}

func TestImportedAndDirectExternalPSKHandshakes(t *testing.T) {
	direct, err := NewDirectExternalPSK([]byte("direct-device"), bytes.Repeat([]byte{0x31}, 32), crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	directSHA384, err := NewDirectExternalPSK([]byte("direct-sha384-device"), bytes.Repeat([]byte{0x32}, 32), crypto.SHA384)
	if err != nil {
		t.Fatal(err)
	}
	imported, err := ImportExternalPSK([]byte("imported-device"), bytes.Repeat([]byte{0x42}, 32), []byte("client=1;server=2"), crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	unknown, err := NewDirectExternalPSK([]byte("unknown"), bytes.Repeat([]byte{0x53}, 32), crypto.SHA384)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name   string
		client []*ExternalPSK
		server []*ExternalPSK
		suites []uint16
		want   *ExternalPSK
		cipher uint16
	}{
		{"direct", []*ExternalPSK{direct}, []*ExternalPSK{direct}, []uint16{TLS_AES_128_GCM_SHA256}, direct, TLS_AES_128_GCM_SHA256},
		{"direct-sha384-default-suites", []*ExternalPSK{directSHA384}, []*ExternalPSK{directSHA384}, nil, directSHA384, TLS_AES_256_GCM_SHA384},
		{"imported-sha384-later-identity", []*ExternalPSK{unknown, imported}, []*ExternalPSK{imported}, []uint16{TLS_AES_256_GCM_SHA384}, imported, TLS_AES_256_GCM_SHA384},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client, server, clientErr, serverErr := externalPSKHandshake(t,
				&Config{ExternalPSKs: test.client, CipherSuites: test.suites, SessionTicketsDisabled: true},
				&Config{ExternalPSKs: test.server, CipherSuites: test.suites, SessionTicketsDisabled: true})
			if clientErr != nil || serverErr != nil {
				t.Fatalf("external PSK handshake: client=%v server=%v", clientErr, serverErr)
			}
			for side, state := range map[string]ConnectionState{"client": client.ConnectionState(), "server": server.ConnectionState()} {
				if state.DidResume || len(state.PeerCertificates) != 0 || state.CipherSuite != test.cipher || !bytes.Equal(state.ExternalPSKIdentity(), test.want.identity) || !bytes.Equal(state.ExternalPSKContext(), test.want.context) {
					t.Fatalf("%s state=%+v", side, state)
				}
			}
			state := client.ConnectionState()
			identity := state.ExternalPSKIdentity()
			identity[0] ^= 0xff
			if !bytes.Equal(client.ConnectionState().ExternalPSKIdentity(), test.want.identity) {
				t.Fatal("ConnectionState returned mutable external identity storage")
			}
			if _, err = client.WriteDatagram([]byte("external-psk")); err != nil {
				t.Fatal(err)
			}
			buffer := make([]byte, 32)
			n, _, readErr := server.ReadDatagram(buffer)
			if readErr != nil || string(buffer[:n]) != "external-psk" {
				t.Fatalf("application data=%q err=%v", buffer[:n], readErr)
			}
			state = client.ConnectionState()
			if err = client.Close(); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(state.ExternalPSKIdentity(), test.want.identity) || !bytes.Equal(client.ConnectionState().ExternalPSKIdentity(), test.want.identity) {
				t.Fatal("closing the connection discarded external PSK authentication state")
			}
		})
	}
}

func TestExternalPSKHelloRetryRequest(t *testing.T) {
	psk, err := ImportExternalPSK([]byte("hrr-device"), bytes.Repeat([]byte{0x64}, 32), []byte("roles"), 0)
	if err != nil {
		t.Fatal(err)
	}
	client, _, clientErr, serverErr := externalPSKHandshake(t,
		&Config{ExternalPSKs: []*ExternalPSK{psk}, CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}, SessionTicketsDisabled: true},
		&Config{ExternalPSKs: []*ExternalPSK{psk}, CurvePreferences: []tls.CurveID{tls.CurveP256}, SessionTicketsDisabled: true})
	if clientErr != nil || serverErr != nil {
		t.Fatalf("external PSK HRR: client=%v server=%v", clientErr, serverErr)
	}
	if !bytes.Equal(client.ConnectionState().ExternalPSKIdentity(), psk.identity) {
		t.Fatal("HRR lost external PSK authentication")
	}
}

func TestExternalPSKHelloRetryRequestIdentityFiltering(t *testing.T) {
	psk, err := ImportExternalPSK([]byte("filter-device"), bytes.Repeat([]byte{0x6d}, 32), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	config, err := (&Config{ExternalPSKs: []*ExternalPSK{psk}}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	initial := &clientHello{pskIdentities: []pskIdentityEntry{{identity: psk.keys[0].wireIdentity}, {identity: psk.keys[1].wireIdentity}}}
	compatibleRemoved := &clientHello{pskIdentities: []pskIdentityEntry{{identity: psk.keys[1].wireIdentity}}}
	allowed, err := removedPSKsAreIncompatible(config, initial, compatibleRemoved, suite)
	if err != nil {
		t.Fatal(err)
	}
	if allowed {
		t.Fatal("HRR allowed removal of a SHA-256 external PSK for a SHA-256 suite")
	}
	incompatibleRemoved := &clientHello{pskIdentities: []pskIdentityEntry{{identity: psk.keys[0].wireIdentity}}}
	allowed, err = removedPSKsAreIncompatible(config, initial, incompatibleRemoved, suite)
	if err != nil || !allowed {
		t.Fatalf("HRR rejected removal of only an incompatible PSK: allowed=%v err=%v", allowed, err)
	}
}

func TestExternalPSKBinderLabelsAreSeparated(t *testing.T) {
	psk, err := ImportExternalPSK([]byte("binder-device"), bytes.Repeat([]byte{0x75}, 32), nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	suite, _ := cipherSuiteForID(TLS_AES_128_GCM_SHA256)
	key := &psk.keys[0]
	hello := &clientHello{
		cipherSuites: []uint16{suite.id}, supportedGroups: []tls.CurveID{tls.X25519},
		keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}}, signatureSchemes: defaultSignatureSchemes(),
	}
	body, err := marshalClientHelloWithPSKOffers(hello, []clientPSKOffer{{identity: key.wireIdentity, psk: key.key, suite: suite, binderLabel: labelImportedBinder}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = verifyClientHelloPSKBinderAtWithLabel(body, suite, key.key, labelImportedBinder, nil, nil, 0); err != nil {
		t.Fatal(err)
	}
	if err = verifyClientHelloPSKBinderAtWithLabel(body, suite, key.key, labelExternalBinder, nil, nil, 0); err == nil {
		t.Fatal("imported binder was accepted as a direct external binder")
	} else if description, ok := protocolAlert(err); !ok || description != alertDecryptError {
		t.Fatalf("binder mismatch alert=%d ok=%v err=%v", description, ok, err)
	}
}

func TestExternalPSKCertificateFallbackAndFailures(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	clientPSK, _ := NewDirectExternalPSK([]byte("client-id"), bytes.Repeat([]byte{1}, 32), crypto.SHA256)
	serverPSK, _ := NewDirectExternalPSK([]byte("server-id"), bytes.Repeat([]byte{2}, 32), crypto.SHA256)
	client, _, clientErr, serverErr := externalPSKHandshake(t,
		&Config{ExternalPSKs: []*ExternalPSK{clientPSK}, RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true},
		&Config{ExternalPSKs: []*ExternalPSK{serverPSK}, Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true})
	if clientErr != nil || serverErr != nil {
		t.Fatalf("certificate fallback: client=%v server=%v", clientErr, serverErr)
	}
	if len(client.ConnectionState().ExternalPSKIdentity()) != 0 || len(client.ConnectionState().PeerCertificates) == 0 {
		t.Fatal("unknown external identity did not fall back to certificate authentication")
	}

	wrongKey, _ := NewDirectExternalPSK([]byte("same-id"), bytes.Repeat([]byte{3}, 32), crypto.SHA256)
	rightKey, _ := NewDirectExternalPSK([]byte("same-id"), bytes.Repeat([]byte{4}, 32), crypto.SHA256)
	_, _, badClientErr, badServerErr := externalPSKHandshake(t,
		&Config{ExternalPSKs: []*ExternalPSK{wrongKey}, SessionTicketsDisabled: true, HandshakeTimeout: time.Second},
		&Config{ExternalPSKs: []*ExternalPSK{rightKey}, SessionTicketsDisabled: true, HandshakeTimeout: time.Second})
	if badClientErr == nil || badServerErr == nil {
		t.Fatalf("wrong external key succeeded: client=%v server=%v", badClientErr, badServerErr)
	}

	for _, test := range []struct {
		name          string
		clientID      []byte
		clientContext []byte
	}{{"identity", []byte("other-device"), []byte("roles")}, {"context", []byte("imported-device"), []byte("other-roles")}} {
		t.Run("wrong-"+test.name, func(t *testing.T) {
			clientImported, importErr := ImportExternalPSK(test.clientID, bytes.Repeat([]byte{5}, 32), test.clientContext, 0)
			if importErr != nil {
				t.Fatal(importErr)
			}
			serverImported, importErr := ImportExternalPSK([]byte("imported-device"), bytes.Repeat([]byte{5}, 32), []byte("roles"), 0)
			if importErr != nil {
				t.Fatal(importErr)
			}
			_, _, mismatchClientErr, mismatchServerErr := externalPSKHandshake(t,
				&Config{ExternalPSKs: []*ExternalPSK{clientImported}, SessionTicketsDisabled: true, HandshakeTimeout: 200 * time.Millisecond, FlightInterval: 5 * time.Millisecond},
				&Config{ExternalPSKs: []*ExternalPSK{serverImported}, SessionTicketsDisabled: true, HandshakeTimeout: 200 * time.Millisecond, FlightInterval: 5 * time.Millisecond})
			if mismatchClientErr == nil || mismatchServerErr == nil {
				t.Fatalf("wrong %s succeeded: client=%v server=%v", test.name, mismatchClientErr, mismatchServerErr)
			}
		})
	}
}

func TestExternalPSKTicketResumptionAndRevocation(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	key := bytes.Repeat([]byte{0x86}, 32)
	psk, _ := ImportExternalPSK([]byte("ticket-device"), key, []byte("client=7;server=9"), 0)
	cache := NewLRUClientSessionCache(2)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x97}, 32))
	clientConfig := &Config{ExternalPSKs: []*ExternalPSK{psk}, RootCAs: roots, ServerName: "server.test", ClientSessionCache: cache}
	serverConfig := &Config{ExternalPSKs: []*ExternalPSK{psk}, Certificates: []tls.Certificate{certificate}, SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour}

	firstClient, firstServer, clientErr, serverErr := externalPSKHandshake(t, clientConfig, serverConfig)
	if clientErr != nil || serverErr != nil {
		t.Fatalf("initial external handshake: client=%v server=%v", clientErr, serverErr)
	}
	waitForExternalPSKTicket(t, cache)
	_ = firstClient.Close()
	_ = firstServer.Close()

	resumedClient, resumedServer, clientErr, serverErr := externalPSKHandshake(t, clientConfig, serverConfig)
	if clientErr != nil || serverErr != nil {
		t.Fatalf("external-origin resumption: client=%v server=%v", clientErr, serverErr)
	}
	if !resumedClient.ConnectionState().DidResume || !resumedServer.ConnectionState().DidResume || !bytes.Equal(resumedClient.ConnectionState().ExternalPSKIdentity(), psk.identity) {
		t.Fatal("ticket resumption did not preserve external PSK authentication")
	}
	waitForExternalPSKTicket(t, cache)
	_ = resumedClient.Close()
	_ = resumedServer.Close()

	serverWithoutPSK := serverConfig.Clone()
	serverWithoutPSK.ExternalPSKs = nil
	fallbackClient, _, clientErr, serverErr := externalPSKHandshake(t, clientConfig, serverWithoutPSK)
	if clientErr != nil || serverErr != nil {
		t.Fatalf("revoked external PSK fallback: client=%v server=%v", clientErr, serverErr)
	}
	if fallbackClient.ConnectionState().DidResume || len(fallbackClient.ConnectionState().ExternalPSKIdentity()) != 0 || len(fallbackClient.ConnectionState().PeerCertificates) == 0 {
		t.Fatal("server accepted a ticket after its external PSK origin was removed")
	}
}

func TestExternalPSKChangeInvalidatesClientTicket(t *testing.T) {
	identity, context := []byte("rotated-device"), []byte("roles")
	oldPSK, _ := ImportExternalPSK(identity, bytes.Repeat([]byte{0xa8}, 32), context, 0)
	cache := NewLRUClientSessionCache(1)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0xb9}, 32))
	clientConfig := &Config{ExternalPSKs: []*ExternalPSK{oldPSK}, ClientSessionCache: cache}
	serverConfig := &Config{ExternalPSKs: []*ExternalPSK{oldPSK}, SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour}
	client, server, clientErr, serverErr := externalPSKHandshake(t, clientConfig, serverConfig)
	if clientErr != nil || serverErr != nil {
		t.Fatalf("initial handshake: client=%v server=%v", clientErr, serverErr)
	}
	waitForExternalPSKTicket(t, cache)
	_ = client.Close()
	_ = server.Close()

	newPSK, _ := ImportExternalPSK(identity, bytes.Repeat([]byte{0xca}, 32), context, 0)
	clientConfig.ExternalPSKs = []*ExternalPSK{newPSK}
	serverConfig.ExternalPSKs = []*ExternalPSK{newPSK}
	client, server, clientErr, serverErr = externalPSKHandshake(t, clientConfig, serverConfig)
	if clientErr != nil || serverErr != nil {
		t.Fatalf("rotated external PSK handshake: client=%v server=%v", clientErr, serverErr)
	}
	if client.ConnectionState().DidResume || server.ConnectionState().DidResume {
		t.Fatal("client offered a ticket derived from the previous external PSK")
	}
}

func TestExternalPSKEarlyDataUnavailable(t *testing.T) {
	psk, _ := NewDirectExternalPSK([]byte("no-early-data"), bytes.Repeat([]byte{0xdb}, 32), crypto.SHA256)
	left, right := memoryDatagramPair()
	client := Client(left, &Config{ExternalPSKs: []*ExternalPSK{psk}, SessionTicketsDisabled: true})
	server := Server(right, &Config{ExternalPSKs: []*ExternalPSK{psk}, SessionTicketsDisabled: true})
	defer client.Close()
	defer server.Close()
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	if _, err := client.WriteEarlyData([]byte("must-not-send")); !errors.Is(err, ErrEarlyDataUnavailable) {
		t.Fatalf("WriteEarlyData returned %v", err)
	}
	if err := <-serverDone; err != nil {
		t.Fatal(err)
	}
	client.earlyMu.Lock()
	pending := len(client.earlyPending)
	client.earlyMu.Unlock()
	if pending != 0 {
		t.Fatal("unavailable external early data remained buffered")
	}
}

func TestExternalPSKWeakNetwork(t *testing.T) {
	psk, _ := ImportExternalPSK([]byte("weak-network-device"), bytes.Repeat([]byte{0xec}, 32), []byte("roles"), 0)
	left, right := memoryDatagramPair()
	clientWire := &weakNetworkConn{Conn: left, enabled: true}
	serverWire := &weakNetworkConn{Conn: right, enabled: true}
	config := &Config{ExternalPSKs: []*ExternalPSK{psk}, SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second, FlightInterval: 5 * time.Millisecond}
	client := Client(clientWire, config)
	server := Server(serverWire, config)
	defer client.Close()
	defer server.Close()
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	serverErr := <-serverDone
	if clientErr != nil || serverErr != nil {
		t.Fatalf("weak-network external PSK: client=%v server=%v", clientErr, serverErr)
	}
}

func externalPSKHandshake(t *testing.T, clientConfig, serverConfig *Config) (*Conn, *Conn, error, error) {
	t.Helper()
	if clientConfig.HandshakeTimeout == 0 {
		clientConfig = clientConfig.Clone()
		clientConfig.HandshakeTimeout = 2 * time.Second
		clientConfig.FlightInterval = 5 * time.Millisecond
	}
	if serverConfig.HandshakeTimeout == 0 {
		serverConfig = serverConfig.Clone()
		serverConfig.HandshakeTimeout = 2 * time.Second
		serverConfig.FlightInterval = 5 * time.Millisecond
	}
	left, right := memoryDatagramPair()
	client := Client(left, clientConfig)
	server := Server(right, serverConfig)
	t.Cleanup(func() {
		_ = client.Close()
		_ = server.Close()
	})
	serverDone := make(chan error, 1)
	go func() { serverDone <- server.Handshake() }()
	clientErr := client.Handshake()
	return client, server, clientErr, <-serverDone
}

func waitForExternalPSKTicket(t *testing.T, cache ClientSessionCache) *ClientSessionState {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if state, ok := cache.Get("server.test"); ok {
			return state
		}
		if state, ok := cache.Get("remote"); ok {
			return state
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("client did not receive an external-origin session ticket")
	return nil
}
