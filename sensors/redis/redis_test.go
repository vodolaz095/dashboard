package redis

import (
	"context"
	"testing"
	"time"
)

func TestSensor(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	sensor := Sensor{}
	sensor.Name = "test_redis"
	sensor.Type = "redis"
	sensor.DatabaseConnectionString = "redis://localhost:6379"
	sensor.Query = "get a"
	sensor.RefreshRate = time.Second
	sensor.Description = "test redis sensor"
	sensor.Link = "http://redis.io/"
	sensor.Minimum = 0
	sensor.Maximum = 10
	err = sensor.Init(ctx)
	if err != nil {
		t.Errorf("error initializing: %s", err)
	}
	err = sensor.Update(ctx, 0)
	if err != nil {
		t.Errorf("error updating value: %s", err)
	}
	t.Logf("Value: %f", sensor.Value())
	err = sensor.Close(ctx)
	if err != nil {
		t.Errorf("error closing value: %s", err)
	}
}
