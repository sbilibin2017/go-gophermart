package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserBalanceRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserBalanceRepository {
	return &UserBalanceRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserBalanceRepository) FilterOne(
	ctx context.Context, login string,
) (*types.UserBalanceDB, error) {
	var userBalance types.UserBalanceDB
	err := getContext(
		ctx,
		r.db,
		r.txProvider,
		userBalanceQuery,
		&userBalance,
		login,
	)
	if err != nil {
		return nil, err
	}
	return &userBalance, nil
}

const userBalanceQuery = `
	SELECT current, withdrawn
	FROM user_balances
	WHERE login = $1
`
