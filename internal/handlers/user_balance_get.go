package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceGetService interface {
	Get(ctx context.Context, login string) (*types.UserBalanceResponse, *types.APIStatus)
}

func UserBalanceHandler(
	svc UserBalanceGetService,
	loginProvider func(ctx context.Context) (string, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := getLoginFromContext(w, r, loginProvider)
		if err != nil {
			return
		}

		balance, apiErr := svc.Get(r.Context(), login)
		if apiErr != nil {
			sendTextPlainResponse(w, apiErr.Message, apiErr.StatusCode)
			return
		}

		if err := encodeResponseBody(w, balance, http.StatusOK); err != nil {
			sendTextPlainResponse(w, errFailedToEncodeResponseBody, http.StatusInternalServerError)
		}
	}
}
