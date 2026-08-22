package user

import "errors"

var (
	// lookup / uniqueness
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")

	// validation input
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrInvalidUsername = errors.New("invalid username")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidRole     = errors.New("invalid role")

	// auth / credentials
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWrongPassword      = errors.New("wrong password")
	ErrPasswordMismatch   = errors.New("password confirmation does not match")
	ErrPasswordTooWeak    = errors.New("password is too weak")

	// account state
	ErrUserInactive      = errors.New("user is inactive")
	ErrUserBanned        = errors.New("user is banned")
	ErrUserAlreadyActive = errors.New("user is already active")

	// authorization / business rules
	ErrUnauthorized           = errors.New("unauthorized")
	ErrForbidden              = errors.New("forbidden")
	ErrCannotDeleteSelf       = errors.New("cannot delete own account")
	ErrCannotChangeOwnRole    = errors.New("cannot change own role")
	ErrInsufficientPermission = errors.New("insufficient permission")
)
