package dtls13

import (
	"crypto"
	"crypto/ecdh"
	"crypto/mlkem"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"io"
)

type ephemeralKey struct {
	group   tls.CurveID
	private *ecdh.PrivateKey
}

type clientKeyExchange interface {
	groupID() tls.CurveID
	publicBytes() []byte
	fallbackPublicBytes() (tls.CurveID, []byte, bool)
	sharedSecret(tls.CurveID, []byte) ([]byte, error)
}

type hybridEphemeralKey struct {
	group   tls.CurveID
	private *ecdh.PrivateKey
	mlkem   crypto.Decapsulator
}

type hybridKeyExchange struct {
	ecdhGroup           tls.CurveID
	ecdhCurve           ecdh.Curve
	ecdhElementSize     int
	mlkemPublicKeySize  int
	mlkemCiphertextSize int
	newMLKEMPrivateKey  func([]byte) (crypto.Decapsulator, error)
	newMLKEMPublicKey   func([]byte) (crypto.Encapsulator, error)
	mlkemFirst          bool
}

func hybridKeyExchangeForID(group tls.CurveID) (hybridKeyExchange, bool) {
	switch group {
	case tls.X25519MLKEM768:
		return hybridKeyExchange{
			ecdhGroup: tls.X25519, ecdhCurve: ecdh.X25519(), ecdhElementSize: 32,
			mlkemPublicKeySize: mlkem.EncapsulationKeySize768, mlkemCiphertextSize: mlkem.CiphertextSize768,
			newMLKEMPrivateKey: func(seed []byte) (crypto.Decapsulator, error) { return mlkem.NewDecapsulationKey768(seed) },
			newMLKEMPublicKey:  func(key []byte) (crypto.Encapsulator, error) { return mlkem.NewEncapsulationKey768(key) },
			mlkemFirst:         true,
		}, true
	case tls.SecP256r1MLKEM768:
		return hybridKeyExchange{
			ecdhGroup: tls.CurveP256, ecdhCurve: ecdh.P256(), ecdhElementSize: 65,
			mlkemPublicKeySize: mlkem.EncapsulationKeySize768, mlkemCiphertextSize: mlkem.CiphertextSize768,
			newMLKEMPrivateKey: func(seed []byte) (crypto.Decapsulator, error) { return mlkem.NewDecapsulationKey768(seed) },
			newMLKEMPublicKey:  func(key []byte) (crypto.Encapsulator, error) { return mlkem.NewEncapsulationKey768(key) },
		}, true
	case tls.SecP384r1MLKEM1024:
		return hybridKeyExchange{
			ecdhGroup: tls.CurveP384, ecdhCurve: ecdh.P384(), ecdhElementSize: 97,
			mlkemPublicKeySize: mlkem.EncapsulationKeySize1024, mlkemCiphertextSize: mlkem.CiphertextSize1024,
			newMLKEMPrivateKey: func(seed []byte) (crypto.Decapsulator, error) { return mlkem.NewDecapsulationKey1024(seed) },
			newMLKEMPublicKey:  func(key []byte) (crypto.Encapsulator, error) { return mlkem.NewEncapsulationKey1024(key) },
		}, true
	default:
		return hybridKeyExchange{}, false
	}
}

func curveForID(group tls.CurveID) (ecdh.Curve, error) {
	switch group {
	case tls.X25519:
		return ecdh.X25519(), nil
	case tls.CurveP256:
		return ecdh.P256(), nil
	case tls.CurveP384:
		return ecdh.P384(), nil
	default:
		return nil, errors.New("dtls13: unsupported key exchange group")
	}
}

func supportedKeyExchangeGroup(group tls.CurveID) bool {
	if _, err := curveForID(group); err == nil {
		return true
	}
	_, ok := hybridKeyExchangeForID(group)
	return ok
}

func generateEphemeralKey(group tls.CurveID, random io.Reader) (clientKeyExchange, error) {
	if random == nil {
		random = rand.Reader
	}
	curve, err := curveForID(group)
	if err == nil {
		private, err := curve.GenerateKey(random)
		if err != nil {
			return nil, err
		}
		return &ephemeralKey{group: group, private: private}, nil
	}
	hybrid, ok := hybridKeyExchangeForID(group)
	if !ok {
		return nil, err
	}
	private, err := hybrid.ecdhCurve.GenerateKey(random)
	if err != nil {
		return nil, err
	}
	var seed [mlkem.SeedSize]byte
	if _, err = io.ReadFull(random, seed[:]); err != nil {
		return nil, err
	}
	decapsulationKey, err := hybrid.newMLKEMPrivateKey(seed[:])
	clear(seed[:])
	if err != nil {
		return nil, err
	}
	return &hybridEphemeralKey{group: group, private: private, mlkem: decapsulationKey}, nil
}

func (k *ephemeralKey) groupID() tls.CurveID { return k.group }

func (k *ephemeralKey) publicBytes() []byte {
	return append([]byte(nil), k.private.PublicKey().Bytes()...)
}

func (*ephemeralKey) fallbackPublicBytes() (tls.CurveID, []byte, bool) {
	return 0, nil, false
}

func (k *ephemeralKey) sharedSecret(peerGroup tls.CurveID, peerPublic []byte) ([]byte, error) {
	if k == nil || k.private == nil {
		return nil, errors.New("dtls13: missing local ephemeral key")
	}
	if peerGroup != k.group {
		return nil, &ProtocolError{"peer key share group does not match local key"}
	}
	return ecdhSharedSecret(k.private, peerPublic)
}

func (k *hybridEphemeralKey) groupID() tls.CurveID { return k.group }

func (k *hybridEphemeralKey) publicBytes() []byte {
	ecdhPublic := k.private.PublicKey().Bytes()
	hybrid, _ := hybridKeyExchangeForID(k.group)
	mlkemPublic := k.mlkem.Encapsulator().Bytes()
	public := make([]byte, 0, hybrid.ecdhElementSize+hybrid.mlkemPublicKeySize)
	if hybrid.mlkemFirst {
		public = append(public, mlkemPublic...)
		return append(public, ecdhPublic...)
	}
	public = append(public, ecdhPublic...)
	return append(public, mlkemPublic...)
}

func (k *hybridEphemeralKey) fallbackPublicBytes() (tls.CurveID, []byte, bool) {
	if k == nil || k.private == nil || k.mlkem == nil {
		return 0, nil, false
	}
	hybrid, _ := hybridKeyExchangeForID(k.group)
	return hybrid.ecdhGroup, append([]byte(nil), k.private.PublicKey().Bytes()...), true
}

func (k *hybridEphemeralKey) sharedSecret(peerGroup tls.CurveID, peerPublic []byte) ([]byte, error) {
	if k == nil || k.private == nil || k.mlkem == nil {
		return nil, errors.New("dtls13: missing local ephemeral key")
	}
	hybrid, _ := hybridKeyExchangeForID(k.group)
	if peerGroup == hybrid.ecdhGroup {
		return ecdhSharedSecret(k.private, peerPublic)
	}
	if peerGroup != k.group {
		return nil, &ProtocolError{"peer key share group does not match local key"}
	}
	if len(peerPublic) != hybrid.ecdhElementSize+hybrid.mlkemCiphertextSize {
		return nil, &ProtocolError{"invalid hybrid server key share length"}
	}
	ecdhShare, mlkemCiphertext := splitHybridShare(hybrid, peerPublic, hybrid.mlkemCiphertextSize)
	ecdhSecret, err := ecdhSharedSecret(k.private, ecdhShare)
	if err != nil {
		return nil, err
	}
	mlkemSecret, err := k.mlkem.Decapsulate(mlkemCiphertext)
	if err != nil {
		return nil, &ProtocolError{"invalid ML-KEM ciphertext"}
	}
	return joinHybridSecret(hybrid, ecdhSecret, mlkemSecret), nil
}

func generateServerKeyShare(group tls.CurveID, clientShare []byte, random io.Reader) ([]byte, []byte, error) {
	hybrid, ok := hybridKeyExchangeForID(group)
	if !ok {
		curve, err := curveForID(group)
		if err != nil {
			return nil, nil, alertError(alertInternalError, err)
		}
		if random == nil {
			random = rand.Reader
		}
		private, err := curve.GenerateKey(random)
		if err != nil {
			return nil, nil, alertError(alertInternalError, err)
		}
		secret, err := ecdhSharedSecret(private, clientShare)
		return append([]byte(nil), private.PublicKey().Bytes()...), secret, err
	}
	if len(clientShare) != hybrid.ecdhElementSize+hybrid.mlkemPublicKeySize {
		return nil, nil, &ProtocolError{"invalid hybrid client key share length"}
	}
	ecdhShare, mlkemPublic := splitHybridShare(hybrid, clientShare, hybrid.mlkemPublicKeySize)
	peerECDH, err := hybrid.ecdhCurve.NewPublicKey(ecdhShare)
	if err != nil {
		return nil, nil, &ProtocolError{"invalid peer key share"}
	}
	peerMLKEM, err := hybrid.newMLKEMPublicKey(mlkemPublic)
	if err != nil {
		return nil, nil, &ProtocolError{"invalid ML-KEM public key"}
	}
	if random == nil {
		random = rand.Reader
	}
	private, err := hybrid.ecdhCurve.GenerateKey(random)
	if err != nil {
		return nil, nil, alertError(alertInternalError, err)
	}
	ecdhSecret, err := private.ECDH(peerECDH)
	if err != nil {
		return nil, nil, &ProtocolError{"invalid peer key share"}
	}
	mlkemSecret, ciphertext := peerMLKEM.Encapsulate()
	ecdhPublic := private.PublicKey().Bytes()
	serverShare := make([]byte, 0, hybrid.ecdhElementSize+hybrid.mlkemCiphertextSize)
	if hybrid.mlkemFirst {
		serverShare = append(serverShare, ciphertext...)
		serverShare = append(serverShare, ecdhPublic...)
	} else {
		serverShare = append(serverShare, ecdhPublic...)
		serverShare = append(serverShare, ciphertext...)
	}
	return serverShare, joinHybridSecret(hybrid, ecdhSecret, mlkemSecret), nil
}

func ecdhSharedSecret(private *ecdh.PrivateKey, peerPublic []byte) ([]byte, error) {
	public, err := private.Curve().NewPublicKey(peerPublic)
	if err != nil {
		return nil, &ProtocolError{"invalid peer key share"}
	}
	secret, err := private.ECDH(public)
	if err != nil {
		return nil, &ProtocolError{"invalid peer key share"}
	}
	return secret, nil
}

func splitHybridShare(hybrid hybridKeyExchange, share []byte, mlkemSize int) (ecdhShare, mlkemShare []byte) {
	if hybrid.mlkemFirst {
		return share[mlkemSize:], share[:mlkemSize]
	}
	return share[:hybrid.ecdhElementSize], share[hybrid.ecdhElementSize:]
}

func joinHybridSecret(hybrid hybridKeyExchange, ecdhSecret, mlkemSecret []byte) []byte {
	secret := make([]byte, 0, len(ecdhSecret)+len(mlkemSecret))
	if hybrid.mlkemFirst {
		secret = append(secret, mlkemSecret...)
		return append(secret, ecdhSecret...)
	}
	secret = append(secret, ecdhSecret...)
	return append(secret, mlkemSecret...)
}
