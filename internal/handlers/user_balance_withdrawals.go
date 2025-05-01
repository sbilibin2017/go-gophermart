package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawalsService interface {
	ListWithdrawals(ctx context.Context, login string) ([]*types.UserWithdrawResponse, *types.APISuccessStatus, *types.APIErrorStatus)
}

func UserBalanceWithdrawalsHandler(svc UserBalanceWithdrawalsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := getLoginFromContext(w, r)
		if err != nil {
			return
		}

		withdrawals, successStatus, errorStatus := svc.ListWithdrawals(r.Context(), login)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		encodeResponseBody(w, withdrawals, successStatus.StatusCode)
	}
}
