package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Success(t *testing.T) {
	mockTokenDecoder := func(tokenString string, config *configs.JWTConfig) (*types.Claims, error) {
		return &types.Claims{Login: "testuser"}, nil
	}
	mockClaimsContextSetter := func(ctx context.Context, claims *types.Claims) context.Context {
		return context.WithValue(ctx, "login", claims.Login)
	}

	handler := AuthMiddleware(
		&configs.JWTConfig{SecretKey: "secret"},
		mockTokenDecoder,
		mockClaimsContextSetter,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "testuser", r.Context().Value("login"))
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	mockTokenDecoder := func(tokenString string, config *configs.JWTConfig) (*types.Claims, error) {
		return nil, errors.New("invalid token")
	}
	mockClaimsContextSetter := func(ctx context.Context, claims *types.Claims) context.Context {
		return ctx
	}

	handler := AuthMiddleware(
		&configs.JWTConfig{SecretKey: "secret"},
		mockTokenDecoder,
		mockClaimsContextSetter,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	mockTokenDecoder := func(tokenString string, config *configs.JWTConfig) (*types.Claims, error) {
		return nil, errors.New("missing token")
	}
	mockClaimsContextSetter := func(ctx context.Context, claims *types.Claims) context.Context {
		return ctx
	}

	handler := AuthMiddleware(
		&configs.JWTConfig{SecretKey: "secret"},
		mockTokenDecoder,
		mockClaimsContextSetter,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
