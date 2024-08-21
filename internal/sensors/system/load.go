package system

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

// Code is partially based on https://github.com/mikoim/go-loadavg

type LoadAverageSensor struct {
	LoadAverage1     float64
	LoadAverage5     float64
	LoadAverage15    float64
	RunningProcesses int
	TotalProcesses   int
	LastProcessId    int
	sensors.UnimplementedSensor
}

func (lav *LoadAverageSensor) Init(ctx context.Context) error {
	lav.Mutex = &sync.RWMutex{}
	return nil
}

func (lav *LoadAverageSensor) Ping(ctx context.Context) error {
	return nil
}

func (lav *LoadAverageSensor) Close(ctx context.Context) error {
	return nil
}

func (lav *LoadAverageSensor) Update(ctx context.Context) (err error) {
	lav.Mutex.Lock()
	defer lav.Mutex.Unlock()
	lav.UpdatedAt = time.Now()
	lav.Value = 0
	if runtime.GOOS != "linux" {
		err = fmt.Errorf("not implemented on %s", runtime.GOOS)
		lav.Error = err
		return err
	}

	raw, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		lav.Error = fmt.Errorf("error opening /proc/loadavg: %w", err)
		return lav.Error
	}
	_, err = fmt.Sscanf(string(raw), "%f %f %f %d/%d %d",
		&lav.LoadAverage1, &lav.LoadAverage5, &lav.LoadAverage15,
		&lav.RunningProcesses, &lav.TotalProcesses,
		&lav.LastProcessId)
	if err != nil {
		lav.Error = fmt.Errorf("error scanning /proc/loadavg: %w", err)
		return lav.Error
	}
	lav.Error = nil
	return nil
}

type LoadAverage1Sensor struct {
	LoadAverageSensor
}

func (lav1 *LoadAverage1Sensor) GetValue() float64 {
	lav1.Mutex.RLock()
	defer lav1.Mutex.RUnlock()
	return lav1.LoadAverage1
}

type LoadAverage5Sensor struct {
	LoadAverageSensor
}

func (lav5 *LoadAverage5Sensor) GetValue() float64 {
	lav5.Mutex.RLock()
	defer lav5.Mutex.RUnlock()
	return lav5.LoadAverage5
}

type LoadAverage15Sensor struct {
	LoadAverageSensor
}

func (lav15 *LoadAverage15Sensor) GetValue() float64 {
	lav15.Mutex.RLock()
	defer lav15.Mutex.RUnlock()
	return lav15.LoadAverage15
}

type TotalProcessSensor struct {
	LoadAverageSensor
}

func (tps *TotalProcessSensor) GetValue() float64 {
	tps.Mutex.RLock()
	defer tps.Mutex.RUnlock()
	return float64(tps.TotalProcesses)
}
