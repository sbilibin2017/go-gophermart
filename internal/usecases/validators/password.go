package validators

import (
	"errors"
	"regexp"
	"strings"
)

type PasswordValidator struct{}

func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{}
}

func (v *PasswordValidator) Validate(password string) error {
	if strings.Contains(password, " ") {
		return ErrInvalidPassword
	}

	const minLength = 8
	const maxLength = 20
	if len(password) < minLength {
		return ErrInvalidPassword
	}
	if len(password) > maxLength {
		return ErrInvalidPassword
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return ErrInvalidPassword
	}
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return ErrInvalidPassword
	}
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasSpecial {
		return ErrInvalidPassword
	}

	return nil
}

var ErrInvalidPassword = errors.New("invalid password")
