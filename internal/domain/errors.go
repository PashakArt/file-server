package domain

import "errors"

var (
	// validate
	ErrLoginInvalid    = errors.New("Login must be at least 8 chars, digit latin only")
	ErrPasswordInvalid = errors.New("Password must be at least 8 chars, contain lower, upper, digit, and special char")
	ErrTokenInvalid    = errors.New("Invalid admin token")

	ErrUserAlreadyExists  = errors.New("User with this login already exists")
	ErrUserNotFound       = errors.New("User not found")
	ErrInvalidCredentials = errors.New("Invalid login or password")
)
