package redis

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/model"
	"github.com/vodolaz095/dashboard/sensors"
)

type SubscribeSensor struct {
	sensors.UnimplementedSensor
	mu        *sync.Mutex
	pubsub    *redis.PubSub
	Client    *redis.Client
	Channel   string
	ValueOnly bool
}

func (s *SubscribeSensor) Init(ctx context.Context) error {
	s.mu = &sync.Mutex{}
	if s.A == 0 {
		s.A = 1
	}
	return s.Ping(ctx)
}

func (s *SubscribeSensor) Ping(ctx context.Context) error {
	return s.Client.Ping(ctx).Err()
}

func (s *SubscribeSensor) Close(ctx context.Context) error {
	err := s.pubsub.Close()
	if err != nil {
		if errors.Is(err, redis.ErrClosed) {
			return nil
		}
	}
	err = s.Client.Close()
	if err != nil {
		if errors.Is(err, redis.ErrClosed) {
			return nil
		}
	}
	return err
}

func (s *SubscribeSensor) Update(_ context.Context) error {
	return nil
}

func (s *SubscribeSensor) parseValue(msg *redis.Message) {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Trace().Msgf("Parsing %s from channel %s", msg.Payload, msg.Channel)

	var val float64
	var err error
	if s.ValueOnly {
		val, err = strconv.ParseFloat(strings.TrimSpace(msg.Payload), 64)
		s.Value = val
		s.Error = err
		s.UpdatedAt = time.Now()
		return
	}
	var payload model.Update
	err = json.Unmarshal([]byte(msg.Payload), &payload)
	if err != nil {
		s.Value = 0
		s.Error = err
		s.UpdatedAt = time.Now()
		return
	}
	s.Value = payload.Value
	s.Error = errors.New(payload.Error)
	s.UpdatedAt = payload.Timestamp
	return
}

func (s *SubscribeSensor) Start(ctx context.Context) {
	log.Warn().Msgf("Starting redis subscriber on %s...", s.Channel)
	s.pubsub = s.Client.Subscribe(ctx, s.Channel)
	ch := s.pubsub.Channel()
	for msg := range ch {
		s.parseValue(msg)
	}
	log.Warn().Msgf("Stopping redis subscriber on %s...", s.Channel)
}
