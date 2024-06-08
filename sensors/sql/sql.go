package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor

	mu  *sync.Mutex
	Con *sql.Conn
}

func (s *Sensor) Init(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (s *Sensor) Ping(ctx context.Context) error {
	return s.Con.PingContext(ctx)
}

func (s *Sensor) Close(ctx context.Context) error {
	err := s.Con.Close()
	if err != nil {
		if errors.Is(err, sql.ErrConnDone) {
			return nil
		}
	}
	return err
}

func (s *Sensor) Update(ctx context.Context) (err error) {
	if s.mu == nil {
		s.mu = &sync.Mutex{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.UpdatedAt = time.Now()
	var val float64
	err = s.Con.QueryRowContext(ctx, s.Query).Scan(&val)
	if err != nil {
		s.Error = err
		return err
	}
	s.Value = val
	return nil
}
