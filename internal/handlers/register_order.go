package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/handlers/helpers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RegisterOrderService interface {
	Register(ctx context.Context, order *types.RegisterOrderRequest) error
}

type RegisterOrderValidator interface {
	Struct(s any) error
}

func RegisterOrderHandler(
	val RegisterOrderValidator,
	svc RegisterOrderService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterOrderRequest

		if err := helpers.DecodeJSONBody(w, r, &req); err != nil {
			helpers.ErrorInternalServerResponse(w, err)
			return
		}

		if err := val.Struct(req); err != nil {
			helpers.ErrorBadRequestResponse(w, err)
			return
		}

		err := svc.Register(r.Context(), &req)
		if err != nil {
			switch err {
			case services.ErrRegisterOrderAlreadyExists:
				helpers.ErrorConflictResponse(w, err)
			case services.ErrRegisterOrderIsNotRegistered:
				helpers.ErrorBadRequestResponse(w, err)
			default:
				helpers.ErrorInternalServerResponse(w, err)
			}
			return
		}

		helpers.SendTextResponse(w, http.StatusOK, "Order registered successfully")
	}
}
