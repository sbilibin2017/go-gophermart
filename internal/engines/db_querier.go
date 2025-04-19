package engines

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DBQuerier struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewDBQuerier(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *DBQuerier {
	return &DBQuerier{db: db, txProvider: txProvider}
}

func (q *DBQuerier) Query(
	ctx context.Context, dest any, query string, args map[string]any,
) error {
	tx := q.txProvider(ctx)
	logQuery := prepareQuery(query)
	logQueryExecution(tx, logQuery, args)

	var (
		rows *sqlx.Rows
		err  error
	)

	if tx != nil {
		rows, err = tx.NamedQuery(query, args)
	} else {
		rows, err = q.db.NamedQuery(query, args)
	}

	if err != nil {
		logQueryError(logQuery, args, err)
		return err
	}
	defer rows.Close()

	if rows.Next() {
		if err = rows.Scan(dest); err != nil {
			logQueryError(logQuery, args, err)
			return err
		}
	}

	if err = rows.Err(); err != nil {
		logQueryError(logQuery, args, err)
		return err
	}

	return nil
}
