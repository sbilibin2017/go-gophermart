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
	ctx context.Context, number string, login string, status string,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		userOrderSaveQuery,
		number, login, status,
	)
}

const userOrderSaveQuery = `
	INSERT INTO user_orders (number, login, status)
	VALUES ($1, $2, $3)
	ON CONFLICT (number) DO NOTHING
`
