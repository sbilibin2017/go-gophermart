package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func handleValidationError(err error) error {
	if err == nil {
		return nil
	}
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, validationErr := range validationErrors {
			field := validationErr.Field()
			tag := validationErr.Tag()
			return fmt.Errorf("field '%s' failed validation: %s", field, tag)
		}
	}
	return fmt.Errorf("validation failed")
}
