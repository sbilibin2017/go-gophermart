package repositories

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserSaveRepository struct {
	db *sql.DB
}

func NewUserSaveRepository(db *sql.DB) *UserSaveRepository {
	return &UserSaveRepository{db: db}
}

const userSaveQuery = `
	INSERT INTO users (login, password)
	VALUES ($1, $2);	
`

type UserSave struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *UserSaveRepository) Save(ctx context.Context, u *UserSave) error {
	_, err := r.db.ExecContext(
		ctx,
		userSaveQuery,
		u.Login,
		u.Password,
	)
	if err != nil {
		return err
	}
	return nil
}
