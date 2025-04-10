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

	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false
	}
	if !hasLower {
		return false
	}
	if !hasDigit {
		return false
	}
	if !hasSpecial {
		return false
	}

	return true
}
