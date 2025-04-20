package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

func execContext(
	ctx context.Context,
	db *sql.DB,
	txProvider func(ctx context.Context) *sql.Tx,
	query string, args ...any,
) (sql.Result, error) {
	// Убираем лишние пробелы внутри запроса
	formattedQuery := strings.Join(strings.Fields(query), " ")

	tx := txProvider(ctx)
	logger.Logger.Info("Executing query",
		zap.String("query", formattedQuery),
		zap.Any("args", args),
	)

	var result sql.Result
	var err error
	if tx != nil {
		result, err = tx.ExecContext(ctx, query, args...)
	} else {
		result, err = db.ExecContext(ctx, query, args...)
	}

	if err != nil {
		logger.Logger.Error("Error executing query",
			zap.String("query", formattedQuery),
			zap.Any("args", args),
			zap.Error(err),
		)
	} else {
		rowsAffected, _ := result.RowsAffected()
		logger.Logger.Info("Query executed successfully",
			zap.String("query", formattedQuery),
			zap.Int64("rows_affected", rowsAffected),
		)
	}

	return result, err
}

func queryRowContext(
	ctx context.Context,
	db *sql.DB,
	txProvider func(ctx context.Context) *sql.Tx,
	query string, args ...any,
) (*sql.Row, error) {
	// Убираем лишние пробелы внутри запроса
	formattedQuery := strings.Join(strings.Fields(query), " ")

	tx := txProvider(ctx)
	logger.Logger.Info("Executing query",
		zap.String("query", formattedQuery),
		zap.Any("args", args),
	)

	var row *sql.Row
	if tx != nil {
		row = tx.QueryRowContext(ctx, query, args...)
	} else {
		row = db.QueryRowContext(ctx, query, args...)
	}

	if err := row.Err(); err != nil {
		logger.Logger.Error("Error executing query",
			zap.String("query", formattedQuery),
			zap.Any("args", args),
			zap.Error(err),
		)
		return nil, err
	}

	logger.Logger.Info("Query executed successfully",
		zap.String("query", formattedQuery),
	)

	return row, nil
}

func scanRow(row *sql.Row, dest ...interface{}) error {
	return row.Scan(dest...)
}
