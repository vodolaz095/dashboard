package redis

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/vodolaz095/dashboard/internal/sensors/redis"
	"github.com/vodolaz095/dashboard/internal/service"
)

// Subscriber subscribes to redis channels
type Subscriber struct {
	Service *service.SensorsService
}

func (ss *Subscriber) Start(ctx context.Context) {
	for k := range ss.Service.Sensors {
		casted, ok := ss.Service.Sensors[k].(*redis.SubscribeSensor)
		if !ok {
			continue
		}
		_, found := ss.Service.RedisConnections[casted.DatabaseConnectionName]
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
					ss.Service.Broadcast(casted.Name, casted.Error.Error(), casted.GetStatus(), casted.Value)
				} else {
					ss.Service.Broadcast(casted.Name, "", casted.GetStatus(), casted.Value)
				}
			}
			log.Info().Msgf("Stopping redis subscriber %s on channel %s...",
				casted.Name, casted.Channel)
		}()
	}
}
