package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccrualRouter(t *testing.T) {
	r := NewAccrualRouter(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		},
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
		},
	)
	tests := []struct {
		name               string
		url                string
		method             string
		expectedStatusCode int
	}{
		{
			name:               "Test GET /api/orders/{number}",
			url:                "/api/orders/123",
			method:             http.MethodGet,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Test POST /api/orders",
			url:                "/api/orders",
			method:             http.MethodPost,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Test POST /api/goods",
			url:                "/api/goods",
			method:             http.MethodPost,
			expectedStatusCode: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
