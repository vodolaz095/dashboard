package file

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

func TestFileSensorRaw(t *testing.T) {
	expected := 15.3

	sensor := Sensor{}
	sensor.Name = "file_raw"
	sensor.Type = "file"
	sensor.RefreshRate = time.Second
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	sensor.PathToReadingsFile = filepath.Join(dir, "data", "raw.txt")

	err := sensor.Init(context.TODO())
	if err != nil {
		t.Errorf("error initializing endpoint: %s", err)
	}
	err = sensors.DoTestSensor(t, &sensor, expected)
	if err != nil {
		t.Errorf("error testing endpoint: %s", err)
	}
}

func TestFileSensorJSON(t *testing.T) {
	expected := 17.9

	sensor := Sensor{}
	sensor.Name = "file_json"
	sensor.Type = "file"
	sensor.RefreshRate = time.Second
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	sensor.PathToReadingsFile = filepath.Join(dir, "data", "json.json")
	sensor.JsonPath = "@.a.value"

	err := sensor.Init(context.TODO())
	if err != nil {
		t.Errorf("error initializing endpoint: %s", err)
	}
	err = sensors.DoTestSensor(t, &sensor, expected)
	if err != nil {
		t.Errorf("error testing endpoint: %s", err)
	}
}
