package config

// Log defines logging configuration
type Log struct {
	Level      string `yaml:"level" validate:"required,oneof=trace debug info warn error fatal"`
	ToJournald bool   `yaml:"to_journald"`
}
