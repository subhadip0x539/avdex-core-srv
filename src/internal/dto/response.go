package dto

type Response[T any] struct {
	Body    T      `json:"body"`
	Message string `json:"message,omitempty"`
}
