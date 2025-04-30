package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GoodRewardSaveRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewGoodRewardSaveRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *GoodRewardSaveRepository {
	return &GoodRewardSaveRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *GoodRewardSaveRepository) Save(
	ctx context.Context, goodReward *types.GoodReward,
) error {
	e := getExecutor(ctx, r.db, r.txProvider)
	_, err := sqlx.NamedExecContext(ctx, e, goodRewardUpsertQuery, goodReward)
	return err
}

const goodRewardUpsertQuery = `
INSERT INTO good_rewards (match, reward, reward_type)
VALUES (:match, :reward, :reward_type)
ON CONFLICT (match)
DO UPDATE SET
	reward = EXCLUDED.reward,
	reward_type = EXCLUDED.reward_type
`
