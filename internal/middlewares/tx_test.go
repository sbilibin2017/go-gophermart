package middlewares

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"
	"time"

	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestTxMiddleware(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "user",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithPollInterval(2 * time.Second), // Увеличиваем интервал опроса
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer container.Terminate(ctx)

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	var db *sql.DB
	for i := 0; i < 5; i++ {
		db, err = sql.Open("pgx", fmt.Sprintf("postgres://user:password@%s:%s/testdb?sslmode=disable", host, port.Port()))
		if err == nil {
			break
		}
		t.Logf("Retrying connection to database (attempt %d)...", i+1)
		time.Sleep(2 * time.Second)
	}

	require.NoError(t, err, "failed to connect to database after 5 attempts")
	defer db.Close()

	txMiddleware := TxMiddleware(db)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, ok := TxFromContext(r.Context())
		if !ok {
			http.Error(w, "no transaction found", http.StatusInternalServerError)
			return
		}

		if r.URL.Query().Get("fail") == "true" {
			tx.Rollback()
			http.Error(w, "transaction rolled back", http.StatusInternalServerError)
			return
		}

		err := tx.Commit()
		if err != nil {
			http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	middlewareHandler := txMiddleware(handler)

	t.Run("Success", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost/?fail=false", nil)
		require.NoError(t, err)
		w := httptest.NewRecorder()

		middlewareHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Rollback", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "http://localhost/?fail=true", nil)
		require.NoError(t, err)
		w := httptest.NewRecorder()

		middlewareHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "transaction rolled back")
	})

}

func TestTxMiddleware_BeginError(t *testing.T) {	
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	
	mock.ExpectBegin().WillReturnError(fmt.Errorf("failed to begin transaction"))

	txMiddleware := TxMiddleware(db)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "failed to begin transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Commit() 
	})

	
	middlewareHandler := txMiddleware(handler)

	req, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	require.NoError(t, err)
	w := httptest.NewRecorder()

	middlewareHandler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to begin transaction")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
