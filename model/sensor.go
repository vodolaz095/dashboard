package model

import "time"

type Sensor struct {
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Minimum     float64   `json:"minimum"`
	Maximum     float64   `json:"maximum"`
	Value       float64   `json:"value"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
