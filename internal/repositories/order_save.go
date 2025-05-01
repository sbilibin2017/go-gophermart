package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
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
	ctx context.Context, order *types.OrderDB,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		orderSaveQuery,
		order.Number, order.Status, order.Accrual,
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
