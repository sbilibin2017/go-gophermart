package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderAccrualRegisterService interface {
	Register(ctx context.Context, orderAccrual *types.OrderAccrual) error
}

type OrderAccrualRegisterValidator interface {
	Struct(v any) error
}

type OrderAccrualRegisterValidationErrorRegistry interface {
	Get(err error) *types.ValidationWithStatusCode
}

type OrderAccrualRegisterHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type OrderAccrualRegisterHandler struct {
	svc             OrderAccrualRegisterService
	val             OrderAccrualRegisterValidator
	valErrRegistry  OrderAccrualRegisterValidationErrorRegistry
	httpErrRegistry OrderAccrualRegisterHTTPErrorRegistry
}

func (h *OrderAccrualRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.OrderAccrual
	if err := decodeJSONRequest(r, &req); err != nil {
		sendTextPlainResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if handled := handleValidationError(w, req, h.val, h.valErrRegistry); handled {
		return
	}

	err := h.svc.Register(r.Context(), &req)
	if handleServiceError(w, err, h.httpErrRegistry) {
		return
	}

	sendTextPlainResponse(w, http.StatusOK, "Order accrual successfully registered")
}
