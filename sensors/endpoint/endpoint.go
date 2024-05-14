package endpoint

import (
	"context"
	"sync"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	mu        *sync.RWMutex
	val       float64
	updatedAt time.Time
}

func (s *Sensor) Init(ctx context.Context) error {
	s.mu = &sync.RWMutex{}
	return nil
}

func (s *Sensor) Ping(ctx context.Context) error {
	return nil
}

func (s *Sensor) Close(ctx context.Context) error {
	return nil
}

func (s *Sensor) Value() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.val
}

func (s *Sensor) Update(ctx context.Context, newVal float64) error {
	s.mu.Lock()
	s.val = newVal
	s.updatedAt = time.Now()
	s.mu.Unlock()
	return nil
}

func (s *Sensor) UpdatedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updatedAt
}
