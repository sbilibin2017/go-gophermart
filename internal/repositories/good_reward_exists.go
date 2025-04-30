package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type GoodRewardExistsRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewGoodRewardExistsRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *GoodRewardExistsRepository {
	return &GoodRewardExistsRepository{db: db, txProvider: txProvider}
}

func (r *GoodRewardExistsRepository) Exists(ctx context.Context, match string) (bool, error) {
	e := getExecutor(ctx, r.db, r.txProvider)
	var exists bool
	err := sqlx.GetContext(ctx, e, &exists, goodRewardExistsQuery, match)
	return exists, err
}

const goodRewardExistsQuery = `
SELECT EXISTS (
	SELECT 1 FROM good_rewards WHERE match = $1
)
`
