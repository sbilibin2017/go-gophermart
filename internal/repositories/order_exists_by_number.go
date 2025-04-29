package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderExistsByNumberRepository struct {
	db *sqlx.DB
}

func NewOrderExistsByNumberRepository(db *sqlx.DB) *OrderExistsByNumberRepository {
	return &OrderExistsByNumberRepository{db: db}
}

func (repo *OrderExistsByNumberRepository) ExistByNumber(
	ctx context.Context, number string, login *string,
) (bool, error) {
	var exists bool
	query, args := buildOrderExistsByNumberQuery(number, login)
	if err := queryRow(ctx, repo.db, query, &exists, args); err != nil {
		return false, err
	}
	return exists, nil
}

func buildOrderExistsByNumberQuery(number string, login *string) (string, map[string]any) {
	var query string
	var args map[string]any
	if login != nil {
		query = orderExistsByNumberWithLoginQuery
		args = map[string]any{
			"number": number,
			"login":  *login,
		}
	} else {
		query = orderExistsByNumberQuery
		args = map[string]any{
			"number": number,
		}
	}
	return query, args
}

const orderExistsByNumberWithLoginQuery = `
SELECT EXISTS(SELECT 1 FROM orders WHERE number = :number AND login = :login)
`

const orderExistsByNumberQuery = `
SELECT EXISTS(SELECT 1 FROM orders WHERE number = :number)
`
