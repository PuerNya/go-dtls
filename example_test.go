package dtls13_test

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"os"

	dtls13 "github.com/puernya/go-dtls"
)

func ExampleDial() {
	rootPEM, err := os.ReadFile("ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	roots := x509.NewCertPool()
	if ok := roots.AppendCertsFromPEM(rootPEM); !ok {
		log.Fatal("invalid root certificate")
	}

	conn, err := dtls13.Dial("udp", "dtls.example:4433", &dtls13.Config{
		RootCAs:    roots,
		ServerName: "dtls.example",
		NextProtos: []string{"example/1"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if _, err := conn.WriteDatagram([]byte("ping")); err != nil {
		log.Fatal(err)
	}
}

func ExampleListen() {
	certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := dtls13.Listen("udp", ":4433", &dtls13.Config{
		Certificates: []tls.Certificate{certificate},
		NextProtos:   []string{"example/1"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			return
		}
		go func() {
			defer conn.Close()
			if err := conn.Handshake(); err != nil {
				log.Printf("DTLS handshake: %v", err)
				return
			}
			buffer := make([]byte, 1200)
			n, info, err := conn.ReadDatagram(buffer)
			if err != nil {
				log.Printf("read datagram: %v", err)
				return
			}
			if info.Truncated {
				log.Printf("discarded %d-byte datagram", info.FullLength)
				return
			}
			if _, err := conn.WriteDatagram(buffer[:n]); err != nil {
				log.Printf("write datagram: %v", err)
			}
		}()
	}
}

func ExampleConn_ReadDatagram() {
	read := func(conn *dtls13.Conn) ([]byte, error) {
		buffer := make([]byte, 1200)
		n, info, err := conn.ReadDatagram(buffer)
		if err != nil {
			return nil, err
		}
		if info.Truncated {
			// The rest of this datagram has already been discarded.
			return nil, fmt.Errorf("datagram needs %d bytes", info.FullLength)
		}
		return append([]byte(nil), buffer[:n]...), nil
	}

	_ = read
}

func ExampleConn_WriteDatagram() {
	send := func(conn *dtls13.Conn, payload []byte) error {
		if err := conn.Handshake(); err != nil {
			return err
		}
		maximum := conn.PathMTU() - conn.RecordOverhead()
		if len(payload) > maximum {
			return fmt.Errorf("payload is %d bytes; current maximum is %d", len(payload), maximum)
		}

		_, err := conn.WriteDatagram(payload)
		if errors.Is(err, dtls13.ErrDatagramTooLarge) {
			// The path MTU can decrease after the check above. Fragmentation and
			// retry policy belong to the application protocol.
			return fmt.Errorf("resize datagram: %w", err)
		}
		return err
	}

	_ = send
}

func ExampleConn_WriteEarlyData() {
	sendReplaySafe := func(conn *dtls13.Conn, payload []byte) error {
		_, err := conn.WriteEarlyData(payload)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, dtls13.ErrEarlyDataUnavailable),
			errors.Is(err, dtls13.ErrEarlyDataRejected):
			// Retry only because this application operation is safe to repeat.
			_, err = conn.WriteDatagram(payload)
			return err
		default:
			return err
		}
	}

	_ = sendReplaySafe
}
