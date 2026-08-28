package class

type Error string

const (
	ErrInvalidClass       Error = "invalid class"
	ErrClassNotFound      Error = "class not found"
	ErrClassAlreadyExists Error = "class already exists"
	ErrClassNotAuthorized Error = "class not authorized"
	ErrUnauthorized       Error = "unauthorized"
	ErrImageTooLarge      Error = "image too large"
)

func (e Error) Error() string {
	return string(e)
}
