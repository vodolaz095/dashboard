package service

import "errors"

var (
	SensorNotFoundErr                  = errors.New("sensor not found")
	DuplicateConnectionError           = errors.New("duplicate connection")
	DuplicateSubscriberError           = errors.New("duplicate subscriber")
	ConnectionNotFoundError            = errors.New("connection not found")
	UnknownDatabaseConnectionTypeError = errors.New("unknown database connection type")
)
