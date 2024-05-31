package redis

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vodolaz095/dashboard/sensors"
)

// SyncSensor executes synchronous query against redis database, for example,
// `GET A`, `ZCOUNT something 0 10000`, `LLEN some_list` and so on
type SyncSensor struct {
	sensors.UnimplementedSensor
	mu     *sync.Mutex
	Client *redis.Client
}

func (s *SyncSensor) Init(ctx context.Context) error {
	s.mu = &sync.Mutex{}
	if s.A == 0 {
		s.A = 1
	}
	return s.Ping(ctx)
}

func (s *SyncSensor) Ping(ctx context.Context) error {
	return s.Client.Ping(ctx).Err()
}

func (s *SyncSensor) Close(ctx context.Context) error {
	err := s.Client.Close()
	if err != nil {
		if errors.Is(err, redis.ErrClosed) {
			return nil
		}
	}
	return err
}

func (s *SyncSensor) Update(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.UpdatedAt = time.Now()
	args := strings.Split(s.Query, " ")
	b := make([]interface{}, len(args))
	for i := range args {
		b[i] = args[i]
	}
	val, err := s.Client.Do(ctx, b...).Float64()
	if err != nil {
		s.Error = err
		return err
	}
	s.Value = val
	return nil
}
