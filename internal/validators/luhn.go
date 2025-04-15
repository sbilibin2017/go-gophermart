package validators

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateLuhn(fl validator.FieldLevel) bool {
	number := fl.Field().String()
	return validateNumberWithLuhn(number)
}

func validateNumberWithLuhn(number string) bool {
	number = strings.ReplaceAll(number, " ", "")
	for _, c := range number {
		if c < '0' || c > '9' {
			return false
		}
	}
	sum := 0
	shouldDouble := false
	for i := len(number) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(string(number[i]))
		if shouldDouble {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		shouldDouble = !shouldDouble
	}
	return sum%10 == 0
}
