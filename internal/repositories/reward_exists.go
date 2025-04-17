package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/dto"
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
	var exists bool
	query := rewardExistsQuery
	err := r.db.GetContext(ctx, &exists, query, match.Match)
	if err != nil {
		return false, err
	}
	return exists, nil
}

var rewardExistsQuery = `SELECT EXISTS (SELECT 1 FROM rewards WHERE match = $1)`
