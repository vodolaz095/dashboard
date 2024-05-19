package service

import "errors"

var SensorNotFoundErr = errors.New("sensor not found")

var DuplicateConnectionError = errors.New("duplicate connection")

var DuplicateSubscriberError = errors.New("duplicate subscriber")

var UnknownDatabaseConnectionTypeError = errors.New("unkown database connection type")
