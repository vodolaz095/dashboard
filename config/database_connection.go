package config

// DatabaseConnection defines names and types of database connections used by Sensors, Publishers and Subscribbers
type DatabaseConnection struct {
	Name                     string `yaml:"name" validate:"required"`
	Type                     string `yaml:"type" validate:"required, oneof=mysql postgres redis"`
	DatabaseConnectionString string `yaml:"connection_string"`
	MaxOpenCons              int    `yaml:"max_open_cons" validate:"gte=0"`
	MaxIdleCons              int    `yaml:"max_idle_cons" validate:"gte=0"`
}
