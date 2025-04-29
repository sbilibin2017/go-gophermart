package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserSaveRepository struct {
	db *sqlx.DB
}

func NewUserSaveRepository(db *sqlx.DB) *UserSaveRepository {
	return &UserSaveRepository{db: db}
}

func (repo *UserSaveRepository) Save(
	ctx context.Context, user *domain.User,
) error {
	return exec(ctx, repo.db, userSaveQuery, user)
}

const userSaveQuery = `
INSERT INTO users (login, password)
	VALUES (:login, :password)
ON CONFLICT (login) DO UPDATE
	SET password = EXCLUDED.password
`
