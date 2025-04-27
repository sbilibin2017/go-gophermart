package repositories

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AccrualOrderFilterByNumberRepository struct {
	db *sqlx.DB
}

func NewAccrualOrderFilterByNumberRepository(db *sqlx.DB) *AccrualOrderFilterByNumberRepository {
	return &AccrualOrderFilterByNumberRepository{db: db}
}

func (repo *AccrualOrderFilterByNumberRepository) FilterByNumber(
	ctx context.Context,
	number string,
	fields []string,
) (map[string]any, error) {
	q, args := buildFilterByNumberQuery(fields, number)
	var result map[string]any
	err := query(ctx, repo.db, q, &result, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildFilterByNumberQuery(fields []string, number string) (string, map[string]any) {
	columns := buildColumnsString(fields)
	query := fmt.Sprintf(
		accrualOrderFilterByNumberQuery,
		columns,
	)
	args := map[string]any{
		"number": number,
	}
	return query, args
}

const accrualOrderFilterByNumberQuery = `
	SELECT %s FROM accrual_order WHERE number = :number
`
