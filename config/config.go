package config

import "gopkg.in/yaml.v3"

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
}

func (c *Config) Dump() ([]byte, error) {
	return yaml.Marshal(c)
}
