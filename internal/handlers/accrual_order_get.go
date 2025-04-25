package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GetByIDService interface {
	Get(ctx context.Context, number string) (*services.OrderResponse, *types.APIStatus)
}

func OrderGetByIDHandler(svc GetByIDService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		number := getPathParam(r, "number")
		order, err := svc.Get(r.Context(), number)
		if err != nil {
			handleErrorResponse(w, err.Message, err.Status)
			return
		}
		encodeJSONResponse(w, order, http.StatusOK)
	}
}
