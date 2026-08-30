package question

type Error string

const (
	ErrInvalidType      Error = "invalid question type"
	ErrInvalidQuestion  Error = "invalid question"
	ErrInvalidOptions   Error = "invalid options"
	ErrInvalidCorrect   Error = "invalid correct answer"
	ErrInvalidPairs     Error = "invalid matching pairs"
	ErrImageTooLarge    Error = "image too large"
	ErrQuestionNotFound Error = "question not found"
	ErrInvalidSubject   Error = "invalid subject"
	ErrInvalidGrade     Error = "invalid grade"
	ErrUnauthorized     Error = "unauthorized"
	ErrQuestionInUse    Error = "question is used in a set, homework, or submission"
)

func (e Error) Error() string {
	return string(e)
}
