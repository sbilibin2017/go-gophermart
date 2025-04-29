package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserGetByLoginRepository struct {
	db *sqlx.DB
}

func NewUserGetByLoginRepository(db *sqlx.DB) *UserGetByLoginRepository {
	return &UserGetByLoginRepository{db: db}
}

func (repo *UserGetByLoginRepository) GetByLogin(
	ctx context.Context, login string,
) (map[string]any, error) {
	args := map[string]any{
		"login": login,
	}
	result := make(map[string]any)
	err := queryRow(ctx, repo.db, getUserByLoginQuery, &result, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

const getUserByLoginQuery = `
SELECT login, password
FROM users
WHERE login = :login
LIMIT 1
`
