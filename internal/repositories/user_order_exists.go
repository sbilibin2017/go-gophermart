package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserOrderUploadRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewUserOrderUploadRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *UserOrderUploadRepository {
	return &UserOrderUploadRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserOrderUploadRepository) Exists(ctx context.Context, filter *UserOrderExistsFilter) (bool, error) {
	var exists bool
	err := query(ctx, r.db, r.txProvider, &exists, userOrderExistsQuery, filter)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type UserOrderExistsFilter struct {
	Number string `db:"number"`
	Login  string `db:"login"`
}

const userOrderExistsQuery = `
	SELECT EXISTS (
		SELECT 1 
		FROM user_orders 
		WHERE order_number = :order_number AND (:user_login = '' OR user_login = :user_login)
	)
`
