package helpers

import (
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

func MapScan(row *sqlx.Row) map[string]any {
	result := make(map[string]any)
	if err := row.MapScan(result); err != nil {
		logger.Logger.Info("No rows found ")
		return nil
	}
	return result
}

func Scan[T any](row *sqlx.Row) (T, error) {
	var result T
	err := row.Scan(&result)
	if err != nil {
		logger.Logger.Error("Error scanning result", zap.Error(err))
		var zero T
		return zero, err
	}
	return result, nil
}
