package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderFilterOneRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserOrderFilterOneRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserOrderFilterOneRepository {
	return &UserOrderFilterOneRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserOrderFilterOneRepository) FilterOne(
	ctx context.Context, number string, login *string,
) (*types.UserOrderDB, error) {
	var userOrder types.UserOrderDB
	query, args := buildUserOrderFilterOneQuery(number, login)
	err := getContext(
		ctx,
		r.db,
		r.txProvider,
		query,
		&userOrder,
		args...,
	)
	if err != nil {
		return nil, err
	}
	return &userOrder, nil
}

func buildUserOrderFilterOneQuery(number string, login *string) (string, []any) {
	var query string
	var args []any
	query = ` WHERE number = $1`
	args = append(args, number)
	if login != nil {
		query += ` AND login = $2`
		args = append(args, *login)
	}
	return userOrderFilterOneBaseQuery + query, args
}

const userOrderFilterOneBaseQuery = `SELECT * FROM user_orders`
