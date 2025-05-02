package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawListRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserBalanceWithdrawListRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserBalanceWithdrawListRepository {
	return &UserBalanceWithdrawListRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserBalanceWithdrawListRepository) ListOrdered(
	ctx context.Context, login string,
) (*[]types.UserBalanceWithdrawDB, error) {
	var withdrawals []types.UserBalanceWithdrawDB
	err := getContext(
		ctx,
		r.db,
		r.txProvider,
		userBalanceWithdrawListQuery,
		&withdrawals,
		login,
	)
	if err != nil {
		return nil, err
	}
	return &withdrawals, nil
}

const userBalanceWithdrawListQuery = `
	SELECT login, number, sum, processed_at
	FROM user_balance_withdrawals
	WHERE login = $1
	ORDER BY processed_at ASC
`
