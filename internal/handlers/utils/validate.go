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
) error {
	err := validate.Struct(v)
	if err != nil {
		http.Error(w, errors.ErrDataIsNotValid.Error(), http.StatusBadRequest)
		return err
	}
	return nil
}
