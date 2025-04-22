package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

func execContextNamed(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
	query string,
	argMap map[string]any,
) error {
	tx, ok := txProvider(ctx)
	if ok {
		logger.Logger.Info("Executing named query inside transaction")
		_, err := tx.NamedExecContext(ctx, query, argMap)
		if err != nil {
			logger.Logger.Error("Error executing named query in transaction",
				zap.String("query", query),
				zap.Any("args", argMap),
				zap.Error(err),
			)
			return err
		}
		return nil
	}

	logger.Logger.Info("Executing named query outside transaction")
	_, err := db.NamedExecContext(ctx, query, argMap)
	if err != nil {
		logger.Logger.Error("Error executing named query outside transaction",
			zap.String("query", query),
			zap.Any("args", argMap),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func getContextNamed(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
	dest any,
	query string,
	argMap map[string]any,
) error {
	tx, ok := txProvider(ctx)
	if ok {
		logger.Logger.Info("Executing named query inside transaction")
		rows, err := tx.NamedQuery(query, argMap)
		if err != nil {
			logger.Logger.Error("Error executing named query in transaction",
				zap.String("query", query),
				zap.Any("args", argMap),
				zap.Error(err),
			)
			return err
		}
		defer rows.Close()

		if rows.Next() {
			if err := rows.StructScan(dest); err != nil {
				logger.Logger.Error("Error scanning result from named query",
					zap.String("query", query),
					zap.Any("args", argMap),
					zap.Error(err),
				)
				return err
			}
		}
		return nil
	}

	logger.Logger.Info("Executing named query outside transaction")
	rows, err := db.NamedQueryContext(ctx, query, argMap)
	if err != nil {
		logger.Logger.Error("Error executing named query outside transaction",
			zap.String("query", query),
			zap.Any("args", argMap),
			zap.Error(err),
		)
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.StructScan(dest); err != nil {
			logger.Logger.Error("Error scanning result from named query",
				zap.String("query", query),
				zap.Any("args", argMap),
				zap.Error(err),
			)
			return err
		}
	}

	return nil
}
