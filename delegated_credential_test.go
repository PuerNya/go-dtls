package dtls13

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"
)

func delegatedParent(t testing.TB, signer crypto.Signer, server bool) tls.Certificate {
	t.Helper()
	usage := x509.ExtKeyUsageClientAuth
	serial, name := int64(9346), "client"
	if server {
		usage = x509.ExtKeyUsageServerAuth
		serial, name = 9345, "server.test"
	}
	now := time.Now().Truncate(time.Second)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(serial), Subject: pkix.Name{CommonName: name},
		DNSNames: []string{"server.test"}, IPAddresses: nil,
		NotBefore: now.Add(-time.Hour), NotAfter: now.Add(30 * 24 * time.Hour),
		KeyUsage: x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign, ExtKeyUsage: []x509.ExtKeyUsage{usage},
		BasicConstraintsValid: true, IsCA: true,
		ExtraExtensions: []pkix.Extension{{Id: oidDelegationUsage, Value: []byte{5, 0}}},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, signer.Public(), signer)
	if err != nil {
		t.Fatal(err)
	}
	leaf, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return tls.Certificate{Certificate: [][]byte{der}, PrivateKey: signer, Leaf: leaf}
}

func delegatedCredentialConfigs(t testing.TB) (*Config, *Config) {
	t.Helper()
	configs := []*Config{
		{EnableDelegatedCredentials: true, ServerName: "server.test", SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond},
		{EnableDelegatedCredentials: true, ClientAuth: tls.RequireAndVerifyClientCert, SessionTicketsDisabled: true, HandshakeTimeout: 2 * time.Second, FlightInterval: 5 * time.Millisecond},
	}
	for i, config := range configs {
		parentKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{byte(10 + i)}, 32))
		parent := delegatedParent(t, parentKey, i == 1)
		key := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{byte(20 + i)}, 32))
		dc, err := NewDelegatedCredential(parent, key, time.Now().Add(time.Hour), i == 1)
		if err != nil {
			t.Fatal(err)
		}
		config.Certificates = []tls.Certificate{parent}
		config.DelegatedCredentials = []*DelegatedCredential{dc}
	}
	configs[0].RootCAs = certificatePool(configs[1].Certificates[0])
	configs[1].ClientCAs = certificatePool(configs[0].Certificates[0])
	return configs[0], configs[1]
}

func requireDelegatedIdentity(t testing.TB, conn *Conn, want []byte) {
	t.Helper()
	state := conn.ConnectionState()
	if !bytes.Equal(state.PeerDelegatedCredential, want) || len(state.PeerCertificates) != 1 || len(state.VerifiedChains) != 1 || len(state.PeerRawPublicKey) != 0 {
		t.Fatalf("incorrect delegated identity: DC=%x certificates=%d chains=%d", state.PeerDelegatedCredential, len(state.PeerCertificates), len(state.VerifiedChains))
	}
	if len(want) > 0 {
		state.PeerDelegatedCredential[0] ^= 1
		if !bytes.Equal(conn.ConnectionState().PeerDelegatedCredential, want) {
			t.Fatal("ConnectionState aliases delegated credential")
		}
	}
}

func TestDelegatedCredentialHandshake(t *testing.T) {
	for _, mode := range []string{"mutual", "UDP", "offline", "server-only", "client-only", "fallback", "compression-OCSP", "HRR-loss", "ECH", "SHA384", "cached", "callbacks", "delegated-callbacks"} {
		t.Run(mode, func(t *testing.T) {
			cc, sc := delegatedCredentialConfigs(t)
			serverDC, clientDC := sc.DelegatedCredentials[0].Credential, cc.DelegatedCredentials[0].Credential
			var loss *handshakeLoss
			switch mode {
			case "offline":
				cc.Certificates, sc.Certificates = nil, nil
			case "server-only":
				sc.EnableDelegatedCredentials = false
				clientDC = nil
			case "client-only":
				cc.EnableDelegatedCredentials = false
				serverDC = nil
			case "fallback":
				cc.EnableDelegatedCredentials, sc.EnableDelegatedCredentials = false, false
				serverDC, clientDC = nil, nil
			case "compression-OCSP":
				cc.EnableCertificateCompression, sc.EnableCertificateCompression = true, true
				cc.EnableOCSPStapling, sc.EnableOCSPStapling = true, true
				cc.DelegatedCredentials[0].Certificate.OCSPStaple = []byte("client OCSP")
				sc.DelegatedCredentials[0].Certificate.OCSPStaple = []byte("server OCSP")
			case "HRR-loss":
				sc.CurvePreferences = []tls.CurveID{tls.CurveP256}
				cc.MTU, sc.MTU = 256, 256
				loss = &handshakeLoss{typ: handshakeTypeCertificate, epoch: 2}
			case "ECH":
				list, key := testECHConfig(t, "public.test", 7)
				cc.EncryptedClientHelloConfigList, sc.EncryptedClientHelloKeys = list, []EncryptedClientHelloKey{key}
			case "SHA384":
				cc.CipherSuites, sc.CipherSuites = []uint16{TLS_AES_256_GCM_SHA384}, []uint16{TLS_AES_256_GCM_SHA384}
			case "cached":
				cc.CachedInformationCache, sc.EnableCachedInformation = NewCachedInformationCache(1), true
				cachedInformationPair(t, cc, sc, 0, nil)
			case "callbacks", "delegated-callbacks":
				clientConfig, serverConfig := cc, sc
				cc.GetClientCertificate = func(info *CertificateRequestInfo) (*tls.Certificate, error) {
					if len(info.DelegatedCredentialSignatureSchemes) == 0 {
						t.Error("callback missing DC offer")
					}
					if mode == "delegated-callbacks" {
						cert := &clientConfig.DelegatedCredentials[0].Certificate
						return cert, info.SupportsCertificate(cert)
					}
					return &clientConfig.Certificates[0], nil
				}
				sc.GetCertificate = func(info *ClientHelloInfo) (*tls.Certificate, error) {
					if len(info.DelegatedCredentialSignatureSchemes) == 0 {
						t.Error("callback missing DC offer")
					}
					if mode == "delegated-callbacks" {
						cert := &serverConfig.DelegatedCredentials[0].Certificate
						return cert, info.SupportsCertificate(cert)
					}
					return &serverConfig.Certificates[0], nil
				}
				cc, sc = cc.Clone(), sc.Clone()
				if mode == "callbacks" {
					serverDC, clientDC = nil, nil
				}
			}
			var client, server *Conn
			if mode == "cached" || loss != nil {
				want := uint8(0)
				if mode == "cached" {
					want = cachedCert | cachedCertRequest
				}
				client, server = cachedInformationPair(t, cc, sc, want, loss)
			} else if mode == "UDP" {
				listener, err := Listen("udp4", "127.0.0.1:0", sc)
				if err != nil {
					t.Fatal(err)
				}
				defer listener.Close()
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()
				stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
				defer stop()
				wire, err := net.Dial("udp4", listener.Addr().String())
				if err != nil {
					t.Fatal(err)
				}
				client = Client(wire, cc)
				defer client.Close()
				done := make(chan error, 1)
				go func() { done <- client.HandshakeContext(ctx) }()
				server, err = listener.Accept()
				if err != nil {
					t.Fatal(err)
				}
				defer server.Close()
				if err = server.HandshakeContext(ctx); err != nil {
					t.Fatal(err)
				}
				if err = <-done; err != nil {
					t.Fatal(err)
				}
			} else {
				client, server = completeHandshakePair(t, cc, sc)
			}
			requireDelegatedIdentity(t, client, serverDC)
			requireDelegatedIdentity(t, server, clientDC)
			if mode == "ECH" && !client.ConnectionState().ECHAccepted {
				t.Fatal("ECH not accepted")
			}
			if mode == "compression-OCSP" && (!bytes.Equal(client.ConnectionState().OCSPResponse, []byte("server OCSP")) || !bytes.Equal(server.ConnectionState().OCSPResponse, []byte("client OCSP"))) {
				t.Fatal("OCSP lost alongside DC")
			}
			for _, pair := range [][2]*Conn{{client, server}, {server, client}} {
				_ = pair[1].SetReadDeadline(time.Now().Add(time.Second))
				if _, err := pair[0].WriteDatagram([]byte("DC")); err != nil {
					t.Fatal(err)
				}
				var data [4]byte
				if n, _, err := pair[1].ReadDatagram(data[:]); err != nil || string(data[:n]) != "DC" {
					t.Fatalf("datagram: %d %v", n, err)
				}
			}
		})
	}
}

func TestDelegatedCredentialSignatureVector(t *testing.T) {
	parentKey := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{31}, 32))
	key := ed25519.NewKeyFromSeed(bytes.Repeat([]byte{32}, 32))
	parent := delegatedParent(t, parentKey, true)
	expiry := time.Now().Truncate(time.Second).Add(time.Hour)
	dc, err := NewDelegatedCredential(parent, key, expiry, true)
	if err != nil {
		t.Fatal(err)
	}
	spki, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		t.Fatal(err)
	}
	// Independent RFC 9345 §4 wire and signature construction, without the
	// production encoder, parser, context builder or signing functions.
	cred := binary.BigEndian.AppendUint32(nil, uint32(expiry.Sub(parent.Leaf.NotBefore)/time.Second))
	cred = append(cred, 8, 7, 0, 0, byte(len(spki)))
	cred = append(cred, spki...)
	cred = append(cred, 8, 7)
	input := append(bytes.Repeat([]byte{32}, 64), []byte("TLS, server delegated credentials\x00")...)
	input = append(input, parent.Certificate[0]...)
	input = append(input, cred...)
	sig := ed25519.Sign(parentKey, input)
	want := append(cred, 0, byte(len(sig)))
	want = append(want, sig...)
	if !bytes.Equal(dc.Credential, want) {
		t.Fatalf("RFC signature vector mismatch: %x", dc.Credential)
	}
	if !key.Equal(dc.Certificate.PrivateKey) {
		t.Fatal("credential retained parent key")
	}
	for _, curve := range []elliptic.Curve{elliptic.P256(), elliptic.P384(), elliptic.P521()} {
		key, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = NewDelegatedCredential(parent, key, expiry, true); err != nil {
			t.Fatal(err)
		}
	}
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewDelegatedCredential(delegatedParent(t, rsaKey, true), key, expiry, true); err != nil {
		t.Fatal("RSA-PSS delegation:", err)
	}
	if _, err := NewDelegatedCredential(parent, rsaKey, expiry, true); err == nil {
		t.Fatal("accepted rsaEncryption delegated key")
	}
}

func TestDelegatedCredentialValidation(t *testing.T) {
	cc, sc := delegatedCredentialConfigs(t)
	raw := sc.DelegatedCredentials[0].Credential
	parent := sc.Certificates[0].Leaf
	dc, err := parseDelegatedCredential(raw)
	if err != nil {
		t.Fatal(err)
	}
	expiry := parent.NotBefore.Add(time.Duration(dc.validTime) * time.Second)
	for _, name := range []string{"expired", "too-long", "parent-expiry", "future-parent", "missing-usage", "critical-usage", "bad-usage", "no-digital-signature", "unoffered-delegation", "unoffered-key", "role", "signature", "certificate-binding"} {
		t.Run(name, func(t *testing.T) {
			p, b, now := *parent, bytes.Clone(raw), time.Now()
			signatures, delegated, server := defaultSignatureSchemes(), delegatedCredentialSchemes(), true
			switch name {
			case "expired":
				now = expiry.Add(time.Nanosecond)
			case "too-long":
				// Keep NotBefore in the past so only the seven-day bound fails.
				binary.BigEndian.PutUint32(b[:4], uint32(now.Add(7*24*time.Hour+time.Second).Sub(p.NotBefore)/time.Second))
				signed, err := parseDelegatedCredential(b)
				if err != nil {
					t.Fatal(err)
				}
				copy(b[len(b)-ed25519.SignatureSize:], ed25519.Sign(sc.Certificates[0].PrivateKey.(ed25519.PrivateKey), delegatedCredentialInput(p.Raw, signed.signed, true)))
				if _, err := validateDelegatedCredential(b, &p, now.Add(time.Second), true, signatures, delegated); err != nil {
					t.Fatal("control credential is invalid:", err)
				}
			case "parent-expiry":
				p.NotAfter = expiry
			case "future-parent":
				now = p.NotBefore.Add(-time.Nanosecond)
			case "missing-usage":
				p.Extensions = nil
			case "critical-usage":
				p.Extensions = []pkix.Extension{{Id: oidDelegationUsage, Critical: true, Value: []byte{5, 0}}}
			case "bad-usage":
				p.Extensions = []pkix.Extension{{Id: oidDelegationUsage, Value: []byte{5, 1, 0}}}
			case "no-digital-signature":
				p.KeyUsage = 0
			case "unoffered-delegation":
				signatures = nil
			case "unoffered-key":
				delegated = nil
			case "role":
				server = false
			case "signature":
				b[len(b)-1] ^= 1
			case "certificate-binding":
				p.Raw = append(bytes.Clone(p.Raw), 0)
			}
			_, err := validateDelegatedCredential(b, &p, now, server, signatures, delegated)
			if got, ok := protocolAlert(err); !ok || got != alertIllegalParameter {
				t.Fatalf("alert=%d: %v", got, err)
			}
		})
	}
	limit := time.Now().Truncate(time.Second).Add(7 * 24 * time.Hour)
	week, err := NewDelegatedCredential(sc.Certificates[0], sc.DelegatedCredentials[0].Certificate.PrivateKey.(crypto.Signer), limit, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := validateDelegatedCredential(week.Credential, parent, limit.Add(-7*24*time.Hour), true, defaultSignatureSchemes(), delegatedCredentialSchemes()); err != nil {
		t.Fatal("exact seven-day credential:", err)
	}
	if _, err := validateDelegatedCredential(raw, parent, expiry, true, defaultSignatureSchemes(), delegatedCredentialSchemes()); err != nil {
		t.Fatal("credential invalid at inclusive expiry:", err)
	}
	for n := range len(raw) {
		if _, err := parseDelegatedCredential(raw[:n]); err == nil {
			t.Fatalf("accepted prefix %d", n)
		}
	}
	if _, err := parseDelegatedCredential(append(bytes.Clone(raw), 0)); err == nil {
		t.Fatal("accepted trailing data")
	}
	for _, scheme := range []tls.SignatureScheme{tls.PSSWithSHA256, tls.ECDSAWithP256AndSHA256, tls.SignatureScheme(0xffff)} {
		b := bytes.Clone(raw)
		binary.BigEndian.PutUint16(b[4:6], uint16(scheme))
		if _, err := parseDelegatedCredential(b); err == nil {
			t.Fatalf("accepted incompatible DC scheme %x", scheme)
		}
	}
	// Session policy is current, not just the original certificate policy.
	cc.Time = func() time.Time { return expiry.Add(time.Second) }
	if validResumedDelegatedCredential(cc, raw, []*x509.Certificate{parent}, true) {
		t.Fatal("expired DC resumable")
	}
}

func TestDelegatedCredentialCertificateExtensions(t *testing.T) {
	_, sc := delegatedCredentialConfigs(t)
	dc := sc.DelegatedCredentials[0]
	for _, test := range []struct {
		name    string
		entries []certificateEntry
		offered bool
		alert   uint8
	}{
		{"unsolicited", []certificateEntry{{data: dc.Certificate.Certificate[0], extensions: map[uint16][]byte{34: dc.Credential}}}, false, alertUnexpectedMessage},
		{"non-leaf", []certificateEntry{{data: dc.Certificate.Certificate[0]}, {data: dc.Certificate.Certificate[0], extensions: map[uint16][]byte{34: dc.Credential}}}, true, alertIllegalParameter},
	} {
		t.Run(test.name, func(t *testing.T) {
			wire, err := (&certificateMessage{certificates: test.entries}).marshal()
			if err != nil {
				t.Fatal(err)
			}
			message, err := parseCertificateMessage(wire, 1<<20)
			if err != nil {
				t.Fatal(err)
			}
			err = validateCertificateMessageWithRequests(message, nil, false, test.offered)
			if got, ok := protocolAlert(err); !ok || got != test.alert {
				t.Fatalf("alert=%d: %v", got, err)
			}
		})
	}
	verify := &certificateVerifyMessage{algorithm: tls.ECDSAWithP256AndSHA256, signature: []byte{1}}
	err := verifyPeerAuthenticationSignature(nil, nil, dc.Credential, nil, verify, nil, true)
	if got, ok := protocolAlert(err); !ok || got != alertIllegalParameter {
		t.Fatalf("CV algorithm alert=%d: %v", got, err)
	}
	// Build duplicate extensions directly because the sender's map cannot.
	w := newWireBuilder(512)
	w.bytes8(nil)
	list := w.startVector24()
	w.bytes24(dc.Certificate.Certificate[0])
	extensions := w.startVector16()
	for range 2 {
		w.u16(int(extDelegatedCredential))
		w.bytes16(dc.Credential)
	}
	w.endVector16(extensions)
	w.endVector24(list)
	_, err = parseCertificateMessage(w.b, 1<<20)
	if got, ok := protocolAlert(err); !ok || got != alertIllegalParameter {
		t.Fatalf("duplicate DC alert=%d: %v", got, err)
	}
}

func TestDelegatedCredentialResumptionAndPHA(t *testing.T) {
	cc, sc := delegatedCredentialConfigs(t)
	cc.SessionTicketsDisabled, sc.SessionTicketsDisabled = false, false
	cc.ClientSessionCache, cc.PostHandshakeAuth = NewLRUClientSessionCache(2), true
	sc.SessionTicketKey = [32]byte{1}
	sc.MaxEarlyData, sc.AllowEarlyDataWithoutCookie = 256, true
	cached := issueEarlyDataTicket(t, cc, sc)
	clone := cloneClientSessionState(cached)
	clone.peerDelegatedCredential[0] ^= 1
	if !bytes.Equal(cached.peerDelegatedCredential, sc.DelegatedCredentials[0].Credential) {
		t.Fatal("session clone aliases DC")
	}
	left, right := memoryDatagramPair()
	client, server := Client(left, cc), Server(right, sc)
	t.Cleanup(func() { _ = client.Close(); _ = server.Close() })
	done := make(chan error, 1)
	go func() { done <- server.Handshake() }()
	if n, err := client.WriteEarlyData([]byte("early DC")); err != nil || n != 8 {
		t.Fatalf("early data write: %d %v", n, err)
	}
	if err := <-done; err != nil {
		t.Fatal(err)
	}
	_ = server.SetReadDeadline(time.Now().Add(time.Second))
	var data [32]byte
	if n, _, err := server.ReadDatagram(data[:]); err != nil || string(data[:n]) != "early DC" {
		t.Fatalf("early data read: %d %v", n, err)
	}
	if !client.ConnectionState().DidResume || !server.ConnectionState().DidResume {
		t.Fatal("DC session did not resume")
	}
	requireDelegatedIdentity(t, client, sc.DelegatedCredentials[0].Credential)
	requireDelegatedIdentity(t, server, cc.DelegatedCredentials[0].Credential)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := server.RequestClientCertificate(ctx); err != nil {
		t.Fatal(err)
	}
	requireDelegatedIdentity(t, server, cc.DelegatedCredentials[0].Credential)
	clientDC, err := parseDelegatedCredential(cc.DelegatedCredentials[0].Credential)
	if err != nil {
		t.Fatal(err)
	}
	sc.Time = func() time.Time {
		return cc.Certificates[0].Leaf.NotBefore.Add(time.Duration(clientDC.validTime)*time.Second + time.Second)
	}
	cc.Time = sc.Time
	client, server = completeHandshakePair(t, cc, sc)
	if client.ConnectionState().DidResume || server.ConnectionState().DidResume {
		t.Fatal("expired DC session resumed")
	}
	requireDelegatedIdentity(t, client, nil)
	requireDelegatedIdentity(t, server, nil)
}

func TestDelegatedCredentialSelectionRejectsInvalid(t *testing.T) {
	for _, serverRole := range []bool{false, true} {
		for _, reason := range []string{"unoffered", "expired", "wrong-key", "wrong-role"} {
			t.Run(fmt.Sprintf("server=%t/%s", serverRole, reason), func(t *testing.T) {
				cc, sc := delegatedCredentialConfigs(t)
				local, peer := cc, sc
				if serverRole {
					local, peer = sc, cc
				}
				d := local.DelegatedCredentials[0]
				schemes := delegatedCredentialSchemes()
				switch reason {
				case "unoffered":
					peer.EnableDelegatedCredentials = false
					schemes = nil
				case "expired":
					local.Time = func() time.Time { return time.Now().Add(2 * time.Hour) }
				case "wrong-key":
					d.Certificate.PrivateKey = ed25519.NewKeyFromSeed(make([]byte, 32))
				case "wrong-role":
					other, err := NewDelegatedCredential(local.Certificates[0], d.Certificate.PrivateKey.(crypto.Signer), time.Now().Add(time.Hour), !serverRole)
					if err != nil {
						t.Fatal(err)
					}
					d.Credential = other.Credential
				}
				left, right := memoryDatagramPair()
				client, server := Client(left, cc), Server(right, sc)
				defer client.Close()
				defer server.Close()
				conn := client
				if serverRole {
					conn = server
				}
				var err error
				conn.config, err = conn.config.normalized()
				if err != nil {
					t.Fatal(err)
				}
				if conn.selectDelegatedCredential(defaultSignatureSchemes(), nil, schemes, serverRole, "", nil, nil) != nil {
					t.Fatal("automatic selection accepted invalid DC")
				}
				// Callbacks must not bypass peer offers, time or key checks.
				if serverRole {
					conn.config.GetCertificate = func(*ClientHelloInfo) (*tls.Certificate, error) { return &d.Certificate, nil }
				} else {
					conn.config.GetClientCertificate = func(*CertificateRequestInfo) (*tls.Certificate, error) { return &d.Certificate, nil }
				}
				done := make(chan error, 1)
				go func() { done <- server.Handshake() }()
				clientErr, serverErr := client.Handshake(), <-done
				if clientErr == nil || serverErr == nil {
					t.Fatalf("invalid callback credential accepted: %v / %v", clientErr, serverErr)
				}
			})
		}
	}
}

func TestDelegatedCredentialCachedPolicy(t *testing.T) {
	cc, sc := delegatedCredentialConfigs(t)
	cc.CachedInformationCache, sc.EnableCachedInformation = NewCachedInformationCache(1), true
	cachedInformationPair(t, cc, sc, 0, nil)
	for _, reason := range []string{"expired", "untrusted", "callback", "disabled"} {
		t.Run(reason, func(t *testing.T) {
			changed := cc.Clone()
			switch reason {
			case "expired":
				changed.Time = func() time.Time { return time.Now().Add(2 * time.Hour) }
			case "untrusted":
				changed.RootCAs = x509.NewCertPool()
			case "disabled":
				changed.EnableDelegatedCredentials = false
			case "callback":
				changed.VerifyPeerCertificate = func([][]byte, [][]*x509.Certificate) error { return errors.New("revoked parent") }
			}
			if reason == "disabled" {
				client, _ := cachedInformationPair(t, changed, sc, cachedCertRequest, nil)
				requireDelegatedIdentity(t, client, nil)
				return
			}
			left, right := memoryDatagramPair()
			client, server := Client(left, changed), Server(right, sc)
			defer client.Close()
			defer server.Close()
			done := make(chan error, 1)
			go func() { done <- server.Handshake() }()
			if err := client.Handshake(); err == nil {
				t.Fatal("cached DC bypassed current policy")
			}
			if err := <-done; err == nil {
				t.Fatal("server accepted failed authentication")
			}
		})
	}
}

func TestInteropWolfSSLDelegatedCredentialFallback(t *testing.T) {
	for _, role := range []string{"server", "client"} {
		t.Run(role, func(t *testing.T) {
			options := wolfSSLInteropOptions{
				configure: func(_ *testing.T, _ string, config *Config) {
					config.EnableDelegatedCredentials = true
					if role == "client" {
						_, delegated := delegatedCredentialConfigs(t)
						config.DelegatedCredentials = delegated.DelegatedCredentials
					}
				},
				connected: func(t *testing.T, conn *Conn, _ int) {
					if len(conn.ConnectionState().PeerDelegatedCredential) != 0 {
						t.Fatal("wolfSSL unexpectedly negotiated DC")
					}
				},
			}
			if role == "server" {
				testInteropWolfSSLServerOptions(t, options)
			} else {
				testInteropWolfSSLClientOptions(t, options)
			}
		})
	}
}

func FuzzDelegatedCredential(f *testing.F) {
	_, sc := delegatedCredentialConfigs(f)
	f.Add(sc.DelegatedCredentials[0].Credential)
	f.Add([]byte{})
	f.Fuzz(func(t *testing.T, data []byte) { _, _ = parseDelegatedCredential(data) })
}
