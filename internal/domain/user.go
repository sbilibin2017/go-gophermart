package domain

import "errors"

type User struct {
	Login    string
	Password string
}

var (
	ErrLoginAlreadyTaken      = errors.New("login is already taken")
	ErrInvalidUserCredentials = errors.New("invalid user credentials")
)
