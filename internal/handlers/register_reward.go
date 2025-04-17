package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
)

type RegisterRewardService interface {
	Register(ctx context.Context, reward *domain.Reward) error
}

type RegisterRewardValidator interface {
	Struct(v interface{}) error
}

func RegisterRewardHandler(
	svc RegisterRewardService, val RegisterRewardValidator,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var reward domain.Reward
		err := utils.DecodeJSON(r, &reward)
		if err != nil {
			http.Error(w, utils.Capitalize(err.Error()), http.StatusBadRequest)
			return
		}

		err = val.Struct(&reward)
		if err != nil {
			http.Error(w, utils.Capitalize(utils.ValidationErrorsToString(utils.FormatValidationError(err))), http.StatusBadRequest)
			return
		}

		err = svc.Register(r.Context(), &reward)
		if err != nil {
			switch err {
			case domain.ErrRewardKeyAlreadyRegistered:
				http.Error(w, utils.Capitalize(err.Error()), http.StatusConflict)
				return
			default:
				http.Error(w, utils.Capitalize(domain.ErrFailedToRegisterReward.Error()), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		utils.EncodeJSON(w, map[string]string{"message": "Reward registered successfully"})

	}
}
