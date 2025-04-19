package handlers

import (
	"errors"
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type TestRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     int64  `json:"reward" validate:"required,gt=0"`
	RewardType string `json:"reward_type" validate:"required,oneof=% pt"`
}

func TestBuildValidationErrorMessage(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name     string
		request  TestRequest
		expected string
	}{
		{
			name: "Test with missing required fields",
			request: TestRequest{
				Match:      "",
				Reward:     0,
				RewardType: "",
			},
			expected: "Match: cannot be blank, Reward: cannot be blank, RewardType: cannot be blank",
		},
		{
			name: "Test with invalid Reward value",
			request: TestRequest{
				Match:      "Match1",
				Reward:     -1,
				RewardType: "pt",
			},
			expected: "Reward: must be greater than 0",
		},
		{
			name: "Test with invalid RewardType value",
			request: TestRequest{
				Match:      "Match1",
				Reward:     10,
				RewardType: "invalid",
			},
			expected: "RewardType: must be one of % pt",
		},
		{
			name: "Test with all valid fields",
			request: TestRequest{
				Match:      "Match1",
				Reward:     10,
				RewardType: "%",
			},
			expected: "",
		},
		{
			name: "Test with multiple errors on the same field",
			request: TestRequest{
				Match:      "",
				Reward:     -1,
				RewardType: "invalid",
			},
			expected: "Match: cannot be blank, Reward: must be greater than 0, RewardType: must be one of % pt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.request)
			var validationErrStr string
			if err != nil {
				validationErrStr = buildValidationError(err).Error()
			}
			assert.Equal(t, tt.expected, validationErrStr)
		})
	}
}

func TestBuildValidationErrorMessage_GeneralError(t *testing.T) {
	err := errors.New("general error occurred")
	expected := "general error occurred"
	result := buildValidationError(err).Error()
	assert.Equal(t, expected, result)

	err = fmt.Errorf("formatted error: %s", "something went wrong")
	expected = "formatted error: something went wrong"
	result = buildValidationError(err).Error()
	assert.Equal(t, expected, result)
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Test with a regular string",
			input:    "hello",
			expected: "Hello",
		},
		{
			name:     "Test with an empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Test with a single character",
			input:    "a",
			expected: "A",
		},
		{
			name:     "Test with already capitalized string",
			input:    "Hello",
			expected: "Hello",
		},
		{
			name:     "Test with string containing numbers",
			input:    "123abc",
			expected: "123abc",
		},
		{
			name:     "Test with string containing special characters",
			input:    "!hello",
			expected: "!hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := capitalize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
