package service

import "context"

func (ss *SensorsService) MakeConnection(ctx context.Context, name string,
	kind DatabaseConnectionType, databaseConnectionString string) (err error) {
	switch kind {
	case DatabaseConnectionTypeMysql, DatabaseConnectionTypeMariadb:
		err = ss.initMysqlConnection(ctx, name, databaseConnectionString)
		break
	case DatabaseConnectionTypePostgres:
		err = ss.initPostgresConnection(ctx, name, databaseConnectionString)
		break
	case DatabaseConnectionTypeRedis:
		err = ss.initRedisConnection(ctx, name, databaseConnectionString)
		break
	default:
		return UnknownDatabaseConnectionTypeError
	}
	return err
}
