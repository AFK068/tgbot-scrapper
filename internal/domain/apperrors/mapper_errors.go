package apperrors

type LinkValidateError struct {
	Message string
}

func (e *LinkValidateError) Error() string {
	return e.Message
}

type LinkTypeError struct {
	Message string
}

func (e *LinkTypeError) Error() string {
	return e.Message
}
