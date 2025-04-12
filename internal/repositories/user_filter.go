package repositories

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type UserFilterRepository struct {
	db *sqlx.DB
}

func NewUserFilterRepository(db *sqlx.DB) *UserFilterRepository {
	return &UserFilterRepository{db: db}
}

const userFilterQuery = `
	SELECT login, password
	FROM users
	WHERE login = $1
	LIMIT 1
`

type UserFilter struct {
	Login string `db:"login"`
}

type UserFiltered struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}

func (r *UserFilterRepository) Filter(
	ctx context.Context, filter *UserFilter,
) (*UserFiltered, error) {
	var userFiltered UserFiltered
	err := r.db.GetContext(ctx, &userFiltered, userFilterQuery, filter.Login)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &userFiltered, nil
}
