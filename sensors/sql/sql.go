package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
)

type Sensor struct {
	sensors.UnimplementedSensor
	Con       *sql.Conn
	val       float64
	updatedAt time.Time
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

func (s *Sensor) Value() float64 {
	return s.val
}

func (s *Sensor) Update(ctx context.Context, _ float64) (err error) {
	var val float64
	err = s.Con.QueryRowContext(ctx, s.Query).Scan(&val)
	if err != nil {
		return err
	}
	s.updatedAt = time.Now()
	s.val = val
	return nil
}

func (s *Sensor) UpdatedAt() time.Time {
	return s.updatedAt
}
