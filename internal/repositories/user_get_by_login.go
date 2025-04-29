package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserGetByLoginRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewUserGetByLoginRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *UserGetByLoginRepository {
	return &UserGetByLoginRepository{db: db, txProvider: txProvider}
}

func (repo *UserGetByLoginRepository) GetByLogin(
	ctx context.Context, login string,
) (*domain.User, error) {
	params := map[string]any{"login": login}
	var user domain.User
	err := queryStruct(ctx, repo.db, repo.txProvider, getUserByLoginQuery, params, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

const getUserByLoginQuery = `
SELECT login, password
FROM users
WHERE login = :login
LIMIT 1
`
