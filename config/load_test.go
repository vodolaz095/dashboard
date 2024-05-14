package config

import (
	"encoding/json"
	"testing"
)

var testConfig Config

func TestLoadFromFile(t *testing.T) {
	cfg, err := LoadFromFile("./../contrib/dashboard.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	raw, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(raw))
}

func TestConfig_Dump(t *testing.T) {
	data, err := testConfig.Dump()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(data))
}
