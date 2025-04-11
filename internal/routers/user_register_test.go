package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/stretchr/testify/assert"
)

func TestRegisterUserRegisterHandler(t *testing.T) {
	r := chi.NewRouter()
	config := &configs.GophermartConfig{}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User registered"))
	}
	RegisterUserRegisterHandler(r, config, handler)
	req, err := http.NewRequest(http.MethodPost, "/api/user/register", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "User registered", rr.Body.String())
}
