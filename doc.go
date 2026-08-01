// Package dtls13 implements Datagram Transport Layer Security 1.3 as defined
// by [RFC 9147]. It provides authenticated, encrypted datagrams over UDP and
// supports the TLS 1.3 handshake inherited from [RFC 8446]. The package does
// not implement DTLS 1.2 or protocol-version fallback.
//
// # Datagram semantics
//
// DTLS preserves application message boundaries. A [Conn] is therefore not a
// [net.Conn] byte stream and intentionally implements neither [net.Conn] nor
// [net.PacketConn]. Applications use [Conn.WriteDatagram] and
// [Conn.ReadDatagram] instead.
//
// Each call to WriteDatagram sends exactly one DTLS Application Data record.
// Application data is not fragmented, reordered, or retransmitted by this
// package. By default, a datagram that exceeds the current path MTU or the
// DTLS record limit is rejected with [ErrDatagramTooLarge] before any partial
// record is sent. [Config.IgnorePathMTU] skips the library's PMTU check for
// application data, but the record limit and transport errors still apply.
// Applications that need messages larger than the available payload must
// define their own fragmentation and recovery protocol.
//
// Endpoints advertise the [RFC 8449] record_size_limit extension by default.
// [Config.RecordSizeLimit] sets the complete protected plaintext this endpoint
// accepts independently of PMTU; the peer's separately advertised value limits
// outgoing records. [ConnectionState] reports whether the extension was
// negotiated and the effective limit in each direction.
//
// Each successful call to ReadDatagram consumes exactly one authenticated
// Application Data record. If the destination is too small, the unread bytes
// are discarded and [DatagramInfo.Truncated] is true. [DatagramInfo.FullLength]
// reports the original length. Empty application datagrams are valid. Replay
// protection discards duplicate and stale records, but records may otherwise
// be delivered out of order because the underlying transport is unreliable.
//
// Handshake and post-handshake messages have different reliability rules.
// The package fragments them to the configured MTU, acknowledges records, and
// retransmits unacknowledged flight data with bounded exponential backoff as
// required by RFC 9147. Those mechanisms do not make application data
// reliable.
//
// # Clients and servers
//
// [Dial] and [DialWithDialer] create a connected UDP client and complete its
// handshake before returning. [Listen] owns a UDP socket and returns a
// [Listener], whose [Listener.Accept] method yields one [Conn] per DTLS
// association. [NewListener] places the same association demultiplexer over an
// existing [net.PacketConn].
//
// [Client] and [Server] wrap an existing connected datagram [net.Conn], usually
// a [net.UDPConn]. The transport must preserve datagram boundaries: a stream
// transport such as TCP is not valid. These constructors defer the handshake
// until [Conn.Handshake], [Conn.ReadDatagram], or [Conn.WriteDatagram] is first
// called. Closing a Conn closes its underlying transport.
//
// A Config may be shared by multiple connections, but it must not be modified
// after first use. Use [Config.Clone] to derive an independent configuration.
// Conn methods synchronize protocol state; a reader and a writer may operate
// concurrently. Multiple ReadDatagram calls are serialized, as are multiple
// writes.
//
// # Authentication
//
// Certificate and ALPN configuration follows the corresponding TLS 1.3
// concepts in [crypto/tls.Config]. Clients normally set [Config.RootCAs] and
// [Config.ServerName]. Servers set [Config.Certificates] or
// [Config.GetCertificate], and may configure [Config.ClientAuth] and
// [Config.ClientCAs] for client authentication.
//
// [Config.GetCertificate] receives [ClientHelloInfo], including the client's
// optional [Config.ServerCertificateAuthorities] hint and offered signature
// schemes. Without a callback, multiple server certificates are tried in
// configuration order and the first certificate compatible with SNI, the
// client's signature schemes, and its CA hint is selected; when no certificate
// matches, the first configured certificate remains the compatibility fallback.
//
// A client can configure multiple certificates and the first chain compatible
// with the server's [CertificateRequestInfo] is selected. [Config.GetClientCertificate]
// can make that decision explicitly. The request exposes acceptable CAs,
// signature schemes, and RFC 9846 [CertificateOIDFilter] values through
// [CertificateRequestInfo.SupportsCertificate]. Recognized Key Usage and
// Extended Key Usage filters are enforced; unknown OIDs are ignored as required
// by RFC 9846. [Config.ClientCAs] populates the server's
// certificate_authorities hint, while [Config.ClientCertificateOIDFilters]
// adds oid_filters to initial and post-handshake requests.
//
// [Config.InsecureSkipVerify] disables the built-in peer identity check and
// should not be enabled in production unless [Config.VerifyPeerCertificate]
// performs an equivalent check. Encryption without authentication does not
// prevent an active attacker from impersonating a peer.
//
// # Hybrid key exchange
//
// [Config.CurvePreferences] accepts the RFC 9954 hybrid groups
// X25519MLKEM768, SecP256r1MLKEM768, and SecP384r1MLKEM1024 defined by
// crypto/tls. X25519MLKEM768 is the recommended general-purpose choice.
// Hybrid groups are explicit opt-in: an empty CurvePreferences keeps the
// package default of X25519 followed by P-256.
//
// The wire encoding follows the current ECDHE-MLKEM profile: each key share
// is a fixed-length concatenation with no inner length fields, and the shared
// secret uses the same component order. If a client configures both a hybrid
// group and its traditional ECDH component, its first ClientHello reuses that
// ECDH key for the fallback share. A HelloRetryRequest causes the second
// ClientHello to contain only the requested share. Large hybrid handshake
// messages use the ordinary RFC 9147 fragmentation, ACK, and retransmission
// machinery.
//
// # Encrypted ClientHello
//
// [Config.EncryptedClientHelloConfigList] enables [RFC 9849] Encrypted
// ClientHello on a client. The value is the complete ECHConfigList wire value,
// including its two-byte length. Obtaining that value through an [RFC 9848]
// DNS SVCB or HTTPS record and decoding its Base64 presentation form are
// application responsibilities.
//
// Servers configure [Config.EncryptedClientHelloKeys] or
// [Config.GetEncryptedClientHelloKeys]. The callback runs against the outer
// ClientHello before SNI, ALPN, PSK, or certificate selection. Once configured,
// ECH is mandatory for the client: it completes the connection only after a
// valid HelloRetryRequest or ServerHello acceptance confirmation. Inner and
// outer ClientHello processing, HPKE encryption, padding, HRR context reuse,
// retry configurations, resumption, and 0-RTT follow RFC 9849.
//
// An authenticated rejection returns [ECHRejectionError]. The rejection
// connection is verified against the ECH public_name using RootCAs, regardless
// of [Config.InsecureSkipVerify], or by
// [Config.EncryptedClientHelloRejectionVerify]. Normal
// [Config.VerifyPeerCertificate] is not called, and the client does not send a
// client certificate on a rejected ECH connection. Callers may retry with the
// returned configuration list only within the same configuration-source and
// transport-endpoint scope.
//
// [Config.EncryptedClientHelloGrease] sends a GREASE ECH extension when no real
// ECH configuration is present; rejection does not fail that ordinary
// connection. [ConnectionState.ECHAccepted] distinguishes accepted ECH from
// GREASE or rejection.
//
// # External pre-shared keys
//
// [Config.ExternalPSKs] enables certificate-free authentication with external
// PSKs following [RFC 9257]. [ImportExternalPSK] is the recommended constructor:
// it implements the [RFC 9258] importer, binds the key to DTLS 1.3 and the
// SHA-256 or SHA-384 target KDF, and uses the required "imp binder" label.
// [NewDirectExternalPSK] is an explicit compatibility path for an already
// TLS-specific PSK and uses the TLS 1.3 "ext binder" label.
//
// Both constructors require at least 128 bits of key material. The package
// offers only psk_dhe_ke, so external-PSK handshakes retain ephemeral key
// exchange and forward secrecy. A Config may contain multiple identities; a
// server ignores unknown identities and can fall back to certificate
// authentication when a certificate is configured.
//
// External PSK identities and importer contexts are opaque bytes transmitted
// in plaintext in ClientHello. Reuse makes connections linkable. Applications
// must not put secrets in either value, should provision one PSK per client and
// server role pair, and must include the provisioning channel binding and both
// peer roles in the importer context when a key is shared by a group.
//
// Base TLS 1.3 forbids combining an external PSK with certificate
// authentication, so [Config.ClientAuth] cannot be enabled on a Config with
// ExternalPSKs. External-PSK 0-RTT is deliberately unavailable because its
// protocol settings and replay lifetime cannot be inferred safely. A
// NewSessionTicket issued after external-PSK authentication may be resumed and
// may use the ordinary ticket-based 0-RTT policy. Removing or changing the
// configured external PSK invalidates tickets derived from it.
// [ConnectionState.ExternalPSKIdentity] and
// [ConnectionState.ExternalPSKContext] report the authenticated origin on the
// initial connection and its ticket resumptions.
//
// # Certificate compression
//
// [Config.EnableCertificateCompression] explicitly enables [RFC 8879] with
// the standard zlib algorithm. A client advertises zlib in ClientHello for the
// server certificate, while a server advertises zlib in CertificateRequest for
// client certificates sent during mTLS or post-handshake authentication.
// Enable the option on both endpoints to allow compression in both directions.
//
// The package sends CompressedCertificate only when the complete compressed
// message is smaller than Certificate and otherwise uses the uncompressed
// message. The peer's advertised algorithm is mandatory, the compressed wire
// message is included in the handshake transcript, and both the declared
// uncompressed length and actual decompression output are bounded by
// [Config.MaxHandshakeMessage].
//
// # Session resumption and 0-RTT
//
// Set [Config.ClientSessionCache] to retain NewSessionTicket state and enable
// client-side session resumption. The built-in cache returned by
// [NewLRUClientSessionCache] is bounded and safe for concurrent use. Servers
// can share [Config.SessionTicketKey] when tickets must remain valid across
// server instances.
//
// [Config.SessionTicketRequest] explicitly enables the [RFC 9149]
// ticket_request extension and requests separate ticket counts after full and
// resumed handshakes. [Config.MaxSessionTickets] bounds the server response.
// Without the extension, the server preserves the RFC 9147 behavior of issuing
// one ticket. The built-in cache atomically consumes distinct requested tickets
// for concurrent resumptions; custom caches keep their existing single-state
// behavior.
//
// A server that authenticated a client certificate stores the presented and
// verified chains in the encrypted ticket. RFC 9846 PSK resumption does not
// repeat CertificateRequest, Certificate, or CertificateVerify. Before accepting
// the ticket, the server enforces the current [Config.ClientAuth] policy, checks
// certificate validity and [Config.ClientCAs], and limits total authentication
// age with [Config.SessionTicketLifetime]. Renewing a ticket does not reset the
// original client CertificateVerify time. [Config.VerifyPeerCertificate] is not
// called again on resumption, matching [crypto/tls]. Rotate
// [Config.SessionTicketKey] or disable tickets when existing sessions must be
// revoked for an application-specific policy change.
//
// [Conn.WriteEarlyData] attempts one client 0-RTT datagram using a cached
// session. The server must configure [Config.MaxEarlyData] and an
// appropriate replay policy. [Config.EarlyDataReplayCache] controls replay
// admission; a nil cache selects a bounded process-wide cache. On an untrusted
// UDP listener, leave
// [Config.AllowEarlyDataWithoutCookie] false; accepting early data before
// return-routability validation increases amplification exposure.
//
// Successful replay-cache admission reduces accidental or malicious reuse of
// a ticket identity within one cache domain, but it cannot give 0-RTT the same
// replay guarantees as 1-RTT data. Early data can be replayed across failures,
// cache domains, or server deployments. It must contain only operations that
// are safe to repeat. Callers must handle [ErrEarlyDataUnavailable] and
// [ErrEarlyDataRejected] and decide whether to retry after the handshake.
//
// # Connection IDs and network paths
//
// The package implements DTLS Connection IDs from [RFC 9146], including the
// RFC 9147 post-handshake update messages. [Config.ConnectionID] and
// [Config.GetConnectionID] configure local IDs; [Conn.SendNewConnectionIDs],
// [Conn.RequestConnectionIDs], and [Conn.UseNextConnectionID] manage updates.
// A Listener can use authenticated CIDs to route records whose source address
// has changed.
//
// A CID authenticates a connection, not a network path. When both peers
// negotiate CID, the package also negotiates the [RFC 9853] Return Routability
// Check by default. A Listener association uses the enhanced procedure: it
// first challenges the old path, falls back to a challenge on the candidate
// path only when needed, enforces the three-times amplification limit, and
// changes [Conn.RemoteAddr] only after validation. If a spare peer CID is
// available, the candidate-path probe uses it without changing old-path
// traffic, and validation activates it before traffic uses the new path.
// Set [Config.DisableReturnRoutabilityCheck] only when the application supplies
// equivalent address validation. A connected UDP client cannot receive packets
// from a changed server address until its underlying transport permits them.
//
// # Post-handshake operations
//
// [Conn.SendKeyUpdate] performs the ACK-gated DTLS 1.3 KeyUpdate procedure;
// the sending epoch changes only after the update record is acknowledged. The
// package also initiates KeyUpdate automatically before an AEAD usage limit is
// reached.
//
// A client sets [Config.PostHandshakeAuth] to advertise post-handshake client
// authentication. A server whose [Config.ClientAuth] requests a certificate
// can then call [Conn.RequestClientCertificate]. Exporters are available after
// a completed handshake through [ConnectionState.ExportKeyingMaterial].
//
// # Resource limits and errors
//
// Config includes explicit bounds for handshake reassembly, queued application
// datagrams, Listener associations, per-association input, replay state, and
// Connection IDs. The defaults suit ordinary Internet-sized datagrams and
// certificate chains. Increase a limit only when the application also accepts
// the corresponding memory and denial-of-service cost.
//
// Local configuration failures are reported as [ConfigError], protocol and
// state-machine failures as [ProtocolError], and fatal alerts received from a
// peer as [AlertError]. Datagram-size errors can be tested with [errors.Is]
// against [ErrDatagramTooLarge]. A valid peer close_notify causes
// ReadDatagram to return [io.EOF]. Network errors and deadlines follow the
// standard [net] package conventions. Applications should inspect errors with
// [errors.Is] and [errors.As], not by matching error strings.
//
// [RFC 9147]: https://www.rfc-editor.org/rfc/rfc9147
// [RFC 8446]: https://www.rfc-editor.org/rfc/rfc8446
// [RFC 8449]: https://www.rfc-editor.org/rfc/rfc8449
// [RFC 8879]: https://www.rfc-editor.org/rfc/rfc8879
// [RFC 9146]: https://www.rfc-editor.org/rfc/rfc9146
// [RFC 9149]: https://www.rfc-editor.org/rfc/rfc9149
// [RFC 9257]: https://www.rfc-editor.org/rfc/rfc9257
// [RFC 9258]: https://www.rfc-editor.org/rfc/rfc9258
// [RFC 9848]: https://www.rfc-editor.org/rfc/rfc9848
// [RFC 9849]: https://www.rfc-editor.org/rfc/rfc9849
// [RFC 9853]: https://www.rfc-editor.org/rfc/rfc9853
//
// [RFC 9954]: https://www.rfc-editor.org/rfc/rfc9954
// [RFC 9846]: https://www.rfc-editor.org/rfc/rfc9846
package dtls13
