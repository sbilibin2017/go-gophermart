package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderListOrderedRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserOrderListOrderedRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserOrderListOrderedRepository {
	return &UserOrderListOrderedRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserOrderListOrderedRepository) List(
	ctx context.Context,
) ([]*types.UserOrderDB, error) {
	var userOrders []*types.UserOrderDB
	err := selectContext(
		ctx,
		r.db,
		r.txProvider,
		userOrderListOrderedQuery,
		&userOrders,
	)
	if err != nil {
		return nil, err
	}
	return userOrders, nil
}

const userOrderListOrderedQuery = `
	SELECT * 
	FROM user_orders
	ORDER BY uploaded_at DESC
`
