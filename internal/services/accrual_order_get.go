package services

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type AccrualOrderGetRepository interface {
	FilterByNumber(ctx context.Context, number string, fields []string) (map[string]any, error)
}

type AccrualOrderGetService struct {
	v  *validator.Validate
	ro AccrualOrderGetRepository
}

func NewAccrualOrderGetService(
	v *validator.Validate,
	ro AccrualOrderGetRepository,
) *AccrualOrderGetService {
	return &AccrualOrderGetService{
		v:  v,
		ro: ro,
	}
}

func (svc *AccrualOrderGetService) Get(
	ctx context.Context, req *types.AccrualOrderGetRequest,
) (*types.AccrualOrderGetResponse, *types.APIStatus, *types.APIStatus) {
	if err := svc.v.Struct(req); err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid order number format",
		}
	}
	order, err := svc.ro.FilterByNumber(
		ctx, req.Order, []string{"number", "status", "accrual"},
	)
	if err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error retrieving accrual data",
		}
	}
	if order == nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusNotFound,
			Message:    "Order is not registered",
		}
	}
	var response types.AccrualOrderGetResponse
	err = mapToStruct(&response, order)
	if err != nil {
		return nil, nil, &types.APIStatus{
			StatusCode: http.StatusInternalServerError,
			Message:    "Error mapping accrual data",
		}
	}
	return &response, &types.APIStatus{
		StatusCode: http.StatusOK,
		Message:    "Success",
	}, nil
}
