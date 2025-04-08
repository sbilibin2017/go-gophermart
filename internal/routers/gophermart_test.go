package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/configs"

	"github.com/stretchr/testify/assert"
)

func mockGophermartHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(name))
	}
}

func TestNewGophermartRouter(t *testing.T) {
	config := &configs.GophermartConfig{
		JWTSecretKey: "test-secret",
		JWTExp:       time.Hour,
	}

	router := NewGophermartRouter(
		config,
		mockGophermartHandler("register"),
		mockGophermartHandler("login"),
		mockGophermartHandler("uploadOrder"),
		mockGophermartHandler("getOrders"),
		mockGophermartHandler("getBalance"),
		mockGophermartHandler("withdraw"),
		mockGophermartHandler("getWithdrawals"),
	)

	tests := []struct {
		method       string
		path         string
		expectedCode int
		authRequired bool
	}{
		{"POST", "/api/user/register", http.StatusOK, false},
		{"POST", "/api/user/login", http.StatusOK, false},
		{"POST", "/api/user/orders", http.StatusUnauthorized, true},
		{"GET", "/api/user/orders", http.StatusUnauthorized, true},
		{"GET", "/api/user/balance", http.StatusUnauthorized, true},
		{"POST", "/api/user/balance/withdraw", http.StatusUnauthorized, true},
		{"GET", "/api/user/withdrawals", http.StatusUnauthorized, true},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, tt.expectedCode, rr.Code, "unexpected status for %s %s", tt.method, tt.path)
	}
}
