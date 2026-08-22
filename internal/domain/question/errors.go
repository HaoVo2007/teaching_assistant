package question

import "errors"

var (
	ErrInvalidType      = errors.New("invalid question type")
	ErrInvalidQuestion  = errors.New("invalid question")
	ErrInvalidOptions   = errors.New("invalid options")
	ErrInvalidCorrect   = errors.New("invalid correct answer")
	ErrInvalidPairs     = errors.New("invalid matching pairs")
	ErrImageTooLarge    = errors.New("image too large")
	ErrQuestionNotFound = errors.New("question not found")
	ErrInvalidSubject   = errors.New("invalid subject")
	ErrInvalidGrade     = errors.New("invalid grade")
)
