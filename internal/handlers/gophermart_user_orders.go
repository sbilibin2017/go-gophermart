package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrdersService interface {
	List(ctx context.Context, login string) ([]*services.UserOrdersResponse, *types.APIStatus)
}

func UserOrdersHandler(svc UserOrdersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := middlewares.GetLogin(r)
		if err != nil {
			handleErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
		orders, errStatus := svc.List(r.Context(), login)
		if errStatus != nil {
			handleErrorResponse(w, errStatus.Message, errStatus.Status)
			return
		}
		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		encodeJSONResponse(w, orders, http.StatusOK)
	}
}
