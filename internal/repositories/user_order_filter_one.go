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
	var order types.UserOrderDB
	query, args := buildUserOrderFilterOneQuery(number, login)
	err := getContext(ctx, r.db, r.txProvider, query, &order, args...)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func buildUserOrderFilterOneQuery(
	number string, login *string,
) (string, []any) {
	query := `
		SELECT number, login, status, accrual, uploaded_at
		FROM user_orders
		WHERE number = $1
	`
	args := []any{number}
	if login != nil {
		query += " AND login = $2"
		args = append(args, *login)
	}
	return query, args
}
