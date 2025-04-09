package errors

import "errors"

var (
	ErrInvalidLogin      = errors.New("invalid login")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserAlreadyExists = errors.New("user already exists")
)
