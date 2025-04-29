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
	var exists bool
	args := map[string]any{"login": login}
	if err := queryRow(ctx, repo.db, userExistsByLoginQuery, &exists, args); err != nil {
		return false, err
	}
	return exists, nil
}

const userExistsByLoginQuery = `SELECT EXISTS(SELECT 1 FROM users WHERE login = :login)`
