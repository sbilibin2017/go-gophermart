package utils

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrUnprocessableJSON = errors.New("unprocessable json")

type Decoder struct{}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return ErrUnprocessableJSON
	}
	return nil
}
