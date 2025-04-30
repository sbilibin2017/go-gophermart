package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserOrderUploadService interface {
	Upload(ctx context.Context, order *types.UserOrder) error
}

type UserOrderUploadValidator interface {
	Struct(v any) error
}

type UserOrderUploadValidationErrorRegistry interface {
	Get(err error) *types.ValidationWithStatusCode
}

type UserOrderUploadHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type UserOrderUploadHandler struct {
	val             UserOrderUploadValidator
	valErrRegistry  UserOrderUploadValidationErrorRegistry
	httpErrRegistry UserOrderUploadHTTPErrorRegistry
	svc             UserOrderUploadService
}

func (h *UserOrderUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	login, err := contextutils.GetLogin(r.Context())
	if err != nil {
		sendTextPlainResponse(w, http.StatusUnauthorized, types.ErrUnauthorized.Error())
		return
	}

	var req struct {
		Number string `json:"number" validate:"required,luhn"`
	}
	req.Number = getURLParam(r, "number")

	if handled := handleValidationError(w, req, h.val, h.valErrRegistry); handled {
		return
	}

	order := &types.UserOrder{
		Number: req.Number,
		Login:  login,
	}

	if handled := handleServiceError(w, h.svc.Upload(r.Context(), order), h.httpErrRegistry); handled {
		return
	}

	sendTextPlainResponse(w, http.StatusAccepted, "Order number successfully uploaded")
}
