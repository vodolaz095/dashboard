package service

import (
	"time"

	"github.com/vodolaz095/dashboard/model"
)

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
			s.Value = sensor.GetValue()
			s.UpdatedAt = sensor.GetUpdatedAt()
			s.Tags = sensor.GetTags()
			if sensor.GetLastError() != nil {
				s.Error = sensor.GetLastError().Error()
			} else {
				s.Error = ""
			}
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
