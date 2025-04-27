package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GophermartUserBalanceWithdrawService interface {
	Withdraw(ctx context.Context, req *types.GophermartUserBalanceWithdrawRequest, login string) (*types.APIStatus, *types.APIStatus)
}

func GophermartUserBalanceWithdrawHandler(svc GophermartUserBalanceWithdrawService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login := getUserLoginFromContext(w, r)
		if login == nil {
			return
		}
		var req types.GophermartUserBalanceWithdrawRequest
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
