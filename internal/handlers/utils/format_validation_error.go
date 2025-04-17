package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]string {
	errorsMap := make(map[string]string)
	if err == nil {
		errorsMap["error"] = "unknown validation error: <nil>"
		return errorsMap
	}
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errorsMap[fieldErr.Field()] = fmt.Sprintf("failed validation for '%s'", fieldErr.Tag())
		}
		return errorsMap
	}
	errorsMap["error"] = fmt.Sprintf("unknown validation error: %s", err.Error())
	return errorsMap
}

func ValidationErrorsToString(errors map[string]string) string {
	var keys []string
	for key := range errors {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var result strings.Builder
	for _, key := range keys {
		result.WriteString(key + ": " + errors[key] + "\n")
	}
	return result.String()
}
