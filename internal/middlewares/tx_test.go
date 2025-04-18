package middlewares

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTxMiddleware_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectCommit()
	txFactory := func(db *sql.DB, op func(tx *sql.Tx) error) error {
		tx, err := db.Begin()
		if err != nil {
			log.Println("Error starting transaction:", err)
			return err
		}
		if err := op(tx); err != nil {
			log.Println("Error during transaction:", err)
			tx.Rollback()
			return err
		}
		log.Println("Committing transaction")
		return tx.Commit()
	}
	req := httptest.NewRequest(http.MethodGet, "/some-endpoint", nil)
	w := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})
	middleware := TxMiddleware(db, txFactory)
	middleware(next).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTxMiddleware_ErrorInTxFactory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()
	mock.ExpectBegin().WillReturnError(errors.New("begin transaction error"))
	txFactory := func(db *sql.DB, op func(tx *sql.Tx) error) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		return op(tx)
	}
	req := httptest.NewRequest(http.MethodGet, "/some-endpoint", nil)
	w := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})
	middleware := TxMiddleware(db, txFactory)
	middleware(next).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "begin transaction error")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestTxMiddleware_ErrorInCommit(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a mock database connection", err)
	}
	defer db.Close()
	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(errors.New("commit error"))
	txFactory := func(db *sql.DB, op func(tx *sql.Tx) error) error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if err := op(tx); err != nil {
			tx.Rollback()
			return err
		}
		return tx.Commit()
	}
	req := httptest.NewRequest(http.MethodGet, "/some-endpoint", nil)
	w := httptest.NewRecorder()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("success"))
	})
	middleware := TxMiddleware(db, txFactory)
	middleware(next).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "commit error")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
