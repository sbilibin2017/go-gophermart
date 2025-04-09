package validators

import (
	"regexp"
	"unicode"

	"github.com/sbilibin2017/go-gophermart/internal/errors"
)

func ValidateUserLogin(login string) error {
	if !regexp.MustCompile("^[a-zA-Z0-9]{3,}$").MatchString(login) {
		return errors.ErrInvalidLogin
	}
	return nil
}

func ValidateUserPassword(password string) error {
	flags := []bool{false, false, false, false}
	if len(password) >= 8 {
		flags[0] = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			flags[1] = true
		case unicode.IsDigit(char):
			flags[2] = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			flags[3] = true
		}
	}
	for _, ok := range flags {
		if !ok {
			return errors.ErrInvalidPassword
		}
	}
	return nil
}
