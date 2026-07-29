package dtls13

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func wolfSSLPaths(t *testing.T) (root, server, client string) {
	t.Helper()
	root = os.Getenv("WOLFSSL_ROOT")
	if root == "" {
		t.Skip("set WOLFSSL_ROOT to run wolfSSL DTLS 1.3 interoperability tests")
	}
	server = filepath.Join(root, "build-zig", "examples", "server", "server.exe")
	client = filepath.Join(root, "build-zig", "examples", "client", "client.exe")
	for _, path := range []string{server, client} {
		if _, err := os.Stat(path); err != nil {
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
	testInteropWolfSSLServer(t, "", nil)
}

func TestInteropWolfSSLServerAES128CCM(t *testing.T) {
	testInteropWolfSSLServer(t, "TLS13-AES128-CCM-SHA256", []uint16{TLS_AES_128_CCM_SHA256})
}

func testInteropWolfSSLServer(t *testing.T, cipherName string, suites []uint16) {
	t.Helper()
	root, serverPath, _ := wolfSSLPaths(t)
	port := unusedUDPPort(t)
	args := []string{"-u", "-v", "4", "-e", "-d", "-p", strconv.Itoa(port)}
	if cipherName != "" {
		args = append(args, "-l", cipherName)
	}
	cmd := exec.Command(serverPath, args...)
	cmd.Dir = root
	var output bytes.Buffer
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
	conn, err := DialWithDialer(dialer, "udp4", fmt.Sprintf("127.0.0.1:%d", port), &Config{
		InsecureSkipVerify: true, CipherSuites: suites, HandshakeTimeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("Go client to wolfSSL 5.9.2 server: %v\n%s", err, output.String())
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	payload := []byte("go-dtls to wolfSSL")
	if _, err = conn.WriteDatagram(payload); err != nil {
		t.Fatal(err)
	}
	reply := make([]byte, 64)
	n, _, err := conn.ReadDatagram(reply)
	if err != nil {
		t.Fatalf("read wolfSSL echo: %v\n%s", err, output.String())
	}
	if !bytes.Equal(reply[:n], payload) {
		t.Fatalf("wolfSSL echo = %q, want %q", reply[:n], payload)
	}
}

func TestInteropWolfSSLClient(t *testing.T) {
	testInteropWolfSSLClient(t, "", nil)
}

func TestInteropWolfSSLClientAES128CCM(t *testing.T) {
	testInteropWolfSSLClient(t, "TLS13-AES128-CCM-SHA256", []uint16{TLS_AES_128_CCM_SHA256})
}

func testInteropWolfSSLClient(t *testing.T, cipherName string, suites []uint16) {
	t.Helper()
	root, _, clientPath := wolfSSLPaths(t)
	certificate, err := tls.LoadX509KeyPair(filepath.Join(root, "certs", "server-cert.pem"), filepath.Join(root, "certs", "server-key.pem"))
	if err != nil {
		t.Fatal(err)
	}
	listener, err := Listen("udp4", "127.0.0.1:0", &Config{
		Certificates: []tls.Certificate{certificate}, CipherSuites: suites, HandshakeTimeout: 5 * time.Second, SessionTicketsDisabled: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.UDPAddr).Port

	args := []string{"-u", "-v", "4", "-d", "-x", "-h", "127.0.0.1", "-p", strconv.Itoa(port)}
	if cipherName != "" {
		args = append(args, "-l", cipherName)
	}
	cmd := exec.Command(clientPath, args...)
	cmd.Dir = root
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	if err = cmd.Start(); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	type acceptResult struct {
		conn *Conn
		err  error
	}
	accepted := make(chan acceptResult, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		accepted <- acceptResult{conn: conn, err: acceptErr}
	}()
	var conn *Conn
	select {
	case result := <-accepted:
		conn, err = result.conn, result.err
	case <-time.After(5 * time.Second):
		_ = listener.Close()
		_ = cmd.Process.Kill()
		t.Fatalf("accept wolfSSL client timed out\n%s", output.String())
	}
	if err != nil {
		_ = cmd.Process.Kill()
		t.Fatalf("accept wolfSSL client: %v\n%s", err, output.String())
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
	type readResult struct {
		n   int
		err error
	}
	request := make([]byte, 64)
	readDone := make(chan readResult, 1)
	go func() {
		n, _, readErr := conn.ReadDatagram(request)
		readDone <- readResult{n: n, err: readErr}
	}()
	var n int
	select {
	case result := <-readDone:
		n, err = result.n, result.err
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		if conn != nil {
			conn.receiveEpochs.mu.RLock()
			failures := conn.receiveEpochs.ciphers[3].authFailures
			conn.receiveEpochs.mu.RUnlock()
			transcriptHash := conn.postHandshakeTranscript.sum()
			t.Fatalf("wolfSSL client completed the handshake but sent no application data (epoch 3 auth failures: %d, full transcript hash: %x)\n%s", failures, transcriptHash, output.String())
		}
		t.Fatalf("wolfSSL client completed the handshake but sent no application data\n%s", output.String())
	}
	if err != nil {
		_ = cmd.Process.Kill()
		t.Fatalf("read wolfSSL request: %v\n%s", err, output.String())
	}
	if !strings.Contains(string(request[:n]), "hello wolfssl") {
		t.Fatalf("unexpected wolfSSL request %q", request[:n])
	}
	if _, err = conn.WriteDatagram([]byte("hello from go-dtls")); err != nil {
		t.Fatal(err)
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
}
