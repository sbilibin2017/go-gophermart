package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type GophermartUserBalanceWithdrawalsListService interface {
	List(ctx context.Context, login string) ([]*services.GophermartUserBalanceWithdrawalsResponse, *services.APIStatus, *services.APIStatus)
}

func GophermartUserBalanceWithdrawalsListHandler(
	svc GophermartUserBalanceWithdrawalsListService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getUserLoginFromContext(w, r, jwt.GetJWTPayload)
		if login == nil {
			return
		}
		withdrawals, status, err := svc.List(r.Context(), *login)
		if err != nil {
			handleError(w, err)
			return
		}
		if len(withdrawals) == 0 {
			sendTextPlainResponse(w, status)
			return
		}
		encodeJSONResponse(w, withdrawals, status)
	}
}
