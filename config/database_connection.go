package config

type DatabaseConnection struct {
	Name                     string `yaml:"name" validate:"required"`
	Type                     string `yaml:"type" validate:"required, oneof=mysql postgres redis"`
	DatabaseConnectionString string `yaml:"database_connection_string"`
}
