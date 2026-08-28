package questionset

type Error string

const (
	ErrQuestionSetNotFound            Error = "question set not found"
	ErrQuestionSetAlreadyExists       Error = "question set already exists"
	ErrQuestionSetNotAuthorized       Error = "question set not authorized"
	ErrInvalidTitle                   Error = "invalid title"
	ErrInvalidQuestionType            Error = "invalid question type"
	ErrInvalidQuestions               Error = "invalid questions"
	ErrInvalidQuestionTypeForQuestion Error = "invalid question type for question"
	ErrUnauthorized                   Error = "unauthorized"
)
