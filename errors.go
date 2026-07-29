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

var (
	// ErrDatagramTooLarge indicates that an application datagram exceeds the
	// current path MTU or the DTLS record size limit. A transport can still
	// return it when [Config.IgnorePathMTU] is enabled. [Conn.WriteDatagram] and
	// [Conn.WriteEarlyData] may wrap it in a net.OpError; use [errors.Is] to test
	// for it. No partial application record is sent when this error is returned.
	ErrDatagramTooLarge = errors.New("dtls13: datagram too large")
	// ErrEarlyDataUnavailable indicates that no usable PSK ticket with an
	// early-data allowance was available, the connection is not a client, or
	// WriteEarlyData was already attempted. The handshake may still have
	// completed successfully and 1-RTT data may still be sent.
	ErrEarlyDataUnavailable = errors.New("dtls13: early data unavailable")
	// ErrEarlyDataRejected indicates that the peer completed the handshake but
	// did not accept the queued 0-RTT record, for example because of a
	// HelloRetryRequest, replay detection, or server policy. The caller decides
	// whether the operation is safe to retry as 1-RTT data.
	ErrEarlyDataRejected = errors.New("dtls13: early data rejected")
)

func datagramTooLargeError(addr net.Addr) error {
	return &net.OpError{Op: "write", Net: "dtls", Addr: addr, Err: ErrDatagramTooLarge}
}
