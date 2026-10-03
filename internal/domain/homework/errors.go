package homework

type Error string

const (
	ErrHomeworkNotFound      Error = "homework not found"
	ErrHomeworkNotAuthorized Error = "homework not authorized"
	ErrInvalidTitle          Error = "invalid title"
	ErrInvalidClassID        Error = "invalid class id"
	ErrInvalidQuestions      Error = "homework must have at least one question"
	ErrInvalidDueDate        Error = "invalid due date"
	ErrHomeworkInUse         Error = "homework has submissions"
)

func (e Error) Error() string {
	return string(e)
}
