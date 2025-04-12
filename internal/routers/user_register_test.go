package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegisterGophermartRouter(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("registered"))
	}
	mockMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	r := chi.NewRouter()
	RegisterUserRegisterRoute(r, "/api/user", mockHandler, mockMiddleware, mockMiddleware)
	req := httptest.NewRequest(http.MethodPost, "/api/user/register", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Equal(t, "registered", rec.Body.String())
}
