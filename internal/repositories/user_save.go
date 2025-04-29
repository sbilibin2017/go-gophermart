package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserSaveRepository struct {
	db *sqlx.DB
}

func NewUserSaveRepository(db *sqlx.DB) *UserSaveRepository {
	return &UserSaveRepository{db: db}
}

func (repo *UserSaveRepository) Save(
	ctx context.Context, login string, password string,
) error {
	args := map[string]any{
		"login":    login,
		"password": password,
	}
	return exec(ctx, repo.db, userSaveQuery, args)
}

const userSaveQuery = `
INSERT INTO users (login, password)
	VALUES (:login, :password)
ON CONFLICT (login) DO UPDATE
	SET password = EXCLUDED.password
`
