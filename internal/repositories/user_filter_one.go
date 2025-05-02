package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserFilterOneRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewUserFilterOneRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *UserFilterOneRepository {
	return &UserFilterOneRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *UserFilterOneRepository) FilterOne(
	ctx context.Context, login string,
) (*types.UserDB, error) {
	var user types.UserDB
	err := getContext(
		ctx,
		r.db,
		r.txProvider,
		userFilterOneQuery,
		&user,
		login,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

const userFilterOneQuery = `
	SELECT * 
	FROM users 
	WHERE login = $1
`
