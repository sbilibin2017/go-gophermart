package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
)

func ValidateStruct(
	w http.ResponseWriter,
	validate *validator.Validate,
	v interface{},
	errMap map[string]error,
) error {
	err := validate.Struct(v)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				field := fieldErr.Field()
				if err, found := errMap[field]; found {
					return err
				}
			}
		}
		return errors.ErrDataIsNotValid
	}
	return nil
}
