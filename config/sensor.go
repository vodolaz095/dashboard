package config

import "time"

type Sensor struct {
	/*
	 * Shared parameters used by all sensor types
	 */

	// Name is used to distinguish sensors from other ones
	Name string `yaml:"name" validate:"required,alphanum"`
	// Type is used to define strategy to load sensor value
	Type string `yaml:"type" validate:"required"`
	// Description is used to explain meaning of this sensor
	Description string `yaml:"description"`
	// Link is used to help visitor read more about sensor
	Link string `yaml:"link" validate:"http_url"`
	// Tags helps to group sensors
	Tags map[string]string `yaml:"tags"`
	// RefreshRate is used to define how often we reload data
	RefreshRate time.Duration `yaml:"refresh_rate"`
	// Minimum is used to warn, when something is below safe value
	Minimum float64 `yaml:"minimum"`
	// Maximum is used to warn, when something is above safe value
	Maximum float64 `yaml:"maximum"`
	// A is coefficient in linear transformation Y=A*X+B used to, for example, convert
	// Fahrenheit degrees into Celsius degrees
	A float64 `yaml:"a"`
	// B is constant in linear transformation Y=A*X+B used to, for example, convert
	// Fahrenheit degrees into Celsius degrees
	B float64 `yaml:"b"`

	/*
	 * Parameters used for mysql, redis and postgres sensors
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

	/*
	 * Parameters used for file sensor
	 */

	// PathToReading used for file sensor reading constant from file, for example
	// cat /sys/class/thermal/thermal_zone1/temp
	// gives temperature sensor reading. It is worth notice that JsonPath parameter is taken into account
	PathToReading string `yaml:"path_to_reading"`

	/*
	 * Parameters used for subscribing to redis change feed
	 */
	// Channel defines redis channel name to use for getting new values, something like `PUBLISH vodolaz095/dashboard/channel_name 15.3`.
	Channel string `yaml:"channel"`
	// ValueOnly defines, if data is sent via redis channel as raw float64 or as model.Update encoded in json.
	ValueOnly bool `yaml:"value_only"`
}
