package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/vodolaz095/dashboard/sensors"
	"github.com/vodolaz095/dqueue"
)

type testSensor struct {
	sensors.UnimplementedSensor
	T      *testing.T
	Alive  bool
	value  float64
	Closed bool
}

func (ts *testSensor) Init(_ context.Context) error {
	return nil
}

func (ts *testSensor) Ping(_ context.Context) error {
	if ts.Alive {
		return nil
	}
	return errors.New("sensor is dead")
}

func (ts *testSensor) Close(_ context.Context) error {
	ts.Closed = true
	return nil
}

func (ts *testSensor) Value() float64 {
	return ts.value
}

func (ts *testSensor) UpdatedAt() time.Time {
	return time.Now()
}

func (ts *testSensor) Update(_ context.Context, val float64) error {
	ts.T.Logf("Updating sensor to %v", val)
	ts.value = val
	return nil
}

func TestSensorsService(t *testing.T) {
	ts := testSensor{}
	ts.T = t
	q := dqueue.New()
	service := SensorsService{
		ListOfSensors: []string{"testSensor"},
		Sensors: map[string]sensors.ISensor{
			"testSensor": &ts,
		},
		UpdateInterval: 100 * time.Millisecond,
		UpdateQueue:    &q,
	}
	err := service.Ping(context.TODO())
	if err != nil {
		if err.Error() != "sensor is dead" {
			t.Error("wrong error")
			return
		}
	} else {
		t.Error("sensor is not alive")
		return
	}
	ts.Alive = true
	err = service.Ping(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}

}
