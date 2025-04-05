package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewServerWuthRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockConfig(ctrl)
	mockConfig.EXPECT().GetRunAddress().Return(":8080").Times(1)
	server := NewServerWithRouter(mockConfig)
	assert.NotNil(t, server)
	assert.Equal(t, ":8080", server.Addr)
	assert.Equal(t, server.Handler, server.rtr)
}

func TestAddRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockConfig := NewMockConfig(ctrl)
	mockConfig.EXPECT().GetRunAddress().Return(":8080").Times(1)
	server := NewServerWithRouter(mockConfig)
	newRouter := chi.NewRouter()
	newRouter.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Test"))
	})
	server.AddRouter(newRouter)
	assert.NotNil(t, server.rtr)
	assert.Len(t, server.rtr.Routes(), 1)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.NoError(t, err)
	rr := httptest.NewRecorder()
	server.rtr.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Test", rr.Body.String())
}
