package utils

import (
	"net/http"
	"unicode"
)

func RespondWithError(w http.ResponseWriter, err error, status int) {
	http.Error(w, capitalizeFirstLetter(err.Error()), status)
}

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	return string(unicode.ToUpper(rune(s[0]))) + s[1:]
}
