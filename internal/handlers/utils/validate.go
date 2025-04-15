package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Struct(i interface{}) error
}

func Validate(v Validator, obj interface{}) error {
	if err := v.Struct(obj); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var validationMessages []string
			for _, ve := range validationErrors {
				validationMessages = append(validationMessages, ve.Error())
			}
			return fmt.Errorf("validation failed: %s", strings.Join(validationMessages, ", "))
		}
		return fmt.Errorf("unexpected validation error: %v", err)
	}
	return nil
}
