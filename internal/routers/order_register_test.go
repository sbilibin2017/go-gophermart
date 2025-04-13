package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestRegisterOrderRegisterRoute(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	gm := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-1", "Executed")
			next.ServeHTTP(w, r)
		})
	}

	lm := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Middleware-2", "Executed")
			next.ServeHTTP(w, r)
		})
	}

	r := chi.NewRouter()
	RegisterOrderRegisterRoute(r, "/api", handler, gm, lm)

	req, err := http.NewRequest("POST", "/api/orders", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	assert.Equal(t, "Executed", rr.Header().Get("X-Middleware-1"))
	assert.Equal(t, "Executed", rr.Header().Get("X-Middleware-2"))
}
