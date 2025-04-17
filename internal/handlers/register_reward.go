package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
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
		var reward domain.Reward
		if err := json.NewDecoder(r.Body).Decode(&reward); err != nil {
			http.Error(w, "Unprocessable JSON", http.StatusBadRequest)
			return
		}

		if err := val.Struct(&reward); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		err := svc.Register(r.Context(), &reward)

		if err != nil {
			switch err {
			case domain.ErrRewardKeyAlreadyRegistered:
				http.Error(w, "Reward key already registered", http.StatusConflict)
				return
			default:
				http.Error(w, "Reward is not registered", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reward is registered successfully"))

	}
}
