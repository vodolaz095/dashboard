package system

import (
	"context"
	"sync"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

type FreeRAMSensor struct {
	sensors.UnimplementedSensor
}

func (frs *FreeRAMSensor) Init(context.Context) error {
	frs.Mutex = &sync.RWMutex{}
	return nil
}

func (frs *FreeRAMSensor) Ping(context.Context) error {
	return nil
}

func (frs *FreeRAMSensor) Close(context.Context) error {
	return nil
}
