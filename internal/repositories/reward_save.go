package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/models"
)

const rewardSaveQuery = `
INSERT INTO rewards (match, reward, reward_type, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
ON CONFLICT (match) 
DO UPDATE SET reward = $2, reward_type = $3, updated_at = NOW()
`

type RewardSaveRepository struct {
	db *sqlx.DB
}

func NewRewardSaveRepository(db *sqlx.DB) *RewardSaveRepository {
	return &RewardSaveRepository{db: db}
}

func (repo *RewardSaveRepository) Save(
	ctx context.Context, reward *models.RewardDB,
) error {
	_, err := repo.db.ExecContext(
		ctx,
		rewardSaveQuery,
		reward.Match,
		reward.Reward,
		reward.RewardType,
	)
	if err != nil {
		return err
	}
	return nil
}
