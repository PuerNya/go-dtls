package dtls13

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Only the three workers carry packets; the parent coordinates over stdio.
type packetProcessSpec struct {
	Role, Address, Network, Output, Observe, CPUs string
	Connections, Operations, Workers, Payload     int
	CID, Burst, Raw                               bool
	Resume, Early, SlowBound                      bool
	Seed                                          uint64
	CIDLengths                                    []int
	Certificate, PrivateKey                       []byte
}

type packetProcessStats struct {
	Elapsed                                          time.Duration
	Allocated, Allocs, Heap, HeapObjects, Stack, Sys uint64
	GCs                                              uint32
	GCPause                                          time.Duration
	Goroutines                                       int
	Resources                                        packetProcessResources
	Wire                                             packetProcessWireStats
	KernelUDP                                        map[string]uint64
	Runtime                                          map[string]any
	Offered, Received, Reordered                     uint64
	LatencySamples, LatencyZeroSamples               int
	Peer                                             int
	Payload                                          int
	P50, P95, P99, Max                               time.Duration
}

type packetProcessResources struct {
	User, System                                                                                                          time.Duration
	RSS, PeakRSS, PageFaults, MinorFaults, MajorFaults, VoluntarySwitches, InvoluntarySwitches, IOReadBytes, IOWriteBytes uint64
}

func (r packetProcessResources) delta(before packetProcessResources) packetProcessResources {
	r.User -= before.User
	r.System -= before.System
	r.PageFaults -= before.PageFaults
	r.MinorFaults -= before.MinorFaults
	r.MajorFaults -= before.MajorFaults
	r.VoluntarySwitches -= before.VoluntarySwitches
	r.InvoluntarySwitches -= before.InvoluntarySwitches
	r.IOReadBytes -= before.IOReadBytes
	r.IOWriteBytes -= before.IOWriteBytes
	return r
}

func packetProcessCPUs(value string) ([]int, error) {
	if value == "" {
		return nil, nil
	}
	var cpus []int
	for entry := range strings.SplitSeq(value, ",") {
		cpu, err := strconv.Atoi(entry)
		if err != nil || cpu < 0 || cpu >= 1024 || slices.Contains(cpus, cpu) {
			return nil, fmt.Errorf("invalid CPU set %q", value)
		}
		cpus = append(cpus, cpu)
	}
	return cpus, nil
}

type packetProcessEvent struct {
	Phase, Address string
	PID            int
	Stats          packetProcessStats
	Links          []packetRelayStats
}

type packetProcessWorker struct {
	cmd    *exec.Cmd
	in     io.WriteCloser
	events chan packetProcessEvent
	done   chan error
	logs   lockedBuffer
}

func startPacketProcess(t testing.TB, spec packetProcessSpec) *packetProcessWorker {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	w := &packetProcessWorker{events: make(chan packetProcessEvent, 1), done: make(chan error, 1)}
	//nolint:gosec // Re-exec only this test binary; no shell or external command input.
	w.cmd = exec.CommandContext(ctx, executable, "-test.run=^TestPacketProcessWorker$", "-test.timeout=7m")
	w.cmd.Env = append(os.Environ(), "GO_DTLS_PACKET_WORKER=1", "GOMAXPROCS="+strconv.Itoa(runtime.GOMAXPROCS(0)))
	w.cmd.Stderr = &w.logs
	w.in, err = w.cmd.StdinPipe()
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	out, err := w.cmd.StdoutPipe()
	if err != nil {
		cancel()
		t.Fatal(err)
	}
	if err = w.cmd.Start(); err != nil {
		cancel()
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = w.in.Close()
		cancel()
		select {
		case <-w.done:
		case <-time.After(5 * time.Second):
			t.Errorf("%s process did not exit", spec.Role)
		}
		if t.Failed() {
			logs := w.logs.String()
			if len(logs) > 8192 {
				logs = logs[:8192]
			}
			t.Logf("%s worker: %s", spec.Role, logs)
		}
	})
	go func() {
		scanner := bufio.NewScanner(out)
		scanner.Buffer(make([]byte, 4096), 4<<20)
		for scanner.Scan() {
			var event packetProcessEvent
			line := scanner.Bytes()
			if len(line) > 0 && line[0] == '{' && json.Unmarshal(line, &event) == nil && event.Phase != "" {
				select {
				case w.events <- event:
				case <-ctx.Done():
				}
			} else {
				_, _ = w.logs.Write(append(bytes.Clone(line), '\n'))
			}
		}
		if err := scanner.Err(); err != nil {
			_, _ = fmt.Fprintln(&w.logs, err)
		}
		close(w.events)
		w.done <- w.cmd.Wait()
	}()
	//nolint:gosec // The generated test key goes only to an owned worker stdin pipe.
	if err = json.NewEncoder(w.in).Encode(spec); err != nil {
		t.Fatal(err)
	}
	return w
}

func (w *packetProcessWorker) send(t testing.TB, phase string) {
	t.Helper()
	if err := json.NewEncoder(w.in).Encode(phase); err != nil {
		t.Fatal(err)
	}
}

func (w *packetProcessWorker) receive(t testing.TB, phase string) packetProcessEvent {
	t.Helper()
	select {
	case event, ok := <-w.events:
		if !ok {
			t.Fatalf("worker exited during %s: %s", phase, w.logs.String())
		}
		if event.Phase != phase {
			t.Fatalf("worker phase %q, want %q", event.Phase, phase)
		}
		return event
	case <-time.After(3 * time.Minute):
		t.Fatalf("worker timed out during %s: %s", phase, w.logs.String())
	}
	return packetProcessEvent{}
}

type packetProcessRun struct {
	workers map[string]*packetProcessWorker
}

func newPacketProcessRun(t testing.TB, spec packetProcessSpec) *packetProcessRun {
	t.Helper()
	certificate, _ := testServerCertificate(t)
	key, err := x509.MarshalPKCS8PrivateKey(certificate.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	spec.Certificate, spec.PrivateKey = certificate.Certificate[0], key
	if spec.Workers == 0 {
		spec.Workers = 16
	}
	if spec.Payload == 0 {
		spec.Payload = 1200
	}
	if spec.Output == "" {
		spec.Output = os.Getenv("GO_DTLS_PROCESS_OUTPUT")
	}
	if spec.Observe == "" {
		spec.Observe = os.Getenv("GO_DTLS_PROCESS_OBSERVE")
	}
	if spec.Observe != "" && spec.Output == "" {
		t.Fatal("GO_DTLS_PROCESS_OBSERVE requires GO_DTLS_PROCESS_OUTPUT")
	}
	if spec.Output != "" {
		spec.Output = filepath.Join(spec.Output, strings.ReplaceAll(t.Name(), "/", "_"))
		//nolint:gosec // The local operator selects the artifact directory.
		if err := os.MkdirAll(spec.Output, 0700); err != nil {
			t.Fatal(err)
		}
	}
	run := &packetProcessRun{workers: make(map[string]*packetProcessWorker)}
	for _, role := range []string{"server", "relay", "client"} {
		spec.Role = role
		spec.CPUs = ""
		for entry := range strings.SplitSeq(os.Getenv("GO_DTLS_PROCESS_CPUS"), ";") {
			if name, cpus, ok := strings.Cut(entry, "="); ok && name == role {
				spec.CPUs = cpus
			}
		}
		worker := startPacketProcess(t, spec)
		run.workers[role] = worker
		ready := worker.receive(t, "ready")
		t.Logf("packet_process role=%s pid=%d address=%s cpus=%s", role, ready.PID, ready.Address, spec.CPUs)
		if role != "client" {
			spec.Address = ready.Address
		}
	}
	handshake := run.stage(t, "handshake")
	payload := min(handshake["client"].Payload, handshake["server"].Payload)
	for _, role := range []string{"client", "server"} {
		worker := run.workers[role]
		worker.send(t, fmt.Sprintf("payload:%d", payload))
		worker.receive(t, "payload")
	}
	run.stage(t, "warmup")
	return run
}

func (r *packetProcessRun) stage(t testing.TB, phase string) map[string]packetProcessStats {
	t.Helper()
	roles := []string{"relay", "server", "client"}
	for _, role := range roles {
		r.workers[role].send(t, "begin:"+phase)
	}
	for _, role := range roles {
		r.workers[role].receive(t, "begin:"+phase)
	}
	for _, role := range []string{"server", "client"} {
		r.workers[role].send(t, phase)
	}
	for _, role := range []string{"server", "client"} {
		r.workers[role].receive(t, phase)
	}
	for _, role := range roles {
		r.workers[role].send(t, "finish:"+phase)
	}
	stats := make(map[string]packetProcessStats)
	for _, role := range roles {
		event := r.workers[role].receive(t, "finish:"+phase)
		stats[role] = event.Stats
		// Full per-link accounting is retained in the worker artifact.
		event.Links = nil
		encoded, err := json.Marshal(event)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("packet_measurement role=%s %s", role, encoded)
	}
	return stats
}

func (r *packetProcessRun) close(t testing.TB) {
	t.Helper()
	r.stage(t, "recovery")
	// Keep the relay alive while both endpoints exchange close_notify.
	r.stage(t, "close")
	for _, role := range []string{"client", "server", "relay"} {
		worker := r.workers[role]
		worker.send(t, "stop")
		worker.receive(t, "stop")
		if err := <-worker.done; err != nil {
			t.Fatalf("%s process: %v: %s", role, err, worker.logs.String())
		}
		worker.done <- nil
	}
}

func TestPacketProcessWorker(t *testing.T) {
	if os.Getenv("GO_DTLS_PACKET_WORKER") != "1" {
		t.Skip("subprocess entry point")
	}
	decoder, encoder := json.NewDecoder(os.Stdin), json.NewEncoder(os.Stdout)
	var spec packetProcessSpec
	if err := decoder.Decode(&spec); err != nil {
		t.Fatal(err)
	}
	if spec.Connections < 1 || spec.Connections > 4096 || spec.Operations < 1 || spec.Workers < 1 || spec.Payload < 24 || spec.Payload > 1200 {
		t.Fatal("invalid process workload")
	}
	if err := setPacketProcessAffinity(spec.CPUs); err != nil {
		t.Fatal(err)
	}
	initialGoroutines := runtime.NumGoroutine()
	config := &Config{ServerName: "server.test", SessionTicketsDisabled: true, MaxPendingConnections: spec.Connections,
		MTU: 1500, HandshakeTimeout: 20 * time.Second, FlightInterval: 40 * time.Millisecond, MaxFlightInterval: 500 * time.Millisecond,
		MaxBufferedApplicationDatagrams: spec.Operations/spec.Connections + 128, MaxBufferedApplicationData: (spec.Operations/spec.Connections + 128) * 1200}
	certificate, err := x509.ParseCertificate(spec.Certificate)
	if err != nil {
		t.Fatal(err)
	}
	config.RootCAs = x509.NewCertPool()
	config.RootCAs.AddCert(certificate)
	key, err := x509.ParsePKCS8PrivateKey(spec.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}
	config.Certificates = []tls.Certificate{{Certificate: [][]byte{spec.Certificate}, PrivateKey: key}}
	if spec.CID {
		var next atomic.Uint64
		config.GetConnectionID = func() ([]byte, error) {
			n := next.Add(1)
			if len(spec.CIDLengths) == 0 {
				return binary.BigEndian.AppendUint64(nil, n), nil
			}
			length := spec.CIDLengths[(n-1)%uint64(len(spec.CIDLengths))]
			if length < 4 || length > 32 {
				return nil, errors.New("invalid workload CID length")
			}
			cid := make([]byte, length)
			cid[0] = byte(length)
			var encoded [8]byte
			binary.BigEndian.PutUint64(encoded[:], n)
			copy(cid[max(1, length-8):], encoded[max(0, 9-length):])
			return cid, nil
		}
	}
	if spec.Resume {
		config.SessionTicketsDisabled = false
		config.ClientSessionCache = NewLRUClientSessionCache(4)
		config.MaxEarlyData = 1024
		config.EnableEarlyDataIO = true
		config.EarlyDataReplayCache = NewLRUEarlyDataReplayCache(16)
		config.AllowEarlyDataWithoutCookie = true
	}
	if spec.SlowBound && spec.Role == "server" {
		config.MaxBufferedApplicationDatagrams, config.MaxBufferedApplicationData = 4, 64
	}
	var wire packetProcessWire
	var rawServer net.PacketConn
	var rawClients []*net.UDPConn
	var listener *Listener
	var relay *packetProcessRelay
	conns := make([]*Conn, spec.Connections)
	address := spec.Address
	switch spec.Role {
	case "server":
		raw, listenErr := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1)})
		if listenErr != nil {
			t.Fatal(listenErr)
		}
		if err := raw.SetReadBuffer(4 << 20); err != nil {
			t.Fatal(err)
		}
		var socket net.PacketConn = raw
		if spec.Observe == "wire" {
			socket = &countedPacketProcessSocket{PacketConn: raw, wire: &wire}
		}
		if spec.Raw {
			rawServer = socket
			t.Cleanup(func() { _ = rawServer.Close() })
		} else {
			listener = NewListener(socket, config)
			t.Cleanup(func() { _ = listener.Close() })
		}
		address = raw.LocalAddr().String()
	case "relay":
		relay = newPacketProcessRelay(t, spec)
		address = relay.front.LocalAddr().String()
	case "client":
		remote, resolveErr := net.ResolveUDPAddr("udp4", address)
		if resolveErr != nil {
			t.Fatal(resolveErr)
		}
		if spec.Raw {
			rawClients = make([]*net.UDPConn, spec.Connections)
		}
		for i := range conns {
			inner, dialErr := net.DialUDP("udp4", nil, remote)
			if dialErr != nil {
				t.Fatal(dialErr)
			}
			if err := inner.SetReadBuffer(1 << 20); err != nil {
				t.Fatal(err)
			}
			if spec.Raw {
				rawClients[i] = inner
			} else {
				var transport net.Conn = inner
				if spec.Observe == "wire" {
					transport = &countedPacketProcessConn{Conn: inner, wire: &wire}
				}
				conns[i] = Client(transport, config)
			}
		}
	default:
		t.Fatal("invalid process role")
	}
	t.Cleanup(func() {
		for _, raw := range rawClients {
			_ = raw.Close()
		}
		for _, conn := range conns {
			if conn != nil {
				_ = conn.Close()
			}
		}
	})
	send := func(event packetProcessEvent) {
		if err := encoder.Encode(event); err != nil {
			t.Fatal(err)
		}
		if spec.Output != "" && (strings.HasPrefix(event.Phase, "finish:") || event.Phase == "stop") {
			data, err := json.MarshalIndent(event, "", "  ")
			if err != nil {
				t.Fatal(err)
			}
			path := filepath.Join(spec.Output, spec.Role+"-"+strings.TrimPrefix(event.Phase, "finish:")+".json")
			if err := os.WriteFile(path, data, 0600); err != nil {
				t.Fatal(err)
			}
		}
	}
	send(packetProcessEvent{Phase: "ready", Address: address, PID: os.Getpid()})
	observation := packetProcessObservation{wire: &wire}
	var result packetProcessStats
	activePhase := ""
	round := 0
	for {
		var phase string
		if err := decoder.Decode(&phase); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			t.Fatal(err)
		}
		if stage, ok := strings.CutPrefix(phase, "begin:"); ok {
			activePhase = stage
			if relay != nil {
				relay.phase(t, stage)
			}
			observation.begin(t, spec, stage)
			result = packetProcessStats{}
		} else if value, ok := strings.CutPrefix(phase, "payload:"); ok {
			payload, err := strconv.Atoi(value)
			if err != nil || payload < 24 || payload > spec.Payload {
				t.Fatal("invalid negotiated load payload")
			}
			spec.Payload = payload
			phase = "payload"
		} else if strings.HasPrefix(phase, "finish:") {
			if relay != nil || activePhase == "miss" {
				result = observation.finish(t, result)
			}
		} else {
			if phase != "stop" {
				observation.reset(t)
			}
			switch phase {
			case "handshake":
				result.Payload = spec.Payload
				if spec.Raw {
					break
				}
				packetProcessHandshake(t, spec, listener, conns)
				for _, conn := range conns {
					result.Payload = min(result.Payload, conn.PathMTU()-conn.RecordOverhead())
				}
				if spec.Resume {
					for _, conn := range conns {
						state := conn.ConnectionState()
						if state.DidResume != (round > 0) {
							t.Fatal("unexpected resumed handshake state")
						}
						if spec.Early && round > 0 && spec.Role == "client" && state.EarlyData != EarlyDataAccepted {
							t.Fatal("resumption did not accept early data")
						}
						if spec.Role == "client" {
							deadline := time.Now().Add(5 * time.Second)
							for {
								if _, ok := config.ClientSessionCache.Get("server.test"); ok {
									break
								}
								if time.Now().After(deadline) {
									t.Fatal("session ticket was not recovered")
								}
								time.Sleep(time.Millisecond)
							}
						} else if spec.Early && round > 0 {
							if err := conn.SetReadDeadline(time.Now().Add(50 * time.Millisecond)); err != nil {
								t.Fatal(err)
							}
							var p [16]byte
							n, info, err := conn.ReadDatagram(p[:])
							if err == nil {
								if !info.EarlyData || string(p[:n]) != "early" {
									t.Fatal("invalid early datagram")
								}
							} else if e, ok := errors.AsType[net.Error](err); !ok || !e.Timeout() {
								t.Fatal(err)
							}
							if err := conn.SetReadDeadline(time.Time{}); err != nil {
								t.Fatal(err)
							}
							conn.startRecordReader()
						}
					}
				}
			case "reconnect":
				for i, conn := range conns {
					local := conn.LocalAddr()
					_ = conn.Close()
					if spec.Role == "client" {
						remote, err := net.ResolveUDPAddr("udp4", spec.Address)
						if err != nil {
							t.Fatal(err)
						}
						inner, err := net.DialUDP("udp4", local.(*net.UDPAddr), remote)
						if err != nil {
							t.Fatal(err)
						}
						if spec.Early {
							conns[i] = ClientEarly(inner, config)
							if _, err := conns[i].WriteDatagram([]byte("early")); err != nil {
								t.Fatal(err)
							}
						} else {
							conns[i] = Client(inner, config)
						}
					} else {
						conns[i] = nil
					}
				}
				round++
			case "bound":
				conn := conns[0]
				if spec.Role == "client" {
					for range 32 {
						if _, err := conn.WriteDatagram([]byte("unread")); err != nil {
							t.Fatal(err)
						}
					}
				} else {
					deadline := time.Now().Add(3 * time.Second)
					for {
						conn.inputMu.Lock()
						err, queued, queuedBytes := conn.readErr, len(conn.applicationDatagrams), conn.bufferedApplicationBytes
						conn.inputMu.Unlock()
						if queued > 4 || queuedBytes > 64 {
							t.Fatal("application queue exceeded bound")
						}
						if err != nil {
							if _, ok := errors.AsType[*ProtocolError](err); !ok {
								t.Fatal(err)
							}
							break
						}
						if time.Now().After(deadline) {
							t.Fatal("slow reader did not reach its bound")
						}
						time.Sleep(time.Millisecond)
					}
				}
			case "miss":
				if spec.Role == "client" {
					inner := conns[0].conn
					payload := bytes.Repeat([]byte{0xff}, 32)
					payload[0] = unifiedHeaderFixed | unifiedHeaderCID
					for i := range spec.Operations {
						if _, err := inner.Write(payload); err != nil {
							t.Fatal(err)
						}
						if i%16 == 15 {
							time.Sleep(time.Millisecond)
						}
					}
				} else {
					time.Sleep(time.Duration(spec.Operations/16+1)*time.Millisecond + 100*time.Millisecond)
				}
			case "warmup", "traffic":
				if spec.Raw {
					result = packetProcessRawTraffic(t, spec, phase, rawServer, rawClients)
				} else {
					result = packetProcessTraffic(t, spec, phase, conns)
				}
			case "recovery":
				if spec.Raw || spec.SlowBound {
					break
				}
				for _, conn := range conns {
					if err := conn.SetReadDeadline(time.Time{}); err != nil {
						t.Fatal(err)
					}
					conn.startRecordReader()
				}
				if spec.Burst {
					deadline := time.Now().Add(5 * time.Second)
					for _, conn := range conns {
						for {
							conn.writeMu.Lock()
							updated := conn.sendCipher != nil && conn.sendCipher.epoch >= 4
							conn.writeMu.Unlock()
							if updated {
								break
							}
							if time.Now().After(deadline) {
								t.Fatal("KeyUpdate did not recover")
							}
							time.Sleep(time.Millisecond)
						}
					}
				}
			case "close":
				for _, raw := range rawClients {
					_ = raw.Close()
				}
				if rawServer != nil {
					_ = rawServer.Close()
				}
				for _, conn := range conns {
					if conn == nil {
						continue
					}
					if err := conn.Close(); err != nil && !spec.SlowBound {
						t.Fatal(err)
					}
				}
				if listener != nil {
					_ = listener.Close()
				}
				clear(conns)
				deadline := time.Now().Add(3 * time.Second)
				for runtime.NumGoroutine() > initialGoroutines+4 && time.Now().Before(deadline) {
					time.Sleep(time.Millisecond)
				}
				runtime.GC()
				runtime.GC()
				if runtime.NumGoroutine() > initialGoroutines+4 {
					t.Fatal("endpoint goroutines did not converge")
				}
			case "stop":
				if relay != nil {
					relay.close(t)
				}
				event := packetProcessEvent{Phase: phase}
				if relay != nil {
					event.Links = relay.snapshot(t)
				}
				send(event)
				return
			default:
				t.Fatalf("unknown worker phase %s", phase)
			}
			if activePhase != "miss" {
				result = observation.finish(t, result)
			}
		}
		event := packetProcessEvent{Phase: phase, Stats: result}
		if relay != nil && strings.HasPrefix(phase, "finish:") {
			event.Links = relay.snapshot(t)
		}
		send(event)
	}
}

func packetProcessHandshake(t *testing.T, spec packetProcessSpec, listener *Listener, conns []*Conn) {
	t.Helper()
	var workers sync.WaitGroup
	if listener != nil {
		for i := range conns {
			conn, err := listener.Accept()
			if err != nil {
				t.Fatal(err)
			}
			conns[i] = conn
			workers.Go(func() {
				if err := conn.Handshake(); err != nil {
					t.Error(err)
				}
			})
		}
	} else {
		for worker := range min(spec.Workers, len(conns)) {
			workers.Go(func() {
				for i := worker; i < len(conns); i += min(spec.Workers, len(conns)) {
					if err := conns[i].Handshake(); err != nil {
						t.Error(err)
						return
					}
				}
			})
		}
	}
	workers.Wait()
	if t.Failed() {
		t.FailNow()
	}
}

func packetProcessTraffic(t *testing.T, spec packetProcessSpec, phase string, conns []*Conn) packetProcessStats {
	t.Helper()
	operations := spec.Operations
	if phase == "warmup" {
		operations = len(conns)
	}
	burst := spec.Burst && phase == "traffic"
	var workers sync.WaitGroup
	results := make([]packetProcessStats, len(conns))
	latencies := make([][]time.Duration, len(conns))
	if spec.Role == "server" || burst {
		for i, conn := range conns {
			workers.Go(func() {
				results[i], latencies[i] = packetProcessExchange(t, spec, conn, i, operations, phase, spec.Role == "client")
			})
		}
	} else {
		for worker := range min(spec.Workers, len(conns)) {
			workers.Go(func() {
				for i := worker; i < len(conns); i += min(spec.Workers, len(conns)) {
					results[i], latencies[i] = packetProcessExchange(t, spec, conns[i], i, operations, phase, true)
				}
			})
		}
	}
	workers.Wait()
	if t.Failed() {
		t.FailNow()
	}
	if spec.Role == "server" && phase == "warmup" {
		ordered := make([]*Conn, len(conns))
		for i, stats := range results {
			if stats.Peer < 0 || stats.Peer >= len(conns) || ordered[stats.Peer] != nil {
				t.Fatal("relay collapsed or duplicated peer identity")
			}
			ordered[stats.Peer] = conns[i]
		}
		copy(conns, ordered)
	}
	result := packetProcessStats{Payload: spec.Payload}
	for _, stats := range results {
		result.Offered += stats.Offered
		result.Received += stats.Received
		result.Reordered += stats.Reordered
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

func packetProcessExchange(t *testing.T, spec packetProcessSpec, conn *Conn, index, operations int, phase string, client bool) (packetProcessStats, []time.Duration) {
	t.Helper()
	burst := spec.Burst && phase == "traffic"
	size := min(spec.Payload, conn.PathMTU()-conn.RecordOverhead())
	if size < 24 {
		t.Error("insufficient datagram MTU")
		return packetProcessStats{}, nil
	}
	buffer, payload := make([]byte, size), make([]byte, size)
	count := operations / spec.Connections
	if index < operations%spec.Connections {
		count++
	}
	deadline := time.Now().Add(3 * time.Minute)
	if phase == "warmup" {
		deadline = time.Now().Add(10 * time.Second)
	}
	if burst {
		deadline = time.Now().Add(3*time.Second + time.Duration(count/16)*10*time.Millisecond)
	}
	if err := conn.SetReadDeadline(deadline); err != nil {
		t.Error(err)
		return packetProcessStats{}, nil
	}
	defer func() {
		// Idle connections must not retain the previous phase's deadline.
		if err := conn.SetReadDeadline(time.Time{}); err != nil {
			t.Error(err)
		}
	}()
	result := packetProcessStats{Peer: index}
	samples := make([]time.Duration, 0, count/16+1)
	seen := make([]bool, count)
	source := conn.RemoteAddr().String()
	var roundTripStart time.Time
	read := func() (int, uint64, bool) {
		n, info, err := conn.ReadDatagram(buffer)
		if err != nil {
			if e, ok := errors.AsType[net.Error](err); !burst || !ok || !e.Timeout() {
				t.Errorf("datagram read: %v", err)
			}
			return 0, 0, false
		}
		if n != size || info.Truncated || info.FullLength != size || info.EarlyData || info.Source == nil || info.Source.String() != source {
			t.Errorf("invalid datagram: n=%d info=%+v", n, info)
			return 0, 0, false
		}
		peer, sequence := int(binary.BigEndian.Uint64(buffer[:8])), binary.BigEndian.Uint64(buffer[8:16])
		if !client && result.Received == 0 && phase == "warmup" {
			index = peer
			result.Peer = peer
			count = operations / spec.Connections
			if index < operations%spec.Connections {
				count++
			}
			seen = make([]bool, count)
		}
		if peer != index || sequence >= uint64(count) || seen[sequence] {
			t.Errorf("invalid/duplicate datagram peer=%d want=%d sequence=%d count=%d", peer, index, sequence, count)
			return 0, 0, false
		}
		for _, value := range buffer[24:] {
			if value != byte(sequence) {
				t.Error("corrupted datagram reached application")
				return 0, 0, false
			}
		}
		seen[sequence] = true
		result.Received++
		if sequence%16 == 0 {
			latency := time.Since(time.Unix(0, int64(binary.BigEndian.Uint64(buffer[16:24]))))
			if client && !burst {
				latency = time.Since(roundTripStart)
			}
			samples = append(samples, latency)
		}
		return n, sequence, true
	}
	write := func(sequence int) bool {
		for i := range payload {
			payload[i] = byte(sequence)
		}
		binary.BigEndian.PutUint64(payload[:8], uint64(index))
		binary.BigEndian.PutUint64(payload[8:16], uint64(sequence))
		binary.BigEndian.PutUint64(payload[16:24], uint64(time.Now().UnixNano()))
		if _, err := conn.WriteDatagram(payload); err != nil {
			t.Errorf("datagram write: %v", err)
			return false
		}
		return true
	}
	if burst {
		written := make(chan uint64, 1)
		go func() {
			var sent uint64
			for sequence := range count {
				if sequence == count/2 {
					if err := conn.SendKeyUpdate(false); err != nil {
						t.Error(err)
						break
					}
				}
				if !write(sequence) {
					break
				}
				sent++
				if sequence%16 == 15 {
					pause := 2 * time.Millisecond
					if _, traffic := packetProcessNetwork(spec.Network); traffic != nil {
						pause = 10 * time.Millisecond
					}
					time.Sleep(pause)
				}
			}
			written <- sent
		}()
		highest := uint64(0)
		for {
			_, sequence, ok := read()
			if !ok {
				break
			}
			if sequence < highest {
				result.Reordered++
			}
			highest = max(highest, sequence)
			if !client && index%2 == 0 {
				time.Sleep(time.Millisecond)
			}
		}
		result.Offered = <-written
	} else {
		for sequence := 0; sequence < count; sequence++ {
			if client {
				roundTripStart = time.Now()
				if !write(sequence) {
					break
				}
				result.Offered++
			}
			n, got, ok := read()
			if !ok {
				break
			}
			if got != uint64(sequence) {
				t.Error("round-trip sequence mismatch")
				break
			}
			if !client {
				if _, err := conn.WriteDatagram(buffer[:n]); err != nil {
					t.Error(err)
					break
				}
				result.Offered++
			}
		}
	}
	return result, samples
}

func TestPacketProcessUDP(t *testing.T) {
	for _, cid := range []bool{false, true} {
		t.Run(fmt.Sprintf("CID%t", cid), func(t *testing.T) {
			run := newPacketProcessRun(t, packetProcessSpec{Connections: 3, Operations: 19, Workers: 2, Payload: 64, CID: cid})
			stats := run.stage(t, "traffic")
			for _, role := range []string{"client", "server"} {
				if stats[role].Offered != 19 || stats[role].Received != 19 {
					t.Fatalf("%s datagrams: %+v", role, stats[role])
				}
			}
			run.close(t)
		})
	}
}

func TestPacketProcessRelayCapacity(t *testing.T) {
	connections, operations := 4, 4096
	if os.Getenv("GO_DTLS_STRESS") == "1" {
		connections, operations = 1024, 1024*1024
	}
	run := newPacketProcessRun(t, packetProcessSpec{Connections: connections, Operations: operations, Workers: 16, Payload: 1200, Raw: true})
	stats := run.stage(t, "traffic")
	if stats["client"].Received != uint64(operations) || stats["server"].Received != uint64(operations) {
		t.Fatal("raw relay lost packets")
	}
	run.close(t)
}

func TestPacketProcessCIDScaling(t *testing.T) {
	if os.Getenv("GO_DTLS_STRESS") != "1" {
		t.Skip("set GO_DTLS_STRESS=1 for CID scaling")
	}
	for _, count := range []int{1, 64, 1024} {
		t.Run(fmt.Sprintf("Connections%d", count), func(t *testing.T) {
			run := newPacketProcessRun(t, packetProcessSpec{Connections: count, Operations: 65536, Workers: 16, Payload: 1200, CID: true, CIDLengths: []int{4, 8, 16}, Seed: 1})
			stats := run.stage(t, "traffic")
			if stats["client"].Received != 65536 {
				t.Fatal("CID scaling lost application datagrams")
			}
			run.stage(t, "miss")
			run.close(t)
		})
	}
}
