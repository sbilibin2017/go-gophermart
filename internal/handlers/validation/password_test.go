package validation

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("password", ValidatePassword)
	tests := []struct {
		password string
		isValid  bool
	}{
		{"Valid1Password@", true},
		{"short1@", false},
		{"NoSpecialChar123", false},
		{"NOLOWERCASE123@", false},
		{"validpassword123@", false},
		{"ValidPass123", false},
		{"@Valid123", true},
		{"Valid@1234", true},
		{"12345678", false},
		{"PASSWORD@123", false},
		{"password@123", false},
	}
	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			type User struct {
				Password string `validate:"password"`
			}
			user := User{Password: tt.password}
			err := validate.Struct(user)
			if tt.isValid {
				assert.NoError(t, err, "expected no error for password: %s", tt.password)
			} else {
				assert.Error(t, err, "expected error for password: %s", tt.password)
			}
		})
	}
}
