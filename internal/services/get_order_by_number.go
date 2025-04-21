package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderFilterRepository interface {
	Filter(ctx context.Context, filter map[string]any, fields []string) (map[string]any, error)
}

type OrderFilterValidator interface {
	Struct(s any) error
}

type GetOrderByNumberService struct {
	v  OrderFilterValidator
	of OrderFilterRepository
}

func NewGetOrderByNumberService(
	v OrderFilterValidator,
	of OrderFilterRepository,
) *GetOrderByNumberService {
	return &GetOrderByNumberService{
		v:  v,
		of: of,
	}
}

func (svc *GetOrderByNumberService) GetOrderByNumber(
	ctx context.Context, number string,
) (*types.OrderResponse, *types.APIError) {
	if err := svc.v.Struct(number); err != nil {
		return nil, types.NewValidationErrorResponse(err)
	}
	filter := map[string]any{
		"order_id": number,
	}
	fields := []string{"order", "status", "accrual"}

	result, err := svc.of.Filter(ctx, filter, fields)
	if err != nil {
		return nil, types.NewInternalError()
	}
	if result == nil {
		return nil, nil
	}

	return types.NewOrderResponse(result), nil
}
