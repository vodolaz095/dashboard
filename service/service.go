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

	subscribers map[string]chan model.Update
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
	err = sensor.Update(ctx, val)
	n := ss.Broadcast(name, err.Error(), val)
	if err != nil {
		log.Error().Err(err).
			Str("name", name).
			Float64("value", val).
			Int("notified", n).
			Msgf("Error updating sensor %s with value %v and %v notified: %s", name, val, n, err)
		return err
	}
	log.Debug().
		Str("name", name).
		Float64("value", sensor.Value()).
		Int("notified", n).
		Msgf("Sensor %s updated with value %v - %v notified", name, val, n)
	return nil
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

func (ss *SensorsService) Subscribe(ctx context.Context, name string) (chan model.Update, error) {
	_, found := ss.subscribers[name]
	if found {
		return nil, errors.New("duplicate subscriber name")
	}
	log.Debug().Msgf("Creating subscription channel for %s...", name)
	ch := make(chan model.Update, 10)
	ss.subscribers[name] = ch
	go func() {
		<-ctx.Done()
		log.Debug().Msgf("Closing subscription channel for %s...", name)
		close(ch)
		delete(ss.subscribers, name)
		log.Debug().Msgf("Subscription channel for %s is closed", name)

	}()

	return ch, nil
}

func (ss *SensorsService) Broadcast(name, error string, value float64) (subscribersNotified int) {
	upd := model.Update{
		Name:      name,
		Value:     value,
		Error:     error,
		Timestamp: time.Now(),
	}
	for k := range ss.subscribers {
		subscribersNotified += 1
		ss.subscribers[k] <- upd
	}
	return subscribersNotified
}
