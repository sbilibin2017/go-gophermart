package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateNumberWithLouna(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid number",
			value:   "4539 1488 0343 6467", // Example Louna number
			wantErr: false,
		},
		{
			name:    "valid number without spaces",
			value:   "4539148803436467",
			wantErr: false,
		},
		{
			name:    "invalid number with incorrect checksum",
			value:   "4539 1488 0343 6468", // Incorrect checksum
			wantErr: true,
		},
		{
			name:    "empty string",
			value:   "",
			wantErr: true,
		},
		{
			name:    "non-numeric string",
			value:   "abcd efgh ijkl",
			wantErr: true,
		},
		{
			name:    "non-numeric characters with spaces",
			value:   "4539 abcd 0343 6467",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNumberWithLouna(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, ErrInvalidLounaNumber, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegisterLounaValidator(t *testing.T) {
	v := validator.New()

	RegisterLounaValidator(v)

	t.Run("valid number", func(t *testing.T) {
		err := v.Var("4539 1488 0343 6467", "louna")
		assert.NoError(t, err)
	})

	t.Run("invalid number", func(t *testing.T) {
		err := v.Var("4539 1488 0343 6468", "louna")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "louna")
	})

	t.Run("empty string", func(t *testing.T) {
		err := v.Var("", "louna")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "louna")
	})

	t.Run("non-numeric string", func(t *testing.T) {
		err := v.Var("abcd efgh ijkl", "louna")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "louna")
	})
}
