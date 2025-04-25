package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserExistsRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) *sqlx.Tx
}

func NewUserExistsRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) *sqlx.Tx,
) *UserExistsRepository {
	return &UserExistsRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserExistsRepository) Exists(ctx context.Context, filter *UserExistsLogin) (bool, error) {
	var exists bool
	err := query(ctx, r.db, r.txProvider, &exists, userExistsByLoginQuery, filter)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type UserExistsLogin struct {
	Login string `db:"login"`
}

const userExistsByLoginQuery = `
	SELECT EXISTS (
		SELECT 1 FROM users WHERE login = :login
	)
`
