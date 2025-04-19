package engines

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DBExecutor struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewDBExecutor(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *DBExecutor {
	return &DBExecutor{db: db, txProvider: txProvider}
}

func (e *DBExecutor) Execute(ctx context.Context, query string, args map[string]any) error {
	tx := e.txProvider(ctx)
	queryPrepared := prepareQuery(query)
	logQueryExecution(tx, queryPrepared, args)

	var err error

	if tx != nil {
		_, err = tx.NamedExecContext(ctx, query, args)
	} else {
		_, err = e.db.NamedExecContext(ctx, query, args)
	}

	if err != nil {
		logQueryError(queryPrepared, args, err)
		return err
	}

	return nil
}
