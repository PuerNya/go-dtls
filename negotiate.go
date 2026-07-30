package dtls13

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/tls"
	"errors"
)

func selectCipherSuite(preferences, offered []uint16) (*cipherSuite, error) {
	for _, preferred := range preferences {
		for _, candidate := range offered {
			if candidate == preferred {
				return cipherSuiteForID(candidate)
			}
		}
	}
	return nil, alertError(alertHandshakeFailure, errors.New("dtls13: no mutually supported cipher suite"))
}

func preferExternalPSKCipherSuite(config *Config, hello *clientHello, fallback *cipherSuite) *cipherSuite {
	if len(config.ExternalPSKs) == 0 || len(hello.pskIdentities) == 0 {
		return fallback
	}
	for _, preferred := range config.CipherSuites {
		for _, offered := range hello.cipherSuites {
			if preferred != offered {
				continue
			}
			suite, _ := cipherSuiteForID(preferred)
			for _, identity := range hello.pskIdentities {
				if findExternalPSK(config, identity.identity, suite.hash) != nil {
					return suite
				}
			}
		}
	}
	return fallback
}

func selectKeyShare(preferences []tls.CurveID, shares []keyShareEntry) (keyShareEntry, error) {
	for _, preferred := range preferences {
		for _, share := range shares {
			if share.group == preferred {
				if supportedKeyExchangeGroup(share.group) {
					return share, nil
				}
			}
		}
	}
	return keyShareEntry{}, alertError(alertHandshakeFailure, errors.New("dtls13: no mutually supported key share"))
}
func signatureSchemeCompatible(signer crypto.Signer, scheme tls.SignatureScheme) bool {
	switch key := signer.Public().(type) {
	case *rsa.PublicKey:
		return scheme == tls.PSSWithSHA256 || scheme == tls.PSSWithSHA384 || scheme == tls.PSSWithSHA512
	case *ecdsa.PublicKey:
		switch key.Curve.Params().Name {
		case "P-256":
			return scheme == tls.ECDSAWithP256AndSHA256
		case "P-384":
			return scheme == tls.ECDSAWithP384AndSHA384
		case "P-521":
			return scheme == tls.ECDSAWithP521AndSHA512
		}
		return false
	case ed25519.PublicKey:
		return scheme == tls.Ed25519
	default:
		return false
	}
}
func selectSignatureScheme(signer crypto.Signer, offered []tls.SignatureScheme) (tls.SignatureScheme, error) {
	for _, preferred := range []tls.SignatureScheme{tls.Ed25519, tls.ECDSAWithP256AndSHA256, tls.PSSWithSHA256, tls.ECDSAWithP384AndSHA384, tls.PSSWithSHA384, tls.ECDSAWithP521AndSHA512, tls.PSSWithSHA512} {
		for _, candidate := range offered {
			if candidate == preferred && signatureSchemeCompatible(signer, candidate) {
				return candidate, nil
			}
		}
	}
	return 0, alertError(alertHandshakeFailure, errors.New("dtls13: no compatible certificate signature scheme"))
}
func negotiateALPN(server, client []string) (string, error) {
	if len(server) == 0 || len(client) == 0 {
		return "", nil
	}
	for _, serverProto := range server {
		for _, clientProto := range client {
			if serverProto == clientProto {
				return serverProto, nil
			}
		}
	}
	return "", alertError(alertNoApplicationProtocol, errors.New("dtls13: no mutually supported ALPN protocol"))
}
