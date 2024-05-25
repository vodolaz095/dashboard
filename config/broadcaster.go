package config

type Broadcaster struct {
	ConnectionName string `yaml:"connection_name"`
	Subject        string `yaml:"subject"`
}
