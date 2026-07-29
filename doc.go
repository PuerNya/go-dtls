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
// [Config.InsecureSkipVerify] disables the built-in peer identity check and
// should not be enabled in production unless [Config.VerifyPeerCertificate]
// performs an equivalent check. Encryption without authentication does not
// prevent an active attacker from impersonating a peer.
//
// # Session resumption and 0-RTT
//
// Set [Config.ClientSessionCache] to retain NewSessionTicket state and enable
// client-side session resumption. The built-in cache returned by
// [NewLRUClientSessionCache] is bounded and safe for concurrent use. Servers
// can share [Config.SessionTicketKey] when tickets must remain valid across
// server instances.
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
// [RFC 9146]: https://www.rfc-editor.org/rfc/rfc9146
// [RFC 9853]: https://www.rfc-editor.org/rfc/rfc9853
//
// [RFC 9846]: https://www.rfc-editor.org/rfc/rfc9846
package dtls13
