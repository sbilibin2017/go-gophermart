package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserService interface {
	Register(ctx context.Context, u *types.User) (*types.Token, error)
}

func UserRegisterHandler(
	config *configs.GophermartConfig, svc UserService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user types.User
		if !utils.DecodeJSON(w, r, &user) {
			return
		}
		token, err := svc.Register(r.Context(), &user)
		if err != nil {
			switch err {
			case errors.ErrInvalidLogin:
				utils.RespondBadRequest(w, err)
				return
			case errors.ErrUserAlreadyExists:
				utils.RespondConflict(w, err)
				return
			default:
				utils.RespondInternalServerError(w, err)
				return
			}
		}
		utils.SetJSONResponseHeader(w)
		utils.SetStatusOKResponseHeader(w)
		if !utils.EncodeJSON(w, token) {
			return
		}
	}
}
