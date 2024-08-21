package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) Refresh(ctx context.Context, name string) (next time.Time, err error) {
	sensor, found := ss.Sensors[name]
	if !found {
		return time.Now(), SensorNotFoundErr
	}

	next = time.Now().Add(sensor.GetRefreshRate())
	log.Trace().
		Str("sensor", name).
		Float64("reading", sensor.GetValue()).
		Time("next", next).
		Msgf("Preparing to update sensor %s...", name)

	err = sensor.Update(ctx)
	if err != nil {
		n := ss.Broadcast(name, err.Error(), sensor.GetValue())
		log.Error().Err(err).
			Str("sensor", name).
			Float64("reading", sensor.GetValue()).
			Int("notified", n).
			Time("next", next).
			Msgf("Error updating sensor %s with value %v and %v notified: %s",
				name, sensor.GetValue(), n, err)
		return next, err
	}
	n := ss.Broadcast(name, "", sensor.GetValue())
	log.Debug().
		Str("sensor", name).
		Float64("reading", sensor.GetValue()).
		Time("next", next).
		Int("notified", n).
		Msgf("Sensor %s updated with value %v - %v notified", name, sensor.GetValue(), n)
	return next, nil
}
