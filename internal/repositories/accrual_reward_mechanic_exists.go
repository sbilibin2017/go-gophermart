package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type AccrualRewardMechanicExistsRepository struct {
	db *sqlx.DB
}

func NewAccrualRewardMechanicExistsRepository(db *sqlx.DB) *AccrualRewardMechanicExistsRepository {
	return &AccrualRewardMechanicExistsRepository{db: db}
}

func (r *AccrualRewardMechanicExistsRepository) Exists(
	ctx context.Context,
	match string,
) (bool, error) {
	var exists bool
	err := query(ctx, r.db, checkAccrualRewardMechanicExistsQuery, &exists, match)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const checkAccrualRewardMechanicExistsQuery = `
	SELECT EXISTS (
		SELECT 1
		FROM accrual_reward_mechanic
		WHERE match = $1
	)
`
