package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserOrderSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserOrderSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserOrderSaveRepository {
	return &UserOrderSaveRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserOrderSaveRepository) Save(
	ctx context.Context, number string, login string, status string, accrual *int64,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		userOrderSaveQuery,
		number, login, status, accrual,
	)
}

const userOrderSaveQuery = `
	INSERT INTO user_orders (number, login, status, accrual)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (number) DO UPDATE
	SET status = EXCLUDED.status,
	    accrual = COALESCE(EXCLUDED.accrual, user_orders.accrual)
`
