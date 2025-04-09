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

type UserGetParam struct {
	Login string
}

type UserGet struct {
	Login    string
	Password string
}

func (r *UserGetByParamRepository) GetByParam(
	ctx context.Context, p *UserGetParam,
) (*UserGet, error) {
	var user UserGet
	err := r.db.QueryRowContext(
		ctx,
		userGetByParamQuery,
		p.Login,
	).Scan(&user.Login, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
