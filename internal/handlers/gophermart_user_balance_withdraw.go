package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type GophermartUserBalanceWithdrawService interface {
	Withdraw(ctx context.Context, req *services.GophermartUserBalanceWithdrawRequest, login string) (*services.APIStatus, *services.APIStatus)
}

func GophermartUserBalanceWithdrawHandler(svc GophermartUserBalanceWithdrawService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getUserLoginFromContext(w, r, jwt.GetJWTPayload)
		if login == nil {
			return
		}
		var req services.GophermartUserBalanceWithdrawRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		status, err := svc.Withdraw(r.Context(), &req, *login)
		if err != nil {
			handleError(w, err)
			return
		}
		sendTextPlainResponse(w, status)
	}
}
