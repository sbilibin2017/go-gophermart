package handlers

// import (
// 	"context"
// 	"net/http"

// 	"github.com/sbilibin2017/go-gophermart/internal/services/domain"
// )

// type RegisterUserService interface {
// 	Register(ctx context.Context, user *domain.User) (string, error)
// }

// type RegisterUserValidator interface {
// 	Struct(v any) error
// }

// type RegisterUserRequest struct {
// 	Login    string `json:"login" validate:"required"`
// 	Password string `json:"password" validate:"required"`
// }

// func RegisterUserHandler(svc RegisterUserService, v RegisterUserValidator) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req RegisterUserRequest
// 		if err := decodeJSONRequest(w, r, &req, v); err != nil {
// 			return
// 		}

// 		user := &domain.User{
// 			Login:    req.Login,
// 			Password: req.Password,
// 		}

// 		token, err := svc.Register(r.Context(), user)
// 		if err != nil {
// 			switch err {
// 			case domain.ErrLoginAlreadyExists:
// 				handleConflictError(w, err)
// 				return
// 			default:
// 				handleInternalError(w, err)
// 				return
// 			}
// 		}

// 		setTokenHeader(w, token)
// 		writeTextPlainResponse(w, "User successfully registered and authenticated", http.StatusOK)
// 	}
// }
