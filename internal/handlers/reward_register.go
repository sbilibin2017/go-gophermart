package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardRegisterService interface {
	Register(ctx context.Context, req *types.RewardRequest) (*types.APISuccessStatus, *types.APIErrorStatus)
}

func RewardRegisterHandler(svc RewardRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RewardRequest
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
