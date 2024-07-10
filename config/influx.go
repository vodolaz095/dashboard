package config

// Influx defines credentials used to send Sensor readings into Influxdb of 2nd version via wire protocol
type Influx struct {
	Endpoint     string `yaml:"endpoint" validate:"http_url"`
	Token        string `yaml:"token"`
	Organization string `yaml:"organization"`
	Bucket       string `yaml:"bucket"`
}

// Valid used to check, if Influx config is sane
func (i Influx) Valid() bool {
	if i.Endpoint == "" {
		return false
	}
	if i.Token == "" {
		return false
	}
	if i.Organization == "" {
		return false
	}
	if i.Bucket == "" {
		return false
	}
	return true
}
