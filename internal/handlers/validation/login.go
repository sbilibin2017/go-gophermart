package validation

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateLogin(fl validator.FieldLevel) bool {
	login := fl.Field().String()
	match := regexp.MustCompile(`^[a-zA-Z0-9]{3,}$`)
	return match.MatchString(login)
}
