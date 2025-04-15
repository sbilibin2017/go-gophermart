package options

import (
	"flag"
	"os"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/log"
)

// Combine возвращает значение из переменной окружения, затем флага, затем значение по умолчанию.
// Приоритет: ENV > FLAG > DEFAULT
func Combine(envKey string, flagKey string, defaultValue string) string {
	if val, ok := os.LookupEnv(envKey); ok && strings.TrimSpace(val) != "" {
		log.Info("Используется значение из переменной окружения",
			"envKey", envKey,
			"value", val,
		)
		return val
	}

	flagVal := flag.Lookup(flagKey)
	if flagVal != nil && strings.TrimSpace(flagVal.Value.String()) != "" {
		log.Info("Используется значение из флага",
			"flagKey", flagKey,
			"value", flagVal.Value.String(),
		)
		return flagVal.Value.String()
	}

	log.Info("Используется значение по умолчанию",
		"envKey", envKey,
		"flagKey", flagKey,
		"defaultValue", defaultValue,
	)
	return defaultValue
}
