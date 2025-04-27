package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserOrdersListRepository interface {
	ListOrdered(ctx context.Context, login string, fields []string) ([]map[string]any, error)
}

type GophermartUserOrderService struct {
	ro GophermartUserOrdersListRepository
}

func NewGophermartUserOrderService(
	ro GophermartUserOrdersListRepository,
) *GophermartUserOrderService {
	return &GophermartUserOrderService{
		ro: ro,
	}
}

func (svc *GophermartUserOrderService) List(
	ctx context.Context, login string,
) ([]*types.GophermartUserOrdersResponse, *types.APIStatus, *types.APIStatus) {
	orders, err := svc.ro.ListOrdered(
		ctx, login, []string{"number", "status", "accrual", "uploaded_at"},
	)
	if err != nil {
		logger.Logger.Errorf("Error getting orders for user %v: %v", login, err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Orders are not retrieved",
		}
	}
	if len(orders) == 0 {
		return nil, &types.APIStatus{
			StatusCode: http.StatusNoContent,
			Message:    "No orders found",
		}, nil
	}
	var response []*types.GophermartUserOrdersResponse
	err = mapListToStruct(&response, orders)
	if err != nil {
		logger.Logger.Errorf("Error mapping orders: %v", err)
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Orders are not retrieved",
		}
	}
	return response, &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved orders",
	}, nil
}
