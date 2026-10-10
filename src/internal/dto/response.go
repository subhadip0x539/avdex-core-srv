package dto

type Response[T any] struct {
	Body    T      `json:"body,omitempty"`
	Message string `json:"message,omitempty"`
}
