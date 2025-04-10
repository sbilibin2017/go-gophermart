package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidLogin(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("login", ValidateLogin)

	validLogins := []string{
		"abc123",
		"xyz4567",
	}

	for _, login := range validLogins {
		t.Run(login, func(t *testing.T) {
			err := v.Var(login, "login")
			assert.Nil(t, err, "expected valid login but got an error")
		})
	}
}

func TestInvalidLoginShort(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("login", ValidateLogin)

	invalidLogin := "ab"
	err := v.Var(invalidLogin, "login")
	assert.NotNil(t, err, "expected invalid login but got no error")
}

func TestInvalidLoginSpecialChars(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("login", ValidateLogin)

	invalidLogin := "@#$"
	err := v.Var(invalidLogin, "login")
	assert.NotNil(t, err, "expected invalid login but got no error")
}

func TestInvalidLoginWithSpecialChars(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("login", ValidateLogin)

	invalidLogin := "abc@123"
	err := v.Var(invalidLogin, "login")
	assert.NotNil(t, err, "expected invalid login but got no error")
}

func TestValidPassword(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("password", ValidatePassword)

	validPasswords := []string{
		"Password1@",
		"Pass1234@",
		"Strong#123",
	}

	for _, password := range validPasswords {
		t.Run(password, func(t *testing.T) {
			err := v.Var(password, "password")
			assert.Nil(t, err, "expected valid password but got an error")
		})
	}
}

func TestInvalidPasswordShort(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("password", ValidatePassword)

	invalidPassword := "short"
	err := v.Var(invalidPassword, "password")
	assert.NotNil(t, err, "expected invalid password but got no error")
}

func TestInvalidPasswordNoUppercase(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("password", ValidatePassword)

	invalidPassword := "nopper123@"
	err := v.Var(invalidPassword, "password")
	assert.NotNil(t, err, "expected invalid password but got no error")
}

func TestInvalidPasswordNoLowercase(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("password", ValidatePassword)

	invalidPassword := "NOLOWER123@"
	err := v.Var(invalidPassword, "password")
	assert.NotNil(t, err, "expected invalid password but got no error")
}

func TestInvalidPasswordNoSpecialChar(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("password", ValidatePassword)

	invalidPassword := "password123"
	err := v.Var(invalidPassword, "password")
	assert.NotNil(t, err, "expected invalid password but got no error")
}

func TestInvalidPasswordOnlyUppercase(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("password", ValidatePassword)

	invalidPassword := "PASSWORD@"
	err := v.Var(invalidPassword, "password")
	assert.NotNil(t, err, "expected invalid password but got no error")
}

func TestInvalidPasswordNoDigit(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("password", ValidatePassword)

	invalidPassword := "Password@"
	err := v.Var(invalidPassword, "password")

	assert.NotNil(t, err, "expected invalid password due to missing digit but got no error")
}

func TestInvalidPasswordNoSspecialChar(t *testing.T) {
	v := validator.New()
	v.RegisterValidation("password", ValidatePassword)

	invalidPassword := "Password123"
	err := v.Var(invalidPassword, "password")

	assert.NotNil(t, err, "expected invalid password due to missing special character but got no error")
}
