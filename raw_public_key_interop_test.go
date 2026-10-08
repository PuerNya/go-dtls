package dtls13

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestInteropWolfSSLRawPublicKey(t *testing.T) {
	peer := os.Getenv("WOLFSSL_RPK_PEER")
	if peer == "" {
		t.Skip("set WOLFSSL_RPK_PEER to the testdata/wolfssl-rpk peer built with HAVE_RPK")
	}
	root, _, _ := wolfSSLPaths(t)
	for _, role := range []string{"client", "server"} {
		for _, mode := range []string{"mutual", "server-only", "raw-client", "raw-server", "sha384", "ccm", "chacha", "hrr", "fragmented", "wrong-pin", "keyupdate", "pha", "resume", "resume-pha", "early", "ech"} {
			t.Run(role+"/"+mode, func(t *testing.T) {
				if mode == "ech" {
					if role == "server" && os.Getenv("GO_DTLS_PROBE_WOLFSSL_RPK_ECH_SERVER") != "1" {
						t.Skip("wolfSSL ECH server parses the DTLS inner ClientHello as TLS; set GO_DTLS_PROBE_WOLFSSL_RPK_ECH_SERVER=1 to probe")
					}
					if role == "client" && os.Getenv("GO_DTLS_PROBE_WOLFSSL_SKIPS") != "1" {
						t.Skip("wolfSSL DTLS ECH client requires the transcript/offset fixes; set GO_DTLS_PROBE_WOLFSSL_SKIPS=1 for a corrected peer")
					}
				}
				testWolfSSLRPK(t, root, peer, role, mode)
			})
		}
	}
}

func testWolfSSLRPK(t *testing.T, root, peer, role, mode string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	goServer := role == "client"
	localName, peerName := "server", "client"
	if !goServer {
		localName, peerName = peerName, localName
	}
	local := wolfSSLCertificate(t, root, localName)
	remote := wolfSSLCertificate(t, root, peerName)
	localDER := rawPublicKeyDER(t, local.PrivateKey.(crypto.Signer))
	peerDER := rawPublicKeyDER(t, remote.PrivateKey.(crypto.Signer))
	var verifications atomic.Int32
	config := &Config{
		ServerName:             "example.com",
		ServerCertificateTypes: []CertificateType{CertificateTypeRawPublicKey},
		ClientCertificateTypes: []CertificateType{CertificateTypeRawPublicKey},
		RawPublicKeySigner:     local.PrivateKey.(crypto.Signer),
		CipherSuites:           []uint16{TLS_AES_128_GCM_SHA256},
		VerifyPeerRawPublicKey: func(raw []byte) error {
			verifications.Add(1)
			return pinRawPublicKey(peerDER)(raw)
		},
		ClientAuth: tls.RequireAndVerifyClientCert, SessionTicketsDisabled: true,
		HandshakeTimeout: 5 * time.Second, FlightInterval: 20 * time.Millisecond,
	}
	credential := filepath.Join("certs", "rpk", peerName+"-cert-rpk.der")
	trust := filepath.Join("certs", "rpk", localName+"-cert-rpk.der")
	peerRaw := true
	resume := mode == "resume" || mode == "resume-pha" || mode == "early"
	echFile := "-"
	switch mode {
	case "server-only":
		config.ClientAuth = tls.NoClientCert
		if goServer {
			peerRaw = false
		} else {
			config.RawPublicKeySigner = nil
		}
	case "raw-client", "raw-server":
		if mode == "raw-client" {
			config.ServerCertificateTypes = nil
		} else {
			config.ClientCertificateTypes = nil
		}
		localRaw := goServer == (mode == "raw-server")
		peerRaw = !localRaw
		if localRaw {
			credential = filepath.Join("certs", peerName+"-cert.pem")
			config.RootCAs = wolfSSLRootCAs(t, root, filepath.Join("certs", "ca-cert.pem"))
			config.ClientCAs = wolfSSLClientCAs(t, root)
		} else {
			config.Certificates = []tls.Certificate{local}
			trust = filepath.Join("certs", localName+"-cert.pem")
		}
	case "hrr":
		if goServer {
			config.CurvePreferences = []tls.CurveID{tls.CurveP256}
		} else {
			config.CurvePreferences = []tls.CurveID{tls.X25519, tls.CurveP256}
		}
	case "fragmented":
		config.MTU = 256
	case "sha384":
		config.CipherSuites = []uint16{TLS_AES_256_GCM_SHA384}
	case "ccm":
		config.CipherSuites = []uint16{TLS_AES_128_CCM_SHA256}
	case "chacha":
		config.CipherSuites = []uint16{TLS_CHACHA20_POLY1305_SHA256}
	case "wrong-pin":
		config.VerifyPeerRawPublicKey = pinRawPublicKey(localDER)
	case "pha", "resume-pha":
		config.PostHandshakeAuth = true
	case "ech":
		echFile = filepath.Join(t.TempDir(), "ech.bin")
		if goServer {
			list, key := testECHConfig(t, "public.test", 1)
			config.EncryptedClientHelloKeys = []EncryptedClientHelloKey{key}
			if err := os.WriteFile(echFile, list, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	if resume {
		config.SessionTicketsDisabled = false
		config.SessionTicketKey = [32]byte{9}
		config.ClientSessionCache = NewLRUClientSessionCache(2)
	}
	if mode == "early" {
		config.MaxEarlyData = 256
		config.AllowEarlyDataWithoutCookie = true
		config.CurvePreferences = []tls.CurveID{tls.CurveP256}
	}
	port := 0
	var listener *Listener
	if goServer {
		var err error
		listener, err = Listen("udp4", "127.0.0.1:0", config)
		if err != nil {
			t.Fatal(err)
		}
		defer listener.Close()
		stop := context.AfterFunc(ctx, func() { _ = listener.Close() })
		defer stop()
		port = listener.Addr().(*net.UDPAddr).Port
	}
	cmd := exec.CommandContext(ctx, peer, role, strconv.Itoa(port), credential,
		filepath.Join("certs", peerName+"-key.pem"), trust, mode, echFile) // #nosec G204 G702 -- explicitly supplied local interoperability fixture.
	cmd.Dir = root
	var output lockedBuffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	waited := false
	defer func() {
		cancel()
		if !waited {
			<-done
		}
	}()
	rounds := 1
	if resume {
		rounds = 2
	}
	for round := range rounds {
		if !goServer {
			ready := "READY " + strconv.Itoa(round) + " "
			for {
				_, address, found := strings.Cut(output.String(), ready)
				line, _, complete := strings.Cut(address, "\n")
				if found && complete {
					var err error
					port, err = strconv.Atoi(strings.TrimSpace(line))
					if err != nil || port <= 0 || port > 65535 {
						t.Fatalf("invalid peer port %q: %v", line, err)
					}
					break
				}
				select {
				case err := <-done:
					waited = true
					t.Fatalf("peer exited before ready: %v\n%s", err, output.String())
				case <-ctx.Done():
					t.Fatalf("peer readiness timed out: %s", output.String())
				case <-time.After(10 * time.Millisecond):
				}
			}
			if mode == "ech" {
				var err error
				config.EncryptedClientHelloConfigList, err = os.ReadFile(echFile) // #nosec G304 -- fixed name in this test's private temporary directory.
				if err != nil {
					t.Fatal(err)
				}
			}
		}
		var conn *Conn
		var err error
		if goServer {
			conn, err = listener.Accept()
		} else {
			var wire net.Conn
			wire, err = net.Dial("udp4", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)))
			if err == nil {
				conn = Client(wire, config)
			}
		}
		if err != nil {
			t.Fatalf("connect: %v\n%s", err, output.String())
		}
		t.Cleanup(func() { _ = conn.Close() })
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if mode == "early" && round == 1 && !goServer {
			_, err = conn.WriteEarlyData([]byte("early"))
		} else {
			err = conn.HandshakeContext(ctx)
		}
		if mode == "wrong-pin" {
			if alert, ok := protocolAlert(err); !ok || alert != alertAccessDenied {
				t.Fatalf("wrong pin error: %v\n%s", err, output.String())
			}
			return
		}
		if err != nil {
			t.Fatalf("handshake %d: %v\n%s", round, err, output.String())
		}
		requireResumptionOnSecondConnection(t, conn, round)
		requireWolfSSLConnection(t, conn, wolfSSLInteropOptions{suites: config.CipherSuites})
		if peerRaw {
			requireRawPublicKey(t, conn, peerDER)
		} else if state := conn.ConnectionState(); len(state.PeerRawPublicKey) != 0 || (mode != "server-only" && len(state.VerifiedChains) == 0) {
			t.Fatal("unexpected X.509 or anonymous peer identity")
		}
		if mode == "ech" && !conn.ConnectionState().ECHAccepted {
			t.Fatal("ECH was not accepted")
		}
		if mode == "keyupdate" && !goServer {
			if err := conn.SendKeyUpdate(true); err != nil {
				t.Fatal(err)
			}
		}
		if (mode == "pha" || mode == "resume-pha") && goServer {
			if err := conn.RequestClientCertificate(ctx); err != nil {
				t.Fatalf("PHA: %v\n%s", err, output.String())
			}
			wantVerifications := int32(2)
			if mode == "resume-pha" {
				wantVerifications = int32(2 * (round + 1))
			}
			if verifications.Load() != wantVerifications {
				t.Fatalf("PHA verification count = %d, want %d", verifications.Load(), wantVerifications)
			}
		}
		var buffer [2048]byte
		if mode == "early" && round == 1 && goServer {
			n, _, err := conn.ReadDatagram(buffer[:])
			if err != nil || string(buffer[:n]) != "early" {
				t.Fatalf("early data: %q %v\n%s", buffer[:n], err, output.String())
			}
		}
		for range 2 {
			if !goServer {
				if _, err := conn.WriteDatagram([]byte("ping")); err != nil {
					t.Fatal(err)
				}
			}
			n, _, err := conn.ReadDatagram(buffer[:])
			if err != nil || !bytes.Equal(buffer[:n], []byte("ping")) {
				t.Fatalf("exchange: %q %v\n%s", buffer[:n], err, output.String())
			}
			if goServer {
				if _, err := conn.WriteDatagram(buffer[:n]); err != nil {
					t.Fatal(err)
				}
			}
		}
		if mode == "keyupdate" {
			requireWolfSSLKeyUpdate(t, conn, round)
			conn.writeMu.Lock()
			sendEpoch := conn.sendCipher.epoch
			conn.writeMu.Unlock()
			if sendEpoch < 4 {
				t.Fatalf("KeyUpdate did not advance the send epoch: %d", sendEpoch)
			}
		}
		if (mode == "pha" || mode == "resume-pha") && !goServer {
			conn.writeMu.Lock()
			requested := conn.hasClientAuthRequestSeq
			conn.writeMu.Unlock()
			if !requested {
				t.Fatal("peer did not request PHA")
			}
		}
		_ = conn.Close()
	}
	select {
	case err := <-done:
		waited = true
		if err != nil || !strings.Contains(output.String(), "DONE "+strconv.Itoa(rounds-1)) {
			t.Fatalf("peer: %v\n%s", err, output.String())
		}
	case <-ctx.Done():
		t.Fatalf("peer did not finish: %s", output.String())
	}
}
