package services

import (
	"context"
	"net/http"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrdersRepository interface {
	ListByUser(ctx context.Context, login string) ([]*repositories.UserOrderDB, error)
}

type UserOrdersService struct {
	repo UserOrdersRepository
}

func NewUserOrdersService(repo UserOrdersRepository) *UserOrdersService {
	return &UserOrdersService{
		repo: repo,
	}
}

func (svc *UserOrdersService) List(ctx context.Context, login string) ([]*UserOrdersResponse, *types.APIStatus) {
	orders, err := svc.repo.ListByUser(ctx, login)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}
	if len(orders) == 0 {
		return nil, &types.APIStatus{
			Status:  http.StatusNoContent,
			Message: "No orders found",
		}
	}

	var response []*UserOrdersResponse
	for _, o := range orders {
		resp := &UserOrdersResponse{
			Number:     o.Number,
			Status:     o.Status,
			UploadedAt: o.UploadedAt.Format(time.RFC3339),
		}
		if o.Accrual != nil {
			resp.Accrual = *o.Accrual
		}
		response = append(response, resp)
	}

	return response, nil
}

type UserOrdersResponse struct {
	Number     string `json:"number"`
	Status     string `json:"status"`
	Accrual    int64  `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}
