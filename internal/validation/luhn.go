package validation

import (
	"github.com/go-playground/validator/v10"
)

func IsLuhnValidationError(err validator.FieldError) bool {
	return err.Tag() == "luhn"
}

func RegisterLuhnValidation(v *validator.Validate) {
	v.RegisterValidation("luhn", luhnValidation)
}

func luhnValidation(fl validator.FieldLevel) bool {
	orderNumber, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return validateNumber(orderNumber)
}

func validateNumber(orderNumber string) bool {
	if len(orderNumber) == 0 {
		return false
	}
	var sum int
	for i, digit := range orderNumber {
		n := int(digit - '0')
		if (i+1)%2 != 0 {
			n = n * 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
	}
	return sum%10 == 0
}
