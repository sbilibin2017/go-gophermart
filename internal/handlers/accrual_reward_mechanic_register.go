package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type AccrualRewardMechanicRegisterService interface {
	Register(ctx context.Context, req *types.AccrualRewardMechanicRegisterRequest) (*types.APIStatus, *types.APIStatus)
}

func AccrualRewardMechanicRegisterHandler(svc AccrualRewardMechanicRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccrualRewardMechanicRegisterRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		resp, err := svc.Register(r.Context(), &req)
		if err != nil {
			handleError(w, err)
			return
		}
		sendTextPlainResponse(w, resp)
	}
}
