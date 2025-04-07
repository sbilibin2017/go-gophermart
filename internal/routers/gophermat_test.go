package routers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/stretchr/testify/assert"
)

func TestNewGophermartRouter_RegisterRoute(t *testing.T) {
	handlerCalled := false
	mockRegisterHandler := func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"status":"registered"}`))
	}
	router := routers.NewGophermartRouter(mockRegisterHandler)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.True(t, handlerCalled, "Handler should be called")
	assert.Contains(t, rec.Header().Get("Content-Encoding"), "gzip", "Gzip middleware should compress the response")
}
