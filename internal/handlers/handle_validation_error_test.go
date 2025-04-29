package handlers

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestHandleValidationError(t *testing.T) {
	validate := validator.New()

	type TestStruct struct {
		Login    string `validate:"required"`
		Password string `validate:"required"`
	}

	tests := []struct {
		name           string
		input          TestStruct
		expectedOutput string
	}{
		{
			name: "Both fields are missing",
			input: TestStruct{
				Login:    "",
				Password: "",
			},
			expectedOutput: "field 'Login' failed validation: required",
		},
		{
			name: "Only Login is missing",
			input: TestStruct{
				Login:    "",
				Password: "validPassword",
			},
			expectedOutput: "field 'Login' failed validation: required",
		},
		{
			name: "Only Password is missing",
			input: TestStruct{
				Login:    "validLogin",
				Password: "",
			},
			expectedOutput: "field 'Password' failed validation: required",
		},
		{
			name: "Both fields are valid",
			input: TestStruct{
				Login:    "validLogin",
				Password: "validPassword",
			},
			expectedOutput: "",
		},
		{
			name: "Unexpected error (not validation error)",
			input: TestStruct{
				Login:    "validLogin",
				Password: "validPassword",
			},
			expectedOutput: "validation failed", // Для других ошибок возвращаем "Validation failed"
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			if tt.name == "Unexpected error (not validation error)" {
				err = fmt.Errorf("some unexpected error")
			} else {
				err = validate.Struct(tt.input)
			}

			actualOutput := handleValidationError(err)

			if actualOutput == nil {
				assert.Equal(t, tt.expectedOutput, "")
			} else {
				assert.Equal(t, tt.expectedOutput, actualOutput.Error())
			}
		})
	}
}
