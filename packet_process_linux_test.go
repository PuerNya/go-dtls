package dtls13

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

func readPacketProcessResources(t *testing.T) packetProcessResources {
	t.Helper()
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		t.Fatal(err)
	}
	r := packetProcessResources{User: time.Duration(usage.Utime.Sec)*time.Second + time.Duration(usage.Utime.Usec)*time.Microsecond,
		System:  time.Duration(usage.Stime.Sec)*time.Second + time.Duration(usage.Stime.Usec)*time.Microsecond,
		PeakRSS: uint64(usage.Maxrss) * 1024, PageFaults: uint64(usage.Minflt + usage.Majflt), MinorFaults: uint64(usage.Minflt), MajorFaults: uint64(usage.Majflt), VoluntarySwitches: uint64(usage.Nvcsw), InvoluntarySwitches: uint64(usage.Nivcsw)}
	for _, path := range []string{"/proc/self/status", "/proc/self/io"} {
		//nolint:gosec // Both paths are fixed Linux process counters.
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for line := range strings.SplitSeq(string(data), "\n") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			var target *uint64
			switch fields[0] {
			case "VmRSS:":
				target = &r.RSS
			case "read_bytes:":
				target = &r.IOReadBytes
			case "write_bytes:":
				target = &r.IOWriteBytes
			}
			if target != nil {
				value, err := strconv.ParseUint(fields[1], 10, 64)
				if err != nil {
					t.Fatal(err)
				}
				if fields[0] == "VmRSS:" {
					value *= 1024
				}
				*target = value
			}
		}
	}
	return r
}

func setPacketProcessAffinity(value string) error {
	cpus, err := packetProcessCPUs(value)
	if err != nil || len(cpus) == 0 {
		return err
	}
	var mask [128]byte
	for _, cpu := range cpus {
		mask[cpu/8] |= 1 << uint(cpu%8)
	}
	// Linux affinity belongs to threads; apply it to every existing runtime thread.
	for range 2 {
		tasks, err := os.ReadDir("/proc/self/task")
		if err != nil {
			return err
		}
		for _, task := range tasks {
			tid, err := strconv.Atoi(task.Name())
			if err != nil {
				return err
			}
			//nolint:gosec // The kernel synchronously reads the fixed 1024-bit CPU mask.
			_, _, errno := syscall.RawSyscall(syscall.SYS_SCHED_SETAFFINITY, uintptr(tid), uintptr(len(mask)), uintptr(unsafe.Pointer(&mask[0])))
			if errno != 0 && errno != syscall.ESRCH {
				return fmt.Errorf("set runtime thread CPU affinity: %w", errno)
			}
		}
	}
	return nil
}

func readPacketKernelUDP(t *testing.T) map[string]uint64 {
	t.Helper()
	data, err := os.ReadFile("/proc/self/net/snmp")
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	result := make(map[string]uint64)
	for line := range strings.SplitSeq(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 0 || fields[0] != "Udp:" {
			continue
		}
		if names == nil {
			names = fields
			continue
		}
		if len(names) != len(fields) {
			t.Fatal("invalid kernel UDP statistics")
		}
		for i, name := range names[1:] {
			value, err := strconv.ParseUint(fields[i+1], 10, 64)
			if err != nil {
				t.Fatal(err)
			}
			result["namespace_"+name] = value
		}
	}
	return result
}
