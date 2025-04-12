package middlewares

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sbilibin2017/go-gophermart/internal/engines/jwt"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockParser := NewMockJWTPArser(ctrl)
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, ok := GetLoginFromContext(r.Context())
		if !ok {
			http.Error(w, "no login in context", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, " + login))
	})
	middleware := AuthMiddleware(mockParser)
	handler := middleware(finalHandler)

	t.Run("missing auth header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), ErrAuthHeaderMissing.Error())
	})

	t.Run("invalid auth header format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidTokenHere")
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), ErrInvalidAuthHeaderFormat.Error())
	})

	t.Run("parser returns error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalidtoken")
		rr := httptest.NewRecorder()
		mockParser.EXPECT().Parse("invalidtoken").Return(nil, errors.New("invalid token"))
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "invalid token")
	})

	t.Run("valid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer validtoken")
		rr := httptest.NewRecorder()
		mockParser.EXPECT().Parse("validtoken").Return(&jwt.Claims{Login: "testuser"}, nil)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Hello, testuser")
	})
}

func TestGetLoginFromContext(t *testing.T) {
	t.Run("no value in context", func(t *testing.T) {
		ctx := context.Background()
		login, ok := GetLoginFromContext(ctx)
		assert.False(t, ok)
		assert.Empty(t, login)
	})

	t.Run("wrong type in context", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), loginKey, "not-claims")
		login, ok := GetLoginFromContext(ctx)
		assert.False(t, ok)
		assert.Empty(t, login)
	})

	t.Run("nil claims in context", func(t *testing.T) {
		var nilClaims *jwt.Claims = nil
		ctx := context.WithValue(context.Background(), loginKey, nilClaims)
		login, ok := GetLoginFromContext(ctx)
		assert.False(t, ok)
		assert.Empty(t, login)
	})

	t.Run("valid claims in context", func(t *testing.T) {
		claims := &jwt.Claims{Login: "validuser"}
		ctx := context.WithValue(context.Background(), loginKey, claims)
		login, ok := GetLoginFromContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, "validuser", login)
	})
}
