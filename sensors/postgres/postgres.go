package postgres

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	sqlSensor "github.com/vodolaz095/dashboard/sensors/sql"
)

type Sensor struct {
	sqlSensor.Sensor
}
