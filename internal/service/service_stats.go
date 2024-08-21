package service

import "github.com/vodolaz095/dashboard/model"

func (ss *SensorsService) Stats() model.Stats {
	return model.Stats{
		SensorsUpdatedNow: int(ss.SensorsBeingUpdated),
		QueueLength:       ss.UpdateQueue.Len(),
		Subscribers:       len(ss.subscribers),
	}
}
