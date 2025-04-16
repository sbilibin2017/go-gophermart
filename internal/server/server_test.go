package server

import (
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestNewServer(t *testing.T) {
	addr := ":8080"
	h := chi.NewRouter()
	srv := NewServer(addr, h)
	if srv == nil {
		t.Fatal("Expected server to be non-nil")
	}
	if srv.Addr != addr {
		t.Errorf("Expected server address %s, but got %s", addr, srv.Addr)
	}
}
