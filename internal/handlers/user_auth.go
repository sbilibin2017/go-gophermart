package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/jwt"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserAuthService interface {
	Authenticate(ctx context.Context, user *types.User) error
}

type UserAuthValidator interface {
	Struct(v any) error
}

type UserAuthValidationErrorRegistry interface {
	Get(err error) *types.ValidationWithStatusCode
}

type UserAuthHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type UserLoginHandler struct {
	jwtConfig       *configs.JWTConfig
	val             UserAuthValidator
	valErrRegistry  UserAuthValidationErrorRegistry
	httpErrRegistry UserAuthHTTPErrorRegistry
	svc             UserAuthService
}

func (h *UserLoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		sendPlainTextResponse(w, http.StatusBadRequest, capitalize(types.ErrJSONDecode.Error()))
		return
	}

	if err := h.val.Struct(req); err != nil {
		valErr := h.valErrRegistry.Get(err)
		if valErr != nil {
			sendPlainTextResponse(w, valErr.StatusCode, capitalize(valErr.Error.Error()))
			return
		}
	}

	var user = &types.User{
		Login:    req.Login,
		Password: req.Password,
	}

	err := h.svc.Authenticate(r.Context(), user)
	if err != nil {
		errHTTP := h.httpErrRegistry.Get(err)
		if errHTTP != nil {
			sendPlainTextResponse(w, errHTTP.StatusCode, capitalize(errHTTP.Error.Error()))
			return
		}
	}

	tokenString, err := jwt.GenerateTokenString(h.jwtConfig, user)
	if err != nil {
		sendPlainTextResponse(w, http.StatusInternalServerError, capitalize(types.ErrTokenEncode.Error()))
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	sendPlainTextResponse(w, http.StatusOK, "User successfully authenticated")
}
