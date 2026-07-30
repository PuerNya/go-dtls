package dtls13

import (
	"bytes"
	"sync"
	"testing"
	"time"
)

func TestCookieBindsAddressAndHash(t *testing.T) {
	now := time.Unix(1000, 0)
	p, err := newCookieProtector(1, bytes.Repeat([]byte{1}, 32), time.Minute, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	hash := bytes.Repeat([]byte{7}, 32)
	cookie, err := p.seal([]byte("192.0.2.1:1234"), hash)
	if err != nil {
		t.Fatal(err)
	}
	got, err := p.open([]byte("192.0.2.1:1234"), cookie)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, hash) {
		t.Fatal("hash mismatch")
	}
	if _, err = p.open([]byte("192.0.2.2:1234"), cookie); err == nil {
		t.Fatal("cookie accepted for another address")
	}
	cookie[len(cookie)-1] ^= 1
	if _, err = p.open([]byte("192.0.2.1:1234"), cookie); err == nil {
		t.Fatal("tampered cookie accepted")
	}
}

func TestCookieAutomaticRotationAndConcurrentUse(t *testing.T) {
	now := time.Unix(1000, 0)
	p, err := newCookieProtector(1, bytes.Repeat([]byte{1}, 32), time.Minute, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	p.rand = bytes.NewReader(bytes.Repeat([]byte{2}, 32))
	address, hash := []byte("peer"), []byte("hash")
	old, err := p.seal(address, hash)
	if err != nil {
		t.Fatal(err)
	}
	now = now.Add(31 * time.Second)
	current, err := p.seal(address, hash)
	if err != nil {
		t.Fatal(err)
	}
	if old[0] == current[0] {
		t.Fatal("cookie key did not rotate automatically")
	}
	if _, err = p.open(address, old); err != nil {
		t.Fatalf("previous cookie rejected during overlap: %v", err)
	}

	var wg sync.WaitGroup
	for range 32 {
		wg.Go(func() {
			for range 100 {
				cookie, sealErr := p.seal(address, hash)
				if sealErr != nil {
					t.Errorf("seal: %v", sealErr)
					return
				}
				if _, openErr := p.open(address, cookie); openErr != nil {
					t.Errorf("open: %v", openErr)
					return
				}
			}
		})
	}
	wg.Wait()
}

func TestNormalizedConfigSharesCookieProtector(t *testing.T) {
	config, err := (&Config{}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	if err = ensureCookieProtector(config); err != nil {
		t.Fatal(err)
	}
	clone, err := config.normalized()
	if err != nil {
		t.Fatal(err)
	}
	if clone.state != config.state {
		t.Fatal("normalized listener connection did not share cookie keys")
	}
}
func TestCookieRotationAndExpiry(t *testing.T) {
	now := time.Unix(1000, 0)
	p, _ := newCookieProtector(1, bytes.Repeat([]byte{1}, 32), time.Minute, func() time.Time { return now })
	old, _ := p.seal([]byte("a"), []byte("hash"))
	if err := p.rotate(2, bytes.Repeat([]byte{2}, 32)); err != nil {
		t.Fatal(err)
	}
	if _, err := p.open([]byte("a"), old); err != nil {
		t.Fatal("previous key rejected during overlap")
	}
	now = now.Add(2 * time.Minute)
	if _, err := p.open([]byte("a"), old); err == nil {
		t.Fatal("expired cookie accepted")
	}
}
