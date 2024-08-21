package system

import (
	"context"
	"sync"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

// Good read - https://stackoverflow.com/questions/20108520/get-amount-of-free-disk-space-using-go

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
