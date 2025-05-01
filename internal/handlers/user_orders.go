package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrdersService interface {
	ListOrdered(ctx context.Context, req *types.UserOrdersRequest) ([]*types.UserOrdersResponse, *types.APISuccessStatus, *types.APIErrorStatus)
}

func UserOrdersHandler(svc UserOrdersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := contextutils.GetClaims(r.Context())
		if err != nil {
			sendTextPlainResponse(w, types.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return

		}

		req := types.UserOrdersRequest(claims.Login)

		orders, successStatus, errorStatus := svc.ListOrdered(r.Context(), &req)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		encodeResponseBody(w, orders, successStatus.StatusCode)
	}
}
