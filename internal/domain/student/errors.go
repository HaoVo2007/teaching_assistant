package student

type Error string

const (
	ErrInvalidStudentCode    Error = "invalid student code"
	ErrInvalidParentID       Error = "invalid parent id"
	ErrStudentNotFound       Error = "student not found"
	ErrGuardianAlreadyExists Error = "guardian already exists"
	ErrGuardianNotFound      Error = "guardian not found"
	ErrStudentNotInClass     Error = "student does not belong to this class"
	ErrStudentInactive       Error = "student is inactive"
	ErrParentNotFound        Error = "parent not found"
	ErrNotAParent            Error = "user is not a parent"
)

func (e Error) Error() string {
	return string(e)
}
