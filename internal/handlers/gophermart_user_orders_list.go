package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type GophermartUserOrderService interface {
	List(ctx context.Context, login string) ([]*services.GophermartUserOrdersResponse, *services.APIStatus, *services.APIStatus)
}

func GophermartUserOrdersListHandler(svc GophermartUserOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getUserLoginFromContext(w, r, jwt.GetJWTPayload)
		if login == nil {
			return
		}
		orders, status, err := svc.List(r.Context(), *login)
		if err != nil {
			handleError(w, err)
			return
		}
		if len(orders) == 0 {
			handleError(w, status)
			return
		}
		encodeJSONResponse(w, orders, status)
	}
}
