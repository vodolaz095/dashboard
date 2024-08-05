package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) StartCheckingForgottenSensors(ctx context.Context) {
	pacemaker := time.NewTicker(5 * ss.UpdateInterval) // not so fast
	var shouldBeUpdatedAt time.Time
	var k string
	for {
		select {
		case <-ctx.Done():
			pacemaker.Stop()
			return
		case <-pacemaker.C:
			for k = range ss.Sensors {
				// if sensor was not updated yet, we do not check it
				if ss.Sensors[k].GetUpdatedAt().IsZero() {
					continue
				}
				// some sensors are only updated by external event
				if ss.Sensors[k].GetType() == "subscriber" {
					continue
				}
				if ss.Sensors[k].GetType() == "endpoint" {
					continue
				}

				// if sensor have missed 2 updates it should be queued for refreshing
				shouldBeUpdatedAt = ss.Sensors[k].GetUpdatedAt().Add(2 * ss.Sensors[k].GetRefreshRate())
				if shouldBeUpdatedAt.Before(time.Now()) {
					next := time.Now().Add(ss.Sensors[k].GetRefreshRate())
					log.Warn().
						Str("sensor", ss.Sensors[k].GetName()).
						Time("next", next).
						Msgf("Sensor %s update forgotten - requeue on %s",
							ss.Sensors[k].GetName(), next.Format("15:04:05.000"),
						)
					ss.UpdateQueue.ExecuteAt(k, next)
				}
			}
		}
	}
}
