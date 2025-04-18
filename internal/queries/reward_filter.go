package queries

import (
	"strconv"
	"strings"
)

const RewardExistsQuery = `SELECT 1 FROM rewards WHERE %s LIMIT 1`

func BuildRewardExistsFilters(
	filter map[string]any,
) (string, []any) {
	var builder strings.Builder
	var args []any
	paramIndex := 1

	if len(filter) == 0 {
		return "", args
	}

	for key, val := range filter {
		if paramIndex > 1 {
			builder.WriteString(" AND ")
		}

		builder.WriteByte('"')
		builder.WriteString(key)
		builder.WriteByte('"')
		builder.WriteString(" = $")
		builder.WriteString(strconv.Itoa(paramIndex))

		args = append(args, val)
		paramIndex++
	}

	return builder.String(), args
}
