package system

import (
	"testing"

	"github.com/vodolaz095/dashboard/internal/sensors"
)

func TestFreeRAMSensor(t *testing.T) {
	s := FreeRAMSensor{}
	val, err := sensors.DoGetSensorValue(t, &s)
	if err != nil {
		t.Errorf("error getting sensor value: %s", err)
		return
	}
	// assert.Greater(t, val, float64(0))
	t.Logf("Free ram in Mbytes : %v", val)
}
