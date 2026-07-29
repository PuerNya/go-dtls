package dtls13

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"time"
)

const cookieMACSize = sha256.Size

type cookieKey struct {
	id     byte
	secret []byte
}
type cookieProtector struct {
	mu       sync.Mutex
	current  cookieKey
	previous *cookieKey
	lifetime time.Duration
	now      func() time.Time
	rand     io.Reader
	rotated  time.Time
}

func newCookieProtector(id byte, secret []byte, lifetime time.Duration, now func() time.Time) (*cookieProtector, error) {
	if len(secret) < 32 {
		return nil, errors.New("dtls13: cookie secret must be at least 32 bytes")
	}
	if lifetime <= 0 {
		return nil, errors.New("dtls13: cookie lifetime must be positive")
	}
	if now == nil {
		now = time.Now
	}
	return &cookieProtector{current: cookieKey{id: id, secret: append([]byte(nil), secret...)}, lifetime: lifetime, now: now, rotated: now()}, nil
}
func (p *cookieProtector) rotate(id byte, secret []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.rotateLocked(id, secret)
}
func (p *cookieProtector) rotateLocked(id byte, secret []byte) error {
	if len(secret) < 32 {
		return errors.New("dtls13: cookie secret must be at least 32 bytes")
	}
	old := p.current
	p.previous = &old
	p.current = cookieKey{id: id, secret: append([]byte(nil), secret...)}
	p.rotated = p.now()
	return nil
}

func (p *cookieProtector) rotateIfNeededLocked() error {
	if p.rand == nil || p.now().Sub(p.rotated) < p.lifetime/2 {
		return nil
	}
	secret := make([]byte, 32)
	if _, err := io.ReadFull(p.rand, secret); err != nil {
		return err
	}
	return p.rotateLocked(p.current.id+1, secret)
}
func cookieMAC(key, address, prefix []byte) []byte {
	m := hmac.New(sha256.New, key)
	m.Write(address)
	m.Write(prefix)
	return m.Sum(nil)
}
func (p *cookieProtector) seal(address, clientHelloHash []byte) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.rotateIfNeededLocked(); err != nil {
		return nil, err
	}
	if len(clientHelloHash) == 0 || len(clientHelloHash) > 255 {
		return nil, errors.New("dtls13: invalid ClientHello hash length")
	}
	prefix := make([]byte, 10+len(clientHelloHash))
	prefix[0] = p.current.id
	binary.BigEndian.PutUint64(prefix[1:9], uint64(p.now().Unix()))
	prefix[9] = byte(len(clientHelloHash))
	copy(prefix[10:], clientHelloHash)
	return append(prefix, cookieMAC(p.current.secret, address, prefix)...), nil
}
func (p *cookieProtector) open(address, cookie []byte) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.rotateIfNeededLocked(); err != nil {
		return nil, err
	}
	if len(cookie) < 10+cookieMACSize {
		return nil, errors.New("dtls13: invalid cookie")
	}
	hashLen := int(cookie[9])
	if hashLen == 0 || len(cookie) != 10+hashLen+cookieMACSize {
		return nil, errors.New("dtls13: invalid cookie")
	}
	var key []byte
	if cookie[0] == p.current.id {
		key = p.current.secret
	} else if p.previous != nil && cookie[0] == p.previous.id {
		key = p.previous.secret
	} else {
		return nil, errors.New("dtls13: invalid cookie")
	}
	prefix, mac := cookie[:10+hashLen], cookie[10+hashLen:]
	if !hmac.Equal(cookieMAC(key, address, prefix), mac) {
		return nil, errors.New("dtls13: invalid cookie")
	}
	created := time.Unix(int64(binary.BigEndian.Uint64(cookie[1:9])), 0)
	now := p.now()
	if created.After(now.Add(time.Minute)) || now.Sub(created) > p.lifetime {
		return nil, errors.New("dtls13: expired cookie")
	}
	return append([]byte(nil), cookie[10:10+hashLen]...), nil
}
