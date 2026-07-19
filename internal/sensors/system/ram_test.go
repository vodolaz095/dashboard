package system

import (
	"testing"
)

func TestFreeRAMSensor(t *testing.T) {
	sensor := FreeRAMSensor{}

	err := sensor.Init(t.Context())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = sensor.Ping(t.Context())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = sensor.Update(t.Context())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if sensor.Error != nil {
		t.Errorf("sensor error: %v", sensor.Error)
	}
	result := sensor.GetValue()
	if result == 0 {
		t.Errorf("expected non-zero value, got 0")
	}
	if result < 0 {
		t.Errorf("expected non-negative value, got %f", result)
	}
	t.Logf("Free RAM: %f", result)

	err = sensor.Close(t.Context())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
