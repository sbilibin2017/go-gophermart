package utils

import (
	"encoding/json"
	"net/http"
)

func Decode[T any](w http.ResponseWriter, r *http.Request, v *T) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}
