package sensors

import (
	"context"
	"sync"
	"testing"
	"time"
)

type ISensor interface {
	/*
		Implemented in abstract class of UnimplementedSensor
	*/
	GetName() string
	GetType() string
	GetDescription() string
	GetLink() string
	GetMinimum() float64
	GetMaximum() float64
	GetTags() map[string]string
	GetValue() float64
	GetUpdatedAt() time.Time
	GetLastError() error
	Next() time.Time
	/*
		To be implemented in custom sensors
	*/
	Init(context.Context) error
	Ping(context.Context) error
	Close(context.Context) error
	Update(context.Context) error
}

type UnimplementedSensor struct {
	/*
	 * Shared parameters
	 */
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
	// Value is used to store of value of sensor
	Value float64 `yaml:"-"`
	// UpdatedAt is used to store moment when sensor was updated last time
	UpdatedAt time.Time `yaml:"-"`
	// Error is used to store most recent error of sensor update
	Error error
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
	// Mutex protects
	Mutex *sync.RWMutex `yaml:"-"`
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
	if u.Tags == nil {
		u.Tags = make(map[string]string, 0)
	}
	_, ok := u.Tags["type"]
	if !ok {
		u.Tags["type"] = u.Type
	}
	return u.Tags
}

func (u *UnimplementedSensor) GetValue() float64 {
	u.Mutex.RLock()
	defer u.Mutex.RUnlock()
	return u.A*u.Value + u.B
}

func (u *UnimplementedSensor) GetUpdatedAt() time.Time {
	u.Mutex.RLock()
	defer u.Mutex.RUnlock()
	return u.UpdatedAt
}

func (u *UnimplementedSensor) GetLastError() error {
	u.Mutex.RLock()
	defer u.Mutex.RUnlock()
	return u.Error
}

func (u *UnimplementedSensor) Next() time.Time {
	u.Mutex.RLock()
	defer u.Mutex.RUnlock()
	a := time.Now().Add(u.RefreshRate)
	b := u.UpdatedAt.Add(u.RefreshRate)
	if a.After(b) {
		return b
	}
	return a
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
	err = sensor.Update(ctx)
	if err != nil {
		t.Errorf("error updating: %s", err)
		return
	}
	t.Logf("Sensor updated with %.4f...", expected)
	for i = 0; i < readAttempts; i++ {
		val = sensor.GetValue()
		if val != expected {
			t.Errorf("unexpected value - %.4f vs %.4f on %v run",
				sensor.GetValue(), val, i)
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
