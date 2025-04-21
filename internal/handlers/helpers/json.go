package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrInvalidRequestBody = errors.New("invalid request body")

func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(dst)
	if err != nil {
		return ErrInvalidRequestBody
	}
	return nil
}
