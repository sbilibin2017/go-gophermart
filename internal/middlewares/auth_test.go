package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockJWTConfig(ctrl)
	mockClaims := NewMockClaims(ctrl)
	mockConfig.EXPECT().GetJWTSecretKey().Return("secret-key").AnyTimes()
	mockClaims.EXPECT().GetLogin().Return("user123").AnyTimes()
	decoder := func(tokenStr string, config JWTConfig) (Claims, error) {
		if tokenStr == "valid-token" {
			return mockClaims, nil
		}
		return nil, errors.New("invalid token")
	}
	req := httptest.NewRequest(http.MethodGet, "/some-path", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	authMiddleware := AuthMiddleware(mockConfig, decoder)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login, ok := r.Context().Value(loginKey).(string)
		assert.True(t, ok, "expected login in context")
		assert.Equal(t, "user123", login)
		w.WriteHeader(http.StatusOK)
	})
	handler := authMiddleware(next)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockJWTConfig(ctrl)
	decoder := func(tokenStr string, config JWTConfig) (Claims, error) {
		return nil, errors.New("invalid token")
	}
	req := httptest.NewRequest(http.MethodGet, "/some-path", nil)
	authMiddleware := AuthMiddleware(mockConfig, decoder)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := authMiddleware(next)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAuthMiddleware_InvalidAuthHeaderFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockJWTConfig(ctrl)
	decoder := func(tokenStr string, config JWTConfig) (Claims, error) {
		return nil, errors.New("invalid token")
	}
	req := httptest.NewRequest(http.MethodGet, "/some-path", nil)
	req.Header.Set("Authorization", "InvalidHeaderFormat")
	authMiddleware := AuthMiddleware(mockConfig, decoder)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := authMiddleware(next)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
