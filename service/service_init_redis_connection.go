package service

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func (ss *SensorsService) initRedisConnection(ctx context.Context, name, dsn string) (err error) {
	_, found := ss.RedisConnections[name]
	if found {
		return DuplicateConnectionError
	}
	opts, err := redis.ParseURL(dsn)
	if err != nil {
		return err
	}

	client := redis.NewClient(opts)
	err = client.Ping(ctx).Err()
	if err != nil {
		return err
	}
	ss.RedisConnections[name] = client
	return nil
}
