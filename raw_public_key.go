package dtls13

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"slices"
)

// CertificateType identifies the credential carried in a Certificate message.
type CertificateType uint8

const (
	// CertificateTypeX509 carries a DER-encoded X.509 certificate chain.
	CertificateTypeX509 CertificateType = 0
	// CertificateTypeRawPublicKey carries one DER-encoded SubjectPublicKeyInfo.
	CertificateTypeRawPublicKey CertificateType = 2
)

func supportsCertificateType(types []CertificateType, typ CertificateType) bool {
	if len(types) == 0 {
		return typ == CertificateTypeX509
	}
	return slices.Contains(types, typ)
}

func validateCertificateTypes(types []CertificateType) error {
	for i, typ := range types {
		if typ != CertificateTypeX509 && typ != CertificateTypeRawPublicKey {
			return &ConfigError{"unsupported certificate type"}
		}
		if slices.Contains(types[:i], typ) {
			return &ConfigError{"duplicate certificate type"}
		}
	}
	return nil
}

func parseCertificateTypes(raw []byte) ([]byte, error) {
	if len(raw) < 2 || int(raw[0]) != len(raw)-1 {
		return nil, alertError(alertDecodeError, &ProtocolError{"invalid certificate type list"})
	}
	return raw[1:], nil
}

func (c *Conn) offerCertificateTypes(hello *clientHello) error {
	for _, offer := range []struct {
		typ     uint16
		types   []CertificateType
		sending bool
	}{
		{extServerCertificateType, c.config.ServerCertificateTypes, false},
		{extClientCertificateType, c.config.ClientCertificateTypes, true},
	} {
		if !slices.Contains(offer.types, CertificateTypeRawPublicKey) {
			continue
		}
		if !offer.sending && c.config.VerifyPeerRawPublicKey == nil {
			if !supportsCertificateType(offer.types, CertificateTypeX509) {
				return &ConfigError{"RPK server authentication requires VerifyPeerRawPublicKey"}
			}
			continue
		}
		var types []byte
		for _, typ := range offer.types {
			if offer.sending && ((typ == CertificateTypeRawPublicKey && c.config.RawPublicKeySigner == nil) ||
				(typ == CertificateTypeX509 && len(c.config.Certificates) == 0 && c.config.GetClientCertificate == nil)) {
				continue
			}
			types = append(types, byte(typ))
		}
		if !slices.Contains(types, byte(CertificateTypeRawPublicKey)) {
			continue
		}
		if hello.unknownExtensions == nil {
			hello.unknownExtensions = make(map[uint16][]byte)
		}
		hello.unknownExtensions[offer.typ] = append([]byte{byte(len(types))}, types...)
	}
	return nil
}

func selectCertificateType(raw []byte, configured []CertificateType, x509Available, rawAvailable bool) (CertificateType, error) {
	offered := []byte{byte(CertificateTypeX509)}
	if raw != nil {
		var err error
		offered, err = parseCertificateTypes(raw)
		if err != nil {
			return 0, err
		}
	}
	if len(configured) == 0 {
		configured = []CertificateType{CertificateTypeX509}
	}
	for _, typ := range configured {
		if (typ == CertificateTypeX509 && !x509Available) || (typ == CertificateTypeRawPublicKey && !rawAvailable) {
			continue
		}
		if slices.Contains(offered, byte(typ)) {
			return typ, nil
		}
	}
	return 0, alertError(alertUnsupportedCertificate, &ProtocolError{"no common certificate type"})
}

func validateSelectedCertificateType(hello *clientHello, typ uint16, raw []byte) (CertificateType, error) {
	offer, ok := hello.unknownExtensions[typ]
	if !ok {
		return 0, alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited certificate type selection"})
	}
	if len(raw) != 1 {
		return 0, alertError(alertDecodeError, &ProtocolError{"invalid certificate type selection"})
	}
	types, err := parseCertificateTypes(offer)
	if err != nil {
		return 0, err
	}
	if !slices.Contains(types, raw[0]) || (raw[0] != byte(CertificateTypeX509) && raw[0] != byte(CertificateTypeRawPublicKey)) {
		return 0, alertError(alertIllegalParameter, &ProtocolError{"unoffered certificate type selection"})
	}
	return CertificateType(raw[0]), nil
}

func validateRawPublicKey(key crypto.PublicKey) error {
	switch key := key.(type) {
	case *rsa.PublicKey:
		if key.N == nil || key.N.BitLen() < 2048 {
			return errors.New("dtls13: RSA raw public key is smaller than 2048 bits")
		}
	case *ecdsa.PublicKey:
		if key.Curve != elliptic.P256() && key.Curve != elliptic.P384() && key.Curve != elliptic.P521() {
			return errors.New("dtls13: unsupported raw public key curve")
		}
	case ed25519.PublicKey:
	default:
		return errors.New("dtls13: unsupported raw public key algorithm")
	}
	return nil
}

func (c *Conn) rawPublicKeyCertificate() (*tls.Certificate, error) {
	signer := c.config.RawPublicKeySigner
	if signer == nil {
		return nil, nil
	}
	public := signer.Public()
	if err := validateRawPublicKey(public); err != nil {
		return nil, err
	}
	der, err := x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return nil, err
	}
	return &tls.Certificate{Certificate: [][]byte{der}, PrivateKey: signer}, nil
}

func verifyRawPublicKey(config *Config, raw []byte) error {
	key, err := x509.ParsePKIXPublicKey(raw)
	if err == nil {
		err = validateRawPublicKey(key)
	}
	if err != nil {
		return alertError(alertBadCertificate, err)
	}
	if config.VerifyPeerRawPublicKey == nil {
		return alertError(alertBadCertificate, errors.New("dtls13: no raw public key trust verifier"))
	}
	if err = config.VerifyPeerRawPublicKey(bytes.Clone(raw)); err != nil {
		return alertError(alertAccessDenied, err)
	}
	return nil
}

func verifyRawPublicKeyMessage(config *Config, message *certificateMessage) ([]byte, error) {
	if len(message.certificates) == 0 {
		return nil, alertError(alertDecodeError, &ProtocolError{"server sent an empty Certificate"})
	}
	if len(message.certificates) != 1 {
		return nil, alertError(alertBadCertificate, &ProtocolError{"RPK Certificate must contain exactly one public key"})
	}
	raw := message.certificates[0].data
	if err := verifyRawPublicKey(config, raw); err != nil {
		return nil, err
	}
	return bytes.Clone(raw), nil
}

func verifyPeerCertificateSignature(certificates []*x509.Certificate, raw []byte, scheme tls.SignatureScheme, transcript, signature []byte, server bool) error {
	var key crypto.PublicKey
	if len(raw) > 0 {
		var err error
		key, err = x509.ParsePKIXPublicKey(raw)
		if err != nil {
			return err
		}
	} else if len(certificates) > 0 {
		key = certificates[0].PublicKey
	} else {
		return errors.New("dtls13: peer authentication key is missing")
	}
	return verifyCertificateVerify(key, scheme, transcript, signature, server)
}
