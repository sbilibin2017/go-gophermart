package routers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/stretchr/testify/assert"
)

func TestNewAccrualRouter_RoutesRegistered(t *testing.T) {
	router := routers.NewAccrualRouter(nil, nil, nil, nil)

	tests := []struct {
		method string
		path   string
	}{
		{method: http.MethodOptions, path: "/goods"},
		{method: http.MethodOptions, path: "/orders"},
		{method: http.MethodOptions, path: "/orders/12345"},
	}

	for _, tt := range tests {
		t.Run(tt.method+" "+tt.path, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)
			assert.NotEqualf(t, http.StatusNotFound, rr.Code, "Route %s %s should be registered", tt.method, tt.path)
		})
	}
}
