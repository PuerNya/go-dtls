package dtls13

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/crypto/ocsp"
)

// Peers are built from pinned sources outside the checkout. These tests never
// download or build dependencies, and never count X.509 fallback as DC success.
func TestInteropDTLS13Peers(t *testing.T) {
	for _, peer := range []string{"NSS", "BORINGSSL", "OPENSSL"} {
		t.Run(peer, func(t *testing.T) {
			binary := os.Getenv("DTLS13_" + peer + "_PEER")
			if binary == "" {
				t.Skip("set DTLS13_" + peer + "_PEER to the testdata/interop peer")
			}
			for _, role := range []string{"server", "client"} {
				roleName := role
				if peer == "OPENSSL" && role == "server" {
					roleName = "server-ACKLoss"
				}
				for _, mode := range []string{"basic", "resume", "early", "early-io", "sha384", "chacha", "ccm-fallback", "hybrid", "rpk-fallback", "cid-rrc-fallback", "external-psk-fallback", "p256", "p384", "hrr", "fragmented", "mutual", "pha", "ocsp", "alpn", "grease", "compression-fallback", "cached-fallback", "keyupdate", "dc", "dc-resume", "dc-early", "dc-fragmented", "dc-sha384", "dc-fallback", "dc-bad-signature", "dc-expired"} {
					if (mode == "ccm-fallback" && peer == "OPENSSL") ||
						(mode == "rpk-fallback" && peer != "NSS") ||
						(mode == "external-psk-fallback" && peer != "BORINGSSL") {
						continue
					}
					if (strings.HasSuffix(mode, "early") || mode == "early-io") && (peer == "BORINGSSL" || role == "server") {
						continue
					}
					if mode == "pha" && peer != "OPENSSL" {
						continue
					}
					if role == "client" && (mode == "dc-bad-signature" || mode == "dc-expired") {
						continue
					}
					if strings.HasPrefix(mode, "dc") && peer != "NSS" && (peer != "BORINGSSL" || role != "server") && mode != "dc-fallback" {
						continue
					}
					t.Run(roleName+"/"+mode, func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, mode) })
				}
				if peer == "OPENSSL" {
					t.Run(roleName+"/ccm", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "ccm") })
				}
				if peer != "NSS" {
					t.Run(roleName+"/rpk", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "rpk") })
				}
				if peer == "BORINGSSL" {
					t.Run(roleName+"/ech", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "ech") })
				}
				if peer == "OPENSSL" && role == "server" {
					t.Run(roleName+"/ech-server-rejection", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "ech-server-rejection") })
				}
				if peer == "OPENSSL" && role == "client" {
					t.Run(roleName+"/ech-client-unavailable", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "ech-client-unavailable") })
				}
				if peer == "NSS" {
					t.Run(roleName+"/ech-unavailable", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "ech-unavailable") })
				}
				if peer == "NSS" || peer == "OPENSSL" {
					t.Run(roleName+"/external-psk", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "external-psk") })
				}
				if peer == "BORINGSSL" {
					t.Run(roleName+"/imported-psk", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "imported-psk") })
				}
				if role == "server" {
					t.Run(roleName+"/ticket-request-fallback", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "ticket-request-fallback") })
				}
				metadataMode := "metadata-unoffered"
				if role == "server" {
					metadataMode = "metadata-ech-rejected"
				}
				t.Run(roleName+"/"+metadataMode, func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, metadataMode) })
				if peer == "OPENSSL" && role == "server" {
					t.Run("server/initial-ACK-rejection", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "initial-ACK-rejection") })
				}
				if peer == "NSS" && role == "client" {
					t.Run("client/DC-in-CertificateRequest-rejection", func(t *testing.T) { testDTLS13Peer(t, peer, binary, role, "DC-CR-rejection") })
				}
			}
		})
	}
}

func writePeerCredential(t testing.TB, directory, name string, certificate tls.Certificate) (string, string, string, string) {
	t.Helper()
	certPEM, keyPEM := filepath.Join(directory, name+".pem"), filepath.Join(directory, name+"-key.pem")
	certDER, keyDER := filepath.Join(directory, name+".der"), filepath.Join(directory, name+"-key.der")
	var chain []byte
	for _, der := range certificate.Certificate {
		chain = append(chain, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})...)
	}
	key, err := x509.MarshalPKCS8PrivateKey(certificate.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	for path, data := range map[string][]byte{certPEM: chain, keyPEM: pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key}), certDER: certificate.Certificate[0], keyDER: key} {
		if err := os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	return certPEM, keyPEM, certDER, keyDER
}

func testDTLS13Peer(t *testing.T, peer, binary, role, mode string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	serverCert, clientCert := delegatedParent(t, serverKey, true), delegatedParent(t, clientKey, false)
	serverRoot := serverCert
	if mode == "ocsp" {
		leaf := *serverRoot.Leaf
		leaf.SerialNumber, leaf.Subject.CommonName = big.NewInt(9347), "server-leaf"
		leaf.RawSubject = nil
		leaf.IsCA, leaf.KeyUsage = false, x509.KeyUsageDigitalSignature
		for _, ext := range leaf.Extensions {
			if ext.Id.Equal(oidDelegationUsage) {
				leaf.ExtraExtensions = append(leaf.ExtraExtensions, ext)
			}
		}
		der, err := x509.CreateCertificate(rand.Reader, &leaf, serverRoot.Leaf, serverKey.Public(), serverKey)
		if err != nil {
			t.Fatal(err)
		}
		serverCert.Leaf, err = x509.ParseCertificate(der)
		if err != nil {
			t.Fatal(err)
		}
		serverCert.Certificate = [][]byte{der, serverRoot.Certificate[0]}
	}
	delegatedKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	dc, err := NewDelegatedCredential(serverCert, delegatedKey, time.Now().Add(time.Hour), true)
	if err != nil {
		t.Fatal(err)
	}
	goServer := role == "client"
	config := &Config{
		ServerName: "server.test", RootCAs: certificatePool(serverRoot), ClientCAs: certificatePool(clientCert),
		Certificates: []tls.Certificate{clientCert}, EnableDelegatedCredentials: true,
		CipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second, FlightInterval: 20 * time.Millisecond,
	}
	if goServer {
		config.Certificates = []tls.Certificate{serverCert}
	}
	var echList []byte
	var echKey EncryptedClientHelloKey
	if mode == "ech" || mode == "ech-server-rejection" || mode == "ech-client-unavailable" {
		echList, echKey = testECHConfig(t, "public.test", 17)
		if mode != "ech-client-unavailable" {
			if goServer {
				config.EncryptedClientHelloKeys = []EncryptedClientHelloKey{echKey}
			} else {
				config.EncryptedClientHelloConfigList = echList
			}
		}
	}
	if mode == "metadata-ech-rejected" {
		config.EncryptedClientHelloConfigList, _ = testECHConfig(t, "server.test", 23)
		config.HandshakeMetadata = []byte("private handshake metadata")
	}
	if mode == "metadata-unoffered" {
		config.AcceptHandshakeMetadata = func(*ClientHelloInfo, []byte) ([]byte, error) {
			return nil, errors.New("metadata callback without an offer")
		}
	}
	if mode == "sha384" || mode == "dc-sha384" {
		config.CipherSuites = []uint16{TLS_AES_256_GCM_SHA384}
	}
	if mode == "chacha" {
		config.CipherSuites = []uint16{TLS_CHACHA20_POLY1305_SHA256}
	}
	if mode == "ccm" || mode == "ccm-fallback" {
		config.CipherSuites = []uint16{TLS_AES_128_CCM_SHA256}
		if mode == "ccm-fallback" {
			config.CipherSuites = append(config.CipherSuites, TLS_AES_128_GCM_SHA256)
		}
	}
	if mode == "hybrid" {
		config.CurvePreferences = []tls.CurveID{tls.X25519MLKEM768}
	}
	var verifiedRPK atomic.Bool
	if mode == "rpk" || mode == "rpk-fallback" {
		config.ServerCertificateTypes = []CertificateType{CertificateTypeRawPublicKey}
		config.ClientCertificateTypes = []CertificateType{CertificateTypeRawPublicKey}
		if mode == "rpk-fallback" {
			config.ServerCertificateTypes = append(config.ServerCertificateTypes, CertificateTypeX509)
			config.ClientCertificateTypes = append(config.ClientCertificateTypes, CertificateTypeX509)
		}
		config.RawPublicKeySigner = config.Certificates[0].PrivateKey.(crypto.Signer)
		peerCert := serverCert
		if goServer {
			peerCert = clientCert
		}
		verify := pinRawPublicKey(rawPublicKeyDER(t, peerCert.PrivateKey.(crypto.Signer)))
		config.VerifyPeerRawPublicKey = func(raw []byte) error {
			verifiedRPK.Store(true)
			return verify(raw)
		}
	}
	if mode == "cid-rrc-fallback" {
		config.ConnectionID = []byte("go-peer")
	}
	if mode == "external-psk" || mode == "external-psk-fallback" || mode == "imported-psk" {
		var psk *ExternalPSK
		var err error
		if mode == "imported-psk" {
			psk, err = ImportExternalPSK([]byte("interop-psk"), bytes.Repeat([]byte{0x42}, 32), []byte("interop-context"), crypto.SHA256)
		} else {
			psk, err = NewDirectExternalPSK([]byte("interop-psk"), bytes.Repeat([]byte{0x42}, 32), crypto.SHA256)
		}
		if err != nil {
			t.Fatal(err)
		}
		config.ExternalPSKs = []*ExternalPSK{psk}
	}
	if mode == "ticket-request-fallback" {
		config.SessionTicketRequest = SessionTicketRequest{Enabled: true}
	}
	if mode == "p256" || (mode == "hrr" && goServer) {
		config.CurvePreferences = []tls.CurveID{tls.CurveP256}
	}
	if mode == "p384" {
		config.CurvePreferences = []tls.CurveID{tls.CurveP384}
	}
	if mode == "fragmented" || mode == "dc-fragmented" {
		config.MTU = 256
	}
	if (mode == "mutual" || mode == "pha" || mode == "DC-CR-rejection") && goServer {
		config.ClientAuth = tls.RequireAndVerifyClientCert
		if peer == "NSS" && mode != "DC-CR-rejection" {
			config.EnableDelegatedCredentials = false
		}
	}
	var clientAuthCalls atomic.Int32
	if mode == "pha" {
		config.PostHandshakeAuth = true
		if !goServer {
			config.GetClientCertificate = func(*CertificateRequestInfo) (*tls.Certificate, error) {
				clientAuthCalls.Add(1)
				return &clientCert, nil
			}
		}
	}
	if mode == "ocsp" {
		config.EnableOCSPStapling = true
		staple, err := ocsp.CreateResponse(serverRoot.Leaf, serverRoot.Leaf, ocsp.Response{Status: ocsp.Good, SerialNumber: serverCert.Leaf.SerialNumber, ThisUpdate: time.Now().Add(-time.Minute), NextUpdate: time.Now().Add(time.Hour)}, serverKey)
		if err != nil {
			t.Fatal(err)
		}
		config.Certificates[0].OCSPStaple = staple
	}
	if mode == "alpn" {
		config.NextProtos = []string{"coap"}
	}
	if mode == "grease" {
		config.EnableGREASE, config.EncryptedClientHelloGrease = true, true
	}
	if mode == "compression-fallback" {
		config.EnableCertificateCompression = true
	}
	if mode == "cached-fallback" {
		config.EnableCachedInformation, config.CachedInformationCache = true, NewCachedInformationCache(1)
	}
	rounds := 1
	earlyIO := mode == "early-io"
	early := strings.HasSuffix(mode, "early") || earlyIO
	if mode == "resume" || mode == "dc-resume" || mode == "ticket-request-fallback" || early {
		rounds = 2
		config.SessionTicketsDisabled = false
		config.ClientSessionCache = NewLRUClientSessionCache(2)
		config.SessionTicketKey = [32]byte{9}
	}
	if early {
		config.MaxEarlyData, config.AllowEarlyDataWithoutCookie = 256, true
		config.EnableEarlyDataIO = earlyIO
	}
	useDC := strings.HasPrefix(mode, "dc") && mode != "dc-fallback"
	if goServer && useDC {
		config.DelegatedCredentials = []*DelegatedCredential{dc}
	}
	if mode == "dc-fallback" {
		config.EnableDelegatedCredentials = false
		if goServer {
			config.DelegatedCredentials = []*DelegatedCredential{dc}
		}
	}
	if mode == "dc-bad-signature" && !goServer {
		dc.Credential[len(dc.Credential)-1] ^= 1
	}
	if mode == "dc-expired" && !goServer {
		config.Time = func() time.Time { return time.Now().Add(2 * time.Hour) }
	}
	directory := t.TempDir()
	remoteCert := serverCert
	if goServer {
		remoteCert = clientCert
	}
	certPEM, keyPEM, certDER, keyDER := writePeerCredential(t, directory, "peer", remoteCert)
	if mode == "ocsp" {
		if err := os.WriteFile(certPEM+".ocsp", config.Certificates[0].OCSPStaple, 0600); err != nil {
			t.Fatal(err)
		}
	}
	_, delegatedPEM, _, delegatedDER := writePeerCredential(t, directory, "delegated", dc.Certificate)
	dcPath := filepath.Join(directory, "dc.bin")
	if err := os.WriteFile(dcPath, dc.Credential, 0600); err != nil {
		t.Fatal(err)
	}
	if mode == "ech" || mode == "ech-server-rejection" || mode == "ech-client-unavailable" {
		data := echKey.Config
		if goServer || mode == "ech-client-unavailable" {
			data = echList
		}
		if err := os.WriteFile(dcPath, data, 0600); err != nil {
			t.Fatal(err)
		}
		if !goServer && mode != "ech-client-unavailable" {
			private := echKey.PrivateKey
			if peer == "OPENSSL" {
				key, err := ecdh.X25519().NewPrivateKey(private)
				if err != nil {
					t.Fatal(err)
				}
				der, err := x509.MarshalPKCS8PrivateKey(key)
				if err != nil {
					t.Fatal(err)
				}
				private = append(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}), pem.EncodeToMemory(&pem.Block{Type: "ECHCONFIG", Bytes: echList})...)
			}
			if err := os.WriteFile(delegatedPEM, private, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	trustCert := clientCert
	if goServer {
		trustCert = serverCert
	}
	trustRoot := trustCert
	if goServer {
		trustRoot = serverRoot
	}
	if err := os.WriteFile(certPEM+".trust", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: trustRoot.Certificate[0]}), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(certDER+".trust", trustCert.Certificate[0], 0600); err != nil {
		t.Fatal(err)
	}
	port := unusedUDPPort(t)
	var listener *Listener
	if goServer {
		listener, err = Listen("udp4", "127.0.0.1:"+strconv.Itoa(port), config)
		if err != nil {
			t.Fatal(err)
		}
		defer listener.Close()
	}
	cmd := exec.CommandContext(ctx, binary, role, strconv.Itoa(port), mode, certPEM, keyPEM, certDER, keyDER, dcPath, delegatedPEM, delegatedDER) // #nosec G204 G702 -- explicitly configured local opt-in test peer, without a shell.
	var output lockedBuffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	defer func() {
		cancel()
		if done != nil {
			err := <-done
			if t.Failed() {
				t.Logf("peer exit: %v\n%s", err, output.String())
			}
		}
	}()
	for round := range rounds {
		var conn *Conn
		var droppedACK, observedHRR atomic.Bool
		if goServer {
			accepted := make(chan *Conn, 1)
			go func() { c, _ := listener.Accept(); accepted <- c }()
			select {
			case err := <-done:
				done = nil
				t.Fatalf("peer exited before connecting: %v\n%s", err, output.String())
			case conn = <-accepted:
				if conn == nil {
					t.Fatalf("accept failed: %s", output.String())
				}
			case <-ctx.Done():
				t.Fatalf("accept timeout: %s", output.String())
			}
		} else {
			deadline := time.Now().Add(3 * time.Second)
			for !strings.Contains(output.String(), fmt.Sprintf("READY %d", round)) && time.Now().Before(deadline) {
				time.Sleep(time.Millisecond)
			}
			if !strings.Contains(output.String(), fmt.Sprintf("READY %d", round)) {
				t.Fatalf("peer did not become ready: %s", output.String())
			}
			udp, e := net.Dial("udp4", "127.0.0.1:"+strconv.Itoa(port))
			if e != nil {
				t.Fatal(e)
			}
			if early && round != 0 {
				conn = ClientEarly(udp, config)
			} else {
				conn = Client(udp, config)
			}
		}
		defer conn.Close()
		if mode == "hrr" || peer == "OPENSSL" && !goServer {
			conn.conn = &observeRecordsConn{Conn: conn.conn, owner: conn, observe: func(r record) (bool, error) {
				// Stock OpenSSL rejects ACKs while waiting for the client's Finished.
				if peer == "OPENSSL" && !goServer && mode != "initial-ACK-rejection" && r.typ == recordTypeACK && r.epoch == 2 {
					droppedACK.Store(true)
					return true, nil
				}
				if r.epoch == 0 && r.typ == recordTypeHandshake {
					fragments, err := parseHandshakeFragments(r.payload)
					if err != nil {
						return false, err
					}
					for _, f := range fragments {
						if (goServer && f.typ == handshakeTypeServerHello && f.offset == 0 && len(f.body) >= 34 && bytes.Equal(f.body[2:34], helloRetryRequestRandom[:])) || (!goServer && f.typ == handshakeTypeClientHello && f.messageSequence > 0) {
							observedHRR.Store(true)
						}
					}
				}
				return false, nil
			}}
		}
		if earlyIO && round != 0 {
			exchangeEarlyBeforeFinished(t, conn, goServer)
		} else if early && round != 0 && !goServer {
			if n, err := conn.WriteDatagram([]byte("early")); err != nil || n != 5 {
				t.Fatalf("early write: %d %v\n%s", n, err, output.String())
			}
		}
		if err = conn.HandshakeContext(ctx); err != nil {
			if mode == "ech-server-rejection" && errors.Is(err, AlertError(alertIllegalParameter)) && !conn.ConnectionState().ECHAccepted {
				return
			}
			if mode == "metadata-ech-rejected" {
				if _, ok := errors.AsType[*ECHRejectionError](err); ok && conn.ConnectionState().HandshakeMetadata == nil {
					return
				}
			}
			if mode == "DC-CR-rejection" && errors.Is(err, AlertError(alertIllegalParameter)) {
				return
			}
			if mode == "initial-ACK-rejection" {
				if errors.Is(err, AlertError(alertUnexpectedMessage)) {
					return
				}
				// UDP can report the peer's exit before delivering its fatal alert.
				select {
				case peerErr := <-done:
					done = nil
					if peerErr != nil && strings.Contains(output.String(), "ossl_statem_server_read_transition:unexpected message") {
						return
					}
				case <-ctx.Done():
				}
			}
			if !goServer && (mode == "dc-bad-signature" || mode == "dc-expired") {
				if alert, ok := protocolAlert(err); ok && alert == alertIllegalParameter {
					return
				}
			}
			t.Fatalf("%s %s %s handshake: %v\n%s", peer, role, mode, err, output.String())
		}
		if mode == "initial-ACK-rejection" || mode == "DC-CR-rejection" || mode == "metadata-ech-rejected" || mode == "ech-server-rejection" {
			t.Fatal("peer limitation no longer reproduces; update the supported matrix")
		}
		if mode == "hrr" && !observedHRR.Load() {
			t.Fatal("key-share HRR did not occur")
		}
		if peer == "OPENSSL" && !goServer && !droppedACK.Load() {
			t.Fatal("ACK-loss scenario did not drop an ACK")
		}
		if !goServer && (mode == "dc-bad-signature" || mode == "dc-expired") {
			t.Fatal("accepted invalid peer DC")
		}
		state := conn.ConnectionState()
		if mode == "metadata-unoffered" && state.HandshakeMetadata != nil {
			t.Fatal("metadata negotiated without an offer")
		}
		if early && round != 0 && !goServer && state.EarlyData != EarlyDataAccepted {
			t.Fatalf("early data status = %v\n%s", state.EarlyData, output.String())
		}
		if state.DidResume != (round != 0) {
			t.Fatal("session was not resumed")
		}
		wantCipher := config.CipherSuites[0]
		if mode == "ccm-fallback" {
			wantCipher = TLS_AES_128_GCM_SHA256
		}
		if state.Version != VersionDTLS13 || state.CipherSuite != wantCipher {
			t.Fatalf("negotiated %#v", state)
		}
		if mode == "rpk-fallback" && (state.ServerCertificateType != CertificateTypeX509 || state.ClientCertificateType != CertificateTypeX509 || len(state.PeerRawPublicKey) != 0) {
			t.Fatalf("RPK fallback negotiated unexpected credential: %#v", state)
		}
		if mode == "rpk" && (state.ServerCertificateType != CertificateTypeRawPublicKey || (!goServer && len(state.PeerRawPublicKey) == 0)) {
			t.Fatalf("RPK not negotiated: %#v", state)
		}
		if mode == "rpk" && !goServer && !verifiedRPK.Load() {
			t.Fatal("RPK verification callback was not invoked")
		}
		if mode == "cid-rrc-fallback" && (len(state.LocalConnectionID) != 0 || len(state.PeerConnectionID) != 0 || state.ReturnRoutabilityCheck) {
			t.Fatalf("CID/RRC fallback negotiated unexpected feature: %#v", state)
		}
		if mode == "external-psk-fallback" && len(state.ExternalPSKIdentity()) != 0 {
			t.Fatalf("external PSK fallback negotiated unexpected PSK: %#v", state)
		}
		if (mode == "external-psk" || mode == "imported-psk") && !bytes.Equal(state.ExternalPSKIdentity(), []byte("interop-psk")) {
			t.Fatalf("external PSK identity = %x", state.ExternalPSKIdentity())
		}
		if mode == "ech" && !state.ECHAccepted {
			t.Fatal("ECH not accepted")
		}
		if !goServer {
			want := []byte(nil)
			if useDC {
				want = dc.Credential
			}
			if !bytes.Equal(state.PeerDelegatedCredential, want) {
				t.Fatalf("DC negotiation: got %x want %x", state.PeerDelegatedCredential, want)
			}
			if mode == "ocsp" && !bytes.Equal(state.OCSPResponse, config.Certificates[0].OCSPStaple) {
				t.Fatal("OCSP staple not received")
			}
		}
		if mode == "alpn" && state.NegotiatedProtocol != "coap" {
			t.Fatal("ALPN not negotiated")
		}
		if goServer && state.ServerName != "server.test" {
			t.Fatal("SNI not received")
		}
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if !goServer {
			if _, err := conn.WriteDatagram([]byte("interop")); err != nil {
				t.Fatal(err)
			}
		}
		var data [64]byte
		if early && !earlyIO && round != 0 && goServer {
			if n, info, err := conn.ReadDatagram(data[:]); err != nil || string(data[:n]) != "early" || !info.EarlyData || !conn.earlyAccepted {
				t.Fatalf("early read: %d %v\n%s", n, err, output.String())
			}
		}
		n, _, err := conn.ReadDatagram(data[:])
		if err != nil || string(data[:n]) != "interop" {
			t.Fatalf("datagram %d %v\n%s", n, err, output.String())
		}
		if goServer {
			if mode == "pha" {
				if err := conn.RequestClientCertificate(ctx); err != nil {
					t.Fatal(err)
				}
				if len(conn.ConnectionState().VerifiedChains) == 0 {
					t.Fatal("PHA identity was not verified")
				}
			}
			if _, err = conn.WriteDatagram(data[:n]); err != nil {
				t.Fatal(err)
			}
		}
		if mode == "pha" && !goServer && clientAuthCalls.Load() == 0 {
			t.Fatal("peer did not request PHA")
		}
		// Compare independently derived exporter bytes, not just handshake completion.
		exporter, err := state.ExportKeyingMaterial("interop", nil, 32)
		if err != nil {
			t.Fatal(err)
		}
		deadline := time.Now().Add(time.Second)
		for !strings.Contains(output.String(), "EXPORTER="+hex.EncodeToString(exporter)) && time.Now().Before(deadline) {
			time.Sleep(time.Millisecond)
		}
		if !strings.Contains(output.String(), "EXPORTER="+hex.EncodeToString(exporter)) {
			t.Fatalf("exporter mismatch: %s", output.String())
		}
		if goServer && peer == "NSS" && useDC && !strings.Contains(output.String(), "PEER_DC=1") {
			t.Fatalf("NSS did not authenticate DC: %s", output.String())
		}
		if goServer && peer == "NSS" && mode == "dc-fallback" && !strings.Contains(output.String(), "PEER_DC=0") {
			t.Fatal("fallback unexpectedly used DC")
		}
		if _, err := conn.WriteDatagram([]byte("done")); err != nil {
			t.Fatal(err)
		}
		if rounds > 1 && !goServer && round == 0 {
			deadline := time.Now().Add(time.Second)
			for {
				if cached, ok := config.ClientSessionCache.Get("server.test"); ok && len(cached.ticket) != 0 {
					break
				}
				if time.Now().After(deadline) {
					t.Fatalf("peer issued no usable ticket: %s", output.String())
				}
				time.Sleep(time.Millisecond)
			}
		}
	}
	select {
	case err := <-done:
		done = nil
		if err != nil {
			t.Fatalf("peer exit: %v\n%s", err, output.String())
		}
	case <-ctx.Done():
		t.Fatalf("peer did not exit: %s", output.String())
	}
}
