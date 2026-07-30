package dtls13

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
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

func wolfSSLPaths(t *testing.T) (root, server, client string) {
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

func unusedUDPPort(t *testing.T) int {
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

func wolfSSLCertificate(t *testing.T, root, name string) tls.Certificate {
	t.Helper()
	certificate, err := tls.LoadX509KeyPair(filepath.Join(root, "certs", name+"-cert.pem"), filepath.Join(root, "certs", name+"-key.pem"))
	if err != nil {
		t.Fatal(err)
	}
	return certificate
}

func wolfSSLClientCAs(t *testing.T, root string) *x509.CertPool {
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

func wolfSSLExternalPSK(t *testing.T) *ExternalPSK {
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

func requireWolfSSLPSK(t *testing.T, root, executable string) {
	t.Helper()
	cmd := exec.Command(executable, "-?") // #nosec G204 -- validated local WOLFSSL_ROOT executable in this opt-in test.
	cmd.Dir = root
	output, _ := cmd.CombinedOutput()
	if !bytes.Contains(output, []byte("Use pre Shared keys")) {
		t.Skip("wolfSSL interoperability build does not enable PSK callbacks")
	}
}
