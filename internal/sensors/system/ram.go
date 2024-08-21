package system

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

// /proc/meminfo
// MemFree:         4337732 kB

type FreeRAMSensor struct {
	sensors.UnimplementedSensor
}

func (frs *FreeRAMSensor) Init(ctx context.Context) error {
	frs.Mutex = &sync.RWMutex{}
	return nil
}

func (frs *FreeRAMSensor) Ping(ctx context.Context) error {
	return nil
}

func (frs *FreeRAMSensor) Close(ctx context.Context) error {
	return nil
}

func (frs *FreeRAMSensor) Update(ctx context.Context) (err error) {
	frs.Mutex.Lock()
	defer frs.Mutex.Unlock()
	frs.UpdatedAt = time.Now()
	frs.Value = 0
	if runtime.GOOS != "linux" {
		err = fmt.Errorf("not implemented on %s", runtime.GOOS)
		frs.Error = err
		return err
	}
	raw, err := os.OpenFile("/proc/meminfo", os.O_RDONLY, 0444)
	if err != nil {
		frs.Error = fmt.Errorf("error opening /proc/meminfo: %w", err)
		return err
	}
	defer func() {
		err1 := raw.Close()
		if err1 != nil {
			frs.Error = fmt.Errorf("error closing /proc/meminfo: %w", err1)
		}
	}()
	var line string
	var val float64
	scanner := bufio.NewScanner(raw)
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
			return
		}
		frs.Value = val / 1024 // MBytes!
		frs.Error = nil
		break
	}
	return
}
