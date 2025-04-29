package domain

import "errors"

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var (
	ErrLoginAlreadyTaken  = errors.New("login is already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
