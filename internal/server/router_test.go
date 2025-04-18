package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	assert.NotNil(t, router, "router must be not nil")
}
