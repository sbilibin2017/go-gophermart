package helpers

import (
	"context"
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

func ExecContext(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string, args map[string]any,
) (sql.Result, error) {
	// Убираем лишние пробелы из строки запроса
	formattedQuery := strings.Join(strings.Fields(query), " ")

	// Получаем транзакцию, если она нужна
	tx := txProvider(ctx)

	// Логируем запрос и передаваемые аргументы
	logger.Logger.Info("Executing named exec",
		zap.String("query", formattedQuery),
		zap.Any("args", args),
	)

	var (
		result sql.Result
		err    error
	)

	// Выполняем запрос с транзакцией, если она есть
	if tx != nil {
		result, err = tx.NamedExecContext(ctx, query, args)
	} else {
		result, err = db.NamedExecContext(ctx, query, args)
	}

	// Если произошла ошибка, логируем её
	if err != nil {
		logger.Logger.Error("Error executing named query",
			zap.String("query", formattedQuery),
			zap.Any("args", args),
			zap.Error(err),
		)
		return nil, err
	}

	rowsAffected, _ := result.RowsAffected()
	logger.Logger.Info("Named query executed successfully",
		zap.String("query", formattedQuery),
		zap.Int64("rows_affected", rowsAffected),
	)

	return result, nil
}
