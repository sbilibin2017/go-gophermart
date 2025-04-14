package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type OrderRegisterUsecase interface {
	Execute(ctx context.Context, req *types.OrderRegisterRequest) (*types.OrderRegisterResponse, error)
}

type OrderRegisterRequestDecoder interface {
	Decode(w http.ResponseWriter, r *http.Request, v any) error
}

func OrderRegisterHandler(
	uc OrderRegisterUsecase,
	rd OrderRegisterRequestDecoder,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderRegisterRequest
		if err := rd.Decode(w, r, &req); err != nil {
			handleOrderRegisterError(w, err)
			return
		}
		resp, err := uc.Execute(r.Context(), &req)
		if err != nil {
			handleOrderRegisterError(w, err)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(resp.Message))
	}
}

func handleOrderRegisterError(w http.ResponseWriter, err error) {
	switch err {
	case errors.ErrOrderRegisterInvalidRequest:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case errors.ErrOrderAlreadyRegistered:
		http.Error(w, err.Error(), http.StatusConflict)
		return
	case errors.ErrRequestDecoderUnprocessableJSON: // Here we match the error from the decoder
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
