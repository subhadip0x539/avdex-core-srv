package domain

type ErrorMessage string

const (
	ErrorMessageInternalServerError ErrorMessage = "Internal Server Error"
	ErrorMessageBadRequest          ErrorMessage = "Bad Request"
	ErrorMessageNotFound            ErrorMessage = "Not Found"
)

type Error struct {
	Code    int
	Message string
	Trace   error
}
