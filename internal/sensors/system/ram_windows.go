//go:build windows
// +build windows

package system

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type memoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

var (
	kernel32             *windows.LazyDLL
	globalMemoryStatusEx *windows.LazyProc
)

func (frs *FreeRAMSensor) Update(ctx context.Context) (err error) {
	frs.Mutex.Lock()
	defer frs.Mutex.Unlock()
	frs.UpdatedAt = time.Now()
	frs.Value = 0

	if kernel32 == nil {
		kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	}
	if globalMemoryStatusEx == nil {
		globalMemoryStatusEx = kernel32.NewProc("GlobalMemoryStatusEx")
	}

	var memStat memoryStatusEx
	memStat.Length = uint32(unsafe.Sizeof(memStat))

	ret, _, err := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStat)))
	if ret == 0 {
		frs.Error = fmt.Errorf("error calling GlobalMemoryStatusEx: %w", err)
		return frs.Error
	}

	frs.Value = float64(memStat.AvailPhys) / (1024 * 1024) // MBytes!
	frs.Error = nil

	return nil
}
