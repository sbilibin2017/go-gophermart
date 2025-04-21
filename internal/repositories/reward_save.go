package repositories

import (
	"context"
	"database/sql"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
)

type RewardSaveRepository struct {
	db *sql.DB
}

func NewRewardSaveRepository(db *sql.DB) *RewardSaveRepository {
	return &RewardSaveRepository{db: db}
}

func (r *RewardSaveRepository) Save(ctx context.Context, data map[string]any) error {
	executor := contextutils.GetExecutor(ctx, r.db)
	_, err := executor.ExecContext(ctx, rewardSaveQuery,
		data["reward_id"],
		data["reward"],
		data["reward_type"],
	)
	return err
}

const rewardSaveQuery = `
	INSERT INTO rewards (reward_id, reward, reward_type, created_at, updated_at)
	VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	ON CONFLICT (reward_id) DO UPDATE
	SET reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type,
		updated_at = CURRENT_TIMESTAMP
`
