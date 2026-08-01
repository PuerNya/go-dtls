package dtls13

import (
	"bytes"
	"crypto/ecdh"
	"crypto/hpke"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	extECH                   uint16 = 0xfe0d
	extECHOuterExtensions    uint16 = 0xfd00
	echOuterType             uint8  = 0
	echInnerType             uint8  = 1
	echAcceptConfirmationLen        = 8
	echAEADExpansion                = 16
)

// EncryptedClientHelloKey associates a serialized ECHConfig with its HPKE
// private key. Config is one ECHConfig, not an ECHConfigList, and must match
// PrivateKey byte-for-byte at the public-key level.
type EncryptedClientHelloKey struct {
	// Config is exactly one serialized ECHConfig, including its version and
	// length fields. It must not contain the outer ECHConfigList length.
	Config []byte
	// PrivateKey is the HPKE KEM private key corresponding to the public key in
	// Config, encoded using crypto/hpke's KEM private-key representation.
	PrivateKey []byte
	// SendAsRetry includes Config in retry_configs when ECH is rejected. At
	// least one configured key must be marked for retry.
	SendAsRetry bool
}

// ECHRejectionError reports an authenticated ECH rejection. RetryConfigList
// contains the server's syntactically validated ECHConfigList, when supplied.
// It is scoped to the configuration source and transport endpoint used for the
// failed connection; callers must not persist or share it outside that scope.
type ECHRejectionError struct {
	// RetryConfigList is the authenticated ECHConfigList supplied by the
	// server, including its two-byte length. It may be nil. The list is valid
	// only for the configuration source and endpoint used by this connection.
	RetryConfigList []byte
}

func (e *ECHRejectionError) Error() string { return "dtls13: server rejected ECH" }

type echCipher struct {
	kdfID  uint16
	aeadID uint16
}

type echConfigExtension struct {
	typ  uint16
	data []byte
}

type echConfig struct {
	raw          []byte
	configID     uint8
	kemID        uint16
	publicKey    []byte
	cipherSuites []echCipher
	maxNameLen   uint8
	publicName   string
	extensions   []echConfigExtension
}

type echClientContext struct {
	config          *echConfig
	hpkeContext     *hpke.Sender
	encapsulatedKey []byte
	kdfID           uint16
	aeadID          uint16
	innerHello      *clientHello
	outerHello      *clientHello
	innerBody       []byte
	innerTranscript *transcriptHash
	acceptedAtHRR   bool
	rejected        bool
	retryConfigs    []byte
}

type echServerContext struct {
	hpkeContext *hpke.Recipient
	configID    uint8
	cipher      echCipher
	inner       bool
}

var (
	errMalformedECHConfigList = errors.New("dtls13: malformed ECHConfigList")
	errMalformedECHExt        = errors.New("dtls13: malformed encrypted_client_hello extension")
	errInvalidECHExt          = errors.New("dtls13: invalid encrypted_client_hello extension")
)

func parseECHConfig(raw []byte) (skip bool, config echConfig, err error) {
	if len(raw) < 4 {
		return false, echConfig{}, errMalformedECHConfigList
	}
	version := uint16(raw[0])<<8 | uint16(raw[1])
	length := int(raw[2])<<8 | int(raw[3])
	if length != len(raw)-4 {
		return false, echConfig{}, errMalformedECHConfigList
	}
	if version != extECH {
		return true, echConfig{}, nil
	}
	p := wireParser{b: raw[4:]}
	config.raw = raw
	config.configID = uint8(p.u8())
	config.kemID = uint16(p.u16())
	config.publicKey = p.bytes16()
	cipherBytes := p.bytes16()
	config.maxNameLen = uint8(p.u8())
	config.publicName = string(p.bytes8())
	extensionBytes := p.bytes16()
	if err := p.done(); err != nil || len(config.publicKey) == 0 || len(cipherBytes) < 4 || len(cipherBytes)%4 != 0 || config.publicName == "" {
		return false, echConfig{}, errMalformedECHConfigList
	}
	for len(cipherBytes) > 0 {
		config.cipherSuites = append(config.cipherSuites, echCipher{
			kdfID:  uint16(cipherBytes[0])<<8 | uint16(cipherBytes[1]),
			aeadID: uint16(cipherBytes[2])<<8 | uint16(cipherBytes[3]),
		})
		cipherBytes = cipherBytes[4:]
	}
	seen := make(map[uint16]struct{})
	extensions := wireParser{b: extensionBytes}
	for extensions.off < len(extensions.b) {
		typ := uint16(extensions.u16())
		data := extensions.bytes16()
		if extensions.err != nil {
			return false, echConfig{}, errMalformedECHConfigList
		}
		if _, duplicate := seen[typ]; duplicate {
			return false, echConfig{}, errMalformedECHConfigList
		}
		seen[typ] = struct{}{}
		config.extensions = append(config.extensions, echConfigExtension{typ: typ, data: data})
	}
	return false, config, nil
}

func parseECHConfigList(data []byte) ([]echConfig, error) {
	if len(data) < 6 || int(data[0])<<8|int(data[1]) != len(data)-2 {
		return nil, errMalformedECHConfigList
	}
	var configs []echConfig
	for offset := 2; offset < len(data); {
		if len(data)-offset < 4 {
			return nil, errMalformedECHConfigList
		}
		length := int(data[offset+2])<<8 | int(data[offset+3])
		end := offset + 4 + length
		if end > len(data) {
			return nil, errMalformedECHConfigList
		}
		skip, config, err := parseECHConfig(data[offset:end])
		if err != nil {
			return nil, err
		}
		if !skip {
			configs = append(configs, config)
		}
		offset = end
	}
	return configs, nil
}

func pickECHConfig(configs []echConfig) (*echConfig, hpke.PublicKey, hpke.KDF, hpke.AEAD) {
	for i := range configs {
		config := &configs[i]
		if !validECHPublicName(config.publicName) {
			continue
		}
		mandatory := false
		for _, extension := range config.extensions {
			mandatory = mandatory || extension.typ&0x8000 != 0
		}
		if mandatory {
			continue
		}
		kem, err := hpke.NewKEM(config.kemID)
		if err != nil {
			continue
		}
		publicKey, err := kem.NewPublicKey(config.publicKey)
		if err != nil {
			continue
		}
		for _, cipher := range config.cipherSuites {
			if cipher.aeadID < 1 || cipher.aeadID > 3 {
				continue
			}
			kdf, kdfErr := hpke.NewKDF(cipher.kdfID)
			aead, aeadErr := hpke.NewAEAD(cipher.aeadID)
			if kdfErr == nil && aeadErr == nil {
				return config, publicKey, kdf, aead
			}
		}
	}
	return nil, nil, nil, nil
}

func validECHPublicName(name string) bool {
	if len(name) == 0 || len(name) > 253 || strings.HasPrefix(name, ".") || strings.HasSuffix(name, ".") {
		return false
	}
	labels := strings.Split(name, ".")
	if len(labels) < 2 {
		return false
	}
	for _, label := range labels {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, c := range []byte(label) {
			if (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') && c != '-' {
				return false
			}
		}
	}
	last := labels[len(labels)-1]
	allDigits := true
	for _, c := range []byte(last) {
		allDigits = allDigits && c >= '0' && c <= '9'
	}
	if allDigits {
		return false
	}
	if len(last) >= 2 && last[0] == '0' && (last[1] == 'x' || last[1] == 'X') {
		hex := true
		for _, c := range []byte(last[2:]) {
			hex = hex && ((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F'))
		}
		if hex {
			return false
		}
	}
	return true
}

func validateECHKeys(keys []EncryptedClientHelloKey) error {
	retry := false
	for i := range keys {
		key := &keys[i]
		skip, config, err := parseECHConfig(key.Config)
		if err != nil || skip {
			return fmt.Errorf("EncryptedClientHelloKeys[%d] has an invalid Config", i)
		}
		kem, err := hpke.NewKEM(config.kemID)
		if err != nil {
			return fmt.Errorf("EncryptedClientHelloKeys[%d] uses an unsupported KEM", i)
		}
		privateKey, err := kem.NewPrivateKey(key.PrivateKey)
		if err != nil || !bytes.Equal(privateKey.PublicKey().Bytes(), config.publicKey) {
			return fmt.Errorf("EncryptedClientHelloKeys[%d] PrivateKey does not match Config", i)
		}
		if picked, _, _, _ := pickECHConfig([]echConfig{config}); picked == nil {
			return fmt.Errorf("EncryptedClientHelloKeys[%d] contains no supported cipher suite", i)
		}
		retry = retry || key.SendAsRetry
	}
	if len(keys) > 0 && !retry {
		return errors.New("EncryptedClientHelloKeys must include at least one SendAsRetry key")
	}
	return nil
}

func buildRetryConfigList(keys []EncryptedClientHelloKey) ([]byte, error) {
	length := 0
	for _, key := range keys {
		if key.SendAsRetry {
			if length > 65535-len(key.Config) {
				return nil, errors.New("ECH retry configuration list exceeds 65535 bytes")
			}
			length += len(key.Config)
		}
	}
	w := newWireBuilder(2 + length)
	start := w.startVector16()
	for _, key := range keys {
		if key.SendAsRetry {
			w.b = append(w.b, key.Config...)
		}
	}
	w.endVector16(start)
	return w.b, w.err
}

func parseECHExt(raw []byte) (typ uint8, cipher echCipher, configID uint8, enc, payload []byte, err error) {
	p := wireParser{b: raw}
	typ = uint8(p.u8())
	if typ == echInnerType {
		if doneErr := p.done(); doneErr != nil {
			return 0, echCipher{}, 0, nil, nil, errMalformedECHExt
		}
		return typ, echCipher{}, 0, nil, nil, nil
	}
	if typ != echOuterType {
		return 0, echCipher{}, 0, nil, nil, errInvalidECHExt
	}
	cipher.kdfID = uint16(p.u16())
	cipher.aeadID = uint16(p.u16())
	configID = uint8(p.u8())
	enc = p.bytes16()
	payload = p.bytes16()
	if doneErr := p.done(); doneErr != nil || len(payload) == 0 {
		return 0, echCipher{}, 0, nil, nil, errMalformedECHExt
	}
	return typ, cipher, configID, enc, payload, nil
}

type rawClientHelloExtension struct {
	typ         uint16
	value       []byte
	valueOffset int
}

func extractRawClientHelloExtensions(body []byte) ([]rawClientHelloExtension, error) {
	p := wireParser{b: body}
	p.take(2 + 32)
	p.bytes8()
	p.bytes8()
	p.bytes16()
	p.bytes8()
	vectorOffset := p.off
	extensionBytes := p.bytes16()
	if err := p.done(); err != nil {
		return nil, err
	}
	extensions := wireParser{b: extensionBytes}
	seen := make(map[uint16]struct{})
	var out []rawClientHelloExtension
	for extensions.off < len(extensions.b) {
		typ := uint16(extensions.u16())
		value := extensions.bytes16()
		if extensions.err != nil {
			return nil, extensions.err
		}
		if _, duplicate := seen[typ]; duplicate {
			return nil, &ProtocolError{"duplicate ClientHello extension"}
		}
		seen[typ] = struct{}{}
		out = append(out, rawClientHelloExtension{typ: typ, value: value, valueOffset: vectorOffset + 2 + extensions.off - len(value)})
	}
	return out, nil
}

func echOuterAAD(body []byte) ([]byte, error) {
	extensions, err := extractRawClientHelloExtensions(body)
	if err != nil {
		return nil, err
	}
	for _, extension := range extensions {
		if extension.typ != extECH {
			continue
		}
		typ, _, _, _, payload, parseErr := parseECHExt(extension.value)
		if parseErr != nil || typ != echOuterType {
			return nil, errInvalidECHExt
		}
		p := wireParser{b: extension.value}
		p.take(1 + 2 + 2 + 1)
		p.bytes16()
		p.bytes16()
		payloadOffset := extension.valueOffset + p.off - len(payload)
		aad := bytes.Clone(body)
		clear(aad[payloadOffset : payloadOffset+len(payload)])
		return aad, nil
	}
	return nil, errInvalidECHExt
}

func encodeInnerClientHello(inner *clientHello, maxNameLength int) ([]byte, error) {
	encoded := cloneClientHello(inner)
	encoded.sessionID = nil
	body, err := encoded.marshal()
	if err != nil {
		return nil, err
	}
	padding := maxNameLength + 9
	if inner.serverName != "" {
		padding = max(0, maxNameLength-len(inner.serverName))
	}
	padding += (32 - (len(body)+padding)%32) % 32
	if len(body)+padding+echAEADExpansion > 65535 {
		return nil, &ProtocolError{"EncodedClientHelloInner exceeds ECH payload limit"}
	}
	return append(body, make([]byte, padding)...), nil
}

func decodeInnerClientHello(outer *clientHello, outerBody, encoded []byte) (*clientHello, []byte, error) {
	p := wireParser{b: encoded}
	versionRandom := p.take(2 + 32)
	sessionID := p.bytes8()
	legacyCookie := p.bytes8()
	cipherSuites := p.bytes16()
	compression := p.bytes8()
	extensionsOffset := p.off
	p.bytes16()
	if p.err != nil || len(sessionID) != 0 {
		return nil, nil, errInvalidECHExt
	}
	for _, value := range encoded[p.off:] {
		if value != 0 {
			return nil, nil, errInvalidECHExt
		}
	}
	var storage [16]orderedExtension
	innerExtensions, err := parseOrderedExtensionsView(encoded[extensionsOffset:p.off], storage[:0])
	if err != nil {
		return nil, nil, errInvalidECHExt
	}
	outerExtensions, err := extractRawClientHelloExtensions(outerBody)
	if err != nil {
		return nil, nil, errInvalidECHExt
	}
	outerIndex := make(map[uint16]int, len(outerExtensions))
	for i, extension := range outerExtensions {
		outerIndex[extension.typ] = i
	}
	w := newWireBuilder(len(encoded) + len(outerBody))
	w.b = append(w.b, versionRandom...)
	w.bytes8(outer.sessionID)
	w.bytes8(legacyCookie)
	w.bytes16(cipherSuites)
	w.bytes8(compression)
	start := w.startVector16()
	for _, extension := range innerExtensions {
		if extension.typ != extECHOuterExtensions {
			w.u16(int(extension.typ))
			w.bytes16(extension.value)
			continue
		}
		refs := wireParser{b: extension.value}
		types := refs.bytes8()
		if err := refs.done(); err != nil || len(types) < 2 || len(types) > 254 || len(types)%2 != 0 {
			return nil, nil, errInvalidECHExt
		}
		last := -1
		for len(types) > 0 {
			typ := uint16(types[0])<<8 | uint16(types[1])
			types = types[2:]
			index, ok := outerIndex[typ]
			if !ok || typ == extECH || index <= last {
				return nil, nil, errInvalidECHExt
			}
			last = index
			w.u16(int(typ))
			w.bytes16(outerExtensions[index].value)
		}
	}
	w.endVector16(start)
	if w.err != nil {
		return nil, nil, errInvalidECHExt
	}
	inner, err := parseClientHello(w.b)
	if err != nil {
		return nil, nil, errInvalidECHExt
	}
	typ, _, _, _, _, err := parseECHExt(inner.encryptedClientHello())
	if err != nil || typ != echInnerType {
		return nil, nil, errInvalidECHExt
	}
	return inner, w.b, nil
}

func generateOuterECHExt(configID uint8, cipher echCipher, enc, payload []byte) ([]byte, error) {
	w := newWireBuilder(10 + len(enc) + len(payload))
	w.u8(int(echOuterType))
	w.u16(int(cipher.kdfID))
	w.u16(int(cipher.aeadID))
	w.u8(int(configID))
	w.bytes16(enc)
	w.bytes16(payload)
	return w.b, w.err
}

func computeOuterECH(outer, inner *clientHello, context *echClientContext, initial bool) ([]byte, error) {
	encoded, err := encodeInnerClientHello(inner, int(context.config.maxNameLen))
	if err != nil {
		return nil, err
	}
	enc := context.encapsulatedKey
	if !initial {
		enc = nil
	}
	cipher := echCipher{kdfID: context.kdfID, aeadID: context.aeadID}
	outerECH, err := generateOuterECHExt(context.config.configID, cipher, enc, make([]byte, len(encoded)+echAEADExpansion))
	if err != nil {
		return nil, err
	}
	outer.setEncryptedClientHello(outerECH)
	aad, err := outer.marshal()
	if err != nil {
		return nil, err
	}
	payload, err := context.hpkeContext.Seal(aad, encoded)
	if err != nil {
		return nil, err
	}
	outerECH, err = generateOuterECHExt(context.config.configID, cipher, enc, payload)
	if err != nil {
		return nil, err
	}
	outer.setEncryptedClientHello(outerECH)
	return outer.marshal()
}

func decryptECHPayload(context *hpke.Recipient, outerBody, payload []byte) ([]byte, error) {
	aad, err := echOuterAAD(outerBody)
	if err != nil {
		return nil, err
	}
	return context.Open(aad, payload)
}

func newECHClientContext(list []byte) (*echClientContext, error) {
	configs, err := parseECHConfigList(list)
	if err != nil {
		return nil, err
	}
	config, publicKey, kdf, aead := pickECHConfig(configs)
	if config == nil {
		return nil, errors.New("dtls13: ECHConfigList contains no supported configuration")
	}
	info := append([]byte("tls ech\x00"), config.raw...)
	enc, sender, err := hpke.NewSender(publicKey, kdf, aead, info)
	if err != nil {
		return nil, err
	}
	return &echClientContext{config: config, hpkeContext: sender, encapsulatedKey: enc, kdfID: kdf.ID(), aeadID: aead.ID()}, nil
}

func makeECHOuter(inner *clientHello, config *echConfig, random io.Reader) (*clientHello, error) {
	outer := cloneClientHello(inner)
	outer.serverName = config.publicName
	outer.alpn = nil
	_ = outer.setCertificateAuthorities(nil)
	if outer.grease {
		if outer.unknownExtensions == nil {
			outer.unknownExtensions = make(map[uint16][]byte, 1)
		}
		outer.unknownExtensions[greaseValue(outer.random[0])] = nil
		outer.grease = false
	}
	if _, err := io.ReadFull(random, outer.random[:]); err != nil {
		return nil, err
	}
	if err := greaseOuterPSK(outer, random); err != nil {
		return nil, err
	}
	return outer, nil
}

func greaseOuterPSK(outer *clientHello, random io.Reader) error {
	fillIdentity := func(identity []byte) (pskIdentityEntry, error) {
		grease := pskIdentityEntry{identity: make([]byte, len(identity))}
		var age [4]byte
		if _, err := io.ReadFull(random, grease.identity); err != nil {
			return pskIdentityEntry{}, err
		}
		if _, err := io.ReadFull(random, age[:]); err != nil {
			return pskIdentityEntry{}, err
		}
		grease.obfuscatedAge = uint32(age[0])<<24 | uint32(age[1])<<16 | uint32(age[2])<<8 | uint32(age[3])
		return grease, nil
	}
	if len(outer.pskIdentities) > 0 {
		identities := make([]pskIdentityEntry, len(outer.pskIdentities))
		binders := make([][]byte, len(outer.pskBinders))
		for i := range identities {
			var err error
			identities[i], err = fillIdentity(outer.pskIdentities[i].identity)
			if err != nil {
				return err
			}
			binders[i] = make([]byte, len(outer.pskBinders[i]))
			if _, err = io.ReadFull(random, binders[i]); err != nil {
				return err
			}
		}
		outer.pskIdentities, outer.pskBinders = identities, binders
		outer.pskIdentity, outer.pskBinder = nil, nil
		return nil
	}
	if len(outer.pskIdentity) > 0 {
		identity, err := fillIdentity(outer.pskIdentity)
		if err != nil {
			return err
		}
		outer.pskIdentity, outer.obfuscatedAge = identity.identity, identity.obfuscatedAge
		outer.pskBinder = make([]byte, len(outer.pskBinder))
		_, err = io.ReadFull(random, outer.pskBinder)
		return err
	}
	return nil
}

func generateGREASEECH(hello *clientHello, random io.Reader) ([]byte, error) {
	inner := cloneClientHello(hello)
	inner.setEncryptedClientHello([]byte{echInnerType})
	encoded, err := encodeInnerClientHello(inner, 0)
	if err != nil {
		return nil, err
	}
	privateKey, err := ecdh.X25519().GenerateKey(random)
	if err != nil {
		return nil, err
	}
	var configID [1]byte
	if _, err = io.ReadFull(random, configID[:]); err != nil {
		return nil, err
	}
	payload := make([]byte, len(encoded)+echAEADExpansion)
	if _, err = io.ReadFull(random, payload); err != nil {
		return nil, err
	}
	return generateOuterECHExt(configID[0], echCipher{kdfID: 1, aeadID: 1}, privateKey.PublicKey().Bytes(), payload)
}

func processECHClientHello(outer *clientHello, outerBody []byte, keys []EncryptedClientHelloKey) (*clientHello, []byte, *echServerContext, error) {
	typ, cipher, configID, enc, payload, err := parseECHExt(outer.encryptedClientHello())
	if err != nil {
		description := uint8(alertDecodeError)
		if errors.Is(err, errInvalidECHExt) {
			description = alertIllegalParameter
		}
		return nil, nil, nil, alertError(description, err)
	}
	if typ == echInnerType {
		return outer, outerBody, &echServerContext{inner: true}, nil
	}
	if len(keys) == 0 {
		return outer, outerBody, nil, nil
	}
	if err = validateECHKeys(keys); err != nil {
		return nil, nil, nil, &ConfigError{err.Error()}
	}
	for i := range keys {
		key := &keys[i]
		_, config, _ := parseECHConfig(key.Config)
		if config.configID != configID || !echConfigOffersCipher(&config, cipher) {
			continue
		}
		kem, _ := hpke.NewKEM(config.kemID)
		privateKey, keyErr := kem.NewPrivateKey(key.PrivateKey)
		kdf, kdfErr := hpke.NewKDF(cipher.kdfID)
		aead, aeadErr := hpke.NewAEAD(cipher.aeadID)
		if keyErr != nil || kdfErr != nil || aeadErr != nil || cipher.aeadID < 1 || cipher.aeadID > 3 {
			continue
		}
		info := append([]byte("tls ech\x00"), config.raw...)
		recipient, setupErr := hpke.NewRecipient(enc, privateKey, kdf, aead, info)
		if setupErr != nil {
			continue
		}
		encoded, openErr := decryptECHPayload(recipient, outerBody, payload)
		if openErr != nil {
			continue
		}
		inner, innerBody, decodeErr := decodeInnerClientHello(outer, outerBody, encoded)
		if decodeErr != nil {
			return nil, nil, nil, alertError(alertIllegalParameter, decodeErr)
		}
		return inner, innerBody, &echServerContext{hpkeContext: recipient, configID: configID, cipher: cipher}, nil
	}
	return outer, outerBody, nil, nil
}

func processSecondECHClientHello(outer *clientHello, outerBody []byte, context *echServerContext) (*clientHello, []byte, error) {
	if len(outer.encryptedClientHello()) == 0 {
		return nil, nil, alertError(alertMissingExtension, errors.New("dtls13: second ClientHello is missing encrypted_client_hello"))
	}
	typ, cipher, configID, enc, payload, err := parseECHExt(outer.encryptedClientHello())
	if err != nil {
		return nil, nil, alertError(alertDecodeError, err)
	}
	if context.inner {
		if typ != echInnerType {
			return nil, nil, alertError(alertIllegalParameter, errors.New("dtls13: second ClientHello changed ECH type"))
		}
		return outer, outerBody, nil
	}
	if typ != echOuterType || cipher != context.cipher || configID != context.configID || len(enc) != 0 {
		return nil, nil, alertError(alertIllegalParameter, errors.New("dtls13: second ClientHello changed ECH parameters"))
	}
	encoded, err := decryptECHPayload(context.hpkeContext, outerBody, payload)
	if err != nil {
		return nil, nil, alertError(alertDecryptError, errors.New("dtls13: failed to decrypt second ClientHelloInner"))
	}
	inner, innerBody, err := decodeInnerClientHello(outer, outerBody, encoded)
	if err != nil {
		return nil, nil, alertError(alertIllegalParameter, err)
	}
	return inner, innerBody, nil
}

func echConfigOffersCipher(config *echConfig, cipher echCipher) bool {
	for _, offered := range config.cipherSuites {
		if offered == cipher {
			return true
		}
	}
	return false
}

func echAcceptConfirmation(suite *cipherSuite, random [32]byte, label string, transcriptHash []byte) []byte {
	secret := make([]byte, suite.hash.Size())
	hkdfExtractInto(suite.hash.New, random[:], nil, secret)
	return expandLabel(suite, secret, label, transcriptHash, echAcceptConfirmationLen)
}

func cloneClientHello(source *clientHello) *clientHello {
	if source == nil {
		return nil
	}
	clone := *source
	clone.sessionID = bytes.Clone(source.sessionID)
	clone.encryptedClientHelloExtension = bytes.Clone(source.encryptedClientHelloExtension)
	clone.cookie = bytes.Clone(source.cookie)
	clone.cipherSuites = append([]uint16(nil), source.cipherSuites...)
	clone.keyShares = make([]keyShareEntry, len(source.keyShares))
	for i, share := range source.keyShares {
		clone.keyShares[i] = keyShareEntry{group: share.group, data: bytes.Clone(share.data)}
	}
	clone.keyShareStorage = [1]keyShareEntry{}
	clone.signatureSchemes = append([]tls.SignatureScheme(nil), source.signatureSchemes...)
	clone.certificateSignatureSchemes = append([]tls.SignatureScheme(nil), source.certificateSignatureSchemes...)
	clone.supportedGroups = append([]tls.CurveID(nil), source.supportedGroups...)
	clone.alpn = append([]string(nil), source.alpn...)
	clone.pskIdentity = bytes.Clone(source.pskIdentity)
	clone.pskBinder = bytes.Clone(source.pskBinder)
	clone.pskIdentities = make([]pskIdentityEntry, len(source.pskIdentities))
	for i, identity := range source.pskIdentities {
		clone.pskIdentities[i] = pskIdentityEntry{identity: bytes.Clone(identity.identity), obfuscatedAge: identity.obfuscatedAge}
	}
	clone.pskBinders = make([][]byte, len(source.pskBinders))
	for i, binder := range source.pskBinders {
		clone.pskBinders[i] = bytes.Clone(binder)
	}
	clone.connectionID = bytes.Clone(source.connectionID)
	clone.unknownExtensions = nil
	if source.unknownExtensions != nil {
		clone.unknownExtensions = make(map[uint16][]byte, len(source.unknownExtensions))
		for typ, value := range source.unknownExtensions {
			clone.unknownExtensions[typ] = bytes.Clone(value)
		}
	}
	return &clone
}
