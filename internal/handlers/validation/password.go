package validation

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check for minimum length
	if len(password) < 8 {
		return false
	}

	// Flags to track the presence of required characters
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	// Iterate through the password characters
	for _, char := range password {
		// Check for uppercase
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		// Check for lowercase
		if unicode.IsLower(char) {
			hasLower = true
		}
		// Check for digit
		if unicode.IsDigit(char) {
			hasDigit = true
		}
		// Check for special character
		if unicode.IsPunct(char) || unicode.IsSymbol(char) {
			hasSpecial = true
		}
	}

	// All conditions must be met
	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return false
	}

	return true
}
