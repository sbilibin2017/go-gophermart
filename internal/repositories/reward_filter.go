package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/queries"
)

type RewardExistsRepository struct {
	db *sqlx.DB
}

func NewRewardExistsRepository(db *sqlx.DB) *RewardExistsRepository {
	return &RewardExistsRepository{db: db}
}

func (r *RewardExistsRepository) Exists(
	ctx context.Context, filter map[string]any,
) (bool, error) {
	if len(filter) == 0 {
		return false, nil
	}

	filters, args := queries.BuildRewardExistsFilters(filter)
	query := fmt.Sprintf(queries.RewardExistsQuery, filters)

	tx := middlewares.TxFromContext(ctx)

	var exists bool
	var err error

	if tx != nil {
		err = tx.GetContext(ctx, &exists, query, args...)
	} else {
		err = r.db.GetContext(ctx, &exists, query, args...)
	}

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}
