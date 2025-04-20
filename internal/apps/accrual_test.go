package apps_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/apps"
)

func fakeHealthCheck(_ *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("I'm a teapot"))
	}
}

func TestConfigureAccrualApp_Routes(t *testing.T) {
	srv := &http.Server{}

	apps.ConfigureAccrualApp(srv, nil, fakeHealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Errorf("expected status %d, actual status %d", http.StatusTeapot, rec.Code)
	}
	if rec.Body.String() != "I'm a teapot" {
		t.Errorf("expected body 'I'm a teapot', actual body %q", rec.Body.String())
	}
}
