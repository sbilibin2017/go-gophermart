package engines

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type QuerierEngine struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewQuerierEngine(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) (*QuerierEngine, error) {
	return &QuerierEngine{db: db, txProvider: txProvider}, nil
}

func (e *QuerierEngine) Query(
	ctx context.Context, query string, dest any, args any,
) error {
	logger.Logger.Infof("query: %s", query)
	logger.Logger.Infof("args: %+v", args)
	return handleQueryArgs(ctx, e.db, e.txProvider, query, dest, args)
}

func handleQueryArgs(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	dest any,
	args any,
) error {
	switch v := args.(type) {
	case []any:
		return handleQueryVarArgs(ctx, db, txProvider, query, dest, v)
	case map[string]any:
		return handleQueryMapArgs(ctx, db, txProvider, query, dest, v)
	case struct{}:
		return handleQueryStructArgs(ctx, db, txProvider, query, dest, v)
	default:
		return nil
	}
}

func handleQueryVarArgs(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	dest any,
	args []any,
) error {
	tx := txProvider(ctx)
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.QueryxContext(ctx, query, args...)
	} else {
		rows, err = db.QueryxContext(ctx, query, args...)
	}
	if err != nil {
		logger.Logger.Error("Error executing query: ", err)
		return err
	}
	defer rows.Close()
	return scanRows(rows, dest)
}

func handleQueryMapArgs(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	dest any,
	args map[string]any,
) error {
	tx := txProvider(ctx)
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.NamedQuery(query, args)
	} else {
		rows, err = db.NamedQueryContext(ctx, query, args)
	}
	if err != nil {
		logger.Logger.Error("Error executing query: ", err)
		return err
	}
	defer rows.Close()
	return scanRows(rows, dest)
}

func handleQueryStructArgs(
	ctx context.Context,
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
	query string,
	dest any,
	args struct{},
) error {
	tx := txProvider(ctx)
	var rows *sqlx.Rows
	var err error
	if tx != nil {
		rows, err = tx.QueryxContext(ctx, query, args)
	} else {
		rows, err = db.QueryxContext(ctx, query, args)
	}
	if err != nil {
		logger.Logger.Error("Error executing query: ", err)
		return err
	}
	defer rows.Close()
	return scanRows(rows, dest)
}

func scanRows(rows *sqlx.Rows, dest any) error {
	if err := rows.StructScan(dest); err != nil {
		logger.Logger.Error("Error scanning query results: ", err)
		return err
	}
	return nil
}
