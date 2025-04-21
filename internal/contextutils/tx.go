package contextutils

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type txKeyType string

const txKey txKeyType = "tx"

func TxFromContext(ctx context.Context) *sqlx.Tx {
	if tx, ok := ctx.Value(txKey).(*sqlx.Tx); ok {
		return tx
	}
	return nil
}

func TxToContext(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}
