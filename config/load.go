package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func LoadFromFile(pathToFile string) (cfg Config, err error) {
	raw, err := os.ReadFile(pathToFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(raw, &cfg)
	return
}
