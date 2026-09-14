package domain

import "errors"

var (
	// validate
	ErrLoginInvalid      = errors.New("login must be at least 8 chars, digit latin only")
	ErrPasswordInvalid   = errors.New("pswd must be at least 8 chars, contain lower, upper, digit, and special char")
	ErrTokenAdminInvalid = errors.New("Invalid admin token")
	ErrTokenInvalid      = errors.New("Invalid token")

	ErrUserAlreadyExists  = errors.New("User with this login already exists")
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Invalid login or password")
)
