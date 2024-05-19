package service

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) Ping(ctx context.Context) (err error) {
	for k := range ss.Sensors {
		err = ss.Sensors[k].Ping(ctx)
		if err != nil {
			return
		}
		log.Trace().Msgf("Sensor %s online!", k)
	}
	log.Debug().Msgf("Sensors online")
	for k := range ss.MysqlConnections {
		err = ss.MysqlConnections[k].PingContext(ctx)
		if err != nil {
			return
		}
		log.Trace().Msgf("Mysql connection %s online!", k)
	}
	log.Debug().Msgf("Mysql connections online")
	for k := range ss.PostgresqlConnections {
		err = ss.PostgresqlConnections[k].PingContext(ctx)
		if err != nil {
			return
		}
		log.Trace().Msgf("Postgres connection %s online!", k)
	}
	log.Debug().Msgf("Postgres connections online")
	for k := range ss.RedisConnections {
		err = ss.RedisConnections[k].Ping(ctx).Err()
		if err != nil {
			return
		}
		log.Trace().Msgf("Redis connection %s online!", k)
	}
	log.Debug().Msgf("Redis connections online")
	return nil
}
