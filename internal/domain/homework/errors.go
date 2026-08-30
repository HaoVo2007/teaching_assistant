package homework

type Error string

const (
	ErrInvalidHomework       Error = "invalid homework"
	ErrHomeworkNotFound      Error = "homework not found"
	ErrHomeworkAlreadyExists Error = "homework already exists"
	ErrHomeworkNotAuthorized Error = "homework not authorized"
	ErrUnauthorized          Error = "unauthorized"
	ErrInvalidTitle          Error = "invalid title"
	ErrInvalidClassID        Error = "invalid class id"
	ErrInvalidQuestions      Error = "should be at least one question"
	ErrInvalidDueDate        Error = "invalid due date"
	ErrHomeworkInUse         Error = "homework has submissions"
)

func (e Error) Error() string {
	return string(e)
}