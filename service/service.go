package service

import (
	"context"
	"errors"
	"time"

	"github.com/vodolaz095/dashboard/model"
	"github.com/vodolaz095/dashboard/sensors"
)

type SensorsService struct {
	ListOfSensors []string
	Sensors       map[string]sensors.ISensor
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
