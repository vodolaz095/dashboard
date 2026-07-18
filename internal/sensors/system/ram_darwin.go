//go:build darwin
// +build darwin

package system

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func (frs *FreeRAMSensor) Update(ctx context.Context) (err error) {
	frs.Mutex.Lock()
	defer frs.Mutex.Unlock()
	frs.UpdatedAt = time.Now()
	frs.Value = 0

	cmd := exec.Command("sysctl", "-n", "hw.pagesize")
	output, err := cmd.Output()
	if err != nil {
		frs.Error = fmt.Errorf("error getting page size: %w", err)
		return err
	}

	pageSize, err := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		frs.Error = fmt.Errorf("error parsing page size: %w", err)
		return err
	}

	cmd = exec.Command("vm_stat")
	output, err = cmd.Output()
	if err != nil {
		frs.Error = fmt.Errorf("error running vm_stat: %w", err)
		return err
	}

	var freePages, inactivePages, speculativePages uint64

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "Pages free:"):
			freePages = parsePagesValue(line)
		case strings.HasPrefix(line, "Pages inactive:"):
			inactivePages = parsePagesValue(line)
		case strings.HasPrefix(line, "Pages speculative:"):
			speculativePages = parsePagesValue(line)
		}
	}

	freeBytes := (freePages + inactivePages + speculativePages) * pageSize
	frs.Value = float64(freeBytes) / (1024 * 1024) // MBytes!
	frs.Error = nil

	return nil
}

func parsePagesValue(line string) uint64 {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		return 0
	}
	valueStr := strings.TrimRight(parts[3], ".")
	value, err := strconv.ParseUint(valueStr, 10, 64)
	if err != nil {
		return 0
	}
	return value
}
