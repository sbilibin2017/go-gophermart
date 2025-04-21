package types

import (
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewAPIError(message string, status int) *APIError {
	return &APIError{
		Status:  status,
		Message: capitalize(message),
	}
}

func NewInternalError() *APIError {
	return &APIError{
		Status:  http.StatusInternalServerError,
		Message: "Internal error",
	}
}

func NewValidationErrorResponse(err error) *APIError {
	return &APIError{
		Status:  http.StatusUnprocessableEntity,
		Message: capitalize(buildValidationError(err).Error()),
	}
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

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
