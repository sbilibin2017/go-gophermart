package json

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ResponseEncoder struct{}

func NewResponseEncoder() *ResponseEncoder {
	return &ResponseEncoder{}
}

func (e *ResponseEncoder) Encode(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if v == nil {
		return nil
	}

	return json.NewEncoder(w).Encode(v)
}

var ErrResponseEncoderUnprocessableJson = errors.New("unprocessable json")
