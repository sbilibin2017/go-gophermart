package helpers

import (
	"errors"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
)

func buildValidationError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := groupErrorsByField(validationErrors)
		errorMessages := buildErrorMessages(fieldErrors)
		return errors.New(errorMessages)
	}
	return err
}

func groupErrorsByField(validationErrors validator.ValidationErrors) map[string][]string {
	fieldErrors := make(map[string][]string)
	for _, ve := range validationErrors {
		errorMessage := createErrorMessage(ve)
		fieldErrors[ve.Field()] = append(fieldErrors[ve.Field()], errorMessage)
	}
	return fieldErrors
}

func createErrorMessage(ve validator.FieldError) string {
	var errorMessage string
	switch ve.Tag() {
	case "required":
		errorMessage = ve.Field() + ": cannot be blank"
	case "oneof":
		errorMessage = ve.Field() + ": must be one of " + ve.Param()
	case "gt":
		errorMessage = ve.Field() + ": must be greater than " + ve.Param()
	}
	return errorMessage
}

func buildErrorMessages(fieldErrors map[string][]string) string {
	var sb strings.Builder
	fields := make([]string, 0, len(fieldErrors))
	for field := range fieldErrors {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, field := range fields {
		errors := fieldErrors[field]
		for _, errorMessage := range errors {
			sb.WriteString(errorMessage)
		}
		sb.WriteString(", ")
	}
	result := sb.String()
	if len(result) > 2 {
		result = result[:len(result)-2]
	}
	return result
}
