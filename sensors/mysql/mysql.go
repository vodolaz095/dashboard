package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	sqlSensor "github.com/vodolaz095/dashboard/sensors/sql"
)

type Sensor struct {
	sqlSensor.Sensor
}
