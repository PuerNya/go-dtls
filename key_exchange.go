package dtls13

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"io"
)

type ephemeralKey struct {
	group   tls.CurveID
	private *ecdh.PrivateKey
}

func curveForID(group tls.CurveID) (ecdh.Curve, error) {
	switch group {
	case tls.X25519:
		return ecdh.X25519(), nil
	case tls.CurveP256:
		return ecdh.P256(), nil
	default:
		return nil, errors.New("dtls13: unsupported key exchange group")
	}
}

func generateEphemeralKey(group tls.CurveID, random io.Reader) (*ephemeralKey, error) {
	curve, err := curveForID(group)
	if err != nil {
		return nil, err
	}
	if random == nil {
		random = rand.Reader
	}
	private, err := curve.GenerateKey(random)
	if err != nil {
		return nil, err
	}
	return &ephemeralKey{group: group, private: private}, nil
}
func (k *ephemeralKey) publicBytes() []byte {
	return append([]byte(nil), k.private.PublicKey().Bytes()...)
}
func (k *ephemeralKey) sharedSecret(peerGroup tls.CurveID, peerPublic []byte) ([]byte, error) {
	if k == nil || k.private == nil {
		return nil, errors.New("dtls13: missing local ephemeral key")
	}
	if peerGroup != k.group {
		return nil, errors.New("dtls13: peer key share group does not match local key")
	}
	curve, err := curveForID(peerGroup)
	if err != nil {
		return nil, err
	}
	public, err := curve.NewPublicKey(peerPublic)
	if err != nil {
		return nil, &ProtocolError{"invalid peer key share"}
	}
	secret, err := k.private.ECDH(public)
	if err != nil {
		return nil, &ProtocolError{"invalid peer key share"}
	}
	return secret, nil
}
