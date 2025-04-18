package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func dummyGophermartHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}
}

func dummyGophermartMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func TestRegisterUserRegisterRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterUserRegisterRoute(
		router,
		dummyGophermartHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyGophermartMiddleware},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/user/register", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRegisterUserLoginRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterUserLoginRoute(
		router,
		dummyGophermartHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyGophermartMiddleware},
	)
	req := httptest.NewRequest(http.MethodPost, "/api/user/login", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRegisterUserOrdersUploadRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterUserOrdersUploadRoute(
		router,
		dummyGophermartHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyGophermartMiddleware},
	)
	req := httptest.NewRequest(http.MethodPost, "/api/user/orders", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRegisterUserOrderListRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterUserOrderListRoute(
		router,
		dummyGophermartHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyGophermartMiddleware},
	)
	req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRegisterUserBalanceRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterUserBalanceRoute(
		router,
		dummyGophermartHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyGophermartMiddleware},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRegisterUserBalanceWithdrawRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterUserBalanceWithdrawRoute(
		router,
		dummyGophermartHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyGophermartMiddleware},
	)

	req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRegisterUserWithdrawalsRoute(t *testing.T) {
	router := chi.NewRouter()
	RegisterUserWithdrawalsRoute(
		router,
		dummyGophermartHandler(http.StatusOK),
		[]func(http.Handler) http.Handler{dummyGophermartMiddleware},
	)

	req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}
