package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"h", "H"},
		{"", ""},
		{"Hello", "Hello"},
		{"1hello", "1hello"},
		{"привет", "Привет"},
		{"!bang", "!bang"},
		{string([]byte{0xff, 0xfe, 0xfd}), string([]byte{0xff, 0xfe, 0xfd})},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := Capitalize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
