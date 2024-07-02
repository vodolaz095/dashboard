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
	Client *redis.Client

	// DatabaseConnectionName is used to dial database
	DatabaseConnectionName string `yaml:"database_connection_name"`
	// Query is send to remote database in order to receive data from it
	Query string `yaml:"query"`
}

func (s *SyncSensor) Init(ctx context.Context) error {
	s.Mutex = &sync.RWMutex{}
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
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.UpdatedAt = time.Now()
	args := strings.Split(s.Query, " ")
	b := make([]interface{}, len(args))
	for i := range args {
		b[i] = args[i]
	}
	val, err := s.Client.Do(ctx, b...).Float64()
	if err != nil {
		s.Value = 0
		s.Error = err
		return err
	}
	s.Value = val
	s.Error = nil
	return nil
}
