package repositories

import (
	"context"
	"database/sql"

	"github.com/sbilibin2017/go-gophermart/internal/types"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const userGetByParamQuery = `
	SELECT login, password
	FROM users
	WHERE login = $1
	LIMIT 1
`

type UserGetByParamRepository struct {
	db *sql.DB
}

func NewUserGetByParamRepository(db *sql.DB) *UserGetByParamRepository {
	return &UserGetByParamRepository{db: db}
}

func (r *UserGetByParamRepository) GetByParam(
	ctx context.Context, p *types.UserGetParam,
) (*types.User, error) {
	var user types.User
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
