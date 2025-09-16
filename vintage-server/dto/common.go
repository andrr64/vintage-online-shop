package dto

type CommonResponse[T any] struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    *T     `json:"data,omitempty"`
}
