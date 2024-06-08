package service

import (
	"fmt"

	"github.com/vodolaz095/dashboard/config"
	"github.com/vodolaz095/dashboard/sensors"
	"github.com/vodolaz095/dashboard/sensors/curl"
	"github.com/vodolaz095/dashboard/sensors/endpoint"
	"github.com/vodolaz095/dashboard/sensors/file"
	"github.com/vodolaz095/dashboard/sensors/mysql"
	"github.com/vodolaz095/dashboard/sensors/postgres"
	"github.com/vodolaz095/dashboard/sensors/redis"
	"github.com/vodolaz095/dashboard/sensors/shell"
)

func populateBaseSensorParams(sensor *sensors.UnimplementedSensor, params config.Sensor) {
	sensor.Name = params.Name
	sensor.RefreshRate = params.RefreshRate
	sensor.Description = params.Description
	sensor.Link = params.Link
	sensor.Minimum = params.Minimum
	sensor.Maximum = params.Maximum
	sensor.Type = params.Type
	sensor.Tags = params.Tags
	if params.A == 0 {
		sensor.A = 1
	} else {
		sensor.A = params.A
	}
	sensor.B = params.B
}

func (ss *SensorsService) MakeSensor(params config.Sensor) (sensor sensors.ISensor, err error) {
	var connectionIsFound bool
	switch params.Type {
	case "mysql", "mariadb":
		ms := &mysql.Sensor{}
		populateBaseSensorParams(&ms.UnimplementedSensor, params)
		ms.Query = params.Query
		ms.DatabaseConnectionName = params.ConnectionName
		ms.Con, connectionIsFound = ss.MysqlConnections[params.ConnectionName]
		if !connectionIsFound {
			return ms, fmt.Errorf("unknown mysql connection: %s", params.ConnectionName)
		}
		ss.UpdateQueue.ExecuteAfter(ms.Name, DefaultWarmUpDelay)
		return ms, nil

	case "redis":
		rs := &redis.SyncSensor{}
		populateBaseSensorParams(&rs.UnimplementedSensor, params)
		rs.Query = params.Query
		rs.DatabaseConnectionName = params.ConnectionName
		rs.Client, connectionIsFound = ss.RedisConnections[params.ConnectionName]
		if !connectionIsFound {
			return rs, fmt.Errorf("unknown redis connection: %s", params.ConnectionName)
		}
		ss.UpdateQueue.ExecuteAfter(rs.Name, DefaultWarmUpDelay)
		return rs, nil

	case "subscriber":
		rss := &redis.SubscribeSensor{}
		populateBaseSensorParams(&rss.UnimplementedSensor, params)
		rss.ValueOnly = params.ValueOnly
		rss.Channel = params.Channel

		rss.DatabaseConnectionName = params.ConnectionName
		rss.Client, connectionIsFound = ss.RedisConnections[params.ConnectionName]
		if !connectionIsFound {
			return rss, fmt.Errorf("unknown redis connection: %s", params.ConnectionName)
		}
		return rss, nil

	case "postgres":
		ps := &postgres.Sensor{}
		populateBaseSensorParams(&ps.UnimplementedSensor, params)
		ps.Query = params.Query
		ps.DatabaseConnectionName = params.ConnectionName
		ps.Con, connectionIsFound = ss.PostgresqlConnections[params.ConnectionName]
		if !connectionIsFound {
			return ps, fmt.Errorf("unknown postgres connection: %s", params.ConnectionName)

		}
		ss.UpdateQueue.ExecuteAfter(ps.Name, DefaultWarmUpDelay)
		return ps, nil

	case "curl":
		cs := &curl.Sensor{}
		populateBaseSensorParams(&cs.UnimplementedSensor, params)
		cs.HttpMethod = params.HttpMethod
		cs.Endpoint = params.Endpoint
		cs.Headers = params.Headers
		cs.Body = params.Body
		cs.JsonPath = params.JsonPath
		ss.UpdateQueue.ExecuteAfter(cs.Name, DefaultWarmUpDelay)
		return cs, nil

	case "shell":
		shs := &shell.Sensor{}
		populateBaseSensorParams(&shs.UnimplementedSensor, params)
		shs.Command = params.Command
		shs.Environment = params.Environment
		shs.JsonPath = params.JsonPath
		ss.UpdateQueue.ExecuteAfter(shs.Name, DefaultWarmUpDelay)
		return shs, nil

	case "endpoint":
		es := &endpoint.Sensor{}
		populateBaseSensorParams(&es.UnimplementedSensor, params)
		return es, nil

	case "file":
		fs := &file.Sensor{}
		populateBaseSensorParams(&fs.UnimplementedSensor, params)
		fs.PathToReadingsFile = params.PathToReading
		fs.JsonPath = params.JsonPath
		ss.UpdateQueue.ExecuteAfter(fs.Name, DefaultWarmUpDelay)
		return fs, nil

	default:
		return nil, fmt.Errorf("unknown sensor type: %s", params.Type)
	}

}
