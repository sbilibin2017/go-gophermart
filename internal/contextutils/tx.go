package contextutils

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type txContextKey struct{}

func GetTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(txContextKey{}).(*sqlx.Tx)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}
	return tx, nil
}

func SetTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txContextKey{}, tx)
}
