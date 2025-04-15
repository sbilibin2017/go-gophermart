package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

var (
	ErrOrderIsNotRegistered = errors.New("order is not registered")
)

type OrderRegisterRequest struct {
	Number uint64              `json:"order" validate:"required,louna"`
	Goods  []OrderRegisterGood `json:"goods" validate:"required,dive"`
}

type OrderRegisterGood struct {
	Description string `json:"description" validate:"required,min=1,max=255"`
	Price       uint64 `json:"price" validate:"required,gt=0"`
}

type OrderRegisterService interface {
	Register(ctx context.Context, order *services.Order) error
}

type OrderRegisterDecoder interface {
	Decode(r *http.Request, v any) error
}

type OrderRegisterValidator interface {
	Struct(i interface{}) error
}

func OrderRegisterHandler(
	svc OrderRegisterService,
	dec OrderRegisterDecoder,
	val OrderRegisterValidator,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req OrderRegisterRequest

		if err := dec.Decode(r, &req); err != nil {
			logger.Logger.Warnw("failed to decode request", "error", err)
			http.Error(w, utils.ErrUnprocessableJSON.Error(), http.StatusBadRequest)
			return
		}

		if err := utils.Validate(val, req); err != nil {
			logger.Logger.Warnw("validation failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var order services.Order
		order.Number = req.Number
		for _, g := range req.Goods {
			order.Goods = append(order.Goods, services.Good{
				Description: g.Description,
				Price:       g.Price,
			})
		}

		logger.Logger.Infow("registering order", "order_number", order.Number, "goods_count", len(order.Goods))

		err := svc.Register(r.Context(), &order)
		if err != nil {
			switch err {
			case services.ErrOrderAlreadyRegistered:
				logger.Logger.Infow("order already registered", "order_number", order.Number)
				http.Error(w, err.Error(), http.StatusConflict)
			default:
				logger.Logger.Errorw("failed to register order", "error", err, "order_number", order.Number)
				http.Error(w, ErrOrderIsNotRegistered.Error(), http.StatusInternalServerError)
			}
			return
		}

		logger.Logger.Infow("order registered successfully", "order_number", order.Number)

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		w.Write(newOrderRegisterResponse())
	}
}

func newOrderRegisterResponse() []byte {
	return []byte("Order registered successfully")
}
