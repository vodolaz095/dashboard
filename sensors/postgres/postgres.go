package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	sqlSensor "github.com/vodolaz095/dashboard/sensors/sql"
)

type Sensor struct {
	sqlSensor.Sensor
}

func (s *Sensor) Init(ctx context.Context) error {
	db, err := sql.Open("pgx", s.DatabaseConnectionString)
	if err != nil {
		return err
	}
	con, err := db.Conn(ctx)
	if err != nil {
		return err
	}
	s.Con = con
	return nil
}
