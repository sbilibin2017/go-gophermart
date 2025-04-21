package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/handlers/helpers"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RegisterOrderService interface {
	Register(ctx context.Context, order *types.RegisterOrderRequest) (*string, *types.APIError)
}

func RegisterOrderHandler(svc RegisterOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterOrderRequest
		if err := helpers.DecodeRequestBody(w, r, &req); err != nil {
			return
		}
		msg, err := svc.Register(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Message, err.Status)
			return
		}
		helpers.SendTextResponse(w, *msg, http.StatusOK)
	}
}
