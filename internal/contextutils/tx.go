package contextutils

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type contextTxKey struct{}

func SetTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, contextTxKey{}, tx)
}

func GetTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(contextTxKey{}).(*sqlx.Tx)
	if !ok {
		return nil, errors.New("transaction is not in context")
	}
	return tx, nil
}

func GetDBExecutor(
	ctx context.Context,
	fallback *sqlx.DB,
) sqlx.ExtContext {
	tx, err := GetTx(ctx)
	if err != nil {
		return fallback
	}
	return tx
}
