package mysql

import (
	"os"
	"testing"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
)

func TestMysqlSensor(t *testing.T) {
	mysqlConnectionString := os.Getenv("MYSQL_URL")
	if mysqlConnectionString == "" {
		mysqlConnectionString = "root:dashboard@/dashboard"
	}
	expected := 5.3

	sensor := Sensor{}
	sensor.Name = "test_mysql"
	sensor.Type = "mysql"
	sensor.DatabaseConnectionString = mysqlConnectionString
	sensor.Query = "SELECT 3+2.3"
	sensor.RefreshRate = time.Second
	sensor.Description = "mysql/mariadb sensor"
	sensor.Link = "https://www.mysql.com/"
	sensor.Minimum = 0
	sensor.Maximum = 10

	err := sensors.DoTestSensor(t, &sensor, expected)
	if err != nil {
		t.Errorf("error testing mysql/mariadb via %s: %s", mysqlConnectionString, err)
	}
}
