package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type validationError struct {
	Field   string
	Tag     string
	Message string
}

func formatValidationError(err error) *validationError {
	if errs, ok := err.(validator.ValidationErrors); ok {
		e := errs[0]

		field := camelToSnake(e.Field())

		var message string
		switch e.Tag() {
		case "required":
			message = fmt.Sprintf("Field %s is required", field)
		case "gt":
			message = fmt.Sprintf("Field %s must be greater than %s", field, e.Param())
		case "oneof":
			message = fmt.Sprintf("Field %s must be one of [%s]", field, e.Param())
		default:
			message = fmt.Sprintf("Field %s is invalid (%s)", field, e.Tag())
		}

		return &validationError{
			Field:   field,
			Tag:     e.Tag(),
			Message: message,
		}
	}
	return nil
}

func camelToSnake(s string) string {
	snake := regexp.MustCompile("([a-z0-9])([A-Z])").ReplaceAllString(s, "${1}_${2}")
	return strings.ToLower(snake)
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

var (
	errInvalidJSONFormat = "Invalid JSON format"
)
