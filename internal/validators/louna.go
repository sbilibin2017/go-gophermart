package validators

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateNumberWithLouna(value string) error {
	value = strings.ReplaceAll(value, " ", "")
	if len(value) == 0 {
		return ErrInvalidLounaNumber
	}
	for i := 0; i < len(value); i++ {
		if value[i] < '0' || value[i] > '9' {
			return ErrInvalidLounaNumber
		}
	}
	sum := 0
	shouldDouble := false
	for i := len(value) - 1; i >= 0; i-- {
		num, _ := strconv.Atoi(string(value[i]))
		if shouldDouble {
			num *= 2
			if num > 9 {
				num -= 9
			}
		}
		sum += num
		shouldDouble = !shouldDouble
	}
	if sum%10 == 0 {
		return nil
	} else {
		return ErrInvalidLounaNumber
	}
}

func RegisterLounaValidator(v *validator.Validate) {
	v.RegisterValidation("louna", func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		err := ValidateNumberWithLouna(value)
		return err == nil
	})
}

var (
	ErrInvalidLounaNumber = errors.New("invalid louna number")
)
