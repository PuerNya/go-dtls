package dtls13

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	alertLevelWarning          uint8 = 1
	alertLevelFatal            uint8 = 2
	alertCloseNotify           uint8 = 0
	alertUnexpectedMessage     uint8 = 10
	alertRecordOverflow        uint8 = 22
	alertHandshakeFailure      uint8 = 40
	alertBadCertificate        uint8 = 42
	alertIllegalParameter      uint8 = 47
	alertUnknownCA             uint8 = 48
	alertAccessDenied          uint8 = 49
	alertDecodeError           uint8 = 50
	alertDecryptError          uint8 = 51
	alertTooManyCIDsRequest    uint8 = 52
	alertProtocolVersion       uint8 = 70
	alertUserCanceled          uint8 = 90
	alertMissingExtension      uint8 = 109
	alertUnsupportedExtension  uint8 = 110
	alertCertificateRequired   uint8 = 116
	alertGeneralError          uint8 = 117
	alertNoApplicationProtocol uint8 = 120
	alertECHRequired           uint8 = 121
)

type alertMessage struct{ level, description uint8 }

// AlertError reports a fatal TLS alert received from the peer. Its numeric
// value is the alert description assigned by TLS. Compare a known description
// with [errors.Is], or use [errors.As] to obtain the value without depending on
// an error string.
type AlertError uint8

func (e AlertError) Error() string {
	return fmt.Sprintf("dtls13: received fatal alert %d", uint8(e))
}

type localAlertError struct {
	description uint8
	err         error
}

func (e *localAlertError) Error() string { return e.err.Error() }
func (e *localAlertError) Unwrap() error { return e.err }

func alertError(description uint8, err error) error {
	return &localAlertError{description: description, err: err}
}

func protocolAlert(err error) (uint8, bool) {
	var local *localAlertError
	if errors.As(err, &local) {
		return local.description, true
	}
	var protocol *ProtocolError
	if !errors.As(err, &protocol) {
		return 0, false
	}
	reason := strings.ToLower(protocol.Reason)
	for _, marker := range []string{"unexpected", " before ", "duplicate", "omitted"} {
		if strings.Contains(reason, marker) {
			return alertUnexpectedMessage, true
		}
	}
	for _, marker := range []string{"truncated", "malformed", "decode", "invalid record length", "invalid handshake fragment"} {
		if strings.Contains(reason, marker) {
			return alertDecodeError, true
		}
	}
	return alertIllegalParameter, true
}

func outboundAlert(err error) (uint8, bool) {
	if err == nil {
		return 0, false
	}
	if description, ok := protocolAlert(err); ok {
		return description, true
	}
	var peer AlertError
	var network net.Error
	if errors.As(err, &peer) || errors.As(err, &network) ||
		errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) ||
		errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return 0, false
	}
	return alertGeneralError, true
}

func (a alertMessage) marshal() ([]byte, error) {
	if a.level != alertLevelWarning && a.level != alertLevelFatal {
		return nil, &ProtocolError{"invalid alert level"}
	}
	return []byte{a.level, a.description}, nil
}
func parseAlert(b []byte) (alertMessage, error) {
	if len(b) != 2 {
		return alertMessage{}, alertError(alertDecodeError, &ProtocolError{"invalid alert length"})
	}
	return alertMessage{level: b[0], description: b[1]}, nil
}

func (a alertMessage) isCloseNotify() bool {
	return (a.level == alertLevelWarning || a.level == alertLevelFatal) && a.description == alertCloseNotify
}

func (a alertMessage) isUserCanceled() bool {
	return (a.level == alertLevelWarning || a.level == alertLevelFatal) && a.description == alertUserCanceled
}

type closureState struct {
	received bool
	number   recordNumber
}

func (c *closureState) receive(number recordNumber) {
	if !c.received || recordNumberLess(c.number, number) {
		c.received = true
		c.number = number
	}
}
func (c *closureState) ignore(number recordNumber) bool {
	return c.received && recordNumberLess(c.number, number)
}
func recordNumberLess(a, b recordNumber) bool {
	return a.epoch < b.epoch || (a.epoch == b.epoch && a.sequence < b.sequence)
}
