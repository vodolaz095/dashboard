package model

import (
	"encoding/json"
	"time"
)

type Update struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

func (u *Update) Pack() string {
	raw, _ := json.Marshal(u)
	return string(raw)
}
