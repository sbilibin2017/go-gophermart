package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GoodRewardService interface {
	Register(ctx context.Context, reward *types.RewardRegisterRequest) (*types.APIResponse[any], *types.APIError)
}

func GoodRewardHandler(svc GoodRewardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RewardRegisterRequest
		if err := decodeJSONRequest(w, r, &req); err != nil {
			return
		}

		resp, err := svc.Register(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Message, err.Status)
			return
		}

		sendTextResponse(w, resp)
	}
}
