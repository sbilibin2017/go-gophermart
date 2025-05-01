package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawalsService interface {
	ListWithdrawals(ctx context.Context, login string) ([]*types.UserWithdrawResponse, *types.APISuccessStatus, *types.APIErrorStatus)
}

func UserBalanceWithdrawalsHandler(svc UserBalanceWithdrawalsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := contextutils.GetClaims(r.Context())
		if err != nil {
			sendTextPlainResponse(w, types.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		withdrawals, successStatus, errorStatus := svc.ListWithdrawals(r.Context(), claims.Login)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		encodeResponseBody(w, withdrawals, successStatus.StatusCode)
	}
}
