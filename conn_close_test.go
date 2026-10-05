package dtls13

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net"
	"sync"
	"testing"
	"time"
)

type closeHandshakeConn struct {
	net.Conn
	read          func([]byte) (int, error)
	resetDeadline func()
}

func (c *closeHandshakeConn) Read(p []byte) (int, error) {
	if c.read != nil {
		return c.read(p)
	}
	return c.Conn.Read(p)
}

func (c *closeHandshakeConn) SetDeadline(deadline time.Time) error {
	if deadline.IsZero() && c.resetDeadline != nil {
		c.resetDeadline()
	}
	return c.Conn.SetDeadline(deadline)
}

func waitCloseHandshakeResult(t *testing.T, result <-chan error) error {
	t.Helper()
	select {
	case err := <-result:
		return err
	case <-time.After(3 * time.Second):
		t.Fatal("Close or Handshake did not return")
		return nil
	}
}

func TestConnCloseBeforeHandshake(t *testing.T) {
	for _, constructor := range []struct {
		name string
		new  func(net.Conn, *Config) *Conn
	}{{"client", Client}, {"server", Server}} {
		t.Run(constructor.name, func(t *testing.T) {
			left, right := memoryDatagramPair()
			defer left.Close()
			defer right.Close()
			conn := constructor.new(left, &Config{})
			if err := conn.Close(); err != nil {
				t.Fatal(err)
			}
			for range 2 {
				if err := conn.Handshake(); !errors.Is(err, net.ErrClosed) {
					t.Fatalf("Handshake after Close = %v; want net.ErrClosed", err)
				}
			}
		})
	}
}

func TestConnCloseDuringHandshake(t *testing.T) {
	certificate, roots := testServerCertificate(t)
	for _, side := range []string{"client", "server"} {
		for _, stage := range []string{"blocked_read", "encrypted_flight", "encrypted_race", "callback", "finalized"} {
			t.Run(side+"/"+stage, func(t *testing.T) {
				left, right := memoryDatagramPair()
				transport := &closeHandshakeConn{Conn: left}
				clientConfig := &Config{RootCAs: roots, ServerName: "server.test", SessionTicketsDisabled: true, HandshakeTimeout: 10 * time.Second}
				serverConfig := &Config{Certificates: []tls.Certificate{certificate}, SessionTicketsDisabled: true, HandshakeTimeout: 10 * time.Second}
				var conn, peer *Conn
				if side == "client" {
					conn, peer = Client(transport, clientConfig), Server(right, serverConfig)
				} else {
					conn, peer = Server(transport, serverConfig), Client(right, clientConfig)
				}

				ready, release := make(chan struct{}), make(chan struct{})
				var readyOnce, releaseOnce sync.Once
				resume := func() { releaseOnce.Do(func() { close(release) }) }
				pause := func() {
					readyOnce.Do(func() {
						close(ready)
						<-release
					})
				}
				closeDone := make(chan error, 1)
				switch stage {
				case "blocked_read":
					transport.read = func(p []byte) (int, error) {
						readyOnce.Do(func() { close(ready) })
						return left.Read(p)
					}
				case "encrypted_flight", "encrypted_race":
					transport.read = func(p []byte) (int, error) {
						n, err := left.Read(p)
						if err == nil && n > 0 && isUnifiedRecord(p[:n]) {
							pause()
						}
						return n, err
					}
				case "callback":
					// Close must return even when called by the handshake itself.
					if side == "client" {
						clientConfig.VerifyPeerCertificate = func([][]byte, [][]*x509.Certificate) error {
							closeDone <- conn.Close()
							pause()
							return nil
						}
					} else {
						serverConfig.GetCertificate = func(*ClientHelloInfo) (*tls.Certificate, error) {
							closeDone <- conn.Close()
							pause()
							return &certificate, nil
						}
					}
				case "finalized":
					// runHandshake clears its deadline after publishing keys and state.
					transport.resetDeadline = pause
				}

				var handshakes sync.WaitGroup
				allDone := make(chan struct{})
				t.Cleanup(func() {
					_ = left.Close()
					_ = right.Close()
					resume()
					select {
					case <-allDone:
					case <-time.After(3 * time.Second):
						t.Error("handshake goroutine did not exit during cleanup")
					}
				})
				handshakeDone := make(chan error, 1)
				handshakes.Go(func() { handshakeDone <- conn.Handshake() })
				if stage != "blocked_read" {
					handshakes.Go(func() { _ = peer.Handshake() })
				}
				go func() { handshakes.Wait(); close(allDone) }()
				select {
				case <-ready:
				case <-time.After(3 * time.Second):
					t.Fatal("handshake did not reach the shutdown point")
				}

				var secrets [][]byte
				if stage == "finalized" {
					if !conn.ConnectionState().HandshakeComplete || conn.sendingTraffic == nil || conn.receivingTraffic == nil {
						t.Fatal("handshake did not install application keys")
					}
					secrets = [][]byte{conn.sendingTraffic.secret, conn.receivingTraffic.secret, conn.resumptionMasterSecret, conn.state.exporter.secret}
				}
				if stage != "callback" {
					go func() {
						if stage == "encrypted_race" {
							<-release
						}
						closeDone <- conn.Close()
					}()
				}
				if stage == "encrypted_race" {
					resume()
				}
				if err := waitCloseHandshakeResult(t, closeDone); err != nil {
					t.Fatalf("Close: %v", err)
				}
				// Close must not wait for the paused handshake or callback.
				resume()
				err := waitCloseHandshakeResult(t, handshakeDone)
				if !errors.Is(err, net.ErrClosed) && (stage != "encrypted_race" || err != nil) {
					t.Fatalf("Handshake = %v; want net.ErrClosed (or success if it won the race)", err)
				}
				if err := conn.Close(); !errors.Is(err, net.ErrClosed) {
					t.Fatalf("repeated Close = %v; want net.ErrClosed", err)
				}
				conn.dispatchMu.Lock()
				conn.writeMu.Lock()
				retained := conn.sendCipher != nil || conn.sendingTraffic != nil || conn.receivingTraffic != nil || conn.finishedACKCipher != nil || conn.resumptionMasterSecret != nil || conn.resumptionSuite != nil
				conn.writeMu.Unlock()
				if conn.receiveEpochs != nil {
					conn.receiveEpochs.mu.RLock()
					retained = retained || len(conn.receiveEpochs.ciphers) != 0
					conn.receiveEpochs.mu.RUnlock()
				}
				conn.dispatchMu.Unlock()
				if retained || conn.ConnectionState().exporter != nil {
					t.Error("closed handshake retained traffic, resumption, or exporter state")
				}
				for _, secret := range secrets {
					if !bytes.Equal(secret, make([]byte, len(secret))) {
						t.Error("closed handshake did not erase secret backing storage")
					}
				}
				if stage != "encrypted_race" {
					conn.readerMu.Lock()
					running := conn.readerRunning
					conn.readerMu.Unlock()
					if running {
						t.Error("closed handshake started the record reader")
					}
				}
			})
		}
	}
}
