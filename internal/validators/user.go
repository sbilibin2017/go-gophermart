package validators

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func ValidateLogin(fl validator.FieldLevel) bool {
	login := fl.Field().String()
	match := regexp.MustCompile(`^[a-zA-Z0-9]{3,}$`)
	return match.MatchString(login)
}

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
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
			return false
		}
	}
	return true
}
