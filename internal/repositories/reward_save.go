package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
)

type RewardSaveRepository struct {
	db *sqlx.DB
}

func NewRewardSaveRepository(db *sqlx.DB) *RewardSaveRepository {
	return &RewardSaveRepository{
		db: db,
	}
}

func (r *RewardSaveRepository) Save(
	ctx context.Context, reward *dto.RewardDB,
) error {
	query := rewardSaveQuery
	args := []interface{}{reward.Match, reward.Reward, reward.RewardType}
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
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
