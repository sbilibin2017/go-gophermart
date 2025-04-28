package engines

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type ExecutorEngine struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewExecutorEngine(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) (*ExecutorEngine, error) {
	return &ExecutorEngine{db: db, txProvider: txProvider}, nil
}

func (e *ExecutorEngine) Exec(
	ctx context.Context,
	query string,
	args any,
) error {
	logger.Logger.Infof("query: %s", query)
	logger.Logger.Infof("args: %+v", args)
	return handleExecArgs(ctx, e.db, e.txProvider, query, args)
}

func handleExecArgs(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	args any,
) error {
	switch v := args.(type) {
	case []any:
		return handleExecVarArgs(ctx, db, txProvider, query, v)
	case map[string]any:
		return handleExecMapArgs(ctx, db, txProvider, query, v)
	case struct{}:
		return handleExecStructArgs(ctx, db, txProvider, query, v)
	default:
		return nil
	}
}

func handleExecVarArgs(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	args []any,
) error {
	tx := txProvider(ctx)
	var err error
	if tx != nil {
		_, err = db.ExecContext(ctx, query, args...)
	} else {
		_, err = db.ExecContext(ctx, query, args...)
	}
	return err
}

func handleExecMapArgs(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	args map[string]any,
) error {
	tx := txProvider(ctx)
	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, args)
	} else {
		_, err = db.NamedExecContext(ctx, query, args)
	}
	return err
}

func handleExecStructArgs(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	args struct{},
) error {
	tx := txProvider(ctx)
	var err error
	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, args)
	} else {
		_, err = db.NamedExecContext(ctx, query, args)
	}
	return err
}
