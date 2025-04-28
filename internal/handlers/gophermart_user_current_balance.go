package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type GophermartUserCurrentBalanceService interface {
	Get(ctx context.Context, login string) (*services.GophermartUserCurrentBalanceResponse, *services.APIStatus, *services.APIStatus)
}

func GophermartUserCurrentBalanceHandler(svc GophermartUserCurrentBalanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getUserLoginFromContext(w, r, jwt.GetJWTPayload)
		if login == nil {
			return
		}
		resp, status, err := svc.Get(r.Context(), *login)
		if err != nil {
			handleError(w, status)
			return
		}
		encodeJSONResponse(w, resp, status)
	}
}
