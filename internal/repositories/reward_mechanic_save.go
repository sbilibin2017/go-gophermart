package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type RewardMechanicSaveRepository struct {
	db *sqlx.DB
}

func NewRewardMechanicSaveRepository(
	db *sqlx.DB,
) *RewardMechanicSaveRepository {
	return &RewardMechanicSaveRepository{
		db: db,
	}
}

func (r *RewardMechanicSaveRepository) Save(
	ctx context.Context, match string, reward int64, rewardType string,
) error {
	_, err := r.db.ExecContext(ctx, rewardMechanicSaveQuery, match, reward, rewardType)
	if err != nil {
		logQuery(rewardMechanicFilterOneQuery, nil, err)
		return err
	}

	logQuery(rewardMechanicFilterOneQuery, nil, err)

	return nil
}

const rewardMechanicSaveQuery = `
	INSERT INTO reward_mechanics (match, reward, reward_type)
	VALUES ($1, $2, $3)
	ON CONFLICT (match)
	DO UPDATE SET
		reward = EXCLUDED.reward,
		reward_type = EXCLUDED.reward_type
`
