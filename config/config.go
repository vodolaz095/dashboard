package config

type Config struct {
	Listen  string   `yaml:"listen" validate:"required,hostname_port"`
	Domain  string   `yaml:"domain" validate:"hostname_rfc1123"`
	Sensors []Sensor `yaml:"sensors"`
}
