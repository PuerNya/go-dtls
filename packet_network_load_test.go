package dtls13

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestPacketNetworkScenarios(t *testing.T) {
	count, messages, seeds := 4, 128, 2
	stress := os.Getenv("GO_DTLS_STRESS") == "1"
	if stress {
		count, messages, seeds = 64, 2048, 5
	}
	seedBase := uint64(1)
	if value := os.Getenv("GO_DTLS_NETWORK_SEED"); value != "" {
		var err error
		seedBase, err = strconv.ParseUint(value, 10, 64)
		if err != nil {
			t.Fatal("GO_DTLS_NETWORK_SEED must be a uint64")
		}
	} else if stress {
		var entropy [8]byte
		if _, err := rand.Read(entropy[:]); err != nil {
			t.Fatal(err)
		}
		seedBase = binary.LittleEndian.Uint64(entropy[:])
	}
	t.Logf("network_seed_base=%d", seedBase)
	names := []string{"Variable", "Constrained", "ExtremeBursts", "LossStorm", "RateCollapse"}
	if !stress {
		names = names[:3]
	}
	for networkIndex, name := range names {
		for cidIndex, cid := range []bool{false, true} {
			for seed := 1; seed <= seeds; seed++ {
				t.Run(fmt.Sprintf("%s/CID%t/Seed%d", name, cid, seed), func(t *testing.T) {
					connections := count
					_, traffic := packetProcessNetwork(name)
					if traffic != nil {
						connections = min(count, 8)
					}
					spec := packetProcessSpec{Network: name, Connections: connections, Operations: connections * messages, Workers: connections, Payload: 1200, CID: cid, Burst: true, Seed: seedBase + uint64(networkIndex*10000+cidIndex*1000+seed*100)}
					run := newPacketProcessRun(t, spec)
					stats := run.stage(t, "traffic")
					if stats["client"].Offered != uint64(spec.Operations) || stats["server"].Offered != uint64(spec.Operations) || stats["client"].Received+stats["server"].Received == 0 {
						t.Fatalf("invalid application accounting: %+v", stats)
					}
					run.close(t)
				})
			}
		}
	}
}

func TestPacketNetworkResumption(t *testing.T) {
	for _, early := range []bool{false, true} {
		t.Run(fmt.Sprintf("Early%t", early), func(t *testing.T) {
			run := newPacketProcessRun(t, packetProcessSpec{Connections: 1, Operations: 1, Workers: 1, Payload: 64, Network: "Constrained", Seed: 90, Resume: true, Early: early})
			for round := 1; round <= 2; round++ {
				run.stage(t, "reconnect")
				run.stage(t, "handshake")
				run.stage(t, "warmup")
			}
			run.close(t)
		})
	}
}

func TestPacketNetworkSlowReaderBound(t *testing.T) {
	run := newPacketProcessRun(t, packetProcessSpec{Connections: 1, Operations: 1, Workers: 1, Payload: 32, Network: "Variable", Seed: 71, SlowBound: true})
	run.stage(t, "bound")
	run.close(t)
}
