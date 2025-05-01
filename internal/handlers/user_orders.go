package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrdersService interface {
	ListOrdered(ctx context.Context, req *types.UserOrdersRequest) ([]*types.UserOrdersResponse, *types.APISuccessStatus, *types.APIErrorStatus)
}

func UserOrdersHandler(svc UserOrdersService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, err := getLoginFromContext(w, r)
		if err != nil {
			return
		}

		req := types.UserOrdersRequest(login)

		orders, successStatus, errorStatus := svc.ListOrdered(r.Context(), &req)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		encodeResponseBody(w, orders, successStatus.StatusCode)
	}
}
