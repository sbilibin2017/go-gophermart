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
	ctx context.Context, number string, login string,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		userOrderSaveQuery,
		number, login,
	)
}

const userOrderSaveQuery = `
	INSERT INTO user_orders (number, login)
	VALUES ($1, $2)
	ON CONFLICT (number)
	DO UPDATE SET
		login = EXCLUDED.login,		
`
