package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	addr := ":8080"
	srv := NewServer(addr)
	assert.NotNil(t, srv)
	assert.IsType(t, &http.Server{}, srv)
	assert.Equal(t, addr, srv.Addr)
}
