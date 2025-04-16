package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/models"
)

const rewardExistsQuery = `SELECT EXISTS (SELECT 1 FROM rewards WHERE match = $1)`

type RewardExistsRepository struct {
	db *sqlx.DB
}

func NewRewardExistsRepository(db *sqlx.DB) *RewardExistsRepository {
	return &RewardExistsRepository{db: db}
}

func (repo *RewardExistsRepository) Exists(
	ctx context.Context, filter *models.RewardFilter,
) (bool, error) {
	var exists bool
	err := repo.db.GetContext(ctx, &exists, rewardExistsQuery, filter.Match)
	if err != nil {
		return false, err
	}
	return exists, nil
}
