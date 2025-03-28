package stackoverflow

import "errors"

var (
	ErrFailedToGetQuestion = errors.New("failed to get question")
	ErrQuestionNotFound    = errors.New("question not found")
	ErrNoAnswersFound      = errors.New("no answers found")
	ErrFailedToGetItems    = errors.New("failed to get items")
	ErrInvalidQuestionURL  = errors.New("invalid question url")
)
