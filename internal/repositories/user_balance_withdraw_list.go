package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawListOrderedRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserBalanceWithdrawListOrderedRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserBalanceWithdrawListOrderedRepository {
	return &UserBalanceWithdrawListOrderedRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserBalanceWithdrawListOrderedRepository) List(
	ctx context.Context,
) ([]*types.UserBalanceWithdrawDB, error) {
	var userWithdraws []*types.UserBalanceWithdrawDB
	err := selectContext(
		ctx,
		r.db,
		r.txProvider,
		userBalanceWithdrawListOrderedQuery,
		&userWithdraws,
	)
	if err != nil {
		return nil, err
	}
	return userWithdraws, nil
}

const userBalanceWithdrawListOrderedQuery = `
	SELECT * 
	FROM user_withdraws
	ORDER BY processed_at DESC
`
