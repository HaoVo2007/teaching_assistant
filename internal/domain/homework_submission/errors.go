package homeworksubmission

type Error string

const (
	ErrInvalidHomeworkSubmission       Error = "invalid homework submission"
	ErrHomeworkSubmissionNotFound      Error = "homework submission not found"
	ErrHomeworkSubmissionAlreadyExists Error = "homework submission already exists"
	ErrHomeworkSubmissionInvalid       Error = "homework submission invalid"
	ErrHomeworkSubmissionInternal      Error = "homework submission internal"
	ErrQuestionMismatch                Error = "student answers do not match homework questions"
	ErrInvalidStudentAnswer            Error = "invalid student answer for question type"
)
