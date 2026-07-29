package dtls13

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding"
	"hash"
)

const handshakeTypeMessageHash uint8 = 254

// transcriptHash adds a complete handshake message in the original TLS 1.3
// Handshake form required by RFC 9147 section 5.2. DTLS-only message sequence
// and fragmentation fields are never included.
type transcriptHash struct {
	h hash.Hash
}

func newTranscriptHash(h hash.Hash) *transcriptHash { return &transcriptHash{h: h} }
func (t *transcriptHash) add(typ uint8, _ uint16, body []byte) error {
	if len(body) >= 1<<24 {
		return &ProtocolError{"handshake message exceeds 24-bit length"}
	}
	var header [4]byte
	header[0] = typ
	putUint24(header[1:4], uint32(len(body)))
	_, _ = t.h.Write(header[:])
	t.h.Write(body)
	return nil
}
func (t *transcriptHash) sum() []byte               { return t.h.Sum(nil) }
func (t *transcriptHash) sumInto(dst []byte) []byte { return t.h.Sum(dst[:0]) }
func (t *transcriptHash) clone() *transcriptHash {
	var h hash.Hash
	switch t.h.Size() {
	case sha256.Size:
		h = sha256.New()
	case sha512.Size384:
		h = sha512.New384()
	default:
		panic("dtls13: unsupported transcript hash size")
	}
	marshaler, ok := t.h.(encoding.BinaryMarshaler)
	if !ok {
		panic("dtls13: transcript hash does not support state snapshots")
	}
	state, err := marshaler.MarshalBinary()
	if err != nil {
		panic("dtls13: failed to snapshot transcript hash: " + err.Error())
	}
	unmarshaler, ok := h.(encoding.BinaryUnmarshaler)
	if !ok {
		panic("dtls13: transcript hash does not support state restoration")
	}
	if err = unmarshaler.UnmarshalBinary(state); err != nil {
		panic("dtls13: failed to restore transcript hash: " + err.Error())
	}
	return newTranscriptHash(h)
}

// addHelloRetryRequest applies RFC 8446's synthetic message_hash rewrite used
// by RFC 9147 section 5.1. It must be called on a fresh transcript.
func (t *transcriptHash) addHelloRetryRequest(clientHelloHash []byte, hrrBody []byte) error {
	if err := t.add(handshakeTypeMessageHash, 0, clientHelloHash); err != nil {
		return err
	}
	return t.add(2, 0, hrrBody)
}
