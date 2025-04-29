package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserBalanceRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewUserBalanceRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *UserBalanceRepository {
	return &UserBalanceRepository{db: db, txProvider: txProvider}
}

func (repo *UserBalanceRepository) GetByLogin(
	ctx context.Context, login string,
) (*domain.UserBalance, error) {
	params := map[string]any{"login": login}
	var userBalance domain.UserBalance

	// Передаем txProvider в queryStruct для обработки транзакции
	err := queryStruct(ctx, repo.db, repo.txProvider, getUserBalanceByLoginQuery, params, &userBalance)
	if err != nil {
		return nil, err
	}
	return &userBalance, nil
}

const getUserBalanceByLoginQuery = `
SELECT current, withdrawn
FROM user_balances
WHERE login = :login
LIMIT 1
`
