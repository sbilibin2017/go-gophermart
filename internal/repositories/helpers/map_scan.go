package helpers

import (
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func MapScan(row *sqlx.Row) map[string]any {
	result := make(map[string]any)
	if err := row.MapScan(result); err != nil {
		logger.Logger.Info("No rows found ")
		return nil
	}
	return result
}
