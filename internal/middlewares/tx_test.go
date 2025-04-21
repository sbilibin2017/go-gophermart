package middlewares

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx" // Импортируем sqlx
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func createTestDB() (*sqlx.DB, sqlmock.Sqlmock) {
	// Создаем мок с помощью sqlmock
	db, mock, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	return sqlxDB, mock
}

func TestHandleTransactionEnd_TxNil(t *testing.T) {
	// Создаем HTTP-респондер для захвата ответа
	w := httptest.NewRecorder()

	// В данном случае транзакция равна nil
	var tx *sqlx.Tx

	// Вызываем функцию с nil транзакцией
	handleTransactionEnd(tx, http.StatusOK, w)

	// Проверяем, что функция просто вернулась без ошибок, т.к. tx == nil
	// В этом случае не должны быть сделаны ни commit, ни rollback
	// Проверим, что статус ответа остался в исходном состоянии, т.е. не был изменен
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleTransactionEnd_Rollback(t *testing.T) {
	// Создаем мок базы данных и мок
	db, mock := createTestDB()

	// Ожидаем вызов Begin
	mock.ExpectBegin()

	// Создаем мок транзакции (будет вызываться из мокируемой базы данных)
	tx, err := db.Beginx()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	// Мокаем ожидание вызова Rollback
	mock.ExpectRollback()

	// Создаем HTTP-респондер для захвата ответа
	w := httptest.NewRecorder()

	// Вызываем handleTransactionEnd с кодом ошибки
	handleTransactionEnd(tx, http.StatusBadRequest, w)

	// Проверка, что транзакция откатилась (Rollback)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleTransactionEnd_Commit(t *testing.T) {
	// Создаем мок базы данных и мок
	db, mock := createTestDB()

	// Ожидаем вызов Begin
	mock.ExpectBegin()

	// Создаем мок транзакции (будет вызываться из мокируемой базы данных)
	tx, err := db.Beginx()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	// Мокаем ожидание вызова Commit
	mock.ExpectCommit()

	// Создаем HTTP-респондер для захвата ответа
	w := httptest.NewRecorder()

	// Вызываем handleTransactionEnd с успешным кодом статуса
	handleTransactionEnd(tx, http.StatusOK, w)

	// Проверка, что транзакция была зафиксирована (Commit)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHandleCommit_Error(t *testing.T) {
	db, mock := createTestDB()
	mock.ExpectBegin()
	tx, err := db.Beginx()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}
	commitErr := assert.AnError
	mock.ExpectCommit().WillReturnError(commitErr)
	w := httptest.NewRecorder()
	handleCommit(tx, w)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBeginTx(t *testing.T) {
	db, mock := createTestDB()
	defer db.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)
	rw := &responseWriter{
		ResponseWriter: httptest.NewRecorder(),
	}
	mock.ExpectBegin()
	tx, err := beginTx(db, rw, req)
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestBeginTx_Error(t *testing.T) {
	db, mock := createTestDB()
	defer db.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)
	req, err := http.NewRequest("GET", "http://example.com", nil)
	require.NoError(t, err)
	rw := httptest.NewRecorder()
	tx, err := beginTx(db, rw, req)
	assert.Nil(t, tx)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, http.StatusInternalServerError, rw.Code)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestTxMiddleware_BeginError(t *testing.T) {
	// Создаем мок базы данных и мок
	db, mock := createTestDB()

	// Мокаем ошибку при вызове Begin
	mock.ExpectBegin().WillReturnError(assert.AnError)

	// Создаем HTTP-респондер для захвата ответа
	w := httptest.NewRecorder()

	// Создаем тестовый запрос
	r, _ := http.NewRequest("GET", "/test", nil)

	// Создаем миддлвару и оборачиваем ее с мокированной базой данных
	middleware := TxMiddleware(db)

	// Создаем обработчик, который просто отвечает на запрос
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("This handler should not be called due to Begin error")
	})

	// Вызываем миддлвару с запросом и обработчиком
	middleware(next).ServeHTTP(w, r)

	// Проверяем, что статус ответа остался пустым, так как транзакция не была начата
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
