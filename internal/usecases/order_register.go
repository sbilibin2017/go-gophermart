package usecases

import (
	"context"
	"errors"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderValidator interface {
	Struct(i interface{}) error
}

type OrderRegisterService interface {
	Register(ctx context.Context, order *types.Order) error
}

type OrderRegisterUsecase struct {
	ov  OrderValidator
	svc OrderRegisterService
}

func NewOrderRegisterUsecase(
	ov OrderValidator,
	svc OrderRegisterService,
) *OrderRegisterUsecase {
	return &OrderRegisterUsecase{
		ov:  ov,
		svc: svc,
	}
}

func (uc *OrderRegisterUsecase) Execute(
	ctx context.Context, req *types.OrderRegisterRequest,
) (*types.OrderRegisterResponse, error) {
	err := uc.ov.Struct(req)
	if err != nil {
		return nil, ErrOrderRegisterInvalidRequest
	}
	var goods []types.Good
	for _, g := range req.Goods {
		goods = append(goods, types.Good{
			Description: g.Description,
			Price:       g.Price,
		})
	}
	err = uc.svc.Register(ctx, &types.Order{Number: req.Number, Goods: goods})
	if err != nil {
		return nil, err
	}
	return &types.OrderRegisterResponse{
		Message: "order registered successfully",
	}, nil
}

var (
	ErrOrderRegisterInvalidRequest = errors.New("invalid order register request")
)
