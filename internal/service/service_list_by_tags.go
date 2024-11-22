package service

import (
	"github.com/vodolaz095/dashboard/model"
)

func (ss *SensorsService) ListByTags(tags map[string]string) (ret []model.Sensor) {
	var s model.Sensor
	var allTagsMatched bool
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
			s.Status = sensor.GetStatus()
			if sensor.GetLastError() != nil {
				s.Error = sensor.GetLastError().Error()
			} else {
				s.Error = ""
			}
			allTagsMatched = false
			tagsAvailable := sensor.GetTags()
			for k := range tags {
				present, ok := tagsAvailable[k]
				if ok {
					allTagsMatched = present == tags[k]
				} else {
					allTagsMatched = false
				}
			}
			if allTagsMatched {
				ret = append(ret, s)
			}
		}
	}
	return ret
}
