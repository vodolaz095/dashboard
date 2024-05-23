package redis

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vodolaz095/dashboard/sensors"
)

func TestRedisSensor(t *testing.T) {
	redisConnectionString := os.Getenv("REDIS_URL")
	if redisConnectionString == "" {
		redisConnectionString = "redis://localhost:6379"
	}
	expected := 5.3

	var err error
	sensor := Sensor{}
	sensor.Name = "test_redis"
	sensor.Type = "redis"
	sensor.Query = "get a"
	sensor.RefreshRate = time.Second
	sensor.Description = "test redis sensor"
	sensor.Link = "http://redis.io/"
	sensor.Minimum = 0
	sensor.Maximum = 10

	opts, err := redis.ParseURL(redisConnectionString)
	if err != nil {
		t.Errorf("error parsing redis connection string: %s", err)
		return
	}
	client := redis.NewClient(opts)
	err = client.Set(context.Background(), "a", fmt.Sprintf("%.2f", expected), time.Second).Err()
	if err != nil {
		t.Errorf("error setting redis key: %s", err)
		return
	}

	sensor.Client = client
	err = sensors.DoTestSensor(t, &sensor, 5.3)
	if err != nil {
		t.Errorf("error executing test: %s", err)
	}
}
