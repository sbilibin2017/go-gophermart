package repositories

import (
	"context"
	"database/sql"
)

type RewardSaveRepository interface {
	Save(ctx context.Context, match string, reward uint, rewardType string) error
}

type RewardSaveRepositoryImpl struct {
	db *sql.DB
}

func NewRewardSaveRepository(db *sql.DB) *RewardSaveRepositoryImpl {
	return &RewardSaveRepositoryImpl{
		db: db,
	}
}

func (r *RewardSaveRepositoryImpl) Save(
	ctx context.Context, match string, reward uint, rewardType string,
) error {
	_, err := r.db.ExecContext(ctx, rewardSaveQuery, match, reward, rewardType)
	if err != nil {
		return err
	}
	return nil
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
