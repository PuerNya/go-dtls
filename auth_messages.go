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
	extensions               map[uint16][]byte
	recordSizeLimit          uint16
	hasRecordSizeLimit       bool
	parsedStorage            [8]orderedExtension
	parsedOverflow           []orderedExtension
	parsedCount              int
	expectedTicketCount      uint8
	hasTicketRequest         bool
	serverCertificateType    CertificateType
	clientCertificateType    CertificateType
	hasClientCertificateType bool
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
		return nil, alertError(alertIllegalParameter, &ProtocolError{"duplicate record_size_limit extension"})
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
		found := slices.Contains(hello.alpn, protocols[0])
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
		if typ == extServerCertificateType || typ == extClientCertificateType {
			selected, err := validateSelectedCertificateType(hello, typ, raw)
			if err != nil {
				return err
			}
			if typ == extServerCertificateType {
				message.serverCertificateType = selected
			} else {
				message.clientCertificateType = selected
				message.hasClientCertificateType = true
			}
			return nil
		}
		if typ == extTicketRequest {
			if !hello.ticketRequest.Enabled {
				return alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited ticket_request response"})
			}
			if len(raw) != 1 {
				return alertError(alertDecodeError, &ProtocolError{"invalid ticket_request response length"})
			}
			message.expectedTicketCount = raw[0]
			message.hasTicketRequest = true
			return nil
		}
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

// certificateEntryExtensionVerdict is the CertificateEntry extension
// classification computed once while parsing, so the entry never needs a
// per-entry extension map on the receive path.
type certificateEntryExtensionVerdict uint8

const (
	// certificateEntryExtensionsRecognized means a recognized extension type
	// appeared in CertificateEntry, where RFC 9846 does not permit it. Zero
	// is reserved for entries without extensions.
	certificateEntryExtensionsRecognized certificateEntryExtensionVerdict = iota + 1
	// certificateEntryExtensionsUnsolicited means an extension appeared
	// without a corresponding request from this endpoint.
	certificateEntryExtensionsUnsolicited
)

type certificateEntry struct {
	data          []byte
	extensions    map[uint16][]byte
	ocspResponse  []byte
	statusRequest bool
	// peerVerdict is populated by parseCertificateMessage; entries built
	// locally for sending leave it zero and use extensions instead.
	peerVerdict certificateEntryExtensionVerdict
}
type certificateMessage struct {
	requestContext []byte
	certificates   []certificateEntry
}

// classifyCertificateEntryExtension maps one CertificateEntry extension type
// to its verdict. Recognized-but-misplaced wins over unsolicited so the RFC
// 9846 §4.3 illegal_parameter rule is reported even when both kinds appear.
// status_request is handled separately because it is permitted on the first
// entry only after the peer requested it in ClientHello or CertificateRequest.
func classifyCertificateEntryExtension(typ uint16) certificateEntryExtensionVerdict {
	if knownExtensionType(typ) {
		return certificateEntryExtensionsRecognized
	}
	return certificateEntryExtensionsUnsolicited
}

func mergeCertificateEntryVerdict(current, next certificateEntryExtensionVerdict) certificateEntryExtensionVerdict {
	if next == certificateEntryExtensionsRecognized || current == certificateEntryExtensionsRecognized {
		return certificateEntryExtensionsRecognized
	}
	if next == certificateEntryExtensionsUnsolicited {
		return certificateEntryExtensionsUnsolicited
	}
	return current
}

func certificateEntryVerdictError(verdict certificateEntryExtensionVerdict) error {
	switch verdict {
	case certificateEntryExtensionsRecognized:
		return alertError(alertIllegalParameter, &ProtocolError{"recognized extension is not permitted in CertificateEntry"})
	case certificateEntryExtensionsUnsolicited:
		return alertError(alertUnsupportedExtension, &ProtocolError{"unsolicited CertificateEntry extension"})
	default:
		return nil
	}
}

func validateCertificateMessage(message *certificateMessage, expectedContext []byte) error {
	return validateCertificateMessageWithStatusRequest(message, expectedContext, false)
}

func validateCertificateMessageWithStatusRequest(message *certificateMessage, expectedContext []byte, requested bool) error {
	if message == nil || !equalBytes(message.requestContext, expectedContext) {
		return alertError(alertIllegalParameter, &ProtocolError{"Certificate request context mismatch"})
	}
	for index, certificate := range message.certificates {
		verdict := certificate.peerVerdict
		if certificate.statusRequest {
			if !requested {
				verdict = mergeCertificateEntryVerdict(verdict, certificateEntryExtensionsUnsolicited)
			} else if index != 0 {
				verdict = mergeCertificateEntryVerdict(verdict, certificateEntryExtensionsRecognized)
			}
		}
		for typ, raw := range certificate.extensions {
			if typ == extStatusRequest {
				if _, err := parseOCSPResponse(raw); err != nil {
					return err
				}
				if !requested {
					verdict = mergeCertificateEntryVerdict(verdict, certificateEntryExtensionsUnsolicited)
				} else if index != 0 {
					verdict = mergeCertificateEntryVerdict(verdict, certificateEntryExtensionsRecognized)
				}
				continue
			}
			verdict = mergeCertificateEntryVerdict(verdict, classifyCertificateEntryExtension(typ))
		}
		if err := certificateEntryVerdictError(verdict); err != nil {
			return err
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
			return nil, alertError(alertDecodeError, &ProtocolError{"empty certificate entry"})
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
		for i := range exts {
			extension := exts[i]
			if extension.typ == extStatusRequest {
				entry.ocspResponse, err = parseOCSPResponse(extension.value)
				if err != nil {
					return nil, err
				}
				entry.statusRequest = true
				continue
			}
			entry.peerVerdict = mergeCertificateEntryVerdict(entry.peerVerdict, classifyCertificateEntryExtension(extension.typ))
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
