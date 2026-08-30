package user

type Error string

const (
	ErrUserNotFound       Error = "user not found"
	ErrUserAlreadyExists  Error = "user already exists"
	ErrEmailAlreadyExists Error = "email already exists"

	ErrInvalidEmail    Error = "invalid email"
	ErrInvalidPassword Error = "invalid password"
	ErrInvalidUsername Error = "invalid username"
	ErrInvalidName     Error = "invalid name"
	ErrInvalidRole     Error = "invalid role"

	ErrInvalidCredentials Error = "invalid credentials"
	ErrWrongPassword      Error = "wrong password"
	ErrPasswordMismatch   Error = "password confirmation does not match"
	ErrPasswordTooWeak    Error = "password is too weak"

	ErrUserInactive      Error = "user is inactive"
	ErrUserBanned        Error = "user is banned"
	ErrUserAlreadyActive Error = "user is already active"

	ErrUnauthorized           Error = "unauthorized"
	ErrForbidden              Error = "forbidden"
	ErrCannotDeleteSelf       Error = "cannot delete own account"
	ErrCannotChangeOwnRole    Error = "cannot change own role"
	ErrInsufficientPermission Error = "insufficient permission"
)
