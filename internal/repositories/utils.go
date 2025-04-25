package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

func query(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	dest any,
	query string,
	arg any,
) error {
	if logger.Logger != nil {
		logger.Logger.Info("Executing query", zap.String("query", query), zap.Any("args", arg))
	}
	tx := txProvider(ctx)
	if tx != nil {
		namedStmt, err := tx.PrepareNamedContext(ctx, query)
		if err != nil {
			if logger.Logger != nil {
				logger.Logger.Error("Failed to prepare named statement", zap.Error(err))
			}
			return err
		}
		return namedStmt.GetContext(ctx, dest, arg)
	}
	namedStmt, err := db.PrepareNamedContext(ctx, query)
	if err != nil {
		if logger.Logger != nil {
			logger.Logger.Error("Failed to prepare named statement", zap.Error(err))
		}
		return err
	}
	return namedStmt.GetContext(ctx, dest, arg)
}

func command(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	arg interface{},
) error {
	if logger.Logger != nil {
		logger.Logger.Info("Executing command", zap.String("query", query), zap.Any("args", arg))
	}
	tx := txProvider(ctx)
	if tx != nil {
		_, err := tx.NamedExecContext(ctx, query, arg)
		if err != nil {
			if logger.Logger != nil {
				logger.Logger.Error("Failed to execute named statement", zap.Error(err))
			}
		}
		return err
	}
	_, err := db.NamedExecContext(ctx, query, arg)
	if err != nil {
		if logger.Logger != nil {
			logger.Logger.Error("Failed to execute named statement", zap.Error(err))
		}
	}
	return err
}
