package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
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
	logger.Logger.Info("Попытка сохранить или обновить награду для матча:", reward.Match)
	logger.Logger.Info("Параметры для сохранения: reward =", reward.Reward, ", reward_type =", reward.RewardType)

	query := rewardSaveQuery
	args := []interface{}{reward.Match, reward.Reward, reward.RewardType}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Logger.Error("Ошибка при выполнении запроса для сохранения награды:", err)
		return err
	}

	logger.Logger.Info("Награда успешно сохранена или обновлена для матча:", reward.Match)
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
