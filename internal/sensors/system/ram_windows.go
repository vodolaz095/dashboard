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

// MemoryStat represents Windows memory statistics
type MemoryStat struct {
	PageFileSize        uint64  // Page file size in bytes
	PageUsagePercentage float64 // Page file usage percentage
	PhysMemSize         uint64  // Physical memory size in bytes
	PhysMemUsage        uint64  // Physical memory usage in bytes
}

// getWindowsRAMStats retrieves RAM statistics on Windows using Windows API
func getWindowsRAMStats() (*MemoryStat, error) {
	// Get system memory information
	var memInfo windows.MemoryBasicInformation
	memInfo.Length = uint32(unsafe.Sizeof(memInfo))

	// Call GlobalMemoryStatusEx to get memory information
	err := windows.GlobalMemoryStatusEx(&memInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory status: %w", err)
	}

	// Calculate physical memory usage in bytes
	physMemSize := memInfo.TotalPhys
	physMemFree := memInfo.AvailPhys
	physMemUsage := physMemSize - physMemFree

	// Calculate page file usage
	pageFileSize := memInfo.TotalPageFile - memInfo.TotalPhys // Page file size beyond physical memory
	pageFileUsage := memInfo.TotalPageFile - memInfo.AvailPageFile
	var pageUsagePercentage float64
	if pageFileSize > 0 {
		pageUsagePercentage = float64(pageFileUsage) / float64(pageFileSize) * 100
	} else {
		pageUsagePercentage = 0
	}

	return &MemoryStat{
		PageFileSize:        pageFileSize,
		PageUsagePercentage: pageUsagePercentage,
		PhysMemSize:         physMemSize,
		PhysMemUsage:        physMemUsage,
	}, nil
}

// Update method for Windows - this would be part of the FreeRAMSensor struct
func (frs *FreeRAMSensor) updateWindows(ctx context.Context) (err error) {
	frs.Mutex.Lock()
	defer frs.Mutex.Unlock()
	frs.UpdatedAt = time.Now()
	frs.Value = 0

	memStat, err := getWindowsRAMStats()
	if err != nil {
		frs.Error = fmt.Errorf("error getting Windows RAM stats: %w", err)
		return frs.Error
	}

	// Convert free physical memory to MB
	freeMB := float64(memStat.PhysMemSize-memStat.PhysMemUsage) / 1024 / 1024
	frs.Value = freeMB
	frs.Error = nil
	return nil
}
