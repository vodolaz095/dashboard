package service

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/sensors/redis"
)

func (ss *SensorsService) StartRedisSubscribers(ctx context.Context) {
	for k := range ss.Sensors {
		casted, ok := ss.Sensors[k].(*redis.SubscribeSensor)
		if !ok {
			continue
		}
		_, found := ss.RedisConnections[casted.DatabaseConnectionName]
		if !found {
			log.Fatal().Msgf("Redis subscriber sensor %s uses unknown connection %s",
				casted.Name, casted.DatabaseConnectionName,
			)
			return
		}
		go func() {
			log.Info().Msgf("Starting redis subscriber %s on channel %s...",
				casted.Name, casted.Channel)
			sub := casted.Client.Subscribe(ctx, casted.Channel)
			defer sub.Close()
			ch := sub.Channel()
			for msg := range ch {
				casted.ParseValue(msg)
				if casted.Error != nil {
					ss.Broadcast(casted.Name, casted.Error.Error(), casted.Value)
				} else {
					ss.Broadcast(casted.Name, "", casted.Value)
				}
			}
			log.Info().Msgf("Stopping redis subscriber %s on %s...",
				casted.Name, casted.Channel)
		}()
	}
}
