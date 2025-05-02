package dto

type UpdateSensorRequest struct {
	Name        string  `json:"name" form:"name"`
	Value       float64 `json:"value" form:"value"`
	Description string  `json:"description" form:"description"`
}
