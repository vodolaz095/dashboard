package service

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/config"
)

func (ss *SensorsService) initRedisConnection(ctx context.Context, params config.DatabaseConnection) (err error) {
	_, found := ss.RedisConnections[params.Name]
	if found {
		return DuplicateConnectionError
	}
	opts, err := redis.ParseURL(params.DatabaseConnectionString)
	if err != nil {
		return err
	}
	opts.MaxIdleConns = params.MaxIdleCons
	opts.MaxActiveConns = params.MaxOpenCons
	client := redis.NewClient(opts)
	err = client.Ping(ctx).Err()
	if err != nil {
		return err
	}
	ss.RedisConnections[params.Name] = client
	log.Info().Msgf("Redis database connection %s is established", params.Name)
	return nil
}
