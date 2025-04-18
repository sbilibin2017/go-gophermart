package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RegisterRewardSaveService interface {
	Register(ctx context.Context, reward *types.Reward) error
}

func RegisterRewardSaveHandler(
	val *validator.Validate,
	svc RegisterRewardSaveService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var req *types.Reward

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusInternalServerError)
			return
		}

		if val != nil {
			if err := val.Struct(req); err != nil {
				http.Error(w, capitalize(buildValidationError(err).Error()), http.StatusBadRequest)
				return
			}
		}

		err := svc.Register(r.Context(), req)
		if err != nil {
			switch err {
			case services.ErrRewardAlreadyExists:
				http.Error(w, "Reward with the same search key already exists", http.StatusConflict)
			case services.ErrRewardIsNotRegistered:
				http.Error(w, "Reward is not registered", http.StatusBadRequest)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reward registered successfully"))
	}
}
