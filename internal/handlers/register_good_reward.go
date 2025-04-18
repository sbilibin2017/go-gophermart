package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
)

type RegisterGoodRewardService interface {
	Register(ctx context.Context, reward *domain.Reward) error
}

type RegisterGoodRewardRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     uint64 `json:"reward" validate:"required,gt=0"`
	RewardType string `json:"reward_type" validate:"required,oneof=% pt"`
}

func RegisterGoodRewardHandler(
	val *validator.Validate,
	svc RegisterGoodRewardService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req RegisterGoodRewardRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, utils.Capitalize(utils.ErrInvalidRequestBody.Error()), http.StatusBadRequest)
			return
		}

		if val != nil {
			if err := val.Struct(req); err != nil {
				http.Error(w, utils.BuildValidationErrorMessage(err), http.StatusBadRequest)
				return
			}
		}

		err := svc.Register(r.Context(), domain.NewReward(
			req.Match,
			req.Reward,
			req.RewardType,
		))

		if err != nil {
			switch err {
			case domain.ErrRewardSearchKeyAlreadyRegistered:
				http.Error(w, utils.Capitalize(err.Error()), http.StatusConflict)
				return
			default:
				http.Error(w, utils.Capitalize(utils.ErrInternal.Error()), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reward registered successfully"))
	}
}
