package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type RewardFilterDB struct {
	Description string `db:"description"`
}

type RewardFilteredDB struct {
	Match      string `db:"match"`
	Reward     uint64 `db:"reward"`
	RewardType string `db:"reward_type"`
}

type RewardFilterRepository struct {
	db *sqlx.DB
}

func NewRewardFilterRepository(db *sqlx.DB) *RewardFilterRepository {
	return &RewardFilterRepository{db: db}
}

func (r *RewardFilterRepository) Filter(
	ctx context.Context, filter *RewardFilterDB,
) (*RewardFilteredDB, error) {
	placeholder := newDescriptionPlaceholder(filter)

	logger.Logger.Debugf("Filtering reward with description: %s (placeholder: %s)", filter.Description, placeholder)

	var reward RewardFilteredDB
	err := r.db.GetContext(ctx, &reward, rewardFilterQuery, placeholder)
	if err != nil {
		logger.Logger.Errorf("Failed to filter reward for description=%s: %v", filter.Description, err)
		return nil, err
	}

	logger.Logger.Debugf("Reward matched: %+v", reward)
	return &reward, nil
}

var rewardFilterQuery = `
	SELECT match, reward, reward_type
	FROM rewards
	WHERE match ILIKE $1
`

func newDescriptionPlaceholder(filter *RewardFilterDB) string {
	return "%" + filter.Description + "%"
}
