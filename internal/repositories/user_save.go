package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type UserSaveRepository struct {
	db *sqlx.DB
}

func NewUserSaveRepository(db *sqlx.DB) *UserSaveRepository {
	return &UserSaveRepository{db: db}
}

const userSaveQuery = `
	INSERT INTO users (login, password)
	VALUES ($1, $2);
`

type UserSave struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}

func (r *UserSaveRepository) Save(ctx context.Context, u *UserSave) error {
	_, err := r.db.ExecContext(ctx, userSaveQuery, u.Login, u.Password)
	return err
}
