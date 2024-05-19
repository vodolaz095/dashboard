package mysql

import (
	"context"

	_ "github.com/go-sql-driver/mysql"
	sqlSensor "github.com/vodolaz095/dashboard/sensors/sql"
)

type Sensor struct {
	sqlSensor.Sensor
}

func (s *Sensor) Init(ctx context.Context) error {
	return s.Con.PingContext(ctx)
}
