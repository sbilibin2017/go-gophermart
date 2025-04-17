package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"unicode"
)

var (
	ErrInvalidJSON    = errors.New("неверный формат json")
	ErrInternalServer = errors.New("внутренняя ошибка сервера")
)

func capitalizeError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if len(msg) == 0 {
		return err
	}
	runes := []rune(msg)
	runes[0] = unicode.ToUpper(runes[0])
	return errors.New(string(runes))
}

func decodeJSON[T any](r *http.Request, v *T) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(v)
	if err != nil {
		return ErrInvalidJSON
	}
	return nil
}

func writeTextResponse(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg))
}
