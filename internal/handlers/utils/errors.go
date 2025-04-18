package utils

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrInternal           = errors.New("internal error")
)
