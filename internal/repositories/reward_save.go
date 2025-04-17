package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
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
	ctx context.Context, match map[string]any,
) error {
	query := rewardSaveQuery
	args := []any{match["match"], match["reward"], match["reward_type"]}
	_, err := r.db.ExecContext(ctx, query, args...)
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
