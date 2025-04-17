package utils

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestFormatValidationError(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name          string
		input         interface{}
		expectedError map[string]string
	}{
		{
			name: "single validation error",
			input: &struct {
				Name  string `validate:"required"`
				Email string `validate:"required,email"`
			}{
				Name:  "",
				Email: "valid@email.com",
			},
			expectedError: map[string]string{
				"Name": "failed validation for 'required'",
			},
		},
		{
			name: "multiple validation errors",
			input: &struct {
				Name  string `validate:"required"`
				Email string `validate:"required,email"`
			}{
				Name:  "",
				Email: "invalid-email", // Email некорректный
			},
			expectedError: map[string]string{
				"Name":  "failed validation for 'required'",
				"Email": "failed validation for 'email'",
			},
		},
		{
			name:          "non-validation error",
			input:         fmt.Errorf("some non-validation error"),
			expectedError: map[string]string{"error": "unknown validation error: some non-validation error"},
		},
		{
			name:          "no validation errors",
			input:         &struct{ Name, Email string }{Name: "John Doe", Email: "john.doe@example.com"},
			expectedError: map[string]string{"error": "unknown validation error: <nil>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch v := tt.input.(type) {
			case error:
				err = v
			default:
				err = validate.Struct(v)
			}

			formattedError := FormatValidationError(err)
			assert.Equal(t, tt.expectedError, formattedError)
		})
	}
}

func TestValidationErrorsToString(t *testing.T) {
	tests := []struct {
		name     string
		errors   map[string]string
		expected string
	}{
		{
			name: "Single error",
			errors: map[string]string{
				"Name": "Name is required",
			},
			expected: "Name: Name is required\n",
		},
		{
			name: "Multiple errors",
			errors: map[string]string{
				"Name":     "Name is required",
				"Email":    "Email is invalid",
				"Password": "Password must be at least 8 characters",
			},
			expected: "Email: Email is invalid\nName: Name is required\nPassword: Password must be at least 8 characters\n",
		},
		{
			name:     "Empty errors",
			errors:   map[string]string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidationErrorsToString(tt.errors)
			assert.Equal(t, tt.expected, result)
		})
	}
}
