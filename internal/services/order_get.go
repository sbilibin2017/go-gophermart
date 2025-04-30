package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetOrderFindRepository interface {
	Find(ctx context.Context, number string) (*types.Order, error)
}

type OrderGetService struct {
	of OrderGetOrderFindRepository
}

func NewOrderGetService(
	of OrderGetOrderFindRepository,
) *OrderGetService {
	return &OrderGetService{
		of: of,
	}
}

func (s *OrderGetService) Get(
	ctx context.Context, number string,
) (*types.Order, error) {
	order, err := s.of.Find(ctx, number)
	if err != nil {
		return nil, err
	}
	return order, nil
}
