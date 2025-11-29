//go:build darwin
// +build darwin

package src

import (
	"sync"

	"unsafe"

	"github.com/ebitengine/purego"
	"golang.org/x/sys/unix"
)

type Library struct {
	addr  uintptr
	path  string
	close func()
}

type HostStatisticsFunc func(host uint32, flavor int32, hostInfoOut uintptr, hostInfoOutCnt *uint32) int

type MachHostSelfFunc func() uint32

type vmStatisticsData struct {
	freeCount     uint32
	activeCount   uint32
	inactiveCount uint32
	wireCount     uint32
}

func (lib *Library) Dlsym(symbol string) (uintptr, error) {
	return purego.Dlsym(lib.addr, symbol)
}

func (lib *Library) Close() {
	lib.close()
}

func NewLibrary(path string) (*Library, error) {
	lib, err := purego.Dlopen(path, purego.RTLD_LAZY|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, err
	}

	closeFunc := func() {
		purego.Dlclose(lib)
	}

	return &Library{
		addr:  lib,
		path:  path,
		close: closeFunc,
	}, nil
}

func GetFunc[T any](lib *Library, symbol string) T {
	var fptr T
	purego.RegisterLibFunc(&fptr, lib.addr, symbol)
	return fptr
}

func getHwMemsize() (uint64, error) {
	total, err := unix.SysctlUint64("hw.memsize")
	if err != nil {
		return 0, err
	}
	return total, nil
}

func operatingSystemAvailableMemory(c chan<- error, wg *sync.WaitGroup) {
	machLib, err := NewLibrary(system)
	if wg != nil {
		defer wg.Done()
	}
	defer close(c)
	if err != nil {
		_ = osavailableMemory.Swap(0)
		c <- errNewLibray
		return
	}
	defer machLib.Close()

	hostStatistics := GetFunc[HostStatisticsFunc](machLib, hostStatisticsSym)
	machHostSelf := GetFunc[MachHostSelfFunc](machLib, machHostSelfSym)

	count := uint32(host_vm_info_count)
	var vmstat vmStatisticsData

	status := hostStatistics(machHostSelf(), host_vm_info,
		uintptr(unsafe.Pointer(&vmstat)), &count)

	if status != kern_success {
		_ = osavailableMemory.Swap(0)
		c <- errKernelfail
		return
	}

	pageSizeAddr, _ := machLib.Dlsym("vm_kernel_page_size")
	pageSize := **(**uint64)(unsafe.Pointer(&pageSizeAddr))
	availableCount := vmstat.inactiveCount + vmstat.freeCount
	available := pageSize * uint64(availableCount)
	_ = osavailableMemory.Swap(available)
	c <- nil
	return
}
