package service

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) Refresh(ctx context.Context, name string) (err error) {
	sensor, found := ss.Sensors[name]
	if !found {
		return SensorNotFoundErr
	}
	err = sensor.Update(ctx)
	if err != nil {
		n := ss.Broadcast(name, err.Error(), sensor.Value())
		log.Error().Err(err).
			Str("name", name).
			Float64("value", sensor.Value()).
			Int("notified", n).
			Msgf("Error updating sensor %s with value %v and %v notified: %s",
				name, sensor.Value(), n, err)
		return err
	}
	n := ss.Broadcast(name, "", sensor.Value())
	log.Debug().
		Str("name", name).
		Float64("value", sensor.Value()).
		Int("notified", n).
		Msgf("Sensor %s updated with value %v - %v notified", name, sensor.Value(), n)
	return nil
}
