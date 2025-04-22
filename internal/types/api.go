package types

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type APIResponse[T any] struct {
	Data    T      `json:"data,omitempty"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
