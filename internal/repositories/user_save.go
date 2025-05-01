package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserSaveRepository {
	return &UserSaveRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserSaveRepository) Save(
	ctx context.Context, login string, password string,
) error {
	return execContext(
		ctx,
		r.db,
		r.txProvider,
		userSaveQuery,
		login, password,
	)
}

const userSaveQuery = `
	INSERT INTO users (login, password)
	VALUES ($1, $2)
	ON CONFLICT (login)
	DO UPDATE SET
		password = EXCLUDED.password
`
