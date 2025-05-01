package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceFilterOneRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserBalanceFilterOneRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserBalanceFilterOneRepository {
	return &UserBalanceFilterOneRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserBalanceFilterOneRepository) FilterOne(
	ctx context.Context, login string,
) (*types.UserBalanceDB, error) {
	var userBalance types.UserBalanceDB
	err := getContext(
		ctx,
		r.db,
		r.txProvider,
		userBalanceFilterOneQuery,
		&userBalance,
		login,
	)
	if err != nil {
		return nil, err
	}
	return &userBalance, nil
}

const userBalanceFilterOneQuery = `
	SELECT * 
	FROM user_balances
	WHERE login = $1
`
