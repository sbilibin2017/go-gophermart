package repositories

import (
	"context"

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
	ctx context.Context, number string, sum int64, processedAt string,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		userBalanceWithdrawSaveQuery,
		number, sum, processedAt,
	)
}

const userBalanceWithdrawSaveQuery = `
	INSERT INTO user_withdraws (number, sum, processed_at)
	VALUES ($1, $2, $3)
	ON CONFLICT (number)
	DO NOTHING
`
