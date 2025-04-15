package dtos

type HealthResponse struct {
	Active bool   `json:"active"`
	Status string `json:"status"`
}
