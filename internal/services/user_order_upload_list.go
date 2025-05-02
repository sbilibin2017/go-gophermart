package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderUploadListUserOrderListRepository interface {
	ListOrdered(ctx context.Context, login string) (*[]types.UserOrderDB, error)
}

type UserOrderListService struct {
	uol UserOrderUploadListUserOrderListRepository
}

func NewUserOrderListService(
	uol UserOrderUploadListUserOrderListRepository,
) *UserOrderListService {
	return &UserOrderListService{
		uol: uol,
	}
}

func (svc *UserOrderListService) List(
	ctx context.Context, login string,
) ([]*types.UserOrderUploadedListResponse, *types.APIStatus) {
	orders, err := svc.uol.ListOrdered(ctx, login)
	if err != nil {
		return nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to fetch user orders",
		}
	}

	if len(*orders) == 0 {
		return nil, &types.APIStatus{
			StatusCode: http.StatusNoContent,
			Message:    "No orders found",
		}
	}

	var response []*types.UserOrderUploadedListResponse
	for _, order := range *orders {
		response = append(response, &types.UserOrderUploadedListResponse{
			Number:     order.Number,
			Status:     order.Status,
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt,
		})
	}

	return response, &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Orders fetched successfully",
	}
}
