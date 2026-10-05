//go:build !windows

package dtls13

import (
	"errors"
	"syscall"
)

// isMessageTooLongErrno reports whether err is the platform "message too long"
// failure. On Unix-like systems an oversized datagram send returns EMSGSIZE.
func isMessageTooLongErrno(err error) bool {
	return errors.Is(err, syscall.EMSGSIZE)
}
