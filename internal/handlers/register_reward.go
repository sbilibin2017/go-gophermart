package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
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
		var req RegisterRewardRequest

		if err := utils.DecodeJSONBody(w, r, &req); err != nil {
			utils.ErrorInternalServerResponse(w, err)
			return
		}

		if err := val.Struct(req); err != nil {
			utils.ErrorBadRequestResponse(w, err)
			return
		}

		err := svc.Register(r.Context(), req.Match, req.Reward, req.RewardType)
		if err != nil {
			switch err {
			case services.ErrRewardAlreadyExists:
				utils.ErrorConflictResponse(w, err)
			case services.ErrRewardIsNotRegistered:
				utils.ErrorBadRequestResponse(w, err)
			default:
				utils.ErrorInternalServerResponse(w, err)
			}
			return
		}

		utils.SendTextResponse(w, http.StatusOK, "Reward registered successfully")
	}
}
