package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardService interface {
	Register(ctx context.Context, req *types.RewardRegisterRequest) (*types.APIStatus, error)
}

func RewardHandler(svc RewardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RewardRegisterRequest
		if err := decodeRequest(w, r, &req); err != nil {
			return
		}
		resp, err := svc.Register(r.Context(), &req)
		if err != nil {
			http.Error(w, resp.Message, resp.Status)
			return
		}
		writeTextPlainResponse(w, resp.Status, resp.Message)
	}
}
