package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGophermartRouter(t *testing.T) {
	authMiddleware := func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
	}
	registerHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	})
	loginHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	uploadOrderHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	getOrdersHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	getBalanceHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	withdrawHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	getWithdrawalsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	r := NewGophermartRouter(
		authMiddleware,
		registerHandler,
		loginHandler,
		uploadOrderHandler,
		getOrdersHandler,
		getBalanceHandler,
		withdrawHandler,
		getWithdrawalsHandler,
	)

	tests := []struct {
		name           string
		method         string
		url            string
		token          string
		expectedStatus int
	}{
		{
			name:           "Register route",
			method:         http.MethodPost,
			url:            "/api/user/register",
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Login route",
			method:         http.MethodPost,
			url:            "/api/user/login",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Upload Order route with valid token",
			method:         http.MethodPost,
			url:            "/api/user/orders",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Orders route with valid token",
			method:         http.MethodGet,
			url:            "/api/user/orders",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Balance route with valid token",
			method:         http.MethodGet,
			url:            "/api/user/balance",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Withdraw route with valid token",
			method:         http.MethodPost,
			url:            "/api/user/balance/withdraw",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Get Withdrawals route with valid token",
			method:         http.MethodGet,
			url:            "/api/user/withdrawals",
			token:          "valid-token",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Unauthorized for missing token on uploadOrderHandler",
			method:         http.MethodPost,
			url:            "/api/user/orders",
			expectedStatus: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			require.NoError(t, err)
			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
