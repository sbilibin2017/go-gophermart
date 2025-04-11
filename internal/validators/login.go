package validators

import (
	"errors"
	"regexp"
	"strings"
)

type LoginValidator struct{}

func NewLoginValidator() *LoginValidator {
	return &LoginValidator{}
}

func (v *LoginValidator) Validate(login string) error {
	const minLength = 3
	const maxLength = 20
	if len(login) < minLength {
		return ErrInvalidLogin
	}
	if len(login) > maxLength {
		return ErrInvalidLogin
	}
	if strings.Contains(login, " ") {
		return ErrInvalidLogin
	}
	match, _ := regexp.MatchString("^[a-zA-Z0-9_-]+$", login)
	if !match {
		return ErrInvalidLogin
	}
	return nil
}

var ErrInvalidLogin = errors.New("invalid login")
