package types

import "errors"

var (
	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrUnauthorized       = errors.New("unauthorized")
)
