package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) StartClock(ctx context.Context) {
	timer := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			timer.Stop()
			log.Info().Msgf("Clock is stopping...")
			return
		case t := <-timer.C:
			ss.Broadcast("clock", "", float64(t.Unix()))
		}
	}
}
