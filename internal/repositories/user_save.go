package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewUserSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *UserSaveRepository {
	return &UserSaveRepository{db: db, txProvider: txProvider}
}

type UserSave struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}

func (r *UserSaveRepository) Save(ctx context.Context, user *UserSave) error {
	return command(ctx, r.db, r.txProvider, userSaveQuery, user)
}

const userSaveQuery = `
	INSERT INTO users ("login", "password", created_at, updated_at)
	VALUES (:login, :password, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
`
