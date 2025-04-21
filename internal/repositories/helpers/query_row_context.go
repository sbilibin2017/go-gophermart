package helpers

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

func QueryRowContext(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string, args map[string]any,
) (*sqlx.Row, error) {
	formattedQuery := strings.Join(strings.Fields(query), " ")
	tx := txProvider(ctx)
	logger.Logger.Info("Executing named query",
		zap.String("query", formattedQuery),
		zap.Any("args", args),
	)

	var row *sqlx.Row
	if tx != nil {
		namedStmt, err := tx.PrepareNamedContext(ctx, query)
		if err != nil {
			logger.Logger.Error("Failed to prepare named statement (tx)",
				zap.Error(err),
			)
			return nil, err
		}
		row = namedStmt.QueryRowxContext(ctx, args)
	} else {
		namedStmt, err := db.PrepareNamedContext(ctx, query)
		if err != nil {
			logger.Logger.Error("Failed to prepare named statement (db)",
				zap.Error(err),
			)
			return nil, err
		}
		row = namedStmt.QueryRowxContext(ctx, args)
	}

	if err := row.Err(); err != nil {
		logger.Logger.Error("Error executing named query",
			zap.String("query", formattedQuery),
			zap.Any("args", args),
			zap.Error(err),
		)
		return nil, err
	}
	logger.Logger.Info("Named query executed successfully",
		zap.String("query", formattedQuery),
	)
	return row, nil
}
