package config

import "time"

type Sensor struct {
	// Name is used to distinguish sensors from other ones
	Name string `yaml:"name" validate:"required,alphanum"`
	// Type is used to define strategy to load sensor value
	Type string `yaml:"type" validate:"required, oneof=mysql redis postgres curl shell endpoint"`
	// Description is used to explain meaning of this sensor
	Description string `yaml:"description"`
	// Link is used to help visitor read more about sensor
	Link string `yaml:"link" validate:"http_url"`
	// Tags helps to group sensors
	Tags map[string]string `yaml:"tags"`

	/*
	 * Parameters used for mysql, redis and postgres
	 */

	// ConnectionName is used to dial database
	ConnectionName string `yaml:"connection_name"`
	// Query is send to remote database in order to receive data from it
	Query string `yaml:"query"`

	/*
	 * Parameters used for curl sensor, which sends HTTP requests to remote servers to obtain data
	 */
	// HttpMethod defines request type being send to remote http(s) endpoint via HTTP protocol
	HttpMethod string `yaml:"http_method" validate:"oneof=GET HEAD POST PUT PATCH DELETE CONNECT OPTIONS TRACE"`
	// Endpoint defines URL where sensor sends request to recieve data
	Endpoint string `yaml:"endpoint" validate:"http_url"`
	// Headers are HTTP request headers being send with any HTTP request
	Headers map[string]string `yaml:"headers"`
	// Body is send with any HTTP request as payload
	Body string `yaml:"body"`
	// JsonPath is used to extract elements from json response of remote endpoint or shell command output using https://jsonpath.com/ syntax
	JsonPath string `yaml:"json_path"`

	/*
	 * Parameters used for shell sensor - which executes scripts or commands
	 */

	// Command is shell command being executed by shell sensor
	Command string `yaml:"command"`
	// Environment is POSIX environment used by shell sensor to execute commands into
	Environment map[string]string `yaml:"environment"`

	/*
	 * Parameters used for endpoint sensor
	 */
	// Token is Bearer strategy token used to send metrics for endpoint sensor
	Token string `json:"token"`

	// RefreshRate is used to define how often we reload data
	RefreshRate time.Duration `yaml:"refresh_rate"`
	// Minimum is used to warn, when something is below safe value
	Minimum float64 `yaml:"minimum"`
	// Maximum is used to warn, when something is above safe value
	Maximum float64 `yaml:"maximum"`
}
