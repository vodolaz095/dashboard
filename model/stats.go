package model

type Stats struct {
	SensorsUpdatedNow int `json:"sensors_updated_now"`
	QueueLength       int `json:"queue_length"`
	Subscribers       int `json:"subscribers"`
}
