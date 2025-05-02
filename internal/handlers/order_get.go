package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderGetService interface {
	Get(ctx context.Context, req *types.OrderGetRequest) (*types.OrderGetResponse, *types.APIStatus, *types.APIStatus)
}

func OrderGetHandler(svc OrderGetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := types.OrderGetRequest{
			Number: getURLParam(r, "number"),
		}

		resp, successStatus, errorStatus := svc.Get(r.Context(), &req)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		encodeResponseBody(w, resp, successStatus.StatusCode)
	}
}
