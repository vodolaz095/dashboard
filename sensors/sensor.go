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

	Update(context.Context, float64) error

	Value() float64
	UpdatedAt() time.Time
}

type UnimplementedSensor struct {
	// Name is used to distinguish sensors from other ones
	Name string `yaml:"name" validate:"required,alphanum"`
	// Type is used to define strategy to load sensor value
	Type string `yaml:"type" validate:"required, oneof=mysql redis postgres sqlite curl shell endpoint"`
	// DatabaseConnectionString is used to dial database, remote url
	DatabaseConnectionString string `yaml:"database_connection_string"`
	// Query is used to either execute against remote resource or to process raw data
	Query string `yaml:"query"`
	// RefreshRate is used to define how often we reload data
	RefreshRate time.Duration `yaml:"refresh_rate"`
	// Description is used to explain meaning of this sensor
	Description string `yaml:"description"`
	// Link is used to help visitor read more about sensor
	Link string `yaml:"link" validate:"http_url"`
	// Minimum is used to warn, when something is below safe value
	Minimum float64 `yaml:"minimum"`
	// Maximum is used to warn, when something is above safe value
	Maximum float64 `yaml:"maximum"`
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
