package dtls13

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"slices"
	"time"
)

const extDelegatedCredential uint16 = 34

var oidDelegationUsage = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 44363, 44}

// DelegatedCredential is an RFC 9345 credential and its parent certificate
// chain. Certificate.PrivateKey is the delegated signing key, not the parent
// key. Credential contains the encoded DelegatedCredential, including its
// delegation signature. The parent private key need not be available here.
// Like tls.Certificate, its contents must not be modified after first use.
type DelegatedCredential struct {
	Certificate tls.Certificate
	Credential  []byte
}

// NewDelegatedCredential issues a server or client credential using the parent
// certificate's private key. The certificate must permit delegation. key is
// the delegated signing key (ECDSA or Ed25519); validUntil must be within seven
// days of now and strictly before the parent certificate expires. The returned
// credential retains key, but not the parent private key.
func NewDelegatedCredential(parent tls.Certificate, key crypto.Signer, validUntil time.Time, server bool) (*DelegatedCredential, error) {
	chain, err := validateConfiguredCertificateChain(&parent, defaultSignatureSchemes(), server)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return nil, errors.New("dtls13: missing delegated private key")
	}
	scheme, err := selectSignatureScheme(key, delegatedCredentialSchemes())
	if err != nil {
		return nil, err
	}
	signer := parent.PrivateKey.(crypto.Signer)
	algorithm, err := selectSignatureScheme(signer, defaultSignatureSchemes())
	if err != nil {
		return nil, err
	}
	seconds := validUntil.Sub(chain[0].NotBefore) / time.Second
	if seconds < 0 || seconds > 1<<32-1 {
		return nil, errors.New("dtls13: delegated credential expiry cannot be encoded")
	}
	public, err := x509.MarshalPKIXPublicKey(key.Public())
	if err != nil {
		return nil, err
	}
	w := newWireBuilder(9 + len(public))
	w.u32(uint32(seconds))
	w.u16(int(scheme))
	w.bytes24(public)
	w.u16(int(algorithm))
	if w.err != nil {
		return nil, w.err
	}
	signature, err := signMessage(rand.Reader, signer, algorithm, delegatedCredentialInput(chain[0].Raw, w.b, server))
	if err != nil {
		return nil, err
	}
	w.bytes16(signature)
	if w.err != nil {
		return nil, w.err
	}
	if _, err = validateDelegatedCredential(w.b, chain[0], time.Now(), server, defaultSignatureSchemes(), delegatedCredentialSchemes()); err != nil {
		return nil, err
	}
	parent.PrivateKey = key
	return &DelegatedCredential{Certificate: parent, Credential: w.b}, nil
}

func delegatedCredentialSchemes() []tls.SignatureScheme {
	// RSA DC keys require RSASSA-PSS SubjectPublicKeyInfo, which crypto/x509
	// does not support. rsa_pss_rsae schemes are forbidden by RFC 9345 §4.
	return []tls.SignatureScheme{tls.Ed25519, tls.ECDSAWithP256AndSHA256, tls.ECDSAWithP384AndSHA384, tls.ECDSAWithP521AndSHA512}
}

type parsedDelegatedCredential struct {
	validTime uint32
	scheme    tls.SignatureScheme
	algorithm tls.SignatureScheme
	public    crypto.PublicKey
	signed    []byte
	signature []byte
}

func parseDelegatedCredential(raw []byte) (*parsedDelegatedCredential, error) {
	p := wireParser{b: raw}
	d := &parsedDelegatedCredential{validTime: p.u32(), scheme: tls.SignatureScheme(p.u16())}
	spki := p.bytes24()
	d.algorithm = tls.SignatureScheme(p.u16())
	d.signed = raw[:p.off]
	d.signature = p.bytes16()
	if err := p.done(); err != nil {
		return nil, err
	}
	if len(spki) == 0 || len(d.signature) == 0 || !slices.Contains(delegatedCredentialSchemes(), d.scheme) {
		return nil, errors.New("dtls13: invalid delegated credential encoding or algorithm")
	}
	var err error
	d.public, err = x509.ParsePKIXPublicKey(spki)
	if err != nil || !signatureSchemePublicKeyCompatible(d.public, d.scheme) {
		return nil, errors.New("dtls13: delegated public key does not match its signature scheme")
	}
	return d, nil
}

func delegatedCredentialInput(certificate, signed []byte, server bool) []byte {
	context := "TLS, client delegated credentials"
	if server {
		context = "TLS, server delegated credentials"
	}
	input := bytes.Repeat([]byte{0x20}, 64)
	input = append(input, context...)
	input = append(input, 0)
	input = append(input, certificate...)
	return append(input, signed...)
}

func validateDelegatedCredential(raw []byte, parent *x509.Certificate, now time.Time, server bool, signatures, delegated []tls.SignatureScheme) (*parsedDelegatedCredential, error) {
	d, err := parseDelegatedCredential(raw)
	if err != nil {
		return nil, alertError(alertIllegalParameter, err)
	}
	expiry := parent.NotBefore.Add(time.Duration(d.validTime) * time.Second)
	if now.Before(parent.NotBefore) || now.After(expiry) || expiry.After(now.Add(7*24*time.Hour)) || !expiry.Before(parent.NotAfter) {
		return nil, alertError(alertIllegalParameter, errors.New("dtls13: invalid delegated credential validity period"))
	}
	permitted := false
	for _, extension := range parent.Extensions {
		if extension.Id.Equal(oidDelegationUsage) {
			permitted = !extension.Critical && bytes.Equal(extension.Value, []byte{5, 0})
		}
	}
	if !permitted || parent.KeyUsage&x509.KeyUsageDigitalSignature == 0 {
		return nil, alertError(alertIllegalParameter, errors.New("dtls13: certificate does not permit delegation"))
	}
	if !slices.Contains(signatures, d.algorithm) || !slices.Contains(delegated, d.scheme) {
		return nil, alertError(alertIllegalParameter, errors.New("dtls13: unoffered delegated credential signature scheme"))
	}
	if err = verifySignedMessage(parent.PublicKey, d.algorithm, delegatedCredentialInput(parent.Raw, d.signed, server), d.signature); err != nil {
		return nil, alertError(alertIllegalParameter, fmt.Errorf("dtls13: invalid delegation signature: %w", err))
	}
	return d, nil
}

func (c *Conn) configuredDelegatedCredential(certificate *tls.Certificate) *DelegatedCredential {
	if c == nil || c.config == nil {
		return nil
	}
	for _, d := range c.config.DelegatedCredentials {
		if d != nil && certificate == &d.Certificate {
			return d
		}
	}
	return nil
}

func (c *Conn) validateLocalDelegatedCredential(d *DelegatedCredential, signatures, certificateSignatures, delegated []tls.SignatureScheme, server bool) ([]*x509.Certificate, error) {
	if d == nil {
		return nil, errors.New("dtls13: nil delegated credential")
	}
	if len(certificateSignatures) == 0 {
		certificateSignatures = signatures
	}
	chain, err := configuredCertificateChain(&d.Certificate, certificateSignatures, server)
	if err != nil {
		return nil, err
	}
	parsed, err := validateDelegatedCredential(d.Credential, chain[0], c.config.Time(), server, signatures, delegated)
	if err != nil {
		return nil, err
	}
	signer, ok := d.Certificate.PrivateKey.(crypto.Signer)
	if !ok {
		return nil, errors.New("dtls13: delegated private key is not a signer")
	}
	public, ok := parsed.public.(interface{ Equal(crypto.PublicKey) bool })
	if !ok || !public.Equal(signer.Public()) {
		return nil, errors.New("dtls13: delegated credential and private key do not match")
	}
	return chain, nil
}

func (c *Conn) selectDelegatedCredential(signatures, certificateSignatures, delegated []tls.SignatureScheme, server bool, name string, authorities [][]byte, filters []CertificateOIDFilter) *tls.Certificate {
	if len(delegated) == 0 {
		return nil
	}
	for _, d := range c.config.DelegatedCredentials {
		chain, err := c.validateLocalDelegatedCredential(d, signatures, certificateSignatures, delegated, server)
		if err != nil || (name != "" && chain[0].VerifyHostname(name) != nil) || !certificateSignedBy(chain, authorities) || matchCertificateOIDFilters(chain[0], filters) != nil {
			continue
		}
		return &d.Certificate
	}
	return nil
}

func (c *Conn) certificateSigningScheme(certificate *tls.Certificate, offered []tls.SignatureScheme) (tls.SignatureScheme, error) {
	if delegated := c.configuredDelegatedCredential(certificate); delegated != nil {
		d, err := parseDelegatedCredential(delegated.Credential)
		if err != nil {
			return 0, err
		}
		return d.scheme, nil
	}
	return selectSignatureScheme(certificate.PrivateKey.(crypto.Signer), offered)
}

func (c *Conn) addDelegatedCredential(message *certificateMessage, certificate *tls.Certificate) {
	if d := c.configuredDelegatedCredential(certificate); d != nil && len(message.certificates) != 0 {
		entry := &message.certificates[0]
		if entry.extensions == nil {
			entry.extensions = make(map[uint16][]byte, 1)
		}
		entry.extensions[extDelegatedCredential] = d.Credential
	}
}

func peerDelegatedCredential(message *certificateMessage, certificates []*x509.Certificate, config *Config, server bool, signatures, delegated []tls.SignatureScheme) ([]byte, error) {
	if len(message.certificates) == 0 || !message.certificates[0].hasDelegatedCredential {
		return nil, nil
	}
	if len(certificates) == 0 {
		return nil, alertError(alertIllegalParameter, errors.New("dtls13: delegated credential requires X.509"))
	}
	raw := message.certificates[0].delegatedCredential
	if _, err := validateDelegatedCredential(raw, certificates[0], config.Time(), server, signatures, delegated); err != nil {
		return nil, err
	}
	return bytes.Clone(raw), nil
}

func verifyPeerAuthenticationSignature(certificates []*x509.Certificate, raw, delegated []byte, schemes []tls.SignatureScheme, verify *certificateVerifyMessage, transcript []byte, server bool) error {
	if len(delegated) != 0 {
		d, err := parseDelegatedCredential(delegated)
		if err != nil || d.scheme != verify.algorithm {
			return alertError(alertIllegalParameter, errors.New("dtls13: CertificateVerify does not match delegated credential"))
		}
		if err = verifyCertificateVerify(d.public, verify.algorithm, transcript, verify.signature, server); err != nil {
			return alertError(alertDecryptError, err)
		}
		return nil
	}
	if !slices.Contains(schemes, verify.algorithm) {
		return alertError(alertIllegalParameter, &ProtocolError{"peer selected an unoffered signature scheme"})
	}
	if err := verifyPeerCertificateSignature(certificates, raw, verify.algorithm, transcript, verify.signature, server); err != nil {
		return alertError(alertDecryptError, err)
	}
	return nil
}

func validResumedDelegatedCredential(config *Config, raw []byte, certificates []*x509.Certificate, server bool) bool {
	if len(raw) == 0 {
		return true
	}
	if !config.EnableDelegatedCredentials || len(certificates) == 0 {
		return false
	}
	_, err := validateDelegatedCredential(raw, certificates[0], config.Time(), server, defaultSignatureSchemes(), delegatedCredentialSchemes())
	return err == nil
}

// validateLocalCertificate checks the signing key and both signature offers.
func (c *Conn) validateLocalCertificate(certificate *tls.Certificate, signatures, certificateSignatures, delegated []tls.SignatureScheme, server bool) ([]*x509.Certificate, error) {
	if d := c.configuredDelegatedCredential(certificate); d != nil {
		return c.validateLocalDelegatedCredential(d, signatures, certificateSignatures, delegated, server)
	}
	if len(certificateSignatures) == 0 {
		certificateSignatures = signatures
	}
	chain, err := validateConfiguredCertificateChain(certificate, certificateSignatures, server)
	if err != nil {
		return nil, err
	}
	_, err = selectSignatureScheme(certificate.PrivateKey.(crypto.Signer), signatures)
	return chain, err
}
