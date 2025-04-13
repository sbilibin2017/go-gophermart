package repositories

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type RewardFilterRepository struct {
	db *sqlx.DB
}

func NewRewardFilterRepository(db *sqlx.DB) *RewardFilterRepository {
	return &RewardFilterRepository{db: db}
}

type RewardFilter struct {
	Description string `db:"description"`
}

type RewardFiltered struct {
	Match      string `db:"match"`       // Название товара
	Reward     uint64 `db:"reward"`      // Баллы награды
	RewardType string `db:"reward_type"` // Тип награды
}

var rewardFilterQuery = `
	SELECT reward, reward_type, match
	FROM rewards
	WHERE match ILIKE ANY($1)
`

func getRewardFilterQueryArgs(filter []*RewardFilter) []interface{} {
	matches := make([]string, 0, len(filter))
	for _, f := range filter {
		matches = append(matches, "%"+f.Description+"%")
	}
	return []interface{}{matches} // Можно использовать напрямую с pgx
}

func (r *RewardFilterRepository) Filter(ctx context.Context, filter []*RewardFilter) ([]*RewardFiltered, error) {
	args := getRewardFilterQueryArgs(filter)
	var rewards []*RewardFiltered
	err := r.db.SelectContext(ctx, &rewards, rewardFilterQuery, args...)
	if err != nil {
		return nil, err
	}
	return rewards, nil
}
