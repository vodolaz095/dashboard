package postgres

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
)

func TestPostgresSensor(t *testing.T) {
	pgConnectionString := os.Getenv("PG_URL")
	if pgConnectionString == "" {
		pgConnectionString = "postgres://dashboard:dashboard@localhost:5432/dashboard"
	}
	expected := 5.3

	db, err := sql.Open("pgx", pgConnectionString)
	if err != nil {
		t.Errorf("error dialing pg via %s: %s", pgConnectionString, err)
		return
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)
	con, err := db.Conn(context.TODO())
	if err != nil {
		t.Errorf("error instantiating connection: %s", err)
		return
	}

	sensor := Sensor{}
	sensor.Name = "test_pg"
	sensor.Type = "postgres"
	sensor.Con = con
	sensor.Query = "SELECT 3+2.3"
	sensor.RefreshRate = time.Second
	sensor.Description = "postgres sensor"
	sensor.Link = "https://www.postgresql.org/"
	sensor.Minimum = 0
	sensor.Maximum = 10

	err = sensors.DoTestSensor(t, &sensor, expected)
	if err != nil {
		t.Errorf("error testing postgres via %s: %s", pgConnectionString, err)
	}
}
