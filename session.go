package dtls13

import (
	"bytes"
	"container/list"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"errors"
	"hash"
	"io"
	"net"
	"sync"
	"time"
)

const (
	handshakeTypeNewSessionTicket uint8  = 4
	extEarlyData                  uint16 = 42
	maxTicketLifetime                    = 7 * 24 * time.Hour
)

type newSessionTicketMessage struct {
	lifetime     uint32
	ageAdd       uint32
	nonce        []byte
	ticket       []byte
	maxEarlyData uint32
}

func (m *newSessionTicketMessage) marshal() ([]byte, error) {
	if m.lifetime == 0 || time.Duration(m.lifetime)*time.Second > maxTicketLifetime {
		return nil, &ProtocolError{"invalid NewSessionTicket lifetime"}
	}
	if len(m.ticket) == 0 {
		return nil, &ProtocolError{"empty NewSessionTicket ticket"}
	}
	if len(m.nonce) > 255 {
		return nil, &ProtocolError{"8-bit vector overflow"}
	}
	if len(m.ticket) > 65535 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	extensionsLength := 2
	if m.maxEarlyData != 0 {
		extensionsLength += 8
	}
	w := newWireBuilder(8 + 1 + len(m.nonce) + 2 + len(m.ticket) + extensionsLength)
	w.u32(m.lifetime)
	w.u32(m.ageAdd)
	w.bytes8(m.nonce)
	w.bytes16(m.ticket)
	extensionsStart := w.startVector16()
	if m.maxEarlyData != 0 {
		var earlyData [4]byte
		binary.BigEndian.PutUint32(earlyData[:], m.maxEarlyData)
		w.u16(int(extEarlyData))
		w.bytes16(earlyData[:])
	}
	w.endVector16(extensionsStart)
	return w.b, w.err
}

func parseNewSessionTicket(b []byte) (*newSessionTicketMessage, error) {
	p := wireParser{b: b}
	m := &newSessionTicketMessage{lifetime: p.u32(), ageAdd: p.u32()}
	m.nonce = append([]byte(nil), p.bytes8()...)
	m.ticket = append([]byte(nil), p.bytes16()...)
	extensionWire := p.take(len(p.b) - p.off)
	if err := p.done(); err != nil {
		return nil, err
	}
	if len(m.ticket) == 0 {
		return nil, alertError(alertDecodeError, &ProtocolError{"empty NewSessionTicket ticket"})
	}
	var extensionStorage [1]orderedExtension
	extensions, err := parseOrderedExtensionsView(extensionWire, extensionStorage[:0])
	if err != nil {
		return nil, err
	}
	for _, extension := range extensions {
		if extension.typ != extEarlyData {
			if knownExtensionType(extension.typ) {
				return nil, alertError(alertIllegalParameter, &ProtocolError{"recognized extension is not permitted in NewSessionTicket"})
			}
			continue
		}
		if len(extension.value) != 4 {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid NewSessionTicket early_data extension length"})
		}
		m.maxEarlyData = binary.BigEndian.Uint32(extension.value)
	}
	return m, nil
}

// ClientSessionCache stores opaque DTLS client sessions by a package-defined
// server cache key. Implementations must be safe for concurrent use. A client
// consumes a selected ticket when beginning a resumption attempt so that
// concurrent connections do not intentionally reuse the same identity.
type ClientSessionCache interface {
	// Get returns the session associated with sessionKey. The returned state is
	// read-only and must not be modified.
	Get(sessionKey string) (*ClientSessionState, bool)
	// Put associates cs with sessionKey. A nil cs removes any existing entry.
	// The package does not modify cs after the call returns.
	Put(sessionKey string, cs *ClientSessionState)
}

// EarlyDataReplayCache controls server acceptance of PSK identities for 0-RTT.
// Implementations must be safe for concurrent use and should share a replay
// domain with every server process that shares SessionTicketKey and accepts
// early data.
//
// A replay cache is a required mitigation, not a guarantee that early data
// cannot be replayed. Applications must still give 0-RTT operations idempotent
// or otherwise replay-safe semantics.
type EarlyDataReplayCache interface {
	// CheckAndStore atomically admits identity until expires. It returns true
	// only when the identity was not already live and the cache retained the new
	// entry for its validity period within the cache's replay domain. It must
	// fail closed by returning false when it cannot store the entry.
	CheckAndStore(identity string, expires time.Time) bool
}

type earlyDataReplayCache struct {
	mu      sync.Mutex
	entries map[string]time.Time
	order   *list.List
	index   map[string]*list.Element
	max     int
}

type earlyReplayEntry struct {
	identity string
	expires  time.Time
}

// NewLRUEarlyDataReplayCache creates a bounded, concurrency-safe in-process
// 0-RTT replay cache. Expired entries are removed as new identities are
// checked. Live entries are never evicted to make room: a full cache fails
// closed and rejects new early data until entries expire.
//
// A non-positive capacity is rejected by returning nil. The cache covers only
// one process; distributed servers sharing ticket keys need a replay cache with
// a correspondingly shared consistency domain.
func NewLRUEarlyDataReplayCache(capacity int) EarlyDataReplayCache {
	if capacity <= 0 {
		return nil
	}
	return &earlyDataReplayCache{entries: make(map[string]time.Time), order: list.New(), index: make(map[string]*list.Element), max: capacity}
}

func (c *earlyDataReplayCache) CheckAndStore(identity string, expires time.Time) bool {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, value := range c.entries {
		if !value.After(now) {
			delete(c.entries, key)
			if element := c.index[key]; element != nil {
				c.order.Remove(element)
				delete(c.index, key)
			}
		}
	}
	if existing, ok := c.entries[identity]; ok && existing.After(now) {
		return false
	}
	if element := c.index[identity]; element != nil {
		c.order.Remove(element)
		delete(c.index, identity)
	}
	if len(c.entries) >= c.max {
		return false
	}
	c.entries[identity] = expires
	element := c.order.PushFront(earlyReplayEntry{identity: identity, expires: expires})
	c.index[identity] = element
	return true
}

var defaultEarlyDataReplayCache EarlyDataReplayCache = NewLRUEarlyDataReplayCache(4096)

// ClientSessionState is opaque state for one resumable DTLS session. Values are
// created and consumed by this package through ClientSessionCache and must not
// be modified by applications or cache implementations.
type ClientSessionState struct {
	ticket           []byte
	psk              []byte
	nonce            []byte
	suite            uint16
	receivedAt       time.Time
	lifetime         uint32
	ageAdd           uint32
	serverName       string
	protocol         string
	maxEarlyData     uint32
	recordSizeLimit  uint16
	peerCertificates []*x509.Certificate
	verifiedChains   [][]*x509.Certificate
}

type lruSessionCache struct {
	mu       sync.Mutex
	capacity int
	entries  map[string]*list.Element
	order    *list.List
}

type sessionCacheEntry struct {
	key   string
	state *ClientSessionState
}

// NewLRUClientSessionCache returns a concurrency-safe, fixed-capacity client
// session cache. The least recently used entry is removed when the capacity is
// exceeded. State is cloned on insertion and retrieval so callers cannot
// mutate cached ticket secrets.
//
// If capacity is less than one, a default capacity of 64 is used.
func NewLRUClientSessionCache(capacity int) ClientSessionCache {
	if capacity < 1 {
		capacity = 64
	}
	return &lruSessionCache{capacity: capacity, entries: make(map[string]*list.Element), order: list.New()}
}

func cloneClientSessionState(state *ClientSessionState) *ClientSessionState {
	if state == nil {
		return nil
	}
	clone := *state
	clone.ticket = append([]byte(nil), state.ticket...)
	clone.psk = append([]byte(nil), state.psk...)
	clone.nonce = append([]byte(nil), state.nonce...)
	clone.peerCertificates = append([]*x509.Certificate(nil), state.peerCertificates...)
	clone.verifiedChains = make([][]*x509.Certificate, len(state.verifiedChains))
	for i := range state.verifiedChains {
		clone.verifiedChains[i] = append([]*x509.Certificate(nil), state.verifiedChains[i]...)
	}
	return &clone
}

func (c *lruSessionCache) Get(key string) (*ClientSessionState, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	element := c.entries[key]
	if element == nil {
		return nil, false
	}
	c.order.MoveToFront(element)
	return cloneClientSessionState(element.Value.(*sessionCacheEntry).state), true
}

func (c *lruSessionCache) Take(key string) (*ClientSessionState, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	element := c.entries[key]
	if element == nil {
		return nil, false
	}
	delete(c.entries, key)
	c.order.Remove(element)
	return cloneClientSessionState(element.Value.(*sessionCacheEntry).state), true
}

func (c *lruSessionCache) Put(key string, state *ClientSessionState) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if state == nil {
		if element := c.entries[key]; element != nil {
			delete(c.entries, key)
			c.order.Remove(element)
		}
		return
	}
	if element := c.entries[key]; element != nil {
		element.Value.(*sessionCacheEntry).state = cloneClientSessionState(state)
		c.order.MoveToFront(element)
		return
	}
	element := c.order.PushFront(&sessionCacheEntry{key: key, state: cloneClientSessionState(state)})
	c.entries[key] = element
	if c.order.Len() > c.capacity {
		oldest := c.order.Back()
		delete(c.entries, oldest.Value.(*sessionCacheEntry).key)
		c.order.Remove(oldest)
	}
}

type sessionTicketState struct {
	createdAt        int64
	lifetime         uint32
	suite            uint16
	psk              []byte
	serverName       string
	protocol         string
	ageAdd           uint32
	maxEarlyData     uint32
	recordSizeLimit  uint16
	clientAuthAt     int64
	peerCertificates []*x509.Certificate
	verifiedChains   [][]*x509.Certificate
}

func (s *sessionTicketState) marshal() ([]byte, error) {
	if len(s.psk) == 0 || s.lifetime == 0 {
		return nil, errors.New("dtls13: invalid session ticket state")
	}
	if len(s.psk) > 255 {
		return nil, &ProtocolError{"8-bit vector overflow"}
	}
	if len(s.serverName) > 65535 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	if len(s.protocol) > 255 {
		return nil, &ProtocolError{"8-bit vector overflow"}
	}
	if len(s.verifiedChains) > 0 && len(s.peerCertificates) == 0 {
		return nil, errors.New("dtls13: verified client chain without a certificate")
	}
	version := 1
	if s.maxEarlyData != 0 {
		version = 4
	}
	if len(s.peerCertificates) > 0 {
		if s.clientAuthAt == 0 {
			return nil, errors.New("dtls13: client authentication time is missing")
		}
		version = 3
		if s.maxEarlyData != 0 {
			version = 5
		}
	}
	recordSizeLimit := effectiveRecordSizeLimit(s.recordSizeLimit)
	capacity := 2 + 8 + 4 + 4 + 2 + 1 + len(s.psk) + 2 + len(s.serverName) + 1 + len(s.protocol)
	if version >= 2 {
		capacity += 4
	}
	if version >= 4 {
		capacity += 2
	}
	if version == 3 || version == 5 {
		capacity += 8 + 2 + 2
		for _, cert := range s.peerCertificates {
			if cert == nil || len(cert.Raw) == 0 {
				return nil, errors.New("dtls13: invalid client certificate in session ticket")
			}
			capacity += 3 + len(cert.Raw)
		}
		for _, chain := range s.verifiedChains {
			if len(chain) == 0 || chain[0] == nil || !bytes.Equal(chain[0].Raw, s.peerCertificates[0].Raw) {
				return nil, errors.New("dtls13: invalid verified client chain in session ticket")
			}
			capacity += 2
			for _, cert := range chain[1:] {
				if cert == nil || len(cert.Raw) == 0 {
					return nil, errors.New("dtls13: invalid verified client chain in session ticket")
				}
				capacity += 3 + len(cert.Raw)
			}
		}
	}
	w := newWireBuilder(capacity)
	w.u16(version)
	w.b = binary.BigEndian.AppendUint64(w.b, uint64(s.createdAt))
	w.u32(s.lifetime)
	w.u32(s.ageAdd)
	w.u16(int(s.suite))
	w.bytes8(s.psk)
	w.string16(s.serverName)
	w.string8(s.protocol)
	if version >= 2 {
		w.u32(s.maxEarlyData)
	}
	if version >= 4 {
		w.u16(int(recordSizeLimit))
	}
	if version == 3 || version == 5 {
		w.b = binary.BigEndian.AppendUint64(w.b, uint64(s.clientAuthAt))
		certificates := w.startVector16()
		for _, cert := range s.peerCertificates {
			w.bytes24(cert.Raw)
		}
		w.endVector16(certificates)
		chains := w.startVector16()
		for _, chain := range s.verifiedChains {
			encoded := w.startVector16()
			for _, cert := range chain[1:] {
				w.bytes24(cert.Raw)
			}
			w.endVector16(encoded)
		}
		w.endVector16(chains)
	}
	return w.b, w.err
}

func parseSessionTicketState(b []byte) (*sessionTicketState, error) {
	p := wireParser{b: b}
	version := p.u16()
	if version < 1 || version > 5 {
		return nil, errors.New("dtls13: unsupported session ticket version")
	}
	created := p.take(8)
	state := &sessionTicketState{}
	if len(created) == 8 {
		state.createdAt = int64(binary.BigEndian.Uint64(created))
	}
	state.lifetime = p.u32()
	state.ageAdd = p.u32()
	state.suite = uint16(p.u16())
	state.psk = append([]byte(nil), p.bytes8()...)
	state.serverName = string(p.bytes16())
	state.protocol = string(p.bytes8())
	if version >= 2 {
		state.maxEarlyData = uint32(p.u32())
	}
	state.recordSizeLimit = defaultRecordSizeLimit
	if version >= 4 {
		state.recordSizeLimit = uint16(p.u16())
		if state.recordSizeLimit < minRecordSizeLimit || state.recordSizeLimit > defaultRecordSizeLimit {
			return nil, errors.New("dtls13: invalid record size limit in session ticket")
		}
	}
	if version == 3 || version == 5 {
		authenticated := p.take(8)
		if len(authenticated) == 8 {
			state.clientAuthAt = int64(binary.BigEndian.Uint64(authenticated))
		}
		certificates := wireParser{b: p.bytes16()}
		for certificates.off < len(certificates.b) {
			cert, err := x509.ParseCertificate(certificates.bytes24())
			if err != nil {
				return nil, errors.New("dtls13: invalid client certificate in session ticket")
			}
			state.peerCertificates = append(state.peerCertificates, cert)
		}
		if err := certificates.done(); err != nil {
			return nil, err
		}
		chains := wireParser{b: p.bytes16()}
		for chains.off < len(chains.b) {
			if len(state.peerCertificates) == 0 {
				return nil, errors.New("dtls13: verified client chain without a certificate")
			}
			encoded := wireParser{b: chains.bytes16()}
			chain := []*x509.Certificate{state.peerCertificates[0]}
			for encoded.off < len(encoded.b) {
				cert, err := x509.ParseCertificate(encoded.bytes24())
				if err != nil {
					return nil, errors.New("dtls13: invalid verified client chain in session ticket")
				}
				chain = append(chain, cert)
			}
			if err := encoded.done(); err != nil {
				return nil, err
			}
			state.verifiedChains = append(state.verifiedChains, chain)
		}
		if err := chains.done(); err != nil {
			return nil, err
		}
	}
	if err := p.done(); err != nil {
		return nil, err
	}
	if state.lifetime == 0 || len(state.psk) == 0 || ((version == 3 || version == 5) && (state.clientAuthAt == 0 || len(state.peerCertificates) == 0)) {
		return nil, errors.New("dtls13: invalid session ticket state")
	}
	return state, nil
}

type sessionTicketProtector struct {
	aead cipher.AEAD
	rand io.Reader
	now  func() time.Time
}

func newSessionTicketProtector(key [32]byte, random io.Reader, now func() time.Time) (*sessionTicketProtector, error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if random == nil {
		random = rand.Reader
	}
	if now == nil {
		now = time.Now
	}
	return &sessionTicketProtector{aead: aead, rand: random, now: now}, nil
}

var sessionTicketAdditionalData = []byte("DTLS 1.3 session ticket v1")

func ensureSessionTicketKey(config *Config) error {
	if config == nil || config.SessionTicketsDisabled {
		return nil
	}
	var zero [32]byte
	if config.SessionTicketKey != zero {
		return nil
	}
	_, err := io.ReadFull(config.Rand, config.SessionTicketKey[:])
	return err
}

func (p *sessionTicketProtector) seal(state *sessionTicketState) ([]byte, error) {
	plain, err := state.marshal()
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, p.aead.NonceSize())
	if _, err = io.ReadFull(p.rand, nonce); err != nil {
		return nil, err
	}
	return append(nonce, p.aead.Seal(nil, nonce, plain, sessionTicketAdditionalData)...), nil
}

func (p *sessionTicketProtector) open(ticket []byte) (*sessionTicketState, error) {
	if len(ticket) < p.aead.NonceSize()+p.aead.Overhead() {
		return nil, errors.New("dtls13: invalid session ticket")
	}
	nonce := ticket[:p.aead.NonceSize()]
	plain, err := p.aead.Open(nil, nonce, ticket[p.aead.NonceSize():], sessionTicketAdditionalData)
	if err != nil {
		return nil, errors.New("dtls13: invalid session ticket")
	}
	state, err := parseSessionTicketState(plain)
	if err != nil {
		return nil, err
	}
	created := time.Unix(state.createdAt, 0)
	lifetime := time.Duration(state.lifetime) * time.Second
	now := p.now()
	if created.After(now.Add(time.Minute)) || now.Sub(created) > lifetime {
		return nil, errors.New("dtls13: expired session ticket")
	}
	return state, nil
}

func clientSessionCacheKey(config *Config, conn net.Conn) string {
	if config.ServerName != "" {
		return config.ServerName
	}
	return conn.RemoteAddr().String()
}

func usableClientSession(config *Config, conn net.Conn) (*ClientSessionState, *cipherSuite) {
	if config.SessionTicketsDisabled || config.ClientSessionCache == nil {
		return nil, nil
	}
	key := clientSessionCacheKey(config, conn)
	var state *ClientSessionState
	var ok bool
	if cache, atomic := config.ClientSessionCache.(interface {
		Take(string) (*ClientSessionState, bool)
	}); atomic {
		state, ok = cache.Take(key)
	} else {
		state, ok = config.ClientSessionCache.Get(key)
		if ok {
			config.ClientSessionCache.Put(key, nil)
		}
	}
	if !ok || state == nil || len(state.ticket) == 0 || len(state.psk) == 0 {
		return nil, nil
	}
	if validateCertificateSecurityPolicy(state.peerCertificates, true) != nil {
		return nil, nil
	}
	for _, chain := range state.verifiedChains {
		if validateCertificateSecurityPolicy(chain, true) != nil {
			return nil, nil
		}
	}
	age := config.Time().Sub(state.receivedAt)
	if age < 0 || age > time.Duration(state.lifetime)*time.Second || state.serverName != config.ServerName {
		config.ClientSessionCache.Put(key, nil)
		return nil, nil
	}
	if state.protocol != "" {
		found := false
		for _, protocol := range config.NextProtos {
			if protocol == state.protocol {
				found = true
				break
			}
		}
		if !found {
			return nil, nil
		}
	}
	suite, err := cipherSuiteForID(state.suite)
	compatible := false
	if err == nil {
		for _, id := range config.CipherSuites {
			candidate, candidateErr := cipherSuiteForID(id)
			compatible = compatible || (candidateErr == nil && candidate.hash == suite.hash)
		}
	}
	if err != nil || !compatible || len(state.psk) != suite.hash.Size() {
		return nil, nil
	}
	return state, suite
}

func validClientAuthenticationTicket(config *Config, state *sessionTicketState) bool {
	// Match TLS 1.3 resumption: restore the authenticated identity without
	// rerunning VerifyPeerCertificate, but enforce the current built-in policy.
	hasCertificate := len(state.peerCertificates) != 0
	if hasCertificate && config.ClientAuth == tls.NoClientCert {
		return false
	}
	if !hasCertificate {
		// Legacy tickets cannot distinguish an anonymous handshake from client
		// authentication state that their format did not preserve.
		return config.ClientAuth == tls.NoClientCert
	}
	if validateCertificateSecurityPolicy(state.peerCertificates, false) != nil {
		return false
	}
	now := config.Time()
	if state.clientAuthAt == 0 || now.Before(time.Unix(state.clientAuthAt, 0)) || now.Sub(time.Unix(state.clientAuthAt, 0)) > config.SessionTicketLifetime {
		return false
	}
	if config.ClientAuth < tls.VerifyClientCertIfGiven {
		return !now.Before(state.peerCertificates[0].NotBefore) && !now.After(state.peerCertificates[0].NotAfter)
	}
	for _, chain := range state.verifiedChains {
		if len(chain) == 0 {
			continue
		}
		valid := true
		for _, cert := range chain {
			if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
				valid = false
				break
			}
		}
		if !valid {
			continue
		}
		intermediates := x509.NewCertPool()
		for i := 1; i+1 < len(chain); i++ {
			intermediates.AddCert(chain[i])
		}
		verified, err := chain[0].Verify(x509.VerifyOptions{
			Roots: config.ClientCAs, Intermediates: intermediates, CurrentTime: now,
			KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		})
		if err != nil {
			continue
		}
		for _, verifiedChain := range verified {
			if validateCertificateSecurityPolicy(verifiedChain, false) == nil {
				return true
			}
		}
	}
	return false
}

func (c *Conn) acceptSessionTicket(hello *clientHello, clientHelloBody, initialClientHelloHash, hrrBody []byte, suite *cipherSuite, protocol string) (*sessionTicketState, uint16, error) {
	c.earlyMu.Lock()
	c.earlyAccepted = false
	c.earlyDataLimit = 0
	c.earlyMu.Unlock()
	if c.config.SessionTicketsDisabled || len(hello.pskIdentity) == 0 || !hello.pskDHE {
		return nil, 0, nil
	}
	if err := ensureSessionTicketKey(c.config); err != nil {
		return nil, 0, err
	}
	protector, err := newSessionTicketProtector(c.config.SessionTicketKey, c.config.Rand, c.config.Time)
	if err != nil {
		return nil, 0, err
	}
	identities := hello.pskIdentities
	if len(identities) == 0 {
		identities = []pskIdentityEntry{{identity: hello.pskIdentity, obfuscatedAge: hello.obfuscatedAge}}
	}
	var state *sessionTicketState
	selected := -1
	for i, identity := range identities {
		candidate, openErr := protector.open(identity.identity)
		var candidateSuite *cipherSuite
		var suiteErr error
		if openErr == nil {
			candidateSuite, suiteErr = cipherSuiteForID(candidate.suite)
		}
		if openErr == nil && suiteErr == nil && candidateSuite.hash == suite.hash && candidate.serverName == hello.serverName && candidate.protocol == protocol && len(candidate.psk) == suite.hash.Size() && validClientAuthenticationTicket(c.config, candidate) {
			state = candidate
			selected = i
			break
		}
	}
	if selected < 0 || selected > int(^uint16(0)) {
		return nil, 0, nil
	}
	clientAge := identities[selected].obfuscatedAge - state.ageAdd
	actualAge := c.config.Time().Sub(time.Unix(state.createdAt, 0))
	if actualAge < 0 || actualAge > time.Duration(state.lifetime)*time.Second {
		return nil, 0, nil
	}
	actualAgeMillis := uint64(actualAge / time.Millisecond)
	ageWithinTolerance := true
	if actualAgeMillis <= uint64(^uint32(0)) {
		difference := int64(clientAge) - int64(actualAgeMillis)
		if difference < 0 {
			difference = -difference
		}
		if difference > int64(10*time.Minute/time.Millisecond) {
			ageWithinTolerance = false
		}
	}
	if err = verifyClientHelloPSKBinderAt(clientHelloBody, suite, state.psk, initialClientHelloHash, hrrBody, selected); err != nil {
		return nil, 0, err
	}
	if selected == 0 && state.suite == suite.id && ageWithinTolerance && hello.earlyData && state.maxEarlyData > 0 && c.config.MaxEarlyData >= state.maxEarlyData && c.config.RecordSizeLimit >= effectiveRecordSizeLimit(state.recordSizeLimit) && hrrBody == nil && c.config.AllowEarlyDataWithoutCookie {
		limit := state.maxEarlyData
		cache := c.config.EarlyDataReplayCache
		if cache == nil {
			cache = defaultEarlyDataReplayCache
		}
		expires := time.Unix(state.createdAt, 0).Add(time.Duration(state.lifetime) * time.Second)
		if cache != nil && cache.CheckAndStore(string(identities[selected].identity), expires) {
			c.earlyMu.Lock()
			c.earlyAccepted = true
			c.earlyDataLimit = limit
			c.earlyMu.Unlock()
		}
	}
	return state, uint16(selected), nil
}

func obfuscatedTicketAge(state *ClientSessionState, now time.Time) uint32 {
	age := now.Sub(state.receivedAt)
	if age < 0 {
		age = 0
	}
	return uint32(age/time.Millisecond) + state.ageAdd
}

func writeBinderTranscriptMessage(h hash.Hash, typ uint8, body []byte, fullLength int) {
	header := []byte{typ, 0, 0, 0}
	putUint24(header[1:], uint32(fullLength))
	h.Write(header)
	h.Write(body)
}

func truncateClientHelloForBinder(body []byte, binderLength int) ([]byte, error) {
	return truncateClientHelloForBinderEntries(body, 1+binderLength)
}

func truncateClientHelloForBinderEntries(body []byte, binderEntriesLength int) ([]byte, error) {
	suffix := 2 + binderEntriesLength
	if binderEntriesLength < 2 || len(body) < suffix {
		return nil, &ProtocolError{"invalid ClientHello binder truncation"}
	}
	at := len(body) - suffix
	if int(binary.BigEndian.Uint16(body[at:at+2])) != binderEntriesLength {
		return nil, &ProtocolError{"invalid ClientHello binders vector"}
	}
	return body[:at], nil
}

func calculatePSKBinder(suite *cipherSuite, psk, transcriptHash []byte) []byte {
	binder := make([]byte, suite.hash.Size())
	calculatePSKBinderInto(suite, psk, transcriptHash, binder)
	return binder
}

func calculatePSKBinderInto(suite *cipherSuite, psk, transcriptHash, out []byte) {
	if len(out) != suite.hash.Size() {
		panic("dtls13: PSK binder output must equal Hash.length")
	}
	hkdfExtractInto(suite.hash.New, psk, nil, out)
	deriveSecretInto(suite, out, labelResumptionBinder, emptyTranscriptHash(suite), out)
	computeFinishedVerifyDataInto(suite, out, transcriptHash, out)
}

func marshalClientHelloWithPSKBinder(hello *clientHello, suite *cipherSuite, psk, initialClientHelloHash, hrrBody []byte) ([]byte, error) {
	hello.pskBinder = make([]byte, suite.hash.Size())
	body, err := hello.marshal()
	if err != nil {
		return nil, err
	}
	truncated, err := truncateClientHelloForBinder(body, suite.hash.Size())
	if err != nil {
		return nil, err
	}
	h := suite.hash.New()
	if len(hrrBody) > 0 {
		writeBinderTranscriptMessage(h, handshakeTypeMessageHash, initialClientHelloHash, len(initialClientHelloHash))
		writeBinderTranscriptMessage(h, handshakeTypeServerHello, hrrBody, len(hrrBody))
	}
	writeBinderTranscriptMessage(h, handshakeTypeClientHello, truncated, len(body))
	calculatePSKBinderInto(suite, psk, h.Sum(nil), hello.pskBinder)
	return hello.marshal()
}

func verifyClientHelloPSKBinder(body []byte, suite *cipherSuite, psk, initialClientHelloHash, hrrBody []byte) error {
	return verifyClientHelloPSKBinderAt(body, suite, psk, initialClientHelloHash, hrrBody, 0)
}

func verifyClientHelloPSKBinderAt(body []byte, suite *cipherSuite, psk, initialClientHelloHash, hrrBody []byte, selected int) error {
	hello, err := parseClientHello(body)
	if err != nil {
		return err
	}
	if selected < 0 || selected >= len(hello.pskBinders) || len(hello.pskBinders[selected]) != suite.hash.Size() {
		return &ProtocolError{"PSK binder length does not match cipher suite"}
	}
	binderBytes := 0
	for _, binder := range hello.pskBinders {
		binderBytes += 1 + len(binder)
	}
	truncated, err := truncateClientHelloForBinderEntries(body, binderBytes)
	if err != nil {
		return err
	}
	h := suite.hash.New()
	if len(hrrBody) > 0 {
		writeBinderTranscriptMessage(h, handshakeTypeMessageHash, initialClientHelloHash, len(initialClientHelloHash))
		writeBinderTranscriptMessage(h, handshakeTypeServerHello, hrrBody, len(hrrBody))
	}
	writeBinderTranscriptMessage(h, handshakeTypeClientHello, truncated, len(body))
	want := calculatePSKBinder(suite, psk, h.Sum(nil))
	if !hmac.Equal(want, hello.pskBinders[selected]) {
		return alertError(alertDecryptError, &ProtocolError{"invalid PSK binder"})
	}
	return nil
}

func (c *Conn) sendNewSessionTicket(schedule *keySchedule, suite *cipherSuite, serverName, protocol string, clientAuthAt int64, peerCertificates []*x509.Certificate, verifiedChains [][]*x509.Certificate) error {
	if c.config.SessionTicketsDisabled {
		return nil
	}
	if err := ensureSessionTicketKey(c.config); err != nil {
		return err
	}
	nonce := make([]byte, 8)
	if _, err := io.ReadFull(c.config.Rand, nonce); err != nil {
		return err
	}
	psk, err := schedule.resumptionPSK(nonce)
	if err != nil {
		return err
	}
	ageBytes := make([]byte, 4)
	if _, err = io.ReadFull(c.config.Rand, ageBytes); err != nil {
		return err
	}
	ageAdd := binary.BigEndian.Uint32(ageBytes)
	protector, err := newSessionTicketProtector(c.config.SessionTicketKey, c.config.Rand, c.config.Time)
	if err != nil {
		return err
	}
	lifetime := uint32(c.config.SessionTicketLifetime / time.Second)
	ticket, err := protector.seal(&sessionTicketState{
		createdAt: c.config.Time().Unix(), lifetime: lifetime, suite: suite.id, psk: psk,
		serverName: serverName, protocol: protocol, ageAdd: ageAdd, maxEarlyData: c.config.MaxEarlyData, recordSizeLimit: c.localRecordSizeLimit,
		clientAuthAt: clientAuthAt, peerCertificates: peerCertificates, verifiedChains: verifiedChains,
	})
	if err != nil {
		var protocolErr *ProtocolError
		if len(peerCertificates) > 0 && errors.As(err, &protocolErr) && protocolErr.Reason == "16-bit vector overflow" {
			return nil
		}
		return err
	}
	if len(ticket) > 65535 {
		return nil
	}
	body, err := (&newSessionTicketMessage{
		lifetime: lifetime, ageAdd: ageAdd, nonce: nonce, ticket: ticket, maxEarlyData: c.config.MaxEarlyData,
	}).marshal()
	if err != nil {
		return err
	}
	c.writeMu.Lock()
	if c.sendingTraffic == nil {
		c.writeMu.Unlock()
		return &ProtocolError{"application traffic state is not installed"}
	}
	if err = c.sendingTraffic.canAllocateMessageSequences(1); err != nil {
		c.writeMu.Unlock()
		return err
	}
	sequence := c.sendingTraffic.messageSequence
	flight, err := buildProtectedFlight([]handshakeMessage{{
		typ: handshakeTypeNewSessionTicket, sequence: sequence, body: body,
	}}, c.currentMTU(), c.sendCipher)
	if err != nil {
		c.writeMu.Unlock()
		return err
	}
	flight.setIntervals(c.config.FlightInterval, c.config.MaxFlightInterval)
	if err = c.writeFlight(c.conn, flight); err != nil {
		c.writeMu.Unlock()
		return err
	}
	c.sendingTraffic.commitMessageSequences(1)
	c.ticketFlight = flight
	c.writeMu.Unlock()
	c.startTicketRetransmission()
	return nil
}

func (c *Conn) ackProtectedRecord(number recordNumber) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	var ackScratch [1][]byte
	acks, _, err := buildACKRecordsInto(ackScratch[:0], []recordNumber{number}, c.currentMTU(), 0, c.sendCipher)
	if err != nil {
		return err
	}
	for _, wire := range acks {
		if err = c.writeRecord(wire); err != nil {
			return err
		}
	}
	return nil
}

func (c *Conn) processNewSessionTicket(body []byte) error {
	message, err := parseNewSessionTicket(body)
	if err != nil {
		return err
	}
	if message.lifetime == 0 || time.Duration(message.lifetime)*time.Second > maxTicketLifetime {
		return nil
	}
	if c.config.SessionTicketsDisabled || c.config.ClientSessionCache == nil || c.resumptionSuite == nil || len(c.resumptionMasterSecret) == 0 {
		return nil
	}
	psk := deriveResumptionPSK(c.resumptionSuite, c.resumptionMasterSecret, message.nonce)
	state := c.ConnectionState()
	c.config.ClientSessionCache.Put(clientSessionCacheKey(c.config, c.conn), &ClientSessionState{
		ticket: append([]byte(nil), message.ticket...), psk: psk, nonce: append([]byte(nil), message.nonce...),
		suite: c.resumptionSuite.id, receivedAt: c.config.Time(), lifetime: message.lifetime,
		ageAdd: message.ageAdd, serverName: c.config.ServerName, protocol: state.NegotiatedProtocol,
		maxEarlyData: message.maxEarlyData, recordSizeLimit: state.PeerRecordSizeLimit, peerCertificates: state.PeerCertificates, verifiedChains: state.VerifiedChains,
	})
	return nil
}

func (c *Conn) startTicketRetransmission() {
	go func() {
		interval := c.flightInterval()
		if interval <= 0 {
			interval = time.Second
		}
		max := c.config.MaxFlightInterval
		if max < interval {
			max = interval
		}
		timer := time.NewTimer(interval)
		defer timer.Stop()
		timeoutCount := 0
		for {
			<-timer.C
			c.readerMu.Lock()
			closed := c.readerClosed
			c.readerMu.Unlock()
			if closed {
				return
			}
			c.writeMu.Lock()
			flight := c.ticketFlight
			if flight == nil || flight.complete() {
				c.ticketFlight = nil
				c.writeMu.Unlock()
				return
			}
			timeoutCount++
			resized, err := c.prepareFlightRetransmission(flight, timeoutCount)
			if err == nil && resized {
				err = c.writeFlight(c.conn, flight)
			} else if err == nil {
				err = c.retransmitFlight(c.conn, flight)
			}
			c.writeMu.Unlock()
			if err != nil {
				c.failConnection(err)
				return
			}
			if interval < max {
				interval *= 2
				if interval > max {
					interval = max
				}
			}
			timer.Reset(interval)
		}
	}()
}
