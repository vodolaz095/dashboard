package service

import "context"

func (ss *SensorsService) InitSensors(ctx context.Context) (err error) {
	for k := range ss.Sensors {
		err = ss.Sensors[k].Init(ctx)
		if err != nil {
			return
		}
	}
	return nil
}
