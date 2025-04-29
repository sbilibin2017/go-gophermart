package contextutils

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type contextTxKey string

const txContextKey contextTxKey = "tx"

func SetTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey, tx)
}

func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(txContextKey).(*sqlx.Tx)
	return tx, ok
}
