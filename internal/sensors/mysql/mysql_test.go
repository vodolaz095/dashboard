package mysql

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

func TestMysqlSensor(t *testing.T) {
	mysqlConnectionString := os.Getenv("MYSQL_URL")
	if mysqlConnectionString == "" {
		mysqlConnectionString = "root:dashboard@/dashboard"
	}
	expected := 5.3

	db, err := sql.Open("mysql", mysqlConnectionString)
	if err != nil {
		t.Errorf("error dialing mysql via %s: %s", mysqlConnectionString, err)
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
	sensor.Name = "test_mysql"
	sensor.Type = "mysql"
	sensor.Con = con
	sensor.Query = "SELECT 3+2.3"
	sensor.RefreshRate = time.Second
	sensor.Description = "mysql/mariadb sensor"
	sensor.Link = "https://www.mysql.com/"
	sensor.Minimum = 0
	sensor.Maximum = 10

	err = sensors.DoTestSensor(t, &sensor, expected)
	if err != nil {
		t.Errorf("error testing mysql/mariadb via %s: %s", mysqlConnectionString, err)
	}
}
