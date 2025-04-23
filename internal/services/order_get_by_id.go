package services

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

const (
	ErrInvalidOrderNumber     = "Invalid order number"
	ErrInternalServerErrorGet = "Internal server error while fetching order"
	ErrOrderNotRegistered     = "Order is not registered"
	SuccessOrderFetched       = "Order fetched successfully"
)

type OrderRepository interface {
	GetByID(ctx context.Context, orderID string, fields []string) (*types.OrderDB, error)
}

type Validator interface {
	Struct(v any) error
}

type OrderGetService struct {
	validator  Validator
	repository OrderRepository
}

func NewOrderGetService(
	validator Validator,
	repository OrderRepository,
) *OrderGetService {
	return &OrderGetService{
		validator:  validator,
		repository: repository,
	}
}

func (svc *OrderGetService) GetByID(
	ctx context.Context, req *types.OrderGetByIDRequest,
) (*types.OrderGetByIDResponse, *types.APIStatus, error) {
	if err := svc.validator.Struct(req); err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusBadRequest,
			Message: ErrInvalidOrderNumber,
		}, err
	}

	order, err := svc.repository.GetByID(
		ctx,
		req.Number,
		[]string{"order_id", "status", "accrual"},
	)
	if err != nil {
		return nil, &types.APIStatus{
			Status:  http.StatusInternalServerError,
			Message: ErrInternalServerErrorGet,
		}, err
	}
	if order == nil {
		return nil, &types.APIStatus{
			Status:  http.StatusNoContent,
			Message: ErrOrderNotRegistered,
		}, nil
	}

	response := &types.OrderGetByIDResponse{
		Order:   order.OrderID,
		Status:  types.OrderStatus(order.Status),
		Accrual: order.Accrual,
	}

	status := &types.APIStatus{
		Status:  http.StatusOK,
		Message: SuccessOrderFetched,
	}

	return response, status, nil
}
