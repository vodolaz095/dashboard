package postgres

import (
	_ "github.com/jackc/pgx/v5/stdlib"

	sqlSensor "github.com/vodolaz095/dashboard/internal/sensors/sql"
)

type Sensor struct {
	sqlSensor.Sensor
}
