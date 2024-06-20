package config

type Influx struct {
	Endpoint     string `yaml:"endpoint" validate:"http_url"`
	Token        string `yaml:"token"`
	Organization string `yaml:"organization"`
	Bucket       string `yaml:"bucket"`
}
