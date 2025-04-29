package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserGetByLoginRepository struct {
	db *sqlx.DB
}

func NewUserGetByLoginRepository(db *sqlx.DB) *UserGetByLoginRepository {
	return &UserGetByLoginRepository{db: db}
}

func (repo *UserGetByLoginRepository) GetByLogin(
	ctx context.Context, login string,
) (*domain.User, error) {
	params := map[string]any{"login": login}
	var user domain.User
	err := queryStruct(ctx, repo.db, getUserByLoginQuery, params, &user)
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
