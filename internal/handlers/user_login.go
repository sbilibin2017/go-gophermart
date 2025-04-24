package handlers

// import (
// 	"context"
// 	"net/http"

// 	"github.com/sbilibin2017/go-gophermart/internal/services/domain"
// )

// type LoginUserService interface {
// 	Login(ctx context.Context, user *domain.User) (string, error)
// }

// type LoginUserValidator interface {
// 	Struct(v any) error
// }

// type LoginUserRequest struct {
// 	Login    string `json:"login" validate:"required"`
// 	Password string `json:"password" validate:"required"`
// }

// func LoginUserHandler(svc LoginUserService, v LoginUserValidator) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req LoginUserRequest
// 		if err := decodeJSONRequest(w, r, &req, v); err != nil {
// 			handleBadRequestError(w, err)
// 			return
// 		}

// 		user := &domain.User{
// 			Login:    req.Login,
// 			Password: req.Password,
// 		}

// 		token, err := svc.Login(r.Context(), user)
// 		if err != nil {
// 			switch err {
// 			case domain.ErrInvalidCredentials:
// 				handleUnauthorizedError(w, err)
// 				return
// 			default:
// 				handleInternalError(w, err)
// 				return
// 			}
// 		}

// 		setTokenHeader(w, token)
// 		writeTextPlainResponse(w, "User successfully authenticated", http.StatusOK)
// 	}
// }
