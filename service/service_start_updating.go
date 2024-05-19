package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) StartKeepingSensorsUpToDate(ctx context.Context) {
	var name string
	pacemaker := time.NewTicker(ss.UpdateInterval)
	for {
		select {
		case <-ctx.Done():
			pacemaker.Stop()
			return
		case <-pacemaker.C:
			task, ready := ss.UpdateQueue.Get()
			if ready {
				rtc, cancel := context.WithTimeout(ctx, DefaultSensorTimeout)
				name = task.Payload.(string)
				log.Debug().Msgf("Updating sensor %s...", name)
				err := ss.Update(rtc, name, 0)
				cancel()
				if err != nil {
					log.Error().Err(err).Msgf("Sensor %s update failed with %s",
						task.Payload.(string), err,
					)
				} else {
					log.Debug().Msgf("Sensor %s is updated!", name)
				}
				ss.UpdateQueue.ExecuteAfter(name, time.Second) // TODO
			}
		}
	}
}
