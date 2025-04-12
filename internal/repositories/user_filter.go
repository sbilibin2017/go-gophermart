package repositories

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type UserFilterRepository struct {
	db *sql.DB
}

func NewUserFilterRepository(db *sql.DB) *UserFilterRepository {
	return &UserFilterRepository{db: db}
}

const userFilterQuery = `
	SELECT login, password
	FROM users
	WHERE login = $1
	LIMIT 1
`

type UserFilter struct {
	Login string `json:"login"`
}

type UserFiltered struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *UserFilterRepository) Filter(
	ctx context.Context, tx *sql.Tx, filter *UserFilter,
) (*UserFiltered, error) {
	var userFiltered UserFiltered
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, userFilterQuery, filter.Login).
			Scan(&userFiltered.Login, &userFiltered.Password)
	} else {
		err = r.db.QueryRowContext(ctx, userFilterQuery, filter.Login).
			Scan(&userFiltered.Login, &userFiltered.Password)
	}

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &userFiltered, nil
}
