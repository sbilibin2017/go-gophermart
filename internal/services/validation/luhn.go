package validation

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func LuhnValidator(fl validator.FieldLevel) bool {
	order := fl.Field().String()
	if order == "" {
		return false
	}
	if !isValidLuhn(order) {
		return false
	}
	return true
}

func isValidLuhn(order string) bool {
	order = strings.ReplaceAll(order, " ", "")
	for _, c := range order {
		if c < '0' || c > '9' {
			return false
		}
	}
	var sum int
	shouldDouble := false
	for i := len(order) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(order[i]))
		if shouldDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		shouldDouble = !shouldDouble
	}
	return sum%10 == 0
}
