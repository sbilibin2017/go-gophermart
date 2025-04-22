package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetFilterRepository interface {
	GetByID(ctx context.Context, orderID string, fields []string) (map[string]any, error)
}

type OrderGetValidator interface {
	Struct(v any) error
}

type OrderGetService struct {
	v  OrderGetValidator
	of OrderGetFilterRepository
}

func NewOrderGetService(
	v OrderGetValidator,
	of OrderGetFilterRepository,
) *OrderGetService {
	return &OrderGetService{
		v:  v,
		of: of,
	}
}

func (svc *OrderGetService) Get(
	ctx context.Context, req *types.OrderGetRequest,
) (*types.APIResponse[types.OrderGetResponse], *types.APIError) {
	if err := svc.v.Struct(req); err != nil {
		return nil, &types.APIError{
			Status:  http.StatusBadRequest,
			Message: "Invalid order number",
		}
	}

	order, err := svc.of.GetByID(
		ctx,
		req.Number,
		[]string{"order_id", "status", "accrual"},
	)
	if err != nil {
		return nil, &types.APIError{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error while fetching order",
		}
	}
	if order == nil {
		return nil, &types.APIError{
			Status:  http.StatusNoContent,
			Message: "Order is not registered",
		}
	}

	accrual := order["accrual"].(int64)
	response := types.OrderGetResponse{
		Order:   order["order_id"].(string),
		Status:  types.OrderStatus(order["status"].(string)),
		Accrual: &accrual,
	}

	return &types.APIResponse[types.OrderGetResponse]{
		Data:    response,
		Status:  http.StatusOK,
		Message: "Order found",
	}, nil
}
