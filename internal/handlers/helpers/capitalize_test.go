package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
