package dtls13

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"io"
	"time"
)

// VersionDTLS13 is the DTLS 1.3 version value carried in supported_versions.
const VersionDTLS13 uint16 = 0xfefc

const defaultMaxSessionTickets = 4

// SessionTicketRequest configures the RFC 9149 ticket_request extension.
// Enabled sends the extension. NewSessionCount is used when the server
// performs a full handshake, and ResumptionCount is used when it accepts a
// ticket-based resumption. Setting Enabled with both counts at zero explicitly
// asks the server not to issue tickets for that handshake.
type SessionTicketRequest struct {
	Enabled         bool
	NewSessionCount uint8
	ResumptionCount uint8
}

// Config configures a DTLS client or server. Fields that represent TLS 1.3
// concepts follow [crypto/tls.Config] where practical.
//
// A Config may be reused by multiple connections. It must not be modified
// after it has been passed to a DTLS function. Call [Config.Clone] before
// changing a configuration that is already in use.
type Config struct {
	// Rand provides cryptographically secure randomness for key exchange,
	// signatures, cookies, tickets, and protocol nonces. It must be safe for
	// concurrent use. If nil, crypto/rand.Reader is used.
	Rand io.Reader
	// Time returns the current time for certificate validation, ticket expiry,
	// replay policy, and protocol timers. It must be safe for concurrent use.
	// If nil, time.Now is used.
	Time func() time.Time

	// Certificates contains certificate chains to present to the peer. A
	// server uses the first certificate unless GetCertificate is set. A client
	// considers the first certificate when the server requests client
	// authentication and sends it when compatible. The leaf certificate must be
	// followed by any intermediates, as in tls.Certificate. SHA-1 and MD5
	// certificate signatures are rejected. An RSA server leaf must have a
	// modulus of at least 2048 bits.
	Certificates []tls.Certificate
	// GetCertificate selects a server certificate after the ClientHello has
	// been parsed. It must return a non-nil certificate or an error. When set,
	// it takes precedence over Certificates. It is not used by clients. The
	// callback must be safe for concurrent use when Config is shared.
	GetCertificate func(*ClientHelloInfo) (*tls.Certificate, error)
	// RootCAs defines the roots used by a client to verify the server
	// certificate. If nil, the host's root CA set is used.
	RootCAs *x509.CertPool
	// ClientCAs defines the roots used by a server when ClientAuth requires
	// client-certificate verification. If nil, the host's root CA set is used.
	ClientCAs *x509.CertPool
	// ServerName is sent by clients as SNI and is used to verify the server
	// certificate hostname. Dial derives it from the target address when it is
	// empty; Client does not.
	ServerName string
	// VerifyPeerCertificate, when non-nil, is called after normal certificate
	// parsing and verification. rawCerts contains the peer-provided DER chain.
	// verifiedChains contains the chains built by normal verification, or is nil
	// when built-in verification was skipped. Returning an error aborts the
	// handshake. Like crypto/tls, it is not called on resumed connections. A
	// server instead restores authenticated client state from the protected
	// ticket and checks it against the current ClientAuth, ClientCAs, certificate
	// validity, and SessionTicketLifetime. Disable tickets or rotate
	// SessionTicketKey when a changed callback must invalidate existing sessions.
	VerifyPeerCertificate func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error
	// ClientAuth controls a server's policy for client certificates during the
	// handshake and for RequestClientCertificate. The zero value is
	// tls.NoClientCert. A resumable mutual-TLS session carries the authenticated
	// client certificate state in its encrypted ticket; the abbreviated PSK
	// handshake does not request a fresh certificate. It is ignored by clients.
	ClientAuth tls.ClientAuthType
	// InsecureSkipVerify disables built-in certificate-chain and hostname
	// verification. Certificate signatures within the DTLS handshake are still
	// checked, and the RFC 9325 RSA and SHA-1/MD5 security floors still apply.
	// This should be used only for tests or together with
	// VerifyPeerCertificate that implements equivalent identity verification.
	InsecureSkipVerify bool
	// PostHandshakeAuth makes a client advertise support for post-handshake
	// client authentication and permits it to answer CertificateRequest using
	// Certificates. It has no effect on server configurations.
	PostHandshakeAuth bool
	// EncryptedClientHelloGrease sends an RFC 9849 GREASE ECH extension when no
	// EncryptedClientHelloConfigList is configured. It does not require ECH to
	// be accepted and preserves ordinary certificate and resumption semantics.
	// The zero value keeps the default handshake wire and cost unchanged.
	EncryptedClientHelloGrease bool
	// SessionTicketRequest configures the RFC 9149 ticket_request extension on
	// a client. Its zero value preserves RFC 9147 behavior. The built-in LRU
	// cache pools requested tickets and atomically consumes a distinct ticket
	// per concurrent connection.
	SessionTicketRequest SessionTicketRequest
	// NextProtos lists supported ALPN protocol names in preference order. The
	// server's order controls selection. If either peer provides no list, no
	// protocol is negotiated; otherwise the handshake fails when there is no
	// common protocol.
	NextProtos []string
	// CipherSuites lists supported TLS 1.3 cipher suites in preference order.
	// Empty selects AES-128-GCM, AES-256-GCM, ChaCha20-Poly1305, then AES-128-CCM.
	// TLS_AES_128_CCM_8_SHA256 is intentionally not supported.
	CipherSuites []uint16
	// CurvePreferences contains key exchange groups in preference order. Empty
	// selects X25519 followed by P-256. The ECDHE-MLKEM groups can be enabled
	// explicitly; X25519MLKEM768 is the recommended hybrid choice.
	CurvePreferences []tls.CurveID
	// ExternalPSKs contains immutable externally provisioned PSKs used for
	// certificate-free authentication. Create entries with [ImportExternalPSK]
	// or [NewDirectExternalPSK]. Clients offer all entries compatible with their
	// cipher suites; servers select the first offered identity they recognize.
	// Only psk_dhe_ke is used, so every handshake retains forward secrecy.
	//
	// Identities and importer contexts are plaintext and linkable. External PSK
	// 0-RTT is not offered; a later ticket-based resumption may use 0-RTT under
	// the ordinary MaxEarlyData and replay-cache policy. External PSKs cannot be
	// combined with ClientAuth without RFC 9973, which this package does not
	// implement.
	ExternalPSKs []*ExternalPSK
	// EncryptedClientHelloConfigList is the complete RFC 9849 ECHConfigList
	// obtained from the configuration source, including the two-byte list
	// length. RFC 9848 DNS SVCB/HTTPS lookup and Base64 presentation decoding
	// are application responsibilities. A non-nil value makes ECH mandatory:
	// the connection succeeds only when the server accepts ECH.
	EncryptedClientHelloConfigList []byte
	// EncryptedClientHelloRejectionVerify optionally verifies the public-name
	// connection used to authenticate an ECH rejection. Normal certificate
	// verification and VerifyPeerCertificate are not run for that connection.
	// The callback must verify the ECHConfig public_name and return an error on
	// failure. If nil, RootCAs and the public_name are used regardless of
	// InsecureSkipVerify.
	EncryptedClientHelloRejectionVerify func(ConnectionState) error
	// EncryptedClientHelloKeys contains server ECH configurations and their
	// HPKE private keys. At least one key must have SendAsRetry set. Clients
	// ignore this field. GetEncryptedClientHelloKeys takes precedence.
	EncryptedClientHelloKeys []EncryptedClientHelloKey
	// GetEncryptedClientHelloKeys selects server ECH keys after the outer
	// ClientHello is parsed and before any SNI, ALPN, PSK, or certificate
	// negotiation. It must be safe for concurrent use.
	GetEncryptedClientHelloKeys func(*ClientHelloInfo) ([]EncryptedClientHelloKey, error)

	// MTU is the initial maximum UDP payload, including DTLS record framing,
	// generated by Conn. Values below 256 are rejected. Zero selects 1200 bytes,
	// which avoids IP fragmentation on the Internet paths discussed by RFC 9147
	// section 4.1.1. The effective path MTU may decrease after write errors or
	// repeated handshake timeouts; see Conn.PathMTU.
	MTU int
	// IgnorePathMTU skips the library's current PMTU payload limit for
	// Application Data, including 0-RTT. The zero value is false. When enabled,
	// each application datagram is still limited to one DTLS record with at
	// most 2^14 bytes of content, but the complete record is handed directly to
	// the transport. The network may fragment or drop it, or the transport may
	// return [ErrDatagramTooLarge].
	//
	// This option does not change handshake or post-handshake flight
	// fragmentation, retransmission, or PMTU backoff.
	IgnorePathMTU bool
	// EnableCertificateCompression enables RFC 8879 certificate compression
	// with the standard zlib algorithm. A client advertises support in its
	// ClientHello, and a server advertises support for client certificates in
	// CertificateRequest. A Certificate is compressed only after the peer
	// advertises zlib and when the complete CompressedCertificate message is
	// smaller than the original.
	//
	// The zero value is false. Decompressed messages remain subject to
	// MaxHandshakeMessage.
	EnableCertificateCompression bool
	// RecordSizeLimit is the maximum complete protected plaintext, including
	// the inner content type, that this endpoint accepts. Zero selects the
	// DTLS 1.3 maximum of 2^14+1 bytes. Values from 64 through 2^14+1 are
	// advertised with the RFC 8449 record_size_limit extension. The peer's
	// independently advertised limit controls outgoing records.
	//
	// This limit is independent of MTU and does not include record headers or
	// AEAD expansion.
	RecordSizeLimit uint16
	// FlightInterval is the initial handshake and post-handshake retransmission
	// timeout. Zero selects one second.
	FlightInterval time.Duration
	// MaxFlightInterval caps exponential retransmission backoff. Zero selects
	// 60 seconds.
	MaxFlightInterval time.Duration
	// HandshakeTimeout bounds the complete handshake. Zero selects 30 seconds.
	// A shorter HandshakeContext deadline takes precedence.
	HandshakeTimeout time.Duration
	// ReplayWindow is the number of record sequence numbers tracked per epoch.
	// It must be between 1 and 64. Zero selects 64. Larger values accept more
	// reordering at the cost of retaining a wider anti-replay window.
	ReplayWindow int
	// MaxHandshakeMessage bounds memory allocated while reassembling one
	// handshake message. Zero selects 1 MiB. RFC 9147 permits larger messages,
	// so applications with unusually large certificate chains may raise it up
	// to the protocol limit of 2^24-1 bytes.
	MaxHandshakeMessage int
	// MaxBufferedHandshakeMessages bounds the number of incomplete handshake
	// messages retained across message sequence numbers. Zero selects 8.
	MaxBufferedHandshakeMessages int
	// MaxBufferedHandshakeBytes bounds the total bytes retained for incomplete
	// handshake reassembly. Zero selects four times MaxHandshakeMessage. It must
	// be at least MaxHandshakeMessage.
	MaxBufferedHandshakeBytes int
	// MaxBufferedApplicationData bounds decrypted application data waiting for
	// ReadDatagram, in bytes. Zero selects 1 MiB. Exceeding the limit terminates
	// the association instead of allowing unbounded memory growth.
	MaxBufferedApplicationData int
	// MaxBufferedApplicationDatagrams bounds complete application records
	// waiting for ReadDatagram, including zero-length records. Zero selects 1024.
	// Exceeding the limit terminates the association.
	MaxBufferedApplicationDatagrams int
	// MaxPendingConnections bounds active address-demultiplexed Listener
	// sessions, including sessions returned by Accept. Zero selects 128. New
	// unrecognized peers are ignored while this limit is full.
	MaxPendingConnections int
	// MaxSessionQueueDatagrams bounds unread UDP datagrams per Listener
	// session. Zero selects 64; excess datagrams are dropped.
	MaxSessionQueueDatagrams int

	// ClientSessionCache enables client-side session resumption. A cache may be
	// shared by concurrent clients and therefore must be safe for concurrent use.
	ClientSessionCache ClientSessionCache
	// SessionTicketKey authenticates and encrypts server ticket state. Servers
	// sharing resumable sessions must use the same key. A zero key is replaced
	// with random key material when the configuration is first used. Applications
	// must rotate explicitly configured keys regularly; changing the key
	// immediately invalidates tickets issued with the previous key.
	SessionTicketKey [32]byte
	// SessionTicketLifetime controls ticket validity. Zero selects 24 hours;
	// values from one second through seven days are accepted, matching RFC 9846.
	// For mutual TLS it also bounds the total age of the original online client
	// CertificateVerify authentication, so renewed tickets cannot extend it
	// indefinitely.
	SessionTicketLifetime time.Duration
	// MaxEarlyData advertises and permits at most this many 0-RTT application
	// bytes per resumed connection. Zero disables early data. It does not by
	// itself make early data replay-safe; see EarlyDataReplayCache.
	MaxEarlyData uint32
	// SessionTicketsDisabled disables server NewSessionTicket messages and
	// client ticket use. It also disables session resumption and 0-RTT.
	SessionTicketsDisabled bool
	// MaxSessionTickets limits how many tickets a server sends for one RFC 9149
	// request. Zero selects 4. The limit does not affect clients that omit the
	// extension, for which the server sends one ticket as before.
	MaxSessionTickets uint8
	// EarlyDataReplayCache is shared by server connections to prevent reuse of
	// a PSK identity for 0-RTT. It must be safe for concurrent use. Nil uses a
	// bounded process-wide cache. Deployments with more than one process must
	// provide a cache whose replay domain covers every server that shares ticket
	// keys and accepts 0-RTT.
	EarlyDataReplayCache EarlyDataReplayCache
	// AllowEarlyDataWithoutCookie permits a server to accept 0-RTT on the
	// initial address before a cookie exchange. Keep it false on untrusted UDP
	// listeners: accepting application bytes before return-routability validation
	// increases amplification exposure. A false value causes early data to fall
	// back to a 1-RTT handshake through HelloRetryRequest.
	AllowEarlyDataWithoutCookie bool
	// DisableConnectionID suppresses the client's recommended empty-CID offer
	// when ConnectionID is nil. It has no effect when ConnectionID is non-nil.
	DisableConnectionID bool
	// DisableReturnRoutabilityCheck disables RFC 9853 negotiation. By default,
	// a client offering Connection ID also offers Return Routability Check, and
	// a capable server negotiates it with Connection ID. Set this only when the
	// application validates changed peer addresses itself.
	DisableReturnRoutabilityCheck bool
	// ConnectionID is the CID that the peer places in protected records sent
	// to this endpoint. A non-nil empty slice negotiates support without
	// requesting a CID in this direction. Clients with nil ConnectionID offer
	// an empty CID by default, as recommended by RFC 9147 section 5.1.
	ConnectionID []byte
	// GetConnectionID optionally creates a local CID for each accepted Listener
	// session and in response to RequestConnectionId. It overrides ConnectionID
	// for a newly accepted Listener session. Returned IDs must be at most 255
	// bytes and must be unique and prefix-free among active sessions. The
	// callback may be invoked concurrently for different connections.
	GetConnectionID func() ([]byte, error)
	// MaxConnectionIDs bounds the number of local or peer-provided CIDs kept
	// for one connection. It must be between 1 and 255. Zero selects 8.
	MaxConnectionIDs int

	state *configState
}

type configState struct {
	cookieProtector
	certificateCompression *certificateCompressionCache
}

func (c *Config) clientConnectionIDOffer() ([]byte, bool) {
	if c.ConnectionID != nil {
		return append([]byte(nil), c.ConnectionID...), true
	}
	if !c.DisableConnectionID {
		return []byte{}, true
	}
	return nil, false
}

func ensureCookieProtector(config *Config) error {
	state := ensureConfigState(config)
	protector := &state.cookieProtector
	protector.mu.Lock()
	defer protector.mu.Unlock()
	if len(protector.current.secret) >= 32 {
		return nil
	}
	secret := make([]byte, 32)
	if _, err := io.ReadFull(config.Rand, secret); err != nil {
		return err
	}
	protector.current = cookieKey{id: 1, secret: secret}
	protector.lifetime = time.Minute
	protector.now = config.Time
	protector.rand = config.Rand
	protector.rotated = config.Time()
	return nil
}

// ClientHelloInfo contains information from a ClientHello for a
// Config.GetCertificate callback. Its fields must not be modified.
type ClientHelloInfo struct {
	// ServerName is the SNI name requested by the client, or the empty string
	// when the client sent no server_name extension.
	ServerName string
	// SupportedProtos lists the ALPN protocols offered by the client, in client
	// preference order. The callback must not retain or modify this slice.
	SupportedProtos []string
	// Conn is the server connection processing the ClientHello. Its handshake is
	// not complete while GetCertificate is running.
	Conn *Conn
}

func (c *Config) normalized() (*Config, error) {
	if c == nil {
		c = &Config{}
	}
	// Public slices are immutable after first use. The per-connection copy can
	// share them; Config.Clone retains independent backing arrays for callers.
	x := new(Config)
	if c.EnableCertificateCompression {
		c.ensureCertificateCompressionCache()
		configStateMu.Lock()
		*x = *c
		configStateMu.Unlock()
	} else {
		*x = *c
	}
	if x.Rand == nil {
		x.Rand = rand.Reader
	}
	if x.Time == nil {
		x.Time = time.Now
	}
	if x.MTU == 0 {
		x.MTU = 1200
	}
	if x.MTU < 256 {
		return nil, &ConfigError{"MTU must be at least 256"}
	}
	if x.RecordSizeLimit == 0 {
		x.RecordSizeLimit = defaultRecordSizeLimit
	}
	if x.RecordSizeLimit < minRecordSizeLimit || x.RecordSizeLimit > defaultRecordSizeLimit {
		return nil, &ConfigError{"RecordSizeLimit must be between 64 and 2^14+1"}
	}
	if x.FlightInterval == 0 {
		x.FlightInterval = time.Second
	}
	if x.MaxFlightInterval == 0 {
		x.MaxFlightInterval = 60 * time.Second
	}
	if x.HandshakeTimeout == 0 {
		x.HandshakeTimeout = 30 * time.Second
	}
	if x.ReplayWindow == 0 {
		x.ReplayWindow = 64
	}
	if x.ReplayWindow < 1 || x.ReplayWindow > 64 {
		return nil, &ConfigError{"ReplayWindow must be between 1 and 64"}
	}
	if len(x.CipherSuites) == 0 {
		x.CipherSuites = defaultCipherSuites()
	}
	seenSuites := make(map[uint16]bool, len(x.CipherSuites))
	for _, id := range x.CipherSuites {
		if _, err := cipherSuiteForID(id); err != nil {
			return nil, &ConfigError{"CipherSuites contains an unsupported cipher suite"}
		}
		if seenSuites[id] {
			return nil, &ConfigError{"CipherSuites contains a duplicate cipher suite"}
		}
		seenSuites[id] = true
	}
	if len(x.CurvePreferences) == 0 {
		x.CurvePreferences = []tls.CurveID{tls.X25519, tls.CurveP256}
	}
	seenGroups := make(map[tls.CurveID]bool, len(x.CurvePreferences))
	for _, group := range x.CurvePreferences {
		if !supportedKeyExchangeGroup(group) {
			return nil, &ConfigError{"CurvePreferences contains an unsupported group"}
		}
		if seenGroups[group] {
			return nil, &ConfigError{"CurvePreferences contains a duplicate group"}
		}
		seenGroups[group] = true
	}
	if len(x.ExternalPSKs) > 0 {
		seenExternalIdentities := make(map[string]bool)
		for _, psk := range x.ExternalPSKs {
			if psk == nil || len(psk.identity) == 0 || len(psk.keys) == 0 {
				return nil, &ConfigError{"ExternalPSKs contains an invalid entry"}
			}
			for i := range psk.keys {
				key := &psk.keys[i]
				if len(key.wireIdentity) == 0 || len(key.wireIdentity) > 65535 || len(key.key) < externalPSKMinimumLength || key.binderLabel == nil {
					return nil, &ConfigError{"ExternalPSKs contains an invalid entry"}
				}
				if seenExternalIdentities[string(key.wireIdentity)] {
					return nil, &ConfigError{"ExternalPSKs contains a duplicate wire identity"}
				}
				seenExternalIdentities[string(key.wireIdentity)] = true
			}
		}
		if x.ClientAuth != tls.NoClientCert {
			return nil, &ConfigError{"ExternalPSKs cannot be combined with ClientAuth"}
		}
	}
	if x.EncryptedClientHelloConfigList != nil {
		configs, err := parseECHConfigList(x.EncryptedClientHelloConfigList)
		if err != nil {
			return nil, &ConfigError{"EncryptedClientHelloConfigList is malformed"}
		}
		config, _, _, _ := pickECHConfig(configs)
		if config == nil {
			return nil, &ConfigError{"EncryptedClientHelloConfigList contains no supported configuration"}
		}
	}
	if len(x.EncryptedClientHelloKeys) > 0 {
		if err := validateECHKeys(x.EncryptedClientHelloKeys); err != nil {
			return nil, &ConfigError{err.Error()}
		}
	}
	if x.MaxHandshakeMessage == 0 {
		x.MaxHandshakeMessage = 1 << 20
	}
	if x.MaxHandshakeMessage < 1 || x.MaxHandshakeMessage >= 1<<24 {
		return nil, &ConfigError{"MaxHandshakeMessage must be between 1 and 2^24-1"}
	}
	if x.MaxBufferedHandshakeMessages == 0 {
		x.MaxBufferedHandshakeMessages = 8
	}
	if x.MaxBufferedHandshakeMessages < 1 {
		return nil, &ConfigError{"MaxBufferedHandshakeMessages must be positive"}
	}
	if x.MaxBufferedHandshakeBytes == 0 {
		x.MaxBufferedHandshakeBytes = 4 * x.MaxHandshakeMessage
	}
	if x.MaxBufferedHandshakeBytes < x.MaxHandshakeMessage {
		return nil, &ConfigError{"MaxBufferedHandshakeBytes must be at least MaxHandshakeMessage"}
	}
	if x.MaxBufferedApplicationData == 0 {
		x.MaxBufferedApplicationData = 1 << 20
	}
	if x.MaxBufferedApplicationData < 1 {
		return nil, &ConfigError{"MaxBufferedApplicationData must be positive"}
	}
	if x.MaxBufferedApplicationDatagrams == 0 {
		x.MaxBufferedApplicationDatagrams = 1024
	}
	if x.MaxBufferedApplicationDatagrams < 1 {
		return nil, &ConfigError{"MaxBufferedApplicationDatagrams must be positive"}
	}
	if x.MaxPendingConnections == 0 {
		x.MaxPendingConnections = 128
	}
	if x.MaxPendingConnections < 1 {
		return nil, &ConfigError{"MaxPendingConnections must be positive"}
	}
	if x.MaxSessionQueueDatagrams == 0 {
		x.MaxSessionQueueDatagrams = 64
	}
	if x.MaxSessionQueueDatagrams < 1 {
		return nil, &ConfigError{"MaxSessionQueueDatagrams must be positive"}
	}
	if x.SessionTicketLifetime == 0 {
		x.SessionTicketLifetime = 24 * time.Hour
	}
	if x.SessionTicketLifetime < time.Second || x.SessionTicketLifetime > 7*24*time.Hour {
		return nil, &ConfigError{"SessionTicketLifetime must be between one second and seven days"}
	}
	if x.MaxSessionTickets == 0 {
		x.MaxSessionTickets = defaultMaxSessionTickets
	}
	if len(x.ConnectionID) > 255 {
		return nil, &ConfigError{"ConnectionID must not exceed 255 bytes"}
	}
	if x.MaxConnectionIDs == 0 {
		x.MaxConnectionIDs = 8
	}
	if x.MaxConnectionIDs < 1 || x.MaxConnectionIDs > 255 {
		return nil, &ConfigError{"MaxConnectionIDs must be between 1 and 255"}
	}
	return x, nil
}

// Clone returns a shallow clone of c whose slice-valued configuration fields
// have independent backing arrays. It is safe to clone a Config that is not
// being concurrently modified. Clone returns nil when c is nil.
//
// Callback values, certificate internals, certificate pools, session caches,
// and the time and randomness sources are shared with c.
func (c *Config) Clone() *Config {
	if c == nil {
		return nil
	}
	if c.EnableCertificateCompression {
		return cloneConfigWithStateLock(c)
	}
	return cloneConfig(c)
}

func cloneConfig(c *Config) *Config {
	x := *c
	x.Certificates = append([]tls.Certificate(nil), c.Certificates...)
	x.NextProtos = append([]string(nil), c.NextProtos...)
	x.CipherSuites = append([]uint16(nil), c.CipherSuites...)
	x.CurvePreferences = append([]tls.CurveID(nil), c.CurvePreferences...)
	x.ExternalPSKs = append([]*ExternalPSK(nil), c.ExternalPSKs...)
	if c.EncryptedClientHelloConfigList != nil {
		x.EncryptedClientHelloConfigList = append([]byte{}, c.EncryptedClientHelloConfigList...)
	}
	x.EncryptedClientHelloKeys = append([]EncryptedClientHelloKey(nil), c.EncryptedClientHelloKeys...)
	if c.ConnectionID != nil {
		x.ConnectionID = append([]byte{}, c.ConnectionID...)
	}
	return &x
}

func cloneConfigWithStateLock(c *Config) *Config {
	configStateMu.Lock()
	x := cloneConfig(c)
	configStateMu.Unlock()
	return x
}
