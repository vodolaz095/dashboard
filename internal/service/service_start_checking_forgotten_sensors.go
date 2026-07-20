package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) StartCheckingForgottenSensors(ctx context.Context) {
	pacemaker := time.NewTicker(time.Second) // not so fast
	var shouldBeUpdatedAt time.Time
	for {
		select {
		case <-ctx.Done():
			pacemaker.Stop()
			return
		case <-pacemaker.C:
			for k := range ss.Sensors {
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

				// if sensor have missed 2 updates it should be updated
				shouldBeUpdatedAt = ss.Sensors[k].GetUpdatedAt().Add(2 * ss.Sensors[k].GetRefreshRate())
				if shouldBeUpdatedAt.Before(time.Now()) {
					oneOfSensors := ss.Sensors[k]
					go func() {
						ctx2, cancel := context.WithTimeout(ctx, oneOfSensors.GetRefreshRate()/2)
						defer cancel()
						err := oneOfSensors.Update(ctx2)
						if err != nil {
							log.Error().Err(err).Msgf("Error updating sensod %v %s: %s",
								k, oneOfSensors.GetName(), err)
						}
					}()
				}
			}
		}
	}
}
