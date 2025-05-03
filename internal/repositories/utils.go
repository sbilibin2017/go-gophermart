package repositories

import "github.com/sbilibin2017/go-gophermart/internal/logger"

func logQuery(query string, args interface{}, err error) {
	logger.Info("Executed query", "query", query, "args", args)
	if err != nil {
		logger.Error("Error executing query", "error", err)
	}
}
