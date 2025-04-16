package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type RegisterRewardService interface {
	Register(ctx context.Context, reward *types.RegisterRewardRequest) error
}

func RegisterRewardHandler(
	svc RegisterRewardService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reward types.RegisterRewardRequest

		err := json.NewDecoder(r.Body).Decode(&reward)
		if err != nil {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		if reward.Match == "" {
			http.Error(w, "Поле match не может быть пустым", http.StatusBadRequest)
			return
		}

		if reward.RewardType == "" {
			http.Error(w, "Поле reward_type не может быть пустым", http.StatusBadRequest)
			return
		}

		if reward.RewardType != "%" && reward.RewardType != "pt" {
			http.Error(w, "Некорректный тип вознаграждения", http.StatusBadRequest)
			return
		}

		err = svc.Register(r.Context(), &reward)

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
