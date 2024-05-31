package postgres

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	sqlSensor "github.com/vodolaz095/dashboard/sensors/sql"
)

type Sensor struct {
	sqlSensor.Sensor
}

func (s *Sensor) Init(ctx context.Context) error {
	if s.A == 0 {
		s.A = 1
	}
	return s.Con.PingContext(ctx)
}
