package model

import "time"

type Update struct {
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}
