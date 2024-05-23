package endpoint

import (
	"context"
	"testing"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
)

func TestEndpointSensor(t *testing.T) {
	expected := 5.3

	sensor := Sensor{}
	sensor.Name = "endpoint"
	sensor.Type = "endpoint"
	sensor.RefreshRate = time.Second
	err := sensor.Init(context.TODO())
	if err != nil {
		t.Errorf("error initializing endpoint: %s", err)
	}
	sensor.Set(expected)
	err = sensors.DoTestSensor(t, &sensor, expected)
	if err != nil {
		t.Errorf("error testing endpoint: %s", err)
	}
}
