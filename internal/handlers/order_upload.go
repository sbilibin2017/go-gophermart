package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderUploadService interface {
	Upload(ctx context.Context, order *types.Order) error
}

type OrderUploadValidator interface {
	Struct(v any) error
}

type OrderUploadValidationErrorRegistry interface {
	Get(err error) *types.ValidationWithStatusCode
}

type OrderUploadHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type OrderUploadHandler struct {
	valErrRegistry  OrderUploadValidationErrorRegistry
	httpErrRegistry OrderUploadHTTPErrorRegistry
	svc             OrderUploadService
	val             OrderUploadValidator
}

func (h *OrderUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	login, err := contextutils.GetLogin(r.Context())
	if err != nil {
		sendPlainTextResponse(w, http.StatusUnauthorized, capitalize(types.ErrUnauthorized.Error()))
		return
	}

	var req struct {
		Number string `json:"number" validate:"required,luhn"`
	}
	req.Number = getURLParam(r, "number")

	if err := h.val.Struct(req); err != nil {
		valErr := h.valErrRegistry.Get(err)
		if valErr != nil {
			sendPlainTextResponse(w, valErr.StatusCode, capitalize(valErr.Error.Error()))
			return
		}
	}

	var order = &types.Order{
		Number: req.Number,
		Login:  login,
	}

	err = h.svc.Upload(r.Context(), order)

	errHTTP := h.httpErrRegistry.Get(err)
	if errHTTP != nil {
		sendPlainTextResponse(w, errHTTP.StatusCode, capitalize(errHTTP.Error.Error()))
		return
	}

	sendPlainTextResponse(w, http.StatusAccepted, "Order number successfully uploaded")
}
