package system

import (
	"context"
	"sync"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

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
