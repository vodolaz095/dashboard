package sql

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	Con *sql.Conn

	// DatabaseConnectionName is used to dial database
	DatabaseConnectionName string `yaml:"database_connection_name"`
	// Query is send to remote database in order to receive data from it
	Query string `yaml:"query"`
}

func (s *Sensor) Init(ctx context.Context) error {
	if s.A == 0 {
		s.A = 1
	}
	s.Mutex = &sync.RWMutex{}
	return s.Con.PingContext(ctx)
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
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.UpdatedAt = time.Now()
	s.Value = 0
	var val float64
	err = s.Con.QueryRowContext(ctx, s.Query).Scan(&val)
	if err != nil {
		s.Error = err
		return err
	}
	s.Value = val
	s.Error = nil
	return nil
}
