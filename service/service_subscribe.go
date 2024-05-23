package service

import (
	"context"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/model"
)

func (ss *SensorsService) Subscribe(ctx context.Context, name string) (chan model.Update, error) {
	_, found := ss.subscribers[name]
	if found {
		return nil, DuplicateSubscriberError
	}
	log.Debug().Msgf("Creating subscription channel for %s...", name)
	ch := make(chan model.Update, DefaultSubscriptionChannelChannelDepth)
	if ss.subscribers == nil {
		ss.subscribers = make(map[string]chan model.Update, 0)
	}
	ss.subscribers[name] = ch
	go func() {
		<-ctx.Done()
		log.Debug().Msgf("Closing subscription channel for %s...", name)
		close(ch)
		delete(ss.subscribers, name)
		log.Debug().Msgf("Subscription channel for %s is closed", name)

	}()

	return ch, nil
}
