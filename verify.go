package dtls13

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
)

var oidExtensionKeyUsage = asn1.ObjectIdentifier{2, 5, 29, 15}

func verifyCertificateChain(config *Config, message *certificateMessage, peerIsServer bool, signatureSchemes []tls.SignatureScheme) ([]*x509.Certificate, [][]*x509.Certificate, error) {
	if len(message.certificates) == 0 {
		if peerIsServer {
			return nil, nil, alertError(alertDecodeError, errors.New("dtls13: server sent an empty certificate chain"))
		}
		return nil, nil, errors.New("dtls13: peer sent an empty certificate chain")
	}
	raw := make([][]byte, len(message.certificates))
	certs := make([]*x509.Certificate, len(message.certificates))
	for i, entry := range message.certificates {
		raw[i] = entry.data
		cert, err := x509.ParseCertificate(entry.data)
		if err != nil {
			return nil, nil, alertError(alertBadCertificate, fmt.Errorf("dtls13: parse peer certificate %d: %w", i, err))
		}
		certs[i] = cert
	}
	for _, extension := range certs[0].Extensions {
		if extension.Id.Equal(oidExtensionKeyUsage) && certs[0].KeyUsage&x509.KeyUsageDigitalSignature == 0 {
			return nil, nil, alertError(alertBadCertificate, errors.New("dtls13: peer certificate does not permit digital signatures"))
		}
	}
	if err := validateCertificateSignatureAlgorithms(certs, signatureSchemes); err != nil {
		return nil, nil, alertError(alertBadCertificate, err)
	}
	var chains [][]*x509.Certificate
	if !config.InsecureSkipVerify {
		intermediates := x509.NewCertPool()
		for _, cert := range certs[1:] {
			intermediates.AddCert(cert)
		}
		opts := x509.VerifyOptions{Intermediates: intermediates, CurrentTime: config.Time()}
		if peerIsServer {
			opts.Roots = config.RootCAs
			opts.DNSName = config.ServerName
			opts.KeyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		} else {
			opts.Roots = config.ClientCAs
			opts.KeyUsages = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
		}
		verified, err := certs[0].Verify(opts)
		if err != nil {
			description := uint8(alertBadCertificate)
			var unknownAuthority x509.UnknownAuthorityError
			if errors.As(err, &unknownAuthority) {
				description = alertUnknownCA
			}
			return nil, nil, alertError(description, fmt.Errorf("dtls13: verify peer certificate: %w", err))
		}
		chains = verified
	}
	if config.VerifyPeerCertificate != nil {
		if err := config.VerifyPeerCertificate(raw, chains); err != nil {
			return nil, nil, alertError(alertAccessDenied, err)
		}
	}
	return certs, chains, nil
}

func certificateSignatureScheme(algorithm x509.SignatureAlgorithm) (tls.SignatureScheme, bool) {
	switch algorithm {
	case x509.SHA256WithRSA:
		return tls.PKCS1WithSHA256, true
	case x509.SHA384WithRSA:
		return tls.PKCS1WithSHA384, true
	case x509.SHA512WithRSA:
		return tls.PKCS1WithSHA512, true
	case x509.SHA256WithRSAPSS:
		return tls.PSSWithSHA256, true
	case x509.SHA384WithRSAPSS:
		return tls.PSSWithSHA384, true
	case x509.SHA512WithRSAPSS:
		return tls.PSSWithSHA512, true
	case x509.ECDSAWithSHA256:
		return tls.ECDSAWithP256AndSHA256, true
	case x509.ECDSAWithSHA384:
		return tls.ECDSAWithP384AndSHA384, true
	case x509.ECDSAWithSHA512:
		return tls.ECDSAWithP521AndSHA512, true
	case x509.PureEd25519:
		return tls.Ed25519, true
	case x509.SHA1WithRSA:
		return tls.PKCS1WithSHA1, true
	case x509.ECDSAWithSHA1:
		return tls.ECDSAWithSHA1, true
	default:
		return 0, false
	}
}

func validateCertificateSignatureAlgorithms(certificates []*x509.Certificate, offered []tls.SignatureScheme) error {
	for _, certificate := range certificates {
		selfSigned := bytes.Equal(certificate.RawIssuer, certificate.RawSubject) &&
			certificate.CheckSignature(certificate.SignatureAlgorithm, certificate.RawTBSCertificate, certificate.Signature) == nil
		if selfSigned {
			continue
		}
		scheme, ok := certificateSignatureScheme(certificate.SignatureAlgorithm)
		if !ok {
			return errors.New("dtls13: certificate uses a prohibited signature algorithm")
		}
		accepted := false
		for _, candidate := range offered {
			accepted = accepted || candidate == scheme
		}
		if !accepted {
			return fmt.Errorf("dtls13: certificate signature algorithm %v was not offered", certificate.SignatureAlgorithm)
		}
	}
	return nil
}

func validateConfiguredCertificate(certificate *tls.Certificate, offered []tls.SignatureScheme) error {
	if certificate == nil || len(certificate.Certificate) == 0 {
		return errors.New("dtls13: configured certificate chain is empty")
	}
	parsed := make([]*x509.Certificate, len(certificate.Certificate))
	for i, der := range certificate.Certificate {
		cert, err := x509.ParseCertificate(der)
		if err != nil {
			return fmt.Errorf("dtls13: parse configured certificate %d: %w", i, err)
		}
		parsed[i] = cert
	}
	for _, extension := range parsed[0].Extensions {
		if extension.Id.Equal(oidExtensionKeyUsage) && parsed[0].KeyUsage&x509.KeyUsageDigitalSignature == 0 {
			return errors.New("dtls13: configured certificate does not permit digital signatures")
		}
	}
	signer, ok := certificate.PrivateKey.(interface{ Public() crypto.PublicKey })
	if !ok {
		return errors.New("dtls13: configured certificate private key is not a signer")
	}
	certificatePublic, err := x509.MarshalPKIXPublicKey(parsed[0].PublicKey)
	if err != nil {
		return err
	}
	signerPublic, err := x509.MarshalPKIXPublicKey(signer.Public())
	if err != nil || !bytes.Equal(certificatePublic, signerPublic) {
		return errors.New("dtls13: configured certificate and private key do not match")
	}
	return validateCertificateSignatureAlgorithms(parsed, offered)
}
