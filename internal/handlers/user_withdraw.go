package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserWithdrawService interface {
	Withdraw(ctx context.Context, withdraw *types.UserWithdrawal) error
}

type UserWithdrawValidator interface {
	Struct(v any) error
}

type UserWithdrawValidationErrorRegistry interface {
	Get(err error) *types.ValidationWithStatusCode
}

type UserWithdrawHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type UserWithdrawHandler struct {
	val             UserWithdrawValidator
	httpErrRegistry UserWithdrawHTTPErrorRegistry
	valErrRegistry  UserWithdrawValidationErrorRegistry
	svc             UserWithdrawService
}

func (h *UserWithdrawHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	login, err := contextutils.GetLogin(r.Context())
	if err != nil {
		sendTextPlainResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Order string `json:"order" validate:"required,luhn"`
		Sum   int64  `json:"sum" validate:"required,min=1"`
	}

	if err := decodeJSONRequest(r, &req); err != nil {
		sendTextPlainResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if handled := handleValidationError(w, req, h.val, h.valErrRegistry); handled {
		return
	}

	withdraw := &types.UserWithdrawal{
		Login:  login,
		Number: req.Order,
		Sum:    req.Sum,
	}

	if handled := handleServiceError(w, h.svc.Withdraw(r.Context(), withdraw), h.httpErrRegistry); handled {
		return
	}

	sendTextPlainResponse(w, http.StatusOK, "Withdraw request successful")
}
