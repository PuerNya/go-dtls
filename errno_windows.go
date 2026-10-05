//go:build windows

package dtls13

import (
	"errors"
	"syscall"
)

// isMessageTooLongErrno reports whether err is the Windows "message too long"
// failure. On Windows syscall.EMSGSIZE is not the Winsock error: an oversized
// UDP send returns WSAEMSGSIZE, whose Go spelling is syscall.Errno(10040), so
// the constant is named here rather than in shared code.
func isMessageTooLongErrno(err error) bool {
	return errors.Is(err, syscall.Errno(10040))
}
