package dtls13

import (
	"bytes"
	"container/list"
	"crypto/sha256"
	"sync"
)

const (
	extCachedInfo            uint16 = 25
	cachedCert               uint8  = 1
	cachedCertRequest        uint8  = 2
	maxCachedInformationSize        = 64 << 10
)

// CachedInformationCache retains authenticated server Certificate and initial
// CertificateRequest messages for RFC 7924. It is safe to share across clients.
// Cache hits still undergo current certificate and signature verification.
// Use NewCachedInformationCache to create a cache; the zero value is disabled.
// A cache must not be copied after first use.
type CachedInformationCache struct {
	mu       sync.Mutex
	capacity int
	entries  map[string]*list.Element
	order    list.List
}

type cachedInformationEntry struct {
	key    string
	bodies [2][]byte
}

// NewCachedInformationCache creates an LRU cache of server entries. Capacity
// below one selects 64. Each entry retains at most two 64 KiB messages; larger
// messages are not cached. Sharing a cache can make connections linkable via
// the fingerprints advertised in ClientHello. ECH keeps them in the inner hello.
func NewCachedInformationCache(capacity int) *CachedInformationCache {
	if capacity < 1 {
		capacity = 64
	}
	return &CachedInformationCache{capacity: capacity, entries: make(map[string]*list.Element)}
}

func (c *CachedInformationCache) get(key string) [2][]byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e := c.entries[key]; e != nil {
		c.order.MoveToFront(e)
		return e.Value.(*cachedInformationEntry).bodies
	}
	return [2][]byte{}
}

func (c *CachedInformationCache) put(key string, bodies [2][]byte) {
	for i, body := range bodies {
		if len(body) > maxCachedInformationSize {
			bodies[i] = nil
		} else {
			bodies[i] = bytes.Clone(body)
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.capacity == 0 {
		return
	}
	if e := c.entries[key]; e != nil {
		e.Value.(*cachedInformationEntry).bodies = bodies
		c.order.MoveToFront(e)
		return
	}
	c.entries[key] = c.order.PushFront(&cachedInformationEntry{key: key, bodies: bodies})
	if c.order.Len() > c.capacity {
		e := c.order.Back()
		delete(c.entries, e.Value.(*cachedInformationEntry).key)
		c.order.Remove(e)
	}
}

// RFC 7924 Appendix A hashes the TLS Handshake header as well as the body.
// DTLS sequence and fragmentation fields are excluded, as in the transcript.
func cachedInformationHash(typ uint8, body []byte) [sha256.Size]byte {
	t := newTranscriptHash(sha256.New())
	_ = t.add(typ, 0, body)
	var hash [sha256.Size]byte
	t.sumInto(hash[:0])
	return hash
}

type cachedInformationOffer struct {
	typ  uint8
	hash []byte
}

func parseCachedInformationOffer(raw []byte) ([]cachedInformationOffer, error) {
	p := wireParser{b: raw}
	objects := wireParser{b: p.bytes16()}
	if err := p.done(); err != nil {
		return nil, err
	}
	if len(objects.b) == 0 {
		return nil, alertError(alertDecodeError, &ProtocolError{"empty cached_info offer"})
	}
	var offers []cachedInformationOffer
	for objects.off < len(objects.b) && objects.err == nil {
		typ := uint8(objects.u8())
		hash := objects.bytes8()
		if len(hash) == 0 || ((typ == cachedCert || typ == cachedCertRequest) && len(hash) != sha256.Size) {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid cached_info fingerprint length"})
		}
		if typ == cachedCert || typ == cachedCertRequest {
			offers = append(offers, cachedInformationOffer{typ, hash})
		}
	}
	return offers, objects.done()
}

func validateCachedInformationSelection(hello *clientHello, raw []byte) (uint8, error) {
	offer, ok := hello.unknownExtensions[extCachedInfo]
	if !ok {
		return 0, alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited cached_info response"})
	}
	offers, err := parseCachedInformationOffer(offer)
	if err != nil {
		return 0, err
	}
	p := wireParser{b: raw}
	types := p.bytes16()
	if err = p.done(); err != nil {
		return 0, err
	}
	if len(types) == 0 {
		return 0, alertError(alertDecodeError, &ProtocolError{"empty cached_info response"})
	}
	var selected uint8
	for _, typ := range types {
		found := false
		for _, object := range offers {
			found = found || object.typ == typ
		}
		if !found {
			return 0, alertError(alertIllegalParameter, &ProtocolError{"unoffered cached_info type"})
		}
		selected |= typ
	}
	return selected, nil
}

type clientCachedInformation struct {
	offered  [2][]byte
	received [2][]byte
	selected uint8
}

func (c *Conn) offerCachedInformation(s *clientHandshakeState, hello *clientHello) {
	if c.config.CachedInformationCache == nil {
		return
	}
	s.cachedInfo = &clientCachedInformation{offered: c.config.CachedInformationCache.get(clientSessionCacheKey(c.config, c.conn))}
	w := newWireBuilder(70)
	start := w.startVector16()
	for i, typ := range []uint8{handshakeTypeCertificate, handshakeTypeCertificateRequest} {
		body := s.cachedInfo.offered[i]
		if len(body) == 0 || len(body) > c.config.MaxHandshakeMessage {
			continue
		}
		hash := cachedInformationHash(typ, body)
		w.u8(i + 1)
		w.bytes8(hash[:])
	}
	w.endVector16(start)
	if len(w.b) > 2 {
		if hello.unknownExtensions == nil {
			hello.unknownExtensions = make(map[uint16][]byte)
		}
		hello.unknownExtensions[extCachedInfo] = w.b
	}
}

func (s *clientCachedInformation) resolve(typ uint8, body []byte, maxSize int) ([]byte, error) {
	if s == nil {
		return body, nil
	}
	i := 0
	if typ == handshakeTypeCertificateRequest {
		i = 1
	}
	if s.selected&uint8(i+1) != 0 {
		if len(body) != 1+sha256.Size || body[0] != sha256.Size {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid cached handshake fingerprint"})
		}
		cached := s.offered[i]
		hash := cachedInformationHash(typ, cached)
		if len(cached) == 0 || len(cached) > maxSize || !bytes.Equal(body[1:], hash[:]) {
			return nil, alertError(alertIllegalParameter, &ProtocolError{"cached handshake fingerprint mismatch"})
		}
		body = cached
	}
	if len(body) <= maxCachedInformationSize {
		s.received[i] = bytes.Clone(body)
	}
	return body, nil
}

func cachedHandshakeBody(typ uint8, body []byte, offers []cachedInformationOffer) ([]byte, uint8) {
	if len(offers) == 0 {
		return body, 0
	}
	kind := cachedCert
	if typ == handshakeTypeCertificateRequest {
		kind = cachedCertRequest
	}
	hash := cachedInformationHash(typ, body)
	for _, offer := range offers {
		if offer.typ == kind && bytes.Equal(offer.hash, hash[:]) {
			return append([]byte{sha256.Size}, hash[:]...), kind
		}
	}
	return body, 0
}

func (c *Conn) parseCachedServerCertificate(s *clientHandshakeState, message completedHandshake) (*certificateMessage, error) {
	algorithms := s.hello.certificateCompressionAlgorithms()
	if s.cachedInfo == nil {
		return parseCertificateHandshakeMessage(message.typ, message.body, algorithms, c.config.MaxHandshakeMessage)
	}
	body := message.body
	if message.typ == handshakeTypeCompressedCertificate {
		if s.cachedInfo.selected&cachedCert != 0 {
			return nil, alertError(alertUnexpectedMessage, &ProtocolError{"compressed Certificate after cached_info selection"})
		}
		decoded, pooled, err := decompressCertificate(body, algorithms, c.config.MaxHandshakeMessage)
		if err != nil {
			return nil, err
		}
		if pooled {
			defer releaseCertificateDecompressionBuffer(decoded)
		}
		body = decoded
	}
	body, err := s.cachedInfo.resolve(handshakeTypeCertificate, body, c.config.MaxHandshakeMessage)
	if err != nil {
		return nil, err
	}
	return parseCertificateMessage(body, c.config.MaxHandshakeMessage)
}
