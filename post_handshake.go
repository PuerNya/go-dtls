package dtls13

import (
	"crypto/tls"
	"errors"
	"sync"
)

const handshakeTypeCertificateRequest uint8 = 13
const handshakeTypeKeyUpdate uint8 = 24

type certificateRequestMessage struct {
	requestContext              []byte
	signatureSchemes            []tls.SignatureScheme
	certificateSignatureSchemes []tls.SignatureScheme
}

func (m *certificateRequestMessage) marshal() ([]byte, error) {
	return m.marshalWithCertificateCompression(nil)
}

func (m *certificateRequestMessage) marshalWithCertificateCompression(algorithms *certificateCompressionAlgorithms) ([]byte, error) {
	schemes, err := marshalSignatureSchemes(m.signatureSchemes)
	if err != nil {
		return nil, err
	}
	items := map[uint16][]byte{extSignatureAlgorithms: schemes}
	if len(m.certificateSignatureSchemes) > 0 {
		items[extSignatureAlgorithmsCert], err = marshalSignatureSchemes(m.certificateSignatureSchemes)
		if err != nil {
			return nil, err
		}
	}
	if algorithms != nil {
		items[extCompressCertificate], err = marshalCertificateCompressionAlgorithms(algorithms)
		if err != nil {
			return nil, err
		}
	}
	exts, err := marshalExtensions(items, []uint16{extSignatureAlgorithms, extSignatureAlgorithmsCert, extCompressCertificate})
	if err != nil {
		return nil, err
	}
	var w wireBuilder
	w.bytes8(m.requestContext)
	w.b = append(w.b, exts...)
	return w.b, w.err
}
func parseCertificateRequest(b []byte) (*certificateRequestMessage, error) {
	request, _, err := parseCertificateRequestWithCompression(b)
	return request, err
}

func parseCertificateRequestWithCompression(b []byte) (*certificateRequestMessage, *certificateCompressionAlgorithms, error) {
	p := wireParser{b: b}
	m := &certificateRequestMessage{requestContext: append([]byte(nil), p.bytes8()...)}
	var extensionStorage [4]orderedExtension
	exts, err := parseOrderedExtensionsView(p.take(len(p.b)-p.off), extensionStorage[:0])
	if err != nil {
		return nil, nil, err
	}
	raw, ok := orderedExtensionValue(exts, extSignatureAlgorithms)
	if !ok {
		return nil, nil, alertError(alertMissingExtension, &ProtocolError{"CertificateRequest has no signature_algorithms"})
	}
	m.signatureSchemes, err = parseSignatureSchemes(raw)
	if err != nil {
		return nil, nil, err
	}
	if raw, ok = orderedExtensionValue(exts, extSignatureAlgorithmsCert); ok {
		m.certificateSignatureSchemes, err = parseSignatureSchemes(raw)
		if err != nil {
			return nil, nil, err
		}
	}
	var certificateCompressionAlgorithms *certificateCompressionAlgorithms
	if raw, ok = orderedExtensionValue(exts, extCompressCertificate); ok {
		certificateCompressionAlgorithms, err = parseCertificateCompressionAlgorithms(raw)
		if err != nil {
			return nil, nil, err
		}
	}
	for _, extension := range exts {
		if extension.typ != extSignatureAlgorithms && extension.typ != extSignatureAlgorithmsCert && extension.typ != extCompressCertificate && knownExtensionType(extension.typ) {
			return nil, nil, alertError(alertIllegalParameter, &ProtocolError{"recognized extension is not permitted in CertificateRequest"})
		}
	}
	return m, certificateCompressionAlgorithms, nil
}

type keyUpdateMessage struct{ requestUpdate bool }

func (m keyUpdateMessage) marshal() []byte {
	if m.requestUpdate {
		return []byte{1}
	}
	return []byte{0}
}
func parseKeyUpdate(b []byte) (keyUpdateMessage, error) {
	if len(b) != 1 {
		return keyUpdateMessage{}, alertError(alertDecodeError, &ProtocolError{"invalid KeyUpdate length"})
	}
	if b[0] > 1 {
		return keyUpdateMessage{}, &ProtocolError{"invalid KeyUpdate"}
	}
	return keyUpdateMessage{requestUpdate: b[0] == 1}, nil
}

// keyUpdateState enforces RFC 9147 section 8: new-key records and another
// KeyUpdate cannot be sent until the KeyUpdate record has been acknowledged.
type keyUpdateState struct {
	mu      sync.Mutex
	pending bool
	records []recordNumber
}

func (s *keyUpdateState) begin(record recordNumber) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pending {
		return errors.New("dtls13: KeyUpdate is already awaiting acknowledgement")
	}
	s.pending = true
	if cap(s.records) == 1 {
		s.records = append(s.records[:0], record)
	} else {
		s.records = []recordNumber{record}
	}
	return nil
}
func (s *keyUpdateState) addRetransmission(record recordNumber) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.pending {
		return errors.New("dtls13: KeyUpdate is not awaiting acknowledgement")
	}
	s.records = append(s.records, record)
	return nil
}
func (s *keyUpdateState) ack(numbers []recordNumber) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.pending {
		return false
	}
	for _, number := range numbers {
		for _, record := range s.records {
			if number == record {
				s.pending = false
				if cap(s.records) == 1 {
					clear(s.records)
					s.records = s.records[:0]
				} else {
					s.records = nil
				}
				return true
			}
		}
	}
	return false
}
func (s *keyUpdateState) canUseNewKeys() bool { s.mu.Lock(); defer s.mu.Unlock(); return !s.pending }
