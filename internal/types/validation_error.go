package types

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field string
	Tag   string
	Error error
}

func NewValidationError(err error) *ValidationError {
	var ve validator.ValidationErrors
	if errs, ok := err.(validator.ValidationErrors); ok {
		ve = errs
	} else {
		return nil
	}

	if len(ve) == 0 {
		return nil
	}

	for _, fieldError := range ve {
		validationErr := &ValidationError{
			Field: fieldError.Field(),
			Tag:   fieldError.Tag(),
			Error: fmt.Errorf(
				"field '%s' failed validation: %s",
				fieldError.Field(),
				fieldError.Tag(),
			),
		}
		return validationErr
	}

	return nil
}

type ValidationWithStatusCode struct {
	ValidationError
	StatusCode int
}

func NewValidationWithStatusCode(
	err ValidationError,
	statusCode int,
) *ValidationWithStatusCode {
	return &ValidationWithStatusCode{
		ValidationError: err,
		StatusCode:      statusCode,
	}
}
