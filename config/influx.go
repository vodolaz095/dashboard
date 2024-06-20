package config

type Influx struct {
	Endpoint     string `yaml:"endpoint" validate:"http_url"`
	Token        string `yaml:"token"`
	Organization string `yaml:"organization"`
	Bucket       string `yaml:"bucket"`
}

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
