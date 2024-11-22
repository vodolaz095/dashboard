package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/vodolaz095/dashboard/model"
)

func (ss *SensorsService) StartClock(ctx context.Context) {
	timer := time.NewTicker(500 * time.Millisecond)
	var stats model.Stats
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Info().Msgf("Clock is stopping...")
			return
		case t := <-timer.C:
			ss.Broadcast("clock", "", "", float64(t.Unix()))
			stats = ss.Stats()
			ss.Broadcast("sensors_updated_now", "", "", float64(stats.SensorsUpdatedNow))
			ss.Broadcast("queue_length", "", "", float64(stats.QueueLength))
			ss.Broadcast("subscribers", "", "", float64(stats.Subscribers))
		}
	}
}
