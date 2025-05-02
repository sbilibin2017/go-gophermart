package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceWithdrawListService interface {
	List(ctx context.Context, login string) ([]*types.UserBalanceWithdrawResponse, *types.APIStatus)
}

func UserBalanceWithdrawListHandler(
	svc UserBalanceWithdrawListService,
	loginProvider func(ctx context.Context) (string, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := getLoginFromContext(w, r, loginProvider)
		if err != nil {
			return
		}

		withdrawals, apiErr := svc.List(r.Context(), login)
		if apiErr != nil {
			sendTextPlainResponse(w, apiErr.Message, apiErr.StatusCode)
			return
		}

		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if err := encodeResponseBody(w, withdrawals, http.StatusOK); err != nil {
			sendTextPlainResponse(w, errFailedToEncodeResponseBody, http.StatusInternalServerError)
		}
	}
}
