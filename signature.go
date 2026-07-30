package dtls13

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
)

func certificateVerifyInput(transcriptHash []byte, server bool) []byte {
	contextString := "TLS 1.3, client CertificateVerify"
	if server {
		contextString = "TLS 1.3, server CertificateVerify"
	}
	b := make([]byte, 64, 64+len(contextString)+1+len(transcriptHash))
	for i := range b {
		b[i] = 0x20
	}
	b = append(b, contextString...)
	b = append(b, 0)
	b = append(b, transcriptHash...)
	return b
}

func signatureParameters(scheme tls.SignatureScheme) (crypto.Hash, *rsa.PSSOptions, error) {
	switch scheme {
	case tls.PSSWithSHA256:
		return crypto.SHA256, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA256}, nil
	case tls.PSSWithSHA384:
		return crypto.SHA384, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA384}, nil
	case tls.PSSWithSHA512:
		return crypto.SHA512, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash, Hash: crypto.SHA512}, nil
	case tls.ECDSAWithP256AndSHA256:
		return crypto.SHA256, nil, nil
	case tls.ECDSAWithP384AndSHA384:
		return crypto.SHA384, nil, nil
	case tls.ECDSAWithP521AndSHA512:
		return crypto.SHA512, nil, nil
	case tls.Ed25519:
		return crypto.Hash(0), nil, nil
	default:
		return 0, nil, errors.New("dtls13: unsupported TLS 1.3 signature scheme")
	}
}

func signCertificateVerify(random io.Reader, signer crypto.Signer, scheme tls.SignatureScheme, transcriptHash []byte, server bool) ([]byte, error) {
	if random == nil {
		random = rand.Reader
	}
	hash, pss, err := signatureParameters(scheme)
	if err != nil {
		return nil, err
	}
	input := certificateVerifyInput(transcriptHash, server)
	message := input
	var opts crypto.SignerOpts = hash
	if hash != 0 {
		h := hash.New()
		h.Write(input)
		message = h.Sum(nil)
	}
	if pss != nil {
		opts = pss
	}
	sig, err := signer.Sign(random, message, opts)
	if err != nil {
		return nil, fmt.Errorf("dtls13: sign CertificateVerify: %w", err)
	}
	return sig, nil
}

func verifyCertificateVerify(public crypto.PublicKey, scheme tls.SignatureScheme, transcriptHash, signature []byte, server bool) error {
	hash, pss, err := signatureParameters(scheme)
	if err != nil {
		return err
	}
	input := certificateVerifyInput(transcriptHash, server)
	digest := input
	if hash != 0 {
		h := hash.New()
		h.Write(input)
		digest = h.Sum(nil)
	}
	switch key := public.(type) {
	case *rsa.PublicKey:
		if pss == nil {
			return errors.New("dtls13: RSA key used with non-RSA-PSS signature scheme")
		}
		return rsa.VerifyPSS(key, hash, digest, signature, pss)
	case *ecdsa.PublicKey:
		if pss != nil || hash == 0 {
			return errors.New("dtls13: ECDSA key used with incompatible signature scheme")
		}
		if !ecdsa.VerifyASN1(key, digest, signature) {
			return errors.New("dtls13: invalid ECDSA CertificateVerify signature")
		}
		return nil
	case ed25519.PublicKey:
		if scheme != tls.Ed25519 {
			return errors.New("dtls13: Ed25519 key used with incompatible signature scheme")
		}
		if !ed25519.Verify(key, input, signature) {
			return errors.New("dtls13: invalid Ed25519 CertificateVerify signature")
		}
		return nil
	default:
		return errors.New("dtls13: unsupported certificate public key")
	}
}
