package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type GoodRewardRegisterService interface {
	Register(ctx context.Context, goodReward *types.GoodReward) error
}

type GoodRewardRegisterValidator interface {
	Struct(v any) error
}

type GoodRewardRegisterValidationErrorRegistry interface {
	Get(err error) *types.ValidationWithStatusCode
}

type GoodRewardRegisterHTTPErrorRegistry interface {
	Get(err error) *types.HTTPError
}

type GoodRewardRegisterHandler struct {
	svc             GoodRewardRegisterService
	val             GoodRewardRegisterValidator
	valErrRegistry  GoodRewardRegisterValidationErrorRegistry
	httpErrRegistry GoodRewardRegisterHTTPErrorRegistry
}

func (h *GoodRewardRegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req types.GoodReward

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

	sendTextPlainResponse(w, http.StatusOK, "Reward successfully registered")
}
