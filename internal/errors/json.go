package errors

import "errors"

var (
	ErrRequestDecoderUnprocessableJSON = errors.New("invalid json")
)
