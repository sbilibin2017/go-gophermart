package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RewardExistsRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, bool)
}

func NewRewardExistsRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, bool),
) *RewardExistsRepository {
	return &RewardExistsRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *RewardExistsRepository) Exists(ctx context.Context, filter map[string]any) (bool, error) {
	var exists bool
	query := rewardExistsQuery

	if tx, ok := r.txProvider(ctx); ok {
		err := tx.GetContext(ctx, &exists, query, filter["match"])
		return exists, err
	}

	err := r.db.GetContext(ctx, &exists, query, filter["match"])
	return exists, err
}

var rewardExistsQuery = `SELECT EXISTS (SELECT 1 FROM rewards WHERE match = $1)`
