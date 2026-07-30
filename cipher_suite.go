package dtls13

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hkdf"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"hash"
	"runtime"
	"unsafe"

	"github.com/pion/dtls/v3/pkg/crypto/ccm"
	"golang.org/x/crypto/chacha20"
	"golang.org/x/crypto/chacha20poly1305"
)

const (
	// TLS_AES_128_GCM_SHA256 identifies the TLS_AES_128_GCM_SHA256 TLS 1.3
	// cipher suite.
	TLS_AES_128_GCM_SHA256 uint16 = 0x1301
	// TLS_AES_256_GCM_SHA384 identifies the TLS_AES_256_GCM_SHA384 TLS 1.3
	// cipher suite.
	TLS_AES_256_GCM_SHA384 uint16 = 0x1302
	// TLS_CHACHA20_POLY1305_SHA256 identifies the
	// TLS_CHACHA20_POLY1305_SHA256 TLS 1.3 cipher suite.
	TLS_CHACHA20_POLY1305_SHA256 uint16 = 0x1303
	// TLS_AES_128_CCM_SHA256 identifies the TLS_AES_128_CCM_SHA256 TLS 1.3
	// cipher suite.
	TLS_AES_128_CCM_SHA256 uint16 = 0x1304
	// TLS_AES_128_CCM_8_SHA256 identifies TLS_AES_128_CCM_8_SHA256. The value
	// is exported for protocol identification, but the suite is not supported
	// by this package and is rejected in Config.CipherSuites. RFC 9147 requires
	// additional deployment-level safeguards for its shortened authentication
	// tag that a general-purpose library cannot enforce.
	TLS_AES_128_CCM_8_SHA256 uint16 = 0x1305
	maxSupportedHashSize            = sha512.Size384
)

type cipherSuite struct {
	id               uint16
	hash             crypto.Hash
	keyLen           int
	ivLen            int
	recordLimit      uint64
	authFailureLimit uint64
}

var supportedCipherSuites = [...]cipherSuite{
	{id: TLS_AES_128_GCM_SHA256, hash: crypto.SHA256, keyLen: 16, ivLen: 12, recordLimit: 1 << 24, authFailureLimit: uint64(1) << 36},
	{id: TLS_AES_256_GCM_SHA384, hash: crypto.SHA384, keyLen: 32, ivLen: 12, recordLimit: 1 << 24, authFailureLimit: uint64(1) << 36},
	{id: TLS_CHACHA20_POLY1305_SHA256, hash: crypto.SHA256, keyLen: chacha20poly1305.KeySize, ivLen: chacha20poly1305.NonceSize, recordLimit: 1 << 48, authFailureLimit: uint64(1) << 36},
	{id: TLS_AES_128_CCM_SHA256, hash: crypto.SHA256, keyLen: 16, ivLen: 12, recordLimit: 1 << 23, authFailureLimit: 1 << 23},
}

var (
	emptySHA256TranscriptHash = sha256.Sum256(nil)
	emptySHA384TranscriptHash = sha512.Sum384(nil)
	zeroHashSecretStorage     [maxSupportedHashSize]byte
	hkdfSingleByteFields      = func() (fields [256][1]byte) {
		for i := range fields {
			fields[i][0] = byte(i)
		}
		return fields
	}()
	hkdfSHA256LabelHeaders = func() (headers [256][3]byte) {
		for i := range headers {
			headers[i] = [3]byte{0, sha256.Size, byte(i)}
		}
		return headers
	}()
	hkdfSHA384LabelHeaders = func() (headers [256][3]byte) {
		for i := range headers {
			headers[i] = [3]byte{0, sha512.Size384, byte(i)}
		}
		return headers
	}()
	dtls13LabelPrefix       = []byte("dtls13")
	labelClientEarlyTraffic = newSingleBlockHKDFLabel("c e traffic")
	labelDerived            = newSingleBlockHKDFLabel("derived")
	labelFinished           = newSingleBlockHKDFLabel("finished")
	labelResumptionMaster   = newSingleBlockHKDFLabel("res master")
	labelResumptionBinder   = newSingleBlockHKDFLabel("res binder")
	labelTrafficUpdate      = newSingleBlockHKDFLabel("traffic upd")
	labelResumption         = newSingleBlockHKDFLabel("resumption")
)

// emptyTranscriptHash returns shared immutable storage. Callers must only pass
// it to APIs that read the digest.
func emptyTranscriptHash(suite *cipherSuite) []byte {
	switch suite.hash {
	case crypto.SHA256:
		return emptySHA256TranscriptHash[:]
	case crypto.SHA384:
		return emptySHA384TranscriptHash[:]
	default:
		panic("dtls13: unsupported transcript hash")
	}
}

// zeroHashSecret returns shared immutable Hash.length zero bytes for the TLS
// key schedule.
func zeroHashSecret(suite *cipherSuite) []byte {
	size := suite.hash.Size()
	if size <= 0 || size > len(zeroHashSecretStorage) {
		panic("dtls13: unsupported zero-secret length")
	}
	return zeroHashSecretStorage[:size]
}

func defaultCipherSuites() []uint16 {
	ids := make([]uint16, len(supportedCipherSuites))
	for i := range supportedCipherSuites {
		ids[i] = supportedCipherSuites[i].id
	}
	return ids
}

func cipherSuiteForID(id uint16) (*cipherSuite, error) {
	switch id {
	case TLS_AES_128_GCM_SHA256:
		return &supportedCipherSuites[0], nil
	case TLS_AES_256_GCM_SHA384:
		return &supportedCipherSuites[1], nil
	case TLS_CHACHA20_POLY1305_SHA256:
		return &supportedCipherSuites[2], nil
	case TLS_AES_128_CCM_SHA256:
		return &supportedCipherSuites[3], nil
	case TLS_AES_128_CCM_8_SHA256:
		return nil, errors.New("dtls13: TLS_AES_128_CCM_8_SHA256 requires deployment-specific forgery safeguards")
	default:
		return nil, errors.New("dtls13: unsupported cipher suite")
	}
}

func (s *cipherSuite) newAEAD(key []byte) (cipher.AEAD, error) {
	if s.id == TLS_CHACHA20_POLY1305_SHA256 {
		return chacha20poly1305.New(key)
	}
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if s.id == TLS_AES_128_CCM_SHA256 {
		return ccm.NewCCM(b, 16, s.ivLen)
	}
	return cipher.NewGCM(b)
}

type sequenceMasker interface {
	mask(sample []byte) ([2]byte, error)
}

type aesSequenceMask struct {
	block   cipher.Block
	scratch [aes.BlockSize]byte
}

func (m *aesSequenceMask) mask(sample []byte) ([2]byte, error) {
	if len(sample) < aes.BlockSize {
		return [2]byte{}, errors.New("dtls13: invalid AES sequence mask sample")
	}
	m.block.Encrypt(m.scratch[:], sample[:aes.BlockSize])
	return [2]byte{m.scratch[0], m.scratch[1]}, nil
}

type chachaSequenceMask struct{ key []byte }

func (m chachaSequenceMask) mask(sample []byte) ([2]byte, error) {
	if len(sample) < 16 {
		return [2]byte{}, errors.New("dtls13: invalid ChaCha20 sequence mask sample")
	}
	cipher, err := chacha20.NewUnauthenticatedCipher(m.key, sample[4:16])
	if err != nil {
		return [2]byte{}, err
	}
	cipher.SetCounter(binary.LittleEndian.Uint32(sample[:4]))
	var mask [2]byte
	cipher.XORKeyStream(mask[:], mask[:])
	return mask, nil
}

func (s *cipherSuite) newSequenceMask(key []byte) (sequenceMasker, error) {
	if s.id == TLS_CHACHA20_POLY1305_SHA256 {
		if len(key) != chacha20.KeySize {
			return nil, errors.New("dtls13: invalid ChaCha20 sequence mask key")
		}
		return chachaSequenceMask{key: append([]byte(nil), key...)}, nil
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &aesSequenceMask{block: block}, nil
}

func hkdfExtract(newHash func() hash.Hash, secret, salt []byte) []byte {
	out := make([]byte, newHash().Size())
	hkdfExtractInto(newHash, secret, salt, out)
	return out
}

func hkdfExtractInto(newHash func() hash.Hash, secret, salt, out []byte) {
	mac := hmac.New(newHash, salt)
	if len(out) != mac.Size() {
		panic("dtls13: HKDF-Extract output must equal Hash.length")
	}
	_, _ = mac.Write(secret)
	if sum := mac.Sum(out[:0]); len(sum) != len(out) {
		panic("dtls13: unexpected HKDF-Extract output length")
	}
}
func hkdfExpand(newHash func() hash.Hash, secret, info []byte, length int) []byte {
	out, err := hkdf.Expand(newHash, secret, string(info), length)
	if err != nil {
		panic("dtls13: HKDF-Expand failed: " + err.Error())
	}
	return out
}

type hkdfExpander struct {
	mac      hash.Hash
	hashSize int
	previous [64]byte
	counter  [1]byte
	info     [128]byte
}

func newHKDFExpander(newHash func() hash.Hash, secret []byte) hkdfExpander {
	mac := hmac.New(newHash, secret)
	return hkdfExpander{mac: mac, hashSize: mac.Size()}
}

func (e *hkdfExpander) expand(info []byte, length int) []byte {
	if length < 0 || length > 255*e.hashSize {
		panic("dtls13: HKDF-Expand length exceeds RFC 5869 limit")
	}
	out := make([]byte, length)
	e.expandInto(info, out)
	return out
}

func (e *hkdfExpander) expandInto(info, out []byte) {
	length := len(out)
	if length > 255*e.hashSize {
		panic("dtls13: HKDF-Expand length exceeds RFC 5869 limit")
	}
	previous := e.previous[:0]
	for at, counter := 0, byte(1); at < length; counter++ {
		e.mac.Reset()
		_, _ = e.mac.Write(previous)
		_, _ = e.mac.Write(info)
		e.counter[0] = counter
		_, _ = e.mac.Write(e.counter[:])
		previous = e.mac.Sum(e.previous[:0])
		at += copy(out[at:], previous)
	}
	clear(e.previous[:e.hashSize])
}

// expandLabel is HKDF-Expand-Label from RFC 8446 with the mandatory
// "dtls13" prefix (without a trailing space) from RFC 9147 section 5.9.
func expandLabel(s *cipherSuite, secret []byte, label string, context []byte, length int) []byte {
	fullLabel := "dtls13" + label
	info := make([]byte, 2+1+len(fullLabel)+1+len(context))
	info[0] = byte(length >> 8)
	info[1] = byte(length)
	info[2] = byte(len(fullLabel))
	copy(info[3:], fullLabel)
	at := 3 + len(fullLabel)
	info[at] = byte(len(context))
	copy(info[at+1:], context)
	return hkdfExpand(s.hash.New, secret, info, length)
}

func expandLabelWithInto(expander *hkdfExpander, label string, context, out []byte) {
	length := len(out)
	const prefix = "dtls13"
	fullLabelLen := len(prefix) + len(label)
	infoLen := 2 + 1 + fullLabelLen + 1 + len(context)
	var info []byte
	if infoLen <= len(expander.info) {
		info = expander.info[:infoLen]
	} else {
		info = make([]byte, infoLen)
	}
	info[0] = byte(length >> 8)
	info[1] = byte(length)
	info[2] = byte(fullLabelLen)
	copy(info[3:], prefix)
	copy(info[3+len(prefix):], label)
	at := 3 + fullLabelLen
	info[at] = byte(len(context))
	copy(info[at+1:], context)
	expander.expandInto(info, out)
}

func expandLabelWithScratchInto(s *cipherSuite, secret []byte, label string, context, out, scratch []byte) {
	const prefix = "dtls13"
	fullLabelLen := len(prefix) + len(label)
	infoLen := 2 + 1 + fullLabelLen + 1 + len(context)
	if fullLabelLen > 255 || len(context) > 255 || len(out) > 255*s.hash.Size() || len(scratch) < infoLen || len(scratch) < s.hash.Size() {
		panic("dtls13: invalid caller-storage HKDF-Expand-Label input")
	}
	mac := hmac.New(s.hash.New, secret)
	info := scratch[:infoLen]
	info[0] = byte(len(out) >> 8)
	info[1] = byte(len(out))
	info[2] = byte(fullLabelLen)
	copy(info[3:], prefix)
	copy(info[3+len(prefix):], label)
	at := 3 + fullLabelLen
	info[at] = byte(len(context))
	copy(info[at+1:], context)

	var previous []byte
	for at, counter := 0, 1; at < len(out); counter++ {
		mac.Reset()
		_, _ = mac.Write(previous)
		_, _ = mac.Write(info)
		_, _ = mac.Write(hkdfSingleByteFields[counter][:])
		remaining := len(out) - at
		if remaining < s.hash.Size() {
			block := mac.Sum(scratch[:0])
			copy(out[at:], block)
			clear(scratch[:s.hash.Size()])
			break
		}
		previous = mac.Sum(out[at:at])
		at += s.hash.Size()
	}
}

type singleBlockHKDFLabel struct {
	full         []byte
	sha256Header [3]byte
	sha384Header [3]byte
}

func newSingleBlockHKDFLabel(label string) *singleBlockHKDFLabel {
	full := []byte("dtls13" + label)
	if len(full) > 255 {
		panic("dtls13: HKDF label exceeds uint8 length")
	}
	return &singleBlockHKDFLabel{
		full:         full,
		sha256Header: [3]byte{0, sha256.Size, byte(len(full))},
		sha384Header: [3]byte{0, sha512.Size384, byte(len(full))},
	}
}

// Protocol label descriptors are shared immutable storage.
func expandLabelHashInto(s *cipherSuite, secret []byte, label *singleBlockHKDFLabel, context, out []byte) {
	if len(out) != s.hash.Size() {
		panic("dtls13: single-block HKDF output must equal Hash.length")
	}
	if label == nil || len(context) > 255 {
		panic("dtls13: invalid HKDF label or context")
	}
	var header []byte
	switch len(out) {
	case sha256.Size:
		header = label.sha256Header[:]
	case sha512.Size384:
		header = label.sha384Header[:]
	default:
		panic("dtls13: unsupported HKDF hash length")
	}
	mac := hmac.New(s.hash.New, secret)
	_, _ = mac.Write(header)
	_, _ = mac.Write(label.full)
	_, _ = mac.Write(hkdfSingleByteFields[len(context)][:])
	_, _ = mac.Write(context)
	_, _ = mac.Write(hkdfSingleByteFields[1][:])
	if sum := mac.Sum(out[:0]); len(sum) != len(out) {
		panic("dtls13: unexpected HKDF output length")
	}
}

func expandLabelHashStringInto(s *cipherSuite, secret []byte, label string, context, out []byte) {
	if len(out) != s.hash.Size() {
		panic("dtls13: single-block HKDF output must equal Hash.length")
	}
	fullLabelLen := len(dtls13LabelPrefix) + len(label)
	if fullLabelLen > 255 || len(context) > 255 {
		panic("dtls13: invalid HKDF label or context")
	}
	var header []byte
	switch len(out) {
	case sha256.Size:
		header = hkdfSHA256LabelHeaders[fullLabelLen][:]
	case sha512.Size384:
		header = hkdfSHA384LabelHeaders[fullLabelLen][:]
	default:
		panic("dtls13: unsupported HKDF hash length")
	}
	mac := hmac.New(s.hash.New, secret)
	_, _ = mac.Write(header)
	_, _ = mac.Write(dtls13LabelPrefix)
	writeHashString(mac, label)
	_, _ = mac.Write(hkdfSingleByteFields[len(context)][:])
	_, _ = mac.Write(context)
	_, _ = mac.Write(hkdfSingleByteFields[1][:])
	if sum := mac.Sum(out[:0]); len(sum) != len(out) {
		panic("dtls13: unexpected HKDF output length")
	}
}

func writeHashString(h hash.Hash, value string) {
	if len(value) == 0 {
		return
	}
	// hash.Hash.Write consumes the bytes synchronously and does not retain them.
	view := unsafe.Slice(unsafe.StringData(value), len(value)) // #nosec G103 -- Write consumes this read-only view synchronously.
	_, _ = h.Write(view)
	runtime.KeepAlive(value)
}

func deriveSecret(s *cipherSuite, secret []byte, label *singleBlockHKDFLabel, transcriptHash []byte) []byte {
	out := make([]byte, s.hash.Size())
	deriveSecretInto(s, secret, label, transcriptHash, out)
	return out
}

func deriveSecretInto(s *cipherSuite, secret []byte, label *singleBlockHKDFLabel, transcriptHash, out []byte) {
	expandLabelHashInto(s, secret, label, transcriptHash, out)
}

type trafficKeys struct{ key, iv, sn []byte }

func deriveTrafficKeys(s *cipherSuite, trafficSecret []byte) trafficKeys {
	material := make([]byte, 2*s.keyLen+s.ivLen)
	keys := trafficKeys{
		key: material[:s.keyLen],
		iv:  material[s.keyLen : s.keyLen+s.ivLen],
		sn:  material[s.keyLen+s.ivLen:],
	}
	deriveTrafficKeysInto(s, trafficSecret, keys.key, keys.iv, keys.sn)
	return keys
}

func deriveTrafficKeysInto(s *cipherSuite, trafficSecret, key, iv, sn []byte) {
	if len(key) != s.keyLen || len(iv) != s.ivLen || len(sn) != s.keyLen {
		panic("dtls13: invalid traffic key destination sizes")
	}
	expander := newHKDFExpander(s.hash.New, trafficSecret)
	expandLabelWithInto(&expander, "key", nil, key)
	expandLabelWithInto(&expander, "iv", nil, iv)
	expandLabelWithInto(&expander, "sn", nil, sn)
}
