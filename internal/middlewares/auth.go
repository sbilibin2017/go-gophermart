package middlewares

import (
	"context"
	"net/http"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

type UserRepository interface {
	GetByID(ctx context.Context, id *types.UserID) (*types.User, error)
}

type JWTDecoder interface {
	Decode(tokenString string) (*types.Claims, error)
}

func AuthMiddleware(repo UserRepository, decoder JWTDecoder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
				return
			}
			claims, err := decoder.Decode(tokenString)
			if err != nil {
				http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
				return
			}
			user, err := repo.GetByID(r.Context(), claims.UserID)
			if err != nil || user == nil {
				http.Error(w, "Unauthorized: User not found", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
