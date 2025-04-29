package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewUserSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *UserSaveRepository {
	return &UserSaveRepository{db: db, txProvider: txProvider}
}

func (repo *UserSaveRepository) Save(
	ctx context.Context, user *domain.User,
) error {
	return exec(ctx, repo.db, repo.txProvider, userSaveQuery, user)
}

const userSaveQuery = `
INSERT INTO users (login, password)
	VALUES (:login, :password)
ON CONFLICT (login) DO UPDATE
	SET password = EXCLUDED.password
`
