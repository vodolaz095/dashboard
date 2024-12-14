package service

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/config"
	"github.com/vodolaz095/dashboard/internal/sensors"
	"github.com/vodolaz095/dashboard/internal/sensors/curl"
	"github.com/vodolaz095/dashboard/internal/sensors/endpoint"
	"github.com/vodolaz095/dashboard/internal/sensors/file"
	"github.com/vodolaz095/dashboard/internal/sensors/mysql"
	"github.com/vodolaz095/dashboard/internal/sensors/postgres"
	"github.com/vodolaz095/dashboard/internal/sensors/redis"
	"github.com/vodolaz095/dashboard/internal/sensors/shell"
	"github.com/vodolaz095/dashboard/internal/sensors/system"
	"github.com/vodolaz095/dashboard/internal/sensors/victoriametrics"
)

func populateBaseSensorParams(sensor *sensors.UnimplementedSensor, params config.Sensor) {
	sensor.Name = params.Name
	if params.RefreshRate == 0 {
		sensor.RefreshRate = DefaultRefreshRate
		log.Warn().Msgf("Setting refresh rate for sensor %s to default value of %s!",
			params.Name, DefaultRefreshRate.String(),
		)
	} else {
		sensor.RefreshRate = params.RefreshRate
	}
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
	case "load1":
		l1 := system.LoadAverage1Sensor{}
		populateBaseSensorParams(&l1.UnimplementedSensor, params)
		ss.UpdateQueue.ExecuteAfter(l1.Name, DefaultWarmUpDelay)
		return &l1, nil

	case "load5":
		l5 := system.LoadAverage1Sensor{}
		populateBaseSensorParams(&l5.UnimplementedSensor, params)
		ss.UpdateQueue.ExecuteAfter(l5.Name, DefaultWarmUpDelay)
		return &l5, nil

	case "load15":
		l15 := system.LoadAverage1Sensor{}
		populateBaseSensorParams(&l15.UnimplementedSensor, params)
		ss.UpdateQueue.ExecuteAfter(l15.Name, DefaultWarmUpDelay)
		return &l15, nil

	case "process":
		tps := system.TotalProcessSensor{}
		populateBaseSensorParams(&tps.UnimplementedSensor, params)
		ss.UpdateQueue.ExecuteAfter(tps.Name, DefaultWarmUpDelay)
		return &tps, nil

	case "free_ram":
		frs := system.FreeRAMSensor{}
		populateBaseSensorParams(&frs.UnimplementedSensor, params)
		ss.UpdateQueue.ExecuteAfter(frs.Name, DefaultWarmUpDelay)
		return &frs, nil

	case "used_disk_space":
		uds := system.UsedDiskSpaceSensor{}
		populateBaseSensorParams(&uds.UnimplementedSensor, params)
		uds.Path = params.PathToMountPoint
		ss.UpdateQueue.ExecuteAfter(uds.Name, DefaultWarmUpDelay)
		return &uds, nil

	case "free_disk_space":
		fds := system.FreeDiskSpaceSensor{}
		populateBaseSensorParams(&fds.UnimplementedSensor, params)
		fds.Path = params.PathToMountPoint
		ss.UpdateQueue.ExecuteAfter(fds.Name, DefaultWarmUpDelay)
		return &fds, nil

	case "free_disk_space_ratio":
		fdsr := system.FreeDiskSpaceRatioSensor{}
		populateBaseSensorParams(&fdsr.UnimplementedSensor, params)
		fdsr.Path = params.PathToMountPoint
		ss.UpdateQueue.ExecuteAfter(fdsr.Name, DefaultWarmUpDelay)
		return &fdsr, nil

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
		es := &endpoint.Sensor{
			Token: params.Token,
		}
		populateBaseSensorParams(&es.UnimplementedSensor, params)
		return es, nil

	case "file":
		fs := &file.Sensor{}
		populateBaseSensorParams(&fs.UnimplementedSensor, params)
		fs.PathToReadingsFile = params.PathToReading
		fs.JsonPath = params.JsonPath
		ss.UpdateQueue.ExecuteAfter(fs.Name, DefaultWarmUpDelay)
		return fs, nil

	case "victoria", "victoria_metrics", "victoria metrics":
		vs := &victoriametrics.VMSenor{}
		populateBaseSensorParams(&vs.UnimplementedSensor, params)
		vs.Endpoint = params.Endpoint
		vs.Headers = params.Headers
		vs.Query = params.Query
		vs.Filter = params.Filter
		ss.UpdateQueue.ExecuteAfter(vs.Name, DefaultWarmUpDelay)
		return vs, nil

	default:
		return nil, fmt.Errorf("unknown sensor type: %s", params.Type)
	}

}
