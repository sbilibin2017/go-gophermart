package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RewardRegisterService interface {
	Register(ctx context.Context, req *types.RewardRegisterMechanicRequest) (*types.APIStatus, *types.APIStatus)
}

func RewardRegisterHandler(svc RewardRegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RewardRegisterMechanicRequest
		err := decodeRequestBody(w, r, &req)
		if err != nil {
			return
		}

		successStatus, errorStatus := svc.Register(r.Context(), &req)
		if errorStatus != nil {
			w.WriteHeader(errorStatus.StatusCode)
			return
		}

		w.WriteHeader(successStatus.StatusCode)
	}
}
