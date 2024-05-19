package redis

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	Client    *redis.Client
	val       float64
	updatedAt time.Time
}

func (s *Sensor) Init(ctx context.Context) error {
	return s.Ping(ctx)
}

func (s *Sensor) Ping(ctx context.Context) error {
	return s.Client.Ping(ctx).Err()
}

func (s *Sensor) Close(ctx context.Context) error {
	err := s.Client.Close()
	if err != nil {
		if errors.Is(err, redis.ErrClosed) {
			return nil
		}
	}
	return err
}

func (s *Sensor) Value() float64 {
	return s.val
}

func (s *Sensor) Update(ctx context.Context, _ float64) error {
	s.updatedAt = time.Now()
	args := strings.Split(s.Query, " ")
	b := make([]interface{}, len(args))
	for i := range args {
		b[i] = args[i]
	}
	val, err := s.Client.Do(ctx, b...).Float64()
	if err != nil {
		return err
	}
	s.val = val
	s.updatedAt = time.Now()
	return nil
}

func (s *Sensor) UpdatedAt() time.Time {
	return s.updatedAt
}
