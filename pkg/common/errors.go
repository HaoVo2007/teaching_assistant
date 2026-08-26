package common

type Error string

const (
	ErrUnauthorized        Error = "unauthorized"
	ErrBadRequest          Error = "bad request"
	ErrNotFound            Error = "not found"
	ErrAlreadyExists       Error = "already exists"
	ErrInternalServerError Error = "internal server error"
)

func (e Error) Error() string {
	return string(e)
}
