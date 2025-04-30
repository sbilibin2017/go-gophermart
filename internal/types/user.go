package types

import "errors"

type User struct {
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
}

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
