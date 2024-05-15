package shell

import (
	"fmt"
	"testing"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
)

func TestRawShellSensor(t *testing.T) {
	expected := 5.3

	var err error
	sensor := Sensor{}
	sensor.Name = "test_shell_raw"
	sensor.Type = "redis"
	sensor.Command = fmt.Sprintf("/bin/echo %.2f", expected)
	sensor.JsonPath = ""
	sensor.RefreshRate = time.Second
	sensor.Description = "test raw shell sensor"
	sensor.Link = "http://example.org"
	sensor.Minimum = 0
	sensor.Maximum = 10

	err = sensors.DoTestSensor(t, &sensor, expected)
	if err != nil {
		t.Errorf("error executing test: %s", err)
	}
}

func TestJsonShellSensor(t *testing.T) {
	var err error
	sensor := Sensor{}
	sensor.Name = "test_shell_raw"
	sensor.Type = "redis"
	sensor.Command = `/bin/echo {"a":5.3}`
	sensor.JsonPath = "$.a"
	sensor.RefreshRate = time.Second
	sensor.Description = "test raw shell sensor"
	sensor.Link = "http://example.org"
	sensor.Minimum = 0
	sensor.Maximum = 10

	err = sensors.DoTestSensor(t, &sensor, 5.3)
	if err != nil {
		t.Errorf("error executing test: %s", err)
	}
}
