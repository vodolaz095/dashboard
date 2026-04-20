package service

import "errors"

var (
	ErrSensorNotFound                = errors.New("sensor not found")
	ErrDuplicateConnection           = errors.New("duplicate connection")
	ErrDuplicateSubscriber           = errors.New("duplicate subscriber")
	ErrConnectionNotFound            = errors.New("connection not found")
	ErrUnknownDatabaseConnectionType = errors.New("unknown database connection type")
)
