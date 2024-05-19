package service

import "context"

func (ss *SensorsService) Close(ctx context.Context) (err error) {
	for k := range ss.MysqlConnections {
		err = ss.MysqlConnections[k].Close()
		if err != nil {
			return
		}
	}
	for k := range ss.PostgresqlConnections {
		err = ss.PostgresqlConnections[k].Close()
		if err != nil {
			return
		}
	}
	for k := range ss.RedisConnections {
		err = ss.RedisConnections[k].Close()
		if err != nil {
			return
		}
	}
	for k := range ss.Sensors {
		err = ss.Sensors[k].Close(ctx)
		if err != nil {
			return
		}
	}
	return nil
}
