package dtls13

import (
	"bytes"
	"encoding/binary"
	"net"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type packetProcessWire struct {
	reads, writes, readBytes, writeBytes, readErrors, writeErrors atomic.Uint64
}

type packetProcessWireStats struct{ Reads, Writes, ReadBytes, WriteBytes, ReadErrors, WriteErrors uint64 }

func (w *packetProcessWire) snapshot() packetProcessWireStats {
	return packetProcessWireStats{w.reads.Load(), w.writes.Load(), w.readBytes.Load(), w.writeBytes.Load(), w.readErrors.Load(), w.writeErrors.Load()}
}

func (s packetProcessWireStats) delta(before packetProcessWireStats) packetProcessWireStats {
	return packetProcessWireStats{s.Reads - before.Reads, s.Writes - before.Writes, s.ReadBytes - before.ReadBytes, s.WriteBytes - before.WriteBytes, s.ReadErrors - before.ReadErrors, s.WriteErrors - before.WriteErrors}
}

func (w *packetProcessWire) read(n int, err error) {
	if err != nil {
		w.readErrors.Add(1)
	} else {
		w.reads.Add(1)
		w.readBytes.Add(uint64(n))
	}
}

func (w *packetProcessWire) write(n int, err error) {
	if err != nil {
		w.writeErrors.Add(1)
	} else {
		w.writes.Add(1)
		w.writeBytes.Add(uint64(n))
	}
}

type countedPacketProcessConn struct {
	net.Conn
	wire *packetProcessWire
}

func (c *countedPacketProcessConn) Read(p []byte) (int, error) {
	n, err := c.Conn.Read(p)
	c.wire.read(n, err)
	return n, err
}
func (c *countedPacketProcessConn) Write(p []byte) (int, error) {
	n, err := c.Conn.Write(p)
	c.wire.write(n, err)
	return n, err
}

type countedPacketProcessSocket struct {
	net.PacketConn
	wire *packetProcessWire
}

func (c *countedPacketProcessSocket) ReadFrom(p []byte) (int, net.Addr, error) {
	n, from, err := c.PacketConn.ReadFrom(p)
	c.wire.read(n, err)
	return n, from, err
}
func (c *countedPacketProcessSocket) WriteTo(p []byte, to net.Addr) (int, error) {
	n, err := c.PacketConn.WriteTo(p, to)
	c.wire.write(n, err)
	return n, err
}

// Raw UDP measures the relay's forwarding envelope with the same sockets,
// payload, concurrency and operation count, without DTLS work in the endpoints.
func packetProcessRawTraffic(t *testing.T, spec packetProcessSpec, phase string, server net.PacketConn, clients []*net.UDPConn) packetProcessStats {
	t.Helper()
	operations := spec.Operations
	if phase == "warmup" {
		operations = spec.Connections
	}
	result := packetProcessStats{Payload: spec.Payload}
	if server != nil {
		if err := server.SetDeadline(time.Now().Add(3 * time.Minute)); err != nil {
			t.Fatal(err)
		}
		buffer := make([]byte, spec.Payload)
		sources := make(map[string]bool)
		for range operations {
			n, from, err := server.ReadFrom(buffer)
			if err != nil || n != spec.Payload {
				t.Fatalf("raw server read: %d %v", n, err)
			}
			if phase == "warmup" {
				if sources[from.String()] {
					t.Fatal("relay collapsed raw client sources")
				}
				sources[from.String()] = true
			}
			result.Received++
			if _, err := server.WriteTo(buffer[:n], from); err != nil {
				t.Fatal(err)
			}
			result.Offered++
		}
		return result
	}
	count := min(spec.Workers, len(clients))
	results := make([]packetProcessStats, count)
	latencies := make([][]time.Duration, count)
	var workers sync.WaitGroup
	for worker := range count {
		workers.Go(func() {
			payload, buffer := make([]byte, spec.Payload), make([]byte, spec.Payload)
			for index := worker; index < len(clients); index += count {
				conn := clients[index]
				if err := conn.SetDeadline(time.Now().Add(3 * time.Minute)); err != nil {
					t.Error(err)
					return
				}
				messages := operations / spec.Connections
				if index < operations%spec.Connections {
					messages++
				}
				for sequence := range messages {
					binary.BigEndian.PutUint64(payload, uint64(index))
					binary.BigEndian.PutUint64(payload[8:], uint64(sequence))
					started := time.Now()
					if _, err := conn.Write(payload); err != nil {
						t.Error(err)
						return
					}
					results[worker].Offered++
					n, err := conn.Read(buffer)
					if err != nil || !bytes.Equal(buffer[:n], payload) {
						t.Errorf("raw UDP round trip: n=%d err=%v", n, err)
						return
					}
					results[worker].Received++
					if sequence%16 == 0 {
						latencies[worker] = append(latencies[worker], time.Since(started))
					}
				}
			}
		})
	}
	workers.Wait()
	if t.Failed() {
		t.FailNow()
	}
	for _, stats := range results {
		result.Offered += stats.Offered
		result.Received += stats.Received
	}
	samples := slices.Concat(latencies...)
	slices.Sort(samples)
	result.LatencySamples = len(samples)
	for _, sample := range samples {
		if sample == 0 {
			result.LatencyZeroSamples++
		}
	}
	if len(samples) > 0 {
		result.P50, result.P95, result.P99, result.Max = samples[len(samples)/2], samples[len(samples)*95/100], samples[len(samples)*99/100], samples[len(samples)-1]
	}
	return result
}
