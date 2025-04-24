package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetRepository interface {
	Get(ctx context.Context, number string, fields []string) (map[string]any, error)
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
	ctx context.Context, order string,
) (*types.OrderResponse, *types.APIStatus) {
	err := svc.v.Struct(order)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: "Invalid order ID format",
		}
	}

	data, err := svc.ogr.Get(ctx, order, []string{"order", "status", "accrual"})
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	if data == nil {
		return nil, &types.APIStatus{
			Status:  http.StatusNotFound,
			Message: "Order not found",
		}
	}

	var response types.OrderResponse
	err = convertToStruct(&response, data)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Error converting data to structure",
		}
	}

	return &response, nil
}
