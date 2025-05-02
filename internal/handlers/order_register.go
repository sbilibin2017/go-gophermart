package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderRegisterService interface {
	Register(ctx context.Context, req *types.OrderRequest) (*types.APIStatus, *types.APIStatus)
}

func OrderRegisterHandler(svc OrderRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderRequest

		err := decodeRequestBody(w, r, &req)
		if err != nil {
			return
		}

		successStatus, errorStatus := svc.Register(r.Context(), &req)
		if errorStatus != nil {
			sendTextPlainResponse(w, errorStatus.Message, errorStatus.StatusCode)
			return
		}

		sendTextPlainResponse(w, successStatus.Message, successStatus.StatusCode)
	}
}
