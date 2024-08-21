package endpoint

import (
	"context"
	"sync"
	"time"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor

	// Token is Bearer strategy token used to send metrics for endpoint sensor
	Token string `json:"token"`
}

func (s *Sensor) Init(ctx context.Context) error {
	s.Mutex = &sync.RWMutex{}
	if s.A == 0 {
		s.A = 1
	}
	return nil
}

func (s *Sensor) Ping(ctx context.Context) error {
	return nil
}

func (s *Sensor) Close(ctx context.Context) error {
	return nil
}

func (s *Sensor) Update(_ context.Context) error {
	return nil
}

func (s *Sensor) Set(newVal float64) {
	s.Mutex.Lock()
	s.Value = newVal
	s.UpdatedAt = time.Now()
	s.Mutex.Unlock()
}

func (s *Sensor) Increment(delta float64) {
	s.Mutex.Lock()
	s.Value = s.Value + delta
	s.UpdatedAt = time.Now()
	s.Mutex.Unlock()
}
