package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawRequest struct {
	Order string `json:"order" validate:"required"`
	Sum   int64  `json:"sum" validate:"required"`
}

type UserBalanceWithdrawService interface {
	Withdraw(ctx context.Context, req *UserBalanceWithdrawRequest, login string) (*types.APIStatus, *types.APIStatus)
}

func UserBalanceWithdrawHandler(svc UserBalanceWithdrawService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := middlewares.GetLogin(r)
		if err != nil {
			handleErrorResponse(w, "User not authenticated", http.StatusUnauthorized)
			return
		}
		var req UserBalanceWithdrawRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		status, errStatus := svc.Withdraw(r.Context(), &req, login)
		if errStatus != nil {
			handleErrorResponse(w, errStatus.Message, errStatus.Status)
			return
		}
		writeTextPlainResponse(w, status.Message, status.Status)
	}
}
