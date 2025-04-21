package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/handlers/helpers"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GetOrderByNumberService interface {
	GetOrderByNumber(ctx context.Context, number string) (*types.OrderResponse, *types.APIError)
}

func GetOrderByNumberHandler(svc GetOrderByNumberService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		number := helpers.GetURLParam(r, "number")
		order, err := svc.GetOrderByNumber(r.Context(), number)
		if err != nil {
			http.Error(w, err.Message, err.Status)
			return
		}
		helpers.EncodeResponseBody(w, order, http.StatusOK)
	}
}
