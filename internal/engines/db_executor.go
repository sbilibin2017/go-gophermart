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

func (d *DBExecutor) Execute(
	ctx context.Context,
	query string,
	argMap map[string]any,
) error {
	tx, ok := d.txProvider(ctx)
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
	_, err := d.db.NamedExecContext(ctx, query, argMap)
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
