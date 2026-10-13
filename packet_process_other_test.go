//go:build !linux && !windows

package dtls13

import (
	"errors"
	"testing"
)

func readPacketProcessResources(t *testing.T) packetProcessResources { return packetProcessResources{} }

func readPacketKernelUDP(t *testing.T) map[string]uint64 { return nil }
func setPacketProcessAffinity(value string) error {
	if value != "" {
		return errors.New("process CPU affinity is supported only on Linux and Windows")
	}
	return nil
}
