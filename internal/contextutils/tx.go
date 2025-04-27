package contextutils

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type contextKey string

const txContextKey contextKey = "tx"

func GetTx(ctx context.Context) *sqlx.Tx {
	tx, ok := ctx.Value(txContextKey).(*sqlx.Tx)
	if !ok {
		return nil
	}
	return tx
}

func SetTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}
