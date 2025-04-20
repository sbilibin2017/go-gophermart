package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorInternalServerResponse(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("internal error")

	ErrorInternalServerResponse(w, err)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal error\n", w.Body.String())
}

func TestErrorBadRequestResponse(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("bad request error")

	ErrorBadRequestResponse(w, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Bad request error\n", w.Body.String())
}

func TestErrorConflictResponse(t *testing.T) {
	w := httptest.NewRecorder()
	err := errors.New("conflict happened")

	ErrorConflictResponse(w, err)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Equal(t, "Conflict happened\n", w.Body.String())
}
