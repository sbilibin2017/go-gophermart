package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawService interface {
	Withdraw(ctx context.Context, req *types.UserBalanceWithdrawRequest) *types.APIStatus
}

func UserBalanceWithdrawHandler(
	svc UserBalanceWithdrawService,
	loginProvider func(ctx context.Context) (string, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := getLoginFromContext(w, r, loginProvider)
		if err != nil {
			return 
		}

		var req types.UserBalanceWithdrawRequest
		if err := decodeRequestBody(w, r, &req); err != nil {
			return 
		}

		req.Login = login 

		status := svc.Withdraw(r.Context(), &req)
		if status != nil {
			sendTextPlainResponse(w, status.Message, status.StatusCode)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
