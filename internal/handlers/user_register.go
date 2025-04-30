package handlers

import (
	"context"
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

	if handled := handleServiceError(w, h.svc.Register(r.Context(), user), h.httpErrRegistry); handled {
		return
	}

	tokenString, err := jwt.GenerateTokenString(h.jwtConfig, user)
	if err != nil {
		sendTextPlainResponse(w, http.StatusBadRequest, types.ErrTokenEncode.Error())
		return
	}

	setAuthorizationHeader(w, tokenString)
	sendTextPlainResponse(w, http.StatusOK, "User successfully registered and authenticated")
}
