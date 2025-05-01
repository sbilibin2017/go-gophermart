package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawService interface {
	Withdraw(ctx context.Context, req *types.UserBalanceWithdrawRequest) (*types.APISuccessStatus, *types.APIErrorStatus)
}

func UserBalanceWithdrawHandler(svc UserBalanceWithdrawService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := getLoginFromContext(w, r)
		if err != nil {
			return
		}

		var req types.UserBalanceWithdrawRequest
		err = decodeRequestBody(w, r, &req)
		if err != nil {
			return
		}

		req.Login = login

		successStatus, errorStatus := svc.Withdraw(r.Context(), &req)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		sendTextPlainResponse(w, successStatus.Message, successStatus.StatusCode)
	}
}
