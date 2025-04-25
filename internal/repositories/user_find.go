package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserFindRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewUserFindRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *UserFindRepository {
	return &UserFindRepository{db: db, txProvider: txProvider}
}

func (r *UserFindRepository) Find(ctx context.Context, filter *UserFindFilter) (*UserFindDB, error) {
	var user UserFindDB
	err := query(ctx, r.db, r.txProvider, &user, userFindByLoginQuery, filter)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type UserFindFilter struct {
	Login string `db:"login"`
}

type UserFindDB struct {
	Login    string `db:"login"`
	Password string `db:"password"`
}

const userFindByLoginQuery = `
	SELECT id, login, password
	FROM users
	WHERE login = :login
	LIMIT 1
`
