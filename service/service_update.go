package service

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (ss *SensorsService) Update(ctx context.Context, name string, val float64) (err error) {
	sensor, found := ss.Sensors[name]
	if !found {
		return SensorNotFoundErr
	}
	err = sensor.Update(ctx, val)
	if err != nil {
		n := ss.Broadcast(name, err.Error(), val)
		log.Error().Err(err).
			Str("name", name).
			Float64("value", val).
			Int("notified", n).
			Msgf("Error updating sensor %s with value %v and %v notified: %s", name, val, n, err)
		return err
	}
	n := ss.Broadcast(name, "", val)
	log.Debug().
		Str("name", name).
		Float64("value", sensor.Value()).
		Int("notified", n).
		Msgf("Sensor %s updated with value %v - %v notified", name, sensor.Value(), n)
	return nil
}
