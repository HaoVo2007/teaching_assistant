package class

type Error string

const (
	ErrInvalidClass      Error = "invalid class"
	ErrClassNotFound     Error = "class not found"
	ErrUnauthorized      Error = "unauthorized"
	ErrImageTooLarge     Error = "image too large"
	ErrClassInUse        Error = "class has homeworks"
	ErrStudentCodeExists Error = "student code already exists"
)

func (e Error) Error() string {
	return string(e)
}
