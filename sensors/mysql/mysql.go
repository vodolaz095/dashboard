package mysql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	sqlSensor "github.com/vodolaz095/dashboard/sensors/sql"
)

type Sensor struct {
	sqlSensor.Sensor
}

func (s *Sensor) Init(ctx context.Context) error {
	db, err := sql.Open("mysql", s.DatabaseConnectionString)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	con, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	s.Con = con
	return nil
}
