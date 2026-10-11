package dtls13

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
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

type peerBenchmarkFeature struct {
	name, mode string
}

var peerBenchmarkFeatures = map[string][]peerBenchmarkFeature{
	"NSS": {
		{"CertificateAES128GCM", "basic"}, {"AES256GCM", "sha384"}, {"ChaCha20Poly1305", "chacha"},
		{"P256", "p256"}, {"P384", "p384"}, {"HelloRetryRequest", "hrr"}, {"X25519MLKEM768", "hybrid"},
		{"FragmentedHandshake", "fragmented"}, {"MutualTLS", "mutual"}, {"ALPN", "alpn"},
		{"KeyUpdate", "keyupdate"}, {"OCSPStapling", "ocsp"}, {"ExternalPSK", "external-psk"},
		{"SessionResumption", "resume"}, {"MutualTLSSessionResumption", "mutual-resume"},
		{"EarlyData", "early"}, {"DelegatedCredential", "dc"}, {"DelegatedCredentialResumption", "dc-resume"},
	},
	"BORINGSSL": {
		{"CertificateAES128GCM", "basic"}, {"AES256GCM", "sha384"}, {"ChaCha20Poly1305", "chacha"},
		{"P256", "p256"}, {"P384", "p384"}, {"HelloRetryRequest", "hrr"}, {"X25519MLKEM768", "hybrid"},
		{"FragmentedHandshake", "fragmented"}, {"MutualTLS", "mutual"}, {"ALPN", "alpn"},
		{"KeyUpdate", "keyupdate"}, {"OCSPStapling", "ocsp"}, {"RawPublicKey", "rpk"}, {"ImportedPSK", "imported-psk"},
		{"SessionResumption", "resume"}, {"MutualTLSSessionResumption", "mutual-resume"},
		{"ECH", "ech"}, {"DelegatedCredential", "dc"}, {"DelegatedCredentialResumption", "dc-resume"},
	},
	"OPENSSL": {
		{"CertificateAES128GCM", "basic"}, {"AES256GCM", "sha384"}, {"ChaCha20Poly1305", "chacha"},
		{"AES128CCM", "ccm"}, {"P256", "p256"}, {"P384", "p384"}, {"HelloRetryRequest", "hrr"}, {"X25519MLKEM768", "hybrid"},
		{"FragmentedHandshake", "fragmented"}, {"MutualTLS", "mutual"}, {"ALPN", "alpn"},
		{"KeyUpdate", "keyupdate"}, {"OCSPStapling", "ocsp"}, {"RawPublicKey", "rpk"},
		{"ExternalPSK", "external-psk"}, {"SessionResumption", "resume"}, {"MutualTLSSessionResumption", "mutual-resume"}, {"EarlyData", "early"},
		{"PostHandshakeAuthentication", "pha"},
	},
}

func peerBenchmarkPair(mode string) bool {
	return mode == "resume" || mode == "mutual-resume" || mode == "dc-resume" || mode == "early"
}

type peerBenchmarkCredential struct{ pem, keyPEM, der, keyDER string }

type peerBenchmarkServerResult struct {
	err       error
	exporters []string
}

type peerBenchmarkFixture struct {
	server, serverRoot, client             tls.Certificate
	serverFiles, clientFiles               peerBenchmarkCredential
	dc                                     *DelegatedCredential
	dcFile, delegatedPEM, delegatedDER     string
	echList                                []byte
	echKey                                 EncryptedClientHelloKey
	echListFile, echConfigFile, echKeyFile string
	ocspResponse                           []byte
}

func newPeerBenchmarkFixture(b *testing.B, mode string) peerBenchmarkFixture {
	b.Helper()
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	f := peerBenchmarkFixture{server: delegatedParent(b, serverKey, true), client: delegatedParent(b, clientKey, false)}
	if mode == "fragmented" && len(f.server.Certificate[0]) <= 256 {
		b.Fatal("server certificate cannot prove MTU-256 fragmentation")
	}
	f.serverRoot = f.server
	if mode == "ocsp" {
		leaf := *f.server.Leaf
		leaf.SerialNumber, leaf.Subject.CommonName, leaf.RawSubject = big.NewInt(9347), "server-leaf", nil
		leaf.IsCA, leaf.KeyUsage = false, x509.KeyUsageDigitalSignature
		for _, ext := range leaf.Extensions {
			if ext.Id.Equal(oidDelegationUsage) {
				leaf.ExtraExtensions = append(leaf.ExtraExtensions, ext)
			}
		}
		der, err := x509.CreateCertificate(rand.Reader, &leaf, f.serverRoot.Leaf, serverKey.Public(), serverKey)
		if err != nil {
			b.Fatal(err)
		}
		f.server.Leaf, err = x509.ParseCertificate(der)
		if err != nil {
			b.Fatal(err)
		}
		f.server.Certificate = [][]byte{der, f.serverRoot.Certificate[0]}
	}
	dir := b.TempDir()
	f.serverFiles.pem, f.serverFiles.keyPEM, f.serverFiles.der, f.serverFiles.keyDER = writePeerCredential(b, dir, "server", f.server)
	f.clientFiles.pem, f.clientFiles.keyPEM, f.clientFiles.der, f.clientFiles.keyDER = writePeerCredential(b, dir, "client", f.client)
	for _, item := range []struct {
		files              peerBenchmarkCredential
		trustPEM, trustDER tls.Certificate
	}{
		{f.serverFiles, f.client, f.client}, {f.clientFiles, f.serverRoot, f.server},
	} {
		if err := os.WriteFile(item.files.pem+".trust", pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: item.trustPEM.Certificate[0]}), 0600); err != nil {
			b.Fatal(err)
		}
		if err := os.WriteFile(item.files.der+".trust", item.trustDER.Certificate[0], 0600); err != nil {
			b.Fatal(err)
		}
	}
	delegatedKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		b.Fatal(err)
	}
	f.dc, err = NewDelegatedCredential(f.server, delegatedKey, time.Now().Add(time.Hour), true)
	if err != nil {
		b.Fatal(err)
	}
	_, f.delegatedPEM, _, f.delegatedDER = writePeerCredential(b, dir, "delegated", f.dc.Certificate)
	f.dcFile = filepath.Join(dir, "dc.bin")
	f.echList, f.echKey = testECHConfig(b, "public.test", 17)
	f.echListFile, f.echConfigFile, f.echKeyFile = filepath.Join(dir, "ech-list.bin"), filepath.Join(dir, "ech-config.bin"), filepath.Join(dir, "ech-key.bin")
	response, err := ocsp.CreateResponse(f.serverRoot.Leaf, f.serverRoot.Leaf, ocsp.Response{Status: ocsp.Good, SerialNumber: f.server.Leaf.SerialNumber, ThisUpdate: time.Now().Add(-time.Minute), NextUpdate: time.Now().Add(time.Hour)}, serverKey)
	if err != nil {
		b.Fatal(err)
	}
	f.ocspResponse = response
	for path, data := range map[string][]byte{
		f.dcFile: f.dc.Credential, f.echListFile: f.echList, f.echConfigFile: f.echKey.Config,
		f.echKeyFile: f.echKey.PrivateKey, f.serverFiles.pem + ".ocsp": response, f.clientFiles.pem + ".ocsp": response,
	} {
		if err := os.WriteFile(path, data, 0600); err != nil {
			b.Fatal(err)
		}
	}
	return f
}

func (f *peerBenchmarkFixture) config(b *testing.B, mode string, server bool) *Config {
	b.Helper()
	config := &Config{ServerName: "server.test", RootCAs: certificatePool(f.serverRoot),
		ClientCAs: certificatePool(f.client), Certificates: []tls.Certificate{f.client},
		CipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, CurvePreferences: []tls.CurveID{tls.X25519},
		SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second, FlightInterval: 20 * time.Millisecond}
	if server {
		config.Certificates = []tls.Certificate{f.server}
	}
	switch mode {
	case "sha384":
		config.CipherSuites = []uint16{TLS_AES_256_GCM_SHA384}
	case "chacha":
		config.CipherSuites = []uint16{TLS_CHACHA20_POLY1305_SHA256}
	case "ccm":
		config.CipherSuites = []uint16{TLS_AES_128_CCM_SHA256}
	case "p256":
		config.CurvePreferences = []tls.CurveID{tls.CurveP256}
	case "p384":
		config.CurvePreferences = []tls.CurveID{tls.CurveP384}
	case "hrr":
		config.CurvePreferences = []tls.CurveID{tls.X25519, tls.CurveP256}
		if server {
			config.CurvePreferences = []tls.CurveID{tls.CurveP256}
		}
	case "hybrid":
		config.CurvePreferences = []tls.CurveID{tls.X25519MLKEM768}
	case "fragmented":
		config.MTU = 256
	case "mutual", "mutual-resume":
		if server {
			config.ClientAuth = tls.RequireAndVerifyClientCert
		}
	case "alpn":
		config.NextProtos = []string{"coap"}
	case "ocsp":
		config.EnableOCSPStapling = true
		if server {
			config.Certificates[0].OCSPStaple = f.ocspResponse
		}
	case "rpk":
		config.ServerCertificateTypes = []CertificateType{CertificateTypeRawPublicKey}
		config.ClientCertificateTypes = []CertificateType{CertificateTypeRawPublicKey}
		config.RawPublicKeySigner = config.Certificates[0].PrivateKey.(crypto.Signer)
		trusted := f.server
		if server {
			trusted = f.client
		}
		config.VerifyPeerRawPublicKey = pinRawPublicKey(rawPublicKeyDER(b, trusted.PrivateKey.(crypto.Signer)))
	case "external-psk":
		psk, err := NewDirectExternalPSK([]byte("interop-psk"), bytes.Repeat([]byte{0x42}, 32), crypto.SHA256)
		if err != nil {
			b.Fatal(err)
		}
		config.ExternalPSKs = []*ExternalPSK{psk}
	case "imported-psk":
		psk, err := ImportExternalPSK([]byte("interop-psk"), bytes.Repeat([]byte{0x42}, 32), []byte("interop-context"), crypto.SHA256)
		if err != nil {
			b.Fatal(err)
		}
		config.ExternalPSKs = []*ExternalPSK{psk}
	case "resume", "early":
		config.SessionTicketsDisabled = false
		config.SessionTicketKey = [32]byte{9}
		if !server {
			config.ClientSessionCache = NewLRUClientSessionCache(2)
		}
		if mode == "early" && server {
			config.MaxEarlyData, config.AllowEarlyDataWithoutCookie = 256, true
			config.EarlyDataReplayCache = NewLRUEarlyDataReplayCache(1024)
		}
	case "ech":
		if server {
			config.EncryptedClientHelloKeys = []EncryptedClientHelloKey{f.echKey}
		} else {
			config.EncryptedClientHelloConfigList = f.echList
		}
	case "dc", "dc-resume":
		config.EnableDelegatedCredentials = true
		if server {
			config.DelegatedCredentials = []*DelegatedCredential{f.dc}
		}
	case "pha":
		config.PostHandshakeAuth = true
		if server {
			config.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}
	if mode == "mutual-resume" || mode == "dc-resume" {
		config.SessionTicketsDisabled = false
		config.SessionTicketKey = [32]byte{9}
		if !server {
			config.ClientSessionCache = NewLRUClientSessionCache(2)
		}
	}
	return config
}

func (f *peerBenchmarkFixture) peerArgs(role string, port int, mode string, repeats int) []string {
	files := f.clientFiles
	if role == "server" {
		files = f.serverFiles
	}
	extra, key := f.dcFile, f.delegatedPEM
	if mode == "ech" {
		extra, key = f.echListFile, f.echKeyFile
		if role == "server" {
			extra = f.echConfigFile
		}
	}
	return []string{role, strconv.Itoa(port), mode, files.pem, files.keyPEM, files.der, files.keyDER,
		extra, key, f.delegatedDER, strconv.Itoa(repeats)}
}

func BenchmarkDTLS13PeerRealUDP(b *testing.B) {
	for _, peer := range []string{"NSS", "BORINGSSL", "OPENSSL"} {
		binary := os.Getenv("DTLS13_" + peer + "_PEER")
		b.Run(peer, func(b *testing.B) {
			if binary == "" {
				b.Skip("set DTLS13_" + peer + "_PEER to the testdata/interop peer")
			}
			if _, err := os.Stat(binary); err != nil { // #nosec G703 -- opt-in local test executable path.
				b.Fatal(err)
			}
			for _, feature := range peerBenchmarkFeatures[peer] {
				b.Run(feature.name, func(b *testing.B) {
					fixture := newPeerBenchmarkFixture(b, feature.mode)
					benchmarkPeerDirections(b, peer, binary, &fixture, feature)
				})
			}
		})
	}
}

func benchmarkPeerDirections(b *testing.B, peer, binary string, fixture *peerBenchmarkFixture, feature peerBenchmarkFeature) {
	b.Helper()
	for _, direction := range []string{"GoClient/GoServer", "GoClient/" + peer + "Server", peer + "Client/GoServer", peer + "Client/" + peer + "Server"} {
		b.Run(direction, func(b *testing.B) {
			if peer == "NSS" && feature.mode == "fragmented" && strings.HasSuffix(direction, "/NSSServer") {
				b.Skip("NSS fixture cannot enforce a fragmented handshake without a configurable DTLS MTU")
			}
			if feature.mode == "early" && strings.HasSuffix(direction, peer+"Server") {
				b.Skip(peer + " server does not support DTLS 1.3 early data in this fixture")
			}
			if strings.HasPrefix(feature.mode, "dc") && peer == "BORINGSSL" && strings.HasPrefix(direction, peer+"Client") {
				b.Skip("BoringSSL client does not authenticate delegated credentials in DTLS 1.3")
			}
			benchPeerDirection(b, peer, binary, fixture, feature, direction)
		})
	}
}

func benchPeerDirection(b *testing.B, peer, binary string, f *peerBenchmarkFixture, feature peerBenchmarkFeature, direction string) {
	b.Helper()
	repeats := b.N * 4
	clientIsGo, serverIsGo := strings.HasPrefix(direction, "GoClient/"), strings.HasSuffix(direction, "/GoServer")
	port := unusedUDPPort(b)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(repeats)*15*time.Second+10*time.Second)
	defer cancel()
	var earlySent chan struct{}
	if peer == "NSS" && feature.mode == "early" && !clientIsGo && serverIsGo {
		earlySent = make(chan struct{}, 1)
	}
	var listener *Listener
	var goServerDone <-chan peerBenchmarkServerResult
	if serverIsGo {
		config := f.config(b, feature.mode, true)
		var err error
		listener, err = Listen("udp4", "127.0.0.1:"+strconv.Itoa(port), config)
		if err != nil {
			b.Fatal(err)
		}
		defer listener.Close()
		goServerDone = runPeerBenchmarkGoServer(ctx, listener, feature, f, repeats, earlySent)
	}
	if !clientIsGo {
		benchmarkPeerProcessClient(b, peer, binary, f, feature, port, repeats, serverIsGo, goServerDone, earlySent, ctx)
		return
	}
	var serverCmd *exec.Cmd
	serverOutput := &lockedBuffer{changed: make(chan struct{}, 1)}
	if !serverIsGo {
		serverCmd = exec.CommandContext(ctx, binary, f.peerArgs("server", port, feature.mode, repeats)...) // #nosec G204 G702 -- opt-in local benchmark peer.
		if feature.mode == "keyupdate" {
			serverCmd.Env = append(os.Environ(), "DTLS13_BENCH_SERVER_KEYUPDATE=1")
		}
		serverCmd.Stdout, serverCmd.Stderr = serverOutput, serverOutput
		if err := serverCmd.Start(); err != nil {
			b.Fatal(err)
		}
		defer func() {
			if serverCmd != nil {
				_ = serverCmd.Process.Kill()
				_ = serverCmd.Wait()
			}
		}()
		if err := waitPeerBenchmarkLine(serverOutput, "READY 0", 5*time.Second); err != nil {
			b.Fatal(err)
		}
	}
	config := f.config(b, feature.mode, false)
	b.ResetTimer()
	elapsed, clientExporters := runPeerBenchmarkGoClient(b, ctx, peer, port, feature, f, config, serverOutput, repeats, !serverIsGo)
	b.StopTimer()
	unit := "go_ms/conn"
	if peerBenchmarkPair(feature.mode) {
		unit = "go_ms/pair"
	}
	b.ReportMetric(float64(elapsed)/float64(repeats)/float64(time.Millisecond), unit)
	if goServerDone != nil {
		select {
		case result := <-goServerDone:
			if result.err != nil {
				b.Fatalf("go server: %v", result.err)
			}
			if len(result.exporters) != len(clientExporters) {
				b.Fatal("go server exporter count mismatch")
			}
			for round, exporter := range clientExporters {
				if result.exporters[round] != exporter {
					b.Fatalf("go exporter mismatch at round %d", round)
				}
			}
		case <-ctx.Done():
			b.Fatal("go server timed out")
		}
	}
	if serverCmd != nil {
		if err := serverCmd.Wait(); err != nil {
			b.Fatalf("%s server: %v\n%s", peer, err, serverOutput.String())
		}
		serverCmd = nil
	}
	rounds := repeats
	if peerBenchmarkPair(feature.mode) {
		rounds *= 2
	}
	if !serverIsGo {
		for round := range rounds {
			if _, err := peerBenchmarkRound(serverOutput.String(), round); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func benchmarkPeerProcessClient(b *testing.B, peer, binary string, f *peerBenchmarkFixture, feature peerBenchmarkFeature, port, repeats int, serverIsGo bool, goServerDone <-chan peerBenchmarkServerResult, earlySent chan<- struct{}, ctx context.Context) {
	b.Helper()
	var elapsed time.Duration
	var clientExporters []string
	b.ResetTimer()
	for range repeats {
		b.StopTimer()
		serverOutput := &lockedBuffer{changed: make(chan struct{}, 1)}
		var serverCmd *exec.Cmd
		if !serverIsGo {
			serverCmd = exec.CommandContext(ctx, binary, f.peerArgs("server", port, feature.mode, 1)...) // #nosec G204 G702 -- opt-in local benchmark peer.
			serverCmd.Stdout, serverCmd.Stderr = serverOutput, serverOutput
			if err := serverCmd.Start(); err != nil {
				b.Fatal(err)
			}
			defer func(cmd *exec.Cmd) {
				if cmd.ProcessState == nil {
					_ = cmd.Process.Kill()
					_ = cmd.Wait()
				}
			}(serverCmd)
			if err := waitPeerBenchmarkLine(serverOutput, "READY 0", 5*time.Second); err != nil {
				b.Fatal(err)
			}
		}
		clientOutput := &lockedBuffer{changed: make(chan struct{}, 1)}
		cmd := exec.CommandContext(ctx, binary, f.peerArgs("client", port, feature.mode, 1)...) // #nosec G204 G702 -- opt-in local benchmark peer.
		cmd.Stdout, cmd.Stderr = clientOutput, clientOutput
		var stdin io.WriteCloser
		if serverCmd != nil && peerBenchmarkPair(feature.mode) {
			cmd.Env = append(os.Environ(), "DTLS13_BENCH_SYNC=1")
			var err error
			stdin, err = cmd.StdinPipe()
			if err != nil {
				b.Fatal(err)
			}
		}
		b.StartTimer()
		start := time.Now()
		var err error
		var peerWait time.Duration
		if earlySent != nil {
			err = cmd.Start()
			if err == nil {
				err = waitPeerBenchmarkLine(clientOutput, "EARLY_SENT 1", 5*time.Second)
				if err == nil {
					earlySent <- struct{}{}
				} else {
					_ = cmd.Process.Kill()
				}
				waitErr := cmd.Wait()
				if err == nil {
					err = waitErr
				}
			}
		} else if stdin == nil {
			err = cmd.Run()
		} else {
			err = cmd.Start()
			if err == nil {
				err = waitPeerBenchmarkLine(clientOutput, "WAIT 1", 5*time.Second)
				if err == nil {
					b.StopTimer()
					waitStart := time.Now()
					err = waitPeerBenchmarkLine(serverOutput, "READY 1", 5*time.Second)
					if err == nil {
						_, err = stdin.Write([]byte{1})
					}
					peerWait = time.Since(waitStart)
					b.StartTimer()
				}
				_ = stdin.Close()
				waitErr := cmd.Wait()
				if err == nil {
					err = waitErr
				}
			}
		}
		elapsed += time.Since(start) - peerWait
		b.StopTimer()
		if err != nil {
			var goServerErr error
			if goServerDone != nil {
				select {
				case result := <-goServerDone:
					goServerErr = result.err
				default:
				}
			}
			b.Fatalf("%s client: %v; go server: %v\n%s\nserver:\n%s", peer, err, goServerErr, clientOutput.String(), serverOutput.String())
		}
		if serverCmd != nil {
			if err := serverCmd.Wait(); err != nil {
				b.Fatalf("%s server: %v\n%s", peer, err, serverOutput.String())
			}
		}
		rounds := 1
		if peerBenchmarkPair(feature.mode) {
			rounds = 2
		}
		for round := range rounds {
			clientLine, err := peerBenchmarkRound(clientOutput.String(), round)
			if err != nil {
				b.Fatal(err)
			}
			clientExporter := strings.TrimSpace(strings.SplitN(clientLine, "EXPORTER=", 2)[1])
			clientExporters = append(clientExporters, clientExporter)
			if serverCmd != nil {
				serverLine, err := peerBenchmarkRound(serverOutput.String(), round)
				if err != nil {
					b.Fatal(err)
				}
				if !strings.Contains(serverLine, "EXPORTER="+clientExporter) {
					b.Fatalf("peer exporter mismatch at round %d", round)
				}
			}
		}
	}
	if goServerDone != nil {
		select {
		case result := <-goServerDone:
			if result.err != nil {
				b.Fatalf("go server: %v", result.err)
			}
			if len(result.exporters) != len(clientExporters) {
				b.Fatal("go server exporter count mismatch")
			}
			for round, exporter := range clientExporters {
				if result.exporters[round] != exporter {
					b.Fatalf("peer/go exporter mismatch at round %d", round)
				}
			}
		case <-ctx.Done():
			b.Fatal("go server timed out")
		}
	}
	unit := "peer_process_ms/conn"
	if peerBenchmarkPair(feature.mode) {
		unit = "peer_process_ms/pair"
	}
	b.ReportMetric(float64(elapsed)/float64(repeats)/float64(time.Millisecond), unit)
}

func waitPeerBenchmarkLine(output *lockedBuffer, want string, timeout time.Duration) error {
	timer := time.NewTimer(timeout)
	defer timer.Stop()
	for {
		if strings.Contains(output.String(), want) {
			return nil
		}
		select {
		case <-output.changed:
		case <-timer.C:
			if strings.Contains(output.String(), want) {
				return nil
			}
			return fmt.Errorf("peer did not emit %q: %s", want, output.String())
		}
	}
}

func peerBenchmarkRound(output string, round int) (string, error) {
	prefix := fmt.Sprintf("ROUND=%d ", round)
	for line := range strings.SplitSeq(output, "\n") {
		if strings.HasPrefix(line, prefix) && strings.Contains(line, "VERSION=fefc") && strings.Contains(line, "EXPORTER=") {
			return line, nil
		}
	}
	return "", fmt.Errorf("missing verified peer round %d: %s", round, output)
}

func runPeerBenchmarkGoServer(ctx context.Context, listener *Listener, feature peerBenchmarkFeature, f *peerBenchmarkFixture, repeats int, earlySent <-chan struct{}) <-chan peerBenchmarkServerResult {
	done := make(chan peerBenchmarkServerResult, 1)
	go func() {
		var result peerBenchmarkServerResult
		rounds := repeats
		if peerBenchmarkPair(feature.mode) {
			rounds *= 2
		}
		for round := range rounds {
			conn, err := listener.Accept()
			if err != nil {
				result.err = err
				done <- result
				return
			}
			if earlySent != nil && round%2 != 0 {
				select {
				case <-earlySent:
				case <-ctx.Done():
					_ = conn.Close()
					result.err = ctx.Err()
					done <- result
					return
				}
			}
			_ = conn.SetDeadline(time.Now().Add(8 * time.Second))
			var observedHRR, observedFragment atomic.Bool
			if feature.mode == "hrr" || feature.mode == "fragmented" {
				conn.conn = &observeRecordsConn{Conn: conn.conn, owner: conn, observe: func(r record) (bool, error) {
					if r.typ == recordTypeHandshake {
						fragments, err := parseHandshakeFragments(r.payload)
						if err != nil {
							return false, err
						}
						for _, fragment := range fragments {
							if fragment.length > uint32(len(fragment.body)) {
								observedFragment.Store(true)
							}
							if r.epoch == 0 && fragment.typ == handshakeTypeServerHello && fragment.offset == 0 && len(fragment.body) >= 34 && bytes.Equal(fragment.body[2:34], helloRetryRequestRandom[:]) {
								observedHRR.Store(true)
							}
						}
					}
					return false, nil
				}}
			}
			err = conn.Handshake()
			if err == nil && feature.mode == "hrr" && !observedHRR.Load() {
				err = errors.New("key-share HelloRetryRequest was not sent")
			}
			if err == nil && feature.mode == "fragmented" && !observedFragment.Load() {
				err = errors.New("Go server did not fragment its handshake")
			}
			if err == nil {
				exporter, exportErr := conn.ConnectionState().ExportKeyingMaterial("interop", nil, 32)
				if exportErr != nil {
					err = exportErr
				} else {
					result.exporters = append(result.exporters, hex.EncodeToString(exporter))
				}
			}
			if err == nil {
				err = servePeerBenchmarkGoConnection(conn, feature, f, round)
			}
			closeErr := conn.Close()
			if err == nil {
				err = closeErr
			}
			if err != nil {
				result.err = fmt.Errorf("round %d: %w", round, err)
				done <- result
				return
			}
		}
		done <- result
	}()
	return done
}

func servePeerBenchmarkGoConnection(conn *Conn, feature peerBenchmarkFeature, f *peerBenchmarkFixture, round int) error {
	resumed := peerBenchmarkPair(feature.mode) && round%2 != 0
	state := conn.ConnectionState()
	if err := validatePeerBenchmarkState(state, feature, f, true, resumed); err != nil {
		return err
	}
	var data [64]byte
	if feature.mode == "early" && resumed {
		n, info, err := conn.ReadDatagram(data[:])
		if err != nil || string(data[:n]) != "early" || !info.EarlyData {
			return fmt.Errorf("early data n=%d info=%v err=%w", n, info, err)
		}
	}
	n, _, err := conn.ReadDatagram(data[:])
	if err != nil || string(data[:n]) != "interop" {
		return fmt.Errorf("request n=%d err=%w", n, err)
	}
	if feature.mode == "pha" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = conn.RequestClientCertificate(ctx)
		cancel()
		if err != nil {
			return err
		}
		if len(conn.ConnectionState().VerifiedChains) == 0 {
			return errors.New("PHA did not authenticate client")
		}
	}
	if _, err := conn.WriteDatagram(data[:n]); err != nil {
		return err
	}
	if _, err := conn.WriteDatagram([]byte("done")); err != nil {
		return err
	}
	n, _, err = conn.ReadDatagram(data[:])
	if err != nil || string(data[:n]) != "done" {
		return fmt.Errorf("completion n=%d err=%w", n, err)
	}
	if feature.mode == "keyupdate" && conn.receivingTraffic.current < 4 {
		return errors.New("KeyUpdate did not advance receive epoch")
	}
	return nil
}

func runPeerBenchmarkGoClient(b *testing.B, ctx context.Context, peer string, port int, feature peerBenchmarkFeature, f *peerBenchmarkFixture, config *Config, peerOutput *lockedBuffer, repeats int, peerServer bool) (time.Duration, []string) {
	b.Helper()
	var elapsed time.Duration
	var exporters []string
	rounds := repeats
	if peerBenchmarkPair(feature.mode) {
		rounds *= 2
	}
	for round := range rounds {
		resumed := peerBenchmarkPair(feature.mode) && round%2 != 0
		if round > 0 && !resumed && peerBenchmarkPair(feature.mode) {
			config.ClientSessionCache = NewLRUClientSessionCache(2)
		}
		if peerServer {
			b.StopTimer()
			if err := waitPeerBenchmarkLine(peerOutput, fmt.Sprintf("READY %d", round), 5*time.Second); err != nil {
				b.Fatal(err)
			}
			b.StartTimer()
		}
		start := time.Now()
		raw, err := (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, "udp4", fmt.Sprintf("127.0.0.1:%d", port))
		if err != nil {
			b.Fatal(err)
		}
		var conn *Conn
		if feature.mode == "early" && resumed {
			conn = ClientEarly(raw, config)
			if n, e := conn.WriteDatagram([]byte("early")); e != nil || n != 5 {
				b.Fatalf("early data: n=%d err=%v", n, e)
			}
		} else {
			conn = Client(raw, config)
		}
		var observedHRR atomic.Bool
		if feature.mode == "hrr" || peer == "OPENSSL" && peerServer {
			conn.conn = &observeRecordsConn{Conn: conn.conn, owner: conn, observe: func(r record) (bool, error) {
				if feature.mode == "hrr" && r.epoch == 0 && r.typ == recordTypeHandshake {
					fragments, err := parseHandshakeFragments(r.payload)
					if err != nil {
						return false, err
					}
					for _, fragment := range fragments {
						if fragment.typ == handshakeTypeClientHello && fragment.messageSequence > 0 {
							observedHRR.Store(true)
						}
					}
				}
				return peer == "OPENSSL" && peerServer && r.typ == recordTypeACK && r.epoch == 2, nil
			}}
		}
		_ = conn.SetDeadline(time.Now().Add(8 * time.Second))
		if err = conn.HandshakeContext(ctx); err != nil {
			b.Fatalf("round %d handshake: %v\n%s", round, err, peerOutput.String())
		}
		if feature.mode == "hrr" && !observedHRR.Load() {
			b.Fatal("key-share HelloRetryRequest was not observed")
		}
		state := conn.ConnectionState()
		if err = validatePeerBenchmarkState(state, feature, f, false, resumed); err != nil {
			b.Fatal(err)
		}
		exporter, exportErr := state.ExportKeyingMaterial("interop", nil, 32)
		if exportErr != nil {
			b.Fatal(exportErr)
		}
		exporters = append(exporters, hex.EncodeToString(exporter))
		if feature.mode == "keyupdate" && !peerServer {
			if err = conn.SendKeyUpdate(false); err != nil {
				b.Fatal(err)
			}
		}
		if _, err = conn.WriteDatagram([]byte("interop")); err != nil {
			b.Fatal(err)
		}
		var data [64]byte
		n, _, err := conn.ReadDatagram(data[:])
		if err != nil || string(data[:n]) != "interop" {
			b.Fatalf("round %d echo n=%d err=%v", round, n, err)
		}
		if feature.mode == "keyupdate" {
			if peerServer && conn.receivingTraffic.current < 4 {
				b.Fatal("peer KeyUpdate did not advance receive epoch")
			}
			if !peerServer && conn.sendingTraffic.cipher.epoch < 4 {
				b.Fatal("KeyUpdate did not advance send epoch")
			}
		}
		if _, err = conn.WriteDatagram([]byte("done")); err != nil {
			b.Fatal(err)
		}
		n, _, err = conn.ReadDatagram(data[:])
		if err != nil || string(data[:n]) != "done" {
			b.Fatalf("round %d completion n=%d err=%v", round, n, err)
		}
		if peerServer {
			if e := waitPeerBenchmarkLine(peerOutput, fmt.Sprintf("ROUND=%d ", round), time.Second); e != nil {
				b.Fatal(e)
			}
			line, e := peerBenchmarkRound(peerOutput.String(), round)
			if e != nil || !strings.Contains(line, "EXPORTER="+exporters[round]) {
				b.Fatalf("round %d peer exporter mismatch: %v", round, e)
			}
		}
		if err = conn.Close(); err != nil {
			b.Fatal(err)
		}
		if !resumed && peerBenchmarkPair(feature.mode) {
			if err = waitForBenchmarkTicket(config); err != nil {
				b.Fatal(err)
			}
		}
		elapsed += time.Since(start)
	}
	return elapsed, exporters
}

func validatePeerBenchmarkState(state ConnectionState, feature peerBenchmarkFeature, f *peerBenchmarkFixture, server, resumed bool) error {
	if state.Version != VersionDTLS13 || state.DidResume != resumed {
		return fmt.Errorf("DTLS state: version=%#x resumed=%v, want resumed=%v", state.Version, state.DidResume, resumed)
	}
	wantCipher := uint16(TLS_AES_128_GCM_SHA256)
	switch feature.mode {
	case "sha384":
		wantCipher = TLS_AES_256_GCM_SHA384
	case "chacha":
		wantCipher = TLS_CHACHA20_POLY1305_SHA256
	case "ccm":
		wantCipher = TLS_AES_128_CCM_SHA256
	}
	if state.CipherSuite != wantCipher {
		return fmt.Errorf("cipher suite=%#x, want %#x", state.CipherSuite, wantCipher)
	}
	if feature.mode == "alpn" && state.NegotiatedProtocol != "coap" {
		return errors.New("ALPN was not negotiated")
	}
	if feature.mode == "ech" && !state.ECHAccepted {
		return errors.New("ECH was not accepted")
	}
	if feature.mode == "early" && resumed && state.EarlyData != EarlyDataAccepted {
		return fmt.Errorf("early data status=%v", state.EarlyData)
	}
	if (feature.mode == "external-psk" || feature.mode == "imported-psk") && !bytes.Equal(state.ExternalPSKIdentity(), []byte("interop-psk")) {
		return errors.New("external PSK was not selected")
	}
	if feature.mode == "rpk" && state.ServerCertificateType != CertificateTypeRawPublicKey {
		return errors.New("RPK was not negotiated")
	}
	if strings.HasPrefix(feature.mode, "dc") && !server && !bytes.Equal(state.PeerDelegatedCredential, f.dc.Credential) {
		return errors.New("delegated credential was not authenticated")
	}
	if feature.mode == "ocsp" && !server && !bytes.Equal(state.OCSPResponse, f.ocspResponse) {
		return errors.New("OCSP staple was not received")
	}
	if (feature.mode == "mutual" || feature.mode == "mutual-resume") && server && len(state.PeerCertificates) == 0 {
		return errors.New("client certificate was not authenticated")
	}
	return nil
}

func TestPeerBenchmarkRound(t *testing.T) {
	line, err := peerBenchmarkRound("READY 0\nROUND=1 VERSION=fefc CIPHER=1301 EXPORTER=abc\n", 1)
	if err != nil || !strings.Contains(line, "EXPORTER=abc") {
		t.Fatalf("line=%q err=%v", line, err)
	}
	if _, err := peerBenchmarkRound("ROUND=1 VERSION=fefb EXPORTER=abc\n", 1); err == nil {
		t.Fatal("accepted DTLS 1.2")
	}
}

func TestWaitPeerBenchmarkLine(t *testing.T) {
	output := &lockedBuffer{changed: make(chan struct{}, 1)}
	go func() { _, _ = output.Write([]byte("READY 1\n")) }()
	if err := waitPeerBenchmarkLine(output, "READY 1", time.Second); err != nil {
		t.Fatal(err)
	}
}
