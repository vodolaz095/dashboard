//go:build linux
// +build linux

package system

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func (frs *FreeRAMSensor) Update(context.Context) (err error) {
	frs.Mutex.Lock()
	defer frs.Mutex.Unlock()
	frs.UpdatedAt = time.Now()
	frs.Value = 0

	raw, err := os.OpenFile("/proc/meminfo", os.O_RDONLY, 0444)
	if err != nil {
		frs.Error = fmt.Errorf("error opening /proc/meminfo: %w", err)
		return err
	}

	defer func() {
		closeErr := raw.Close()
		if closeErr != nil {
			if frs.Error == nil {
				frs.Error = fmt.Errorf("error closing /proc/meminfo: %w", closeErr)
			}
			if err == nil {
				err = closeErr
			}
		}
	}()

	var line string
	var val float64
	scanner := bufio.NewScanner(raw)
	found := false

	for scanner.Scan() {
		line = scanner.Text()
		if !strings.HasPrefix(line, "MemFree:") {
			continue
		}
		line = strings.TrimPrefix(line, "MemFree:")
		line = strings.TrimSuffix(line, "kB")
		line = strings.TrimSpace(line)

		val, err = strconv.ParseFloat(line, 64)
		if err != nil {
			frs.Error = err
			return err
		}
		frs.Value = val / 1024 // MBytes!
		frs.Error = nil
		found = true
		break
	}

	if scanErr := scanner.Err(); scanErr != nil {
		frs.Error = fmt.Errorf("error scanning /proc/meminfo: %w", scanErr)
		return frs.Error
	}

	if !found {
		frs.Error = fmt.Errorf("MemFree not found in /proc/meminfo")
		return frs.Error
	}

	return err
}
