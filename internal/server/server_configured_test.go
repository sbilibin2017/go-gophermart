package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfigured(t *testing.T) {
	addr := ":8080"
	serverWithRouter := NewServerConfigured(addr)
	assert.NotNil(t, serverWithRouter, "ServerWithRouter should not be nil")

}
