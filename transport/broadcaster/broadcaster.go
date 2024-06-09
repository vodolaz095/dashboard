package broadcaster

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/service"
)

type redisSink struct {
	Client    *redis.Client
	Subject   string
	ValueOnly bool
}

// Publisher broadcast Sensor values into redis via `pub/sub` channels
type Publisher struct {
	Service    *service.SensorsService
	redisSinks []redisSink
}

// InitConnection initialize redis connections for publishing sensor values
func (p *Publisher) InitConnection(name, subject string, valueOnly bool) error {
	client, found := p.Service.RedisConnections[name]
	if !found {
		return service.ConnectionNotFoundError
	}
	p.redisSinks = append(p.redisSinks, redisSink{
		Client:    client,
		Subject:   subject,
		ValueOnly: valueOnly,
	})
	return nil
}

// Start starts broadcasting new sensor readings into redis channels
func (p *Publisher) Start(ctx context.Context) {
	feed, err := p.Service.Subscribe(ctx, "dashboard.broadcaster")
	if err != nil {
		log.Fatal().Err(err).Msgf("broadcaster failed to subscribe: %s", err)
	}
	for {
		select {
		case <-ctx.Done():
			log.Info().Msgf("Broadcaster is closing...")
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