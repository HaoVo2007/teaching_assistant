package user

type Error string

const (
	ErrUserNotFound       Error = "user not found"
	ErrEmailAlreadyExists Error = "email already exists"
	ErrInvalidEmail       Error = "invalid email"
	ErrInvalidPassword    Error = "invalid password"
	ErrInvalidUsername    Error = "invalid username"
	ErrInvalidRole        Error = "invalid role"
	ErrInvalidCredentials Error = "invalid credentials"
	ErrUnauthorized       Error = "unauthorized"
)

func (e Error) Error() string {
	return string(e)
}
