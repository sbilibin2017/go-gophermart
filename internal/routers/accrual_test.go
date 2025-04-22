package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccrualRouter(t *testing.T) {
	goodRewardHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Good Reward"))
	})

	orderAcceptHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Order Accepted"))
	})

	orderGetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Order Details"))
	})

	middlewares := []func(http.Handler) http.Handler{}

	r := NewAccrualRouter(goodRewardHandler, orderAcceptHandler, orderGetHandler, middlewares)

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "POST /goods",
			method:         http.MethodPost,
			url:            "/goods",
			expectedStatus: http.StatusOK,
			expectedBody:   "Good Reward",
		},
		{
			name:           "POST /orders",
			method:         http.MethodPost,
			url:            "/orders",
			expectedStatus: http.StatusCreated,
			expectedBody:   "Order Accepted",
		},
		{
			name:           "GET /orders/{number}",
			method:         http.MethodGet,
			url:            "/orders/12345",
			expectedStatus: http.StatusOK,
			expectedBody:   "Order Details",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())
		})
	}
}
