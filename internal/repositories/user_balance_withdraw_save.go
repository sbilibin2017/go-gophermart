package repositories

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserBalanceWithdrawSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserBalanceWithdrawSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserBalanceWithdrawSaveRepository {
	return &UserBalanceWithdrawSaveRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserBalanceWithdrawSaveRepository) Save(
	ctx context.Context, login string, number string, sum int64,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		userBalanceWithdrawSaveQuery,
		login, number, sum, time.Now().UTC(),
	)
}

const userBalanceWithdrawSaveQuery = `
	INSERT INTO user_balance_withdrawals (login, number, sum, processed_at)
	VALUES ($1, $2, $3, $4)
`
