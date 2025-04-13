package usecases

import (
	"context"
	"errors"

	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type OrderValidator interface {
	Struct(i interface{}) error
}

type OrderRegisterService interface {
	Register(ctx context.Context, order *services.Order) error
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

type OrderRegisterRequest struct {
	Order uint64 `json:"order" validate:"required"`
	Goods []Good `json:"goods" validate:"dive,required"`
}

type Good struct {
	Description string `json:"description" validate:"required"`
	Price       uint64 `json:"price" validate:"required,min=1"`
}

type OrderRegisterResponse struct {
	Message string `json:"message"`
}

func (uc *OrderRegisterUsecase) Execute(
	ctx context.Context, req *OrderRegisterRequest,
) (*OrderRegisterResponse, error) {
	err := uc.ov.Struct(req)
	if err != nil {
		return nil, ErrOrderRegisterInvalidRequest
	}
	var goods []services.Good
	for _, g := range req.Goods {
		goods = append(goods, services.Good{
			Description: g.Description,
			Price:       g.Price,
		})
	}
	err = uc.svc.Register(ctx, &services.Order{Order: req.Order, Goods: goods})
	if err != nil {
		return nil, err
	}
	return &OrderRegisterResponse{
		Message: "order registered successfully",
	}, nil
}

var (
	ErrOrderRegisterInvalidRequest = errors.New("invalid order register request")
)
