package json

import (
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/errors"
)

type RequestDecoder struct{}

func NewRequestDecoder() *RequestDecoder {
	return &RequestDecoder{}
}

func (d *RequestDecoder) Decode(w http.ResponseWriter, r *http.Request, v any) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return errors.ErrRequestDecoderUnprocessableJSON
	}
	return nil
}
