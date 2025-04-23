package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestLuhnValidator(t *testing.T) {
	validate := validator.New()
	RegisterLuhnValidator(validate)

	tests := []struct {
		name     string
		input    string
		expected error
	}{
		{
			name:     "Valid Luhn number",
			input:    "1234567812345670",
			expected: nil,
		},
		{
			name:     "Invalid Luhn number",
			input:    "1234567812345671",
			expected: assert.AnError, // Подходит под ошибку, если невалидный
		},
		{
			name:     "Empty string",
			input:    "",
			expected: assert.AnError, // Ожидается ошибка
		},
		{
			name:     "Luhn number with special characters",
			input:    "1234-5678@1234#5670",
			expected: nil,
		},
		{
			name:     "Too short number",
			input:    "1",
			expected: assert.AnError, // Ожидается ошибка, так как номер слишком короткий
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.input, "luhn")
			if tt.expected == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "luhn")
			}
		})
	}
}
