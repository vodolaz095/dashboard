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
	Client    *redis.Client
	Channel   string
	ValueOnly bool
}

func (s *SubscribeSensor) Init(ctx context.Context) error {
	s.Mutex = &sync.RWMutex{}
	if s.A == 0 {
		s.A = 1
	}
	return s.Ping(ctx)
}

func (s *SubscribeSensor) Ping(ctx context.Context) error {
	return s.Client.Ping(ctx).Err()
}

func (s *SubscribeSensor) Close(_ context.Context) (err error) {
	// since canceling subscriber closes connection as expected
	return nil
}

func (s *SubscribeSensor) Update(_ context.Context) error {
	return nil
}

func (s *SubscribeSensor) ParseValue(msg *redis.Message) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
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
