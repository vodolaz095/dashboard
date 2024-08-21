package service

import "time"

type DatabaseConnectionType string

const (
	DatabaseConnectionTypeMysql            DatabaseConnectionType = "mysql"
	DatabaseConnectionTypeMariadb          DatabaseConnectionType = "mariadb"
	DatabaseConnectionTypePostgres         DatabaseConnectionType = "postgres"
	DatabaseConnectionTypeRedis            DatabaseConnectionType = "redis"
	DefaultSubscriptionChannelChannelDepth                        = 10
	DefaultWarmUpDelay                                            = 50 * time.Millisecond
	DefaultSensorTimeout                                          = 5 * time.Second
	DefaultRefreshRate                                            = 5 * time.Second
)
