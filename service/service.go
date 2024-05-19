package service

import (
	"database/sql"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vodolaz095/dqueue"

	"github.com/vodolaz095/dashboard/model"
	"github.com/vodolaz095/dashboard/sensors"
)

const DefaultSensorTimeout = 5 * time.Second

type SensorsService struct {
	ListOfSensors  []string
	Sensors        map[string]sensors.ISensor
	UpdateInterval time.Duration
	UpdateQueue    *dqueue.Handler

	subscribers map[string]chan model.Update

	// cached database connections
	MysqlConnections      map[string]*sql.Conn
	PostgresqlConnections map[string]*sql.Conn
	RedisConnections      map[string]*redis.Client
}
