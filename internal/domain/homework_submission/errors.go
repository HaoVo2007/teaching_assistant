package homeworksubmission

type Error string

const (
	ErrInvalidHomeworkSubmission        Error = "invalid homework submission"
	ErrHomeworkSubmissionNotFound       Error = "homework submission not found"
	ErrStudentNotInClass                Error = "student not in class"
	ErrHomeworkSubmissionAlreadyExists  Error = "homework submission already exists"
	ErrQuestionMismatch                 Error = "student answers do not match homework questions"
	ErrInvalidStudentAnswer             Error = "invalid student answer for question type"
	ErrHomeworkSubmissionDueDateExpired Error = "homework submission due date expired"
)

func (e Error) Error() string {
	return string(e)
}
