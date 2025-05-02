package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderListRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserOrderListRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserOrderListRepository {
	return &UserOrderListRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserOrderListRepository) ListOrdered(
	ctx context.Context, login *string,
) ([]*types.UserOrderDB, error) {
	var orders []*types.UserOrderDB
	query, args := buildUserOrdersListQuery(login)
	err := getContext(
		ctx,
		r.db,
		r.txProvider,
		query,
		&orders,
		args...,
	)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func buildUserOrdersListQuery(login *string) (string, []any) {
	if login != nil {
		return baseUserOrdersListQuery + `
		WHERE login = $1
		ORDER BY uploaded_at DESC
		`, []any{*login}
	}
	return baseUserOrdersListQuery + `
	ORDER BY uploaded_at DESC
	`, nil
}

const baseUserOrdersListQuery = `
	SELECT number, login, status, accrual, uploaded_at
	FROM user_orders
`
