package repositories

import (
	"context"
	"database/sql"

	// Importing the standard Go log package
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardMechanicFilterOneRepository struct {
	db *sqlx.DB
}

func NewRewardMechanicFilterOneRepository(
	db *sqlx.DB,
) *RewardMechanicFilterOneRepository {
	return &RewardMechanicFilterOneRepository{
		db: db,
	}
}

func (r *RewardMechanicFilterOneRepository) FilterOne(
	ctx context.Context, match string,
) (*types.RewardMechanicDB, error) {
	var reward types.RewardMechanicDB
	err := r.db.GetContext(ctx, &reward, rewardMechanicFilterOneQuery, match)
	if err != nil {
		logQuery(rewardMechanicFilterOneQuery, nil, err)
		if err == sql.ErrNoRows {
			return nil, nil

		}
		return nil, err
	}
	logQuery(rewardMechanicFilterOneQuery, nil, nil)
	return &reward, nil
}

const rewardMechanicFilterOneQuery = `
	SELECT *
	FROM reward_mechanics
	WHERE match = $1
`
