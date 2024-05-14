package config

import "gopkg.in/yaml.v3"

type Config struct {
	Listen      string   `yaml:"listen" validate:"required,hostname_port"`
	Domain      string   `yaml:"domain" validate:"hostname_rfc1123"`
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Keywords    []string `yaml:"keywords"`
	DoIndex     bool     `yaml:"do_index"`
	Log         Log      `yaml:"log"`
	Sensors     []Sensor `yaml:"sensors"`
}

func (c *Config) Dump() ([]byte, error) {
	return yaml.Marshal(c)
}
