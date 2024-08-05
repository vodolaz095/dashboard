package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) StartRefreshingSensors(ctx context.Context) {
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
				name = task.Payload.(string)
				go func() {
					rtc, cancel := context.WithTimeout(ctx, DefaultSensorTimeout)
					defer cancel()
					nextUpdateOn, err := ss.Refresh(rtc, name)
					if err != nil {
						log.Error().Err(err).
							Str("sensor", name).
							Msgf("Sensor %s update failed with %s",
								task.Payload.(string), err,
							)
					}
					ss.UpdateQueue.ExecuteAt(name, nextUpdateOn)
				}()
			}
		}
	}
}
