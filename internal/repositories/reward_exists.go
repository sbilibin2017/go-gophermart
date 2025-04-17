package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type RewardExistsRepository struct {
	db *sqlx.DB
}

func NewRewardExistsRepository(db *sqlx.DB) *RewardExistsRepository {
	return &RewardExistsRepository{
		db: db,
	}
}

func (r *RewardExistsRepository) Exists(
	ctx context.Context, match *dto.RewardExistsFilterDB,
) (bool, error) {
	logger.Logger.Info("Проверка существования награды для матчевого значения: ", match.Match)

	var exists bool
	query := rewardExistsQuery
	err := r.db.GetContext(ctx, &exists, query, match.Match)
	if err != nil {
		logger.Logger.Error("Ошибка при выполнении запроса: ", err)
		return false, err
	}

	logger.Logger.Info("Запрос выполнен успешно. Награда существует: ", exists)
	return exists, nil
}

var rewardExistsQuery = `SELECT EXISTS (SELECT 1 FROM rewards WHERE match = $1)`
