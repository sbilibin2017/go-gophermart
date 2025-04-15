package server

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	addr := ":8080"
	srv := NewServer(addr)
	if srv == nil {
		t.Fatal("Expected server to be non-nil")
	}
	if srv.Addr != addr {
		t.Errorf("Expected server address %s, but got %s", addr, srv.Addr)
	}
}
