package validators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginValidator_Validate(t *testing.T) {
	validator := NewLoginValidator()

	validLogins := []string{
		"user123",
		"valid_user_1",
		"username-123",
	}

	for _, login := range validLogins {
		t.Run("valid login: "+login, func(t *testing.T) {
			err := validator.Validate(login)
			assert.NoError(t, err, "Expected no error for valid login")
		})
	}

	invalidLogins := []struct {
		login       string
		expectedErr error
	}{
		{"ab", ErrInvalidLogin},
		{"thisisaverylongloginname", ErrInvalidLogin},
		{"user name", ErrInvalidLogin},
		{"user@name", ErrInvalidLogin},
		{"user#123", ErrInvalidLogin},
		{"тест", ErrInvalidLogin},
	}

	for _, test := range invalidLogins {
		t.Run("invalid login: "+test.login, func(t *testing.T) {
			err := validator.Validate(test.login)
			assert.NotNil(t, err)
		})
	}
}
