package server

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockAddresser := NewMockServerAddresser(ctrl)
	expectedAddress := ":8080"
	mockAddresser.EXPECT().GetRunAddress().Return(expectedAddress)

	s := NewServer(mockAddresser)
	assert.NotNil(t, s)
	assert.Equal(t, expectedAddress, s.Addr)
}

func TestSetHandler(t *testing.T) {
	s := &Server{Server: &http.Server{}}
	router := chi.NewRouter()
	s.SetHandler(router)
	assert.Equal(t, router, s.Handler)
}
