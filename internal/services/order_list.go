package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type OrderListRepository interface {
	ListByLoginOrderedDesc(ctx context.Context, login string) ([]*domain.Order, error)
}

type OrderListService struct {
	repo OrderListRepository
}

func NewOrderListService(repo OrderListRepository) *OrderListService {
	return &OrderListService{
		repo: repo,
	}
}

func (svc *OrderListService) List(
	ctx context.Context, login string,
) ([]*domain.Order, error) {
	orders, err := svc.repo.ListByLoginOrderedDesc(ctx, login)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
