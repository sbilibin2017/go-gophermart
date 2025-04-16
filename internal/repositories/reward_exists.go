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

func (r *RewardExistsRepository) Exists(ctx context.Context, tx *sqlx.Tx, match *RewardExistsFilter) (bool, error) {
	var exists bool
	query := rewardExistsQuery

	var err error
	if tx != nil {
		err = tx.GetContext(ctx, &exists, query, match.Match)
	} else {
		err = r.db.GetContext(ctx, &exists, query, match.Match)
	}

	if err != nil {
		return false, err
	}
	return exists, nil
}

type RewardExistsFilter struct {
	Match string `db:"match"`
}

var rewardExistsQuery = `SELECT EXISTS (SELECT 1 FROM rewards WHERE match = $1)`
