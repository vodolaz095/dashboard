package sensors

import (
	"context"
	"time"
)

type ISensor interface {
	Value() float64
	Update(context.Context, float64) error
	UpdatedAt() time.Time
	Close() error
}

type UnimplementedSensor struct {
	// Name is used to distinguish sensors from other ones
	Name string `yaml:"name" validate:"required,alphanum"`
	// Type is used to define strategy to load sensor value
	Type string `yaml:"type" validate:"required, oneof=mysql redis postgres sqlite curl shell endpoint"`
	// DatabaseConnectionString is used to dial database, remote url
	DatabaseConnectionString string `yaml:"database_connection_string"`
	// Query is used to either execute against remote resourse or process raw data
	Query string `yaml:"query"`
	// RefreshRate is used to define how often we reload data
	RefreshRate time.Duration `yaml:"refresh_rate"`
	// Description is used to explain meaning of this sensor
	Description string `yaml:"description"`
	// Link is used to help visitor read more about sensor
	Link string `yaml:"link" validate:"http_url"`
	// Minimum is used to warn, when something is below safe value
	Minimum float64 `yaml:"minimum"`
	// Maximum is used to warn, when something is above safe value
	Maximum float64 `yaml:"maximum"`
}
