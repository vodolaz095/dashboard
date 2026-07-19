package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vodolaz095/dashboard/internal/sensors"
)

func TestFreeDiskSpaceSensor(t *testing.T) {
	s := FreeDiskSpaceSensor{diskSpaceSensor{Path: "/"}}
	val, err := sensors.DoGetSensorValue(t, &s)
	if err != nil {
		t.Errorf("error getting sensor value: %s", err)
		return
	}
	assert.Greater(t, val, float64(0))
	t.Logf("Free : %v", val)
}

func TestUsedDiskSpaceSensor(t *testing.T) {
	s := UsedDiskSpaceSensor{diskSpaceSensor{Path: "/"}}
	val, err := sensors.DoGetSensorValue(t, &s)
	if err != nil {
		t.Errorf("error getting sensor value: %s", err)
		return
	}
	assert.Greater(t, val, float64(0))
	t.Logf("Used : %v", val)
}

func TestFreeDiskSpaceRatioSensor(t *testing.T) {
	s := FreeDiskSpaceRatioSensor{diskSpaceSensor{Path: "/"}}
	val, err := sensors.DoGetSensorValue(t, &s)
	if err != nil {
		t.Errorf("error getting sensor value: %s", err)
		return
	}
	assert.Greater(t, val, float64(0))
	assert.Less(t, val, float64(100))
	t.Logf("Ratio : %.2f%%", val)
}
