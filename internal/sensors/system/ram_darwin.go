//go:build darwin
// +build darwin

package system

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// getDarwinRAMStats executes vm_stat and parses its output to get memory statistics
func getDarwinRAMStats() (freeMB float64, totalMB float64, err error) {
	// Execute vm_stat command
	output, err := executeCommand("vm_stat")
	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute vm_stat: %w", err)
	}

	// Parse vm_stat output
	var pageSize uint64
	var freePages uint64
	var activePages uint64
	var inactivePages uint64
	var wiredPages uint64

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "page size of ") {
			// Extract page size (e.g., "page size of 4096 bytes")
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				pageSize, err = strconv.ParseUint(parts[3], 10, 64)
				if err != nil {
					return 0, 0, fmt.Errorf("failed to parse page size: %w", err)
				}
			}
		} else if strings.HasPrefix(line, "Pages free:") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				// Remove trailing dot and parse
				value := strings.TrimSuffix(parts[2], ".")
				freePages, err = strconv.ParseUint(value, 10, 64)
				if err != nil {
					return 0, 0, fmt.Errorf("failed to parse free pages: %w", err)
				}
			}
		} else if strings.HasPrefix(line, "Pages active:") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				// Remove trailing dot and parse
				value := strings.TrimSuffix(parts[2], ".")
				activePages, err = strconv.ParseUint(value, 10, 64)
				if err != nil {
					return 0, 0, fmt.Errorf("failed to parse active pages: %w", err)
				}
			}
		} else if strings.HasPrefix(line, "Pages inactive:") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				// Remove trailing dot and parse
				value := strings.TrimSuffix(parts[2], ".")
				inactivePages, err = strconv.ParseUint(value, 10, 64)
				if err != nil {
					return 0, 0, fmt.Errorf("failed to parse inactive pages: %w", err)
				}
			}
		} else if strings.HasPrefix(line, "Pages wired down:") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				// Remove trailing dot and parse
				value := strings.TrimSuffix(parts[3], ".")
				wiredPages, err = strconv.ParseUint(value, 10, 64)
				if err != nil {
					return 0, 0, fmt.Errorf("failed to parse wired pages: %w", err)
				}
			}
		}
	}

	// Make sure we got a page size
	if pageSize == 0 {
		return 0, 0, fmt.Errorf("could not determine page size")
	}

	// Calculate memory in bytes
	freeBytes := freePages * pageSize
	usedBytes := (activePages + inactivePages + wiredPages) * pageSize
	totalBytes := freeBytes + usedBytes

	// Convert to MB
	freeMB = float64(freeBytes) / 1024 / 1024
	totalMB = float64(totalBytes) / 1024 / 1024

	return freeMB, totalMB, nil
}

// executeCommand is a helper function to execute system commands
// This would need to be implemented in a separate file or imported from a utility package
type commandExecutor func(string) ([]byte, error)

// Update method for macOS - this would be part of the FreeRAMSensor struct
func (frs *FreeRAMSensor) updateDarwin(ctx context.Context) (err error) {
	frs.Mutex.Lock()
	defer frs.Mutex.Unlock()
	frs.UpdatedAt = time.Now()
	frs.Value = 0

	freeMB, totalMB, err := getDarwinRAMStats()
	if err != nil {
		frs.Error = fmt.Errorf("error getting Darwin RAM stats: %w", err)
		return frs.Error
	}

	frs.Value = freeMB
	frs.Error = nil
	return nil
}
