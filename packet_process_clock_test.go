//go:build !windows

package dtls13

import (
	"testing"
	"time"
)

func startPacketProcessTimer(t *testing.T) func() time.Duration {
	started := time.Now()
	return func() time.Duration { return time.Since(started) }
}
