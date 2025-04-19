package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RegisterRewardSaveService interface {
	Register(ctx context.Context, reward *types.Reward) error
}

type RewardSaveValidator interface {
	Struct(s any) error
}

func RegisterRewardSaveHandler(
	val RewardSaveValidator,
	svc RegisterRewardSaveService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var req *types.Reward

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, ErrInvalidRequestBody.Error(), http.StatusInternalServerError)
			return
		}

		if err := val.Struct(req); err != nil {
			http.Error(w, capitalize(buildValidationError(err).Error()), http.StatusBadRequest)
			return
		}

		err := svc.Register(r.Context(), req)
		if err != nil {
			switch err {
			case services.ErrRewardAlreadyExists:
				http.Error(w, capitalize(err.Error()), http.StatusConflict)
			case services.ErrRewardIsNotRegistered:
				http.Error(w, capitalize(err.Error()), http.StatusBadRequest)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(SuccessRewardRegistered))
	}
}

const (
	SuccessRewardRegistered = "Reward registered successfully"
)
