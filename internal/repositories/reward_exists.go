package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RewardExistsRepository struct {
	db *sqlx.DB
}

func NewRewardExistsRepository(db *sqlx.DB) *RewardExistsRepository {
	return &RewardExistsRepository{
		db: db,
	}
}

func (r *RewardExistsRepository) Exists(
	ctx context.Context, filter map[string]any,
) (bool, error) {
	var exists bool
	query := rewardExistsQuery
	err := r.db.GetContext(ctx, &exists, query, filter["match"])
	if err != nil {
		return false, err
	}
	return exists, nil
}

var rewardExistsQuery = `SELECT EXISTS (SELECT 1 FROM rewards WHERE match = $1)`
