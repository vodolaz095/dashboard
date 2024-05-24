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

func (s *Sensor) GetValue() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.val
}

func (s *Sensor) Update(_ context.Context) error {
	return nil
}

func (s *Sensor) Set(newVal float64) {
	s.mu.Lock()
	s.val = newVal
	s.updatedAt = time.Now()
	s.mu.Unlock()
}

func (s *Sensor) Increment(delta float64) {
	s.mu.Lock()
	s.val = s.val + delta
	s.updatedAt = time.Now()
	s.mu.Unlock()
}

func (s *Sensor) GetUpdatedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updatedAt
}
