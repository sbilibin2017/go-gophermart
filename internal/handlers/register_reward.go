package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/dto"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

type RegisterRewardUsecase interface {
	Execute(ctx context.Context, req *dto.RegisterRewardRequest) (*dto.RegisterRewardResponse, error)
}

func RegisterRewardHandler(
	uc RegisterRewardUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Info("Начало обработки запроса для регистрации награды")

		var req dto.RegisterRewardRequest
		err := decodeJSON(r, &req)
		if err != nil {
			logger.Logger.Info("Ошибка декодирования JSON в запросе:", err)
			handleRegisterRewardError(w, err)
			return
		}
		logger.Logger.Info("JSON успешно декодирован в запросе")

		resp, err := uc.Execute(r.Context(), &req)
		if err != nil {
			logger.Logger.Info("Ошибка при выполнении бизнес-логики регистрации награды:", err)
			handleRegisterRewardError(w, err)
			return
		}

		logger.Logger.Info("Запрос успешно обработан, отправка ответа с сообщением:", resp.Message)
		writeTextResponse(w, resp.Message)
	}
}

func handleRegisterRewardError(w http.ResponseWriter, err error) {
	logger.Logger.Info("Обработка ошибки при регистрации награды")

	if _, ok := err.(*validator.InvalidValidationError); ok {
		logger.Logger.Info("Ошибка валидации данных:", err)
		http.Error(w, dto.GetRegisterRewardRequestError(err).Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, services.ErrGoodRewardAlreadyExists) {
		logger.Logger.Info("Награда для товара уже существует:", err)
		http.Error(w, capitalizeError(err).Error(), http.StatusConflict)
		return
	} else if errors.Is(err, ErrInternalServer) {
		logger.Logger.Info("Внутренняя ошибка сервера:", err)
		http.Error(w, capitalizeError(err).Error(), http.StatusInternalServerError)
		return
	} else {
		logger.Logger.Info("Неизвестная ошибка:", err)
		http.Error(w, capitalizeError(ErrInternalServer).Error(), http.StatusInternalServerError)
	}
}
