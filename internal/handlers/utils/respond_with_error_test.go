package utils

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRespondWithError(t *testing.T) {
	rec := httptest.NewRecorder()
	testErr := errors.New("something went wrong")
	status := http.StatusInternalServerError

	RespondWithError(rec, testErr, status)

	res := rec.Result()
	defer res.Body.Close()

	// Проверяем статус код
	require.Equal(t, status, res.StatusCode)

	// Проверяем тело ответа
	expectedBody := testErr.Error() + "\n"
	body := rec.Body.String()
	assert.Equal(t, expectedBody, body)
}
