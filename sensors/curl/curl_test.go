package curl

import (
	"net/http"
	"testing"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
)

// http://ip-api.com/line/193.41.76.51?fields=lat -> 55.9397
// http://ip-api.com/json/193.41.76.51 -> lat=55.9397

func TestCurlRawSensor(t *testing.T) {
	expected := 55.9397
	sensor := Sensor{}
	sensor.Name = "test_mysql"
	sensor.Type = "curl"
	sensor.Endpoint = "http://ip-api.com/line/193.41.76.51?fields=lat"
	sensor.JsonPath = ""
	sensor.RefreshRate = time.Second
	sensor.Description = "curl sensor"
	sensor.Link = "https://www.mysql.com/"
	sensor.Minimum = 0
	sensor.Maximum = 10

	sensor.Client = http.DefaultClient
	sensor.ExpectedStatusCode = http.StatusOK
	sensor.Headers = make(map[string]string, 0)

	err := sensors.DoTestSensor(t, &sensor, expected)
	if err != nil {
		t.Errorf("error testing curl: %s", err)
	}
}

func TestCurlJSONSensor(t *testing.T) {
	expected := 55.9397
	sensor := Sensor{}
	sensor.Name = "test_mysql"
	sensor.Type = "mysql"
	sensor.Endpoint = "http://ip-api.com/json/193.41.76.51"
	sensor.JsonPath = "@.lat"
	sensor.RefreshRate = time.Second
	sensor.Description = "curl sensor"
	sensor.Link = "https://www.mysql.com/"
	sensor.Minimum = 0
	sensor.Maximum = 10

	sensor.Client = http.DefaultClient
	sensor.Headers = make(map[string]string, 0)
	sensor.ExpectedStatusCode = http.StatusOK

	err := sensors.DoTestSensor(t, &sensor, expected)
	if err != nil {
		t.Errorf("error testing curl: %s", err)
	}
}
