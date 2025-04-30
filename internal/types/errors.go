package types

import "errors"

var (
	ErrJSONDecode          = errors.New("invalid json")
	ErrInternalServerError = errors.New("internal server error")
	ErrTokenEncode         = errors.New("token is not encoded")
	ErrUnauthorized        = errors.New("unauthorized")
)
