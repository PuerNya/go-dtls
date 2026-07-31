package dtls13

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type wolfSSLInteropOptions struct {
	cipherName               string
	suites                   []uint16
	certificateCompression   bool
	externalPSK              *ExternalPSK
	args                     []string
	connections              int
	clientEarlyData          []byte
	disableServerEcho        bool
	requireClientCertificate bool
	loadClientCertificate    bool
	configure                func(*testing.T, string, *Config)
	wrapClientConn           func(net.Conn) net.Conn
	connected                func(*testing.T, *Conn, int)
	exchanged                func(*testing.T, *Conn, int)
	exchange                 func(*testing.T, *Conn, int)
	outputContains           []string
	unsupportedOutput        string
	dropClientFinalACK       bool
}

type lockedBuffer struct {
	mu sync.Mutex
	bytes.Buffer
}

type dropFinalACKConn struct {
	net.Conn
	mu             sync.Mutex
	dropNextRead   bool
	dropped        bool
	finishedWrites int
}

type finalACKProxy struct {
	conn           *net.UDPConn
	server         *net.UDPAddr
	done           chan struct{}
	mu             sync.Mutex
	client         *net.UDPAddr
	dropped        bool
	armed          bool
	finishedWrites int
	err            error
}

func newFinalACKProxy(t *testing.T, server *net.UDPAddr) *finalACKProxy {
	t.Helper()
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	proxy := &finalACKProxy{conn: conn, server: server, done: make(chan struct{})}
	go proxy.run()
	return proxy
}

func (p *finalACKProxy) run() {
	defer close(p.done)
	buffer := make([]byte, 64<<10)
	for {
		n, from, err := p.conn.ReadFromUDP(buffer)
		if err != nil {
			return
		}
		if from.Port == p.server.Port && from.IP.Equal(p.server.IP) {
			p.mu.Lock()
			client := p.client
			drop := p.armed && !p.dropped
			if drop {
				p.dropped = true
				p.armed = false
			}
			p.mu.Unlock()
			if client != nil && !drop {
				if _, err = p.conn.WriteToUDP(buffer[:n], client); err != nil {
					p.setError(err)
					return
				}
			}
			continue
		}

		p.mu.Lock()
		p.client = &net.UDPAddr{IP: append(net.IP(nil), from.IP...), Port: from.Port, Zone: from.Zone}
		if isEpoch2Datagram(buffer[:n]) {
			p.finishedWrites++
			if p.finishedWrites == 1 {
				p.armed = true
			}
		}
		p.mu.Unlock()
		if _, err = p.conn.WriteToUDP(buffer[:n], p.server); err != nil {
			p.setError(err)
			return
		}
	}
}

func (p *finalACKProxy) setError(err error) {
	p.mu.Lock()
	if p.err == nil {
		p.err = err
	}
	p.mu.Unlock()
}

func (p *finalACKProxy) waitForRetransmit(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		p.mu.Lock()
		finishedWrites := p.finishedWrites
		p.mu.Unlock()
		if finishedWrites >= 2 {
			return true
		}
		time.Sleep(5 * time.Millisecond)
	}
	return false
}

func (p *finalACKProxy) result() (dropped bool, finishedWrites int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.dropped, p.finishedWrites, p.err
}

func (p *finalACKProxy) Close() {
	_ = p.conn.Close()
	<-p.done
}

func (c *dropFinalACKConn) Write(p []byte) (int, error) {
	if isEpoch2Datagram(p) {
		c.mu.Lock()
		c.finishedWrites++
		if c.finishedWrites == 1 {
			c.dropNextRead = true
		}
		c.mu.Unlock()
	}
	return c.Conn.Write(p)
}

func (c *dropFinalACKConn) Read(p []byte) (int, error) {
	for {
		n, err := c.Conn.Read(p)
		c.mu.Lock()
		drop := err == nil && c.dropNextRead && !c.dropped
		if drop {
			c.dropNextRead = false
			c.dropped = true
		}
		c.mu.Unlock()
		if !drop {
			return n, err
		}
	}
}

func isEpoch2Datagram(datagram []byte) bool {
	return len(datagram) > 0 && datagram[0]&0xe0 == unifiedHeaderFixed && datagram[0]&unifiedHeaderEpochMask == 2
}

func (c *dropFinalACKConn) result() (dropped bool, finishedWrites int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.dropped, c.finishedWrites
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Buffer.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Buffer.String()
}

func wolfSSLPaths(t testing.TB) (root, server, client string) {
	t.Helper()
	root = os.Getenv("WOLFSSL_ROOT")
	if root == "" {
		t.Skip("set WOLFSSL_ROOT to run wolfSSL DTLS 1.3 interoperability tests")
	}
	server = filepath.Join(root, "build-zig", "examples", "server", "server.exe")
	client = filepath.Join(root, "build-zig", "examples", "client", "client.exe")
	for _, path := range []string{server, client} {
		if _, err := os.Stat(path); err != nil { // #nosec G703 -- path comes from local WOLFSSL_ROOT in this opt-in test.
			t.Fatalf("wolfSSL executable %q: %v", path, err)
		}
	}
	return root, server, client
}

func unusedUDPPort(t testing.TB) int {
	t.Helper()
	conn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
	if err != nil {
		t.Fatal(err)
	}
	port := conn.LocalAddr().(*net.UDPAddr).Port
	if err := conn.Close(); err != nil {
		t.Fatal(err)
	}
	return port
}

func TestInteropWolfSSLServer(t *testing.T) {
	testInteropWolfSSLServer(t, "", nil, false, nil)
}

func TestInteropWolfSSLServerECHGrease(t *testing.T) {
	testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
		configure: func(_ *testing.T, _ string, config *Config) {
			config.EncryptedClientHelloGrease = true
		},
		connected: func(t *testing.T, conn *Conn, _ int) {
			if conn.ConnectionState().ECHAccepted {
				t.Fatal("GREASE ECH was reported as accepted")
			}
		},
	})
}

func TestInteropWolfSSLServerAES128CCM(t *testing.T) {
	testInteropWolfSSLServer(t, "TLS13-AES128-CCM-SHA256", []uint16{TLS_AES_128_CCM_SHA256}, false, nil)
}

func TestInteropWolfSSLServerCertificateCompressionOffer(t *testing.T) {
	testInteropWolfSSLServer(t, "", nil, true, nil)
}

func TestInteropWolfSSLServerExternalPSK(t *testing.T) {
	testInteropWolfSSLServer(t, "TLS13-AES128-GCM-SHA256", []uint16{TLS_AES_128_GCM_SHA256}, false, wolfSSLExternalPSK(t))
}

func TestInteropWolfSSLServerHybridKeyExchange(t *testing.T) {
	for _, test := range wolfSSLHybridGroups {
		t.Run(test.name, func(t *testing.T) {
			if !test.wolfSSLServer {
				t.Skip("wolfSSL server does not complete this DTLS 1.3 hybrid handshake")
			}
			testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
				args: []string{"--pqc", test.name},
				configure: func(_ *testing.T, _ string, config *Config) {
					config.CurvePreferences = []tls.CurveID{test.group}
					config.MTU = 4096
				},
			})
		})
	}
}

func TestInteropWolfSSLServerConnectionID(t *testing.T) {
	nextCID := []byte("go-next")
	testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
		args: []string{"--cid", "wolf-srv"},
		configure: func(_ *testing.T, _ string, config *Config) {
			config.ConnectionID = []byte("go-cli")
			config.DisableReturnRoutabilityCheck = true
		},
		connected: func(t *testing.T, conn *Conn, _ int) {
			state := conn.ConnectionState()
			if len(state.LocalConnectionID) == 0 && len(state.PeerConnectionID) == 0 {
				return
			}
			if !bytes.Equal(state.LocalConnectionID, []byte("go-cli")) || !bytes.Equal(state.PeerConnectionID, []byte("wolf-srv")) {
				t.Fatalf("negotiated CIDs local=%q peer=%q", state.LocalConnectionID, state.PeerConnectionID)
			}
			if err := conn.SendNewConnectionIDs([][]byte{nextCID}, true); err != nil {
				t.Fatalf("send immediate CID update: %v", err)
			}
		},
		exchanged: func(t *testing.T, conn *Conn, _ int) {
			if state := conn.ConnectionState(); !bytes.Equal(state.LocalConnectionID, nextCID) {
				t.Fatalf("wolfSSL did not switch immediately to new CID: %q", state.LocalConnectionID)
			}
		},
	})
}

func TestInteropWolfSSLServerKeyUpdate(t *testing.T) {
	testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
		args: []string{"-U"},
		connected: func(t *testing.T, conn *Conn, _ int) {
			if err := conn.SendKeyUpdate(false); err != nil {
				t.Fatalf("send KeyUpdate: %v", err)
			}
		},
	})
}

func TestInteropWolfSSLServerPostHandshakeAuthentication(t *testing.T) {
	testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
		args:              []string{"-Q"},
		disableServerEcho: true,
		configure: func(t *testing.T, root string, config *Config) {
			config.PostHandshakeAuth = true
			config.Certificates = []tls.Certificate{wolfSSLCertificate(t, root, "client")}
		},
		exchanged: func(t *testing.T, conn *Conn, _ int) {
			conn.writeMu.Lock()
			requested := conn.hasClientAuthRequestSeq
			conn.writeMu.Unlock()
			if !requested {
				t.Fatal("wolfSSL server did not request post-handshake authentication")
			}
		},
	})
}

func TestInteropWolfSSLServerSessionResumption(t *testing.T) {
	testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
		args:        []string{"-r"},
		connections: 2,
		configure: func(_ *testing.T, _ string, config *Config) {
			config.ClientSessionCache = NewLRUClientSessionCache(2)
		},
		connected: requireResumptionOnSecondConnection,
	})
}

func TestInteropWolfSSLServerMutualTLSSessionResumption(t *testing.T) {
	testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
		args:                     []string{"-r"},
		connections:              2,
		requireClientCertificate: true,
		configure: func(t *testing.T, root string, config *Config) {
			config.Certificates = []tls.Certificate{wolfSSLCertificate(t, root, "client")}
			config.ClientSessionCache = NewLRUClientSessionCache(2)
		},
		connected: requireResumptionOnSecondConnection,
	})
}

func TestInteropWolfSSLServerEarlyData(t *testing.T) {
	testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
		args:            []string{"-r", "-0"},
		connections:     2,
		clientEarlyData: []byte("early go-dtls"),
		configure: func(_ *testing.T, _ string, config *Config) {
			config.ClientSessionCache = NewLRUClientSessionCache(2)
		},
		connected: requireResumptionOnSecondConnection,
	})
}

func TestInteropWolfSSLServerRetransmitsFinishedAfterDroppedFinalACK(t *testing.T) {
	var wire *dropFinalACKConn
	testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
		configure: func(_ *testing.T, _ string, config *Config) {
			config.FlightInterval = 20 * time.Millisecond
		},
		wrapClientConn: func(conn net.Conn) net.Conn {
			wire = &dropFinalACKConn{Conn: conn}
			return wire
		},
		exchanged: func(t *testing.T, _ *Conn, _ int) {
			dropped, finishedWrites := wire.result()
			if !dropped || finishedWrites < 2 {
				t.Fatalf("dropped final ACK=%v, client Finished writes=%d", dropped, finishedWrites)
			}
		},
	})
}

func testInteropWolfSSLServer(t *testing.T, cipherName string, suites []uint16, certificateCompression bool, externalPSK *ExternalPSK) {
	t.Helper()
	testInteropWolfSSLServerOptions(t, wolfSSLInteropOptions{
		cipherName: cipherName, suites: suites, certificateCompression: certificateCompression, externalPSK: externalPSK,
	})
}

func testInteropWolfSSLServerOptions(t *testing.T, options wolfSSLInteropOptions) {
	t.Helper()
	root, serverPath, _ := wolfSSLPaths(t)
	port := unusedUDPPort(t)
	args := []string{"-u", "-v", "4", "-p", strconv.Itoa(port)}
	if options.requireClientCertificate {
		args = append(args, "-F")
	} else {
		args = append(args, "-d")
	}
	if !options.disableServerEcho {
		args = append(args, "-e")
	}
	if options.externalPSK != nil {
		requireWolfSSLPSK(t, root, serverPath)
		args = append(args, "-s", "--onlyPskDheKe")
	}
	if options.cipherName != "" {
		args = append(args, "-l", options.cipherName)
	}
	args = append(args, options.args...)
	cmd := exec.Command(serverPath, args...) // #nosec G204 -- validated local WOLFSSL_ROOT executable in this opt-in test.
	cmd.Dir = root
	var output lockedBuffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		_ = cmd.Wait()
	}()
	time.Sleep(100 * time.Millisecond)
	if cmd.ProcessState != nil {
		t.Fatalf("wolfSSL server exited before handshake: %s", output.String())
	}

	dialer := &net.Dialer{Timeout: 5 * time.Second}
	config := &Config{
		InsecureSkipVerify: true, CipherSuites: options.suites, EnableCertificateCompression: options.certificateCompression,
		ExternalPSKs: externalPSKList(options.externalPSK), HandshakeTimeout: 5 * time.Second,
	}
	if options.configure != nil {
		options.configure(t, root, config)
	}
	connections := options.connections
	if connections == 0 {
		connections = 1
	}
	address := fmt.Sprintf("127.0.0.1:%d", port)
	for index := 0; index < connections; index++ {
		var conn *Conn
		var err error
		if index > 0 && len(options.clientEarlyData) > 0 {
			raw, dialErr := dialer.Dial("udp4", address)
			if dialErr != nil {
				t.Fatal(dialErr)
			}
			clientConfig := config.Clone()
			clientConfig.ServerName = "127.0.0.1"
			conn = Client(raw, clientConfig)
			n, earlyErr := conn.WriteEarlyData(options.clientEarlyData)
			if errors.Is(earlyErr, ErrEarlyDataRejected) {
				_ = conn.Close()
				t.Skip("wolfSSL server rejects 0-RTT after its DTLS HelloRetryRequest")
			}
			if earlyErr != nil || n != len(options.clientEarlyData) {
				_ = conn.Close()
				t.Fatalf("write early data: n=%d err=%v", n, earlyErr)
			}
		} else if options.wrapClientConn != nil {
			raw, dialErr := dialer.Dial("udp4", address)
			if dialErr != nil {
				t.Fatal(dialErr)
			}
			clientConfig := config.Clone()
			clientConfig.ServerName = "127.0.0.1"
			conn = Client(options.wrapClientConn(raw), clientConfig)
			err = conn.Handshake()
		} else {
			conn, err = DialWithDialer(dialer, "udp4", address, config)
		}
		if err != nil {
			t.Fatalf("Go client to wolfSSL 5.9.2 server: %v\n%s", err, output.String())
		}
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if options.connected != nil {
			options.connected(t, conn, index)
		}
		if options.exchange != nil {
			options.exchange(t, conn, index)
		} else {
			payload := []byte("go-dtls to wolfSSL")
			if _, err = conn.WriteDatagram(payload); err != nil {
				t.Fatal(err)
			}
			reply := make([]byte, 64)
			n, _, readErr := conn.ReadDatagram(reply)
			if readErr != nil {
				t.Fatalf("read wolfSSL echo: %v\n%s", readErr, output.String())
			}
			if options.disableServerEcho {
				if n == 0 {
					t.Fatal("wolfSSL server returned an empty application datagram")
				}
			} else if !bytes.Equal(reply[:n], payload) {
				t.Fatalf("wolfSSL echo = %q, want %q", reply[:n], payload)
			}
		}
		if options.exchanged != nil {
			options.exchanged(t, conn, index)
		}
		_ = conn.Close()
	}
	requireWolfSSLOutput(t, &output, options.outputContains)
}

func TestInteropWolfSSLClient(t *testing.T) {
	testInteropWolfSSLClient(t, "", nil, false, nil)
}

func TestInteropWolfSSLClientAES128CCM(t *testing.T) {
	testInteropWolfSSLClient(t, "TLS13-AES128-CCM-SHA256", []uint16{TLS_AES_128_CCM_SHA256}, false, nil)
}

func TestInteropWolfSSLClientCertificateCompressionFallback(t *testing.T) {
	testInteropWolfSSLClient(t, "", nil, true, nil)
}

func TestInteropWolfSSLClientExternalPSK(t *testing.T) {
	testInteropWolfSSLClient(t, "TLS13-AES128-GCM-SHA256", []uint16{TLS_AES_128_GCM_SHA256}, false, wolfSSLExternalPSK(t))
}

func TestInteropWolfSSLClientHybridKeyExchange(t *testing.T) {
	for _, test := range wolfSSLHybridGroups {
		t.Run(test.name, func(t *testing.T) {
			testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
				args: []string{"--pqc", test.name},
				configure: func(_ *testing.T, _ string, config *Config) {
					config.CurvePreferences = []tls.CurveID{test.group}
				},
			})
		})
	}
}

func TestInteropWolfSSLClientConnectionID(t *testing.T) {
	testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
		args: []string{"--cid", "wolf-cli"},
		configure: func(_ *testing.T, _ string, config *Config) {
			config.ConnectionID = []byte("go-srv")
			config.DisableReturnRoutabilityCheck = true
		},
		exchanged: func(t *testing.T, conn *Conn, _ int) {
			state := conn.ConnectionState()
			if !bytes.Equal(state.LocalConnectionID, []byte("go-srv")) || !bytes.Equal(state.PeerConnectionID, []byte("wolf-cli")) {
				t.Fatalf("negotiated CIDs local=%q peer=%q", state.LocalConnectionID, state.PeerConnectionID)
			}
		},
		outputContains: []string{"CID extension was negotiated"},
	})
}

func TestInteropWolfSSLClientKeyUpdate(t *testing.T) {
	testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
		args: []string{"-I"},
		connected: func(t *testing.T, conn *Conn, _ int) {
			if err := conn.SendKeyUpdate(false); err != nil {
				t.Fatalf("send KeyUpdate: %v", err)
			}
		},
	})
}

func TestInteropWolfSSLClientPostHandshakeAuthentication(t *testing.T) {
	testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
		args:                  []string{"-Q"},
		loadClientCertificate: true,
		configure: func(t *testing.T, root string, config *Config) {
			config.ClientAuth = tls.RequireAndVerifyClientCert
			config.ClientCAs = wolfSSLClientCAs(t, root)
		},
		connected: func(t *testing.T, conn *Conn, _ int) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := conn.RequestClientCertificate(ctx); err != nil {
				t.Fatalf("request post-handshake client certificate: %v", err)
			}
			if len(conn.ConnectionState().PeerCertificates) == 0 {
				t.Fatal("post-handshake authentication returned no client certificate")
			}
		},
	})
}

func TestInteropWolfSSLClientSessionResumption(t *testing.T) {
	testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
		args:        []string{"-r", "--waitTicket"},
		connections: 2,
		configure: func(_ *testing.T, _ string, config *Config) {
			config.SessionTicketsDisabled = false
		},
		connected:      requireResumptionOnSecondConnection,
		outputContains: []string{"reused session id"},
	})
}

func TestInteropWolfSSLClientMutualTLSSessionResumption(t *testing.T) {
	testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
		args:                  []string{"-r", "--waitTicket"},
		connections:           2,
		loadClientCertificate: true,
		configure: func(t *testing.T, root string, config *Config) {
			config.ClientAuth = tls.RequireAndVerifyClientCert
			config.ClientCAs = wolfSSLClientCAs(t, root)
			config.SessionTicketsDisabled = false
		},
		connected:         requireResumptionOnSecondConnection,
		unsupportedOutput: "wolfSSL_connect resume error -328, malformed buffer input error",
	})
}

func TestInteropWolfSSLClientEarlyData(t *testing.T) {
	testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
		args:        []string{"-r", "--waitTicket", "-0"},
		connections: 2,
		configure: func(_ *testing.T, _ string, config *Config) {
			config.SessionTicketsDisabled = false
			config.MaxEarlyData = 4096
			config.AllowEarlyDataWithoutCookie = true
		},
		connected: requireResumptionOnSecondConnection,
		exchange: func(t *testing.T, conn *Conn, index int) {
			want := 1
			if index == 1 {
				want = 2
			}
			early, application := false, false
			for range want {
				request := make([]byte, 64)
				n, _, err := conn.ReadDatagram(request)
				if err != nil {
					t.Fatalf("read wolfSSL application data: %v", err)
				}
				message := string(request[:n])
				early = early || strings.Contains(message, "A drop of info")
				application = application || strings.Contains(message, "wolfssl")
			}
			if !application || (index == 1 && !early) {
				t.Fatalf("connection %d early=%v application=%v", index+1, early, application)
			}
			if _, err := conn.WriteDatagram([]byte("hello from go-dtls")); err != nil {
				t.Fatal(err)
			}
		},
	})
}

func TestInteropWolfSSLClientRetransmitsFinishedAfterDroppedFinalACK(t *testing.T) {
	testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
		dropClientFinalACK: true,
	})
}

func testInteropWolfSSLClient(t *testing.T, cipherName string, suites []uint16, certificateCompression bool, externalPSK *ExternalPSK) {
	t.Helper()
	testInteropWolfSSLClientOptions(t, wolfSSLInteropOptions{
		cipherName: cipherName, suites: suites, certificateCompression: certificateCompression, externalPSK: externalPSK,
	})
}

func testInteropWolfSSLClientOptions(t *testing.T, options wolfSSLInteropOptions) {
	t.Helper()
	root, _, clientPath := wolfSSLPaths(t)
	certificate := wolfSSLCertificate(t, root, "server")
	config := &Config{
		Certificates: []tls.Certificate{certificate}, CipherSuites: options.suites,
		EnableCertificateCompression: options.certificateCompression, ExternalPSKs: externalPSKList(options.externalPSK),
		HandshakeTimeout: 5 * time.Second, SessionTicketsDisabled: true,
	}
	if options.configure != nil {
		options.configure(t, root, config)
	}
	listener, err := Listen("udp4", "127.0.0.1:0", config)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.UDPAddr).Port
	var proxy *finalACKProxy
	if options.dropClientFinalACK {
		proxy = newFinalACKProxy(t, listener.Addr().(*net.UDPAddr))
		defer proxy.Close()
		port = proxy.conn.LocalAddr().(*net.UDPAddr).Port
	}

	args := []string{"-u", "-v", "4", "-d", "-h", "127.0.0.1", "-p", strconv.Itoa(port)}
	if !options.loadClientCertificate {
		args = append(args, "-x")
	}
	if options.externalPSK != nil {
		requireWolfSSLPSK(t, root, clientPath)
		args = append(args, "-s", "--onlyPskDheKe", "--openssl-psk")
	}
	if options.cipherName != "" {
		args = append(args, "-l", options.cipherName)
	}
	args = append(args, options.args...)
	cmd := exec.Command(clientPath, args...) // #nosec G204 -- validated local WOLFSSL_ROOT executable in this opt-in test.
	cmd.Dir = root
	var output lockedBuffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}()

	connections := options.connections
	if connections == 0 {
		connections = 1
	}
	for index := 0; index < connections; index++ {
		accepted := make(chan struct {
			conn *Conn
			err  error
		}, 1)
		go func() {
			conn, acceptErr := listener.Accept()
			accepted <- struct {
				conn *Conn
				err  error
			}{conn: conn, err: acceptErr}
		}()
		var conn *Conn
		select {
		case result := <-accepted:
			conn, err = result.conn, result.err
		case processErr := <-done:
			text := output.String()
			if options.unsupportedOutput != "" && strings.Contains(text, options.unsupportedOutput) {
				t.Skipf("wolfSSL peer limitation: %s", options.unsupportedOutput)
			}
			t.Fatalf("wolfSSL client exited before connection %d: %v\n%s", index+1, processErr, text)
		case <-time.After(5 * time.Second):
			_ = listener.Close()
			_ = cmd.Process.Kill()
			t.Fatalf("accept wolfSSL client timed out\n%s", output.String())
		}
		if err != nil {
			_ = cmd.Process.Kill()
			t.Fatalf("accept wolfSSL client: %v\n%s", err, output.String())
		}
		_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
		if options.connected != nil {
			options.connected(t, conn, index)
		}
		if proxy != nil {
			if err = conn.Handshake(); err != nil {
				t.Fatalf("handshake through final-ACK proxy: %v", err)
			}
			if !proxy.waitForRetransmit(1500 * time.Millisecond) {
				dropped, finishedWrites, proxyErr := proxy.result()
				if dropped && finishedWrites == 1 && proxyErr == nil {
					t.Skip("wolfSSL client did not retransmit Finished after its final ACK was dropped")
				}
				t.Fatalf("final-ACK proxy dropped=%v Finished writes=%d err=%v", dropped, finishedWrites, proxyErr)
			}
		}
		if options.exchange != nil {
			options.exchange(t, conn, index)
		} else {
			request := make([]byte, 64)
			n, _, readErr := conn.ReadDatagram(request)
			if readErr != nil {
				_ = cmd.Process.Kill()
				t.Fatalf("read wolfSSL request: %v\n%s", readErr, output.String())
			}
			if !strings.Contains(string(request[:n]), "wolfssl") {
				t.Fatalf("unexpected wolfSSL request %q", request[:n])
			}
			if _, err = conn.WriteDatagram([]byte("hello from go-dtls")); err != nil {
				t.Fatal(err)
			}
		}
		if options.exchanged != nil {
			options.exchanged(t, conn, index)
		}
		_ = conn.Close()
	}
	select {
	case err = <-done:
		if err != nil {
			t.Fatalf("wolfSSL client failed: %v\n%s", err, output.String())
		}
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		t.Fatalf("wolfSSL client timed out\n%s", output.String())
	}
	requireWolfSSLOutput(t, &output, options.outputContains)
	if proxy != nil {
		dropped, finishedWrites, proxyErr := proxy.result()
		if proxyErr != nil || !dropped || finishedWrites < 2 {
			t.Fatalf("final-ACK proxy dropped=%v Finished writes=%d err=%v", dropped, finishedWrites, proxyErr)
		}
	}
}

func wolfSSLCertificate(t testing.TB, root, name string) tls.Certificate {
	t.Helper()
	certificate, err := tls.LoadX509KeyPair(filepath.Join(root, "certs", name+"-cert.pem"), filepath.Join(root, "certs", name+"-key.pem"))
	if err != nil {
		t.Fatal(err)
	}
	return certificate
}

func wolfSSLClientCAs(t testing.TB, root string) *x509.CertPool {
	t.Helper()
	pem, err := os.ReadFile(filepath.Join(root, "certs", "client-cert.pem")) // #nosec G304 -- root is the explicit local WOLFSSL_ROOT test fixture.
	if err != nil {
		t.Fatal(err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(pem) {
		t.Fatal("wolfSSL CA file contains no certificates")
	}
	return pool
}

func requireResumptionOnSecondConnection(t *testing.T, conn *Conn, index int) {
	t.Helper()
	if err := conn.Handshake(); err != nil {
		t.Fatalf("connection %d handshake: %v", index+1, err)
	}
	if conn.ConnectionState().DidResume != (index == 1) {
		t.Fatalf("connection %d DidResume=%v, want %v", index+1, conn.ConnectionState().DidResume, index == 1)
	}
}

func requireWolfSSLOutput(t *testing.T, output *lockedBuffer, required []string) {
	t.Helper()
	text := output.String()
	for _, substring := range required {
		if !strings.Contains(text, substring) {
			t.Fatalf("wolfSSL output does not contain %q:\n%s", substring, text)
		}
	}
}

func wolfSSLExternalPSK(t testing.TB) *ExternalPSK {
	t.Helper()
	key := bytes.Repeat([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef}, 4)
	psk, err := NewDirectExternalPSK([]byte("Client_identity"), key, crypto.SHA256)
	if err != nil {
		t.Fatal(err)
	}
	return psk
}

func externalPSKList(psk *ExternalPSK) []*ExternalPSK {
	if psk == nil {
		return nil
	}
	return []*ExternalPSK{psk}
}

var wolfSSLHybridGroups = []struct {
	name          string
	group         tls.CurveID
	wolfSSLServer bool
}{
	{"X25519MLKEM768", tls.X25519MLKEM768, true},
	{"SecP256r1MLKEM768", tls.SecP256r1MLKEM768, true},
	{"SecP384r1MLKEM1024", tls.SecP384r1MLKEM1024, false},
}

const wolfSSLBenchmarkBatch = 20

type wolfSSLBenchmarkOperation uint8

const (
	wolfSSLBenchmarkHandshake wolfSSLBenchmarkOperation = iota
	wolfSSLBenchmarkApplicationData
	wolfSSLBenchmarkKeyUpdate
	wolfSSLBenchmarkPHA
)

type wolfSSLBenchmarkFeature struct {
	name                            string
	configs                         func() (*Config, *Config)
	suite                           uint16
	operation                       wolfSSLBenchmarkOperation
	batch                           int
	resume                          bool
	earlyData                       bool
	externalPSK                     bool
	connectionID                    bool
	mutualTLS                       bool
	wolfClientProcess               bool
	loadWolfClientCertificate       bool
	requireWolfClientCertificate    bool
	wolfClientArgs                  []string
	wolfServerArgs                  []string
	wolfClientOutput                []string
	wolfServerOutput                []string
	goClientWolfServerUnsupported   string
	wolfClientGoServerUnsupported   string
	wolfClientWolfServerUnsupported string
}

func BenchmarkWolfSSLFeatureRealUDP(b *testing.B) {
	root, serverPath, clientPath := wolfSSLPaths(b)
	serverCertificate := wolfSSLCertificate(b, root, "server")
	clientCertificate := wolfSSLCertificate(b, root, "client")
	clientCAs := wolfSSLClientCAs(b, root)
	externalPSK := wolfSSLExternalPSK(b)
	var ticketKey [32]byte
	copy(ticketKey[:], bytes.Repeat([]byte{0x5c}, len(ticketKey)))

	baseConfigs := func() (*Config, *Config) {
		return &Config{
				InsecureSkipVerify: true, ServerName: "127.0.0.1", CipherSuites: []uint16{TLS_AES_128_GCM_SHA256},
				SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second,
			}, &Config{
				Certificates: []tls.Certificate{serverCertificate}, CipherSuites: []uint16{TLS_AES_128_GCM_SHA256},
				SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second,
			}
	}
	mutualTLSConfigs := func(resume bool) func() (*Config, *Config) {
		return func() (*Config, *Config) {
			client := &Config{
				InsecureSkipVerify: true, ServerName: "127.0.0.1", Certificates: []tls.Certificate{clientCertificate},
				CipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, SessionTicketsDisabled: !resume, HandshakeTimeout: 5 * time.Second,
			}
			server := &Config{
				Certificates: []tls.Certificate{serverCertificate}, ClientAuth: tls.RequireAndVerifyClientCert, ClientCAs: clientCAs,
				CipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, SessionTicketsDisabled: !resume,
				SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour, HandshakeTimeout: 5 * time.Second,
			}
			if resume {
				client.ClientSessionCache = NewLRUClientSessionCache(1)
			}
			return client, server
		}
	}
	resumptionConfigs := func(earlyData bool) func() (*Config, *Config) {
		return func() (*Config, *Config) {
			client := &Config{
				InsecureSkipVerify: true, ServerName: "127.0.0.1", CipherSuites: []uint16{TLS_AES_128_GCM_SHA256},
				ClientSessionCache: NewLRUClientSessionCache(1), HandshakeTimeout: 5 * time.Second,
			}
			server := &Config{
				Certificates: []tls.Certificate{serverCertificate}, CipherSuites: []uint16{TLS_AES_128_GCM_SHA256},
				SessionTicketKey: ticketKey, SessionTicketLifetime: time.Hour, HandshakeTimeout: 5 * time.Second,
			}
			if earlyData {
				server.MaxEarlyData = 4096
				server.AllowEarlyDataWithoutCookie = true
				server.EarlyDataReplayCache = NewLRUEarlyDataReplayCache(wolfSSLBenchmarkBatch * 4)
			}
			return client, server
		}
	}

	features := []wolfSSLBenchmarkFeature{
		{name: "CertificateAES128GCM", configs: baseConfigs, suite: TLS_AES_128_GCM_SHA256, wolfClientArgs: []string{"-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"-l", "TLS13-AES128-GCM-SHA256"}},
		{name: "ApplicationDataRoundTrip", configs: baseConfigs, suite: TLS_AES_128_GCM_SHA256, operation: wolfSSLBenchmarkApplicationData, wolfClientProcess: true, wolfClientArgs: []string{"-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"-l", "TLS13-AES128-GCM-SHA256"}},
		{name: "MutualTLS", configs: mutualTLSConfigs(false), suite: TLS_AES_128_GCM_SHA256, mutualTLS: true, loadWolfClientCertificate: true, requireWolfClientCertificate: true, wolfClientArgs: []string{"-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"-l", "TLS13-AES128-GCM-SHA256"}},
		{name: "AES128CCM", suite: TLS_AES_128_CCM_SHA256, configs: func() (*Config, *Config) {
			client, server := baseConfigs()
			client.CipherSuites, server.CipherSuites = []uint16{TLS_AES_128_CCM_SHA256}, []uint16{TLS_AES_128_CCM_SHA256}
			return client, server
		}, wolfClientArgs: []string{"-l", "TLS13-AES128-CCM-SHA256"}, wolfServerArgs: []string{"-l", "TLS13-AES128-CCM-SHA256"}},
		{name: "ExternalPSK", suite: TLS_AES_128_GCM_SHA256, externalPSK: true, configs: func() (*Config, *Config) {
			config := &Config{ExternalPSKs: []*ExternalPSK{externalPSK}, CipherSuites: []uint16{TLS_AES_128_GCM_SHA256}, SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second}
			return config.Clone(), config.Clone()
		}, wolfClientArgs: []string{"-s", "--onlyPskDheKe", "--openssl-psk", "-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"-s", "--onlyPskDheKe", "-l", "TLS13-AES128-GCM-SHA256"}},
		{name: "ConnectionID", suite: TLS_AES_128_GCM_SHA256, connectionID: true, wolfClientProcess: true, configs: func() (*Config, *Config) {
			client, server := baseConfigs()
			client.ConnectionID, server.ConnectionID = []byte("go-cli"), []byte("go-srv")
			client.DisableReturnRoutabilityCheck, server.DisableReturnRoutabilityCheck = true, true
			return client, server
		}, wolfClientArgs: []string{"--cid", "wolf-cli", "-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"--cid", "wolf-srv", "-l", "TLS13-AES128-GCM-SHA256"}, wolfServerOutput: []string{"CID extension was negotiated"}},
		{name: "KeyUpdate", configs: baseConfigs, suite: TLS_AES_128_GCM_SHA256, operation: wolfSSLBenchmarkKeyUpdate, wolfClientProcess: true, wolfClientArgs: []string{"-I", "-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"-U", "-l", "TLS13-AES128-GCM-SHA256"}},
		{name: "PostHandshakeAuthentication", suite: TLS_AES_128_GCM_SHA256, operation: wolfSSLBenchmarkPHA, wolfClientProcess: true, loadWolfClientCertificate: true, configs: func() (*Config, *Config) {
			client, server := baseConfigs()
			client.PostHandshakeAuth = true
			client.Certificates = []tls.Certificate{clientCertificate}
			server.ClientAuth = tls.RequireAndVerifyClientCert
			server.ClientCAs = clientCAs
			return client, server
		}, wolfClientArgs: []string{"-Q", "-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"-Q", "-l", "TLS13-AES128-GCM-SHA256"}, wolfServerOutput: []string{"Successfully requested post-hs certificate"}},
		{name: "SessionResumption", configs: resumptionConfigs(false), suite: TLS_AES_128_GCM_SHA256, resume: true, wolfClientProcess: true, wolfClientArgs: []string{"-r", "--waitTicket", "-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"-r", "-l", "TLS13-AES128-GCM-SHA256"}, wolfClientOutput: []string{"reused session id"}},
		{name: "MutualTLSSessionResumption", configs: mutualTLSConfigs(true), suite: TLS_AES_128_GCM_SHA256, resume: true, mutualTLS: true, wolfClientProcess: true, loadWolfClientCertificate: true, requireWolfClientCertificate: true, wolfClientArgs: []string{"-r", "--waitTicket", "-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"-r", "-l", "TLS13-AES128-GCM-SHA256"}, wolfClientOutput: []string{"reused session id"}, wolfClientGoServerUnsupported: "wolfSSL client cannot parse the go-dtls mTLS session ticket"},
		{name: "EarlyData", configs: resumptionConfigs(true), suite: TLS_AES_128_GCM_SHA256, batch: 1, resume: true, earlyData: true, wolfClientProcess: true, wolfClientArgs: []string{"-r", "--waitTicket", "-0", "-l", "TLS13-AES128-GCM-SHA256"}, wolfServerArgs: []string{"-r", "-0", "-l", "TLS13-AES128-GCM-SHA256"}, wolfClientOutput: []string{"reused session id"}, goClientWolfServerUnsupported: "wolfSSL server rejects go-dtls 0-RTT after HelloRetryRequest", wolfClientWolfServerUnsupported: "wolfSSL server rejects wolfSSL client 0-RTT after HelloRetryRequest"},
	}
	for _, feature := range features {
		b.Run(feature.name, func(b *testing.B) {
			if feature.externalPSK {
				requireWolfSSLPSK(b, root, serverPath)
				requireWolfSSLPSK(b, root, clientPath)
			}
			benchmarkWolfSSLFeatureDirections(b, root, serverPath, clientPath, feature)
		})
	}
}

func BenchmarkHybridKeyExchangeRealUDP(b *testing.B) {
	root, serverPath, clientPath := wolfSSLPaths(b)
	certificate := wolfSSLCertificate(b, root, "server")
	for _, test := range wolfSSLHybridGroups {
		feature := wolfSSLBenchmarkFeature{
			name: test.name,
			configs: func() (*Config, *Config) {
				return &Config{
						InsecureSkipVerify: true, ServerName: "127.0.0.1", CurvePreferences: []tls.CurveID{test.group}, MTU: 4096,
						SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second,
					}, &Config{
						Certificates: []tls.Certificate{certificate}, CurvePreferences: []tls.CurveID{test.group},
						SessionTicketsDisabled: true, HandshakeTimeout: 5 * time.Second,
					}
			},
			wolfClientArgs: []string{"--pqc", test.name}, wolfServerArgs: []string{"--pqc", test.name},
		}
		if !test.wolfSSLServer {
			feature.goClientWolfServerUnsupported = "wolfSSL server does not complete this DTLS 1.3 hybrid handshake"
		}
		b.Run(test.name, func(b *testing.B) {
			benchmarkWolfSSLFeatureDirections(b, root, serverPath, clientPath, feature)
		})
	}
}

func benchmarkWolfSSLFeatureDirections(b *testing.B, root, serverPath, clientPath string, feature wolfSSLBenchmarkFeature) {
	b.Helper()
	batch := feature.batch
	if batch == 0 {
		batch = wolfSSLBenchmarkBatch
	}
	b.Run("GoClient/GoServer", func(b *testing.B) {
		connections := b.N * batch
		clientConfig, serverConfig := feature.configs()
		listener, done := startGoBenchmarkServer(b, serverConfig, feature, connections, false)
		defer listener.Close()
		benchmarkGoClient(b, listener.Addr().String(), clientConfig, feature, connections, false)
		waitForGoBenchmarkServer(b, done)
	})
	b.Run("GoClient/WolfSSLServer", func(b *testing.B) {
		if feature.goClientWolfServerUnsupported != "" {
			b.Skip(feature.goClientWolfServerUnsupported)
		}
		connections := b.N * batch
		clientConfig, _ := feature.configs()
		port := unusedUDPPort(b)
		server, output, done := startWolfSSLBenchmarkServer(b, root, serverPath, port, connections, feature)
		defer func() { _ = server.Process.Kill() }()
		benchmarkGoClient(b, fmt.Sprintf("127.0.0.1:%d", port), clientConfig, feature, connections, true)
		waitForWolfSSLBenchmarkServer(b, done, output, feature.wolfServerOutput)
	})
	b.Run("WolfSSLClient/GoServer", func(b *testing.B) {
		if feature.wolfClientGoServerUnsupported != "" {
			b.Skip(feature.wolfClientGoServerUnsupported)
		}
		connections := b.N * batch
		_, serverConfig := feature.configs()
		listener, done := startGoBenchmarkServer(b, serverConfig, feature, connections, true)
		defer listener.Close()
		benchmarkWolfSSLClient(b, root, clientPath, listener.Addr().(*net.UDPAddr).Port, connections, feature)
		waitForGoBenchmarkServer(b, done)
	})
	b.Run("WolfSSLClient/WolfSSLServer", func(b *testing.B) {
		if feature.wolfClientWolfServerUnsupported != "" {
			b.Skip(feature.wolfClientWolfServerUnsupported)
		}
		connections := b.N * batch
		port := unusedUDPPort(b)
		server, output, done := startWolfSSLBenchmarkServer(b, root, serverPath, port, connections, feature)
		defer func() { _ = server.Process.Kill() }()
		benchmarkWolfSSLClient(b, root, clientPath, port, connections, feature)
		waitForWolfSSLBenchmarkServer(b, done, output, feature.wolfServerOutput)
	})
}

func benchmarkGoServerConnections(feature wolfSSLBenchmarkFeature, connections int) int {
	if feature.resume {
		return 2 * connections
	}
	return connections
}

func startGoBenchmarkServer(b *testing.B, config *Config, feature wolfSSLBenchmarkFeature, connections int, wolfClient bool) (*Listener, <-chan error) {
	b.Helper()
	listener, err := Listen("udp4", "127.0.0.1:0", config)
	if err != nil {
		b.Fatal(err)
	}
	total := benchmarkGoServerConnections(feature, connections)
	done := make(chan error, 1)
	go func() {
		for index := range total {
			conn, acceptErr := listener.Accept()
			if acceptErr != nil {
				done <- acceptErr
				return
			}
			_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
			serveErr := conn.Handshake()
			if serveErr == nil {
				serveErr = serveGoBenchmarkConnection(conn, feature, index, connections, wolfClient)
			}
			closeErr := conn.Close()
			if serveErr != nil {
				b.Logf("go-dtls benchmark server connection %d: %v", index+1, serveErr)
				done <- serveErr
				return
			}
			if closeErr != nil {
				done <- closeErr
				return
			}
		}
		done <- nil
	}()
	return listener, done
}

func serveGoBenchmarkConnection(conn *Conn, feature wolfSSLBenchmarkFeature, index, connections int, wolfClient bool) error {
	if feature.resume {
		full := index%2 == 0
		if wolfClient && !feature.wolfClientProcess {
			full = index < connections
		}
		if feature.earlyData && !full {
			buffer := make([]byte, 256)
			n, _, err := conn.ReadDatagram(buffer)
			if err != nil || n == 0 {
				return fmt.Errorf("read early data: n=%d err=%w", n, err)
			}
			want := "benchmark early data"
			if wolfClient {
				want = "A drop of info"
			}
			if !strings.Contains(string(buffer[:n]), want) {
				return fmt.Errorf("early data %q does not contain %q", buffer[:n], want)
			}
			if !wolfClient {
				if _, _, err = conn.ReadDatagram(buffer); err != nil {
					return fmt.Errorf("read after early data: %w", err)
				}
				if _, err = conn.WriteDatagram([]byte("early data response")); err != nil {
					return fmt.Errorf("write after early data: %w", err)
				}
			}
		}
		if wolfClient && (full || feature.wolfClientProcess) {
			if _, _, err := conn.ReadDatagram(make([]byte, 256)); err != nil {
				return fmt.Errorf("read ticket-producing request: %w", err)
			}
			if _, err := conn.WriteDatagram([]byte("benchmark response")); err != nil {
				return fmt.Errorf("write ticket-producing response: %w", err)
			}
		}
		if err := waitForBenchmarkClose(conn); err != nil {
			return err
		}
		return validateGoBenchmarkConnection(conn, feature, true, !full)
	}

	switch feature.operation {
	case wolfSSLBenchmarkApplicationData:
		if _, _, err := conn.ReadDatagram(make([]byte, 256)); err != nil {
			return fmt.Errorf("read application request: %w", err)
		}
		if _, err := conn.WriteDatagram([]byte("application response")); err != nil {
			return fmt.Errorf("write application response: %w", err)
		}
	case wolfSSLBenchmarkKeyUpdate:
		if _, _, err := conn.ReadDatagram(make([]byte, 256)); err != nil {
			return fmt.Errorf("read after KeyUpdate: %w", err)
		}
		if conn.receivingTraffic == nil || conn.receivingTraffic.current < 4 {
			return errors.New("peer KeyUpdate did not advance the receive epoch")
		}
		if _, err := conn.WriteDatagram([]byte("key update response")); err != nil {
			return err
		}
		if _, _, err := conn.ReadDatagram(make([]byte, 256)); err != nil {
			return fmt.Errorf("read KeyUpdate confirmation: %w", err)
		}
	case wolfSSLBenchmarkPHA:
		if _, _, err := conn.ReadDatagram(make([]byte, 256)); err != nil {
			return fmt.Errorf("read before PHA: %w", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err := conn.RequestClientCertificate(ctx)
		cancel()
		if err != nil {
			return fmt.Errorf("request post-handshake certificate: %w", err)
		}
		if _, err = conn.WriteDatagram([]byte("post-handshake auth response")); err != nil {
			return err
		}
		if _, _, err = conn.ReadDatagram(make([]byte, 256)); err != nil {
			return fmt.Errorf("read PHA confirmation: %w", err)
		}
	default:
		if wolfClient && feature.wolfClientProcess {
			if _, _, err := conn.ReadDatagram(make([]byte, 256)); err != nil {
				return fmt.Errorf("read wolfSSL application request: %w", err)
			}
			if _, err := conn.WriteDatagram([]byte("benchmark response")); err != nil {
				return fmt.Errorf("write wolfSSL application response: %w", err)
			}
		}
	}
	if err := waitForBenchmarkClose(conn); err != nil {
		return err
	}
	return validateGoBenchmarkConnection(conn, feature, true, false)
}

func waitForBenchmarkClose(conn *Conn) error {
	_, _, err := conn.ReadDatagram(nil)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func benchmarkGoClient(b *testing.B, address string, config *Config, feature wolfSSLBenchmarkFeature, connections int, wolfServer bool) {
	b.Helper()
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	skip := connections / 10
	var elapsed time.Duration
	b.ResetTimer()
	for index := range connections {
		clientConfig := config
		if feature.resume {
			clientConfig = config.Clone()
			clientConfig.ClientSessionCache = NewLRUClientSessionCache(1)
			full, err := DialWithDialer(dialer, "udp4", address, clientConfig)
			if err != nil {
				b.Fatal(err)
			}
			_ = full.SetDeadline(time.Now().Add(5 * time.Second))
			if wolfServer {
				if _, err = full.WriteDatagram([]byte("benchmark ticket request")); err == nil {
					_, _, err = full.ReadDatagram(make([]byte, 256))
				}
			}
			if err == nil {
				err = waitForBenchmarkTicket(clientConfig)
			}
			if closeErr := full.Close(); err == nil {
				err = closeErr
			}
			if err != nil {
				b.Fatalf("prepare resumed connection: %v", err)
			}
		}

		start := time.Now()
		var conn *Conn
		var err error
		if feature.earlyData {
			var raw net.Conn
			raw, err = dialer.Dial("udp4", address)
			if err == nil {
				conn = Client(raw, clientConfig)
				var n int
				n, err = conn.WriteEarlyData([]byte("benchmark early data"))
				if err == nil && n == 0 {
					err = errors.New("WriteEarlyData wrote no data")
				}
			}
		} else {
			conn, err = DialWithDialer(dialer, "udp4", address, clientConfig)
		}
		if err == nil {
			_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
			err = runGoBenchmarkClientOperation(conn, feature)
		}
		keyUpdateAdvanced := true
		if err == nil && feature.operation == wolfSSLBenchmarkKeyUpdate {
			keyUpdateAdvanced = conn.sendingTraffic != nil && conn.sendingTraffic.cipher.epoch >= 4
		}
		if err == nil {
			err = conn.Close()
		}
		duration := time.Since(start)
		if err != nil {
			b.Fatal(err)
		}
		if !keyUpdateAdvanced {
			b.Fatal("KeyUpdate did not advance the sending epoch")
		}
		if err = validateGoBenchmarkConnection(conn, feature, false, feature.resume); err != nil {
			b.Fatal(err)
		}
		if index >= skip {
			elapsed += duration
		}
	}
	b.StopTimer()
	b.ReportMetric(float64(elapsed)/float64(connections-skip)/float64(time.Millisecond), "go_ms/conn")
}

func runGoBenchmarkClientOperation(conn *Conn, feature wolfSSLBenchmarkFeature) error {
	if feature.earlyData {
		if _, err := conn.WriteDatagram([]byte("early data 1-RTT request")); err != nil {
			return err
		}
		_, _, err := conn.ReadDatagram(make([]byte, 256))
		return err
	}
	switch feature.operation {
	case wolfSSLBenchmarkApplicationData:
		if _, err := conn.WriteDatagram([]byte("application request")); err != nil {
			return err
		}
		_, _, err := conn.ReadDatagram(make([]byte, 256))
		return err
	case wolfSSLBenchmarkKeyUpdate:
		if err := conn.SendKeyUpdate(false); err != nil {
			return err
		}
		if _, err := conn.WriteDatagram([]byte("key update request")); err != nil {
			return err
		}
		if _, _, err := conn.ReadDatagram(make([]byte, 256)); err != nil {
			return err
		}
		_, err := conn.WriteDatagram([]byte("key update confirmation"))
		return err
	case wolfSSLBenchmarkPHA:
		if _, err := conn.WriteDatagram([]byte("post-handshake auth request")); err != nil {
			return err
		}
		if _, _, err := conn.ReadDatagram(make([]byte, 256)); err != nil {
			return err
		}
		_, err := conn.WriteDatagram([]byte("post-handshake auth confirmation"))
		return err
	default:
		return nil
	}
}

func waitForBenchmarkTicket(config *Config) error {
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		if _, ok := config.ClientSessionCache.Get(config.ServerName); ok {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return errors.New("session ticket was not cached")
}

func validateGoBenchmarkConnection(conn *Conn, feature wolfSSLBenchmarkFeature, server, resumed bool) error {
	state := conn.ConnectionState()
	if state.DidResume != resumed {
		return fmt.Errorf("DidResume=%v, want %v", state.DidResume, resumed)
	}
	if feature.suite != 0 && state.CipherSuite != feature.suite {
		return fmt.Errorf("cipher suite=%#x, want %#x", state.CipherSuite, feature.suite)
	}
	if feature.externalPSK && len(state.ExternalPSKIdentity()) == 0 {
		return errors.New("external PSK was not selected")
	}
	if feature.connectionID && (len(state.LocalConnectionID) == 0 || len(state.PeerConnectionID) == 0) {
		return errors.New("Connection ID was not negotiated")
	}
	if feature.mutualTLS && server && len(state.PeerCertificates) == 0 {
		return errors.New("client certificate was not authenticated")
	}
	if feature.operation == wolfSSLBenchmarkPHA && server && len(state.PeerCertificates) == 0 {
		return errors.New("post-handshake client certificate was not authenticated")
	}
	return nil
}

func startWolfSSLBenchmarkServer(b *testing.B, root, serverPath string, port, connections int, feature wolfSSLBenchmarkFeature) (*exec.Cmd, *lockedBuffer, <-chan error) {
	b.Helper()
	args := []string{"-u", "-v", "4", "-p", strconv.Itoa(port), "-C", strconv.Itoa(connections)}
	if feature.requireWolfClientCertificate {
		args = append(args, "-F")
	} else {
		args = append(args, "-d")
	}
	args = append(args, feature.wolfServerArgs...)
	cmd := exec.Command(serverPath, args...) // #nosec G204 -- validated local WOLFSSL_ROOT executable in this opt-in benchmark.
	cmd.Dir = root
	output := new(lockedBuffer)
	cmd.Stdout, cmd.Stderr = output, output
	if err := cmd.Start(); err != nil {
		b.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	time.Sleep(100 * time.Millisecond)
	select {
	case err := <-done:
		b.Fatalf("wolfSSL server exited before benchmark: %v\n%s", err, output.String())
	default:
	}
	return cmd, output, done
}

func benchmarkWolfSSLClient(b *testing.B, root, clientPath string, port, connections int, feature wolfSSLBenchmarkFeature) {
	b.Helper()
	baseArgs := []string{"-u", "-v", "4", "-d", "-h", "127.0.0.1", "-p", strconv.Itoa(port)}
	if !feature.loadWolfClientCertificate {
		baseArgs = append(baseArgs, "-x")
	}
	baseArgs = append(baseArgs, feature.wolfClientArgs...)
	if !feature.wolfClientProcess {
		args := append(append([]string(nil), baseArgs...), "-b", strconv.Itoa(connections))
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, clientPath, args...) // #nosec G204 -- validated local WOLFSSL_ROOT executable in this opt-in benchmark.
		cmd.Dir = root
		var output lockedBuffer
		cmd.Stdout, cmd.Stderr = &output, &output
		b.ResetTimer()
		err := cmd.Run()
		b.StopTimer()
		if err != nil {
			b.Fatalf("wolfSSL client benchmark failed: %v\n%s", err, output.String())
		}
		if err = requireBenchmarkOutput(output.String(), feature.wolfClientOutput); err != nil {
			b.Fatal(err)
		}
		average, err := parseWolfSSLBenchmarkAverage(output.String(), feature.resume)
		if err != nil {
			b.Fatal(err)
		}
		b.ReportMetric(average, "wolfssl_ms/conn")
		return
	}

	skip := connections / 10
	var elapsed time.Duration
	b.ResetTimer()
	for index := range connections {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		cmd := exec.CommandContext(ctx, clientPath, baseArgs...) // #nosec G204 -- validated local WOLFSSL_ROOT executable in this opt-in benchmark.
		cmd.Dir = root
		var output lockedBuffer
		cmd.Stdout, cmd.Stderr = &output, &output
		start := time.Now()
		err := cmd.Run()
		duration := time.Since(start)
		cancel()
		if err != nil {
			b.Fatalf("wolfSSL client process benchmark failed: %v\n%s", err, output.String())
		}
		if err = requireBenchmarkOutput(output.String(), feature.wolfClientOutput); err != nil {
			b.Fatal(err)
		}
		if index >= skip {
			elapsed += duration
		}
	}
	b.StopTimer()
	metric := "wolf_process_ms/conn"
	if feature.resume {
		metric = "wolf_process_ms/pair"
	}
	b.ReportMetric(float64(elapsed)/float64(connections-skip)/float64(time.Millisecond), metric)
}

func parseWolfSSLBenchmarkAverage(output string, resume bool) (float64, error) {
	prefix := "wolfSSL_connect avg took:"
	if resume {
		prefix = "wolfSSL_resume  avg took:"
	}
	index := strings.LastIndex(output, prefix)
	if index < 0 {
		return 0, fmt.Errorf("wolfSSL benchmark output omitted %q: %q", prefix, output)
	}
	var average float64
	if _, err := fmt.Sscanf(output[index:], prefix+" %f milliseconds", &average); err != nil {
		return 0, fmt.Errorf("parse wolfSSL benchmark average: %w", err)
	}
	return average, nil
}

func requireBenchmarkOutput(output string, required []string) error {
	for _, substring := range required {
		if !strings.Contains(output, substring) {
			return fmt.Errorf("wolfSSL benchmark output does not contain %q:\n%s", substring, output)
		}
	}
	return nil
}

func waitForGoBenchmarkServer(b *testing.B, done <-chan error) {
	b.Helper()
	select {
	case err := <-done:
		if err != nil {
			b.Fatal(err)
		}
	case <-time.After(10 * time.Second):
		b.Fatal("go-dtls benchmark server timed out")
	}
}

func waitForWolfSSLBenchmarkServer(b *testing.B, done <-chan error, output *lockedBuffer, required []string) {
	b.Helper()
	select {
	case err := <-done:
		if err != nil {
			b.Fatalf("wolfSSL benchmark server failed: %v\n%s", err, output.String())
		}
	case <-time.After(10 * time.Second):
		b.Fatalf("wolfSSL benchmark server timed out\n%s", output.String())
	}
	if err := requireBenchmarkOutput(output.String(), required); err != nil {
		b.Fatal(err)
	}
}

func TestParseWolfSSLBenchmarkAverage(t *testing.T) {
	for _, test := range []struct {
		output string
		resume bool
	}{
		{"wolfSSL_connect avg took:  12.345 milliseconds\n", false},
		{"wolfSSL_resume  avg took:   12.345 milliseconds\n", true},
	} {
		average, err := parseWolfSSLBenchmarkAverage(test.output, test.resume)
		if err != nil || average != 12.345 {
			t.Fatalf("average=%v err=%v", average, err)
		}
	}
}

func requireWolfSSLPSK(t testing.TB, root, executable string) {
	t.Helper()
	cmd := exec.Command(executable, "-?") // #nosec G204 -- validated local WOLFSSL_ROOT executable in this opt-in test.
	cmd.Dir = root
	output, _ := cmd.CombinedOutput()
	if !bytes.Contains(output, []byte("Use pre Shared keys")) {
		t.Skip("wolfSSL interoperability build does not enable PSK callbacks")
	}
}
