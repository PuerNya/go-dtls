package dtls13

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

// Decisions replay by link/write ordinal; a changed packet size or queue outcome
// aborts replay rather than claiming that concurrent DTLS timing was reproduced.
type packetNetworkDecision struct {
	Kind, Phase      string
	Write            uint64
	Size, Rate, Loss int
	Lost, Duplicate  bool
	Delay            [2]time.Duration
	Corrupt          [2]byte
	Queued           [2]bool
	Copy             int
	At               int64
}

type packetNetworkTrace struct {
	encode *json.Encoder
	decode *json.Decoder
	phase  string
}

func (c *impairedPacketConn) nextDecision(size int) (packetNetworkDecision, error) {
	d := packetNetworkDecision{Kind: "decision", Write: c.stats.writes + 1, Size: size}
	if c.trace != nil {
		d.Phase = c.trace.phase
	}
	if c.trace != nil && c.trace.decode != nil {
		for {
			if err := c.trace.decode.Decode(&d); err != nil {
				return d, fmt.Errorf("replay decision %d: %w", c.stats.writes+1, err)
			}
			if d.Kind == "decision" {
				break
			}
		}
		if d.Write != c.stats.writes+1 || d.Size != size || d.Phase != c.trace.phase || d.Rate < c.profile.minRate || d.Rate > c.profile.maxRate || d.Loss < 0 || d.Loss > 1000 || d.Delay[0] < 0 || d.Delay[1] < 0 {
			return d, errors.New("packet decision replay diverged or is invalid")
		}
		c.rate = d.Rate
		return d, nil
	}
	if c.stats.writes%16 == 0 {
		if len(c.profile.rateChoices) == 0 {
			c.rate = c.profile.minRate + c.random.IntN(c.profile.maxRate-c.profile.minRate+1)
		} else {
			c.rate = c.profile.rateChoices[c.random.IntN(len(c.profile.rateChoices))]
		}
	}
	d.Rate, d.Loss = c.rate, c.profile.loss
	if c.profile.maxLoss > c.profile.loss {
		if c.lossRemaining == 0 {
			c.lossRate = c.profile.loss + c.random.IntN(c.profile.maxLoss-c.profile.loss+1)
			c.lossRemaining = 1 + c.random.IntN(64)
		}
		d.Loss = c.lossRate
		c.lossRemaining--
	}
	d.Lost = c.random.IntN(1000) < d.Loss
	if !d.Lost {
		d.Duplicate = c.random.IntN(1000) < c.profile.duplicate
		copies := 1
		if d.Duplicate {
			copies++
		}
		for i := range copies {
			if size > 0 && c.random.IntN(1000) < c.corrupt {
				d.Corrupt[i] = byte(1 << c.random.IntN(8))
			}
			d.Delay[i] = c.profile.delay + time.Duration(c.random.Int64N(int64(c.profile.jitter)+1))
		}
	}
	return d, nil
}

func (c *impairedPacketConn) recordDecision(d packetNetworkDecision) error {
	if c.trace == nil || c.trace.encode == nil {
		return nil
	}
	d.At = time.Now().UnixNano()
	return c.trace.encode.Encode(d)
}

func TestImpairedPacketConnReplay(t *testing.T) {
	profile := packetNetworkProfile{delay: time.Hour, minRate: 1 << 20, maxRate: 2 << 20, loss: 100, maxLoss: 900, duplicate: 500, queuePackets: 3, queueBytes: 12}
	var recorded bytes.Buffer
	to := &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 9}
	run := func(seed uint64, trace *packetNetworkTrace) packetNetworkStats {
		conn := newImpairedPacketConn(t, seed, profile)
		conn.trace = trace
		conn.corrupt = 1000
		for range 512 {
			if _, err := conn.WriteTo([]byte("four"), to); err != nil {
				t.Fatal(err)
			}
		}
		if err := conn.Close(); err != nil {
			t.Fatal(err)
		}
		return conn.snapshot()
	}
	before := run(99, &packetNetworkTrace{encode: json.NewEncoder(&recorded), phase: "traffic"})
	after := run(12345, &packetNetworkTrace{decode: json.NewDecoder(bytes.NewReader(recorded.Bytes())), phase: "traffic"})
	if before != after || before.lost == 0 || before.overflow == 0 || before.corrupted == 0 || before.duplicated == 0 || before.canceled == 0 {
		t.Fatalf("replay accounting: %+v -> %+v", before, after)
	}
	for _, input := range []string{"", `{"Kind":"decision","Write":1,"Size":5,"Rate":1048576}`} {
		conn := newImpairedPacketConn(t, 1, profile)
		conn.trace = &packetNetworkTrace{decode: json.NewDecoder(bytes.NewBufferString(input))}
		if _, err := conn.WriteTo([]byte("four"), to); err == nil {
			t.Fatal("invalid replay was accepted")
		}
	}
	// Confirm there is a complete decision for every original write.
	decoder := json.NewDecoder(bytes.NewReader(recorded.Bytes()))
	var decisions int
	for {
		var event packetNetworkDecision
		if err := decoder.Decode(&event); err != nil {
			if !errors.Is(err, io.EOF) {
				t.Fatal(err)
			}
			break
		}
		if event.Kind == "decision" {
			decisions++
		}
	}
	if decisions != 512 {
		t.Fatalf("recorded %d decisions", decisions)
	}
}
