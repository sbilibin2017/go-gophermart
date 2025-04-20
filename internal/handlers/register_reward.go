package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type RegisterRewardService interface {
	Register(ctx context.Context, match string, reward uint64, rewardType string) error
}

type RegisterRewardValidator interface {
	Struct(s any) error
}

type RegisterRewardRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     uint64 `json:"reward" validate:"required,gt=0"`
	RewardType string `json:"reward_type" validate:"required,oneof=% pt"`
}

func RegisterRewardHandler(
	val RegisterRewardValidator,
	svc RegisterRewardService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req *RegisterRewardRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusInternalServerError)
			return
		}

		if err := val.Struct(req); err != nil {
			http.Error(w, capitalize(buildValidationError(err).Error()), http.StatusBadRequest)
			return
		}

		err := svc.Register(r.Context(), req.Match, req.Reward, req.RewardType)
		if err != nil {
			switch err {
			case services.ErrRewardAlreadyExists:
				http.Error(w, capitalize(err.Error()), http.StatusConflict)
			case services.ErrRewardIsNotRegistered:
				http.Error(w, capitalize(err.Error()), http.StatusBadRequest)
			}
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reward registered successfully"))
	}
}
