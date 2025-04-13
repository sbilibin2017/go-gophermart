package json

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RequestDecoder struct{}

func NewRequestDecoder() *RequestDecoder {
	return &RequestDecoder{}
}

func (d *RequestDecoder) Decode(w http.ResponseWriter, r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return ErrRequestDecoderUnprocessableJSON
	}
	return nil
}

var ErrRequestDecoderUnprocessableJSON = errors.New("invalid json")
