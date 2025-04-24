package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardRegisterService interface {
	Register(ctx context.Context, req *types.RewardRegisterRequest) (*types.APIStatus, *types.APIStatus)
}

func RewardRegisterHandler(svc RewardRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc == nil {
			return
		}
		var req types.RewardRegisterRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}
		resp, err := svc.Register(r.Context(), &req)
		if err != nil {
			handleErrorResponse(w, err.Message, err.Status)
			return
		}
		writeTextPlainResponse(w, resp.Message, resp.Status)
	}
}
