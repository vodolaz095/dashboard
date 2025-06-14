package service

import (
	"context"

	"github.com/vodolaz095/dashboard/config"
)

func (ss *SensorsService) MakeConnection(ctx context.Context, opts config.DatabaseConnection) (err error) {
	switch DatabaseConnectionType(opts.Type) {
	case DatabaseConnectionTypeMysql, DatabaseConnectionTypeMariadb:
		err = ss.initMysqlConnection(ctx, opts)
		break
	case DatabaseConnectionTypePostgres:
		err = ss.initPostgresConnection(ctx, opts)
		break
	case DatabaseConnectionTypeRedis:
		err = ss.initRedisConnection(ctx, opts)
		break
	default:
		return UnknownDatabaseConnectionTypeError
	}
	return err
}
