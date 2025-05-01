package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewOrderSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *OrderSaveRepository {
	return &OrderSaveRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *OrderSaveRepository) Save(
	ctx context.Context, number string, status string, accrual int64,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		orderSaveQuery,
		number, status, accrual,
	)
}

const orderSaveQuery = `
	INSERT INTO orders (number, status, accrual)
	VALUES ($1, $2, $3)
	ON CONFLICT (number)
	DO UPDATE SET
		status = EXCLUDED.status,
		accrual = EXCLUDED.accrual
`
