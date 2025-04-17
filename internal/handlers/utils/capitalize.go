package utils

import (
	"strings"
	"unicode/utf8"
)

func Capitalize(msg string) string {
	if len(msg) == 0 {
		return msg
	}
	firstRune, size := utf8.DecodeRuneInString(msg)
	if firstRune == utf8.RuneError {
		return msg
	}
	return strings.ToUpper(string(firstRune)) + msg[size:]
}
