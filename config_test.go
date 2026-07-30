package dtls13

import (
	"crypto/tls"
	"slices"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	c, err := (*Config)(nil).normalized()
	if err != nil {
		t.Fatal(err)
	}
	if c.MTU != 1200 || c.RecordSizeLimit != defaultRecordSizeLimit || c.ReplayWindow != 64 || c.MaxHandshakeMessage != 1<<20 || c.MaxBufferedHandshakeMessages != 8 || c.MaxBufferedHandshakeBytes != 4<<20 || c.MaxBufferedApplicationData != 1<<20 || c.MaxBufferedApplicationDatagrams != 1024 || c.MaxConnectionIDs != 8 {
		t.Fatalf("unexpected defaults: %#v", c)
	}
	if len(c.CurvePreferences) != 2 || c.CurvePreferences[0] != tls.X25519 {
		t.Fatalf("unexpected curve defaults: %v", c.CurvePreferences)
	}
	if len(c.CipherSuites) != 4 || c.CipherSuites[2] != TLS_CHACHA20_POLY1305_SHA256 || c.CipherSuites[3] != TLS_AES_128_CCM_SHA256 {
		t.Fatalf("unexpected cipher suite defaults: %x", c.CipherSuites)
	}
	if cid, offered := c.clientConnectionIDOffer(); !offered || cid == nil || len(cid) != 0 {
		t.Fatalf("default client CID offer=%x, %v", cid, offered)
	}
}

func TestConfigRecordSizeLimitRange(t *testing.T) {
	for _, limit := range []uint16{minRecordSizeLimit, defaultRecordSizeLimit} {
		config, err := (&Config{RecordSizeLimit: limit}).normalized()
		if err != nil || config.RecordSizeLimit != limit {
			t.Fatalf("RecordSizeLimit=%d normalized to %d, %v", limit, config.RecordSizeLimit, err)
		}
	}
	for _, limit := range []uint16{1, minRecordSizeLimit - 1, defaultRecordSizeLimit + 1} {
		if _, err := (&Config{RecordSizeLimit: limit}).normalized(); err == nil {
			t.Fatalf("accepted RecordSizeLimit=%d", limit)
		}
	}
}

func TestConfigCanDisableDefaultClientConnectionIDOffer(t *testing.T) {
	c, err := (&Config{DisableConnectionID: true}).normalized()
	if err != nil {
		t.Fatal(err)
	}
	if cid, offered := c.clientConnectionIDOffer(); offered || cid != nil {
		t.Fatalf("disabled client CID offer=%x, %v", cid, offered)
	}
}

func TestConfigRejectsInvalidCipherSuites(t *testing.T) {
	for _, suites := range [][]uint16{{0xffff}, {TLS_AES_128_CCM_8_SHA256}, {TLS_AES_128_GCM_SHA256, TLS_AES_128_GCM_SHA256}} {
		if _, err := (&Config{CipherSuites: suites}).normalized(); err == nil {
			t.Fatalf("accepted invalid cipher suites: %x", suites)
		}
	}
}

func TestConfigRejectsUnsupportedCurve(t *testing.T) {
	if _, err := (&Config{CurvePreferences: []tls.CurveID{tls.CurveP521}}).normalized(); err == nil {
		t.Fatal("accepted unsupported curve")
	}
	if _, err := (&Config{CurvePreferences: []tls.CurveID{tls.X25519, tls.X25519}}).normalized(); err == nil {
		t.Fatal("accepted duplicate curve")
	}
}

func TestConfigAcceptsHybridKeyExchangeGroups(t *testing.T) {
	groups := []tls.CurveID{tls.X25519MLKEM768, tls.SecP256r1MLKEM768, tls.SecP384r1MLKEM1024, tls.CurveP384}
	config, err := (&Config{CurvePreferences: groups}).normalized()
	if err != nil || !slices.Equal(config.CurvePreferences, groups) {
		t.Fatalf("hybrid groups = %v, %v", config.CurvePreferences, err)
	}
}
func TestConfigRejectsSmallMTU(t *testing.T) {
	if _, err := (&Config{MTU: 255}).normalized(); err == nil {
		t.Fatal("accepted small MTU")
	}
}

func TestConfigRejectsInvalidApplicationBuffer(t *testing.T) {
	if _, err := (&Config{MaxBufferedApplicationData: -1}).normalized(); err == nil {
		t.Fatal("accepted negative application buffer limit")
	}
	if _, err := (&Config{MaxBufferedApplicationDatagrams: -1}).normalized(); err == nil {
		t.Fatal("accepted negative application datagram limit")
	}
}

func TestConfigRejectsInvalidConnectionIDLimit(t *testing.T) {
	for _, limit := range []int{-1, 256} {
		if _, err := (&Config{MaxConnectionIDs: limit}).normalized(); err == nil {
			t.Fatalf("accepted MaxConnectionIDs=%d", limit)
		}
	}
}
