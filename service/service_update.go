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
	err = sensor.Update(ctx)
	if err != nil {
		n := ss.Broadcast(name, err.Error(), sensor.GetValue())
		log.Error().Err(err).
			Str("name", name).
			Float64("value", sensor.GetValue()).
			Int("notified", n).
			Msgf("Error updating sensor %s with value %v and %v notified: %s",
				name, sensor.GetValue(), n, err)
		return sensor.Next(), err
	}
	n := ss.Broadcast(name, "", sensor.GetValue())
	log.Debug().
		Str("name", name).
		Float64("value", sensor.GetValue()).
		Int("notified", n).
		Msgf("Sensor %s updated with value %v - %v notified", name, sensor.GetValue(), n)
	return sensor.Next(), nil
}
