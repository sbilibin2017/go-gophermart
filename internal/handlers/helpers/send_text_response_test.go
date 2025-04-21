package helpers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendTextResponse(t *testing.T) {
	w := httptest.NewRecorder()
	status := http.StatusCreated
	message := "Test response message"

	SendTextResponse(w, status, message)

	assert.Equal(t, status, w.Code)
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Equal(t, message, w.Body.String())
}
