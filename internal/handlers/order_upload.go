package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/domain"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/validation"
)

type OrderUploadService interface {
	Upload(ctx context.Context, order *domain.Order, login string) error
}

type OrderUploadValidator interface {
	Struct(v any) error
}

func OrderUploadHandler(
	val OrderUploadValidator,
	svc OrderUploadService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		login, ok := middlewares.GetLoginFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized) // 401
			return
		}

		number := chi.URLParam(r, "number")

		var requestData struct {
			Number string `json:"number" validate:"required,luhn"`
		}

		requestData.Number = number

		err := val.Struct(requestData)
		if err != nil {
			err = validation.HandleInvalidLuhnNumber(err)
			if errors.Is(err, validation.ErrInvalidLuhnNumber) {
				http.Error(w, "Invalid Luhn number", http.StatusUnprocessableEntity)
				return
			}
			http.Error(w, "Validation failed", http.StatusBadRequest)
			return
		}

		order := &domain.Order{
			Number: requestData.Number,
		}

		err = svc.Upload(r.Context(), order, login)
		if err != nil {
			switch err {
			case domain.ErrUserOrderExists:
				http.Error(w, "Order already uploaded by this user", http.StatusOK) // 200
				return
			case domain.ErrOrderExists:
				http.Error(w, "Order already uploaded by another user", http.StatusConflict) // 409
				return
			default:
				http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
				return
			}
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Order accepted successfully"))
	}
}
