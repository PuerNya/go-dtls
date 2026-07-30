# go-dtls

[简体中文](README.md) | [English](README.en.md) | [Русский](README.ru.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/puernya/go-dtls.svg)](https://pkg.go.dev/github.com/puernya/go-dtls)
[![CI](https://github.com/puernya/go-dtls/actions/workflows/ci.yml/badge.svg)](https://github.com/puernya/go-dtls/actions/workflows/ci.yml)
[![License: GPL v3](https://img.shields.io/badge/license-GPLv3-blue.svg)](../LICENSE)

`go-dtls` is a DTLS 1.3 library implemented in Go. Its protocol behavior follows [RFC 9147](https://www.rfc-editor.org/rfc/rfc9147). The module path is `github.com/puernya/go-dtls`, and the package name is `dtls13`.

The implementation covers the mandatory semantics of RFC 9147 and all 11 of its direct normative references. TLS 1.3 semantics follow [RFC 9846](https://www.rfc-editor.org/rfc/rfc9846). Mandatory and recommended behavior applicable to supported features is implemented; unsupported optional extensions are listed briefly under Scope Boundaries.

> RFC support status is not a third-party compliance certification. RFC 9325 certificate security limits are enforced by default: an RSA server-authentication leaf must be at least 2048 bits, and SHA-1/MD5 certificates cannot bypass this policy through self-signing, trust anchors, `InsecureSkipVerify`, or session resumption.

## Core Semantics

DTLS is an unreliable datagram protocol, not a TLS byte stream:

- `Conn` intentionally implements neither `net.Conn` nor `net.PacketConn`.
- Each `WriteDatagram` sends one DTLS Application Data record. Application datagrams are not internally fragmented, ordered, or retransmitted.
- Each `ReadDatagram` consumes one authenticated Application Data record. If the buffer is too small, the remainder is discarded and reported through `DatagramInfo.Truncated`.
- By default, an application datagram exceeding the current path MTU or the RFC record limit returns `ErrDatagramTooLarge` without a partial write. `IgnorePathMTU` can skip only the former check.
- Handshake messages still use RFC 9147 fragmentation, ACKs, loss recovery, and exponential backoff. These reliability mechanisms do not change application datagram semantics.
- RFC 9954 hybrid key exchange can explicitly enable three standard ECDHE-MLKEM groups; the default remains X25519/P-256.
- RFC 9849 ECH encrypts the real ClientHello; rejection, HRR, resumption, and 0-RTT remain fail-closed.
- `Listener` accepts authenticated DTLS associations from a UDP socket and returns a strongly typed `*dtls13.Conn`.

## Requirements

| Item | Requirement |
| --- | --- |
| Go | Go 1.26 or later; performance data was measured with Go 1.26.3 on `windows/amd64` |
| Transport | `udp`, `udp4`, or `udp6`; TCP is not accepted |
| Windows race | The repository script requires Zig 0.17 and a working CGO toolchain |
| wolfSSL interoperability tests | Optional; set `WOLFSSL_ROOT` to a compatible wolfSSL source/build directory |

Install:

```sh
go get github.com/puernya/go-dtls
```

Because the final segment of the import path contains a hyphen, Go source uses the declared package name `dtls13`:

```go
import dtls13 "github.com/puernya/go-dtls"
```

## Quick Start

### Client

`Dial` creates a UDP connection and completes the handshake before returning. This client verifies the server certificate, sends one datagram, and reads one response:

```go
package main

import (
	"crypto/x509"
	"log"
	"os"
	"time"

	dtls13 "github.com/puernya/go-dtls"
)

func main() {
	caPEM, err := os.ReadFile("ca.pem")
	if err != nil {
		log.Fatal(err)
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caPEM) {
		log.Fatal("ca.pem does not contain a certificate")
	}

	conn, err := dtls13.Dial("udp", "127.0.0.1:4444", &dtls13.Config{
		RootCAs:    roots,
		ServerName: "server.example",
		NextProtos: []string{"example/1"},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
		log.Fatal(err)
	}
	if _, err = conn.WriteDatagram([]byte("ping")); err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1200)
	n, info, err := conn.ReadDatagram(buffer)
	if err != nil {
		log.Fatal(err)
	}
	if info.Truncated {
		log.Fatalf("response needs %d bytes", info.FullLength)
	}
	log.Printf("received %q from %s", buffer[:n], info.Source)
}
```

### Server

`Listen` creates a UDP listener. After receiving a new association, `Accept` returns a `*Conn`. The application may call `Handshake` explicitly or let the first `ReadDatagram` or `WriteDatagram` trigger the handshake.

```go
package main

import (
	"crypto/tls"
	"log"

	dtls13 "github.com/puernya/go-dtls"
)

func main() {
	certificate, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	listener, err := dtls13.Listen("udp", ":4444", &dtls13.Config{
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
			log.Fatal(err)
		}
		go serve(conn)
	}
}

func serve(conn *dtls13.Conn) {
	defer conn.Close()
	if err := conn.Handshake(); err != nil {
		log.Printf("handshake: %v", err)
		return
	}

	buffer := make([]byte, 1200)
	n, info, err := conn.ReadDatagram(buffer)
	if err != nil {
		log.Printf("read: %v", err)
		return
	}
	if info.Truncated {
		log.Printf("discarded truncated datagram: need %d bytes", info.FullLength)
		return
	}
	if _, err = conn.WriteDatagram(buffer[:n]); err != nil {
		log.Printf("write: %v", err)
	}
}
```

## Common API

| API | Purpose |
| --- | --- |
| `Dial` / `DialWithDialer` | Create a UDP client connection and complete the handshake; `DialWithDialer` configures the local address and dial timeout |
| `Listen` | Create and own a UDP listener |
| `NewListener` | Create a DTLS listener over an existing `net.PacketConn`; `Accept` reports configuration errors |
| `Listener.Accept` | Accept a new association and return `*Conn` |
| `Client` / `Server` | Wrap an existing connected `net.Conn` as a DTLS client or server; the handshake can be explicit before first I/O |
| `Handshake` / `HandshakeContext` | Perform the handshake explicitly; later calls return the same result |
| `ReadDatagram` | Read one authenticated datagram and return its source, full length, and truncation status |
| `WriteDatagram` | Send one datagram to the authenticated peer of the association |
| `SetDeadline` / `SetReadDeadline` / `SetWriteDeadline` | Set deadlines for the underlying datagram I/O |
| `ConnectionState` | Get the version, suite, ALPN, certificates, resumption state, external PSK identity/context, active CIDs, RRC state, and negotiated RFC 8449 limits in both directions |
| `Close` | Send `close_notify`, clear traffic/resumption secrets, and close the underlying connection |

### Datagram Size and Truncation

Before writing, applications should use `PathMTU` and `RecordOverhead` to estimate the current maximum payload and account for PMTU decreases during the connection lifetime:

```go
maximum := conn.PathMTU() - conn.RecordOverhead()
if len(payload) > maximum {
	// Fragmentation belongs to the application protocol; each fragment remains a separate datagram.
}
```

The path can change after a successful precheck, so write errors must still be handled:

```go
if _, err := conn.WriteDatagram(payload); errors.Is(err, dtls13.ErrDatagramTooLarge) {
	// Shrink or split the payload according to the application protocol, then retry as a new datagram.
}
```

Set `IgnorePathMTU: true` when the application actively probes PMTU or explicitly relies on IP fragmentation. `WriteDatagram` and `WriteEarlyData` then skip the library PMTU payload check and pass one complete DTLS record directly to the underlying transport. Handshake, ACK, and post-handshake flights still use PMTU fragmentation, retransmission, and backoff. This option does not relax the `2^14`-byte record content limit or a negotiated `record_size_limit`. The transport can still fragment, drop, or return `ErrDatagramTooLarge`; the library does not reduce PMTU or automatically retransmit that application datagram.

A short `ReadDatagram` buffer is not a streaming read. When `Truncated=true`, the unread part of that record has already been discarded, and the next read returns the next record.

### Error Model

| Error | Meaning |
| --- | --- |
| `*ConfigError` | Invalid local configuration, such as an invalid transport, MTU, suite, or resource limit |
| `*ProtocolError` | Received or generated content that violates the protocol state or format |
| `AlertError` | The peer returned a fatal TLS alert; use `errors.As` to obtain the numeric alert description |
| `ErrDatagramTooLarge` | The application datagram exceeds the current PMTU or record limit; the transport may still return it with `IgnorePathMTU`; use `errors.Is` |
| `ErrEarlyDataUnavailable` | No early-data ticket is available, or this connection cannot send 0-RTT |
| `ErrEarlyDataRejected` | The handshake completed, but the peer rejected sent 0-RTT because of HRR, replay, or policy |
| `*ECHRejectionError` | The server rejected ECH after authenticating the `public_name` connection; it may carry a `RetryConfigList` restricted to the same configuration source and endpoint |
| `io.EOF` | The peer sent a valid `close_notify`, closing the read direction |

Deadlines, socket closure, and underlying UDP errors follow Go's `net` error model. Callers must not depend on error strings.

## Advanced Features

| Feature | API / configuration | Description |
| --- | --- | --- |
| External PSK / importer | `ImportExternalPSK`, `NewDirectExternalPSK`, `ExternalPSKs` | RFC 9257/9258 certificate-free authentication; the importer is recommended, only `psk_dhe_ke` is used, and multiple identities, HRR, and ticket resumption are supported |
| Session resumption | `ClientSessionCache`, `NewLRUClientSessionCache` | The client caches NewSessionTicket; the server is controlled by `SessionTicketKey` and ticket settings; mTLS resumption preserves client authentication state |
| 0-RTT | `WriteEarlyData`, `MaxEarlyData`, `EarlyDataReplayCache` | Available only on resumed connections; callers must handle `ErrEarlyDataUnavailable` and `ErrEarlyDataRejected`, and early data must be replay-safe |
| KeyUpdate | `SendKeyUpdate(requestPeer)` | Reliably sent, with the sending epoch switched after ACK; also triggered automatically near AEAD usage limits |
| CID / path validation | `ConnectionID`, `GetConnectionID`, `SendNewConnectionIDs`, `RequestConnectionIDs`, `UseNextConnectionID` | Supports RFC 9146 CID negotiation and updates and negotiates RFC 9853 RRC by default; Listener rebinds only after validating the new path |
| Certificate compression | `EnableCertificateCompression` | Explicitly enables RFC 8879 zlib for server certificates and mTLS/PHA client certificates; sends a plain Certificate when compression is not smaller |
| Encrypted ClientHello | `EncryptedClientHelloConfigList`, `EncryptedClientHelloKeys`, `EncryptedClientHelloGrease` | RFC 9849 Inner/Outer ClientHello, HPKE, HRR, acceptance confirmation, retry configurations, resumption, and 0-RTT |
| Handshake client authentication | `ClientAuth`, `ClientCAs`, `Certificates` | Uses the client-certificate policies from `crypto/tls` |
| Post-handshake client authentication | `PostHandshakeAuth`, `RequestClientCertificate` | The client first advertises support, then the server initiates PHA |
| Exporter | `ConnectionState().ExportKeyingMaterial` | Exports RFC 8446 section 7.5 material with the DTLS `dtls13` label |

### CID Address Changes

A client that offers a CID also offers the RFC 9853 `rrc` extension by default. The server enables path validation only when both CID and RRC are negotiated. After Listener receives an authenticated CID-routed record from a new source, it performs an enhanced check: it challenges the old address first; if the old path remains reachable, the existing binding is retained; after `path_drop` or timeout, it challenges the candidate address. `RemoteAddr` and Listener tuple routing are updated atomically only after the candidate correctly echoes the random cookie. Application writes continue using the old address during validation.

Traffic sent to a candidate path never exceeds three times the number of valid bytes received from that address. The timer is `3xRTT` when an RTT sample exists and one second otherwise. Each challenge uses a fresh CSPRNG cookie; unknown RRC message types and invalid response/drop messages are silently discarded. When the peer has supplied a spare CID, the candidate-path challenge temporarily uses it to avoid reusing the old CID across paths. The spare CID is activated only after validation; application traffic on the old path keeps its original CID during validation.

Applications that already provide equivalent address validation can set `DisableReturnRoutabilityCheck: true`. This disables only RRC, not CID. `Dial` uses connected UDP, and operating systems normally do not deliver packets from a server whose source address changed to that socket. Automatic rebinding therefore primarily applies to Listener associations that can route different sources by CID. An empty CID can negotiate RRC but cannot uniquely demultiplex across five-tuples, so it cannot support Listener address migration.

## Main Configuration

Where TLS 1.3 semantics match, `Config` follows `crypto/tls.Config`. A configuration may be reused but must not be modified after first use; call `Clone` to derive another configuration.

| Configuration | Default / behavior |
| --- | --- |
| `Certificates` / `GetCertificate` | Server certificate; an RSA leaf must be at least 2048 bits, and the complete chain must not use SHA-1/MD5 |
| `RootCAs` / `ServerName` | Client-side server-certificate verification; `Dial` uses the target hostname when `ServerName` is unset |
| `ClientCAs` / `ClientAuth` | Server-side client-certificate verification policy |
| `VerifyPeerCertificate` | Additional verification after standard certificate processing on a full handshake; like `crypto/tls`, it is not called again on resumption |
| `InsecureSkipVerify` | `false` by default; production applications should not rely on it to bypass identity verification |
| `NextProtos` | ALPN protocol list |
| `CipherSuites` | AES-128-GCM, AES-256-GCM, ChaCha20-Poly1305, AES-128-CCM |
| `CurvePreferences` | X25519 and P-256 by default; explicitly accepts `tls.X25519MLKEM768`, `tls.SecP256r1MLKEM768`, and `tls.SecP384r1MLKEM1024` |
| `ExternalPSKs` | Empty by default; configure immutable external PSKs created by `ImportExternalPSK` or `NewDirectExternalPSK`; cannot be combined with `ClientAuth` |
| `EncryptedClientHelloConfigList` | `nil` by default; clients provide the complete RFC 9849 ECHConfigList including its two-byte length; non-nil makes acceptance mandatory |
| `EncryptedClientHelloRejectionVerify` | Optional replacement for built-in `RootCAs` plus `public_name` verification of an ECH rejection connection |
| `EncryptedClientHelloKeys` / `GetEncryptedClientHelloKeys` | Server ECHConfig and HPKE keys; at least one key must set `SendAsRetry`, and the callback runs before SNI, ALPN, or certificate selection |
| `EncryptedClientHelloGrease` | `false` by default; sends GREASE ECH without a real configuration, and rejection does not fail the ordinary connection |
| `MTU` | 1200-byte UDP payload; minimum 256 |
| `IgnorePathMTU` | `false` by default; only Application Data skips the library PMTU check, while handshake behavior is unchanged |
| `RecordSizeLimit` | `0` selects the `2^14+1` default; values from `64..2^14+1` set the complete `DTLSInnerPlaintext` this endpoint accepts and are advertised with RFC 8449 independently of PMTU |
| `EnableCertificateCompression` | `false` by default; enables standard RFC 8879 zlib and sends `CompressedCertificate` only when the peer offered zlib and the complete compressed message is smaller; output is bounded by `MaxHandshakeMessage` |
| `FlightInterval` | One-second initial handshake retransmission interval |
| `MaxFlightInterval` | 60-second exponential-backoff cap |
| `HandshakeTimeout` | 30 seconds |
| `ReplayWindow` | 64 records per epoch |
| `MaxHandshakeMessage` | 1 MiB, configurable up to `2^24-1` |
| `MaxBufferedApplicationData` | 1 MiB |
| `MaxBufferedApplicationDatagrams` | 1024 datagrams |
| `MaxPendingConnections` | 128 Listener sessions |
| `MaxSessionQueueDatagrams` | 64 datagrams per Listener session |
| `SessionTicketLifetime` | 24 hours, maximum 7 days |
| `MaxEarlyData` | 0, so 0-RTT is disabled by default |
| `MaxConnectionIDs` | 8 CIDs per direction |
| `DisableReturnRoutabilityCheck` | `false` by default; disable RRC only when the application provides equivalent address validation |

### Post-Quantum Hybrid Key Exchange

RFC 9954 hybrid key exchange uses the Go standard library ML-KEM implementation and supports three IANA `DTLS-OK=Y` groups. `tls.X25519MLKEM768` is the recommended general-purpose choice. Hybrid key exchange is not enabled by default:

```go
CurvePreferences: []tls.CurveID{
	tls.X25519MLKEM768,
	tls.X25519,
},
```

When the matching traditional group is also listed, the first ClientHello reuses the hybrid ECDH component for its fallback share. After HRR, the client sends only the requested group. Large key shares use the ordinary DTLS handshake fragmentation, ACK, and retransmission machinery.

### Encrypted ClientHello

The client obtains a trusted DNS SVCB/HTTPS `ech` parameter and decodes its Base64 presentation form as required by RFC 9848. The library accepts the complete wire-format ECHConfigList and does not perform DNS lookup:

```go
echConfigList, err := base64.StdEncoding.DecodeString(dnsECHParameter)
if err != nil {
	log.Fatal(err)
}

clientConfig := &dtls13.Config{
	RootCAs:                        roots,
	ServerName:                     "private.example",
	EncryptedClientHelloConfigList: echConfigList,
}
```

The server installs one or more ECHConfig values and corresponding HPKE private keys generated by deployment tooling. `Config` is one ECHConfig without the outer ECHConfigList length; `PrivateKey` uses the corresponding KEM private-key encoding from `crypto/hpke`:

```go
serverConfig := &dtls13.Config{
	Certificates: []tls.Certificate{certificate},
	EncryptedClientHelloKeys: []dtls13.EncryptedClientHelloKey{{
		Config:      echConfig,
		PrivateKey:  echPrivateKey,
		SendAsRetry: true,
	}},
}
```

With real ECH configured, the client succeeds only after validating the HRR or ServerHello acceptance confirmation. An authenticated rejection returns `*ECHRejectionError` and verifies the outer connection against the ECH `public_name`, even when `InsecureSkipVerify` is set. A rejection connection does not invoke ordinary `VerifyPeerCertificate` and does not send a client certificate. Use `RetryConfigList` only with the same DNS configuration source and transport endpoint. Set `EncryptedClientHelloRejectionVerify` for custom public-name verification. `ConnectionState().ECHAccepted` reports acceptance of real ECH; GREASE does not set it.

### External PSKs and the Importer

The RFC 9258 importer is the recommended entry point. It binds an EPSK to DTLS 1.3 with the `dtls13` label and derives separate SHA-256 and SHA-384 target keys. The returned value does not retain the original EPSK:

```go
psk, err := dtls13.ImportExternalPSK(
	[]byte("device-17"),
	provisionedKey, // At least 16 bytes, preferably derived from at least 128 bits of entropy.
	[]byte("client=device-17;server=gateway-2"),
	crypto.SHA256, // Use 0 when the EPSK has no associated hash; SHA-256 is the default.
)
if err != nil {
	log.Fatal(err)
}

config := &dtls13.Config{ExternalPSKs: []*dtls13.ExternalPSK{psk}}
```

Existing deployments with an explicitly TLS-specific PSK can use `NewDirectExternalPSK(identity, key, hash)`. Direct PSKs use `ext binder`; imported PSKs use `imp binder`, and the forms are not interchangeable. Clients may configure multiple identities. A server selects the first known identity compatible with its selected cipher-suite hash, or falls back to certificate authentication when a certificate is configured. Both forms offer only `psk_dhe_ke`. `DidResume` is `false` for the initial external-PSK handshake; a ticket issued afterward can resume normally while retaining the authentication origin through `ConnectionState.ExternalPSKIdentity()` and `ExternalPSKContext()`. Removing or changing the configured external PSK invalidates tickets derived from it.

Identities and importer contexts appear in cleartext in ClientHello, so reuse makes connections linkable and neither field may contain secrets. Provision a PSK for a fixed client/server role pair. Group keys must bind both peer identities and the upstream provisioning channel in the context. Base TLS 1.3 does not combine external PSKs with certificate authentication, so `ClientAuth` cannot be enabled with `ExternalPSKs`. An external PSK does not send 0-RTT directly; only later ticket resumptions may use the ordinary `MaxEarlyData` and replay-cache policy.

### Certificate Compression

With `EnableCertificateCompression: true`, a client offers zlib in ClientHello so the server can compress its certificate. A server offers zlib in CertificateRequest so an enabled client can compress its mTLS or PHA certificate. Enable it on both endpoints when certificates should be compressible in both directions.

The implementation uses only Go's standard-library zlib. It sends `CompressedCertificate` only when the complete message is smaller than a plain Certificate and otherwise falls back safely. Handshake fragmentation, ACK, retransmission, HRR, resumption, and `record_size_limit` semantics are unchanged. Both the declared uncompressed length and actual output are bounded by `MaxHandshakeMessage`.

### Fast mTLS Resumption

The client and server use their normal mTLS settings; additionally enable a client session cache and server tickets:

```go
clientConfig := &dtls13.Config{
	Certificates:       []tls.Certificate{clientCertificate},
	RootCAs:            serverRoots,
	ServerName:         "server.example",
	ClientSessionCache: dtls13.NewLRUClientSessionCache(64),
}
serverConfig := &dtls13.Config{
	Certificates:          []tls.Certificate{serverCertificate},
	ClientAuth:            tls.RequireAndVerifyClientCert,
	ClientCAs:             clientRoots,
	SessionTicketKey:      ticketKey,
	SessionTicketLifetime: time.Hour,
}
```

The first connection performs full `CertificateRequest -> Certificate -> CertificateVerify` client authentication. Later connections resume with the RFC 9147/RFC 9846 PSK handshake without resending certificates. The server recovers the client certificates and verified chains from an AES-256-GCM authenticated and encrypted ticket, then reevaluates them against the active `ClientAuth`, `ClientCAs`, and certificate validity. If policy is not satisfied, the ticket is ignored and a full handshake is used.

A ticket without client identity is used for resumption only when `ClientAuth == tls.NoClientCert`. Any configured client-certificate policy causes a full handshake, preventing anonymous-session resumption.

A renewed ticket retains the time of the original online `CertificateVerify`. `SessionTicketLifetime` limits both the ticket lifetime and the total lifetime of that client authentication. `VerifyPeerCertificate` is not rerun on resumption; rotate `SessionTicketKey` or disable session tickets when application identity policy changes. Applications must periodically rotate explicit ticket keys. This configuration has one active key, so changing it immediately invalidates older tickets. Post-handshake PHA does not issue a supplemental ticket automatically; establish a new full mTLS connection if the PHA identity must be available to later resumptions.

Supported cipher suites:

| Constant | ID | Status |
| --- | --- | --- |
| `TLS_AES_128_GCM_SHA256` | `0x1301` | Supported |
| `TLS_AES_256_GCM_SHA384` | `0x1302` | Supported |
| `TLS_CHACHA20_POLY1305_SHA256` | `0x1303` | Supported |
| `TLS_AES_128_CCM_SHA256` | `0x1304` | Supported |
| `TLS_AES_128_CCM_8_SHA256` | `0x1305` | Explicitly disabled; a general-purpose library cannot guarantee the deployment-level additional forgery protections required by RFC 9147 |

## RFC 9147 Completion

Normative keywords are interpreted according to BCP 14. `MUST`, `MUST NOT`, `REQUIRED`, `SHALL`, and `SHALL NOT` are enforced in both send and receive directions by clients and servers. The client actively implements `SHOULD`-class requirements; the server tolerates peer deviations only when authentication, confidentiality, replay protection, amplification limits, and state consistency are not weakened. Once a `MAY` or `OPTIONAL` capability is negotiated, its conditional mandatory requirements are fully enforced.

### Overall Status

| Specification | Status | Implementation |
| --- | --- | --- |
| [RFC 9147](https://www.rfc-editor.org/rfc/rfc9147) | Complete | Record, Handshake, epochs, ACK, KeyUpdate, CID update, Application Data, and applicable security requirements; recommended behavior is complete for enabled features |
| [RFC 9146](https://www.rfc-editor.org/rfc/rfc9146) | Complete | CID negotiation, directional CIDs, updates, Listener routing, error handling, and address retention; DTLS 1.2-only details do not apply |
| [RFC 8449](https://www.rfc-editor.org/rfc/rfc8449) | Complete | Default CH/EE negotiation, directional limits, minimum 64, fatal `record_overflow`, HRR, resumption, 0-RTT, KeyUpdate, ACK, and PMTU independence |
| [RFC 8879](https://www.rfc-editor.org/rfc/rfc8879) | Complete | Explicit opt-in zlib; directional ClientHello/CertificateRequest negotiation, server and mTLS/PHA client certificates, CompressedCertificate transcripts, safe fallback, and decompression bounds |
| [RFC 9257](https://www.rfc-editor.org/rfc/rfc9257) | Complete | At least 128-bit external PSKs, DHE-only handshakes, opaque identities, multiple identities, certificate fallback, privacy guidance, and pairwise/role deployment requirements are covered |
| [RFC 9258](https://www.rfc-editor.org/rfc/rfc9258) | Complete | `ImportedIdentity`, DTLS `0xfefc`, SHA-256/384 target KDFs, the EPSK source hash, `dtls13derived psk`, and `imp binder` are implemented |
| [RFC 9848](https://www.rfc-editor.org/rfc/rfc9848) | Complete/application integration | Accepts the complete ECHConfigList decoded from a DNS `ech` parameter; SVCB/HTTPS lookup and Base64 decoding belong to the application |
| [RFC 9849](https://www.rfc-editor.org/rfc/rfc9849) | Complete | HPKE, Inner/Outer ClientHello, padding, HRR, acceptance confirmation, retry configurations, authenticated rejection, GREASE, resumption, and 0-RTT |
| [RFC 9954](https://www.rfc-editor.org/rfc/rfc9954) | Complete | Three standard ECDHE-MLKEM groups, traditional-share fallback, HRR, fragmentation, mTLS resumption, 0-RTT, ECH, and strict error semantics; disabled by default |
| [RFC 9846](https://www.rfc-editor.org/rfc/rfc9846) | Complete for enabled features | Ignores `user_canceled(90)` and continues waiting for `close_notify` during the handshake, final-ACK wait, and post-handshake processing; local cryptographic failure without a more specific alert sends `general_error(117)`, while a specific protocol alert always takes precedence |
| [RFC 9325](https://www.rfc-editor.org/rfc/rfc9325) | Partial | PFS, AEAD, SNI/ALPN, tickets, 0-RTT, KeyUpdate, and certificate security limits are covered; OCSP stapling is absent, and this module intentionally does not implement the DTLS 1.2 support required by the BCP |
| [RFC 9525](https://www.rfc-editor.org/rfc/rfc9525) | Partial | Go X.509 and `ServerName` cover DNS-ID/IP-ID; URI-ID, SRV-ID, and application service identities are delegated to caller verification callbacks |
| [RFC 9853](https://www.rfc-editor.org/rfc/rfc9853) | Complete | Extension 61, protected content type 27, all three messages, unknown types, enhanced/basic state machines, the three-times amplification limit, one-second/3xRTT timer, NAT rebinding, off-path protection, and spare-CID cross-path privacy |

### Section Coverage

| RFC 9147 section | Status | Implementation summary |
| --- | --- | --- |
| sections 1-2 Introduction and Terminology | Not applicable | Normative keywords follow RFC 2119/8174; no independent wire feature |
| section 3 Design Goals | Complete | Loss, reordering, duplication, delay, fragmentation, replay protection, dynamic PMTU, and resource reclamation |
| section 4 Record Layer | Complete | DTLSPlaintext, unified header, truncated sequence numbers, epochs, AEAD, sequence-number protection, anti-replay, CID demultiplexing, and usage limits |
| section 5 Handshake | Complete | TLS 1.3 handshake, HRR cookie, authentication, fragmentation/reassembly, flights, ACK, retransmission, timeout, and a new association on the same five-tuple |
| section 6 Epochs | Complete | Epochs 0/1/2/3/4+, 0-RTT, KeyUpdate, and bounded retention and clearing of old keys |
| section 7 ACK | Complete | Content type 26, empty ACK, immediate ACK, partial ACK, sliding window, and reliable post-handshake ACK |
| section 8 KeyUpdate | Complete | ACK-gated updates, retransmission, `update_requested`, old-epoch retention, and limit handling |
| section 9 CID Update | Complete | New/RequestConnectionId, dynamic routing, update ACK, resource limits, and prefix-free validation |
| section 10 Application Data | Complete | Connected-datagram API, datagram boundaries, no ordering/retransmission, explicit truncation, and deadlines |
| section 11 Security | Complete | Cookie rotation, three-times amplification limits for handshake and RRC candidate paths, anti-replay, AEAD limits, post-validation address updates, and bounded state |
| section 12 DTLS 1.2 differences | Complete | All applicable differences are included in record, handshake, epoch, ACK, and CID behavior |
| section 13 DTLS 1.2 updates | Not applicable | This module implements only DTLS 1.3 |
| section 14 IANA | Not applicable | Uses assigned values; the library does not perform registry operations |

### Direct Normative References

This table contains all 11 `Normative References` from the RFC Editor XML for RFC 9147. Lower-layer protocols are implemented by Go and the operating system; this does not mean that the library reimplements the UDP/TCP/IP stack.

| Specification | Use in RFC 9147 | Coverage |
| --- | --- | --- |
| [RFC 8439](https://www.rfc-editor.org/rfc/rfc8439) | Counter, nonce, and block function for ChaCha20 sequence-number protection | Complete |
| [RFC 768](https://www.rfc-editor.org/rfc/rfc768) | UDP transport | Lower layer; the library is restricted to UDP and preserves datagram boundaries |
| [RFC 793](https://www.rfc-editor.org/rfc/rfc793) | Reliable-transport background and MSL reference for old epochs | Applicable semantics complete; the library does not accept TCP transport |
| [RFC 1191](https://www.rfc-editor.org/rfc/rfc1191) | IPv4 PMTU discovery | Lower-layer integration complete; platform MTU errors are normalized to `ErrDatagramTooLarge`; applications can probe actively with `IgnorePathMTU` |
| [RFC 2119](https://www.rfc-editor.org/rfc/rfc2119) | BCP 14 normative keywords | Applicable |
| [RFC 4443](https://www.rfc-editor.org/rfc/rfc4443) | IPv6 ICMP Packet Too Big | Lower-layer integration complete; the OS handles ICMPv6 and the library consumes write-error feedback |
| [RFC 4821](https://www.rfc-editor.org/rfc/rfc4821) | Packetization Layer PMTU Discovery | Complete; write errors and consecutive black-hole timeouts trigger fallback and refragmentation |
| [RFC 6298](https://www.rfc-editor.org/rfc/rfc6298) | Initial RTO, exponential backoff, and cap | Complete; one-second default, 60-second maximum, and RTT samples following Karn's algorithm |
| [RFC 8174](https://www.rfc-editor.org/rfc/rfc8174) | BCP 14 capitalization rules | Applicable |
| [RFC 9146](https://www.rfc-editor.org/rfc/rfc9146) | DTLS Connection ID | Complete; negotiation, updates, routing, isolation, and resource limits |
| [RFC 8446](https://www.rfc-editor.org/rfc/rfc8446) | TLS 1.3 base protocol | Supports handshake, authentication, key schedule, PSK/0-RTT, KeyUpdate, PHA, and exporter; protocol semantics follow RFC 9846 |

### Related Specifications and Extensions

| Specification | Relationship to this implementation | Status |
| --- | --- | --- |
| [RFC 9846](https://www.rfc-editor.org/rfc/rfc9846) | TLS 1.3 KeyShare, PSK/HRR, NST, AEAD limits, KeyUpdate, alerts, and vector bounds | Complete for enabled features; mTLS resumption preserves authentication state, policy/CA/validity, and total authentication lifetime; see Overall Status for `user_canceled` and `general_error` |
| [RFC 8449](https://www.rfc-editor.org/rfc/rfc8449) | TLS/DTLS `record_size_limit` | Clients advertise it by default; servers respond only to an offer. Sending follows the peer limit, receiving follows the local limit, absence restores the protocol maximum, and PMTU remains an independent lower bound |
| [RFC 8879](https://www.rfc-editor.org/rfc/rfc8879) | TLS/DTLS Certificate Compression | Explicit opt-in standard zlib; CH/CR negotiation, server and client certificates, HRR, mTLS, PHA, fragmentation/retransmission, transcripts, and bounded decompression are complete; plain Certificate is used when smaller |
| [RFC 9257](https://www.rfc-editor.org/rfc/rfc9257) | TLS 1.3 external PSK guidance | DHE-only use, multiple identities, unknown-identity fallback, cleartext identity risks, ticket-origin binding, and the external-PSK 0-RTT policy are implemented |
| [RFC 9258](https://www.rfc-editor.org/rfc/rfc9258) | TLS/DTLS 1.3 PSK Importer | SHA-256/384 target derivation, the DTLS label, ImportedIdentity wire encoding, and the distinct binder label are implemented |
| [RFC 9848](https://www.rfc-editor.org/rfc/rfc9848) | ECH DNS configuration bootstrapping | The library parses a complete ECHConfigList; DNS SVCB/HTTPS retrieval and presentation-format Base64 decoding are application integration |
| [RFC 9849](https://www.rfc-editor.org/rfc/rfc9849) | TLS/DTLS Encrypted ClientHello | Client and server, HPKE, Inner/Outer reconstruction, padding, HRR context reuse, acceptance confirmation, retry configurations, rejection authentication, GREASE, PSK resumption, and 0-RTT are complete |
| [RFC 9954](https://www.rfc-editor.org/rfc/rfc9954) | TLS/DTLS Hybrid Key Exchange | Implements the RFC construction and X25519MLKEM768, SecP256r1MLKEM768, and SecP384r1MLKEM1024 from the current ECDHE-MLKEM profile; the concrete profile is still in the RFC Editor queue |
| [RFC 9325](https://www.rfc-editor.org/rfc/rfc9325) | TLS/DTLS deployment security BCP | Tickets use AES-256-GCM and are limited to one second through seven days; RSA 2048-bit and SHA-1/MD5 certificate limits are enforced on full handshakes, trust anchors, and resumption paths; see Overall Status for the OCSP and DTLS 1.2 scope exceptions |
| [RFC 9525](https://www.rfc-editor.org/rfc/rfc9525) | Service identity verification | DNS-ID/IP-ID are verified strictly by default; callers implement application semantics for other reference identifiers |
| [RFC 9853](https://www.rfc-editor.org/rfc/rfc9853) | Return Routability Check for CID address changes | Complete; enhanced check by default, basic check after the old path fails, rebind only after validation, independent candidate-path amplification limit, and spare-CID probing when available |
| [RFC 8701](https://www.rfc-editor.org/rfc/rfc8701) | GREASE anti-ossification | Receivers tolerate valid unknown values while preserving HRR invariants; senders do not actively generate GREASE |

Unsupported optional extensions include RFC 9149 Ticket Requests.

### Scope Boundaries

The following items do not reduce completion of mandatory RFC 9147 semantics, but users should understand the boundaries:

- This module implements only DTLS 1.3 and provides no DTLS 1.2 fallback; it therefore does not claim full RFC 9325 compliance with the BCP requirement that general-purpose implementations support DTLS 1.2.
- Heartbeat record demultiplexing is implemented; the complete Heartbeat protocol is defined by RFC 6520 and is outside RFC 9147 scope.
- The sender uses the valid one-record-per-UDP-datagram mode and exposes no optional multi-record aggregation API.
- Concurrent multiple NewSessionTicket or PHA requests are not exposed; the RFC permits but does not require this capability.
- Automatic RRC rebinding requires a transport that receives from different sources and can send to a selected destination. The standard Listener supports this; connected UDP clients are constrained by operating-system peer filtering. An empty CID cannot uniquely route across five-tuples.
- The wolfSSL master `f699037` build (version string 5.9.2) supports CID, KeyUpdate, PHA, session tickets, 0-RTT, `SESSION_CERTS`, direct external PSKs, and all three hybrid groups, but not RFC 8449, RFC 8879, the RFC 9258 importer, or RFC 9853 RRC. Its ECH/HPKE build cannot complete a DTLS accepted-ECH handshake, so only the successful ordinary handshake with ECH GREASE is recorded; accepted-ECH interoperability is not claimed. For hybrid key exchange, all three groups pass with wolfSSL as the client, while its server passes the X25519 and P-256 hybrids; that server cannot complete a fragmented hybrid ClientHello, and the P-384 hybrid times out even without fragmentation. Other peer limits are explicit: the server's HRR rejects go-dtls client 0-RTT; the client cannot parse the 1421-byte go-dtls mTLS ticket; and the client does not retransmit Finished after losing the final ACK.

## Benchmark

The following representative results were measured on an AMD Ryzen 7 7435H with Go 1.26.3 on Windows/amd64. They are not cross-machine performance guarantees.

| Scenario | Representative result |
| --- | --- |
| Full certificate handshake and close, `BenchmarkConnectionHandshakeLifecycle` | About `619.4 us/op`, `99508 B/op`, `760 allocs/op` |
| RFC 9257/9258 external-PSK handshake and close, `BenchmarkExternalPSKHandshakeLifecycle` | About `355.4 us/op`, `98287 B/op`, `727 allocs/op` |
| Full mTLS, `BenchmarkMutualTLSHandshakeLifecycle/Full` | About `912.9 us/op`, `116237 B/op`, `974 allocs/op` |
| Resumed mTLS, `BenchmarkMutualTLSHandshakeLifecycle/Resumed` | About `457.3 us/op`, `115336 B/op`, `799-800 allocs/op` |
| RFC 9849 ECH full handshake, `BenchmarkECHHandshakeLifecycle/Direct` | About `1.53 ms/op`, `148615 B/op`, `1260 allocs/op` |
| RFC 9849 ECH with HRR, `BenchmarkECHHandshakeLifecycle/HRR` | About `1.63 ms/op`, `151391 B/op`, `1281 allocs/op` |
| RFC 8879 zlib server-certificate handshake, four-certificate chain | About `1.049 ms/op`, `123323 B/op`, `1022 allocs/op` |
| RFC 8879 zlib full mTLS, four-certificate chains in both directions | About `1.746 ms/op`, `160719 B/op`, `1469 allocs/op` |
| RFC 8879 zlib compression / decompression | About `7.2-7.7 us/op`, `4 allocs/op` / `6.3-6.9 us/op`, `4290-4300 B/op`, `6 allocs/op` |
| AES-128-GCM 1200 B seal | About 1.86 GB/s, 1 alloc/op |
| AES-128-GCM 1200 B in-place round trip | About 1.05 GB/s, 1 alloc/op |
| Unauthenticated record error classification | About 12.4-13.7 ns/op, 0 allocs/op |
| Extension marshal | About 316.5-358.0 ns/op, 128 B/op, 1 alloc/op |
| Extension ordered-view parse | About 69.8-80.1 ns/op, 0 B/op, 0 allocs/op |
| Single key_share caller-storage parse | About 33.5-35.0 ns/op, 0 B/op, 0 allocs/op |
| ClientHello marshal | About 436-519 ns/op, 424 B/op, 7 allocs/op |
| ServerHello marshal | About 72-91 ns/op, 112 B/op, 1 alloc/op |
| 1200 B single-fragment handshake reassembly | About 0.47-0.60 us/op, 1280 B/op, 1 alloc/op |
| 4 KiB / MTU 1200 protected-flight construction | About 2.89 us/op, 5616 B/op, 6 allocs/op |
| 4 KiB / MTU 1200 plain-flight construction | About 2.15-2.58 us/op, 5040 B/op, 9 allocs/op |

Full-connection results are medians from repeated `-cpu=1` runs.

Run all benchmarks:

```sh
go test -run '^$' -bench . -benchmem
```

Run full-connection and record-layer benchmarks separately:

```sh
go test -run '^$' -bench '^BenchmarkConnectionHandshakeLifecycle$' -benchmem -benchtime=2000x -cpu=1
go test -run '^$' -bench '^BenchmarkECHHandshakeLifecycle/(Direct|HRR)$' -benchmem -benchtime=2000x -cpu=1
go test -run '^$' -bench '^BenchmarkExternalPSKHandshakeLifecycle$' -benchmem -benchtime=2000x -cpu=1
go test -run '^$' -bench '^BenchmarkMutualTLSHandshakeLifecycle/(Full|Resumed)$' -benchmem -count=10
go test -run '^$' -bench '^BenchmarkCertificateCompression' -benchmem
go test -run '^$' -bench '^BenchmarkProtectedRecord(Seal|RoundTripInPlace)$' -benchmem -count=5
```

The repository also includes focused benchmarks for cipher suites, ACK, records/parsers, transcripts, key schedule, KeyUpdate, handshake messages, reassembly, and flight construction. Keep the Go version, CPU, `-cpu`, and `-benchtime` fixed, and inspect `ns/op`, `B/op`, `allocs/op`, and full-connection profiles together.

## Test Coverage

- RFC 9325 certificate-policy tests cover server configuration, client reception, self-signed certificates, unsent trust anchors, `InsecureSkipVerify`, ordinary resumption, and mTLS resumption, with differential behavior against `crypto/x509` for 1024-bit RSA/SHA-1 trust anchors.
- RFC 8449 tests cover CH/EE, the minimum of 64, directional limits, invalid values and extension combinations, authenticated overflow, HRR, resumption, 0-RTT, KeyUpdate, ACK, PMTU independence, and compatibility with third parties that do not negotiate it.
- RFC 8879 tests cover ClientHello/CertificateRequest negotiation, zlib, CompressedCertificate, transcripts, invalid algorithms/streams/lengths, decompression limits, plain-Certificate fallback, HRR, resumption, mTLS/PHA, record limits, fragmentation/retransmission, weak networks, resource lifecycles, and safe fallback with third parties that do not support it.
- RFC 9257/9258 tests cover independent importer derivation, SHA-256/384 KDF separation, `imp`/`ext` binder separation, direct and imported PSKs, multiple identities, HRR filtering, identity/key/context failures, certificate fallback, connection state, ticket resumption and revocation, the 0-RTT policy, and weak networks.
- RFC 9848/9849 tests cover public configuration vectors, ECHConfig/ECHConfigList, HPKE, Inner/Outer and padding, outer-extension reconstruction, HRR acceptance confirmation and downgrade rejection, authenticated retry configurations, client-certificate suppression, GREASE, resumption, 0-RTT, fragmentation, weak networks, and real UDP.
- Weak-network tests cover bidirectional loss, delay, reordering, and duplication, including CH/SH/Finished/ACK/HRR/mTLS-resumption combinations.
- mTLS tests cover full handshakes, PSK resumption, 0-RTT, CA/policy fallback, and renewed-ticket authentication lifetime.
- RFC 9846 alert tests cover handshake processing, final-ACK waiting, post-handshake reordering, `close_notify`, and local cryptographic failure.
- RFC 9853 tests cover RRC messages/state machines, real UDP NAT rebinding, CID updates, weak-network combinations, and connection resource lifecycles.
- Parser/record fuzzing covers copy and in-place decryption differentials for all four AEADs.
- Bidirectional real-UDP tests with wolfSSL master `f699037` cover HRR, RSA-PSS certificate handshakes, Finished ACK, application data, AES-GCM, AES-128-CCM, direct external PSKs, CID, KeyUpdate, PHA, ordinary session resumption, and all three hybrid groups in the directions supported by the peer. Additional supported directions cover go-dtls-initiated immediate CID switching, mTLS resumption, and Finished retransmission after final-ACK loss, plus wolfSSL-client-initiated 0-RTT. ECH coverage is limited to GREASE fallback because the peer cannot currently complete DTLS accepted-ECH.

See [CONTRIBUTING.en.md](CONTRIBUTING.en.md) for the development environment, required checks, performance validation, and commit rules.

## License

This project is distributed under the [GNU General Public License v3.0](../LICENSE) (`GPL-3.0-only`).
