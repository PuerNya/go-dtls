package dtls13

import (
	"errors"
	"net"
)

// ConfigError reports an invalid local configuration or an operation that
// exceeds a configured resource limit. Use [errors.As] to inspect the Reason.
type ConfigError struct {
	// Reason describes the rejected setting or limit.
	Reason string
}

func (e *ConfigError) Error() string { return "dtls13: invalid configuration: " + e.Reason }

// ProtocolError reports malformed protocol data or an operation that is not
// valid in the connection's current DTLS state. When caused by authenticated
// peer input, the connection sends the corresponding fatal alert when
// possible. Use [errors.As] to inspect the Reason.
type ProtocolError struct {
	// Reason describes the protocol or state-machine violation.
	Reason string
}

func (e *ProtocolError) Error() string { return "dtls13: protocol error: " + e.Reason }

// vectorOverflowError reports that a length-prefixed wire vector does not fit
// the width it was written with. It wraps a [ProtocolError] so callers that
// classify by error type still see the protocol failure, while callers that
// must react to one specific width can match this type instead of comparing
// diagnostic text.
type vectorOverflowError struct {
	bits int
	err  error
}

func (e *vectorOverflowError) Error() string { return e.err.Error() }
func (e *vectorOverflowError) Unwrap() error { return e.err }

// ErrDatagramTooLarge indicates that an application datagram exceeds the
// current path MTU or the DTLS record size limit. A transport can still
// return it when [Config.IgnorePathMTU] is enabled. [Conn.WriteDatagram]
// may wrap it in a net.OpError; use [errors.Is] to test for it. No partial
// application record is sent when this error is returned.
var ErrDatagramTooLarge = errors.New("dtls13: datagram too large")

func datagramTooLargeError(addr net.Addr) error {
	return &net.OpError{Op: "write", Net: "dtls", Addr: addr, Err: ErrDatagramTooLarge}
}
