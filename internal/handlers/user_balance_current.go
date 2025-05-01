package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserBalanceCurrentService interface {
	Get(ctx context.Context, req *types.UserBalanceCurrentRequest) (*types.UserBalanceCurrentResponse, *types.APISuccessStatus, *types.APIErrorStatus)
}

func UserBalanceCurrentHandler(svc UserBalanceCurrentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := contextutils.GetClaims(r.Context())
		if err != nil {
			sendTextPlainResponse(w, types.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return

		}

		req := types.UserBalanceCurrentRequest(claims.Login)

		currentBalance, successStatus, errorStatus := svc.Get(r.Context(), &req)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		encodeResponseBody(w, currentBalance, successStatus.StatusCode)
	}
}
