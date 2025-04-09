package repositories

import (
	"context"
	"database/sql"

	"github.com/sbilibin2017/go-gophermart/internal/types"

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

func (r *UserSaveRepository) Save(ctx context.Context, u *types.User) error {
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
