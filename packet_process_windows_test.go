package dtls13

import (
	"fmt"
	"syscall"
	"testing"
	"time"
	"unsafe"
)

var packetKernel32 = syscall.NewLazyDLL("kernel32.dll")
var packetMemoryInfo = syscall.NewLazyDLL("psapi.dll").NewProc("GetProcessMemoryInfo")

func startPacketProcessTimer(t *testing.T) func() time.Duration {
	t.Helper()
	var frequency, start int64
	//nolint:gosec // Win32 synchronously fills the fixed LARGE_INTEGER buffers.
	if ok, _, err := packetKernel32.NewProc("QueryPerformanceFrequency").Call(uintptr(unsafe.Pointer(&frequency))); ok == 0 || frequency <= 0 {
		t.Fatalf("QueryPerformanceFrequency: %v", err)
	}
	counter := packetKernel32.NewProc("QueryPerformanceCounter")
	//nolint:gosec // Win32 synchronously fills the fixed LARGE_INTEGER buffer.
	if ok, _, err := counter.Call(uintptr(unsafe.Pointer(&start))); ok == 0 {
		t.Fatal(err)
	}
	return func() time.Duration {
		var end int64
		//nolint:gosec // Win32 synchronously fills the fixed LARGE_INTEGER buffer.
		if ok, _, err := counter.Call(uintptr(unsafe.Pointer(&end))); ok == 0 {
			t.Fatal(err)
		}
		ticks := end - start
		return time.Duration(ticks/frequency)*time.Second + time.Duration(ticks%frequency)*time.Second/time.Duration(frequency)
	}
}

func readPacketProcessResources(t *testing.T) packetProcessResources {
	t.Helper()
	var creation, exit, kernel, user syscall.Filetime
	handle := ^uintptr(0)
	if err := syscall.GetProcessTimes(syscall.Handle(handle), &creation, &exit, &kernel, &user); err != nil {
		t.Fatal(err)
	}

	duration := func(value syscall.Filetime) time.Duration {
		return time.Duration(uint64(value.HighDateTime)<<32|uint64(value.LowDateTime)) * 100
	}
	var memory struct {
		Size, PageFaultCount                                                                                                                                                                 uint32
		PeakWorkingSetSize, WorkingSetSize, QuotaPeakPagedPoolUsage, QuotaPagedPoolUsage, QuotaPeakNonPagedPoolUsage, QuotaNonPagedPoolUsage, PagefileUsage, PeakPagefileUsage, PrivateUsage uintptr
	}
	memory.Size = uint32(unsafe.Sizeof(memory))
	//nolint:gosec // Win32 fills the correctly sized PROCESS_MEMORY_COUNTERS_EX buffer synchronously.
	if ok, _, err := packetMemoryInfo.Call(handle, uintptr(unsafe.Pointer(&memory)), uintptr(memory.Size)); ok == 0 {
		t.Fatal(err)
	}
	var io struct{ ReadOperations, WriteOperations, OtherOperations, ReadBytes, WriteBytes, OtherBytes uint64 }
	//nolint:gosec // Win32 fills the fixed IO_COUNTERS buffer synchronously.
	if ok, _, err := packetKernel32.NewProc("GetProcessIoCounters").Call(handle, uintptr(unsafe.Pointer(&io))); ok == 0 {
		t.Fatal(err)
	}
	return packetProcessResources{User: duration(user), System: duration(kernel), RSS: uint64(memory.WorkingSetSize), PeakRSS: uint64(memory.PeakWorkingSetSize), PageFaults: uint64(memory.PageFaultCount), IOReadBytes: io.ReadBytes, IOWriteBytes: io.WriteBytes}
}

func setPacketProcessAffinity(value string) error {
	cpus, err := packetProcessCPUs(value)
	if err != nil || len(cpus) == 0 {
		return err
	}
	var mask uintptr
	for _, cpu := range cpus {
		if cpu >= int(unsafe.Sizeof(mask))*8 {
			return fmt.Errorf("CPU %d exceeds the process affinity mask", cpu)
		}
		mask |= 1 << uint(cpu)
	}
	if ok, _, err := packetKernel32.NewProc("SetProcessAffinityMask").Call(^uintptr(0), mask); ok == 0 {
		return fmt.Errorf("set process CPU affinity: %w", err)
	}
	return nil
}

func readPacketKernelUDP(t *testing.T) map[string]uint64 {
	t.Helper()
	var stats [5]uint32
	//nolint:gosec // GetUdpStatisticsEx fills five DWORDs synchronously for IPv4.
	if status, _, _ := syscall.NewLazyDLL("iphlpapi.dll").NewProc("GetUdpStatisticsEx").Call(uintptr(unsafe.Pointer(&stats[0])), syscall.AF_INET); status != 0 {
		t.Fatalf("GetUdpStatisticsEx: %d", status)
	}
	return map[string]uint64{"host_InDatagrams": uint64(stats[0]), "host_NoPorts": uint64(stats[1]), "host_InErrors": uint64(stats[2]), "host_OutDatagrams": uint64(stats[3])}
}
