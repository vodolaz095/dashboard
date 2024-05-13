package redis

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	client    *redis.Client
	val       float64
	updatedAt time.Time
}

func (s *Sensor) Init(ctx context.Context) error {
	opts, err := redis.ParseURL(s.DatabaseConnectionString)
	if err != nil {
		return err
	}
	s.client = redis.NewClient(opts)
	return s.client.Ping(ctx).Err()
}

func (s *Sensor) Ping(ctx context.Context) error {
	return s.client.Ping(ctx).Err()
}

func (s *Sensor) Close(ctx context.Context) error {
	return s.client.Close()
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
	val, err := s.client.Do(ctx, b...).Float64()
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
