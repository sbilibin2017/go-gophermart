package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetOrderFilterOneRepository interface {
	FilterOne(ctx context.Context, number string) (*types.OrderDB, error)
}

type OrderGetOrderOrderValidator interface {
	Struct(v any) error
}

type OrderGetService struct {
	v  OrderGetOrderOrderValidator
	of OrderGetOrderFilterOneRepository
}

func NewOrderGetService(
	v OrderGetOrderOrderValidator,
	of OrderGetOrderFilterOneRepository,
) *OrderGetService {
	return &OrderGetService{
		v:  v,
		of: of,
	}
}

func (svc *OrderGetService) Get(
	ctx context.Context, req *types.OrderGetRequest,
) (*types.OrderGetResponse, *types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		valErr := formatValidationError(err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    valErr.Message,
		}
	}

	order, err := svc.of.FilterOne(ctx, req.Number)
	if err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Internal server error",
		}
	}
	if order == nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusNotFound,
			Message:    "Order not found",
		}
	}

	response := &types.OrderGetResponse{
		Order:   order.Number,
		Status:  order.Status,
		Accrual: order.Accrual,
	}

	return response, &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Order retrieved successfully",
	}, nil
}
