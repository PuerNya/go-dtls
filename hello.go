package dtls13

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
)

const (
	extServerName              uint16 = 0
	extMaxFragmentLength       uint16 = 1
	extSupportedGroups         uint16 = 10
	extSupportedVersions       uint16 = 43
	extKeyShare                uint16 = 51
	extSignatureAlgorithms     uint16 = 13
	extALPN                    uint16 = 16
	extPadding                 uint16 = 21
	extCompressCertificate     uint16 = 27
	extRecordSizeLimit         uint16 = 28
	extPreSharedKey            uint16 = 41
	extCookie                  uint16 = 44
	extPSKKeyExchangeModes     uint16 = 45
	extPostHandshakeAuth       uint16 = 49
	extSignatureAlgorithmsCert uint16 = 50
	extConnectionID            uint16 = 54
	extReturnRoutability       uint16 = 61
)

func knownExtensionType(typ uint16) bool {
	switch typ {
	case extServerName, extSupportedGroups, extSupportedVersions, extKeyShare,
		extSignatureAlgorithms, extALPN, extPadding, extCompressCertificate, extRecordSizeLimit, extPreSharedKey, extEarlyData,
		extCookie, extPSKKeyExchangeModes, extPostHandshakeAuth,
		extSignatureAlgorithmsCert, extConnectionID, extReturnRoutability,
		extECH, extECHOuterExtensions:
		return true
	default:
		return false
	}
}

type keyShareEntry struct {
	group tls.CurveID
	data  []byte
}
type pskIdentityEntry struct {
	identity      []byte
	obfuscatedAge uint32
}
type clientHello struct {
	random                        [32]byte
	sessionID                     []byte
	encryptedClientHelloExtension []byte
	cookie                        []byte
	cipherSuites                  []uint16
	keyShares                     []keyShareEntry
	keyShareStorage               [1]keyShareEntry
	signatureSchemes              []tls.SignatureScheme
	certificateSignatureSchemes   []tls.SignatureScheme
	supportedGroups               []tls.CurveID
	serverName                    string
	alpn                          []string
	pskIdentity                   []byte
	obfuscatedAge                 uint32
	pskBinder                     []byte
	pskIdentities                 []pskIdentityEntry
	pskBinders                    [][]byte
	pskDHE                        bool
	earlyData                     bool
	connectionID                  []byte
	hasConnectionID               bool
	returnRoutability             bool
	postHandshakeAuth             bool
	recordSizeLimit               uint16
	hasRecordSizeLimit            bool
	certificateCompressionOffered bool
	unknownExtensions             map[uint16][]byte
}

func (h *clientHello) encryptedClientHello() []byte {
	if h == nil {
		return nil
	}
	return h.encryptedClientHelloExtension
}

func (h *clientHello) setEncryptedClientHello(value []byte) {
	h.encryptedClientHelloExtension = value
}

type serverHello struct {
	random            [32]byte
	sessionID         []byte
	cipherSuite       uint16
	keyShare          keyShareEntry
	selectedIdentity  *uint16
	connectionID      []byte
	hasConnectionID   bool
	returnRoutability bool
}

func validateServerHelloConnectionID(hello *clientHello, sh *serverHello) error {
	if sh != nil && sh.hasConnectionID && (hello == nil || !hello.hasConnectionID) {
		return alertError(alertUnsupportedExtension, &ProtocolError{"server negotiated an unoffered connection ID"})
	}
	return nil
}

func validateServerHelloReturnRoutability(hello *clientHello, sh *serverHello) error {
	if sh == nil || !sh.returnRoutability {
		return nil
	}
	if hello == nil || !hello.returnRoutability {
		return alertError(alertUnsupportedExtension, &ProtocolError{"server negotiated an unoffered return routability check"})
	}
	if !hello.hasConnectionID || !sh.hasConnectionID {
		return alertError(alertIllegalParameter, &ProtocolError{"return routability check requires connection ID"})
	}
	return nil
}

type wireBuilder struct {
	b   []byte
	err error
}

func newWireBuilder(capacity int) wireBuilder {
	return wireBuilder{b: make([]byte, 0, capacity)}
}

func (w *wireBuilder) u8(v int) {
	if w.err != nil {
		return
	}
	if v < 0 || v > 255 {
		w.err = &ProtocolError{"8-bit vector overflow"}
		return
	}
	w.b = append(w.b, byte(v))
}
func (w *wireBuilder) u16(v int) {
	if w.err != nil {
		return
	}
	if v < 0 || v > 65535 {
		w.err = &ProtocolError{"16-bit vector overflow"}
		return
	}
	w.b = append(w.b, byte(v>>8), byte(v))
}
func (w *wireBuilder) u24(v int) {
	if w.err != nil {
		return
	}
	if v < 0 || v >= 1<<24 {
		w.err = &ProtocolError{"24-bit vector overflow"}
		return
	}
	w.b = append(w.b, byte(v>>16), byte(v>>8), byte(v))
}
func (w *wireBuilder) u32(v uint32) {
	if w.err != nil {
		return
	}
	w.b = append(w.b, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}
func (w *wireBuilder) bytes8(v []byte) {
	w.u8(len(v))
	if w.err == nil {
		w.b = append(w.b, v...)
	}
}
func (w *wireBuilder) bytes16(v []byte) {
	w.u16(len(v))
	if w.err == nil {
		w.b = append(w.b, v...)
	}
}
func (w *wireBuilder) bytes24(v []byte) {
	w.u24(len(v))
	if w.err == nil {
		w.b = append(w.b, v...)
	}
}

func (w *wireBuilder) string8(v string) {
	w.u8(len(v))
	if w.err == nil {
		w.b = append(w.b, v...)
	}
}

func (w *wireBuilder) string16(v string) {
	w.u16(len(v))
	if w.err == nil {
		w.b = append(w.b, v...)
	}
}

func (w *wireBuilder) startVector16() int {
	if w.err != nil {
		return -1
	}
	start := len(w.b)
	w.b = append(w.b, 0, 0)
	return start
}

func (w *wireBuilder) endVector16(start int) {
	if start < 0 {
		return
	}
	if w.err != nil {
		w.b = w.b[:start]
		return
	}
	length := len(w.b) - start - 2
	if length > 65535 {
		w.b = w.b[:start]
		w.err = &ProtocolError{"16-bit vector overflow"}
		return
	}
	binary.BigEndian.PutUint16(w.b[start:start+2], uint16(length))
}

func (w *wireBuilder) startVector24() int {
	if w.err != nil {
		return -1
	}
	start := len(w.b)
	w.b = append(w.b, 0, 0, 0)
	return start
}

func (w *wireBuilder) endVector24(start int) {
	if start < 0 {
		return
	}
	if w.err != nil {
		w.b = w.b[:start]
		return
	}
	length := len(w.b) - start - 3
	if length >= 1<<24 {
		w.b = w.b[:start]
		w.err = &ProtocolError{"24-bit vector overflow"}
		return
	}
	putUint24(w.b[start:start+3], uint32(length))
}

type wireParser struct {
	b   []byte
	off int
	err error
}

func (p *wireParser) take(n int) []byte {
	if p.err != nil {
		return nil
	}
	if n < 0 || n > len(p.b)-p.off {
		p.err = &ProtocolError{"truncated handshake message"}
		return nil
	}
	v := p.b[p.off : p.off+n]
	p.off += n
	return v
}
func (p *wireParser) u8() int {
	b := p.take(1)
	if b == nil {
		return 0
	}
	return int(b[0])
}
func (p *wireParser) u16() int {
	b := p.take(2)
	if b == nil {
		return 0
	}
	return int(binary.BigEndian.Uint16(b))
}
func (p *wireParser) u24() int {
	b := p.take(3)
	if b == nil {
		return 0
	}
	return int(b[0])<<16 | int(b[1])<<8 | int(b[2])
}
func (p *wireParser) u32() uint32 {
	b := p.take(4)
	if b == nil {
		return 0
	}
	return binary.BigEndian.Uint32(b)
}
func (p *wireParser) bytes8() []byte  { return p.take(p.u8()) }
func (p *wireParser) bytes16() []byte { return p.take(p.u16()) }
func (p *wireParser) bytes24() []byte { return p.take(p.u24()) }
func (p *wireParser) done() error {
	if p.err != nil {
		return p.err
	}
	if p.off != len(p.b) {
		return &ProtocolError{"trailing handshake data"}
	}
	return nil
}

func marshalKeyShares(shares []keyShareEntry, vector bool) ([]byte, error) {
	length := 0
	for _, share := range shares {
		if len(share.data) > 65535 {
			return nil, &ProtocolError{"16-bit vector overflow"}
		}
		length += 4 + len(share.data)
	}
	if vector && length > 65535 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	capacity := length
	if vector {
		capacity += 2
	}
	out := newWireBuilder(capacity)
	appendEntries := func(w *wireBuilder) {
		for _, s := range shares {
			w.u16(int(s.group))
			w.bytes16(s.data)
		}
	}
	if vector {
		start := out.startVector16()
		appendEntries(&out)
		out.endVector16(start)
	} else {
		appendEntries(&out)
	}
	return out.b, out.err
}

func marshalSignatureSchemes(schemes []tls.SignatureScheme) ([]byte, error) {
	if len(schemes) == 0 {
		schemes = defaultSignatureSchemes()
	}
	if len(schemes) > 65535/2 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	out := newWireBuilder(2 + 2*len(schemes))
	start := out.startVector16()
	for _, scheme := range schemes {
		out.u16(int(scheme))
	}
	out.endVector16(start)
	return out.b, out.err
}

var defaultSignatureSchemeStorage = [...]tls.SignatureScheme{
	tls.Ed25519,
	tls.ECDSAWithP256AndSHA256,
	tls.PSSWithSHA256,
	tls.ECDSAWithP384AndSHA384,
	tls.PSSWithSHA384,
	tls.ECDSAWithP521AndSHA512,
	tls.PSSWithSHA512,
	// PKCS#1 schemes are valid for certificate signatures, but are
	// never selected by selectSignatureScheme for CertificateVerify.
	tls.PKCS1WithSHA256,
	tls.PKCS1WithSHA384,
	tls.PKCS1WithSHA512,
}

// defaultSignatureSchemes returns shared immutable storage used only by
// internal handshake state.
func defaultSignatureSchemes() []tls.SignatureScheme {
	return defaultSignatureSchemeStorage[:]
}
func parseSignatureSchemes(b []byte) ([]tls.SignatureScheme, error) {
	p := wireParser{b: b}
	raw := p.bytes16()
	if err := p.done(); err != nil {
		return nil, err
	}
	if len(raw) < 2 || len(raw)%2 != 0 {
		return nil, alertError(alertDecodeError, &ProtocolError{"invalid signature_algorithms extension length"})
	}
	out := make([]tls.SignatureScheme, 0, len(raw)/2)
	for len(raw) > 0 {
		out = append(out, tls.SignatureScheme(binary.BigEndian.Uint16(raw)))
		raw = raw[2:]
	}
	return out, nil
}
func parseClientSupportedVersions(b []byte) (bool, error) {
	if len(b) < 3 || int(b[0]) != len(b)-1 || (len(b)-1)%2 != 0 {
		return false, alertError(alertDecodeError, &ProtocolError{"invalid supported_versions extension length"})
	}
	for b = b[1:]; len(b) > 0; b = b[2:] {
		if binary.BigEndian.Uint16(b[:2]) == VersionDTLS13 {
			return true, nil
		}
	}
	return false, nil
}

func offeredDTLS13(b []byte) bool {
	offered, err := parseClientSupportedVersions(b)
	return err == nil && offered
}

func marshalSupportedGroups(groups []tls.CurveID, shares []keyShareEntry) ([]byte, error) {
	if len(groups) == 0 {
		for _, share := range shares {
			groups = append(groups, share.group)
		}
	}
	if len(groups) == 0 {
		return nil, &ProtocolError{"empty supported_groups"}
	}
	if len(groups) > 65535/2 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	out := newWireBuilder(2 + 2*len(groups))
	start := out.startVector16()
	for _, group := range groups {
		out.u16(int(group))
	}
	out.endVector16(start)
	return out.b, out.err
}
func parseSupportedGroups(b []byte) ([]tls.CurveID, error) {
	p := wireParser{b: b}
	raw := p.bytes16()
	if err := p.done(); err != nil {
		return nil, err
	}
	if len(raw) < 2 || len(raw)%2 != 0 {
		return nil, alertError(alertDecodeError, &ProtocolError{"invalid supported_groups extension length"})
	}
	out := make([]tls.CurveID, 0, len(raw)/2)
	for len(raw) > 0 {
		out = append(out, tls.CurveID(binary.BigEndian.Uint16(raw)))
		raw = raw[2:]
	}
	return out, nil
}
func marshalServerName(name string) ([]byte, error) {
	if name == "" {
		return nil, nil
	}
	if len(name) > 65535-3 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	out := newWireBuilder(2 + 1 + 2 + len(name))
	start := out.startVector16()
	out.u8(0)
	out.string16(name)
	out.endVector16(start)
	return out.b, out.err
}
func parseServerName(b []byte) (string, error) {
	p := wireParser{b: b}
	raw := p.bytes16()
	if err := p.done(); err != nil {
		return "", err
	}
	q := wireParser{b: raw}
	if q.u8() != 0 {
		return "", &ProtocolError{"unsupported server name type"}
	}
	name := q.bytes16()
	if err := q.done(); err != nil {
		return "", err
	}
	if len(name) == 0 {
		return "", alertError(alertDecodeError, &ProtocolError{"empty server name"})
	}
	return string(name), nil
}
func marshalALPN(protocols []string) ([]byte, error) {
	if len(protocols) == 0 {
		return nil, nil
	}
	length := 0
	for _, protocol := range protocols {
		if protocol == "" {
			return nil, &ProtocolError{"empty ALPN protocol"}
		}
		if len(protocol) > 255 {
			return nil, &ProtocolError{"8-bit vector overflow"}
		}
		length += 1 + len(protocol)
	}
	if length > 65535 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	out := newWireBuilder(2 + length)
	start := out.startVector16()
	for _, protocol := range protocols {
		out.string8(protocol)
	}
	out.endVector16(start)
	return out.b, out.err
}
func marshalCookie(cookie []byte) ([]byte, error) {
	if len(cookie) == 0 {
		return nil, &ProtocolError{"empty cookie extension"}
	}
	if len(cookie) > 65535 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	w := newWireBuilder(2 + len(cookie))
	w.bytes16(cookie)
	return w.b, w.err
}
func parseCookie(b []byte) ([]byte, error) {
	p := wireParser{b: b}
	cookie := p.bytes16()
	if err := p.done(); err != nil {
		return nil, err
	}
	if len(cookie) == 0 {
		return nil, alertError(alertDecodeError, &ProtocolError{"empty cookie extension"})
	}
	return append([]byte(nil), cookie...), nil
}
func marshalConnectionID(connectionID []byte) ([]byte, error) {
	if len(connectionID) > 255 {
		return nil, &ProtocolError{"8-bit vector overflow"}
	}
	w := newWireBuilder(1 + len(connectionID))
	w.bytes8(connectionID)
	return w.b, w.err
}
func parseConnectionID(b []byte) ([]byte, error) {
	p := wireParser{b: b}
	connectionID := append([]byte(nil), p.bytes8()...)
	if err := p.done(); err != nil {
		return nil, err
	}
	return connectionID, nil
}

var (
	clientSupportedVersionsStorage = [...]byte{2, byte(VersionDTLS13 >> 8), byte(VersionDTLS13 & 0xff)}
	serverSupportedVersionStorage  = [...]byte{byte(VersionDTLS13 >> 8), byte(VersionDTLS13 & 0xff)}
	pskKeyExchangeModesStorage     = [...]byte{1, 1}
)

func marshalPSKKeyExchangeModes() []byte { return pskKeyExchangeModesStorage[:] }

func parsePSKKeyExchangeModes(b []byte) (bool, error) {
	if len(b) < 2 || int(b[0]) != len(b)-1 {
		return false, alertError(alertDecodeError, &ProtocolError{"invalid psk_key_exchange_modes extension length"})
	}
	for _, mode := range b[1:] {
		if mode == 1 {
			return true, nil
		}
	}
	return false, nil
}

func marshalClientPSK(identity []byte, obfuscatedAge uint32, binder []byte) ([]byte, error) {
	return marshalClientPSKs([]pskIdentityEntry{{identity: identity, obfuscatedAge: obfuscatedAge}}, [][]byte{binder})
}

func marshalClientPSKs(identityList []pskIdentityEntry, binderList [][]byte) ([]byte, error) {
	if len(identityList) == 0 || len(identityList) != len(binderList) {
		return nil, &ProtocolError{"invalid pre_shared_key offer"}
	}
	for i, identity := range identityList {
		if len(identity.identity) == 0 || len(binderList[i]) < 32 {
			return nil, &ProtocolError{"invalid pre_shared_key offer"}
		}
	}
	identitiesLength := 0
	for _, identity := range identityList {
		if len(identity.identity) > 65535 || identitiesLength > 65535-6-len(identity.identity) {
			return nil, &ProtocolError{"16-bit vector overflow"}
		}
		identitiesLength += 6 + len(identity.identity)
	}
	bindersLength := 0
	for _, binder := range binderList {
		if len(binder) > 255 {
			return nil, &ProtocolError{"8-bit vector overflow"}
		}
		if bindersLength > 65535-1-len(binder) {
			return nil, &ProtocolError{"16-bit vector overflow"}
		}
		bindersLength += 1 + len(binder)
	}
	out := newWireBuilder(2 + identitiesLength + 2 + bindersLength)
	identitiesStart := out.startVector16()
	for _, identity := range identityList {
		out.bytes16(identity.identity)
		out.u32(identity.obfuscatedAge)
	}
	out.endVector16(identitiesStart)
	bindersStart := out.startVector16()
	for _, binder := range binderList {
		out.bytes8(binder)
	}
	out.endVector16(bindersStart)
	return out.b, out.err
}

func parseClientPSKs(b []byte) (identityList []pskIdentityEntry, binderList [][]byte, err error) {
	p := wireParser{b: b}
	identities := wireParser{b: p.bytes16()}
	binders := wireParser{b: p.bytes16()}
	if err = p.done(); err != nil {
		return nil, nil, err
	}
	for identities.off < len(identities.b) {
		identity := append([]byte(nil), identities.bytes16()...)
		age := identities.u32()
		if identities.err != nil || len(identity) == 0 {
			return nil, nil, alertError(alertDecodeError, &ProtocolError{"invalid pre_shared_key identity"})
		}
		identityList = append(identityList, pskIdentityEntry{identity: identity, obfuscatedAge: age})
	}
	for binders.off < len(binders.b) {
		binder := append([]byte(nil), binders.bytes8()...)
		if binders.err != nil || len(binder) < 32 {
			return nil, nil, alertError(alertDecodeError, &ProtocolError{"invalid pre_shared_key binder length"})
		}
		binderList = append(binderList, binder)
	}
	if len(identityList) == 0 || len(identityList) != len(binderList) {
		return nil, nil, &ProtocolError{"pre_shared_key identity and binder counts differ"}
	}
	return identityList, binderList, nil
}
func parseALPN(b []byte) ([]string, error) {
	p := wireParser{b: b}
	raw := p.bytes16()
	if err := p.done(); err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return nil, alertError(alertDecodeError, &ProtocolError{"empty ALPN protocol list"})
	}
	q := wireParser{b: raw}
	var out []string
	for q.off < len(q.b) {
		v := q.bytes8()
		if q.err != nil {
			return nil, q.err
		}
		if len(v) == 0 {
			return nil, alertError(alertDecodeError, &ProtocolError{"empty ALPN protocol"})
		}
		out = append(out, string(v))
	}
	return out, nil
}
func parseKeyShares(b []byte, vector bool) ([]keyShareEntry, error) {
	return parseKeySharesMode(b, vector, true, nil)
}

func parseKeySharesView(b []byte, vector bool) ([]keyShareEntry, error) {
	return parseKeySharesMode(b, vector, false, nil)
}

func parseKeySharesViewInto(b []byte, vector bool, dst []keyShareEntry) ([]keyShareEntry, error) {
	return parseKeySharesMode(b, vector, false, dst)
}

func parseKeySharesMode(b []byte, vector, copyData bool, dst []keyShareEntry) ([]keyShareEntry, error) {
	p := wireParser{b: b}
	if vector {
		b = p.bytes16()
		if err := p.done(); err != nil {
			return nil, err
		}
	}
	entries := b
	q := wireParser{b: entries}
	var groupStorage [8]tls.CurveID
	groups := groupStorage[:0]
	var seen map[tls.CurveID]struct{}
	count := 0
	for q.off < len(q.b) {
		group := tls.CurveID(q.u16())
		data := q.bytes16()
		if q.err != nil {
			return nil, q.err
		}
		if len(data) == 0 {
			return nil, alertError(alertDecodeError, &ProtocolError{"empty key share"})
		}
		if seen == nil {
			for _, candidate := range groups {
				if candidate == group {
					return nil, alertError(alertIllegalParameter, &ProtocolError{"duplicate key_share group"})
				}
			}
			if len(groups) < cap(groups) {
				groups = append(groups, group)
			} else {
				seen = make(map[tls.CurveID]struct{}, 16)
				for _, candidate := range groups {
					seen[candidate] = struct{}{}
				}
				seen[group] = struct{}{}
			}
		} else {
			if _, duplicate := seen[group]; duplicate {
				return nil, alertError(alertIllegalParameter, &ProtocolError{"duplicate key_share group"})
			}
			seen[group] = struct{}{}
		}
		count++
	}
	if count == 0 {
		if vector {
			return nil, nil
		}
		return nil, alertError(alertDecodeError, &ProtocolError{"empty key share list"})
	}
	var out []keyShareEntry
	if cap(dst) >= count {
		out = dst[:count]
	} else {
		out = make([]keyShareEntry, count)
	}
	q = wireParser{b: entries}
	for i := range out {
		out[i].group = tls.CurveID(q.u16())
		data := q.bytes16()
		if copyData {
			data = append([]byte(nil), data...)
		}
		out[i].data = data
	}
	return out, nil
}

func extensionsWireLength(items map[uint16][]byte, order []uint16) (int, error) {
	length := 0
	for _, typ := range order {
		if v, ok := items[typ]; ok {
			if len(v) > 65535 || length > 65535-4-len(v) {
				return 0, &ProtocolError{"16-bit vector overflow"}
			}
			length += 4 + len(v)
		}
	}
	return 2 + length, nil
}

func allExtensionsWireLength(items map[uint16][]byte) (int, error) {
	length := 0
	for _, value := range items {
		if len(value) > 65535 || length > 65535-4-len(value) {
			return 0, &ProtocolError{"16-bit vector overflow"}
		}
		length += 4 + len(value)
	}
	return 2 + length, nil
}

func appendExtensions(w *wireBuilder, items map[uint16][]byte, order []uint16) {
	start := w.startVector16()
	for _, typ := range order {
		v, ok := items[typ]
		if !ok {
			continue
		}
		w.u16(int(typ))
		w.bytes16(v)
	}
	w.endVector16(start)
}

type orderedExtension struct {
	typ   uint16
	value []byte
}

func orderedExtensionsWireLength(items []orderedExtension) (int, error) {
	length := 0
	for _, item := range items {
		if len(item.value) > 65535 || length > 65535-4-len(item.value) {
			return 0, &ProtocolError{"16-bit vector overflow"}
		}
		length += 4 + len(item.value)
	}
	return 2 + length, nil
}

func appendOrderedExtensions(w *wireBuilder, items []orderedExtension) {
	start := w.startVector16()
	for _, item := range items {
		w.u16(int(item.typ))
		w.bytes16(item.value)
	}
	w.endVector16(start)
}

func parseOrderedExtensionsView(b []byte, dst []orderedExtension) ([]orderedExtension, error) {
	p := wireParser{b: b}
	raw := p.bytes16()
	if err := p.done(); err != nil {
		return nil, err
	}
	q := wireParser{b: raw}
	items := dst[:0]
	for q.off < len(q.b) {
		typ := uint16(q.u16())
		value := q.bytes16()
		if q.err != nil {
			return nil, q.err
		}
		for _, item := range items {
			if item.typ == typ {
				return nil, alertError(alertIllegalParameter, &ProtocolError{"duplicate extension"})
			}
		}
		items = append(items, orderedExtension{typ: typ, value: value})
	}
	return items, q.done()
}

func orderedExtensionValue(items []orderedExtension, typ uint16) ([]byte, bool) {
	for _, item := range items {
		if item.typ == typ {
			return item.value, true
		}
	}
	return nil, false
}

func marshalExtensions(items map[uint16][]byte, order []uint16) ([]byte, error) {
	length, err := extensionsWireLength(items, order)
	if err != nil {
		return nil, err
	}
	out := newWireBuilder(length)
	appendExtensions(&out, items, order)
	return out.b, out.err
}
func parseExtensions(b []byte) (map[uint16][]byte, error) {
	return parseExtensionsMode(b, true)
}

func parseExtensionsView(b []byte) (map[uint16][]byte, error) {
	return parseExtensionsMode(b, false)
}

func parseExtensionsMode(b []byte, copyValues bool) (map[uint16][]byte, error) {
	var storage [8]orderedExtension
	items, err := parseOrderedExtensionsView(b, storage[:0])
	if err != nil {
		return nil, err
	}
	out := make(map[uint16][]byte, len(items))
	for _, item := range items {
		v := item.value
		if copyValues {
			v = append([]byte(nil), v...)
		}
		out[item.typ] = v
	}
	return out, nil
}

func (h *clientHello) marshal() ([]byte, error) {
	if len(h.cipherSuites) == 0 {
		return nil, &ProtocolError{"ClientHello has no cipher suites"}
	}
	if h.hasRecordSizeLimit && (h.recordSizeLimit < minRecordSizeLimit || h.recordSizeLimit > defaultRecordSizeLimit) {
		return nil, &ProtocolError{"invalid record_size_limit"}
	}
	ks, err := marshalKeyShares(h.keyShares, true)
	if err != nil {
		return nil, err
	}
	signatures, err := marshalSignatureSchemes(h.signatureSchemes)
	if err != nil {
		return nil, err
	}
	groups, err := marshalSupportedGroups(h.supportedGroups, h.keyShares)
	if err != nil {
		return nil, err
	}
	var certificateSignatures []byte
	if len(h.certificateSignatureSchemes) > 0 {
		certificateSignatures, err = marshalSignatureSchemes(h.certificateSignatureSchemes)
		if err != nil {
			return nil, err
		}
	}
	var serverName []byte
	if h.serverName != "" {
		serverName, err = marshalServerName(h.serverName)
		if err != nil {
			return nil, err
		}
	}
	var alpn []byte
	if len(h.alpn) > 0 {
		alpn, err = marshalALPN(h.alpn)
		if err != nil {
			return nil, err
		}
	}
	var cookie []byte
	if len(h.cookie) > 0 {
		cookie, err = marshalCookie(h.cookie)
		if err != nil {
			return nil, err
		}
	}
	var connectionID []byte
	if h.hasConnectionID {
		connectionID, err = marshalConnectionID(h.connectionID)
		if err != nil {
			return nil, err
		}
	}
	var psk []byte
	hasPSK := false
	if len(h.pskIdentities) > 0 {
		psk, err = marshalClientPSKs(h.pskIdentities, h.pskBinders)
		if err != nil {
			return nil, err
		}
		hasPSK = true
	} else if len(h.pskIdentity) > 0 {
		psk, err = marshalClientPSK(h.pskIdentity, h.obfuscatedAge, h.pskBinder)
		if err != nil {
			return nil, err
		}
		hasPSK = true
	} else if len(h.pskBinder) > 0 || len(h.pskBinders) > 0 {
		return nil, &ProtocolError{"PSK binder without identity"}
	}
	var recordSizeLimit [2]byte
	if h.hasRecordSizeLimit {
		binary.BigEndian.PutUint16(recordSizeLimit[:], h.recordSizeLimit)
	}
	var certificateCompression []byte
	if h.certificateCompressionOffered {
		certificateCompression, err = marshalCertificateCompressionAlgorithms(h.certificateCompressionAlgorithms())
		if err != nil {
			return nil, err
		}
	}
	var extensionStorage [18]orderedExtension
	extensions := extensionStorage[:0]
	if serverName != nil {
		extensions = append(extensions, orderedExtension{typ: extServerName, value: serverName})
	}
	extensions = append(extensions,
		orderedExtension{typ: extSupportedGroups, value: groups},
		orderedExtension{typ: extSignatureAlgorithms, value: signatures},
	)
	if certificateSignatures != nil {
		extensions = append(extensions, orderedExtension{typ: extSignatureAlgorithmsCert, value: certificateSignatures})
	}
	if alpn != nil {
		extensions = append(extensions, orderedExtension{typ: extALPN, value: alpn})
	}
	if certificateCompression != nil {
		extensions = append(extensions, orderedExtension{typ: extCompressCertificate, value: certificateCompression})
	}
	if h.hasRecordSizeLimit {
		extensions = append(extensions, orderedExtension{typ: extRecordSizeLimit, value: recordSizeLimit[:]})
	}
	extensions = append(extensions, orderedExtension{typ: extSupportedVersions, value: clientSupportedVersionsStorage[:]})
	if cookie != nil {
		extensions = append(extensions, orderedExtension{typ: extCookie, value: cookie})
	}
	extensions = append(extensions, orderedExtension{typ: extKeyShare, value: ks})
	if h.postHandshakeAuth {
		extensions = append(extensions, orderedExtension{typ: extPostHandshakeAuth})
	}
	if h.hasConnectionID {
		extensions = append(extensions, orderedExtension{typ: extConnectionID, value: connectionID})
	}
	if h.returnRoutability {
		extensions = append(extensions, orderedExtension{typ: extReturnRoutability})
	}
	if h.earlyData {
		extensions = append(extensions, orderedExtension{typ: extEarlyData})
	}
	if ech := h.encryptedClientHello(); ech != nil {
		extensions = append(extensions, orderedExtension{typ: extECH, value: ech})
	}
	if hasPSK {
		extensions = append(extensions,
			orderedExtension{typ: extPSKKeyExchangeModes, value: marshalPSKKeyExchangeModes()},
			orderedExtension{typ: extPreSharedKey, value: psk},
		)
	}
	extsLength, err := orderedExtensionsWireLength(extensions)
	if err != nil {
		return nil, err
	}
	if len(h.sessionID) > 255 {
		return nil, &ProtocolError{"8-bit vector overflow"}
	}
	if len(h.cipherSuites) > 65535/2 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	suitesLength := 2 * len(h.cipherSuites)
	w := newWireBuilder(2 + len(h.random) + 1 + len(h.sessionID) + 1 + 2 + suitesLength + 2 + extsLength)
	w.u16(int(dtlsLegacyVersion))
	w.b = append(w.b, h.random[:]...)
	w.bytes8(h.sessionID)
	w.bytes8(nil)
	suitesStart := w.startVector16()
	for _, suite := range h.cipherSuites {
		w.u16(int(suite))
	}
	w.endVector16(suitesStart)
	w.bytes8([]byte{0})
	appendOrderedExtensions(&w, extensions)
	return w.b, w.err
}

func parseClientHello(b []byte) (*clientHello, error) {
	p := wireParser{b: b}
	if p.u16() != int(dtlsLegacyVersion) {
		return nil, &ProtocolError{"invalid ClientHello legacy version"}
	}
	h := &clientHello{}
	copy(h.random[:], p.take(32))
	h.sessionID = append([]byte(nil), p.bytes8()...)
	legacyCookie := p.bytes8()
	suites := p.bytes16()
	compression := p.bytes8()
	extBytes := p.take(len(p.b) - p.off)
	if err := p.done(); err != nil {
		return nil, err
	}
	if len(legacyCookie) != 0 {
		return nil, alertError(alertIllegalParameter, &ProtocolError{"DTLS 1.3 ClientHello legacy_cookie must be empty"})
	}
	if len(h.sessionID) > 32 || len(suites) < 2 || len(suites)%2 != 0 || len(compression) != 1 || compression[0] != 0 {
		return nil, alertError(alertDecodeError, &ProtocolError{"invalid ClientHello legacy fields"})
	}
	for len(suites) > 0 {
		h.cipherSuites = append(h.cipherSuites, binary.BigEndian.Uint16(suites))
		suites = suites[2:]
	}
	var extensionStorage [16]orderedExtension
	exts, err := parseOrderedExtensionsView(extBytes, extensionStorage[:0])
	if err != nil {
		return nil, err
	}
	for _, extension := range exts {
		switch extension.typ {
		case extServerName, extSupportedGroups, extSignatureAlgorithms, extSignatureAlgorithmsCert, extALPN, extPadding, extCompressCertificate, extRecordSizeLimit,
			extSupportedVersions, extCookie, extKeyShare, extPostHandshakeAuth,
			extConnectionID, extReturnRoutability, extEarlyData, extPSKKeyExchangeModes, extPreSharedKey, extECH:
		default:
			if h.unknownExtensions == nil {
				h.unknownExtensions = make(map[uint16][]byte)
			}
			h.unknownExtensions[extension.typ] = append([]byte(nil), extension.value...)
		}
	}
	if _, ok := orderedExtensionValue(exts, extECHOuterExtensions); ok {
		return nil, alertError(alertIllegalParameter, &ProtocolError{"ech_outer_extensions is only valid in EncodedClientHelloInner"})
	}
	if raw, ok := orderedExtensionValue(exts, extECH); ok {
		if _, _, _, _, _, parseErr := parseECHExt(raw); parseErr != nil {
			description := uint8(alertDecodeError)
			if errors.Is(parseErr, errInvalidECHExt) {
				description = alertIllegalParameter
			}
			return nil, alertError(description, parseErr)
		}
		h.setEncryptedClientHello(append([]byte(nil), raw...))
	}
	versions, ok := orderedExtensionValue(exts, extSupportedVersions)
	if !ok {
		return nil, alertError(alertProtocolVersion, &ProtocolError{"ClientHello has no supported_versions extension"})
	}
	offered, versionsErr := parseClientSupportedVersions(versions)
	if versionsErr != nil {
		return nil, versionsErr
	}
	if !offered {
		return nil, alertError(alertProtocolVersion, &ProtocolError{"ClientHello does not offer DTLS 1.3"})
	}
	if signatures, ok := orderedExtensionValue(exts, extSignatureAlgorithms); ok {
		h.signatureSchemes, err = parseSignatureSchemes(signatures)
		if err != nil {
			return nil, err
		}
	}
	if signatures, ok := orderedExtensionValue(exts, extSignatureAlgorithmsCert); ok {
		h.certificateSignatureSchemes, err = parseSignatureSchemes(signatures)
		if err != nil {
			return nil, err
		}
	}
	groupsRaw, hasGroups := orderedExtensionValue(exts, extSupportedGroups)
	ks, hasKeyShare := orderedExtensionValue(exts, extKeyShare)
	if hasGroups != hasKeyShare {
		return nil, alertError(alertMissingExtension, &ProtocolError{"supported_groups and key_share must be offered together"})
	}
	if hasGroups {
		h.supportedGroups, err = parseSupportedGroups(groupsRaw)
		if err != nil {
			return nil, err
		}
		h.keyShares, err = parseKeySharesViewInto(ks, true, h.keyShareStorage[:0])
		if err != nil {
			return nil, err
		}
	}
	nextGroup := 0
	for _, share := range h.keyShares {
		matched := false
		for nextGroup < len(h.supportedGroups) {
			group := h.supportedGroups[nextGroup]
			nextGroup++
			if group == share.group {
				matched = true
				break
			}
		}
		if !matched {
			return nil, alertError(alertIllegalParameter, &ProtocolError{"key_share groups are not an ordered subset of supported_groups"})
		}
	}
	if raw, ok := orderedExtensionValue(exts, extServerName); ok {
		h.serverName, err = parseServerName(raw)
		if err != nil {
			return nil, err
		}
	}
	if raw, ok := orderedExtensionValue(exts, extALPN); ok {
		h.alpn, err = parseALPN(raw)
		if err != nil {
			return nil, err
		}
	}
	if raw, ok := orderedExtensionValue(exts, extCompressCertificate); ok {
		h.certificateCompressionOffered = true
		if len(raw) != 3 || raw[0] != 2 || raw[1] != 0 || raw[2] != byte(certificateCompressionZlib) {
			if h.unknownExtensions == nil {
				h.unknownExtensions = make(map[uint16][]byte)
			}
			h.unknownExtensions[extCompressCertificate] = append([]byte(nil), raw...)
		}
		_, err = parseCertificateCompressionAlgorithms(raw)
		if err != nil {
			return nil, err
		}
	}
	if raw, ok := orderedExtensionValue(exts, extRecordSizeLimit); ok {
		h.recordSizeLimit, err = parseRecordSizeLimit(raw, false)
		if err != nil {
			return nil, err
		}
		h.hasRecordSizeLimit = true
	}
	if raw, ok := orderedExtensionValue(exts, extCookie); ok {
		h.cookie, err = parseCookie(raw)
		if err != nil {
			return nil, err
		}
	}
	if raw, ok := orderedExtensionValue(exts, extConnectionID); ok {
		h.connectionID, err = parseConnectionID(raw)
		if err != nil {
			return nil, err
		}
		h.hasConnectionID = true
	}
	if raw, ok := orderedExtensionValue(exts, extReturnRoutability); ok {
		if len(raw) != 0 {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid rrc extension length"})
		}
		h.returnRoutability = true
	}
	if raw, ok := orderedExtensionValue(exts, extPostHandshakeAuth); ok {
		if len(raw) != 0 {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid post_handshake_auth extension length"})
		}
		h.postHandshakeAuth = true
	}
	if raw, ok := orderedExtensionValue(exts, extEarlyData); ok {
		if len(raw) != 0 {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid early_data extension length"})
		}
		h.earlyData = true
	}
	if raw, ok := orderedExtensionValue(exts, extPadding); ok {
		for _, value := range raw {
			if value != 0 {
				return nil, alertError(alertIllegalParameter, &ProtocolError{"ClientHello padding contains non-zero data"})
			}
		}
	}
	if modes, ok := orderedExtensionValue(exts, extPSKKeyExchangeModes); ok {
		h.pskDHE, err = parsePSKKeyExchangeModes(modes)
		if err != nil {
			return nil, err
		}
	}
	if raw, ok := orderedExtensionValue(exts, extPreSharedKey); ok {
		if len(exts) == 0 || exts[len(exts)-1].typ != extPreSharedKey {
			return nil, &ProtocolError{"pre_shared_key must be the last ClientHello extension"}
		}
		_, modesOK := orderedExtensionValue(exts, extPSKKeyExchangeModes)
		if !modesOK {
			return nil, alertError(alertMissingExtension, &ProtocolError{"pre_shared_key without psk_key_exchange_modes"})
		}
		h.pskIdentities, h.pskBinders, err = parseClientPSKs(raw)
		if err != nil {
			return nil, err
		}
		h.pskIdentity = append([]byte(nil), h.pskIdentities[0].identity...)
		h.obfuscatedAge = h.pskIdentities[0].obfuscatedAge
		h.pskBinder = append([]byte(nil), h.pskBinders[0]...)
	}
	if h.earlyData && len(h.pskIdentity) == 0 {
		return nil, &ProtocolError{"early_data without pre_shared_key"}
	}
	if len(h.pskIdentity) == 0 && len(h.signatureSchemes) == 0 {
		return nil, alertError(alertMissingExtension, &ProtocolError{"certificate ClientHello has no signature_algorithms"})
	}
	if len(h.pskIdentity) == 0 && !hasGroups {
		return nil, alertError(alertMissingExtension, &ProtocolError{"certificate ClientHello has no supported_groups or key_share"})
	}
	return h, nil
}

func parseRecordSizeLimit(raw []byte, rejectAboveMaximum bool) (uint16, error) {
	if len(raw) != 2 {
		return 0, alertError(alertDecodeError, &ProtocolError{"invalid record_size_limit extension length"})
	}
	limit := binary.BigEndian.Uint16(raw)
	if limit < minRecordSizeLimit || (rejectAboveMaximum && limit > defaultRecordSizeLimit) {
		return 0, alertError(alertIllegalParameter, &ProtocolError{"invalid record_size_limit value"})
	}
	return limit, nil
}

func (h *serverHello) marshal() ([]byte, error) {
	if len(h.sessionID) != 0 {
		return nil, alertError(alertIllegalParameter, &ProtocolError{"DTLS 1.3 ServerHello legacy_session_id must be empty"})
	}
	if len(h.keyShare.data) > 65535 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	if h.hasConnectionID && len(h.connectionID) > 255 {
		return nil, &ProtocolError{"8-bit vector overflow"}
	}
	// Compute the extension vector length without allocating temporary values.
	extsLength := 2
	addExtension := func(valueLength int) bool {
		if valueLength > 65535 || extsLength > 65535-4-valueLength {
			return false
		}
		extsLength += 4 + valueLength
		return true
	}
	if !addExtension(2) || !addExtension(4+len(h.keyShare.data)) {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	if h.hasConnectionID && !addExtension(1+len(h.connectionID)) {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	if h.returnRoutability && !addExtension(0) {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	if h.selectedIdentity != nil && !addExtension(2) {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	w := newWireBuilder(2 + len(h.random) + 1 + len(h.sessionID) + 2 + 1 + extsLength)
	w.u16(int(dtlsLegacyVersion))
	w.b = append(w.b, h.random[:]...)
	w.bytes8(h.sessionID)
	w.u16(int(h.cipherSuite))
	w.u8(0)
	start := w.startVector16()
	w.u16(int(extSupportedVersions))
	w.bytes16(serverSupportedVersionStorage[:])
	w.u16(int(extKeyShare))
	w.u16(4 + len(h.keyShare.data))
	w.u16(int(h.keyShare.group))
	w.bytes16(h.keyShare.data)
	if h.hasConnectionID {
		w.u16(int(extConnectionID))
		w.u16(1 + len(h.connectionID))
		w.bytes8(h.connectionID)
	}
	if h.returnRoutability {
		w.u16(int(extReturnRoutability))
		w.u16(0)
	}
	if h.selectedIdentity != nil {
		w.u16(int(extPreSharedKey))
		w.u16(2)
		w.u16(int(*h.selectedIdentity))
	}
	w.endVector16(start)
	return w.b, w.err
}
func parseServerHello(b []byte) (*serverHello, error) {
	p := wireParser{b: b}
	if p.u16() != int(dtlsLegacyVersion) {
		return nil, &ProtocolError{"invalid ServerHello legacy version"}
	}
	h := &serverHello{}
	copy(h.random[:], p.take(32))
	h.sessionID = append([]byte(nil), p.bytes8()...)
	h.cipherSuite = uint16(p.u16())
	if p.u8() != 0 {
		return nil, &ProtocolError{"invalid ServerHello compression method"}
	}
	extBytes := p.take(len(p.b) - p.off)
	if err := p.done(); err != nil {
		return nil, err
	}
	if len(h.sessionID) != 0 {
		return nil, &ProtocolError{"DTLS 1.3 ServerHello legacy_session_id must be empty"}
	}
	var extensionStorage [5]orderedExtension
	exts, err := parseOrderedExtensionsView(extBytes, extensionStorage[:0])
	if err != nil {
		return nil, err
	}
	for _, extension := range exts {
		switch extension.typ {
		case extSupportedVersions, extKeyShare, extConnectionID, extReturnRoutability, extPreSharedKey:
		default:
			if knownExtensionType(extension.typ) {
				return nil, alertError(alertIllegalParameter, &ProtocolError{"recognized extension is not permitted in ServerHello"})
			}
			return nil, alertError(alertUnsupportedExtension, &ProtocolError{"unsupported ServerHello extension"})
		}
	}
	v, ok := orderedExtensionValue(exts, extSupportedVersions)
	if !ok {
		return nil, alertError(alertProtocolVersion, &ProtocolError{"ServerHello has no supported_versions extension"})
	}
	if len(v) != 2 {
		return nil, alertError(alertDecodeError, &ProtocolError{"invalid ServerHello supported_versions length"})
	}
	if binary.BigEndian.Uint16(v) != VersionDTLS13 {
		return nil, alertError(alertIllegalParameter, &ProtocolError{"server did not select DTLS 1.3"})
	}
	raw, ok := orderedExtensionValue(exts, extKeyShare)
	if !ok {
		return nil, alertError(alertMissingExtension, &ProtocolError{"ServerHello has no key_share"})
	}
	var keyShareStorage [1]keyShareEntry
	shares, err := parseKeySharesViewInto(raw, false, keyShareStorage[:0])
	if err != nil {
		return nil, err
	}
	if len(shares) != 1 {
		return nil, &ProtocolError{"ServerHello must contain one key share"}
	}
	h.keyShare = shares[0]
	if raw, ok := orderedExtensionValue(exts, extConnectionID); ok {
		h.connectionID, err = parseConnectionID(raw)
		if err != nil {
			return nil, err
		}
		h.hasConnectionID = true
	}
	if raw, ok := orderedExtensionValue(exts, extReturnRoutability); ok {
		if len(raw) != 0 {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid ServerHello rrc extension length"})
		}
		h.returnRoutability = true
	}
	if raw, ok := orderedExtensionValue(exts, extPreSharedKey); ok {
		if len(raw) != 2 {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid ServerHello pre_shared_key extension length"})
		}
		selected := binary.BigEndian.Uint16(raw)
		h.selectedIdentity = &selected
	}
	return h, nil
}
