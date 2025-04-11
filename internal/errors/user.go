package errors

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidLogin = errors.New("invalid login")
	ErrInvalidPassword = errors.New("invalid password")
)
