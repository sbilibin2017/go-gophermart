package handlers

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/requests"
)

type RewardService interface {
	Register(ctx context.Context, reward *domain.Reward) error
}

func RegisterGoodRewardHandler(
	svc RewardService,
	decode func(w http.ResponseWriter, r *http.Request, v *requests.RewardRequest) error,
	validate func(w http.ResponseWriter, validate *validator.Validate, v interface{}) error,
	respondWithError func(w http.ResponseWriter, err error, status int),
) http.HandlerFunc {
	v := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.RewardRequest

		err := decode(w, r, &req)
		if err != nil {
			respondWithError(w, err, http.StatusBadRequest)
			return
		}

		err = validate(w, v, &req)
		if err != nil {
			respondWithError(w, err, http.StatusBadRequest)
			return
		}

		err = svc.Register(r.Context(), domain.NewReward(
			req.Match, req.Reward, req.RewardType,
		))

		if err != nil {
			switch err {
			case errors.ErrGoodRewardAlreadyExists:
				respondWithError(w, err, http.StatusConflict)
				return
			default:
				respondWithError(w, err, http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Good reward registered successfully"))
	}
}
