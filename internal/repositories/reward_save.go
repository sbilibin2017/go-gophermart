package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RewardSaveRepository interface {
	Save(ctx context.Context, match string, reward uint, rewardType string) error
}

type RewardSaveRepositoryImpl struct {
	db *sqlx.DB
}

func NewRewardSaveRepository(db *sqlx.DB) *RewardSaveRepositoryImpl {
	return &RewardSaveRepositoryImpl{
		db: db,
	}
}

func (r *RewardSaveRepositoryImpl) Save(
	ctx context.Context, tx *sqlx.Tx, reward *RewardSave,
) error {
	query := rewardSaveQuery
	args := []interface{}{reward.Match, reward.Reward, reward.RewardType}

	var err error
	if tx != nil {
		_, err = tx.ExecContext(ctx, query, args...)
	} else {
		_, err = r.db.ExecContext(ctx, query, args...)
	}

	return err
}

type RewardSave struct {
	Match      string `db:"match"`
	Reward     uint   `db:"reward"`
	RewardType string `db:"reward_type"`
}

var rewardSaveQuery = `
INSERT INTO rewards (match, reward, reward_type, created_at, updated_at) 
VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (match) 
DO UPDATE SET 
    reward = EXCLUDED.reward,
    reward_type = EXCLUDED.reward_type,
    updated_at = CURRENT_TIMESTAMP
`
