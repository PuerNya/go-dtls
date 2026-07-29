package dtls13

import (
	"crypto/tls"
	"encoding/hex"
)

var helloRetryRequestRandom = func() [32]byte {
	raw, _ := hex.DecodeString("cf21ad74e59a6111be1d8c021e65b891c2a211167abb8c5e079e09e2c8a8339c")
	var out [32]byte
	copy(out[:], raw)
	return out
}()

type helloRetryRequest struct {
	sessionID     []byte
	cipherSuite   uint16
	selectedGroup tls.CurveID
	cookie        []byte
}

func (h *helloRetryRequest) marshal() ([]byte, error) {
	if len(h.cookie) == 0 && h.selectedGroup == 0 {
		return nil, &ProtocolError{"HelloRetryRequest does not request a change"}
	}
	if len(h.cookie) > 65535 {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	if len(h.sessionID) > 255 {
		return nil, &ProtocolError{"8-bit vector overflow"}
	}
	// Compute the extension vector length without allocating temporary values.
	extsLength := 2
	addExtension := func(valueLength int) bool {
		if valueLength > 65535 || extsLength > 65535-4-valueLength {
			return false
		}
		extsLength += 4 + valueLength
		return true
	}
	if !addExtension(2) || (len(h.cookie) > 0 && !addExtension(2+len(h.cookie))) ||
		(h.selectedGroup != 0 && !addExtension(2)) {
		return nil, &ProtocolError{"16-bit vector overflow"}
	}
	w := newWireBuilder(2 + len(helloRetryRequestRandom) + 1 + len(h.sessionID) + 2 + 1 + extsLength)
	w.u16(int(dtlsLegacyVersion))
	w.b = append(w.b, helloRetryRequestRandom[:]...)
	w.bytes8(h.sessionID)
	w.u16(int(h.cipherSuite))
	w.u8(0)
	start := w.startVector16()
	w.u16(int(extSupportedVersions))
	w.bytes16(serverSupportedVersionStorage[:])
	if len(h.cookie) > 0 {
		w.u16(int(extCookie))
		w.u16(2 + len(h.cookie))
		w.bytes16(h.cookie)
	}
	if h.selectedGroup != 0 {
		w.u16(int(extKeyShare))
		w.u16(2)
		w.u16(int(h.selectedGroup))
	}
	w.endVector16(start)
	return w.b, w.err
}
func parseHelloRetryRequest(b []byte) (*helloRetryRequest, error) {
	p := wireParser{b: b}
	if p.u16() != int(dtlsLegacyVersion) {
		return nil, &ProtocolError{"invalid HelloRetryRequest legacy version"}
	}
	random := p.take(32)
	if random == nil || string(random) != string(helloRetryRequestRandom[:]) {
		return nil, &ProtocolError{"ServerHello is not a HelloRetryRequest"}
	}
	h := &helloRetryRequest{sessionID: append([]byte(nil), p.bytes8()...)}
	if len(h.sessionID) != 0 {
		return nil, alertError(alertIllegalParameter, &ProtocolError{"DTLS 1.3 HelloRetryRequest legacy_session_id must be empty"})
	}
	h.cipherSuite = uint16(p.u16())
	if p.u8() != 0 {
		return nil, &ProtocolError{"invalid HelloRetryRequest compression method"}
	}
	var extensionStorage [3]orderedExtension
	exts, err := parseOrderedExtensionsView(p.take(len(p.b)-p.off), extensionStorage[:0])
	if err != nil {
		return nil, err
	}
	v, ok := orderedExtensionValue(exts, extSupportedVersions)
	if !ok {
		return nil, alertError(alertProtocolVersion, &ProtocolError{"HelloRetryRequest has no supported_versions extension"})
	}
	if len(v) != 2 {
		return nil, alertError(alertDecodeError, &ProtocolError{"invalid HelloRetryRequest supported_versions length"})
	}
	if uint16(v[0])<<8|uint16(v[1]) != VersionDTLS13 {
		return nil, alertError(alertIllegalParameter, &ProtocolError{"HelloRetryRequest did not select DTLS 1.3"})
	}
	raw, ok := orderedExtensionValue(exts, extKeyShare)
	if ok {
		if len(raw) != 2 {
			return nil, alertError(alertDecodeError, &ProtocolError{"invalid HelloRetryRequest selected group length"})
		}
		h.selectedGroup = tls.CurveID(uint16(raw[0])<<8 | uint16(raw[1]))
	}
	if raw, ok = orderedExtensionValue(exts, extCookie); ok {
		h.cookie, err = parseCookie(raw)
		if err != nil {
			return nil, err
		}
	}
	for _, extension := range exts {
		if extension.typ != extSupportedVersions && extension.typ != extCookie && extension.typ != extKeyShare {
			if knownExtensionType(extension.typ) {
				return nil, alertError(alertIllegalParameter, &ProtocolError{"recognized extension is not permitted in HelloRetryRequest"})
			}
			return nil, alertError(alertUnsupportedExtension, &ProtocolError{"unsupported HelloRetryRequest extension"})
		}
	}
	if len(h.cookie) == 0 && h.selectedGroup == 0 {
		return nil, &ProtocolError{"HelloRetryRequest does not request a change"}
	}
	return h, nil
}
