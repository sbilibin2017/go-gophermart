package repositories

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserGetByParamRepository struct {
	db *sql.DB
}

func NewUserGetByParamRepository(db *sql.DB) *UserGetByParamRepository {
	return &UserGetByParamRepository{db: db}
}

const userGetByParamQuery = `
	SELECT login, password
	FROM users
	WHERE login = $1
	LIMIT 1
`

func (r *UserGetByParamRepository) GetByParam(
	ctx context.Context, p map[string]any,
) (map[string]any, error) {
	var login, password string
	err := r.db.QueryRowContext(
		ctx,
		userGetByParamQuery,
		p["login"],
	).Scan(&login, &password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user := map[string]any{
		"login":    login,
		"password": password,
	}
	return user, nil
}
