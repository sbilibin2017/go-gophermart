package repositories

import "strings"

func buildColumnsString(fields []string) string {
	if len(fields) == 0 {
		return "*"
	}
	return strings.Join(fields, ", ")
}
