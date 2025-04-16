package errors

import "errors"

var (
	ErrInternal       = errors.New("internal error")
	ErrDataIsNotValid = errors.New("data is not valid")
)
