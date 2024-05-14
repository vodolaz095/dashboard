package model

import "time"

// Sensor is data transfer object for http, json and metrics visualization
type Sensor struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Minimum     float64   `json:"minimum"`
	Maximum     float64   `json:"maximum"`
	Value       float64   `json:"value"`
	Error       string    `json:"error"`
	Tags        []string  `json:"tags"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
