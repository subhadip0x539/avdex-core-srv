package domain

import "net/http"

type ErrorMessage string

const (
	ErrorMessageInternalServerError ErrorMessage = "Internal Server Error"
	ErrorMessageBadRequest          ErrorMessage = "Bad Request"
	ErrorMessageNotFound            ErrorMessage = "Not Found"
)

type Error struct {
	Code    int
	Message ErrorMessage
	Trace   error
}

func NewBadRequest(trace error) *Error {
	return &Error{Code: http.StatusBadRequest, Message: ErrorMessageBadRequest, Trace: trace}
}

func NewInternalServerError(trace error) *Error {
	return &Error{Code: http.StatusInternalServerError, Message: ErrorMessageInternalServerError, Trace: trace}
}
