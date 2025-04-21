package repositories

import (
	"context"
	"database/sql"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
)

type RewardExistsRepository struct {
	db *sql.DB
}

func NewRewardExistsRepository(db *sql.DB) *RewardExistsRepository {
	return &RewardExistsRepository{db: db}
}

func (r *RewardExistsRepository) Exists(ctx context.Context, filter map[string]any) (bool, error) {
	executor := contextutils.GetExecutor(ctx, r.db)
	var exists bool
	err := executor.QueryRowContext(ctx, rewardExistsQuery, filter["reward_id"]).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const rewardExistsQuery = `
	SELECT EXISTS(
		SELECT 1 FROM rewards WHERE reward_id = $1
	)
`
