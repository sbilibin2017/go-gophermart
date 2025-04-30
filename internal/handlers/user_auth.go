package handlers

import (
	"context"
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

	if err := decodeJSONRequest(r, &req); err != nil {
		sendTextPlainResponse(w, http.StatusBadRequest, capitalize(types.ErrJSONDecode.Error()))
		return
	}

	if handled := handleValidationError(w, req, h.val, h.valErrRegistry); handled {
		return
	}

	user := &types.User{
		Login:    req.Login,
		Password: req.Password,
	}

	if handled := handleServiceError(w, h.svc.Authenticate(r.Context(), user), h.httpErrRegistry); handled {
		return
	}

	tokenString, err := jwt.GenerateTokenString(h.jwtConfig, user)
	if err != nil {
		sendTextPlainResponse(w, http.StatusInternalServerError, types.ErrTokenEncode.Error())
		return
	}

	setAuthorizationHeader(w, tokenString)
	sendTextPlainResponse(w, http.StatusOK, "User successfully authenticated")
}
