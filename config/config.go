package config

import "gopkg.in/yaml.v3"

// Config defines structure parsed from initial configuration file
type Config struct {
	WebUI WebUI `yaml:"web_ui" validate:"required"`
	// Log sets logging settings
	Log Log `yaml:"log" validate:"required"`
	// DatabaseConnections defines database connections being used by sensors
	DatabaseConnections []DatabaseConnection `yaml:"database_connections"`
	// Sensors defines configuration for methods of acquiring metrics data
	Sensors []Sensor `yaml:"sensors" validate:"required"`
	// Broadcasters defines configuration for sinks where we broadcast Sensors data
	Broadcasters []Broadcaster `yaml:"broadcasters"`
	// Influx defines configuration for Influxdb v2 time-series database used for storing historical sensor values
	Influx Influx `yaml:"influx" validate:"omitempty"`
}

func (c *Config) Dump() ([]byte, error) {
	return yaml.Marshal(c)
}
