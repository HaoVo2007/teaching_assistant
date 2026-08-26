package questionset

type Error string

const (
	ErrQuestionSetNotFound      Error = "question set not found"
	ErrQuestionSetAlreadyExists Error = "question set already exists"
	ErrQuestionSetNotAuthorized Error = "question set not authorized"
	ErrInvalidTitle             Error = "invalid title"
	ErrInvalidQuestions         Error = "invalid questions"
)
