package validation

import (
	"regexp"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func RegisterLuhnValidator(validate *validator.Validate) {
	validate.RegisterValidation("luhn", numberLuhnValidator)
}

func numberLuhnValidator(fl validator.FieldLevel) bool {
	number := fl.Field().String()
	return validateNumberLuhn(number)
}

func validateNumberLuhn(number string) bool {
	re := regexp.MustCompile("[^0-9]")
	number = re.ReplaceAllString(number, "")

	if len(number) < 2 {
		return false
	}

	var sum int
	shouldDouble := false
	for i := len(number) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(number[i]))
		if shouldDouble {
			digit = digit * 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}
