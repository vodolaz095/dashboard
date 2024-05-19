package service

type DatabaseConnectionType string

const (
	DatabaseConnectionTypeMysql    DatabaseConnectionType = "mysql"
	DatabaseConnectionTypeMariadb  DatabaseConnectionType = "mariadb"
	DatabaseConnectionTypePostgres DatabaseConnectionType = "postgres"
	DatabaseConnectionTypeRedis    DatabaseConnectionType = "redis"
)
