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
	return s.Con.PingContext(ctx)
}
