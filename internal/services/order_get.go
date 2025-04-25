package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetRepository interface {
	Get(ctx context.Context, filter *repositories.OrderGetFilter) (*repositories.OrderGetDB, error)
}

type OrderGetValidator interface {
	Struct(v any) error
}

type OrderGetService struct {
	v   OrderGetValidator
	ogr OrderGetRepository
}

func NewOrderGetService(
	v OrderGetValidator,
	ogr OrderGetRepository,
) *OrderGetService {
	return &OrderGetService{
		v:   v,
		ogr: ogr,
	}
}

func (svc *OrderGetService) Get(
	ctx context.Context, number string,
) (*OrderResponse, *types.APIStatus) {
	err := svc.v.Struct(number)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid order ID format",
		}
	}

	order, err := svc.ogr.Get(ctx, &repositories.OrderGetFilter{Number: number})
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	if order == nil {
		return nil, &types.APIStatus{
			Status:  http.StatusNotFound,
			Message: "Order not found",
		}
	}
	return &OrderResponse{
		Order:   order.Number,
		Status:  order.Status,
		Accrual: order.Accrual,
	}, nil
}

type OrderResponse struct {
	Order   string `json:"order" `
	Status  string `json:"status"`
	Accrual *int64 `json:"accrual"`
}
