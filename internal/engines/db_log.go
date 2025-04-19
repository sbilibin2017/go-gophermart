package engines

import (
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	logger, _ = zapConfig.Build()
}

func prepareQuery(query string) string {
	logQuery := strings.ReplaceAll(query, "\n", " ")
	logQuery = strings.ReplaceAll(logQuery, "\t", " ")
	logQuery = regexp.MustCompile(`\s+`).ReplaceAllString(logQuery, " ")
	logQuery = strings.TrimSpace(logQuery)
	return logQuery
}

func logQueryExecution(tx *sqlx.Tx, query string, args map[string]any) {
	msg := "Executing query without transaction"
	if tx != nil {
		msg = "Executing query within transaction"
	}

	logger.Info(msg,
		zap.String("query", query),
		zap.Any("args", args),
	)
}

func logQueryError(query string, args map[string]any, err error) {
	logger.Error("Failed to execute query",
		zap.String("query", query),
		zap.Any("args", args),
		zap.Error(err),
	)
}
