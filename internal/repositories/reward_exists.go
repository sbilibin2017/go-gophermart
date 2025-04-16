package repositories

import (
	"context"
	"database/sql"
)

type RewardExistsRepository struct {
	db *sql.DB
}

func NewRewardExistsRepository(db *sql.DB) *RewardExistsRepository {
	return &RewardExistsRepository{
		db: db,
	}
}

func (r *RewardExistsRepository) Exists(ctx context.Context, match string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, rewardExistsQuery, match).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

var rewardExistsQuery = `SELECT EXISTS (SELECT 1 FROM rewards WHERE match = $1)`
