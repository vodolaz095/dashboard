package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/config"
	"github.com/vodolaz095/dashboard/model"
	"github.com/vodolaz095/dashboard/service"
)

type redisSink struct {
	Client    *redis.Client
	Subject   string
	ValueOnly bool
	Sensors   map[string]bool
}

func (rs *redisSink) CanBroadcast(upd *model.Update) (ok bool) {
	if len(rs.Sensors) == 0 {
		return true
	}
	_, sensorNameWhitelisted := rs.Sensors[upd.Name]
	if sensorNameWhitelisted {
		return true
	}
	return false
}

// Publisher broadcast Sensor values into redis via `pub/sub` channels
type Publisher struct {
	Service    *service.SensorsService
	redisSinks []redisSink
}

// InitConnection initialize redis connections for publishing sensor values
func (p *Publisher) InitConnection(params config.Broadcaster) error {
	client, found := p.Service.RedisConnections[params.ConnectionName]
	if !found {
		return service.ConnectionNotFoundError
	}
	sink := redisSink{
		Client:    client,
		Subject:   params.Subject,
		ValueOnly: params.ValueOnly,
		Sensors:   make(map[string]bool),
	}
	for k := range params.SensorsToListen {
		sink.Sensors[params.SensorsToListen[k]] = true
	}
	p.redisSinks = append(p.redisSinks, sink)
	return nil
}

// Start starts broadcasting new sensor readings into redis channels
func (p *Publisher) Start(ctx context.Context) {
	feed, err := p.Service.Subscribe(ctx, "dashboard.broadcaster.redis")
	if err != nil {
		log.Fatal().Err(err).Msgf("broadcaster failed to subscribe: %s", err)
	}
	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("Redis broadcaster is closing...")
			for i := range p.redisSinks {
				err = p.redisSinks[i].Client.Close()
				if err != nil {
					if !errors.Is(err, redis.ErrClosed) {
						log.Error().Err(err).Msgf("error closing redis sink: %s", err)
					}
				}
			}
			return

		case upd := <-feed:
			for i := range p.redisSinks {
				if !p.redisSinks[i].CanBroadcast(&upd) {
					continue
				}
				if p.redisSinks[i].ValueOnly {
					err = p.redisSinks[i].Client.Publish(ctx,
						fmt.Sprintf(p.redisSinks[i].Subject, upd.Name), upd.Value,
					).Err()
				} else {
					err = p.redisSinks[i].Client.Publish(ctx,
						fmt.Sprintf(p.redisSinks[i].Subject, upd.Name), upd.Pack(),
					).Err()
				}
				if err != nil {
					log.Error().Err(err).Msgf("error publishing into %s: %s",
						fmt.Sprintf(p.redisSinks[i].Subject, upd.Name), err,
					)
				}
			}
		}
	}
}
