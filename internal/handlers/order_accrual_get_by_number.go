package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderAccrualGetByNumberService interface {
	GetByNumber(ctx context.Context, number string) (*types.OrderAccrual, error)
}

type OrderAccrualGetByNumberValidator interface {
	Struct(v any) error
}

type OrderAccrualGetByNumberValidationErrorRegistry interface {
	Get(err error) *types.ValidationWithStatusCode
}

type OrderAccrualGetByNumberHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type OrderAccrualGetByNumberHandler struct {
	svc             OrderAccrualGetByNumberService
	val             OrderAccrualGetByNumberValidator
	valErrRegistry  OrderAccrualGetByNumberValidationErrorRegistry
	httpErrRegistry OrderAccrualGetByNumberHTTPErrorRegistry
}

func (h *OrderAccrualGetByNumberHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Number string `json:"number" validate:"required,luhn"`
	}
	req.Number = getURLParam(r, "number")

	if handled := handleValidationError(w, req, h.val, h.valErrRegistry); handled {
		return
	}

	resp, err := h.svc.GetByNumber(r.Context(), req.Number)
	if handleServiceError(w, err, h.httpErrRegistry) {
		return
	}

	sendJSONResponse(w, http.StatusOK, resp)
}
