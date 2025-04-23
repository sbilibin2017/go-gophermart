package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GoodRewardService interface {
	Register(ctx context.Context, req *types.GoodRewardRegisterRequest) (*types.APIStatus, *types.APIStatus, error)
}

func GoodRewardHandler(svc GoodRewardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GoodRewardRegisterRequest
		if err := decodeRequest(w, r, &req); err != nil {
			return
		}
		resp, status, err := svc.Register(r.Context(), &req)
		if err != nil {
			http.Error(w, status.Message, status.Status)
			return
		}
		writeTextPlainResponse(w, resp.Status, resp.Message)
	}
}
