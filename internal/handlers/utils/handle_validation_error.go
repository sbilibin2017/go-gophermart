package utils

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

func HandleValidationError(w http.ResponseWriter, err error) {
	validationErrors, ok := err.(validator.ValidationErrors)
	if ok {
		for _, validationErr := range validationErrors {
			http.Error(w, validationErr.Error(), http.StatusBadRequest)
			return
		}
	}
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
