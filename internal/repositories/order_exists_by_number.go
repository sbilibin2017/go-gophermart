package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type OrderExistsByNumberRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewOrderExistsByNumberRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *OrderExistsByNumberRepository {
	return &OrderExistsByNumberRepository{db: db, txProvider: txProvider}
}

func (repo *OrderExistsByNumberRepository) ExistByNumber(
	ctx context.Context, number string, login *string,
) (bool, error) {
	query, params := buildOrderExistsQuery(number, login)
	var exists bool
	err := queryValue(ctx, repo.db, repo.txProvider, query, params, &exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func buildOrderExistsQuery(number string, login *string) (string, map[string]any) {
	if login != nil {
		return orderExistsByNumberWithLoginQuery, map[string]any{
			"number": number,
			"login":  *login,
		}
	}
	return orderExistsByNumberQuery, map[string]any{
		"number": number,
	}
}

const orderExistsByNumberWithLoginQuery = `
SELECT EXISTS(SELECT 1 FROM orders WHERE number = :number AND login = :login)
`

const orderExistsByNumberQuery = `
SELECT EXISTS(SELECT 1 FROM orders WHERE number = :number)
`
