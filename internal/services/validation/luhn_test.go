package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestLuhnValidator(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("luhn", LuhnValidator)

	tests := []struct {
		name      string
		order     string
		expectErr bool
	}{
		{
			name:      "Valid card number with no spaces",
			order:     "4532015112830366",
			expectErr: false,
		},
		{
			name:      "Valid card number with spaces",
			order:     "4532 0151 1283 0366",
			expectErr: false,
		},
		{
			name:      "Card number with non-numeric characters",
			order:     "4532-0151-1283-0366",
			expectErr: true,
		},
		{
			name:      "Short invalid card number",
			order:     "12345",
			expectErr: true,
		},
		{
			name:      "Empty string",
			order:     "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Var(tt.order, "luhn")
			if tt.expectErr {
				assert.Error(t, err, "Expected an error for input: %s", tt.order)
			} else {
				assert.NoError(t, err, "Expected no error for input: %s", tt.order)
			}
		})
	}
}
