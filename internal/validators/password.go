package usecases

import (
	"regexp"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/errors"
)

type PasswordValidator struct{}

func NewPasswordValidator() *PasswordValidator {
	return &PasswordValidator{}
}

func (v *PasswordValidator) Validate(password string) error {
	if strings.Contains(password, " ") {
		return errors.ErrInvalidPassword
	}

	const minLength = 8
	const maxLength = 20
	if len(password) < minLength {
		return errors.ErrInvalidPassword
	}
	if len(password) > maxLength {
		return errors.ErrInvalidPassword
	}

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return errors.ErrInvalidPassword
	}
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return errors.ErrInvalidPassword
	}
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasSpecial {
		return errors.ErrInvalidPassword
	}

	return nil
}
