package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserRegisterService interface {
	Register(ctx context.Context, user *types.User) error
}

type UserRegisterValidator interface {
	Struct(v any) error
}

type UserRegisterValidationErrorRegistry interface {
	Get(err error) *types.ValidationWithStatusCode
}

type UserRegisterHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type UserRegisterHandler struct {
	jwtConfig       *configs.JWTConfig
	val             UserRegisterValidator
	valErrRegistry  UserRegisterValidationErrorRegistry
	httpErrRegistry UserRegisterHTTPErrorRegistry
	svc             UserRegisterService
}

func (h *UserRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// Декодируем тело запроса
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		// Используем нашу функцию для отправки ошибки
		sendPlainTextResponse(w, http.StatusBadRequest, capitalize(types.ErrJSONDecode.Error()))
		return
	}

	// Проверка на валидность данных
	if err := h.val.Struct(req); err != nil {
		valErr := h.valErrRegistry.Get(err)
		if valErr != nil {
			// Используем нашу функцию для отправки ошибки
			sendPlainTextResponse(w, valErr.StatusCode, capitalize(valErr.Error.Error()))
			return
		}
	}

	// Создаем объект пользователя
	var user = &types.User{
		Login:    req.Login,
		Password: req.Password,
	}

	// Попытка регистрации
	err := h.svc.Register(r.Context(), user)
	errHTTP := h.httpErrRegistry.Get(err)
	if errHTTP != nil {
		// Используем нашу функцию для отправки ошибки
		sendPlainTextResponse(w, errHTTP.StatusCode, capitalize(errHTTP.Error.Error()))
		return
	}

	// Генерация токена
	tokenString, err := jwt.GenerateTokenString(h.jwtConfig, user)
	if err != nil {
		// Используем нашу функцию для отправки ошибки
		sendPlainTextResponse(w, http.StatusBadRequest, capitalize(types.ErrTokenEncode.Error()))
		return
	}

	// Отправка успешного ответа с токеном
	w.Header().Set("Authorization", "Bearer "+tokenString)
	sendPlainTextResponse(w, http.StatusOK, "User successfully registered and authenticated")
}
