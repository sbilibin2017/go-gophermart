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

func (d *RequestDecoder) Decode(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

var ErrRequestDecoderUnprocessableJson = errors.New("unprocessable json")
