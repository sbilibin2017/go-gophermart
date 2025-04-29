package domain

import "errors"

type Order struct {
	Number string
}

var (
	ErrUserOrderExists = errors.New("user order exists")
	ErrOrderExists     = errors.New("order exists")
)
