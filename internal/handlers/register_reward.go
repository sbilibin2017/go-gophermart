package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/handlers/helpers"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RegisterRewardService interface {
	Register(ctx context.Context, reward *types.RegisterRewardRequest) (*string, *types.APIError)
}

func RegisterRewardHandler(svc RegisterRewardService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRewardRequest
		if err := helpers.DecodeRequestBody(w, r, &req); err != nil {
			return
		}
		msg, err := svc.Register(r.Context(), &req)
		if err != nil {
			http.Error(w, err.Message, err.Status)
			return
		}
		helpers.SendTextResponse(w, *msg, http.StatusOK)

	}
}
