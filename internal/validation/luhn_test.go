package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateNumber(t *testing.T) {
	tests := []struct {
		name          string
		orderNumber   string
		expectedValid bool
	}{
		{
			name:          "Valid Luhn number",
			orderNumber:   "1234567812345670", // Пример действительного номера по алгоритму Луна
			expectedValid: true,
		},
		{
			name:          "Invalid Luhn number (wrong checksum)",
			orderNumber:   "12345678901", // Пример недействительного номера по алгоритму Луна
			expectedValid: false,
		},
		{
			name:          "Valid Luhn number (even length)",
			orderNumber:   "1234567812345670", // Еще один пример действительного номера
			expectedValid: true,
		},
		{
			name:          "Invalid Luhn number (incorrect last digit)",
			orderNumber:   "1234567812345671", // Пример недействительного номера
			expectedValid: false,
		},
		{
			name:          "Empty string",
			orderNumber:   "", // Пустая строка
			expectedValid: false,
		},
		{
			name:          "Single digit",
			orderNumber:   "5", // Один символ (недопустимо)
			expectedValid: false,
		},
		{
			name:          "Valid Luhn number (long even length)",
			orderNumber:   "6011514433546201", // Пример действительного номера
			expectedValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := validateNumber(tt.orderNumber)
			assert.Equal(t, tt.expectedValid, valid, "Expected validity of %s to be %v", tt.orderNumber, tt.expectedValid)
		})
	}
}
