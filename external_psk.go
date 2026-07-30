package dtls13

import (
	"crypto"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

const (
	externalPSKMinimumLength        = 16
	tlsKDFHKDFSHA256         uint16 = 0x0001
	tlsKDFHKDFSHA384         uint16 = 0x0002
)

// ExternalPSK is immutable externally provisioned key material for DTLS 1.3.
// Create values with [ImportExternalPSK] whenever possible. The direct
// constructor exists for deployments that already provision TLS-specific PSKs.
//
// PSK identities and importer contexts are sent in plaintext in ClientHello.
// Reusing either value makes connections linkable, and sensitive information
// must not be placed in them.
type ExternalPSK struct {
	identity []byte
	context  []byte
	keys     []externalPSKKey
}

type externalPSKKey struct {
	wireIdentity []byte
	key          []byte
	hash         crypto.Hash
	binderLabel  *singleBlockHKDFLabel
	digest       [sha256.Size]byte
}

type externalPSKSelection struct {
	psk *ExternalPSK
	key *externalPSKKey
}

// NewDirectExternalPSK configures a non-imported external PSK. hash associates
// the key with a TLS 1.3 cipher-suite hash and must be crypto.SHA256 or
// crypto.SHA384. A zero hash selects SHA-256.
//
// key must contain at least 128 bits, and identity must contain 1..65535 bytes.
// The inputs are copied. Prefer [ImportExternalPSK], which separates keys by
// protocol and target KDF as recommended by RFC 9257 and RFC 9258.
func NewDirectExternalPSK(identity, key []byte, hash crypto.Hash) (*ExternalPSK, error) {
	if err := validateExternalPSKInput(identity, key); err != nil {
		return nil, err
	}
	if hash == 0 {
		hash = crypto.SHA256
	}
	if _, ok := tlsKDFForHash(hash); !ok {
		return nil, errors.New("dtls13: external PSK hash must be SHA-256 or SHA-384")
	}
	keyCopy := append([]byte(nil), key...)
	return &ExternalPSK{
		identity: append([]byte(nil), identity...),
		keys: []externalPSKKey{{
			wireIdentity: append([]byte(nil), identity...),
			key:          keyCopy,
			hash:         hash,
			binderLabel:  labelExternalBinder,
			digest:       sha256.Sum256(keyCopy),
		}},
	}, nil
}

// ImportExternalPSK imports an external PSK for DTLS 1.3 according to RFC
// 9258. It derives separate SHA-256 and SHA-384 target keys with the mandatory
// "dtls13" label prefix and uses the "imp binder" binder key label.
//
// sourceHash is the hash associated with the externally provisioned key. A
// zero value selects SHA-256 as recommended when no hash is associated. key
// must contain at least 128 bits. identity must be non-empty, and the serialized
// ImportedIdentity must fit the TLS 16-bit PSK identity vector. Inputs are
// copied; the original key is not retained after both target keys are derived.
//
// context must include any context and channel binding used to derive key.
// Deployments using a group key should bind the client and server identities in
// context. The context is visible on the wire and must not contain secrets.
func ImportExternalPSK(identity, key, context []byte, sourceHash crypto.Hash) (*ExternalPSK, error) {
	if err := validateExternalPSKInput(identity, key); err != nil {
		return nil, err
	}
	if sourceHash == 0 {
		sourceHash = crypto.SHA256
	}
	if !sourceHash.Available() {
		return nil, errors.New("dtls13: external PSK source hash is unavailable")
	}
	if len(context) > 65535 {
		return nil, errors.New("dtls13: external PSK context exceeds 65535 bytes")
	}
	psk := &ExternalPSK{
		identity: append([]byte(nil), identity...),
		context:  append([]byte(nil), context...),
		keys:     make([]externalPSKKey, 0, 2),
	}
	for _, target := range []crypto.Hash{crypto.SHA256, crypto.SHA384} {
		wireIdentity, err := marshalImportedIdentity(identity, context, target)
		if err != nil {
			return nil, err
		}
		imported := deriveImportedPSK(key, wireIdentity, sourceHash, target.Size())
		psk.keys = append(psk.keys, externalPSKKey{
			wireIdentity: wireIdentity,
			key:          imported,
			hash:         target,
			binderLabel:  labelImportedBinder,
			digest:       sha256.Sum256(imported),
		})
	}
	return psk, nil
}

func validateExternalPSKInput(identity, key []byte) error {
	if len(identity) == 0 || len(identity) > 65535 {
		return errors.New("dtls13: external PSK identity must contain 1..65535 bytes")
	}
	if len(key) < externalPSKMinimumLength {
		return errors.New("dtls13: external PSK must contain at least 128 bits")
	}
	return nil
}

func marshalImportedIdentity(identity, context []byte, targetHash crypto.Hash) ([]byte, error) {
	targetKDF, ok := tlsKDFForHash(targetHash)
	if !ok {
		return nil, errors.New("dtls13: unsupported external PSK target KDF")
	}
	length := 2 + len(identity) + 2 + len(context) + 2 + 2
	if length > 65535 {
		return nil, errors.New("dtls13: imported PSK identity exceeds 65535 bytes")
	}
	w := newWireBuilder(length)
	w.bytes16(identity)
	w.bytes16(context)
	w.u16(int(VersionDTLS13))
	w.u16(int(targetKDF))
	return w.b, w.err
}

func deriveImportedPSK(key, importedIdentity []byte, sourceHash crypto.Hash, length int) []byte {
	h := sourceHash.New()
	_, _ = h.Write(importedIdentity)
	context := h.Sum(nil)
	extracted := hkdfExtract(sourceHash.New, key, nil)
	const label = "dtls13derived psk"
	info := make([]byte, 2+1+len(label)+1+len(context))
	binary.BigEndian.PutUint16(info, uint16(length))
	info[2] = byte(len(label))
	copy(info[3:], label)
	at := 3 + len(label)
	info[at] = byte(len(context))
	copy(info[at+1:], context)
	return hkdfExpand(sourceHash.New, extracted, info, length)
}

func tlsKDFForHash(hash crypto.Hash) (uint16, bool) {
	switch hash {
	case crypto.SHA256:
		return tlsKDFHKDFSHA256, true
	case crypto.SHA384:
		return tlsKDFHKDFSHA384, true
	default:
		return 0, false
	}
}

func hashForTLSKDF(kdf uint16) (crypto.Hash, bool) {
	switch kdf {
	case tlsKDFHKDFSHA256:
		return crypto.SHA256, true
	case tlsKDFHKDFSHA384:
		return crypto.SHA384, true
	default:
		return 0, false
	}
}

func findExternalPSK(config *Config, wireIdentity []byte, hash crypto.Hash) *externalPSKSelection {
	for _, psk := range config.ExternalPSKs {
		for i := range psk.keys {
			key := &psk.keys[i]
			if key.hash == hash && hmac.Equal(key.wireIdentity, wireIdentity) {
				return &externalPSKSelection{psk: psk, key: key}
			}
		}
	}
	return nil
}

func matchExternalPSK(config *Config, wireIdentity []byte, hash crypto.Hash, digest [sha256.Size]byte) *externalPSKSelection {
	selection := findExternalPSK(config, wireIdentity, hash)
	if selection == nil || !hmac.Equal(selection.key.digest[:], digest[:]) {
		return nil
	}
	return selection
}
