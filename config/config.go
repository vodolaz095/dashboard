package config

import "gopkg.in/yaml.v3"

type Config struct {
	// Listen sets address, where application is listening, for example, 127.0.0.1:3000
	Listen string `yaml:"listen" validate:"required,hostname_port"`
	// Domain sets HTTP HOST where application accepts requests
	Domain string `yaml:"domain" validate:"hostname_rfc1123"`
	// Title sets title of index page
	Title string `yaml:"title"`
	// Description sets description of index page
	Description string `yaml:"description"`
	// Keywords sets keywords of index page
	Keywords []string `yaml:"keywords"`
	// DoIndex sets http header equivalents to allow page indexing by search engine crawlers
	DoIndex bool `yaml:"do_index"`
	// Log sets logging settings
	Log Log `yaml:"log"`
	// DatabaseConnections defines database connections being used by sensors
	DatabaseConnections []DatabaseConnection `yaml:"database_connections"`
	// Sensors defines configuration for methods of acquiring metrics data
	Sensors []Sensor `yaml:"sensors"`
}

func (c *Config) Dump() ([]byte, error) {
	return yaml.Marshal(c)
}
