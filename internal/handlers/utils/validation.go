package utils

import (
	"sort"
	"strings"

	"github.com/go-playground/validator/v10"
)

func BuildValidationErrorMessage(err error) string {
	var sb strings.Builder

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		fieldErrors := make(map[string][]string)

		for _, ve := range validationErrors {
			field := ve.Field()
			tag := ve.Tag()
			param := ve.Param()

			var errorMessage string

			switch tag {
			case "required":
				errorMessage = field + ": cannot be blank"
			case "oneof":
				errorMessage = field + ": must be one of " + param
			case "gt":
				errorMessage = field + ": must be greater than " + param
			}

			fieldErrors[field] = append(fieldErrors[field], errorMessage)
		}

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

	return err.Error()
}
