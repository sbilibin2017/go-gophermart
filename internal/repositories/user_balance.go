package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserBalanceRepository struct {
	db *sqlx.DB
}

func NewUserBalanceRepository(db *sqlx.DB) *UserBalanceRepository {
	return &UserBalanceRepository{db: db}
}

func (repo *UserBalanceRepository) GetByLogin(
	ctx context.Context, login string,
) (*domain.UserBalance, error) {
	params := map[string]any{"login": login}
	var userBalance domain.UserBalance
	err := queryStruct(ctx, repo.db, getUserBalanceByLoginQuery, params, &userBalance)
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
