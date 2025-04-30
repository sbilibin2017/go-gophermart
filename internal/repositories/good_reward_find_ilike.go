package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GoodRewardFindILikeRepository struct {
	db         *sqlx.DB
	txProvider func(ctx context.Context) (*sqlx.Tx, error)
}

func NewGoodRewardFindILikeRepository(
	db *sqlx.DB,
	txProvider func(ctx context.Context) (*sqlx.Tx, error),
) *GoodRewardFindILikeRepository {
	return &GoodRewardFindILikeRepository{
		db:         db,
		txProvider: txProvider,
	}
}

func (r *GoodRewardFindILikeRepository) FindILike(ctx context.Context, description string) (*types.GoodReward, error) {
	e := getExecutor(ctx, r.db, r.txProvider)
	var goodReward types.GoodReward
	err := sqlx.GetContext(ctx, e, &goodReward, goodRewardFindILikeQuery, "%"+description+"%")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &goodReward, nil
}

const goodRewardFindILikeQuery = `
SELECT match, reward, reward_type
FROM good_rewards
WHERE LOWER(match) ILIKE LOWER($1)
LIMIT 1
`

// Создаем ошибку для случая, когда товар с таким описанием не найден
var ErrGoodRewardNotFound = errors.New("good reward not found")
