package services

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

type GophermartUserOrderService struct {
	ro ListRepository
}

func NewGophermartUserOrderService(
	ro ListRepository,
) *GophermartUserOrderService {
	return &GophermartUserOrderService{
		ro: ro,
	}
}

func (svc *GophermartUserOrderService) List(
	ctx context.Context, login string,
) ([]*GophermartUserOrdersResponse, *APIStatus, *APIStatus) {
	orders, err := svc.ro.List(
		ctx, map[string]any{"login": login}, []string{"number", "status", "accrual", "uploaded_at"}, "uploaded_at", false,
	)
	if err != nil {
		logger.Logger.Errorf("Error getting orders for user %v: %v", login, err)
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Orders are not retrieved",
		}
	}
	if len(orders) == 0 {
		return nil, &APIStatus{
			StatusCode: http.StatusNoContent,
			Message:    "No orders found",
		}, nil
	}
	var response []*GophermartUserOrdersResponse
	err = mapListToStruct(&response, orders)
	if err != nil {
		logger.Logger.Errorf("Error mapping orders: %v", err)
		return nil, nil, &APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Orders are not retrieved",
		}
	}
	return response, &APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Successfully retrieved orders",
	}, nil
}

type GophermartUserOrdersResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    *int64    `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}
