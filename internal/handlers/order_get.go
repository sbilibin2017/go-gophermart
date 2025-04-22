package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetService interface {
	Get(ctx context.Context, req *types.OrderGetRequest) (*types.APIResponse[types.OrderGetResponse], *types.APIError)
}

func OrderGetHandler(svc OrderGetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		number := getURLParam(r, "name")

		resp, err := svc.Get(r.Context(), &types.OrderGetRequest{Number: number})
		if err != nil {
			http.Error(w, err.Message, err.Status)
			return
		}

		sendJSONResponse(w, *resp)
	}
}
