package engines

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"go.uber.org/zap"
)

type DBExecutor struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewDBExecutor(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *DBExecutor {
	return &DBExecutor{
		db:         db,
		txProvider: txProvider,
	}
}

func (d *DBExecutor) Execute(
	ctx context.Context,
	query string,
	args any,
) error {
	tx, ok := d.txProvider(ctx)
	if ok {
		logger.Logger.Info("Executing named query inside transaction")
		_, err := tx.NamedExecContext(ctx, query, args)
		if err != nil {
			logger.Logger.Error("Error executing named query in transaction",
				zap.String("query", query),
				zap.Any("args", args),
				zap.Error(err),
			)
			return err
		}
		return nil
	}

	logger.Logger.Info("Executing named query outside transaction")
	_, err := d.db.NamedExecContext(ctx, query, args)
	if err != nil {
		logger.Logger.Error("Error executing named query outside transaction",
			zap.String("query", query),
			zap.Any("args", args),
			zap.Error(err),
		)
		return err
	}

	return nil
}
