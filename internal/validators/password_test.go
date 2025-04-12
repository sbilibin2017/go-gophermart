package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordValidator_Validate(t *testing.T) {
	validator := NewPasswordValidator()

	tests := []struct {
		password string
		expected error
	}{
		{"short", ErrInvalidPassword},
		{"ThisIsAVeryLongPassword123!", ErrInvalidPassword},
		{"Password with space1!", ErrInvalidPassword},
		{"password1!", ErrInvalidPassword},
		{"Password!", ErrInvalidPassword},
		{"Password1", ErrInvalidPassword},
		{"ValidPass1!", nil},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			err := validator.Validate(tt.password)
			assert.Equal(t, tt.expected, err)
		})
	}
}
