package repositories

import (
	"context"
	"fmt"
)

type AccrualOrderFilterByNumberRepository struct {
	q Querier
}

func NewAccrualOrderFilterByNumberRepository(
	q Querier,
) *AccrualOrderFilterByNumberRepository {
	return &AccrualOrderFilterByNumberRepository{q: q}
}

// Implementing the Filter method to match the FilterRepository interface
func (repo *AccrualOrderFilterByNumberRepository) Filter(
	ctx context.Context,
	filter map[string]any,
	fields []string,
) (map[string]any, error) {
	q, args := buildFilterByNumberQuery(fields, filter["number"])
	var result map[string]any
	err := repo.q.Query(ctx, q, &result, args)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildFilterByNumberQuery(fields []string, number any) (string, map[string]any) {
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
