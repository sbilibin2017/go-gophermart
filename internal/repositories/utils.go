package repositories

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

func execContext(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
	query string,
	args ...any,
) error {
	e := getExecutor(ctx, db, txProvider)
	_, err := e.ExecContext(ctx, query, args...)
	err = handleExecutorError(err)
	logQuery(query, args, err)
	return err
}

func getContext(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
	query string,
	dest any,
	args ...any,
) error {
	e := getExecutor(ctx, db, txProvider)
	err := sqlx.GetContext(ctx, e, dest, query, args...)
	err = handleExecutorError(err)
	logQuery(query, args, err)
	return err
}

func getExecutor(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) sqlx.ExtContext {
	tx, err := txProvider(ctx)
	if err != nil {
		return db
	}
	return tx
}

func handleExecutorError(err error) error {
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func logQuery(query string, args any, err error) {
	logger.Logger.Infof("Executed query: %s", query)
	logger.Logger.Infof("Query arguments: %v", args)
	if err != nil {
		logger.Logger.Errorf("Error executing query: %v", err)
	}
}
