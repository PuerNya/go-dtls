package dtls13

import (
	"crypto/tls"
	"encoding/binary"
	"slices"
)

const (
	handshakeTypeEncryptedExtensions uint8 = 8
	handshakeTypeCertificate         uint8 = 11
	handshakeTypeCertificateVerify   uint8 = 15
	handshakeTypeFinished            uint8 = 20
)

type encryptedExtensions struct {
	extensions         map[uint16][]byte
	recordSizeLimit    uint16
	hasRecordSizeLimit bool
	parsedStorage      [8]orderedExtension
	parsedOverflow     []orderedExtension
	parsedCount        int
}

func (m *encryptedExtensions) marshal() ([]byte, error) {
	if !m.hasRecordSizeLimit {
		var storage [8]uint16
		order := sortedExtensionTypesInto(storage[:0], m.extensions)
		return marshalExtensions(m.extensions, order)
	}
	if m.recordSizeLimit < minRecordSizeLimit || m.recordSizeLimit > defaultRecordSizeLimit {
		return nil, &ProtocolError{"invalid record_size_limit"}
	}
	if _, duplicate := m.extensions[extRecordSizeLimit]; duplicate {
		return nil, &ProtocolError{"duplicate record_size_limit extension"}
	}
	var limit [2]byte
	binary.BigEndian.PutUint16(limit[:], m.recordSizeLimit)
	var typeStorage [8]uint16
	order := sortedExtensionTypesInto(typeStorage[:0], m.extensions)
	var extensionStorage [9]orderedExtension
	extensions := extensionStorage[:0]
	for _, typ := range order {
		extensions = append(extensions, orderedExtension{typ: typ, value: m.extensions[typ]})
	}
	extensions = append(extensions, orderedExtension{typ: extRecordSizeLimit, value: limit[:]})
	length, err := orderedExtensionsWireLength(extensions)
	if err != nil {
		return nil, err
	}
	w := newWireBuilder(length)
	appendOrderedExtensions(&w, extensions)
	return w.b, w.err
}
func parseEncryptedExtensions(b []byte) (encryptedExtensions, error) {
	var m encryptedExtensions
	exts, err := parseOrderedExtensionsView(b, m.parsedStorage[:0])
	if err != nil {
		return encryptedExtensions{}, err
	}
	if len(exts) > len(m.parsedStorage) {
		m.parsedOverflow = exts
	} else {
		m.parsedCount = len(exts)
	}
	return m, nil
}

func validateEncryptedExtension(hello *clientHello, typ uint16, raw []byte) (protocol string, earlyData bool, err error) {
	switch typ {
	case extServerName:
		if hello.serverName == "" || len(raw) != 0 {
			return "", false, alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited server_name acknowledgement"})
		}
	case extALPN:
		if len(hello.alpn) == 0 {
			return "", false, alertError(alertUnsupportedExtension, &ProtocolError{"server selected unoffered ALPN"})
		}
		protocols, parseErr := parseALPN(raw)
		if parseErr != nil || len(protocols) != 1 {
			return "", false, &ProtocolError{"invalid server ALPN selection"}
		}
		found := false
		for _, offered := range hello.alpn {
			if offered == protocols[0] {
				found = true
				break
			}
		}
		if !found {
			return "", false, alertError(alertNoApplicationProtocol, &ProtocolError{"server selected an unoffered ALPN protocol"})
		}
		protocol = protocols[0]
	case extEarlyData:
		if !hello.earlyData || len(raw) != 0 {
			return "", false, alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited early_data acceptance"})
		}
		earlyData = true
	case extSupportedGroups:
		if len(hello.supportedGroups) == 0 {
			return "", false, alertError(alertUnsupportedExtension, &ProtocolError{"server sent unoffered supported_groups"})
		}
		if _, parseErr := parseSupportedGroups(raw); parseErr != nil {
			return "", false, parseErr
		}
	case extECH:
		if len(hello.encryptedClientHello()) == 0 {
			return "", false, alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited ECH retry configurations"})
		}
		if _, parseErr := parseECHConfigList(raw); parseErr != nil {
			return "", false, alertError(alertDecodeError, parseErr)
		}
	default:
		if knownExtensionType(typ) {
			return "", false, alertError(alertIllegalParameter, &ProtocolError{"recognized extension is not permitted in EncryptedExtensions"})
		}
		return "", false, alertError(alertUnsupportedExtension, &ProtocolError{"unsupported EncryptedExtensions extension"})
	}
	return protocol, earlyData, nil
}

func validateEncryptedExtensions(hello *clientHello, message *encryptedExtensions) (protocol string, earlyData bool, retryConfigs []byte, err error) {
	if hello == nil || message == nil {
		return "", false, nil, &ProtocolError{"missing EncryptedExtensions context"}
	}
	hasRecordSizeLimit, hasMaxFragmentLength := false, false
	if message.extensions != nil {
		_, hasRecordSizeLimit = message.extensions[extRecordSizeLimit]
		_, hasMaxFragmentLength = message.extensions[extMaxFragmentLength]
	} else {
		parsed := message.parsedOverflow
		if parsed == nil {
			parsed = message.parsedStorage[:message.parsedCount]
		}
		for _, extension := range parsed {
			hasRecordSizeLimit = hasRecordSizeLimit || extension.typ == extRecordSizeLimit
			hasMaxFragmentLength = hasMaxFragmentLength || extension.typ == extMaxFragmentLength
		}
	}
	if hasRecordSizeLimit && hasMaxFragmentLength {
		return "", false, nil, alertError(alertIllegalParameter, &ProtocolError{"record_size_limit and max_fragment_length cannot both be negotiated"})
	}
	validate := func(typ uint16, raw []byte) error {
		if typ == extECH {
			if _, _, validateErr := validateEncryptedExtension(hello, typ, raw); validateErr != nil {
				return validateErr
			}
			retryConfigs = raw
			return nil
		}
		if typ == extRecordSizeLimit {
			if !hello.hasRecordSizeLimit {
				return alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited record_size_limit"})
			}
			limit, parseErr := parseRecordSizeLimit(raw, true)
			if parseErr != nil {
				return parseErr
			}
			message.recordSizeLimit = limit
			message.hasRecordSizeLimit = true
			return nil
		}
		selected, accepted, validateErr := validateEncryptedExtension(hello, typ, raw)
		if selected != "" {
			protocol = selected
		}
		earlyData = earlyData || accepted
		return validateErr
	}
	if message.extensions != nil {
		for typ, raw := range message.extensions {
			if err = validate(typ, raw); err != nil {
				return "", false, nil, err
			}
		}
		return protocol, earlyData, retryConfigs, nil
	}
	parsed := message.parsedOverflow
	if parsed == nil {
		parsed = message.parsedStorage[:message.parsedCount]
	}
	for _, extension := range parsed {
		if err = validate(extension.typ, extension.value); err != nil {
			return "", false, nil, err
		}
	}
	return protocol, earlyData, retryConfigs, nil
}

func validateEarlyDataSelection(accepted bool, selectedIdentity *uint16) error {
	if accepted && (selectedIdentity == nil || *selectedIdentity != 0) {
		return alertError(alertIllegalParameter, &ProtocolError{"early_data was accepted without selecting PSK identity 0"})
	}
	return nil
}

type certificateEntry struct {
	data                 []byte
	extensions           map[uint16][]byte
	peerExtensionType    uint16
	peerExtensionPresent bool
}
type certificateMessage struct {
	requestContext []byte
	certificates   []certificateEntry
}

func validateCertificateMessage(message *certificateMessage, expectedContext []byte) error {
	if message == nil || !equalBytes(message.requestContext, expectedContext) {
		return alertError(alertIllegalParameter, &ProtocolError{"Certificate request context mismatch"})
	}
	for _, certificate := range message.certificates {
		if certificate.peerExtensionPresent {
			if knownExtensionType(certificate.peerExtensionType) {
				return alertError(alertIllegalParameter, &ProtocolError{"recognized extension is not permitted in CertificateEntry"})
			}
			return alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited CertificateEntry extension"})
		}
		for typ := range certificate.extensions {
			if knownExtensionType(typ) {
				return alertError(alertIllegalParameter, &ProtocolError{"recognized extension is not permitted in CertificateEntry"})
			}
			return alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited CertificateEntry extension"})
		}
	}
	return nil
}

func (m *certificateMessage) marshal() ([]byte, error) {
	if len(m.requestContext) > 255 {
		return nil, &ProtocolError{"8-bit vector overflow"}
	}
	listLength := 0
	for _, cert := range m.certificates {
		if len(cert.data) == 0 {
			return nil, &ProtocolError{"empty certificate entry"}
		}
		extensionsLength, err := allExtensionsWireLength(cert.extensions)
		if err != nil {
			return nil, err
		}
		entryLength := 3 + len(cert.data) + extensionsLength
		if len(cert.data) >= 1<<24 || entryLength >= 1<<24 || listLength > (1<<24)-1-entryLength {
			return nil, &ProtocolError{"24-bit vector overflow"}
		}
		listLength += 3 + len(cert.data) + extensionsLength
	}
	w := newWireBuilder(1 + len(m.requestContext) + 3 + listLength)
	w.bytes8(m.requestContext)
	start := w.startVector24()
	var extensionTypeStorage [8]uint16
	for _, cert := range m.certificates {
		w.bytes24(cert.data)
		order := sortedExtensionTypesInto(extensionTypeStorage[:0], cert.extensions)
		appendExtensions(&w, cert.extensions, order)
	}
	w.endVector24(start)
	return w.b, w.err
}
func parseCertificateMessage(b []byte, maxSize int) (*certificateMessage, error) {
	if len(b) > maxSize {
		return nil, &ProtocolError{"Certificate message exceeds configured limit"}
	}
	p := wireParser{b: b}
	m := &certificateMessage{requestContext: append([]byte(nil), p.bytes8()...)}
	raw := p.bytes24()
	if err := p.done(); err != nil {
		return nil, err
	}
	q := wireParser{b: raw}
	var extensionStorage [8]orderedExtension
	for q.off < len(q.b) {
		data := q.bytes24()
		if q.err != nil {
			return nil, q.err
		}
		if len(data) == 0 {
			return nil, &ProtocolError{"empty certificate entry"}
		}
		start := q.off
		extLen := q.u16()
		if q.err != nil {
			return nil, q.err
		}
		q.off = start
		extWire := q.take(2 + extLen)
		exts, err := parseOrderedExtensionsView(extWire, extensionStorage[:0])
		if err != nil {
			return nil, err
		}
		entry := certificateEntry{data: append([]byte(nil), data...)}
		if len(exts) > 0 {
			entry.peerExtensionType = exts[0].typ
			entry.peerExtensionPresent = true
		}
		m.certificates = append(m.certificates, entry)
	}
	return m, q.done()
}

type certificateVerifyMessage struct {
	algorithm tls.SignatureScheme
	signature []byte
}

func (m *certificateVerifyMessage) marshal() ([]byte, error) {
	if len(m.signature) == 0 {
		return nil, &ProtocolError{"empty CertificateVerify signature"}
	}
	if len(m.signature) > 65535 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	w := newWireBuilder(4 + len(m.signature))
	w.u16(int(m.algorithm))
	w.bytes16(m.signature)
	return w.b, w.err
}
func parseCertificateVerify(b []byte) (*certificateVerifyMessage, error) {
	p := wireParser{b: b}
	m := &certificateVerifyMessage{algorithm: tls.SignatureScheme(p.u16()), signature: append([]byte(nil), p.bytes16()...)}
	if err := p.done(); err != nil {
		return nil, err
	}
	return m, nil
}

func parseFinished(b []byte, hashSize int) ([]byte, error) {
	if len(b) != hashSize {
		return nil, alertError(alertDecodeError, &ProtocolError{"invalid Finished length"})
	}
	return append([]byte(nil), b...), nil
}

func sortedExtensionTypesInto(types []uint16, exts map[uint16][]byte) []uint16 {
	if cap(types) < len(exts) {
		types = make([]uint16, 0, len(exts))
	} else {
		types = types[:0]
	}
	for typ := range exts {
		types = append(types, typ)
	}
	slices.Sort(types)
	return types
}
