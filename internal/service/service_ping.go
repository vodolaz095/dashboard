package service

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) Ping(ctx context.Context) (err error) {
	for k, con := range ss.MysqlConnections {
		err = con.PingContext(ctx)
		if err != nil {
			return
		}
		log.Trace().Msgf("Mysql connection %s online!", k)
	}
	log.Debug().Msgf("Mysql connections online")
	for k, con := range ss.PostgresqlConnections {
		err = con.PingContext(ctx)
		if err != nil {
			return
		}
		log.Trace().Msgf("Postgres connection %s online!", k)
	}
	log.Debug().Msgf("Postgres connections online")
	for k, con := range ss.RedisConnections {
		err = con.Ping(ctx).Err()
		if err != nil {
			return
		}
		log.Trace().Msgf("Redis connection %s online!", k)
	}
	log.Debug().Msgf("Redis connections online")

	for k := range ss.Sensors {
		err = ss.Sensors[k].Ping(ctx)
		if err != nil {
			return
		}
		log.Trace().Msgf("Sensor %s online!", k)
	}
	log.Debug().Msgf("Sensors online")
	return nil
}
