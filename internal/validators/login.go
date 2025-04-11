package usecases

import (
	"regexp"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/errors"
)

type LoginValidator struct{}

func NewLoginValidator() *LoginValidator {
	return &LoginValidator{}
}

func (v *LoginValidator) Validate(login string) error {
	const minLength = 3
	const maxLength = 20
	if len(login) < minLength {
		return errors.ErrInvalidLogin
	}
	if len(login) > maxLength {
		return errors.ErrInvalidLogin
	}
	if strings.Contains(login, " ") {
		return errors.ErrInvalidLogin
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", login)
	if !match {
		return errors.ErrInvalidLogin
	}
	return nil
}
