package dtls13

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

func rawPublicKeyDER(t testing.TB, signer crypto.Signer) []byte {
	t.Helper()
	der, err := x509.MarshalPKIXPublicKey(signer.Public())
	if err != nil {
		t.Fatal(err)
	}
	return der
}

func pinRawPublicKey(want []byte) func([]byte) error {
	return func(got []byte) error {
		if !bytes.Equal(got, want) {
			return errors.New("untrusted raw public key")
		}
		return nil
	}
}

func rawPublicKeyConfigs(t testing.TB) (*Config, *Config) {
	t.Helper()
	serverKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{1}, ed25519.SeedSize))
	clientKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{2}, ed25519.SeedSize))
	client := &Config{
		ServerName: "server.test", ServerCertificateTypes: []CertificateType{CertificateTypeRawPublicKey},
		ClientCertificateTypes: []CertificateType{CertificateTypeRawPublicKey}, RawPublicKeySigner: clientKey,
		VerifyPeerRawPublicKey: pinRawPublicKey(rawPublicKeyDER(t, serverKey)),
		HandshakeTimeout:       2 * time.Second, FlightInterval: 5 * time.Millisecond, MTU: 256,
	}
	server := &Config{
		ServerCertificateTypes: []CertificateType{CertificateTypeRawPublicKey},
		ClientCertificateTypes: []CertificateType{CertificateTypeRawPublicKey}, RawPublicKeySigner: serverKey,
		VerifyPeerRawPublicKey: pinRawPublicKey(rawPublicKeyDER(t, clientKey)), ClientAuth: tls.RequireAndVerifyClientCert,
		HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond, MTU: 256,
	}
	return client, server
}

func requireRawPublicKey(t testing.TB, conn *Conn, want []byte) {
	t.Helper()
	state := conn.ConnectionState()
	if !bytes.Equal(state.PeerRawPublicKey, want) || len(state.PeerCertificates) != 0 || len(state.VerifiedChains) != 0 || len(state.OCSPResponse) != 0 {
		t.Fatalf("unexpected RPK identity: %#v", state)
	}
	state.PeerRawPublicKey[0] ^= 0xff
	if !bytes.Equal(conn.ConnectionState().PeerRawPublicKey, want) {
		t.Fatal("ConnectionState exposes mutable RPK storage")
	}
}

func TestRawPublicKeyHandshake(t *testing.T) {
	for _, mode := range []string{"mutual", "server-only", "empty-client", "raw-client", "raw-server", "x509-fallback", "x509-without-rpk-verifier", "rpk-without-x509-credential"} {
		t.Run(mode, func(t *testing.T) {
			clientConfig, serverConfig := rawPublicKeyConfigs(t)
			serverDER, clientDER := rawPublicKeyDER(t, serverConfig.RawPublicKeySigner), rawPublicKeyDER(t, clientConfig.RawPublicKeySigner)
			serverType, clientType := CertificateTypeRawPublicKey, CertificateTypeRawPublicKey
			switch mode {
			case "rpk-without-x509-credential":
				clientConfig.ServerCertificateTypes = []CertificateType{CertificateTypeX509, CertificateTypeRawPublicKey}
				serverConfig.ServerCertificateTypes = clientConfig.ServerCertificateTypes
			case "server-only":
				serverConfig.ClientAuth = tls.NoClientCert
				clientConfig.RawPublicKeySigner = nil
				clientType = CertificateTypeX509
			case "empty-client":
				serverConfig.ClientAuth = tls.RequestClientCert
				serverConfig.ClientCertificateTypes = nil
				clientConfig.RawPublicKeySigner = nil
				clientType = CertificateTypeX509
			case "raw-client", "x509-fallback", "x509-without-rpk-verifier":
				certificate, roots := testServerCertificate(t)
				serverConfig.Certificates, clientConfig.RootCAs = []tls.Certificate{certificate}, roots
				serverConfig.ServerCertificateTypes = nil
				clientConfig.ServerCertificateTypes = []CertificateType{CertificateTypeRawPublicKey, CertificateTypeX509}
				serverType = CertificateTypeX509
				if mode == "x509-without-rpk-verifier" {
					clientConfig.VerifyPeerRawPublicKey = nil
					serverConfig.ServerCertificateTypes = []CertificateType{CertificateTypeRawPublicKey, CertificateTypeX509}
				}
			}
			if mode == "raw-server" || mode == "x509-fallback" {
				certificate, roots := testClientCertificate(t)
				clientConfig.Certificates, serverConfig.ClientCAs = []tls.Certificate{certificate}, roots
				serverConfig.ClientCertificateTypes = nil
				clientConfig.ClientCertificateTypes = []CertificateType{CertificateTypeRawPublicKey, CertificateTypeX509}
				clientType = CertificateTypeX509
			}
			client, server := completeHandshakePair(t, clientConfig, serverConfig)
			for _, conn := range []*Conn{client, server} {
				state := conn.ConnectionState()
				if state.ServerCertificateType != serverType || state.ClientCertificateType != clientType {
					t.Fatalf("certificate types = %d/%d, want %d/%d", state.ServerCertificateType, state.ClientCertificateType, serverType, clientType)
				}
			}
			if serverType == CertificateTypeRawPublicKey {
				requireRawPublicKey(t, client, serverDER)
			}
			if clientType == CertificateTypeRawPublicKey {
				requireRawPublicKey(t, server, clientDER)
			}
			if _, err := client.WriteDatagram([]byte("RPK")); err != nil {
				t.Fatal(err)
			}
			var payload [8]byte
			if n, _, err := server.ReadDatagram(payload[:]); err != nil || string(payload[:n]) != "RPK" {
				t.Fatalf("application data: %d %v", n, err)
			}
		})
	}
}

func TestRawPublicKeyResumptionAndPostHandshakeAuth(t *testing.T) {
	clientConfig, serverConfig := rawPublicKeyConfigs(t)
	clientConfig.ClientSessionCache = NewLRUClientSessionCache(2)
	clientConfig.PostHandshakeAuth = true
	serverConfig.SessionTicketKey = [32]byte{1}
	serverConfig.SessionTicketLifetime = time.Hour
	serverDER, clientDER := rawPublicKeyDER(t, serverConfig.RawPublicKeySigner), rawPublicKeyDER(t, clientConfig.RawPublicKeySigner)
	cached := issueEarlyDataTicket(t, clientConfig, serverConfig)
	clone := cloneClientSessionState(cached)
	clone.peerRawPublicKey[0] ^= 0xff
	if !bytes.Equal(cached.peerRawPublicKey, serverDER) {
		t.Fatal("session clone aliases RPK")
	}
	client, server := completeHandshakePair(t, clientConfig, serverConfig)
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume {
		t.Fatal("RPK session did not resume")
	}
	requireRawPublicKey(t, client, serverDER)
	requireRawPublicKey(t, server, clientDER)
	if client.ConnectionState().ClientCertificateType != CertificateTypeRawPublicKey || server.ConnectionState().ClientCertificateType != CertificateTypeRawPublicKey {
		t.Fatal("resumed PSK did not negotiate the RPK type for post-handshake authentication")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.RequestClientCertificate(ctx); err != nil {
		t.Fatalf("RPK post-handshake authentication after resumption: %v", err)
	}
	requireRawPublicKey(t, server, clientDER)
	changed := clientConfig.Clone()
	changed.VerifyPeerRawPublicKey = pinRawPublicKey(clientDER)
	changed, err := changed.normalized()
	if err != nil {
		t.Fatal(err)
	}
	if validateClientSession(changed, cached) != nil {
		t.Fatal("accepted a ticket after the server pin changed")
	}
	changed.VerifyPeerRawPublicKey = clientConfig.VerifyPeerRawPublicKey
	changed.ServerCertificateTypes = nil
	if validateClientSession(changed, cached) != nil {
		t.Fatal("accepted an RPK ticket under an X.509-only policy")
	}
	withoutSigner := clientConfig.Clone()
	withoutSigner.RawPublicKeySigner = nil
	withoutSigner.ClientSessionCache = NewLRUClientSessionCache(1)
	withoutSigner.ClientSessionCache.Put("server.test", cached)
	resumedClient, resumedServer := completeHandshakePair(t, withoutSigner, serverConfig)
	if !resumedClient.ConnectionState().DidResume || !resumedServer.ConnectionState().DidResume {
		t.Fatal("resumption incorrectly required the original signing key")
	}
	requireRawPublicKey(t, resumedServer, clientDER)
	if resumedClient.ConnectionState().ClientCertificateType != CertificateTypeX509 || resumedServer.ConnectionState().ClientCertificateType != CertificateTypeX509 {
		t.Fatal("resumption reused a certificate type absent from EncryptedExtensions")
	}
	if err := resumedServer.RequestClientCertificate(ctx); err == nil {
		t.Fatal("requested an unnegotiated RPK credential")
	}

	// A resumed PSK retains the old identity, but a later Certificate uses
	// the type negotiated on this connection (RFC 9846 sections 4.3 and 4.5.1).
	x509Client, x509Server := withoutSigner.Clone(), serverConfig.Clone()
	certificate, roots := testClientCertificate(t)
	x509Client.ClientCertificateTypes = nil
	x509Client.Certificates = []tls.Certificate{certificate}
	x509Client.ClientSessionCache = NewLRUClientSessionCache(1)
	x509Client.ClientSessionCache.Put("server.test", cached)
	x509Server.ClientCertificateTypes = []CertificateType{CertificateTypeRawPublicKey, CertificateTypeX509}
	x509Server.ClientCAs = roots
	_, reauthenticated := completeHandshakePair(t, x509Client, x509Server)
	if !reauthenticated.ConnectionState().DidResume {
		t.Fatal("certificate type change prevented PSK resumption")
	}
	requireRawPublicKey(t, reauthenticated, clientDER)
	if err := reauthenticated.RequestClientCertificate(ctx); err != nil {
		t.Fatal(err)
	}
	if state := reauthenticated.ConnectionState(); len(state.PeerRawPublicKey) != 0 || len(state.PeerCertificates) == 0 || state.ClientCertificateType != CertificateTypeX509 {
		t.Fatal("PHA did not replace the cached RPK identity with the new X.509 identity")
	}
}

func TestRawPublicKeyAlgorithmsAndRetransmission(t *testing.T) {
	for _, algorithm := range []string{"P256", "P384", "P521", "RSA"} {
		t.Run(algorithm, func(t *testing.T) {
			clientConfig, serverConfig := rawPublicKeyConfigs(t)
			var signer crypto.Signer
			var err error
			if algorithm == "RSA" {
				signer, err = rsa.GenerateKey(rand.Reader, 2048)
			} else {
				curve := elliptic.P256()
				if algorithm == "P384" {
					curve = elliptic.P384()
				}
				if algorithm == "P521" {
					curve = elliptic.P521()
				}
				signer, err = ecdsa.GenerateKey(curve, rand.Reader)
			}
			if err != nil {
				t.Fatal(err)
			}
			serverConfig.RawPublicKeySigner = signer
			clientConfig.VerifyPeerRawPublicKey = pinRawPublicKey(rawPublicKeyDER(t, signer))
			serverConfig.CurvePreferences = []tls.CurveID{tls.CurveP256}
			loss := &handshakeLoss{typ: handshakeTypeCertificate, epoch: 2}
			left, right := memoryDatagramPair()
			wire := &observeRecordsConn{Conn: right, observe: loss.observeMessage}
			clientWire := &observeRecordsConn{Conn: left, observe: loss.observeACK}
			client, server := Client(clientWire, clientConfig), Server(wire, serverConfig)
			clientWire.owner = client
			wire.owner = server
			handshakePair(t, client, server)
			loss.requireRecovery(t)
			requireRawPublicKey(t, client, rawPublicKeyDER(t, signer))
		})
	}
}

func TestRawPublicKeyRejectsUntrustedPeer(t *testing.T) {
	for _, mode := range []string{"wrong-server-pin", "wrong-client-pin", "missing-verifier", "no-common-type", "server-signature", "client-signature"} {
		t.Run(mode, func(t *testing.T) {
			clientConfig, serverConfig := rawPublicKeyConfigs(t)
			clientConfig.InsecureSkipVerify = true
			switch mode {
			case "wrong-server-pin":
				clientConfig.VerifyPeerRawPublicKey = pinRawPublicKey(nil)
			case "wrong-client-pin":
				serverConfig.VerifyPeerRawPublicKey = pinRawPublicKey(nil)
			case "missing-verifier":
				clientConfig.VerifyPeerRawPublicKey = nil
			case "no-common-type":
				serverConfig.ServerCertificateTypes = nil
			case "server-signature":
				serverConfig.RawPublicKeySigner = mismatchedRPKSigner{clientConfig.RawPublicKeySigner, serverConfig.RawPublicKeySigner.Public()}
			case "client-signature":
				clientConfig.RawPublicKeySigner = mismatchedRPKSigner{serverConfig.RawPublicKeySigner, clientConfig.RawPublicKeySigner.Public()}
			}
			left, right := memoryDatagramPair()
			client, server := Client(left, clientConfig), Server(right, serverConfig)
			t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
			done := make(chan error, 1)
			go func() { done <- server.Handshake() }()
			clientErr := client.Handshake()
			_ = client.Close()
			if mode == "missing-verifier" {
				_ = server.Close()
			}
			serverErr := <-done
			if clientErr == nil || (mode != "missing-verifier" && serverErr == nil) {
				t.Fatalf("accepted untrusted peer: client=%v server=%v", clientErr, serverErr)
			}
			if mode == "server-signature" || mode == "client-signature" {
				localErr := clientErr
				if mode == "client-signature" {
					localErr = serverErr
				}
				if alert, ok := protocolAlert(localErr); !ok || alert != alertDecryptError {
					t.Fatalf("bad signature alert: %v", localErr)
				}
			}
		})
	}
}

type mismatchedRPKSigner struct {
	crypto.Signer
	public crypto.PublicKey
}

func (s mismatchedRPKSigner) Public() crypto.PublicKey { return s.public }

func TestRawPublicKeyRealUDP(t *testing.T) {
	for _, suite := range []uint16{TLS_AES_128_GCM_SHA256, TLS_AES_256_GCM_SHA384, TLS_CHACHA20_POLY1305_SHA256, TLS_AES_128_CCM_SHA256} {
		t.Run(tls.CipherSuiteName(suite), func(t *testing.T) {
			clientConfig, serverConfig := rawPublicKeyConfigs(t)
			clientConfig.CipherSuites, serverConfig.CipherSuites = []uint16{suite}, []uint16{suite}
			clientConfig.EnableCertificateCompression, serverConfig.EnableCertificateCompression = true, true
			clientConfig.EnableOCSPStapling, serverConfig.EnableOCSPStapling = true, true
			clientConfig.ConnectionID, serverConfig.ConnectionID = []byte("client"), []byte("server")
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			listener, err := Listen("udp4", "127.0.0.1:0", serverConfig)
			if err != nil {
				t.Fatal(err)
			}
			defer listener.Close()
			stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
			defer stop()
			wire, err := net.Dial("udp4", listener.Addr().String())
			if err != nil {
				t.Fatal(err)
			}
			client := Client(wire, clientConfig)
			defer client.Close()
			done := make(chan error, 1)
			go func() { done <- client.HandshakeContext(ctx) }()
			server, err := listener.Accept()
			if err != nil {
				t.Fatal(err)
			}
			defer server.Close()
			serverErr := server.HandshakeContext(ctx)
			if clientErr := <-done; clientErr != nil || serverErr != nil {
				t.Fatalf("handshake: client=%v server=%v", clientErr, serverErr)
			}
			requireRawPublicKey(t, client, rawPublicKeyDER(t, serverConfig.RawPublicKeySigner))
			requireRawPublicKey(t, server, rawPublicKeyDER(t, clientConfig.RawPublicKeySigner))
			for _, pair := range [][2]*Conn{{client, server}, {server, client}} {
				state := pair[0].ConnectionState()
				if state.CipherSuite != suite || !bytes.Equal(state.LocalConnectionID, pair[1].ConnectionState().PeerConnectionID) {
					t.Fatal("cipher suite or CID negotiation mismatch")
				}
				_ = pair[0].SetDeadline(time.Now().Add(time.Second))
				_ = pair[1].SetDeadline(time.Now().Add(time.Second))
				if err := pair[0].SendKeyUpdate(false); err != nil {
					t.Fatal(err)
				}
				if _, err := pair[0].WriteDatagram([]byte("RPK")); err != nil {
					t.Fatal(err)
				}
				var data [8]byte
				if n, _, err := pair[1].ReadDatagram(data[:]); err != nil || string(data[:n]) != "RPK" {
					t.Fatalf("exchange: %d %v", n, err)
				}
			}
		})
	}
}

func TestRawPublicKeyPostHandshakeTrustRevocation(t *testing.T) {
	clientConfig, serverConfig := rawPublicKeyConfigs(t)
	clientConfig.PostHandshakeAuth = true
	var revoked atomic.Bool
	pin := serverConfig.VerifyPeerRawPublicKey
	serverConfig.VerifyPeerRawPublicKey = func(raw []byte) error {
		if revoked.Load() {
			return errors.New("revoked key")
		}
		return pin(raw)
	}
	_, server := completeHandshakePair(t, clientConfig, serverConfig)
	revoked.Store(true)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.RequestClientCertificate(ctx); err == nil || errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("PHA accepted revoked key or timed out: %v", err)
	}
}

func TestRawPublicKeyWireAndPolicy(t *testing.T) {
	clientConfig, serverConfig := rawPublicKeyConfigs(t)
	raw := rawPublicKeyDER(t, serverConfig.RawPublicKeySigner)
	weakRSA, err := rsa.GenerateKey(rand.Reader, 1024) // #nosec G403 -- verifies rejection of weak RPK credentials.
	if err != nil {
		t.Fatal(err)
	}
	agreementKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	agreementDER, err := x509.MarshalPKIXPublicKey(agreementKey.PublicKey())
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		entries []certificateEntry
		alert   uint8
	}{
		{nil, alertDecodeError},
		{[]certificateEntry{{data: raw}, {data: raw}}, alertBadCertificate},
		{[]certificateEntry{{data: []byte{1, 2, 3}}}, alertBadCertificate},
		{[]certificateEntry{{data: rawPublicKeyDER(t, weakRSA)}}, alertBadCertificate},
		{[]certificateEntry{{data: agreementDER}}, alertBadCertificate},
	} {
		_, err := verifyRawPublicKeyMessage(clientConfig, &certificateMessage{certificates: test.entries})
		if got, ok := protocolAlert(err); !ok || got != test.alert {
			t.Fatalf("malformed RPK Certificate alert = %v, want %d", err, test.alert)
		}
	}
	// A verifier owns its input; mutating it must not replace the key whose
	// CertificateVerify signature is checked or the identity published to callers.
	mutating := clientConfig.Clone()
	mutating.VerifyPeerRawPublicKey = func(input []byte) error { clear(input); return nil }
	verified, err := verifyRawPublicKeyMessage(mutating, &certificateMessage{certificates: []certificateEntry{{data: raw}}})
	if err != nil || !bytes.Equal(verified, rawPublicKeyDER(t, serverConfig.RawPublicKeySigner)) {
		t.Fatalf("trust callback changed the authenticated key: %v", err)
	}
	// No compatible signing algorithm produces an empty client Certificate,
	// even when the connection negotiated RPK.
	client := &Conn{config: clientConfig, clientCertificateType: CertificateTypeRawPublicKey}
	certificate, err := client.selectClientCertificate(&certificateRequestMessage{signatureSchemes: []tls.SignatureScheme{tls.PSSWithSHA256}})
	if err != nil || certificate != nil {
		t.Fatalf("selected incompatible RPK signer: %v", err)
	}
	hello := &clientHello{
		cipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, supportedGroups: []tls.CurveID{tls.X25519},
		signatureSchemes: []tls.SignatureScheme{tls.Ed25519}, keyShares: []keyShareEntry{{group: tls.X25519, data: bytes.Repeat([]byte{1}, 32)}},
	}
	if err := (&Conn{config: clientConfig}).offerCertificateTypes(hello); err != nil {
		t.Fatal(err)
	}
	for _, extension := range []uint16{extServerCertificateType, extClientCertificateType} {
		for _, raw := range [][]byte{{}, {0}, {2, 2}, {1, 2, 0}} {
			malformed := cloneClientHello(hello)
			malformed.unknownExtensions[extension] = raw
			wire, err := malformed.marshal()
			if err != nil {
				t.Fatal(err)
			}
			if _, err := parseClientHello(wire); err == nil {
				t.Fatalf("accepted malformed certificate types %d: %x", extension, raw)
			}
		}
		for _, test := range []struct {
			raw   []byte
			alert uint8
		}{
			{nil, alertDecodeError}, {[]byte{0}, alertIllegalParameter}, {[]byte{2, 2}, alertDecodeError},
		} {
			wire, err := (&encryptedExtensions{extensions: map[uint16][]byte{extension: test.raw}}).marshal()
			if err != nil {
				t.Fatal(err)
			}
			message, err := parseEncryptedExtensions(wire)
			if err != nil {
				t.Fatal(err)
			}
			_, _, _, err = validateEncryptedExtensions(hello, &message)
			if got, ok := protocolAlert(err); !ok || got != test.alert {
				t.Fatalf("selection %d: %x: %v", extension, test.raw, err)
			}
		}
		if _, err := validateSelectedCertificateType(&clientHello{}, extension, []byte{2}); err == nil {
			t.Fatal("accepted unsolicited RPK selection")
		}
		second := cloneClientHello(hello)
		second.unknownExtensions[extension][1] = 0
		if equalClientHelloAfterHRR(hello, second, 0) {
			t.Fatal("HRR changed certificate types")
		}
	}

	clone := clientConfig.Clone()
	clone.ServerCertificateTypes[0] = CertificateTypeX509
	clone.ClientCertificateTypes[0] = CertificateTypeX509
	if clientConfig.ServerCertificateTypes[0] != CertificateTypeRawPublicKey || clientConfig.ClientCertificateTypes[0] != CertificateTypeRawPublicKey {
		t.Fatal("Config.Clone aliases certificate types")
	}
	for _, types := range [][]CertificateType{{1}, {2, 2}} {
		if _, err := (&Config{ServerCertificateTypes: types}).normalized(); err == nil {
			t.Fatal("accepted invalid configured certificate types")
		}
	}
}

func TestRawPublicKeyECHAndEarlyData(t *testing.T) {
	clientConfig, serverConfig := rawPublicKeyConfigs(t)
	list, key := testECHConfig(t, "public.test", 27)
	clientConfig.EncryptedClientHelloConfigList = list
	serverConfig.EncryptedClientHelloKeys = []EncryptedClientHelloKey{key}
	clientConfig.ClientSessionCache = NewLRUClientSessionCache(2)
	serverConfig.SessionTicketKey = [32]byte{3}
	serverConfig.MaxEarlyData = 256
	serverConfig.AllowEarlyDataWithoutCookie = true
	_ = issueEarlyDataTicket(t, clientConfig, serverConfig)
	left, right := memoryDatagramPair()
	client, server := Client(left, clientConfig), Server(right, serverConfig)
	t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
	done := make(chan error, 1)
	go func() { done <- server.Handshake() }()
	if _, err := client.WriteEarlyData([]byte("early RPK")); err != nil {
		t.Fatal(err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	if !client.ConnectionState().DidResume || !client.ConnectionState().ECHAccepted || !server.ConnectionState().ECHAccepted {
		t.Fatal("RPK ECH resumption not negotiated")
	}
	var buffer [32]byte
	if n, _, err := server.ReadDatagram(buffer[:]); err != nil || string(buffer[:n]) != "early RPK" {
		t.Fatalf("early data: %d %v", n, err)
	}
	requireRawPublicKey(t, client, rawPublicKeyDER(t, serverConfig.RawPublicKeySigner))
	requireRawPublicKey(t, server, rawPublicKeyDER(t, clientConfig.RawPublicKeySigner))
	hello := testECHClientHello()
	if err := (&Conn{config: clientConfig}).offerCertificateTypes(hello); err != nil {
		t.Fatal(err)
	}
	configs, err := parseECHConfigList(list)
	if err != nil {
		t.Fatal(err)
	}
	outer, err := makeECHOuter(hello, &configs[0], rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	if outer.unknownExtensions[extClientCertificateType] != nil || outer.unknownExtensions[extServerCertificateType] != nil {
		t.Fatal("ECH outer exposes credential preferences")
	}
}

func TestRawPublicKeyExternalPSKSelection(t *testing.T) {
	clientConfig, serverConfig := rawPublicKeyConfigs(t)
	psk, err := NewDirectExternalPSK([]byte("rpk-alternative"), bytes.Repeat([]byte{4}, 32), crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	clientConfig.ExternalPSKs, serverConfig.ExternalPSKs = []*ExternalPSK{psk}, []*ExternalPSK{psk}
	if _, err := serverConfig.normalized(); err == nil {
		t.Fatal("accepted external PSK with client authentication")
	}
	serverConfig.ClientAuth = tls.NoClientCert
	clientConfig.VerifyPeerRawPublicKey = func([]byte) error { t.Fatal("PSK handshake verified RPK"); return nil }
	client, server := completeHandshakePair(t, clientConfig, serverConfig)
	for _, conn := range []*Conn{client, server} {
		state := conn.ConnectionState()
		if len(state.PeerRawPublicKey) != 0 || len(state.PeerCertificates) != 0 || !bytes.Equal(state.ExternalPSKIdentity(), psk.identity) {
			t.Fatal("external PSK mixed with RPK authentication")
		}
	}
}

func TestRawPublicKeyTicketPolicy(t *testing.T) {
	clientConfig, serverConfig := rawPublicKeyConfigs(t)
	serverConfig.SessionTicketLifetime = time.Hour
	serverConfig, err := serverConfig.normalized()
	if err != nil {
		t.Fatal(err)
	}
	state := &sessionTicketState{
		createdAt: time.Now().Unix(), clientAuthAt: time.Now().Unix(), lifetime: 3600,
		psk: bytes.Repeat([]byte{5}, 32), suite: TLS_AES_128_GCM_SHA256,
		maxEarlyData:     256,
		recordSizeLimit:  512,
		peerRawPublicKey: rawPublicKeyDER(t, clientConfig.RawPublicKeySigner),
	}
	wire, err := state.marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parseSessionTicketState(wire)
	if err != nil {
		t.Fatal(err)
	}
	if !validClientAuthenticationTicket(serverConfig, parsed) {
		t.Fatal("rejected valid RPK ticket")
	}
	if parsed.clientAuthAt != state.clientAuthAt || parsed.recordSizeLimit != state.recordSizeLimit || parsed.maxEarlyData != state.maxEarlyData || !bytes.Equal(parsed.peerRawPublicKey, state.peerRawPublicKey) {
		t.Fatal("RPK ticket lost authentication, early data or record limit")
	}
	for i := range len(wire) {
		if _, err := parseSessionTicketState(wire[:i]); err == nil {
			t.Fatalf("accepted truncated RPK ticket at %d", i)
		}
	}
	mutated := bytes.Clone(wire)
	mutated[len(mutated)-len(state.peerRawPublicKey)] ^= 0xff
	if _, err := parseSessionTicketState(mutated); err == nil {
		t.Fatal("accepted malformed ticket public key")
	}
	if _, err := parseSessionTicketState(append(bytes.Clone(wire), 0)); err == nil {
		t.Fatal("accepted trailing RPK ticket data")
	}
	for _, policy := range []string{"pin", "type", "expired-auth", "future-auth", "no-client-auth"} {
		t.Run(policy, func(t *testing.T) {
			config := serverConfig.Clone()
			candidate := *parsed
			switch policy {
			case "pin":
				config.VerifyPeerRawPublicKey = pinRawPublicKey(nil)
			case "type":
				config.ClientCertificateTypes = nil
			case "expired-auth":
				candidate.clientAuthAt -= 7200
			case "future-auth":
				candidate.clientAuthAt += 7200
			case "no-client-auth":
				config.ClientAuth = tls.NoClientCert
			}
			if validClientAuthenticationTicket(config, &candidate) {
				t.Fatal("resumed invalid RPK authentication")
			}
		})
	}
}
