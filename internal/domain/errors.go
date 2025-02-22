package domain

type ChatIsNotExistError struct {
	Message string
}

func (e *ChatIsNotExistError) Error() string {
	return e.Message
}

type ChatAlreadyExistError struct {
	Message string
}

func (e *ChatAlreadyExistError) Error() string {
	return e.Message
}

type LinkIsNotExistError struct {
	Message string
}

func (e *LinkIsNotExistError) Error() string {
	return e.Message
}
