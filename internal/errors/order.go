package errors

import "errors"

var (
	ErrOrderRegisterInvalidRequest = errors.New("invalid order register request")
	ErrOrderAlreadyRegistered      = errors.New("order already registered")
	ErrOrderRegisterInternal       = errors.New("internal error")
)
