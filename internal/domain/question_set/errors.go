package questionset

type Error string

const (
	ErrQuestionSetNotFound            Error = "question set not found"
	ErrInvalidTitle                   Error = "invalid title"
	ErrInvalidQuestionType            Error = "invalid question type"
	ErrInvalidQuestions               Error = "invalid questions"
	ErrInvalidQuestionTypeForQuestion Error = "invalid question type for question"
	ErrUnauthorized                   Error = "unauthorized"
)

func (e Error) Error() string {
	return string(e)
}
