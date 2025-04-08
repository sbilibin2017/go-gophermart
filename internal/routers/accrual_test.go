package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockAccrualHandler(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(name))
	}
}

func TestNewAccrualRouter(t *testing.T) {
	router := NewAccrualRouter(
		mockAccrualHandler("getOrderAccrual"),
		mockAccrualHandler("registerOrder"),
		mockAccrualHandler("registerGoods"),
	)

	tests := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{"GET", "/api/orders/123456", http.StatusOK},
		{"POST", "/api/orders", http.StatusOK},
		{"POST", "/api/goods", http.StatusOK},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, tt.expectedCode, rr.Code, "unexpected status for %s %s", tt.method, tt.path)
	}
}
