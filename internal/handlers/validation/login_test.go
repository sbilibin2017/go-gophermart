package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateLogin(t *testing.T) {
	validate := validator.New()
	tests := []struct {
		login   string
		isValid bool
	}{
		{"validLogin123", true},
		{"123", true},
		{"abc", true},
		{"abc123", true},
		{"", false},
		{"ab", false},
		{"a@bc123", false},
		{"longerlogin123456", true},
		{"short", true},
		{"short123", true},
	}
	for _, tt := range tests {
		t.Run(tt.login, func(t *testing.T) {
			type User struct {
				Login string `validate:"login"`
			}
			validate.RegisterValidation("login", ValidateLogin)
			user := User{Login: tt.login}
			err := validate.Struct(user)
			if tt.isValid {
				assert.NoError(t, err, "expected no error for login: %s", tt.login)
			} else {
				assert.Error(t, err, "expected error for login: %s", tt.login)
			}
		})
	}
}
