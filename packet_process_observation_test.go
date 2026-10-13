package dtls13

import (
	"os"
	"path/filepath"
	"runtime"
	"runtime/metrics"
	"runtime/pprof"
	"runtime/trace"
	"slices"
	"strconv"
	"testing"
	"time"
)

type packetProcessObservation struct {
	elapsed      func() time.Duration
	before       runtime.MemStats
	resources    packetProcessResources
	wire         *packetProcessWire
	wireBefore   packetProcessWireStats
	kernelBefore map[string]uint64
	runtime      []metrics.Sample
	files        []*os.File
	mode, prefix string
}

func packetRuntimeMetrics() []metrics.Sample {
	names := []string{"/sched/gomaxprocs:threads", "/sched/latencies:seconds", "/gc/pauses:seconds", "/gc/heap/goal:bytes", "/gc/scan/heap:bytes", "/memory/classes/heap/objects:bytes", "/memory/classes/heap/released:bytes", "/sync/mutex/wait/total:seconds", "/cpu/classes/gc/total:cpu-seconds"}
	samples := make([]metrics.Sample, len(names))
	for i, name := range names {
		samples[i].Name = name
	}
	metrics.Read(samples)
	// Histograms are runtime-owned and overwritten by the next Read.
	for i := range samples {
		if samples[i].Value.Kind() == metrics.KindFloat64Histogram {
			h := samples[i].Value.Float64Histogram()
			h.Counts = slices.Clone(h.Counts)
			h.Buckets = slices.Clone(h.Buckets)
		}
	}
	return samples
}

func (o *packetProcessObservation) begin(t *testing.T, spec packetProcessSpec, phase string) {
	t.Helper()
	o.mode = spec.Observe
	o.prefix = filepath.Join(spec.Output, spec.Role+"-"+phase)
	newFile := func(suffix string) *os.File {
		//nolint:gosec // Profile paths derive from the local artifact directory and fixed suffixes.
		file, err := os.Create(o.prefix + suffix)
		if err != nil {
			t.Fatal(err)
		}
		o.files = append(o.files, file)
		return file
	}
	runtime.GC()
	if o.mode == "cpu" || o.mode == "wait" || o.mode == "trace" {
		for _, name := range []string{"allocs", "block", "mutex"} {
			if name != "allocs" && o.mode != "wait" {
				continue
			}
			file := newFile("." + name + "-before")
			if err := pprof.Lookup(name).WriteTo(file, 0); err != nil {
				t.Fatal(err)
			}
		}
	}
	switch o.mode {
	case "":
	case "decisions", "wire":
	case "cpu":
		if err := pprof.StartCPUProfile(newFile(".cpu")); err != nil {
			t.Fatal(err)
		}
	case "wait":
		runtime.SetBlockProfileRate(1000)
		runtime.SetMutexProfileFraction(5)
	case "trace":
		if err := trace.Start(newFile(".trace")); err != nil {
			t.Fatal(err)
		}
	default:
		t.Fatalf("unknown observation mode %q", o.mode)
	}
	o.reset(t)
}

func (o *packetProcessObservation) reset(t *testing.T) {
	t.Helper()
	o.resources = readPacketProcessResources(t)
	o.kernelBefore = readPacketKernelUDP(t)
	o.runtime = packetRuntimeMetrics()
	if o.wire != nil {
		o.wireBefore = o.wire.snapshot()
	}
	runtime.ReadMemStats(&o.before)
	o.elapsed = startPacketProcessTimer(t)
}

func (o *packetProcessObservation) finish(t *testing.T, result packetProcessStats) packetProcessStats {
	t.Helper()
	result.Elapsed = o.elapsed()
	if o.wire != nil {
		result.Wire = o.wire.snapshot().delta(o.wireBefore)
	}
	var after runtime.MemStats
	runtime.ReadMemStats(&after)
	result.Allocated, result.Allocs = after.TotalAlloc-o.before.TotalAlloc, after.Mallocs-o.before.Mallocs
	result.Heap, result.HeapObjects, result.Stack, result.Sys = after.HeapAlloc, after.HeapObjects, after.StackInuse, after.Sys
	result.GCs, result.GCPause = after.NumGC-o.before.NumGC, time.Duration(after.PauseTotalNs-o.before.PauseTotalNs)
	result.Goroutines = runtime.NumGoroutine()
	result.Resources = readPacketProcessResources(t).delta(o.resources)
	result.KernelUDP = readPacketKernelUDP(t)
	for name, before := range o.kernelBefore {
		if result.KernelUDP[name] >= before {
			result.KernelUDP[name] -= before
		}
	}
	result.Runtime = make(map[string]any)
	for i, sample := range packetRuntimeMetrics() {
		switch sample.Value.Kind() {
		case metrics.KindUint64:
			result.Runtime[sample.Name] = sample.Value.Uint64()
		case metrics.KindFloat64:
			result.Runtime[sample.Name] = sample.Value.Float64() - o.runtime[i].Value.Float64()
		case metrics.KindFloat64Histogram:
			h, before := sample.Value.Float64Histogram(), o.runtime[i].Value.Float64Histogram()
			counts := slices.Clone(h.Counts)
			var total uint64
			for j := range counts {
				counts[j] -= before.Counts[j]
				total += counts[j]
			}
			quantiles := make(map[string]any)
			quantiles["samples"] = total
			for _, percentile := range []int{50, 95, 99} {
				var cumulative uint64
				for j, count := range counts {
					cumulative += count
					if total > 0 && cumulative*100 >= total*uint64(percentile) {
						// The last histogram bucket has an infinite upper bound.
						quantiles[strconv.Itoa(percentile)] = h.Buckets[j]
						break
					}
				}
			}
			result.Runtime[sample.Name] = quantiles
		}
	}
	if o.mode == "cpu" {
		pprof.StopCPUProfile()
	}
	if o.mode == "trace" {
		trace.Stop()
	}
	if o.mode == "wait" {
		runtime.SetBlockProfileRate(0)
		runtime.SetMutexProfileFraction(0)
	}
	if o.mode == "cpu" || o.mode == "wait" || o.mode == "trace" {
		runtime.GC()
		for _, name := range []string{"heap", "allocs", "block", "mutex", "goroutine"} {
			if (name == "block" || name == "mutex") && o.mode != "wait" {
				continue
			}
			//nolint:gosec // Profile names are fixed and the artifact directory is operator-owned.
			file, err := os.Create(o.prefix + "." + name)
			if err != nil {
				t.Fatal(err)
			}
			err = pprof.Lookup(name).WriteTo(file, 0)
			closeErr := file.Close()
			if err != nil || closeErr != nil {
				t.Fatalf("profile %s: %v %v", name, err, closeErr)
			}
		}
	}
	for _, file := range o.files {
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}
	o.files = nil
	return result
}
