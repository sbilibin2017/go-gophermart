package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestHealthDBHandler_Healthy(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing().WillReturnError(nil)

	req := httptest.NewRequest(http.MethodGet, "/health/db", nil)
	w := httptest.NewRecorder()

	handler := HealthDBHandler(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHealthDBHandler_Unhealthy(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	mock.ExpectPing().WillReturnError(assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/health/db", nil)
	w := httptest.NewRecorder()

	handler := HealthDBHandler(db)
	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
