package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type UserExistsByLoginRepository struct {
	db *sqlx.DB
}

func NewUserExistsByLoginRepository(db *sqlx.DB) *UserExistsByLoginRepository {
	return &UserExistsByLoginRepository{db: db}
}

func (repo *UserExistsByLoginRepository) ExistsByLogin(
	ctx context.Context, login string,
) (bool, error) {
	params := map[string]any{"login": login}
	var exists bool
	err := queryValue(ctx, repo.db, userExistsByLoginQuery, params, &exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const userExistsByLoginQuery = `SELECT EXISTS(SELECT 1 FROM users WHERE login = :login)`
