package dtls13

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
)

var (
	oidExtensionExtendedKeyUsage = asn1.ObjectIdentifier{2, 5, 29, 37}
	oidAnyExtendedKeyUsage       = asn1.ObjectIdentifier{2, 5, 29, 37, 0}
)

// CertificateOIDFilter is an RFC 9846 certificate-extension filter. Values is
// the DER-encoded extension value, without the X.509 extension wrapper.
type CertificateOIDFilter struct {
	OID    asn1.ObjectIdentifier
	Values []byte
}

// CertificateRequestInfo describes a CertificateRequest received by a client.
// Its slices must not be modified or retained by callbacks.
type CertificateRequestInfo struct {
	AcceptableCAs               [][]byte
	SignatureSchemes            []tls.SignatureScheme
	CertificateSignatureSchemes []tls.SignatureScheme
	OIDFilters                  []CertificateOIDFilter
	Version                     uint16
	Conn                        *Conn
}

// SupportsCertificate reports whether certificate satisfies the request,
// including signature algorithms, acceptable CAs, and recognized OID filters.
func (i *CertificateRequestInfo) SupportsCertificate(certificate *tls.Certificate) error {
	return i.supportsCertificate(certificate, true)
}

func (i *CertificateRequestInfo) supportsCertificate(certificate *tls.Certificate, requireAcceptableCA bool) error {
	if i == nil {
		return errors.New("dtls13: nil CertificateRequestInfo")
	}
	certificateSchemes := i.CertificateSignatureSchemes
	if len(certificateSchemes) == 0 {
		certificateSchemes = i.SignatureSchemes
	}
	parsed, err := validateConfiguredCertificateChain(certificate, certificateSchemes, false)
	if err != nil {
		return err
	}
	signer, ok := certificate.PrivateKey.(crypto.Signer)
	if !ok {
		return errors.New("dtls13: client certificate private key is not a signer")
	}
	if _, err = selectSignatureScheme(signer, i.SignatureSchemes); err != nil {
		return err
	}
	if requireAcceptableCA && !certificateSignedBy(parsed, i.AcceptableCAs) {
		return errors.New("dtls13: certificate chain is not signed by an acceptable CA")
	}
	return matchCertificateOIDFilters(parsed[0], i.OIDFilters)
}

// SupportsCertificate reports whether certificate is compatible with the
// ClientHello, including SNI, signature algorithms, and acceptable CA hints.
func (i *ClientHelloInfo) SupportsCertificate(certificate *tls.Certificate) error {
	if i == nil {
		return errors.New("dtls13: nil ClientHelloInfo")
	}
	certificateSchemes := i.CertificateSignatureSchemes
	if len(certificateSchemes) == 0 {
		certificateSchemes = i.SignatureSchemes
	}
	parsed, err := validateConfiguredCertificateChain(certificate, certificateSchemes, true)
	if err != nil {
		return err
	}
	if i.ServerName != "" {
		if err = parsed[0].VerifyHostname(i.ServerName); err != nil {
			return fmt.Errorf("dtls13: certificate is not valid for %q: %w", i.ServerName, err)
		}
	}
	signer, ok := certificate.PrivateKey.(crypto.Signer)
	if !ok {
		return errors.New("dtls13: server certificate private key is not a signer")
	}
	if _, err = selectSignatureScheme(signer, i.SignatureSchemes); err != nil {
		return err
	}
	if !certificateSignedBy(parsed, i.AcceptableCAs) {
		return errors.New("dtls13: certificate chain is not signed by an acceptable CA")
	}
	return nil
}

func certificateSignedBy(chain []*x509.Certificate, acceptableCAs [][]byte) bool {
	if len(acceptableCAs) == 0 {
		return true
	}
	for _, certificate := range chain {
		for _, ca := range acceptableCAs {
			if bytes.Equal(certificate.RawIssuer, ca) {
				return true
			}
		}
	}
	return false
}

func marshalCertificateAuthorities(authorities [][]byte) ([]byte, error) {
	if err := validateCertificateAuthorities(authorities); err != nil {
		return nil, err
	}
	return marshalCertificateAuthoritiesUnchecked(authorities)
}

func marshalCertificateAuthoritiesUnchecked(authorities [][]byte) ([]byte, error) {
	if len(authorities) == 0 {
		return nil, nil
	}
	w := newWireBuilder(2)
	start := w.startVector16()
	for _, authority := range authorities {
		w.bytes16(authority)
	}
	w.endVector16(start)
	return w.b, w.err
}

func validateCertificateAuthorities(authorities [][]byte) error {
	for _, authority := range authorities {
		if err := validateDistinguishedName(authority); err != nil {
			return err
		}
	}
	return nil
}

func parseCertificateAuthorities(data []byte) ([][]byte, error) {
	var authorities [][]byte
	err := forEachCertificateAuthority(data, func(authority []byte) {
		authorities = append(authorities, append([]byte(nil), authority...))
	})
	return authorities, err
}

func forEachCertificateAuthority(data []byte, yield func([]byte)) error {
	p := wireParser{b: data}
	raw := p.bytes16()
	if err := p.done(); err != nil {
		return err
	}
	if len(raw) < 3 {
		return alertError(alertDecodeError, &ProtocolError{"invalid certificate_authorities extension length"})
	}
	q := wireParser{b: raw}
	for q.off < len(q.b) {
		authority := q.bytes16()
		if q.err != nil {
			return q.err
		}
		if err := validateDistinguishedName(authority); err != nil {
			return alertError(alertDecodeError, err)
		}
		if yield != nil {
			yield(authority)
		}
	}
	return nil
}

func validateDistinguishedName(data []byte) error {
	if len(data) == 0 {
		return &ProtocolError{"empty certificate authority distinguished name"}
	}
	var name pkix.RDNSequence
	rest, err := asn1.Unmarshal(data, &name)
	if err != nil || len(rest) != 0 {
		return &ProtocolError{"malformed certificate authority distinguished name"}
	}
	return nil
}

func marshalOIDFilters(filters []CertificateOIDFilter) ([]byte, error) {
	if len(filters) == 0 {
		return nil, nil
	}
	w := newWireBuilder(2)
	start := w.startVector16()
	for _, filter := range filters {
		oid, err := asn1.Marshal(filter.OID)
		if err != nil || len(oid) == 0 || len(oid) > 255 {
			return nil, &ProtocolError{"invalid oid_filters OID"}
		}
		w.bytes8(oid)
		w.bytes16(filter.Values)
	}
	w.endVector16(start)
	return w.b, w.err
}

func validateOIDFilters(filters []CertificateOIDFilter) error {
	for i, filter := range filters {
		oid, err := asn1.Marshal(filter.OID)
		if err != nil || len(oid) == 0 || len(oid) > 255 {
			return &ProtocolError{"invalid oid_filters OID"}
		}
		for j := range i {
			if filters[j].OID.Equal(filter.OID) {
				return &ProtocolError{"duplicate oid_filters OID"}
			}
		}
		if err = validateOIDFilterValues(filter); err != nil {
			return err
		}
	}
	return nil
}

func parseOIDFilters(data []byte) ([]CertificateOIDFilter, error) {
	p := wireParser{b: data}
	raw := p.bytes16()
	if err := p.done(); err != nil {
		return nil, err
	}
	q := wireParser{b: raw}
	filters := make([]CertificateOIDFilter, 0, 2)
	for q.off < len(q.b) {
		oidDER := q.bytes8()
		values := q.bytes16()
		if q.err != nil || len(oidDER) == 0 {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid oid_filters extension length"})
		}
		var oid asn1.ObjectIdentifier
		rest, err := asn1.Unmarshal(oidDER, &oid)
		canonical, marshalErr := asn1.Marshal(oid)
		if err != nil || marshalErr != nil || len(rest) != 0 || !bytes.Equal(canonical, oidDER) {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid oid_filters OID"})
		}
		for _, existing := range filters {
			if existing.OID.Equal(oid) {
				return nil, alertError(alertIllegalParameter, &ProtocolError{"duplicate oid_filters OID"})
			}
		}
		filter := CertificateOIDFilter{OID: oid, Values: append([]byte(nil), values...)}
		if err = validateOIDFilterValues(filter); err != nil {
			description := uint8(alertDecodeError)
			if oid.Equal(oidExtensionExtendedKeyUsage) && containsAnyExtendedKeyUsage(values) {
				description = alertIllegalParameter
			}
			return nil, alertError(description, err)
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

func validateOIDFilterValues(filter CertificateOIDFilter) error {
	if len(filter.Values) == 0 {
		return nil
	}
	switch {
	case filter.OID.Equal(oidExtensionKeyUsage):
		var usage asn1.BitString
		rest, err := asn1.Unmarshal(filter.Values, &usage)
		if err != nil || len(rest) != 0 {
			return &ProtocolError{"malformed Key Usage oid_filter"}
		}
	case filter.OID.Equal(oidExtensionExtendedKeyUsage):
		var usages []asn1.ObjectIdentifier
		rest, err := asn1.Unmarshal(filter.Values, &usages)
		if err != nil || len(rest) != 0 {
			return &ProtocolError{"malformed Extended Key Usage oid_filter"}
		}
		for _, usage := range usages {
			if usage.Equal(oidAnyExtendedKeyUsage) {
				return &ProtocolError{"anyExtendedKeyUsage is forbidden in oid_filters"}
			}
		}
	}
	return nil
}

func containsAnyExtendedKeyUsage(values []byte) bool {
	var usages []asn1.ObjectIdentifier
	rest, err := asn1.Unmarshal(values, &usages)
	if err != nil || len(rest) != 0 {
		return false
	}
	for _, usage := range usages {
		if usage.Equal(oidAnyExtendedKeyUsage) {
			return true
		}
	}
	return false
}

func matchCertificateOIDFilters(certificate *x509.Certificate, filters []CertificateOIDFilter) error {
	for _, filter := range filters {
		if !filter.OID.Equal(oidExtensionKeyUsage) && !filter.OID.Equal(oidExtensionExtendedKeyUsage) {
			continue
		}
		var certificateValue []byte
		for _, extension := range certificate.Extensions {
			if extension.Id.Equal(filter.OID) {
				certificateValue = extension.Value
				break
			}
		}
		if certificateValue == nil {
			return fmt.Errorf("dtls13: certificate does not contain requested extension %s", filter.OID)
		}
		if len(filter.Values) == 0 {
			continue
		}
		if filter.OID.Equal(oidExtensionKeyUsage) {
			var requested, present asn1.BitString
			requestRest, requestErr := asn1.Unmarshal(filter.Values, &requested)
			presentRest, presentErr := asn1.Unmarshal(certificateValue, &present)
			if requestErr != nil || presentErr != nil || len(requestRest) != 0 || len(presentRest) != 0 {
				return errors.New("dtls13: malformed Key Usage extension")
			}
			for bit := 0; bit < requested.BitLength; bit++ {
				if requested.At(bit) != 0 && present.At(bit) == 0 {
					return errors.New("dtls13: certificate does not satisfy requested Key Usage")
				}
			}
			continue
		}
		var requested, present []asn1.ObjectIdentifier
		requestRest, requestErr := asn1.Unmarshal(filter.Values, &requested)
		presentRest, presentErr := asn1.Unmarshal(certificateValue, &present)
		if requestErr != nil || presentErr != nil || len(requestRest) != 0 || len(presentRest) != 0 {
			return errors.New("dtls13: malformed Extended Key Usage extension")
		}
		for _, want := range requested {
			found := false
			for _, have := range present {
				found = found || have.Equal(want)
			}
			if !found {
				return fmt.Errorf("dtls13: certificate does not contain requested Extended Key Usage %s", want)
			}
		}
	}
	return nil
}

func cloneByteSlices(values [][]byte) [][]byte {
	cloned := make([][]byte, len(values))
	for i, value := range values {
		cloned[i] = append([]byte(nil), value...)
	}
	return cloned
}

func cloneOIDFilters(filters []CertificateOIDFilter) []CertificateOIDFilter {
	cloned := make([]CertificateOIDFilter, len(filters))
	for i, filter := range filters {
		cloned[i] = CertificateOIDFilter{
			OID:    append(asn1.ObjectIdentifier(nil), filter.OID...),
			Values: append([]byte(nil), filter.Values...),
		}
	}
	return cloned
}

func (h *clientHello) setCertificateAuthorities(authorities [][]byte) error {
	if len(authorities) == 0 {
		if h != nil && h.unknownExtensions != nil {
			delete(h.unknownExtensions, extCertificateAuthorities)
		}
		return nil
	}
	raw, err := marshalCertificateAuthoritiesUnchecked(authorities)
	if err != nil {
		return err
	}
	if h.unknownExtensions == nil {
		h.unknownExtensions = make(map[uint16][]byte)
	}
	h.unknownExtensions[extCertificateAuthorities] = raw
	return nil
}

func (h *clientHello) certificateAuthorityNames() [][]byte {
	if h == nil || h.unknownExtensions == nil {
		return nil
	}
	var authorities [][]byte
	_ = forEachCertificateAuthority(h.unknownExtensions[extCertificateAuthorities], func(authority []byte) {
		authorities = append(authorities, authority)
	})
	return authorities
}

func (c *Conn) clientHelloInfo(hello *clientHello) *ClientHelloInfo {
	return &ClientHelloInfo{
		ServerName:                  hello.serverName,
		SupportedProtos:             hello.alpn,
		AcceptableCAs:               hello.certificateAuthorityNames(),
		SignatureSchemes:            hello.signatureSchemes,
		CertificateSignatureSchemes: hello.certificateSignatureSchemes,
		Version:                     VersionDTLS13,
		Conn:                        c,
	}
}

func (c *Conn) certificateRequestInfo(request *certificateRequestMessage) CertificateRequestInfo {
	return CertificateRequestInfo{
		AcceptableCAs:               request.certificateAuthorities,
		SignatureSchemes:            request.signatureSchemes,
		CertificateSignatureSchemes: request.certificateSignatureSchemes,
		OIDFilters:                  request.oidFilters,
		Version:                     VersionDTLS13,
		Conn:                        c,
	}
}

func (c *Conn) newCertificateRequest(context []byte) *certificateRequestMessage {
	request := &certificateRequestMessage{
		requestContext:   context,
		signatureSchemes: defaultSignatureSchemes(),
		oidFilters:       c.config.ClientCertificateOIDFilters,
	}
	if c.config.ClientCAs != nil {
		//nolint:staticcheck // CertPool has no replacement that exposes configured subjects.
		request.certificateAuthorities = c.config.ClientCAs.Subjects()
	}
	return request
}

func (c *Conn) selectClientCertificate(request *certificateRequestMessage) (*tls.Certificate, error) {
	if c.config.GetClientCertificate != nil {
		info := c.certificateRequestInfo(request)
		certificate, err := c.config.GetClientCertificate(&info)
		if err != nil {
			return nil, err
		}
		if certificate == nil {
			return nil, errors.New("dtls13: GetClientCertificate returned nil")
		}
		if len(certificate.Certificate) == 0 {
			return certificate, nil
		}
		if err = info.supportsCertificate(certificate, false); err != nil {
			return nil, err
		}
		return certificate, nil
	}
	info := c.certificateRequestInfo(request)
	for i := range c.config.Certificates {
		certificate := &c.config.Certificates[i]
		if info.SupportsCertificate(certificate) == nil {
			return certificate, nil
		}
	}
	return nil, nil
}
