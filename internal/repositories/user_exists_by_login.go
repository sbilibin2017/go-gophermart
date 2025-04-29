package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserExistsByLoginRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewUserExistsByLoginRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *UserExistsByLoginRepository {
	return &UserExistsByLoginRepository{db: db, txProvider: txProvider}
}

func (repo *UserExistsByLoginRepository) ExistsByLogin(
	ctx context.Context, login string,
) (bool, error) {
	params := map[string]any{"login": login}
	var exists bool

	// Передаем txProvider в queryValue для обработки транзакции
	err := queryValue(ctx, repo.db, repo.txProvider, userExistsByLoginQuery, params, &exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const userExistsByLoginQuery = `SELECT EXISTS(SELECT 1 FROM users WHERE login = :login)`
