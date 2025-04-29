package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type UserRegisterService interface {
	Register(ctx context.Context, user *domain.User) (*string, error)
}

type UserRegisterValidator interface {
	Struct(v any) error
}

func UserRegisterHandler(
	val UserRegisterValidator,
	svc UserRegisterService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData struct {
			Login    string `json:"login" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		err := decodeRequestBody(r, &requestData)
		if err != nil {
			handleBadRequestErrorResponse(w)
			return
		}

		err = val.Struct(requestData)
		if err != nil {
			handleValidationErrorResponse(w, err)
			return
		}
		user := &domain.User{
			Login:    requestData.Login,
			Password: requestData.Password,
		}

		token, err := svc.Register(r.Context(), user)
		if err != nil {
			switch err {
			case domain.ErrLoginAlreadyTaken:
				handleErrorResponse(w, err, http.StatusConflict)
				return
			default:
				handleInternalErrorResponse(w)
				return
			}
		}

		setTokenHeader(w, *token)
		sendTextResponse(w, "User successfully registered", http.StatusOK)
	}
}
