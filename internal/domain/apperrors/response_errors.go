package apperrors

import "fmt"

type ErrorResponse struct {
	Code    int
	Message string
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}
