package sensors

import (
	"context"
	"testing"
	"time"
)

type ISensor interface {
	Init(ctx context.Context) error
	Ping(ctx context.Context) error
	Close(ctx context.Context) error

	GetName() string
	GetType() string
	GetDescription() string
	GetLink() string
	GetMinimum() float64
	GetMaximum() float64
	GetTags() map[string]string
	/*
		To be implemented in custom sensors
	*/
	Update(context.Context, float64) error
	Value() float64
	UpdatedAt() time.Time
}

type UnimplementedSensor struct {
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

	// DatabaseConnectionString is used to dial database
	DatabaseConnectionString string `yaml:"database_connection_string"`
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
	Environment map[string]string `yaml:"headers"`

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

func (u *UnimplementedSensor) GetName() string {
	return u.Name
}

func (u *UnimplementedSensor) GetType() string {
	return u.Type
}

func (u *UnimplementedSensor) GetDescription() string {
	return u.Description
}

func (u *UnimplementedSensor) GetLink() string {
	return u.Link
}

func (u *UnimplementedSensor) GetMinimum() float64 {
	return u.Minimum
}

func (u *UnimplementedSensor) GetMaximum() float64 {
	return u.Maximum
}

func (u *UnimplementedSensor) GetTags() map[string]string {
	return u.Tags
}

func (u *UnimplementedSensor) NextUpdateOn() time.Time {
	return time.Now().Add(u.RefreshRate)
}

const DefaultTestTimeout = time.Second

func DoTestSensor(t *testing.T, sensor ISensor, expected float64) (err error) {
	const readAttempts = 100
	var val float64
	var i int

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTestTimeout)
	defer cancel()
	err = sensor.Init(ctx)
	if err != nil {
		t.Errorf("error initializing: %s", err)
		return
	}
	t.Logf("Sensor initialized!")
	err = sensor.Ping(ctx)
	if err != nil {
		t.Errorf("error pinging: %s", err)
		return
	}
	t.Logf("Sensor pinged!")
	err = sensor.Update(ctx, expected)
	if err != nil {
		t.Errorf("error updating: %s", err)
		return
	}
	t.Logf("Sensor updated with %.4f...", expected)
	for i = 0; i < readAttempts; i++ {
		val = sensor.Value()
		if val != expected {
			t.Errorf("unexpected value - %.4f vs %.4f on %v run",
				sensor.Value(), val, i)
		}
	}
	t.Logf("Value %.4f is retrived %v times", expected, i)
	err = sensor.Close(ctx)
	if err != nil {
		t.Errorf("error closing: %s", err)
		return
	}
	t.Logf("Sensor closed!")
	return nil
}
