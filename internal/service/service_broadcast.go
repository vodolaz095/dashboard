package service

import (
	"time"

	"github.com/vodolaz095/dashboard/model"
)

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
