package dtls13

import "bytes"

// Private-use extension, version 1: one version byte followed by opaque data.
// Both endpoints must agree on its meaning; it is not an IANA-assigned protocol.
const extHandshakeMetadata uint16 = 0xff02

const maxHandshakeMetadata = 4096

func marshalHandshakeMetadata(data []byte) ([]byte, error) {
	if len(data) > maxHandshakeMetadata {
		return nil, &ConfigError{"handshake metadata exceeds 4096 bytes"}
	}
	return append([]byte{1}, data...), nil
}

func parseHandshakeMetadata(raw []byte) ([]byte, error) {
	if len(raw) == 0 || len(raw) > maxHandshakeMetadata+1 {
		return nil, alertError(alertDecodeError, &ProtocolError{"invalid handshake metadata length"})
	}
	if raw[0] != 1 {
		return nil, alertError(alertIllegalParameter, &ProtocolError{"unsupported handshake metadata version"})
	}
	return bytes.Clone(raw[1:]), nil
}

func (c *Conn) serverHandshakeMetadata(s *serverHandshakeState, ee *encryptedExtensions) error {
	if s.ch.handshakeMetadata == nil || c.config.AcceptHandshakeMetadata == nil {
		return nil
	}
	if !s.echAccepted {
		return alertError(alertIllegalParameter, &ProtocolError{"handshake metadata requires ECH"})
	}
	response, err := c.config.AcceptHandshakeMetadata(c.clientHelloInfo(s.ch), bytes.Clone(s.ch.handshakeMetadata))
	if err != nil {
		return err
	}
	wire, err := marshalHandshakeMetadata(response)
	if err != nil {
		return err
	}
	if ee.extensions == nil {
		ee.extensions = make(map[uint16][]byte, 1)
	}
	ee.extensions[extHandshakeMetadata] = wire
	s.handshakeMetadata = s.ch.handshakeMetadata
	return nil
}
