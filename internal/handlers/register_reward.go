package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type RegisterRewardService interface {
	Register(ctx context.Context, reward *domain.Reward) error
}

type RegisterRewardRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     uint   `json:"reward" validate:"required,gt=0"`
	RewardType string `json:"reward_type" validate:"required,oneof=% pt"`
}

func RegisterRewardHandler(
	svc RegisterRewardService,
) http.HandlerFunc {
	validate := validator.New()

	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRewardRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		err = validate.Struct(req)
		if err != nil {
			http.Error(w, "Ошибка валидации полей запроса", http.StatusBadRequest)
			return
		}

		reward := &domain.Reward{
			Match:      req.Match,
			Reward:     req.Reward,
			RewardType: domain.RewardType(req.RewardType),
		}

		err = svc.Register(r.Context(), reward)
		if err != nil {
			switch err {
			case services.ErrGoodRewardAlreadyExists:
				http.Error(w, "Ключ поиска уже зарегистрирован", http.StatusConflict)
				return
			default:
				http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Информация о вознаграждении за товар зарегистрирована"))
	}
}
