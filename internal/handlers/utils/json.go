package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrInvalidJSON = errors.New("invalid json")

func DecodeJSON[T any](r *http.Request, v *T) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return ErrInvalidJSON
	}
	return nil
}

func EncodeJSON[T any](w http.ResponseWriter, v T) error {
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		return ErrInvalidJSON
	}
	return nil
}
