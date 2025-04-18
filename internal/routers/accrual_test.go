package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func testAccrualHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}
}

func dummyAccrualMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func TestRegisterGetOrderByNumberRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterGetOrderByNumberRoute(router,
		testAccrualHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyAccrualMiddleware},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/orders/12345", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRegisterOrdersRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterOrdersRoute(router,
		testAccrualHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyAccrualMiddleware},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/orders", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRegisterGoodsRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterGoodsRoute(router,
		testAccrualHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyAccrualMiddleware},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/goods", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
