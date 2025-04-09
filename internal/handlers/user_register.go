package handlers

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
	"github.com/sbilibin2017/go-gophermart/internal/requests"
	"github.com/sbilibin2017/go-gophermart/internal/responses"
)

type UserRegisterUsecase interface {
	Execute(ctx context.Context, req *requests.UserRegisterRequest) (*responses.UserRegisterResponse, error)
}

func UserRegisterHandler(
	config *configs.GophermartConfig, uc UserRegisterUsecase,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req requests.UserRegisterRequest
		if !utils.DecodeJSON(w, r, &req) {
			return
		}
		resp, err := uc.Execute(r.Context(), &req)
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
		if !utils.EncodeJSON(w, resp.AccessToken) {
			return
		}
	}
}
