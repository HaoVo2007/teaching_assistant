package student

type Error string

const (
	ErrInvalidStudentCode    Error = "invalid student code"
	ErrInvalidParentID       Error = "invalid parent id"
	ErrStudentNotFound       Error = "student not found"
	ErrGuardianAlreadyExists Error = "guardian already exists"
	ErrGuardianNotFound      Error = "guardian not found"
	ErrStudentNotInClass     Error = "student does not belong to this class"
)

func (e Error) Error() string {
	return string(e)
}
