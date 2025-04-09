package errors

import "errors"

var (
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenIsNotParsed        = errors.New("token is not parsed")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")
	ErrTokenIsNotSigned        = errors.New("token is not signed")
)
