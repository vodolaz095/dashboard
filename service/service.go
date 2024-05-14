package service

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dqueue"

	"github.com/vodolaz095/dashboard/model"
	"github.com/vodolaz095/dashboard/sensors"
)

const DefaultSensorTimeout = 5 * time.Second

type SensorsService struct {
	ListOfSensors  []string
	Sensors        map[string]sensors.ISensor
	UpdateInterval time.Duration
	UpdateQueue    *dqueue.Handler
}

var SensorNotFoundErr = errors.New("sensor not found")

func (ss *SensorsService) List() (ret []model.Sensor) {
	ret = make([]model.Sensor, len(ss.ListOfSensors))
	var s model.Sensor
	for i := range ss.ListOfSensors {
		sensor, found := ss.Sensors[ss.ListOfSensors[i]]
		if found {
			s.Name = sensor.GetName()
			s.Type = sensor.GetType()
			s.Description = sensor.GetDescription()
			s.Link = sensor.GetLink()
			s.Minimum = sensor.GetMinimum()
			s.Maximum = sensor.GetMaximum()
			s.Value = sensor.Value()
			s.UpdatedAt = sensor.UpdatedAt()
			s.Tags = sensor.GetTags()
			ret[i] = s
		} else {
			ret[i] = model.Sensor{
				Name:        ss.ListOfSensors[i],
				Type:        "not found",
				Description: "Sensor not found",
				Link:        "",
				Minimum:     0,
				Maximum:     0,
				Value:       0,
				UpdatedAt:   time.Now(),
			}
		}
	}
	return ret
}

func (ss *SensorsService) Ping(ctx context.Context) (err error) {
	for k := range ss.Sensors {
		err = ss.Sensors[k].Ping(ctx)
		if err != nil {
			return
		}
	}
	return nil
}

func (ss *SensorsService) Close(ctx context.Context) (err error) {
	for k := range ss.Sensors {
		err = ss.Sensors[k].Close(ctx)
		if err != nil {
			return
		}
	}
	return nil
}

func (ss *SensorsService) Update(ctx context.Context, name string, val float64) (err error) {
	sensor, found := ss.Sensors[name]
	if !found {
		return SensorNotFoundErr
	}
	return sensor.Update(ctx, val)
}

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
				ss.UpdateQueue.ExecuteAfter(name, time.Second)
			}
		}
	}
}
