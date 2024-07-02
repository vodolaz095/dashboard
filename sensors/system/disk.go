package system

import (
	"context"
	"sync"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
	"golang.org/x/sys/unix"
)

// diskSpaceSensor is a base class for UsedDiskSpaceSensor, FreeDiskSpaceSensor, FreeDiskSpaceRatioSensor
type diskSpaceSensor struct {
	sensors.UnimplementedSensor
	Path      string
	UsedSpase float64
	FreeSpace float64
	Ratio     float64
}

func (ds *diskSpaceSensor) Init(ctx context.Context) error {
	ds.Mutex = &sync.RWMutex{}
	if ds.Tags == nil {
		ds.Tags = make(map[string]string, 0)
	}
	_, ok := ds.Tags["mount_point"]
	if !ok {
		ds.Tags["mount_point"] = ds.Path
	}
	return nil
}

func (ds *diskSpaceSensor) Ping(ctx context.Context) error {
	return nil
}

func (ds *diskSpaceSensor) Close(ctx context.Context) error {
	return nil
}

func (ds *diskSpaceSensor) Update(ctx context.Context) (err error) {
	ds.Mutex.Lock()
	defer ds.Mutex.Unlock()
	ds.UpdatedAt = time.Now()
	var stat unix.Statfs_t
	err = unix.Statfs(ds.Path, &stat)
	if err != nil {
		ds.Error = err
		return err
	}
	ds.FreeSpace = float64(stat.Bavail*uint64(stat.Bsize)) / 1024 / 1024
	ds.UsedSpase = float64((stat.Blocks-stat.Bavail)*uint64(stat.Bsize)) / 1024 / 1024
	ds.Ratio = 100 * ds.FreeSpace / ds.UsedSpase
	return nil
}

type FreeDiskSpaceSensor struct {
	diskSpaceSensor
}

func (fdss *FreeDiskSpaceSensor) GetValue() float64 {
	fdss.Mutex.RLock()
	defer fdss.Mutex.RUnlock()
	return fdss.FreeSpace
}

type UsedDiskSpaceSensor struct {
	diskSpaceSensor
}

func (udss *UsedDiskSpaceSensor) GetValue() float64 {
	udss.Mutex.RLock()
	defer udss.Mutex.RUnlock()
	return udss.UsedSpase
}

type FreeDiskSpaceRatioSensor struct {
	diskSpaceSensor
}

func (sensor *FreeDiskSpaceRatioSensor) GetValue() float64 {
	sensor.Mutex.RLock()
	defer sensor.Mutex.RUnlock()
	return sensor.Ratio
}
