package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServer_SetHandler(t *testing.T) {
	router := chi.NewRouter()
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	s := &Server{
		Server: &http.Server{},
	}
	s.SetHandler(router)
	require.NotNil(t, s.Handler)
	req, err := http.NewRequest("GET", "/ping", nil)
	require.NoError(t, err)
	rr := httptest.NewRecorder()
	s.Handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "pong", rr.Body.String())
}
