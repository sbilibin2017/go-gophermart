package repositories

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type AccrualOrderExistsRepository struct {
	db *sqlx.DB
}

func NewAccrualOrderExistsRepository(db *sqlx.DB) *AccrualOrderExistsRepository {
	return &AccrualOrderExistsRepository{db: db}
}

func (repo *AccrualOrderExistsRepository) Exists(
	ctx context.Context,
	number string,
) (bool, error) {
	var exists bool
	err := query(ctx, repo.db, accrualOrderRegisterExistsQuery, &exists, number)
	if err != nil {
		return false, err
	}
	return exists, nil
}

const accrualOrderRegisterExistsQuery = `
	SELECT EXISTS (
		SELECT 1
		FROM accrual_order
		WHERE number = $1
	)
`
